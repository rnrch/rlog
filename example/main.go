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

	rlog.SwitchMode(rlog.Development)
	logger := rlog.WithName("development").WithValues("foo", "bar")
	logger.Info("hello world", "verbosity", logger.GetVerbosity())

	logger, err := rlog.New(zapr.NewLogger(zap.NewExample()))
	if err != nil {
		rlog.Error(err, "create new logger")
	}
	logger = logger.WithName("example").WithValues("hello", "world").SetVerbosity(10)
	logger.V(7).Info("new logger")
}
