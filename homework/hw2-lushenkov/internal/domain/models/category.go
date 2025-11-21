package models

// CategoryType describes whether category is for income or expense.
type CategoryType string

const (
	CategoryTypeIncome  CategoryType = "income"
	CategoryTypeExpense CategoryType = "expense"
)

// Category groups operations by meaning and type.
type Category struct {
	ID   string       `json:"id"`
	Type CategoryType `json:"type"`
	Name string       `json:"name"`
}

// Accept allows the Category to be visited by an export visitor.
func (c *Category) Accept(v Visitor) error {
	return v.VisitCategory(c)
}
