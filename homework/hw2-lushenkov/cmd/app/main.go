package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/perekoshik/oop-go-concepts/internal/domain/models"
	"github.com/perekoshik/oop-go-concepts/internal/domain/visitor"
	"github.com/perekoshik/oop-go-concepts/internal/services/command"
	"github.com/perekoshik/oop-go-concepts/internal/services/importers"
	"github.com/perekoshik/oop-go-concepts/internal/utils"
)

func main() {
	container, err := utils.NewContainer(".")
	if err != nil {
		log.Fatalf("init container: %v", err)
	}

	reader := bufio.NewReader(os.Stdin)
	logger := func(format string, args ...interface{}) {
		fmt.Printf("[timing] "+format+"\n", args...)
	}

	for {
		fmt.Println("\nHSE Finance Tracker")
		fmt.Println("1. Create bank account")
		fmt.Println("2. List bank accounts")
		fmt.Println("3. Create category")
		fmt.Println("4. List categories")
		fmt.Println("5. Add operation")
		fmt.Println("6. Edit operation")
		fmt.Println("7. Delete operation")
		fmt.Println("8. List operations")
		fmt.Println("9. Analytics: income vs expense difference")
		fmt.Println("10. Analytics: totals by category")
		fmt.Println("11. Analytics: average daily expense")
		fmt.Println("12. Export data")
		fmt.Println("13. Import data")
		fmt.Println("14. Recalculate account balance")
		fmt.Println("0. Exit")
		fmt.Print("Choose option: ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			name := prompt(reader, "Account name: ")
			balance := promptFloat(reader, "Initial balance: ")
			account, err := container.AccountFacade.CreateAccount(name, balance)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			fmt.Println("Account created with id", account.ID)
		case "2":
			accounts, err := container.AccountFacade.ListAccounts()
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			if len(accounts) == 0 {
				fmt.Println("No accounts yet")
				continue
			}
			for _, acc := range accounts {
				fmt.Printf("- %s (%s): balance %.2f\n", acc.Name, acc.ID, acc.Balance)
			}
		case "3":
			name := prompt(reader, "Category name: ")
			catTypeStr := strings.ToLower(prompt(reader, "Type (income/expense): "))
			var catType models.CategoryType
			if catTypeStr == "income" {
				catType = models.CategoryTypeIncome
			} else {
				catType = models.CategoryTypeExpense
			}
			category, err := container.CategoryFacade.CreateCategory(name, catType)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			fmt.Println("Category created with id", category.ID)
		case "4":
			categories, err := container.CategoryFacade.ListCategories()
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			if len(categories) == 0 {
				fmt.Println("No categories yet")
				continue
			}
			for _, cat := range categories {
				fmt.Printf("- %s (%s) [%s]\n", cat.Name, cat.ID, cat.Type)
			}
		case "5":
			typeStr := strings.ToLower(prompt(reader, "Operation type (income/expense): "))
			var opType models.OperationType
			if typeStr == "income" {
				opType = models.OperationTypeIncome
			} else {
				opType = models.OperationTypeExpense
			}
			accountID := prompt(reader, "Account id: ")
			categoryID := prompt(reader, "Category id: ")
			amount := promptFloat(reader, "Amount: ")
			date := promptDate(reader, "Date (YYYY-MM-DD): ")
			description := prompt(reader, "Description (optional): ")
			cmd := &command.AddOperationCommand{
				Facade:      container.OperationFacade,
				Type:        opType,
				AccountID:   accountID,
				CategoryID:  categoryID,
				Amount:      amount,
				Date:        date,
				Description: description,
				OnCreated: func(op *models.Operation) {
					fmt.Println("Operation created with id", op.ID)
				},
			}
			timed := command.NewTimedCommand("Add operation", cmd, logger)
			if err := timed.Execute(); err != nil {
				fmt.Println("Error:", err)
			}
		case "6":
			id := prompt(reader, "Operation id: ")
			typeStr := strings.ToLower(prompt(reader, "New type (income/expense): "))
			var opType models.OperationType
			if typeStr == "income" {
				opType = models.OperationTypeIncome
			} else {
				opType = models.OperationTypeExpense
			}
			accountID := prompt(reader, "New account id: ")
			categoryID := prompt(reader, "New category id: ")
			amount := promptFloat(reader, "New amount: ")
			date := promptDate(reader, "New date (YYYY-MM-DD): ")
			description := prompt(reader, "New description: ")
			cmd := &command.EditOperationCommand{
				Facade:      container.OperationFacade,
				ID:          id,
				Type:        opType,
				AccountID:   accountID,
				CategoryID:  categoryID,
				Amount:      amount,
				Date:        date,
				Description: description,
			}
			timed := command.NewTimedCommand("Edit operation", cmd, logger)
			if err := timed.Execute(); err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Operation updated")
			}
		case "7":
			id := prompt(reader, "Operation id: ")
			cmd := &command.DeleteOperationCommand{
				Facade: container.OperationFacade,
				ID:     id,
			}
			timed := command.NewTimedCommand("Delete operation", cmd, logger)
			if err := timed.Execute(); err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Operation deleted")
			}
		case "8":
			operations, err := container.OperationFacade.ListOperations()
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			if len(operations) == 0 {
				fmt.Println("No operations yet")
				continue
			}
			for _, op := range operations {
				fmt.Printf("- %s (%s) %s %.2f on %s [account=%s category=%s]\n", op.ID, op.Type, op.Description, op.Amount, op.Date.Format(time.RFC3339), op.BankAccountID, op.CategoryID)
			}
		case "9":
			from := promptDate(reader, "From (YYYY-MM-DD): ")
			to := promptDate(reader, "To (YYYY-MM-DD): ")
			diff, err := container.AnalyticsFacade.Difference(from, to)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			fmt.Printf("Net result: %.2f\n", diff)
		case "10":
			from := promptDate(reader, "From (YYYY-MM-DD): ")
			to := promptDate(reader, "To (YYYY-MM-DD): ")
			totals, err := container.AnalyticsFacade.GroupByCategory(from, to)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			if len(totals) == 0 {
				fmt.Println("No data in period")
				continue
			}
			for category, summary := range totals {
				fmt.Printf("- %s: income=%.2f expense=%.2f\n", category, summary.Income, summary.Expense)
			}
		case "11":
			from := promptDate(reader, "From (YYYY-MM-DD): ")
			to := promptDate(reader, "To (YYYY-MM-DD): ")
			avg, err := container.AnalyticsFacade.AverageDailyExpense(from, to)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			fmt.Printf("Average daily expense: %.2f\n", avg)
		case "12":
			formatStr := strings.ToLower(prompt(reader, "Format (json/yaml/csv): "))
			var format visitor.ExportFormat
			switch formatStr {
			case "json":
				format = visitor.ExportFormatJSON
			case "yaml":
				format = visitor.ExportFormatYAML
			case "csv":
				format = visitor.ExportFormatCSV
			default:
				fmt.Println("Unsupported format")
				continue
			}
			path := prompt(reader, "Target file path (relative to project root): ")
			if path == "" {
				path = filepath.Join("data", fmt.Sprintf("export.%s", formatStr))
			}
			if err := container.DataManager.ExportData(format, path); err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Data exported to", path)
			}
		case "13":
			formatStr := strings.ToLower(prompt(reader, "Source format (json/yaml/csv): "))
			var importer importers.Importer
			switch formatStr {
			case "json":
				importer = &importers.JSONImporter{}
			case "yaml":
				importer = &importers.YAMLImporter{}
			case "csv":
				importer = &importers.CSVImporter{}
			default:
				fmt.Println("Unsupported importer")
				continue
			}
			path := prompt(reader, "File path: ")
			if err := container.DataManager.ImportData(importer, path); err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Data imported from", path)
			}
		case "14":
			accountID := prompt(reader, "Account id: ")
			if err := container.AccountFacade.RecalculateBalance(accountID); err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Balance recalculated")
			}
		case "0":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Unknown option")
		}
	}
}

func prompt(reader *bufio.Reader, label string) string {
	fmt.Print(label)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func promptFloat(reader *bufio.Reader, label string) float64 {
	for {
		value := prompt(reader, label)
		if value == "" {
			return 0
		}
		number, err := strconv.ParseFloat(value, 64)
		if err == nil {
			return number
		}
		fmt.Println("Please enter numeric value")
	}
}

func promptDate(reader *bufio.Reader, label string) time.Time {
	for {
		value := prompt(reader, label)
		if value == "" {
			return time.Now()
		}
		parsed, err := time.Parse("2006-01-02", value)
		if err == nil {
			return parsed
		}
		fmt.Println("Please use format YYYY-MM-DD")
	}
}
