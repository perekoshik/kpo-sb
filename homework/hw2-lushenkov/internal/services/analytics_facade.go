package services

import (
	"errors"
	"time"

	"github.com/perekoshik/oop-go-concepts/internal/domain/models"
	"github.com/perekoshik/oop-go-concepts/internal/services/repository"
)

// CategoryTotals aggregates income and expense totals for category.
type CategoryTotals struct {
	Income  float64
	Expense float64
}

// AnalyticsFacade delivers reporting and analytical scenarios.
type AnalyticsFacade struct {
	repo repository.Repository
}

// NewAnalyticsFacade creates analytics facade.
func NewAnalyticsFacade(repo repository.Repository) *AnalyticsFacade {
	return &AnalyticsFacade{repo: repo}
}

// Difference computes net result between income and expenses for period inclusively.
func (f *AnalyticsFacade) Difference(from, to time.Time) (float64, error) {
	if err := validatePeriod(from, to); err != nil {
		return 0, err
	}
	operations, err := f.repo.ListOperations()
	if err != nil {
		return 0, err
	}
	var income, expense float64
	for _, op := range operations {
		if !inPeriod(op.Date, from, to) {
			continue
		}
		switch op.Type {
		case models.OperationTypeIncome:
			income += op.Amount
		case models.OperationTypeExpense:
			expense += op.Amount
		}
	}
	return income - expense, nil
}

// GroupByCategory aggregates totals per category for provided period.
func (f *AnalyticsFacade) GroupByCategory(from, to time.Time) (map[string]CategoryTotals, error) {
	if err := validatePeriod(from, to); err != nil {
		return nil, err
	}
	operations, err := f.repo.ListOperations()
	if err != nil {
		return nil, err
	}
	categories, err := f.repo.ListCategories()
	if err != nil {
		return nil, err
	}
	categoryMap := make(map[string]models.Category, len(categories))
	for _, category := range categories {
		categoryMap[category.ID] = category
	}

	totals := make(map[string]CategoryTotals)
	for _, op := range operations {
		if !inPeriod(op.Date, from, to) {
			continue
		}
		category, ok := categoryMap[op.CategoryID]
		if !ok {
			continue
		}
		entry := totals[category.Name]
		switch op.Type {
		case models.OperationTypeIncome:
			entry.Income += op.Amount
		case models.OperationTypeExpense:
			entry.Expense += op.Amount
		}
		totals[category.Name] = entry
	}
	return totals, nil
}

// AverageDailyExpense calculates average expense per day for period.
func (f *AnalyticsFacade) AverageDailyExpense(from, to time.Time) (float64, error) {
	if err := validatePeriod(from, to); err != nil {
		return 0, err
	}
	operations, err := f.repo.ListOperations()
	if err != nil {
		return 0, err
	}
	var expense float64
	for _, op := range operations {
		if !inPeriod(op.Date, from, to) {
			continue
		}
		if op.Type == models.OperationTypeExpense {
			expense += op.Amount
		}
	}
	days := int(to.Sub(from).Hours()/24) + 1
	if days <= 0 {
		return 0, errors.New("invalid period length")
	}
	return expense / float64(days), nil
}

func validatePeriod(from, to time.Time) error {
	if to.Before(from) {
		return errors.New("period end cannot precede start")
	}
	return nil
}

func inPeriod(date, from, to time.Time) bool {
	return (date.Equal(from) || date.After(from)) && (date.Equal(to) || date.Before(to))
}
