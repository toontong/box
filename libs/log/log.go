package log

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

//log level, from low to high, more high means more serious
const (
	LevelTrace = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

const (
	Ltime  = 1 << iota //time format "2006/01/02 15:04:05"
	Lfile              //file.go:123
	Llevel             //[Trace|Debug|Info...]
)

var LogLevelString = map[string]int{
	"trace": LevelTrace,
	"debug": LevelDebug,
	"info":  LevelInfo,
	"warn":  LevelWarn,
	"error": LevelError,
	"fatal": LevelFatal,
}

var LevelName [6]string = [6]string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "FATAL"}

const TimeFormat = "2006/01/02 15:04:05"

const maxBufPoolSize = 16

type log func(format string, v ...interface{})

type Logger struct {
	Info    log
	Println log
	Debug   log
	Trace   log
	Warn    log
	Error   log
	Fatal   log

	sync.Mutex

	level int
	flag  int

	handler Handler

	quit chan struct{}
	msg  chan []byte

	bufs [][]byte

	wg sync.WaitGroup

	closed bool
}

//new a logger with specified handler and flag
func New(handler Handler, flag int) *Logger {
	var l = new(Logger)

	l.level = LevelInfo
	l.handler = handler

	l.flag = flag

	l.quit = make(chan struct{})
	l.closed = false

	l.msg = make(chan []byte, 1024)

	l.bufs = make([][]byte, 0, 16)

	l.wg.Add(1)

	l.Info = l.Infof
	l.Println = l.Infof
	l.Debug = l.Debugf
	l.Trace = l.Tracef
	l.Warn = l.Warnf
	l.Error = l.Errorf
	l.Fatal = l.Fatalf

	go l.run()

	return l
}

//new a default logger with specified handler and flag: Ltime|Lfile|Llevel
func NewDefault(handler Handler) *Logger {
	return New(handler, Ltime|Lfile|Llevel)
}

func newStdHandler() *StreamHandler {
	h, _ := NewStreamHandler(os.Stdout)
	return h
}

var std = NewDefault(newStdHandler())

type manager struct {
	mapper map[string]interface{}
	mu     sync.RWMutex
}

func newManager() *manager {
	m := new(manager)
	m.mapper = make(map[string]interface{})
	return m
}

func (self *manager) get(name string) *Logger {
	self.mu.Lock()
	defer self.mu.Unlock()

	l, ok := self.mapper[name]
	if ok {
		return l.(*Logger)
	} else {
		l = NewDefault(newStdHandler())
		self.mapper[name] = l
	}
	return l.(*Logger)
}

func (self *manager) close() {
	for _, v := range self.mapper {
		v.(*Logger).Close()
	}
}

var _mgr = newManager()

// like the python logging.getLogger
// return an Gloabl-logger and save in the memory
func GetLogger(name string) *Logger {
	if name == "" || name == "root" {
		return std
	}
	return _mgr.get(name)
}

func Close() {
	std.Close()
	_mgr.close()
}

func (l *Logger) run() {
	defer l.wg.Done()
	for {
		select {
		case msg := <-l.msg:
			l.handler.Write(msg)
			l.putBuf(msg)
		case <-l.quit:
			if len(l.msg) == 0 {
				return
			}
		}
	}
}

func (l *Logger) popBuf() []byte {
	l.Lock()
	var buf []byte
	if len(l.bufs) == 0 {
		buf = make([]byte, 0, 1024)
	} else {
		buf = l.bufs[len(l.bufs)-1]
		l.bufs = l.bufs[0 : len(l.bufs)-1]
	}
	l.Unlock()

	return buf
}

func (l *Logger) putBuf(buf []byte) {
	l.Lock()
	if len(l.bufs) < maxBufPoolSize {
		buf = buf[0:0]
		l.bufs = append(l.bufs, buf)
	}
	l.Unlock()
}

func (l *Logger) Close() {
	if l.closed {
		return
	}
	l.closed = true

	close(l.quit)
	l.wg.Wait()
	l.quit = nil

	l.handler.Close()
}

//set log level, any log level less than it will not log
func (l *Logger) SetLevel(level int) {
	l.level = level
}

func (l *Logger) Level() int {
	return l.level
}

func (l *Logger) SetHandler(h Handler) {
	l.handler = h
}

//a low interface, maybe you can use it for your special log format
//but it may be not exported later......
func (l *Logger) Output(callDepth int, level int, format string, v ...interface{}) {
	if l.level > level {
		return
	}

	buf := l.popBuf()

	if l.flag&Ltime > 0 {
		now := time.Now().Format(TimeFormat)
		buf = append(buf, now...)
		buf = append(buf, " - "...)
	}

	if l.flag&Llevel > 0 {
		buf = append(buf, LevelName[level]...)
		buf = append(buf, " - "...)
	}

	if l.flag&Lfile > 0 {
		_, file, line, ok := runtime.Caller(callDepth)
		if !ok {
			file = "???"
			line = 0
		} else {
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					file = file[i+1:]
					break
				}
			}
		}

		buf = append(buf, file...)
		buf = append(buf, ":["...)

		buf = strconv.AppendInt(buf, int64(line), 10)
		buf = append(buf, "] - "...)
	}

	s := fmt.Sprintf(format, v...)

	buf = append(buf, s...)

	if s[len(s)-1] != '\n' {
		buf = append(buf, '\n')
	}

	l.msg <- buf
}

//log with Trace level
func (l *Logger) Tracef(format string, v ...interface{}) {
	l.Output(2, LevelTrace, format, v...)
}

//log with Debug level
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.Output(2, LevelDebug, format, v...)
}

//log with info level
func (l *Logger) Infof(format string, v ...interface{}) {
	l.Output(2, LevelInfo, format, v...)
}

//log with warn level
func (l *Logger) Warnf(format string, v ...interface{}) {
	l.Output(2, LevelWarn, format, v...)
}

//log with error level
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Output(2, LevelError, format, v...)
}

//log with fatal level
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Output(2, LevelFatal, format, v...)
}

func SetLevel(level int) {
	std.SetLevel(level)
}

func SetLevelS(level string) {
	SetLevel(LogLevelString[strings.ToLower(level)])
}

func Tracef(format string, v ...interface{}) {
	std.Output(2, LevelTrace, format, v...)
}

func Debugf(format string, v ...interface{}) {
	std.Output(2, LevelDebug, format, v...)
}

func Infof(format string, v ...interface{}) {
	std.Output(2, LevelInfo, format, v...)
}

func Warnf(format string, v ...interface{}) {
	std.Output(2, LevelWarn, format, v...)
}

func Errorf(format string, v ...interface{}) {
	std.Output(2, LevelError, format, v...)
}

func Fatalf(format string, v ...interface{}) {
	std.Output(2, LevelFatal, format, v...)
}

func StdLogger() *Logger {
	return std
}

func GetLevel() int {
	return std.level
}

var (
	Info    = Infof
	Println = Infof
	Debug   = Debugf
	Trace   = Tracef
	Warn    = Warnf
	Error   = Errorf
	Fatal   = Fatalf
)
