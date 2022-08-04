package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/millken/golog/internal/log"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const (
	defaultLogLevel   = log.INFO
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
	Level log.Level `json:"level" yaml:"level"`
	// Encoding is the log encoding.  console or json.
	Encoding             string               `json:"encoding" yaml:"encoding"`
	ConsoleEncoderConfig ConsoleEncoderConfig `json:"consoleEncodingConfig" yaml:"consoleEncodingConfig"`
	JSONEncoderConfig    JSONEncoderConfig    `json:"jsonEncoderConfig" yaml:"jsonEncoderConfig"`
	//CallerLevels is the default levels for show caller info.
	CallerLevels []log.Level `json:"callerLevels" yaml:"callerLevels"`
	// StacktraceLevels is the default levels for show stacktrace.
	StacktraceLevels []log.Level  `json:"stacktraceLevels" yaml:"stacktraceLevels"`
	Writer           WriterConfig `json:"handler" yaml:"handler"`
}

type ConsoleEncoderConfig struct {
	// PartsOrder is the order of the parts of the log entry.
	PartsOrder []string `json:"partsOrder" yaml:"partsOrder"`
	// TimeFormat specifies the format for timestamp in output.
	TimeFormat string `json:"timeFormat" yaml:"timeFormat"`
	// DisableTimestamp disables the timestamp in output.
	DisableTimestamp bool `json:"disableTimestamp" yaml:"disableTimestamp"`
}

type JSONEncoderConfig struct {
	// TimeFormat specifies the format for timestamp in output.
	TimeFormat string `json:"timeFormat" yaml:"timeFormat"`
	// DisableTimestamp disables the timestamp in output.
	DisableTimestamp bool `json:"disableTimestamp" yaml:"disableTimestamp"`
}

type WriterConfig struct {
	Type       string     `json:"type" yaml:"type"`
	FileConfig FileConfig `json:"fileConfig" yaml:"fileConfig"`
}

type FileConfig struct {
	Path string `json:"path" yaml:"path"`
}

func newConfigs() *Configs {
	return &Configs{
		Modules: make(map[string]Config),
	}
}

// Load - load config from file.
func Load(path string) error {
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

// GetLevel - getting log level for given module.
func GetLevel(module string) log.Level {
	rwmutex.RLock()
	defer rwmutex.RUnlock()
	cfg, exists := configs.Modules[module]
	if !exists {
		cfg = configs.Default
	}

	return cfg.Level
}

func GetModuleConfig(module string) Config {
	rwmutex.RLock()
	defer rwmutex.RUnlock()
	cfg, exists := configs.Modules[module]
	if !exists {
		cfg = configs.Default
	}

	return cfg
}

// IsEnabledFor - Check if given log level is enabled for given module.
func IsEnabledFor(module string, level log.Level) bool {
	return level <= GetLevel(module)
}

// IsCallerEnabled returns if caller info enabled for given module and level.
func IsCallerEnabled(module string, level log.Level) bool {
	cfg := GetModuleConfig(module)
	for _, l := range cfg.CallerLevels {
		if l == level {
			return true
		}
	}
	return false
}
