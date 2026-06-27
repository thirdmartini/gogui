package log

import (
	"fmt"
	stdlog "log"
	"strings"
)

var Level = 0

const (
	DEBUG = 1 << iota
)

// LogPrintf is a wrapper around log.Printf that adds a prefix and can be replaced if you want to use a different
// logging scheme
var LogPrintf = func(prefix, format string, args ...interface{}) {
	s := strings.TrimSuffix(fmt.Sprintf(format, args...), "\n")
	stdlog.Println(prefix, s)
}

// Println writes a message to the log with an "INFO" prefix.
func Println(args ...interface{}) {
	LogPrintf("INFO ", "%s", fmt.Sprint(args...))
}

// Panicf logs a formatted message with a "PANIC" prefix, then panics.
func Panicf(format string, args ...interface{}) {
	LogPrintf("PANIC", format, args...)
	panic(fmt.Sprintf(format, args...))
}

// Debugf logs a formatted debug message if the DEBUG log level is enabled.
func Debugf(format string, args ...interface{}) {
	if Level&DEBUG == DEBUG {
		LogPrintf("DEBUG", format, args...)
	}
}

// Printf formats and writes a message to the log with an "INFO" prefix using the specified format and arguments.
func Printf(format string, args ...interface{}) {
	LogPrintf("INFO ", format, args...)
}

func Warnf(format string, args ...interface{}) {
	LogPrintf("WARN ", format, args...)
}

// Fatalf logs a formatted message with an "INFO" prefix, then terminates execution by causing a panic.
func Fatalf(format string, args ...interface{}) {
	LogPrintf("FATAL", format, args...)
	panic("fatal")
}
