package models

type SearchHit struct {
	Id        int64    `json:"id"`
	Title     string   `json:"title"`
	Uri       string   `json:"uri"`
	Type      string   `json:"type"`
	Tags      []string `json:"tags"`
	IsStarred bool     `json:"isStarred"`
}

const SearchHitTable = "tpt_dh_search_hit"

func (a *SearchHit) TableName() string {
	return SearchHitTable
}
