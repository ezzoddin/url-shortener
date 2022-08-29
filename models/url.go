package models

import (
	"time"
)

type Url struct {
	Id        int64
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

var table = "short_urls"

func GetTable() string {
	return table
}
