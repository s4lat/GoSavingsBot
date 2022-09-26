package components

import (
	"time"
)

type Spend struct {
	ID uint
	UserID int64
	Name string
	Value float32
	Date time.Time
}

type TimeZone struct {
	UserID int64
	TZ string
}