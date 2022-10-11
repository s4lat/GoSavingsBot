package components

import (
	"gorm.io/gorm"
	"time"
)

// A Spend stores values of spend from db. Used by gorm.DB to save and recieve values.
type Spend struct {
	ID     int64
	UserID int64
	Name   string
	Value  float32
	Date   time.Time
}

// A User stores values of user from db. Used by gorm.DB to save and recieve values.
type User struct {
	ID       int64
	TimeZone string
	Lang     string
}

// Returning slice of spends, ordered from older to newer dates, for year.
func GetSpendsByYear(uid int64, db *gorm.DB, year int, loc *time.Location) []Spend {
	var spends []Spend
	fromDate := time.Date(year, time.Month(1), 1, 0, 0, 0, 0, time.Now().Location()).AddDate(0, 0, -2)
	toDate := fromDate.AddDate(+1, 0, +6)
	db.Order("date").Find(&spends, "user_id = ? AND date BETWEEN ? AND ?", uid, fromDate, toDate)

	sorted_spends := make([]Spend, 0, len(spends))
	for _, spend := range spends {
		spend.Date = spend.Date.In(loc)

		if spend.Date.Year() == year {
			sorted_spends = append(sorted_spends, spend)
		}
	}

	return sorted_spends
}

// Returning slice of spends, ordered from older to newer dates, for month.
func GetSpendsByMonthYear(uid int64, db *gorm.DB, month int, year int, loc *time.Location) []Spend {
	var spends []Spend
	fromDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Now().Location()).AddDate(0, 0, -2)
	toDate := fromDate.AddDate(0, +1, +6)
	db.Order("date").Find(&spends, "user_id = ? AND date BETWEEN ? AND ?", uid, fromDate, toDate)

	sorted_spends := make([]Spend, 0, len(spends))
	for _, spend := range spends {
		spend.Date = spend.Date.In(loc)

		if spend.Date.Month() == time.Month(month) {
			sorted_spends = append(sorted_spends, spend)
		}
	}

	return sorted_spends
}

// Returning slice of spends, ordered from older to newer dates, for day.
func GetSpendsByDayMonthYear(uid int64, db *gorm.DB, day int, month int, year int, loc *time.Location) []Spend {
	var spends []Spend
	fromDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Now().Location()).AddDate(0, 0, -2)
	toDate := fromDate.AddDate(0, 0, +6)
	db.Order("date").Find(&spends, "user_id = ? AND date BETWEEN ? AND ?", uid, fromDate, toDate)

	sorted_spends := make([]Spend, 0, len(spends))
	for _, spend := range spends {
		spend.Date = spend.Date.In(loc)

		if spend.Date.Day() == day {
			sorted_spends = append(sorted_spends, spend)
		}
	}

	return sorted_spends
}
