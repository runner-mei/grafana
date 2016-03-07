package models

type SystemStats struct {
	DashboardCount int
	UserCount      int
	OrgCount       int
	PlaylistCount  int
}

const SystemStatsTable = "tpt_dh_system_stats"

func (a *SystemStats) TableName() string {
	return SystemStatsTable
}

type DataSourceStats struct {
	Count int
	Type  string
}

const DataSourceStatsTable = "tpt_dh_data_source_stats"

func (a *DataSourceStats) TableName() string {
	return DataSourceStatsTable
}

type GetSystemStatsQuery struct {
	Result *SystemStats
}

type GetDataSourceStatsQuery struct {
	Result []*DataSourceStats
}

type AdminStats struct {
	UserCount         int `json:"user_count"`
	OrgCount          int `json:"org_count"`
	DashboardCount    int `json:"dashboard_count"`
	DbSnapshotCount   int `json:"db_snapshot_count"`
	DbTagCount        int `json:"db_tag_count"`
	DataSourceCount   int `json:"data_source_count"`
	PlaylistCount     int `json:"playlist_count"`
	StarredDbCount    int `json:"starred_db_count"`
	GrafanaAdminCount int `json:"grafana_admin_count"`
}

type GetAdminStatsQuery struct {
	Result *AdminStats
}
