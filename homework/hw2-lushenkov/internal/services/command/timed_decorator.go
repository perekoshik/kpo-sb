package command

import (
	"time"
)

// LogFunc describes logging sink used by decorator.
type LogFunc func(format string, args ...interface{})

// TimedCommand decorates command with execution time measurement.
type TimedCommand struct {
	name   string
	inner  Command
	logger LogFunc
}

// NewTimedCommand wraps inner command with logging.
func NewTimedCommand(name string, inner Command, logger LogFunc) *TimedCommand {
	return &TimedCommand{name: name, inner: inner, logger: logger}
}

// Execute runs command and reports duration.
func (c *TimedCommand) Execute() error {
	start := time.Now()
	err := c.inner.Execute()
	if c.logger != nil {
		c.logger("%s completed in %v", c.name, time.Since(start))
	}
	return err
}
