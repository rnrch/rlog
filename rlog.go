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

// Mode is the mode of zapr loggers.
type Mode int

var (
	// Production is the production mode of zapr
	Production Mode = 1
	// Development is the development mode of zapr
	Development Mode = 2
	// Example is the example mode of zapr
	Example Mode = 3
)

type options struct {
	mode   Mode
	logger logr.Logger
	v      int
}

// Options configures how we set up the logger.
type Options func(*options)

// WithMode lets you set the mode of a zapr logger.
func WithMode(mode Mode) func(*options) {
	return func(o *options) {
		o.mode = mode
	}
}

// WithLogger lets you set the logr logger.
func WithLogger(logger logr.Logger) func(*options) {
	return func(o *options) {
		o.logger = logger
	}
}

// WithVerbosity lets you set Logger verbosity.
func WithVerbosity(v int) func(*options) {
	return func(o *options) {
		o.v = v
	}
}

// Logger represents the ability to log messages
type Logger interface {
	Info(msg string, kvPairs ...interface{})
	Error(err error, msg string, kvPairs ...interface{})
	V(v int) Verbose
	WithValues(kvPairs ...interface{}) Logger
	WithName(name string) Logger
	SetLogger(logr logr.Logger)
	SetVerbosity(v int) Logger
}

// NewLogger returns a new Logger
func NewLogger(opts ...Options) (Logger, error) {
	o := options{
		mode: Production,
	}
	for _, opt := range opts {
		opt(&o)
	}
	var err error
	if o.logger == nil {
		o.logger, err = newZaprLogger(o.mode)
		if err != nil {
			return nil, err
		}
	}
	return &loggingT{
		logr: o.logger,
		v:    o.v,
	}, nil
}

var logging loggingT

func init() {
	l, err := newZaprLogger(Production)
	if err != nil {
		panic(err)
	}
	logging.logr = l
}

type loggingT struct {
	logr logr.Logger
	v    int
}

func (l *loggingT) Info(msg string, kvPairs ...interface{}) {
	l.logr.Info(msg, kvPairs...)
}

func (l *loggingT) Error(err error, msg string, kvPairs ...interface{}) {
	l.logr.Error(err, msg, kvPairs...)
}

func (l *loggingT) V(v int) Verbose {
	if l.v >= v {
		return newVerbose(true, l.logr)
	}
	return newVerbose(false, l.logr)
}

func (l *loggingT) WithValues(kvPairs ...interface{}) Logger {
	return &loggingT{
		l.logr.WithValues(kvPairs...),
		l.v,
	}
}

// SetLogger sets the backing logr implementation.
func (l *loggingT) SetLogger(logr logr.Logger) {
	l.logr = logr
}

func (l *loggingT) WithName(name string) Logger {
	return &loggingT{
		l.logr.WithName(name),
		l.v,
	}
}

func (l *loggingT) SetVerbosity(v int) Logger {
	l.v = v
	return l
}

// Verbose is a boolean type that implements logr and records weather it is enabled.
type Verbose struct {
	enabled bool
	logr    logr.Logger
}

func newVerbose(b bool, l logr.Logger) Verbose {
	return Verbose{b, l}
}

// Info logs a non-error message with the given key/value pairs as context when v is enabled.
func (v Verbose) Info(msg string, kvPairs ...interface{}) {
	if v.enabled {
		v.logr.Info(msg, kvPairs...)
	}
}

// Error logs an error, with the given message and key/value pairs as context when v is enabled.
func (v Verbose) Error(err error, msg string, kvPairs ...interface{}) {
	if v.enabled {
		v.logr.Error(err, msg, kvPairs...)
	}
}

// SetVerbosity sets the global verbosity against which all logs will be compared.
func SetVerbosity(v int) {
	logging.v = v
}

// Info logs a non-error message with the given key/value pairs as context.
func Info(msg string, kvPairs ...interface{}) {
	logging.logr.Info(msg, kvPairs...)
}

// Error logs an error, with the given message and key/value pairs as context.
func Error(err error, msg string, kvPairs ...interface{}) {
	logging.logr.Error(err, msg, kvPairs...)
}

// V returns a Verbose struct for a specific verbosity level, relative to this Logger.
func V(v int) Verbose {
	if logging.v >= v {
		return newVerbose(true, logging.logr)
	}
	return newVerbose(false, logging.logr)
}

// SwtichMode replaces the current rlog logger with a new zapr logger guarded by the input mode.
func SwtichMode(mode Mode) {
	l, err := newZaprLogger(mode)
	if err != nil {
		logging.logr.Error(err, "Switch mode", "mode", mode)
		return
	}
	logging = loggingT{
		logr: l,
		v:    logging.v,
	}
}

// SetLogger sets the backing logr implementation for rlog.
func SetLogger(logr logr.Logger) {
	logging.logr = logr
}

func newZaprLogger(mode Mode) (logr.Logger, error) {
	var l logr.Logger
	switch mode {
	case Development:
		logger, err := zap.NewDevelopment()
		if err != nil {
			return nil, err
		}
		l = zapr.NewLogger(logger)
	case Example:
		logger := zap.NewExample()
		l = zapr.NewLogger(logger)
	case Production:
		logger, err := zap.NewProduction()
		if err != nil {
			return nil, err
		}
		l = zapr.NewLogger(logger)
	}
	return l, nil
}
