package models

// Visitor describes visitor capable of handling each domain entity.
type Visitor interface {
	VisitBankAccount(*BankAccount) error
	VisitCategory(*Category) error
	VisitOperation(*Operation) error
}
