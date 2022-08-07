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

const (
	defaultLogLevel   = INFO
	defaultModuleName = ""
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
	// Encoding is the log encoding.  console or json.
	Encoding             string               `json:"encoding" yaml:"encoding"`
	ConsoleEncoderConfig ConsoleEncoderConfig `json:"consoleEncodingConfig" yaml:"consoleEncodingConfig"`
	JSONEncoderConfig    JSONEncoderConfig    `json:"jsonEncoderConfig" yaml:"jsonEncoderConfig"`
	//CallerLevels is the default levels for show caller info.
	CallerLevels []Level `json:"callerLevels" yaml:"callerLevels"`
	// StacktraceLevels is the default levels for show stacktrace.
	StacktraceLevels []Level      `json:"stacktraceLevels" yaml:"stacktraceLevels"`
	Writer           WriterConfig `json:"handler" yaml:"handler"`
}

// ConsoleEncoderConfig is the configuration for the console encoder.
type ConsoleEncoderConfig struct {
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
}

// JSONEncoderConfig is the configuration for the JSONEncoder.
type JSONEncoderConfig struct {
	// TimeFormat specifies the format for timestamp in output.
	TimeFormat string `json:"timeFormat" yaml:"timeFormat"`
	// DisableTimestamp disables the timestamp in output.
	DisableTimestamp bool `json:"disableTimestamp" yaml:"disableTimestamp"`
	// CallerSkipFrame is the number of stack frames to skip when reporting the calling function.
	CallerSkipFrame int `json:"callerSkipFrame" yaml:"callerSkipFrame"`
}

// WriterConfig is a configuration for a writer.
type WriterConfig struct {
	Type         string     `json:"type" yaml:"type"`
	CustomWriter io.Writer  `json:"-" yaml:"-"`
	FileConfig   FileConfig `json:"fileConfig" yaml:"fileConfig"`
}

// FileConfig is a configuration for a file writer.
type FileConfig struct {
	Path string `json:"path" yaml:"path"`
}

func newConfigs() *Configs {
	return &Configs{
		Modules: make(map[string]Config),
	}
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