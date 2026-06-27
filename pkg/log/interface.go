package log

import (
	"fmt"
	"log"
)

var Level = 0

const (
	DEBUG = 1 << iota
)

// LogPrintf is a wrapper around log.Printf that adds a prefix and can be replaced if you want to use a different
// logging scheme
var LogPrintf = func(prefix, format string, args ...interface{}) {
	log.Println(prefix, fmt.Sprintf(format, args...))
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

// Fatalf logs a formatted message with an "INFO" prefix, then terminates execution by causing a panic.
func Fatalf(format string, args ...interface{}) {
	LogPrintf("FATAL", format, args...)
	panic("fatal")
}
