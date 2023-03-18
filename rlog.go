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

const (
	Production Mode = iota + 1
	Development
	Example
)

type options struct {
	v int
}

// Options configures how we set up the logger.
type Options func(*options)

// WithVerbosity sets Logger verbosity.
func WithVerbosity(v int) func(*options) {
	return func(o *options) {
		o.v = v
	}
}

var defaultOptions = options{
	v: 0,
}

// Logger represents implementation of rlog logger
type Logger interface {
	Info(msg string, kvPairs ...interface{})
	Error(err error, msg string, kvPairs ...interface{})
	V(v int) Verbose
	WithValues(kvPairs ...interface{}) Logger
	WithName(name string) Logger
	SetLogger(logr logr.Logger)
	SetVerbosity(v int) Logger
	GetVerbosity() int
}

// New returns a new Logger
func New(logger logr.Logger, opts ...Options) (Logger, error) {
	o := defaultOptions
	for _, opt := range opts {
		opt(&o)
	}
	return &rloggerT{
		logger: logger,
		v:      o.v,
	}, nil
}

var rlogger rloggerT

func init() {
	l, err := newZaprLogger(Production)
	if err != nil {
		panic(err)
	}
	rlogger.logger = l
}

type rloggerT struct {
	logger logr.Logger
	v      int
}

func (l *rloggerT) Info(msg string, kvPairs ...interface{}) {
	l.logger.Info(msg, kvPairs...)
}

func (l *rloggerT) Error(err error, msg string, kvPairs ...interface{}) {
	l.logger.Error(err, msg, kvPairs...)
}

func (l *rloggerT) V(v int) Verbose {
	if l.v >= v {
		return newVerbose(true, l.logger)
	}
	return newVerbose(false, l.logger)
}

func (l *rloggerT) WithValues(kvPairs ...interface{}) Logger {
	return &rloggerT{
		logger: l.logger.WithValues(kvPairs...),
		v:      l.v,
	}
}

func (l *rloggerT) WithName(name string) Logger {
	return &rloggerT{
		logger: l.logger.WithName(name),
		v:      l.v,
	}
}

// SetLogger sets the backing logr implementation.
func (l *rloggerT) SetLogger(logr logr.Logger) {
	l.logger = logr
}

func (l *rloggerT) SetVerbosity(v int) Logger {
	l.v = v
	return l
}

func (l *rloggerT) GetVerbosity() int {
	return l.v
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

// Info logs a non-error message with the given key/value pairs as context.
func Info(msg string, kvPairs ...interface{}) {
	rlogger.logger.Info(msg, kvPairs...)
}

// Error logs an error, with the given message and key/value pairs as context.
func Error(err error, msg string, kvPairs ...interface{}) {
	rlogger.logger.Error(err, msg, kvPairs...)
}

// V returns a Verbose struct for a specific verbosity level, relative to this Logger.
func V(v int) Verbose {
	if rlogger.v >= v {
		return newVerbose(true, rlogger.logger)
	}
	return newVerbose(false, rlogger.logger)
}

func WithValues(kvPairs ...interface{}) Logger {
	return &rloggerT{
		logger: rlogger.logger.WithValues(kvPairs...),
		v:      rlogger.v,
	}
}

func WithName(name string) Logger {
	return &rloggerT{
		logger: rlogger.logger.WithName(name),
		v:      rlogger.v,
	}
}

// SetLogger sets the backing logr implementation for rlog.
func SetLogger(logr logr.Logger) {
	rlogger.logger = logr
}

// SetVerbosity sets the global verbosity against which all logs will be compared.
func SetVerbosity(v int) {
	rlogger.v = v
}

// SetVerbosity prints out the global verbosity
func GetVerbosity() int {
	return rlogger.v
}

// SwitchMode changes the dafult zapr logger mode
func SwitchMode(mode Mode) error {
	l, err := newZaprLogger(mode)
	if err != nil {
		return err
	}
	rlogger.logger = l
	return nil
}

func DefuaultLogger() Logger {
	return &rlogger
}

func newZaprLogger(mode Mode) (logr.Logger, error) {
	var l logr.Logger
	switch mode {
	case Development:
		logger, err := zap.NewDevelopment()
		if err != nil {
			return logr.Logger{}, err
		}
		l = zapr.NewLogger(logger)
	case Example:
		logger := zap.NewExample()
		l = zapr.NewLogger(logger)
	case Production:
		logger, err := zap.NewProduction()
		if err != nil {
			return logr.Logger{}, err
		}
		l = zapr.NewLogger(logger)
	}
	return l, nil
}
