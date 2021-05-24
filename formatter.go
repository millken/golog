package golog

type Formatter interface {
	Format(*Entry) error
}

type NullFormatter struct{}

func (f NullFormatter) Format(entry *Entry) error {
	return nil
}
