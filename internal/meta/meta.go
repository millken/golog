package meta

import (
	"errors"
	"strings"
	"sync"

	"github.com/millken/golog/internal/log"
)

//nolint:gochecknoglobals
var (
	rwmutex     = &sync.RWMutex{}
	levels      = newModuledLevels()
	handlers    = newModuledHandlers()
	callerInfos = newCallerInfo()
)

// SetDefaultLevel - setting default log level for all modules.
func SetDefaultLevel(level log.Level) {
	rwmutex.Lock()
	defer rwmutex.Unlock()
	levels.SetDefaultLevel(level)
}

// SetLevel - setting log level for given module.
func SetLevel(module string, level log.Level) {
	rwmutex.Lock()
	defer rwmutex.Unlock()
	levels.SetLevel(module, level)
}

// GetLevel - getting log level for given module.
func GetLevel(module string) log.Level {
	rwmutex.RLock()
	defer rwmutex.RUnlock()

	return levels.GetLevel(module)
}

// SetDefaultHandler - setting default log handler for all modules.
func SetDefaultHandler(hander log.Handler) {
	rwmutex.Lock()
	defer rwmutex.Unlock()
	handlers.SetDefaultHandler(hander)
}

// SetHandler - setting log handler for given module.
func SetHandler(module string, hander log.Handler) {
	rwmutex.Lock()
	defer rwmutex.Unlock()
	handlers.SetHandler(module, hander)
}

// GetHandler - getting log handler for given module.
func GetHandler(module string) log.Handler {
	rwmutex.RLock()
	defer rwmutex.RUnlock()

	return handlers.GetHandler(module)
}

// IsEnabledFor - Check if given log level is enabled for given module.
func IsEnabledFor(module string, level log.Level) bool {
	rwmutex.RLock()
	defer rwmutex.RUnlock()

	return levels.IsEnabledFor(module, level)
}

// ShowCallerInfo - Show caller info in log lines for given log level and module.
func ShowCallerInfo(module string, level log.Level) {
	rwmutex.Lock()
	defer rwmutex.Unlock()
	callerInfos.ShowCallerInfo(module, level)
}

// HideCallerInfo - Do not show caller info in log lines for given log level and module.
func HideCallerInfo(module string, level log.Level) {
	rwmutex.Lock()
	defer rwmutex.Unlock()
	callerInfos.HideCallerInfo(module, level)
}

// IsCallerInfoEnabled - returns if caller info enabled for given log level and module.
func IsCallerInfoEnabled(module string, level log.Level) bool {
	rwmutex.RLock()
	defer rwmutex.RUnlock()

	return callerInfos.IsCallerInfoEnabled(module, level)
}

// levelNames - log level names in string.
var levelNames = []string{ //nolint:gochecknoglobals
	"PANIC",
	"FATAL",
	"ERROR",
	"WARNING",
	"INFO",
	"DEBUG",
}

// ParseLevel returns the log level from a string representation.
func ParseLevel(level string) (log.Level, error) {
	for i, name := range levelNames {
		if strings.EqualFold(name, level) {
			return log.Level(i), nil
		}
	}

	return log.ERROR, errors.New("invalid log level")
}

// ParseString returns string representation of given log level.
func ParseString(level log.Level) string {
	return levelNames[level]
}
