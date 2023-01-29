/*
 * Copyright (c) 2022. The Greater Heptavirate: programming lodge (https://github.com/TheGreaterHeptavirate). All Rights Reserved.
 *
 * All copies of this software (if not stated otherwise) are dedicated ONLY to personal, non-commercial use.
 */

// Package logger contains an abstraction from the currently used logging library.
package logger

import "github.com/kpango/glg"

func init() {
	// for glg, need to change caller's lenght in order to report real caller path
	// instead of path to this file.
	glg.Get().SetCallerDepth(3)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	if err := glg.Info(args...); err != nil {
		panic(err)
	}
}

// Infof logs a message at level Info on the standard logger.
// It uses fmt.Sprintf to format the message.
func Infof(format string, args ...interface{}) {
	if err := glg.Infof(format, args...); err != nil {
		panic(err)
	}
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	if err := glg.Debug(args...); err != nil {
		panic(err)
	}
}

// Debugf logs a message at level Debug on the standard logger.
// It uses fmt.Sprintf to format the message.
func Debugf(format string, args ...interface{}) {
	if err := glg.Debugf(format, args...); err != nil {
		panic(err)
	}
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	if err := glg.Warn(args...); err != nil {
		panic(err)
	}
}

// Warnf logs a message at level Warn on the standard logger.
// It uses fmt.Sprintf to format the message.
func Warnf(format string, args ...interface{}) {
	if err := glg.Warnf(format, args...); err != nil {
		panic(err)
	}
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	if err := glg.Error(args...); err != nil {
		panic(err)
	}
}

// Errorf logs a message at level Error on the standard logger.
// It uses fmt.Sprintf to format the message.
func Errorf(format string, args ...interface{}) {
	if err := glg.Errorf(format, args...); err != nil {
		panic(err)
	}
}

// Fatal logs a message at level Fatal on the standard logger.
func Fatal(args ...interface{}) {
	glg.Fatal(args...)
}

// Fatalf logs a message at level Fatal on the standard logger.
// It uses fmt.Sprintf to format the message.
func Fatalf(format string, args ...interface{}) {
	glg.Fatalf(format, args...)
}

// Success logs a message at level Info on the standard logger.
// This announces a successful operation.
func Success(args ...interface{}) {
	if err := glg.Success(args...); err != nil {
		panic(err)
	}
}

// Successf logs a message at level Info on the standard logger.
// It uses fmt.Sprintf to format the message.
func Successf(format string, args ...interface{}) {
	if err := glg.Successf(format, args...); err != nil {
		panic(err)
	}
}

// SetLevel sets the logging level.
func SetLevel(l LogLevel) {
	glg.Get().SetLevel(l.Logger())
}
