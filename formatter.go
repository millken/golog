package golog

// Formatter is the interface that must be implemented
type Formatter interface {
	Format(*Entry) error
}

// NullFormatter is a formatter that does nothing.
type NullFormatter struct{}

// Format does nothing.
func (f NullFormatter) Format(entry *Entry) error {
	return nil
}
