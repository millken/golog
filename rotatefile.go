package golog

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	currentTime             = time.Now
	defaultBackupTimeFormat = "20060102"

	_ io.Writer = (*RotateFile)(nil)
)

// RotateFileConfig is the configuration for RotateFile.
type RotateFileConfig struct {
	// Filename is the file to write logs to.  Backup log files will be retained in the same directory.
	// It uses <processname>.log in os.TempDir() if empty.
	Filename string `json:"filename" yaml:"filename"`

	// MaxBackups is the maximum number of old log files to retain.
	MaxBackups int `json:"maxbackups" yaml:"maxbackups"`

	// BackupTimeFormat determines if the time used for formatting the backup file name
	BackupTimeFormat string `json:"backupTimeFormat" yaml:"backupTimeFormat"`

	// LocalTime determines if the time used for formatting the timestamps in
	// backup files is the computer's local time.  The default is to use UTC
	// time.
	LocalTime bool `json:"localtime" yaml:"localtime"`

	// Async determines if the log write should be async
	Async bool `json:"async" yaml:"async"`
}

// RotateFile rotate log to file
type RotateFile struct {
	cfg               RotateFileConfig
	file              *os.File
	bufferWriter      *bufio.Writer
	currentBackupName string
	mu                sync.Mutex
	wokerOnce         sync.Once
	workerCh          chan bool
}

// NewRotateFile creates a new RotateFile.
func NewRotateFile(cfg RotateFileConfig) (*RotateFile, error) {
	if cfg.MaxBackups < 0 {
		return nil, fmt.Errorf("maxbackups cannot be negative")
	}
	f := &RotateFile{
		cfg: cfg,
	}
	if err := f.open(); err != nil {
		return nil, err
	}
	return f, nil
}

// Filename generates the name of the logfile from the current time.
func (f *RotateFile) Filename() string {
	if f.cfg.Filename != "" {
		return f.cfg.Filename
	}
	name := filepath.Base(os.Args[0]) + ".log"
	ex, err := os.Executable()
	if err != nil {
		return filepath.Join(os.TempDir(), name)
	}
	return filepath.Join(filepath.Dir(ex), name)
}

func (f *RotateFile) dir() string {
	return filepath.Dir(f.Filename())
}

func (f *RotateFile) backupTimeFormat() string {
	if f.cfg.BackupTimeFormat != "" {
		return f.cfg.BackupTimeFormat
	}
	return defaultBackupTimeFormat
}

func (f *RotateFile) open() error {
	var err error
	t := currentTime()
	if !f.cfg.LocalTime {
		t = t.UTC()
	}
	f.currentBackupName = t.Format(f.backupTimeFormat())

	f.file, err = os.OpenFile(f.Filename(), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	f.bufferWriter = bufio.NewWriter(f.file)
	return err
}

// Close implements io.Closer, and closes the current logfile.
func (f *RotateFile) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.close()
}

// close closes the file if it is open.
func (f *RotateFile) close() error {
	if f.file == nil {
		return nil
	}
	if err := f.bufferWriter.Flush(); err != nil {
		return err
	}
	err := f.file.Close()
	f.file = nil
	return err
}

// rotate on new day
func (f *RotateFile) reopenIfNeeded() (bool, error) {
	if f.file == nil {
		return false, f.open()
	}
	t := currentTime()
	if !f.cfg.LocalTime {
		t = t.UTC()
	}
	if f.currentBackupName == t.Format(f.backupTimeFormat()) {
		return false, nil
	}
	return true, nil
}

// Write writes data to a file
func (f *RotateFile) Write(d []byte) (int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	rotate, err := f.reopenIfNeeded()
	if err != nil {
		return 0, err
	}
	if rotate {
		if err := f.rotate(); err != nil {
			return 0, err
		}
	}
	n, err := f.bufferWriter.Write(d)
	if err != nil {
		return n, err
	}
	if f.cfg.Async {
		return n, nil
	}
	return n, f.bufferWriter.Flush()
}

// Rotate close the existing log file and create a new one.
func (f *RotateFile) Rotate() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.rotate()
}

func (f *RotateFile) rotate() error {
	if err := f.close(); err != nil {
		return err
	}
	if err := f.openNew(); err != nil {
		return err
	}
	f.wokerOnce.Do(func() {
		f.workerCh = make(chan bool, 1)
		go func() {
			for range f.workerCh {
				f.doWorker()
			}
		}()
	})
	select {
	case f.workerCh <- true:
	default:
	}
	return nil
}

func (f *RotateFile) doWorker() {
	if f.cfg.MaxBackups == 0 {
		return
	}
	files, err := f.oldLogFiles()
	if err != nil {
		log.Println(err)
		return
	}
	if f.cfg.MaxBackups > 0 && f.cfg.MaxBackups < len(files) {
		for _, fi := range files[0 : len(files)-f.cfg.MaxBackups] {
			os.Remove(filepath.Join(f.dir(), fi.Name()))
		}
	}
}

func (f *RotateFile) oldLogFiles() ([]logInfo, error) {
	files, err := ioutil.ReadDir(f.dir())
	if err != nil {
		return nil, fmt.Errorf("can't read log file directory: %s", err)
	}
	logFiles := []logInfo{}

	filename := filepath.Base(f.Filename())
	ext := filepath.Ext(filename)
	prefix := filename[:len(filename)-len(ext)] + "-"

	for _, fi := range files {
		if fi.IsDir() {
			continue
		}
		filename := fi.Name()
		if strings.HasPrefix(filename, prefix) {
			ext := filepath.Ext(filename)
			// ext include unix timestamp of the logfile
			if len(ext) > 10 {
				timestamp, err := strconv.ParseInt(ext[1:], 10, 64)
				if err == nil {
					logFiles = append(logFiles, logInfo{timestamp, fi})
				}
			}
		}
	}

	sort.Sort(byFormatTime(logFiles))

	return logFiles, nil
}

// openNew opens a new log file for writing, moving any old log file out of the
// way.  This methods assumes the file has already been closed.
func (f *RotateFile) openNew() error {
	err := os.MkdirAll(f.dir(), 0755)
	if err != nil {
		return fmt.Errorf("can't make directories for new logfile: %s", err)
	}

	name := f.Filename()
	fi, err := os.Stat(name)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		// move the existing file
		newname := f.backupName()
		if err := os.Rename(name, newname); err != nil {
			return fmt.Errorf("can't rename log file: %s", err)
		}
	}
	return f.open()
}

func (f *RotateFile) backupName() string {
	dir := filepath.Dir(f.Filename())
	filename := filepath.Base(f.Filename())
	ext := filepath.Ext(filename)
	prefix := filename[:len(filename)-len(ext)]
	t := currentTime()
	if !f.cfg.LocalTime {
		t = t.UTC()
	}
	return filepath.Join(dir, fmt.Sprintf("%s-%s%s.%d", prefix, t.Format(f.backupTimeFormat()), ext, t.Unix()))
}

// logInfo is a convenience struct to return the filename and its embedded timestamp.
type logInfo struct {
	timestamp int64
	os.FileInfo
}

// byFormatTime sorts by newest time formatted in the name.
type byFormatTime []logInfo

func (b byFormatTime) Less(i, j int) bool {
	return b[i].timestamp < b[j].timestamp
}

func (b byFormatTime) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b byFormatTime) Len() int {
	return len(b)
}
