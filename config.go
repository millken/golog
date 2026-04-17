package golog

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

// Encoding is the encoding type.
type Encoding string

const (
	defaultLogLevel = INFO
	defaultModule   = "-"
	// JSONEncoding is the json encoding.
	JSONEncoding Encoding = "json"
	// TextEncoding is the text encoding.
	TextEncoding Encoding = "text"
)

var (
	rwmutex = &sync.RWMutex{}
	configs = newConfigs()
)

// Configs holds the default and module-specific configs.
type Configs struct {
	Default Config            `json:"default" yaml:"default"`
	Modules map[string]Config `json:"modules" yaml:"modules"`
}

// Config is a configuration for a logger.
type Config struct {
	// Level is the default log level.
	Level Level `json:"level" yaml:"level"`
	// Encoding is the log encoding.  text or json.
	Encoding    Encoding          `json:"encoding" yaml:"encoding"`
	TextEncoder TextEncoderConfig `json:"textEncoder" yaml:"textEncoder"`
	JSONEncoder JSONEncoderConfig `json:"jsonEncoder" yaml:"jsonEncoder"`
	// CallerLevels is the default levels for show caller info.
	CallerLevels []Level `json:"callerLevels" yaml:"callerLevels"`
	// StacktraceLevels is the default levels for show stacktrace.
	StacktraceLevels []Level       `json:"stacktraceLevels" yaml:"stacktraceLevels"`
	Handler          HandlerConfig `json:"handler" yaml:"handler"`
}

// TextEncoderConfig is the configuration for the text encoder.
type TextEncoderConfig struct {
	// PartsOrder is the order of the parts of the log entry.
	PartsOrder []string `json:"partsOrder" yaml:"partsOrder"`
	// TimeFormat specifies the format for timestamp in output.
	TimeFormat string `json:"timeFormat" yaml:"timeFormat"`
	// DisableTimestamp disables the timestamp in output.
	DisableTimestamp bool `json:"disableTimestamp" yaml:"disableTimestamp"`
	// DisableColor disables the color in output.
	DisableColor bool `json:"disableColor" yaml:"disableColor"`
	// CallerSkipFrame is the number of stack frames to skip when reporting the calling function.
	CallerSkipFrame int `json:"callerSkipFrame" yaml:"callerSkipFrame"`
	// ShowModuleName shows the name of the logger.
	ShowModuleName bool `json:"showModuleName" yaml:"showModuleName"`
}

// JSONEncoderConfig is the configuration for the JSONEncoder.
type JSONEncoderConfig struct {
	// TimeFormat specifies the format for timestamp in output.
	TimeFormat string `json:"timeFormat" yaml:"timeFormat"`
	// DisableTimestamp disables the timestamp in output.
	DisableTimestamp bool `json:"disableTimestamp" yaml:"disableTimestamp"`
	// CallerSkipFrame is the number of stack frames to skip when reporting the calling function.
	CallerSkipFrame int `json:"callerSkipFrame" yaml:"callerSkipFrame"`
	// ShowModuleName shows the name of the logger.
	ShowModuleName bool `json:"showModuleName" yaml:"showModuleName"`
}

// HandlerType defines the type of log handler.
type HandlerType string

const (
	// HandlerTypeFile writes logs to a file (stdout/stderr/path).
	HandlerTypeFile HandlerType = "file"
	// HandlerTypeRotateFile writes logs to a file with time-based rotation.
	HandlerTypeRotateFile HandlerType = "rotateFile"
	// HandlerTypeCustom uses a user-provided io.Writer.
	HandlerTypeCustom HandlerType = "custom"
)

// HandlerConfig is a configuration for a writer.
type HandlerConfig struct {
	Type       HandlerType      `json:"type" yaml:"type"`
	Writer     io.Writer        `json:"-" yaml:"-"`
	File       FileConfig       `json:"file" yaml:"file"`
	RotateFile RotateFileConfig `json:"rotateFile" yaml:"rotateFile"`
}

// FileConfig is a configuration for a file writer.
type FileConfig struct {
	Path string `json:"path" yaml:"path"`
}

func newConfigs() *Configs {
	return &Configs{
		Default: Config{
			Level:    defaultLogLevel,
			Encoding: TextEncoding,
		},
		Modules: make(map[string]Config),
	}
}

// ResetConfigs resets all configs to default values.
func ResetConfigs() {
	rwmutex.Lock()
	defer rwmutex.Unlock()
	loggerProviderFactoryFn.Store(newLoggerProviderFactory()) // reset logger provider
	configs = newConfigs()
}

// GetConfigs returns a deep copy of the current configs.
func GetConfigs() Configs {
	rwmutex.RLock()
	defer rwmutex.RUnlock()
	cp := *configs
	cp.Modules = make(map[string]Config, len(configs.Modules))
	for k, v := range configs.Modules {
		cp.Modules[k] = v
	}
	return cp
}

// SetLevel sets the default log level. Only affects loggers created after this call.
func SetLevel(level Level) {
	rwmutex.Lock()
	defer rwmutex.Unlock()
	configs.Default.Level = level
}

// SetEncoding sets the default log encoding. Only affects loggers created after this call.
func SetEncoding(encoding Encoding) {
	rwmutex.Lock()
	defer rwmutex.Unlock()
	configs.Default.Encoding = encoding
}

// SetTextEncoderConfig sets the text encoder config. Only affects loggers created after this call.
func SetTextEncoderConfig(cfg TextEncoderConfig) {
	rwmutex.Lock()
	defer rwmutex.Unlock()
	configs.Default.TextEncoder = cfg
}

// SetJSONEncoderConfig sets the json encoder config. Only affects loggers created after this call.
func SetJSONEncoderConfig(cfg JSONEncoderConfig) {
	rwmutex.Lock()
	defer rwmutex.Unlock()
	configs.Default.JSONEncoder = cfg
}

// SetCallerLevels sets the caller levels. Only affects loggers created after this call.
func SetCallerLevels(levels ...Level) {
	rwmutex.Lock()
	defer rwmutex.Unlock()
	configs.Default.CallerLevels = levels
}

// SetStacktraceLevels sets the stacktrace levels. Only affects loggers created after this call.
func SetStacktraceLevels(levels ...Level) {
	rwmutex.Lock()
	defer rwmutex.Unlock()
	configs.Default.StacktraceLevels = levels
}

// SetWriter sets the default writer. Only affects loggers created after this call.
func SetWriter(writer io.Writer) {
	rwmutex.Lock()
	defer rwmutex.Unlock()
	configs.Default.Handler.Type = HandlerTypeCustom
	configs.Default.Handler.Writer = writer
}

// LoadConfig loads config from a YAML or JSON file. Only affects loggers created after this call.
func LoadConfig(path string) error {
	var out Configs
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config content: %w", err)
	}
	ext := filepath.Ext(path)
	switch ext {
	case ".json":
		if err := json.Unmarshal(data, &out); err != nil {
			return fmt.Errorf("failed to unmarshal json: %w", err)
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &out); err != nil {
			return fmt.Errorf("failed to unmarshal yaml: %w", err)
		}
	default:
		return fmt.Errorf("unsupported config file extension: %s", ext)
	}

	rwmutex.Lock()
	defer rwmutex.Unlock()
	configs = &out
	return nil
}

// SetModuleConfig sets the config for the given module. Only affects loggers created after this call.
func SetModuleConfig(module string, cfg Config) {
	rwmutex.Lock()
	defer rwmutex.Unlock()
	configs.Modules[module] = cfg
}

// GetModuleConfig returns the config for the given module.
func GetModuleConfig(module string) Config {
	rwmutex.RLock()
	defer rwmutex.RUnlock()
	cfg, exists := configs.Modules[module]
	if !exists {
		cfg = configs.Default
	}

	return cfg
}
