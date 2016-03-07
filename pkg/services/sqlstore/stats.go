package sqlstore

import (
	"github.com/grafana/grafana/pkg/bus"
	m "github.com/grafana/grafana/pkg/models"
)

func init() {
	bus.AddHandler("sql", GetSystemStats)
	bus.AddHandler("sql", GetDataSourceStats)
	bus.AddHandler("sql", GetAdminStats)
}

func GetDataSourceStats(query *m.GetDataSourceStatsQuery) error {
	var rawSql = `SELECT COUNT(*) as count, type FROM ` + m.DataSourceTable + `data_source GROUP BY type`
	query.Result = make([]*m.DataSourceStats, 0)
	err := x.Sql(rawSql).Find(&query.Result)
	if err != nil {
		return err
	}

	return err
}

func GetSystemStats(query *m.GetSystemStatsQuery) error {
	var rawSql = `SELECT
			(
				SELECT COUNT(*)
        FROM ` + dialect.Quote(m.UserTable) + `
      ) AS user_count,
			(
				SELECT COUNT(*)
        FROM ` + dialect.Quote(m.OrgTable) + `
      ) AS org_count,
      (
        SELECT COUNT(*)
        FROM ` + dialect.Quote(m.DashboardTable) + `
      ) AS dashboard_count,
      (
        SELECT COUNT(*)
        FROM ` + dialect.Quote(m.PlaylistTable) + `
      ) AS playlist_count
			`

	var stats m.SystemStats
	_, err := x.Sql(rawSql).Get(&stats)
	if err != nil {
		return err
	}

	query.Result = &stats
	return err
}

func GetAdminStats(query *m.GetAdminStatsQuery) error {
	var rawSql = `SELECT
      (
        SELECT COUNT(*)
        FROM ` + dialect.Quote(m.UserTable) + `
      ) AS user_count,
      (
        SELECT COUNT(*)
        FROM ` + dialect.Quote(m.OrgTable) + `
      ) AS org_count,
      (
        SELECT COUNT(*)
        FROM ` + dialect.Quote(m.DashboardTable) + `
      ) AS dashboard_count,
      (
        SELECT COUNT(*)
        FROM ` + dialect.Quote(m.DashboardSnapshotTable) + `
      ) AS db_snapshot_count,
      (
        SELECT COUNT( DISTINCT ( ` + dialect.Quote("term") + ` ))
        FROM ` + dialect.Quote(DashboardTagTable) + `
      ) AS db_tag_count,
      (
        SELECT COUNT(*)
        FROM ` + dialect.Quote(m.DataSourceTable) + `
      ) AS data_source_count,
      (
        SELECT COUNT(*)
        FROM ` + dialect.Quote(m.PlaylistTable) + `
      ) AS playlist_count,
      (
        SELECT COUNT(DISTINCT ` + dialect.Quote("dashboard_id") + ` )
        FROM ` + dialect.Quote(m.StarTable) + `
      ) AS starred_db_count,
      (
        SELECT COUNT(*)
        FROM ` + dialect.Quote(m.UserTable) + `
        WHERE ` + dialect.Quote("is_admin") + ` = 1
      ) AS grafana_admin_count
      `

	var stats m.AdminStats
	_, err := x.Sql(rawSql).Get(&stats)
	if err != nil {
		return err
	}

	query.Result = &stats
	return err
}
