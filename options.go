package golog

// An Option configures a Logger.
type Option interface {
	apply(*Logger)
}

// optionFunc wraps a func so it satisfies the Option interface.
type optionFunc func(*Logger)

func (f optionFunc) apply(log *Logger) {
	f(log)
}

// WithOptionFieldSize set fields size to the Logger.
func WithOptionFieldSize(size int) Option {
	return optionFunc(func(log *Logger) {
		log.fields = make([]field, 0, size)
	})
}
