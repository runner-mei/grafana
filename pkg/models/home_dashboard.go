package models

import "time"

type HomeDashboard struct {
	Id        int64
	UserId    int64
	AccountId int64

	Created time.Time
	Updated time.Time

	Data map[string]interface{}
}

const HomeDashboardTable = "tpt_dh_home_dashboard"

func (a *HomeDashboard) TableName() string {
	return HomeDashboardTable
}
