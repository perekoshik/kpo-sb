package models

// BankAccount represents a user's account with current balance.
type BankAccount struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
}

// Accept allows the BankAccount to be visited by an export visitor.
func (a *BankAccount) Accept(v Visitor) error {
	return v.VisitBankAccount(a)
}
