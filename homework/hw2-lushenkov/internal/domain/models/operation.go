package models

import "time"

// OperationType defines income or expense operation nature.
type OperationType string

const (
	OperationTypeIncome  OperationType = "income"
	OperationTypeExpense OperationType = "expense"
)

// Operation represents a single financial action linked to account and category.
type Operation struct {
	ID            string        `json:"id"`
	Type          OperationType `json:"type"`
	BankAccountID string        `json:"bank_account_id"`
	Amount        float64       `json:"amount"`
	Date          time.Time     `json:"date"`
	Description   string        `json:"description,omitempty"`
	CategoryID    string        `json:"category_id"`
}

// Accept allows the Operation to be visited by an export visitor.
func (o *Operation) Accept(v Visitor) error {
	return v.VisitOperation(o)
}
