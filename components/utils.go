package components

import (
	"encoding/csv"
	"bytes"
	"fmt"
)

func SpendsToCSV(spends []Spend) *bytes.Buffer{
	records := make([][]string, len(spends) + 1)

	records[0] = []string{"name", "value", "clock", "date"}
	for i := 0; i < len(spends); i++ {
		hours, mins, _ := spends[i].Date.Clock()
		year, month, day := spends[i].Date.Date()
	    records[i + 1] = []string{
	    	spends[i].Name, 
			fmt.Sprintf("%0.2f", spends[i].Value), 
	    	fmt.Sprintf("%02d:%02d", hours, mins),
	    	fmt.Sprintf("%02d.%02d.%d", day, month, year),
	    }
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	w.WriteAll(records)


	return &buf
}