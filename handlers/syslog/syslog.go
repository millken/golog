//go:build !windows && !nacl && !plan9
// +build !windows,!nacl,!plan9

package syslog

import (
	"log/syslog"

	"github.com/millken/golog"
	"github.com/millken/x/buffer"
)

// Syslog to send logs via syslog.
type Syslog struct {
	Writer        *syslog.Writer
	SyslogNetwork string
	SyslogRaddr   string
}

func NewSyslog(network, raddr string, priority syslog.Priority, tag string) (*Syslog, error) {
	w, err := syslog.Dial(network, raddr, priority, tag)
	return &Syslog{w, network, raddr}, err
}

func (s *Syslog) Encode(_ *buffer.Buffer, entry golog.Record) error {
	line := entry.Message
	switch entry.Level {
	case golog.PANIC:
		return s.Writer.Crit(line)
	case golog.FATAL:
		return s.Writer.Crit(line)
	case golog.ERROR:
		return s.Writer.Err(line)
	case golog.WARNING:
		return s.Writer.Warning(line)
	case golog.INFO:
		return s.Writer.Info(line)
	case golog.DEBUG:
		return s.Writer.Debug(line)
	default:
		return nil
	}
}

func (s *Syslog) Write(p []byte) (n int, err error) {
	return 0, nil
}
