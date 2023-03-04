# rlog

[![PkgGoDev](https://pkg.go.dev/badge/github.com/rnrch/rlog)](https://pkg.go.dev/github.com/rnrch/rlog)
[![Go Report Card](https://goreportcard.com/badge/github.com/rnrch/rlog)](https://goreportcard.com/report/github.com/rnrch/rlog)
![Github Actions](https://github.com/rnrch/rlog/workflows/CI/badge.svg)

A minimal wrapper for [logr].

## Usage

```go

package main

import (
	"errors"

	"github.com/go-logr/zapr"
	"github.com/rnrch/rlog"
	"go.uber.org/zap"
)

func main() {
	rlog.Info("hello", "default verbosity", rlog.GetVerbosity())
	rlog.Error(errors.New("error"), "error msg 1", "num", 1)

	rlog.V(2).Info("this is v2 info", "print", false)
	rlog.SetVerbosity(5)
	rlog.V(2).Info("this is v2 info again", "verbosity", rlog.GetVerbosity())
	rlog.V(7).Error(errors.New("error"), "error msg 2", "num", 2, "print", false)

	err := rlog.SwitchMode(rlog.Development)
	if err != nil {
		rlog.Error(err, "switch mode", "mode", rlog.Development)
	}
	logger := rlog.WithName("development").WithValues("foo", "bar")
	logger.Info("hello world", "verbosity", logger.GetVerbosity())

	logger, err = rlog.New(zapr.NewLogger(zap.NewExample()))
	if err != nil {
		rlog.Error(err, "create new logger")
	}
	logger = logger.WithName("example").WithValues("hello", "world").SetVerbosity(10)
	logger.V(7).Info("new logger")
}

```

The default logr implementation is [zapr]. You can use `SetLogger` to change it.

The default zapr mode is production mode. Use `SwtichMode` to change it.

[logr]: https://github.com/go-logr/logr
[zapr]: https://github.com/go-logr/zapr
