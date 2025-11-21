package command

// Command encapsulates a user scenario executable by console.
type Command interface {
	Execute() error
}
