package models

import (
	"time"

	"gorm.io/gorm"
)

// Spend stores values of spend from db. Used by gorm.DB to save and receive values.
type Spend struct {
	Date   time.Time
	Name   string
	ID     int64
	UserID int64
	Value  float32
}

// User stores values of user from db. Used by gorm.DB to save and receive values.
type User struct {
	TimeZone string
	Lang     string
	ID       int64
}

// GetSpendsByYear - returning slice of year spends ordered from older to newer dates.
func GetSpendsByYear(uid int64, db *gorm.DB, year int, loc *time.Location) []Spend {
	var spends []Spend
	fromDate := time.Date(year, time.Month(1), 1, 0, 0, 0, 0, time.Now().Location()).AddDate(0, 0, -2)
	toDate := fromDate.AddDate(+1, 0, +6)
	db.Order("date").Find(&spends, "user_id = ? AND date BETWEEN ? AND ?", uid, fromDate, toDate)

	sortedSpends := make([]Spend, 0, len(spends))
	for _, spend := range spends {
		spend.Date = spend.Date.In(loc)

		if spend.Date.Year() == year {
			sortedSpends = append(sortedSpends, spend)
		}
	}

	return sortedSpends
}

// GetSpendsByMonthYear - returning slice of month spends ordered from older to newer dates.
func GetSpendsByMonthYear(uid int64, db *gorm.DB, month int, year int, loc *time.Location) []Spend {
	var spends []Spend
	fromDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Now().Location()).AddDate(0, 0, -2)
	toDate := fromDate.AddDate(0, +1, +6)
	db.Order("date").Find(&spends, "user_id = ? AND date BETWEEN ? AND ?", uid, fromDate, toDate)

	sortedSpends := make([]Spend, 0, len(spends))
	for _, spend := range spends {
		spend.Date = spend.Date.In(loc)

		if spend.Date.Month() == time.Month(month) {
			sortedSpends = append(sortedSpends, spend)
		}
	}

	return sortedSpends
}

// GetSpendsByDayMonthYear - returning slice of day spends ordered from older to newer dates.
func GetSpendsByDayMonthYear(uid int64, db *gorm.DB, day int, month int, year int, loc *time.Location) []Spend {
	var spends []Spend
	fromDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Now().Location()).AddDate(0, 0, -2)
	toDate := fromDate.AddDate(0, 0, +6)
	db.Order("date").Find(&spends, "user_id = ? AND date BETWEEN ? AND ?", uid, fromDate, toDate)

	sortedSpends := make([]Spend, 0, len(spends))
	for _, spend := range spends {
		spend.Date = spend.Date.In(loc)

		if spend.Date.Day() == day {
			sortedSpends = append(sortedSpends, spend)
		}
	}

	return sortedSpends
}
