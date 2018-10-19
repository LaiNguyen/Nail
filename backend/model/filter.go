package model

import (
	"strings"
)

type FilterInfo struct {
	Status  int    `json:"status"`
	Sort    string `json:"sort"`
	Filter  string `json:"filter"`
	Page    int    `json:"page"`
	PerPage int    `json:"per_page"`
	Total   int    `json:"total"`
}

func (filterInfo *FilterInfo) PrepareSort() string {
	sort := ""
	if filterInfo.Sort != "" {
		sortInfo := strings.Split(filterInfo.Sort, "|")
		sort = sortInfo[0]
		if sortInfo[1] != "asc" {
			sort = "-" + sort
		}
	}
	return sort
}
