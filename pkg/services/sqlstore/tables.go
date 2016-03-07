package sqlstore

// extra tables not required by the core/outside model

type DashboardTag struct {
	Id          int64
	DashboardId int64
	Term        string
}

const DashboardTagTable = "tpt_dh_star"

func (a *DashboardTag) TableName() string {
	return DashboardTagTable
}
