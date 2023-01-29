/*
 * Copyright (c) 2022. The Greater Heptavirate: programming lodge (https://github.com/TheGreaterHeptavirate). All Rights Reserved.
 *
 * All copies of this software (if not stated otherwise) are dedicated ONLY to personal, non-commercial use.
 */

package logger

import (
	"github.com/kpango/glg"
)

//go:generate stringer -type=LogLevel -trimprefix=LogLevel -output=log_level_string.go

// LogLevel is an abstraction from the currently used logging library.
type LogLevel byte

// log levels.
const (
	LogLevelNotSpecified LogLevel = iota
	LogLevelDebug
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
)

// Logger returns the LogLevel converted to currently
// used logging library's log level.
func (l LogLevel) Logger() glg.LEVEL {
	m := map[LogLevel]glg.LEVEL{
		LogLevelDebug: glg.DEBG,
		LogLevelInfo:  glg.INFO,
		LogLevelWarn:  glg.WARN,
		LogLevelError: glg.ERR,
		LogLevelFatal: glg.FATAL,
	}

	if v, ok := m[l]; ok {
		return v
	}

	panic("unknown log level")
}
