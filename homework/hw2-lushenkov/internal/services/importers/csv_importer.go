package importers

import (
	"bytes"
	"encoding/csv"
	"strconv"
	"time"

	"github.com/perekoshik/oop-go-concepts/internal/domain/models"
	"github.com/perekoshik/oop-go-concepts/internal/services/repository"
)

// CSVImporter reads CSV exported by visitor.
type CSVImporter struct {
	templateImporter
}

// Import parses CSV file into storage data.
func (i *CSVImporter) Import(path string) (*repository.StorageData, error) {
	raw, err := i.readFile(path)
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(bytes.NewReader(raw))
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return &repository.StorageData{}, nil
	}
	var data repository.StorageData
	for _, row := range rows[1:] {
		if len(row) < 10 {
			continue
		}
		entity := row[0]
		switch entity {
		case "account":
			balance, _ := strconv.ParseFloat(row[4], 64)
			data.Accounts = append(data.Accounts, models.BankAccount{
				ID:      row[1],
				Name:    row[2],
				Balance: balance,
			})
		case "category":
			data.Categories = append(data.Categories, models.Category{
				ID:   row[1],
				Name: row[2],
				Type: models.CategoryType(row[3]),
			})
		case "operation":
			amount, _ := strconv.ParseFloat(row[7], 64)
			date, _ := time.Parse(time.RFC3339, row[8])
			data.Operations = append(data.Operations, models.Operation{
				ID:            row[1],
				Type:          models.OperationType(row[3]),
				BankAccountID: row[5],
				CategoryID:    row[6],
				Amount:        amount,
				Date:          date,
				Description:   row[9],
			})
		}
	}
	return &data, nil
}
