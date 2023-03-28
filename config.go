package golog

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// Encoding is the encoding type.
type Encoding string

const (
	defaultLogLevel = INFO
	defaultModule   = "-"
	//JSONEncoding is the json encoding.
	JSONEncoding Encoding = "json"
	//TextEncoding is the text encoding.
	TextEncoding Encoding = "text"
)

var (
	rwmutex = &sync.RWMutex{}
	configs = newConfigs()
)

// Configs - configs for golog.
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
	//CallerLevels is the default levels for show caller info.
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

// HandlerConfig is a configuration for a writer.
type HandlerConfig struct {
	Type       string           `json:"type" yaml:"type"`
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

// ResetConfigs - reset configs.
func ResetConfigs() {
	rwmutex.Lock()
	defer rwmutex.Unlock()
	loggerProviderOnce = sync.Once{} // reset logger provider
	configs = newConfigs()
}

func GetConfigs() *Configs {
	rwmutex.RLock()
	defer rwmutex.RUnlock()
	return configs
}

// SetLevel - set log level.
func SetLevel(level Level) {
	rwmutex.Lock()
	defer rwmutex.Unlock()
	configs.Default.Level = level
}

// SetEncoding - set log encoding.
func SetEncoding(encoding Encoding) {
	rwmutex.Lock()
	defer rwmutex.Unlock()
	configs.Default.Encoding = encoding
}

// SetTextEncoderConfig - set text encoder config.
func SetTextEncoderConfig(cfg TextEncoderConfig) {
	rwmutex.Lock()
	defer rwmutex.Unlock()
	configs.Default.TextEncoder = cfg
}

// SetJSONEncoderConfig - set json encoder config.
func SetJSONEncoderConfig(cfg JSONEncoderConfig) {
	rwmutex.Lock()
	defer rwmutex.Unlock()
	configs.Default.JSONEncoder = cfg
}

// SetCallerLevels - set caller levels.
func SetCallerLevels(levels ...Level) {
	rwmutex.Lock()
	defer rwmutex.Unlock()
	configs.Default.CallerLevels = levels
}

// SetStacktraceLevels - set stacktrace levels.
func SetStacktraceLevels(levels ...Level) {
	rwmutex.Lock()
	defer rwmutex.Unlock()
	configs.Default.StacktraceLevels = levels
}

// SetWriter - set writer.
func SetWriter(writer io.Writer) {
	rwmutex.Lock()
	defer rwmutex.Unlock()
	configs.Default.Handler.Type = "custom"
	configs.Default.Handler.Writer = writer
}

// LoadConfig - load config from file.
func LoadConfig(path string) error {
	var out Configs
	data, err := os.ReadFile(path)
	if err != nil {
		return errors.Wrap(err, "failed to read config content")
	}
	ext := filepath.Ext(path)
	switch ext {
	case ".json":
		if err := json.Unmarshal(data, &out); err != nil {
			return errors.Wrap(err, "failed to unmarshal json")
		}
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &out); err != nil {
			return errors.Wrap(err, "failed to unmarshal yaml")
		}
	default:
		return errors.Errorf("unsupported config file extension: %s", ext)
	}

	rwmutex.Lock()
	defer rwmutex.Unlock()
	configs = &out
	return nil
}

// SetModuleConfig - setting config for given module.
func SetModuleConfig(module string, cfg Config) {
	rwmutex.Lock()
	defer rwmutex.Unlock()
	configs.Modules[module] = cfg
}

// GetModuleConfig - getting config for given module.
func GetModuleConfig(module string) Config {
	rwmutex.RLock()
	defer rwmutex.RUnlock()
	cfg, exists := configs.Modules[module]
	if !exists {
		cfg = configs.Default
	}

	return cfg
}
