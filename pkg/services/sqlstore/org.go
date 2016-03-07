package sqlstore

import (
	"time"

	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/events"
	m "github.com/grafana/grafana/pkg/models"
)

func init() {
	bus.AddHandler("sql", GetOrgById)
	bus.AddHandler("sql", CreateOrg)
	bus.AddHandler("sql", UpdateOrg)
	bus.AddHandler("sql", UpdateOrgAddress)
	bus.AddHandler("sql", GetOrgByName)
	bus.AddHandler("sql", SearchOrgs)
	bus.AddHandler("sql", DeleteOrg)
}

func SearchOrgs(query *m.SearchOrgsQuery) error {
	query.Result = make([]*m.OrgDTO, 0)
	sess := x.Table(m.OrgTable)
	if query.Query != "" {
		sess.Where("name LIKE ?", query.Query+"%")
	}
	if query.Name != "" {
		sess.Where("name=?", query.Name)
	}
	sess.Limit(query.Limit, query.Limit*query.Page)
	sess.Cols("id", "name")
	err := sess.Find(&query.Result)
	return err
}

func GetOrgById(query *m.GetOrgByIdQuery) error {
	var org m.Org
	exists, err := x.Id(query.Id).Get(&org)
	if err != nil {
		return err
	}

	if !exists {
		return m.ErrOrgNotFound
	}

	query.Result = &org
	return nil
}

func GetOrgByName(query *m.GetOrgByNameQuery) error {
	var org m.Org
	exists, err := x.Where("name=?", query.Name).Get(&org)
	if err != nil {
		return err
	}

	if !exists {
		return m.ErrOrgNotFound
	}

	query.Result = &org
	return nil
}

func isOrgNameTaken(name string, existingId int64, sess *session) (bool, error) {
	// check if org name is taken
	var org m.Org
	exists, err := sess.Where("name=?", name).Get(&org)

	if err != nil {
		return false, nil
	}

	if exists && existingId != org.Id {
		return true, nil
	}

	return false, nil
}

func CreateOrg(cmd *m.CreateOrgCommand) error {
	return inTransaction2(func(sess *session) error {

		if isNameTaken, err := isOrgNameTaken(cmd.Name, 0, sess); err != nil {
			return err
		} else if isNameTaken {
			return m.ErrOrgNameTaken
		}

		org := m.Org{
			Name:    cmd.Name,
			Created: time.Now(),
			Updated: time.Now(),
		}

		if _, err := sess.Insert(&org); err != nil {
			return err
		}

		user := m.OrgUser{
			OrgId:   org.Id,
			UserId:  cmd.UserId,
			Role:    m.ROLE_ADMIN,
			Created: time.Now(),
			Updated: time.Now(),
		}

		_, err := sess.Insert(&user)
		cmd.Result = org

		sess.publishAfterCommit(&events.OrgCreated{
			Timestamp: org.Created,
			Id:        org.Id,
			Name:      org.Name,
		})

		return err
	})
}

func UpdateOrg(cmd *m.UpdateOrgCommand) error {
	return inTransaction2(func(sess *session) error {

		if isNameTaken, err := isOrgNameTaken(cmd.Name, cmd.OrgId, sess); err != nil {
			return err
		} else if isNameTaken {
			return m.ErrOrgNameTaken
		}

		org := m.Org{
			Name:    cmd.Name,
			Updated: time.Now(),
		}

		if _, err := sess.Id(cmd.OrgId).Update(&org); err != nil {
			return err
		}

		sess.publishAfterCommit(&events.OrgUpdated{
			Timestamp: org.Updated,
			Id:        org.Id,
			Name:      org.Name,
		})

		return nil
	})
}

func UpdateOrgAddress(cmd *m.UpdateOrgAddressCommand) error {
	return inTransaction2(func(sess *session) error {
		org := m.Org{
			Address1: cmd.Address1,
			Address2: cmd.Address2,
			City:     cmd.City,
			ZipCode:  cmd.ZipCode,
			State:    cmd.State,
			Country:  cmd.Country,

			Updated: time.Now(),
		}

		if _, err := sess.Id(cmd.OrgId).Update(&org); err != nil {
			return err
		}

		sess.publishAfterCommit(&events.OrgUpdated{
			Timestamp: org.Updated,
			Id:        org.Id,
			Name:      org.Name,
		})

		return nil
	})
}

func DeleteOrg(cmd *m.DeleteOrgCommand) error {
	return inTransaction2(func(sess *session) error {

		deletes := []string{
			"DELETE FROM " + m.StarTable + " WHERE EXISTS (SELECT 1 FROM " + m.DashboardTable + " WHERE org_id = ? AND " + m.StarTable + ".dashboard_id = " + m.DashboardTable + ".id)",
			"DELETE FROM " + DashboardTagTable + " WHERE EXISTS (SELECT 1 FROM " + m.DashboardTable + " WHERE org_id = ? AND " + DashboardTagTable + ".dashboard_id = " + m.DashboardTable + ".id)",
			"DELETE FROM " + m.DashboardTable + " WHERE org_id = ?",
			"DELETE FROM " + m.ApiKeyTable + " WHERE org_id = ?",
			"DELETE FROM " + m.DataSourceTable + " WHERE org_id = ?",
			"DELETE FROM " + m.OrgUserTable + " WHERE org_id = ?",
			"DELETE FROM " + m.OrgTable + " WHERE id = ?",
			"DELETE FROM " + m.TempUserTable + " WHERE org_id = ?",
		}

		for _, sql := range deletes {
			_, err := sess.Exec(sql, cmd.Id)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
