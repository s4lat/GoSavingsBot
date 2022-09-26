package components

import (
	"time"
	"gorm.io/gorm"
	"log"
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

func GetSpendsByMonthYear(uid int64, db *gorm.DB, month int, year int, loc *time.Location) []Spend{
	var spends []Spend
	fromDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Now().Location()).AddDate(0, 0, -2)
	toDate := fromDate.AddDate(0, +1, +4)
	db.Order("date").Find(&spends, "user_id = ? AND date BETWEEN ? AND ?", uid, fromDate, toDate)

	sorted_spends := make([]Spend, 0, len(spends))
	for _, spend := range spends {
		t := spend.Date.In(loc)

		if t.Month() == time.Month(month) {
			sorted_spends = append(sorted_spends, spend)
		}
	}

	return sorted_spends
}

func GetSpendsByDayMonthYear(uid int64, db *gorm.DB, day int, month int, year int, loc *time.Location) []Spend{
	var spends []Spend
	fromDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Now().Location()).AddDate(0, 0, -2)
	toDate := fromDate.AddDate(0, 0, +4)
	db.Order("date").Find(&spends, "user_id = ? AND date BETWEEN ? AND ?", uid, fromDate, toDate)

	log.Print(len(spends))

	sorted_spends := make([]Spend, 0, len(spends))
	for _, spend := range spends {
		t := spend.Date.In(loc)
		if t.Day() == day {
			sorted_spends = append(sorted_spends, spend)
		}
	}

	return sorted_spends
}
