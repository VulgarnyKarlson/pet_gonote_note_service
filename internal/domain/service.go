package domain

import "time"

type SearchCriteria struct {
	Title    string
	Content  string
	FromDate time.Time
	ToDate   time.Time
}
