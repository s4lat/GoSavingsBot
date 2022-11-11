package export

import (
	"bytes"
	"encoding/csv"
	"fmt"

	"github.com/s4lat/gosavingsbot/models"
)

// SpendsToCSV - converting spends slice to buffer with content represented in CSV format.
func SpendsToCSV(spends []models.Spend) (*bytes.Buffer, error) {
	records := make([][]string, len(spends)+1)

	records[0] = []string{"name", "value", "clock", "date"}
	for i := 0; i < len(spends); i++ {
		hours, mins, _ := spends[i].Date.Clock()
		year, month, day := spends[i].Date.Date()
		records[i+1] = []string{
			spends[i].Name,
			fmt.Sprintf("%0.2f", spends[i].Value),
			fmt.Sprintf("%02d:%02d", hours, mins),
			fmt.Sprintf("%02d.%02d.%d", day, month, year),
		}
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	err := w.WriteAll(records)

	return &buf, err
}
