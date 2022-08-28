package models

import (
	"github.com/lib/pq"
)

type Url struct {
	Id        int64
	Title     string
	CreatedAt pq.NullTime
	UpdatedAt pq.NullTime
}

var table = "short_urls"

func GetTable() string {
	return table
}
