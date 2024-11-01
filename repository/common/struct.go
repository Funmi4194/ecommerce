package common

import (
	"github.com/uptrace/bun"
)

type Pagination struct {
	Page  int      `json:"page"`
	Limit int      `json:"limit"`
	Total int      `json:"total"`
	Pages int      `json:"pages"`
	Tags  []string `json:"tags"`
}

type History struct {
	Act string       `json:"act"`
	By  string       `json:"by"`
	At  bun.NullTime `json:"at"`
}

type Histories []History
