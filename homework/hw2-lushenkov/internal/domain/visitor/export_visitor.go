package visitor

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/perekoshik/oop-go-concepts/internal/domain/models"
)

// ExportFormat enumerates supported export targets.
type ExportFormat string

const (
	ExportFormatJSON ExportFormat = "json"
	ExportFormatYAML ExportFormat = "yaml"
	ExportFormatCSV  ExportFormat = "csv"
)

// FileExportVisitor aggregates domain objects and renders them into chosen format.
type FileExportVisitor struct {
	format     ExportFormat
	accounts   []*models.BankAccount
	categories []*models.Category
	operations []*models.Operation
}

// NewFileExportVisitor prepares visitor to collect entities for specific format.
func NewFileExportVisitor(format ExportFormat) (*FileExportVisitor, error) {
	switch format {
	case ExportFormatJSON, ExportFormatYAML, ExportFormatCSV:
		return &FileExportVisitor{format: format}, nil
	default:
		return nil, fmt.Errorf("unsupported export format: %s", format)
	}
}

// VisitBankAccount collects bank account entities.
func (v *FileExportVisitor) VisitBankAccount(account *models.BankAccount) error {
	v.accounts = append(v.accounts, account)
	return nil
}

// VisitCategory collects categories.
func (v *FileExportVisitor) VisitCategory(category *models.Category) error {
	v.categories = append(v.categories, category)
	return nil
}

// VisitOperation collects operations.
func (v *FileExportVisitor) VisitOperation(operation *models.Operation) error {
	v.operations = append(v.operations, operation)
	return nil
}

// Render marshals collected entities according to configured format.
func (v *FileExportVisitor) Render() ([]byte, error) {
	switch v.format {
	case ExportFormatJSON:
		payload := struct {
			Accounts   []*models.BankAccount `json:"accounts"`
			Categories []*models.Category    `json:"categories"`
			Operations []*models.Operation   `json:"operations"`
		}{v.accounts, v.categories, v.operations}
		return json.MarshalIndent(payload, "", "  ")
	case ExportFormatYAML:
		return v.renderYAML()
	case ExportFormatCSV:
		return v.renderCSV()
	default:
		return nil, fmt.Errorf("format %s not implemented", v.format)
	}
}

func (v *FileExportVisitor) renderCSV() ([]byte, error) {
	buf := &bytes.Buffer{}
	writer := csv.NewWriter(buf)
	header := []string{"entity", "id", "name", "type", "balance", "bank_account_id", "category_id", "amount", "date", "description"}
	if err := writer.Write(header); err != nil {
		return nil, err
	}
	for _, account := range v.accounts {
		err := writer.Write([]string{
			"account",
			account.ID,
			account.Name,
			"",
			strconv.FormatFloat(account.Balance, 'f', 2, 64),
			"",
			"",
			"",
			"",
			"",
		})
		if err != nil {
			return nil, err
		}
	}
	for _, category := range v.categories {
		err := writer.Write([]string{
			"category",
			category.ID,
			category.Name,
			string(category.Type),
			"",
			"",
			"",
			"",
			"",
			"",
		})
		if err != nil {
			return nil, err
		}
	}
	for _, operation := range v.operations {
		err := writer.Write([]string{
			"operation",
			operation.ID,
			"",
			string(operation.Type),
			"",
			operation.BankAccountID,
			operation.CategoryID,
			strconv.FormatFloat(operation.Amount, 'f', 2, 64),
			operation.Date.Format(time.RFC3339),
			operation.Description,
		})
		if err != nil {
			return nil, err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (v *FileExportVisitor) renderYAML() ([]byte, error) {
	buf := &bytes.Buffer{}
	buf.WriteString("accounts:\n")
	for _, account := range v.accounts {
		fmt.Fprintf(buf, "  - id: %s\n", account.ID)
		fmt.Fprintf(buf, "    name: %q\n", account.Name)
		fmt.Fprintf(buf, "    balance: %.2f\n", account.Balance)
	}

	buf.WriteString("categories:\n")
	for _, category := range v.categories {
		fmt.Fprintf(buf, "  - id: %s\n", category.ID)
		fmt.Fprintf(buf, "    name: %q\n", category.Name)
		fmt.Fprintf(buf, "    type: %s\n", category.Type)
	}

	buf.WriteString("operations:\n")
	for _, operation := range v.operations {
		fmt.Fprintf(buf, "  - id: %s\n", operation.ID)
		fmt.Fprintf(buf, "    type: %s\n", operation.Type)
		fmt.Fprintf(buf, "    bank_account_id: %s\n", operation.BankAccountID)
		fmt.Fprintf(buf, "    category_id: %s\n", operation.CategoryID)
		fmt.Fprintf(buf, "    amount: %.2f\n", operation.Amount)
		fmt.Fprintf(buf, "    date: %s\n", operation.Date.Format(time.RFC3339))
		if operation.Description != "" {
			fmt.Fprintf(buf, "    description: %q\n", operation.Description)
		}
	}

	return buf.Bytes(), nil
}
