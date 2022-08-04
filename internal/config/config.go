package config

// Config is a configuration for a logger.
type Config struct {
	// Level is the default log level.
	Level string
	// Encoding is the log encoding.  console or json.
	Encoding string
	//CallerLevels is the default levels for show caller info.
	CallerLevels []string
	// StacktraceLevels is the default levels for show stacktrace.
	StacktraceLevels []string
	// Format is the default log format.
	Format string
	// OutputPath is the default log output path.
	OutputPath string
	// ErrorOutputPaths is the default log error output paths.
	ErrorOutputPaths []string
	// Modules is the default log modules.
}

type ProviderConfig struct {
}

type Encoder struct {
}
