// Copyright 2020 rnrch
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rlog

import (
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

const (
	Production = iota
	Development
	Example
)

type Logger interface {
	Info(msg string, kvPairs ...interface{})
	Error(err error, msg string, kvPairs ...interface{})
	V(level int) Verbose
	WithValues(kvPairs ...interface{}) Logger
	WithName(name string) Logger
	SetLogger(logr logr.Logger)
	SetVerbosity(v int) Logger
}

type loggingT struct {
	logr  logr.Logger
	level int
}

func (l *loggingT) Info(msg string, kvPairs ...interface{}) {
	l.logr.Info(msg, kvPairs...)
}

func (l *loggingT) Error(err error, msg string, kvPairs ...interface{}) {
	l.logr.Error(err, msg, kvPairs...)
}

func (l *loggingT) V(level int) Verbose {
	if l.level >= level {
		return newVerbose(true, l.logr)
	}
	return newVerbose(false, l.logr)
}

func (l *loggingT) WithValues(kvPairs ...interface{}) Logger {
	return &loggingT{
		l.logr.WithValues(kvPairs...),
		l.level,
	}
}

func (l *loggingT) SetLogger(logr logr.Logger) {
	l.logr = logr
}

func (l *loggingT) WithName(name string) Logger {
	return &loggingT{
		l.logr.WithName(name),
		l.level,
	}
}

func (l *loggingT) SetVerbosity(v int) Logger {
	l.level = v
	return l
}

type Verbose struct {
	enabled bool
	logr    logr.Logger
}

func newVerbose(b bool, l logr.Logger) Verbose {
	return Verbose{b, l}
}

func (v Verbose) Info(msg string, kvPairs ...interface{}) {
	if v.enabled {
		v.logr.Info(msg, kvPairs...)
	}
}

func (v Verbose) Error(err error, msg string, kvPairs ...interface{}) {
	if v.enabled {
		v.logr.Error(err, msg, kvPairs...)
	}
}

var logging loggingT

func init() {
	logger, _ := zap.NewProduction()
	logging.logr = zapr.NewLogger(logger)
}

var globalVerbosity int = 0

func SetVerbosity(v int) Logger {
	globalVerbosity = v
	return &logging
}

func Info(msg string, kvPairs ...interface{}) {
	logging.logr.Info(msg, kvPairs...)
}

func Error(err error, msg string, kvPairs ...interface{}) {
	logging.logr.Error(err, msg, kvPairs...)
}

func V(level int) Verbose {
	if globalVerbosity >= level {
		return newVerbose(true, logging.logr)
	}
	return newVerbose(false, logging.logr)
}

func SwtichMode(mode int) {
	switch mode {
	case Development:
		logger, _ := zap.NewDevelopment()
		logging.logr = zapr.NewLogger(logger)
	case Example:
		logger := zap.NewExample()
		logging.logr = zapr.NewLogger(logger)
	case Production:
		logger, _ := zap.NewProduction()
		logging.logr = zapr.NewLogger(logger)
	default:
		logging.Info("failed to change rlog mode", "input mode value", mode)
	}
}

func SetLogger(logr logr.Logger) {
	logging.logr = logr
}

func NewLogger(mode int) Logger {
	var logger *zap.Logger
	switch mode {
	case Development:
		logger, _ = zap.NewDevelopment()
	case Example:
		logger = zap.NewExample()
	default:
		logger, _ = zap.NewProduction()
	}
	return &loggingT{
		logr: zapr.NewLogger(logger),
	}
}
