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
	rlog.SetVerbosity(2)
	rlog.Info("hello", "start", true)
	rlog.Error(errors.New("error1"), "this is err1", "num", 1)
	rlog.V(2).Info("this is info2", "level", 2, "start", false)
	rlog.V(3).Error(errors.New("error2"), "this is err2", "level", 3)
	rlog.SetVerbosity(3)
	rlog.V(3).Info("this is info3", "level", 3, "start", false)

	logger := rlog.NewZaprLogger(rlog.Production)
	logger = logger.WithName("myLogger").WithValues("testLogger", true)
	logger.SetVerbosity(4)
	logger.V(1).Error(errors.New("logger err"), "hello", "v", 1)
	logger.Error(errors.New("logger err2"), "another err")
	logger.V(4).Info("info for logger", "v", 4)
	logger.V(5).Info("info for logger", "v", 5)
	logger.Info("info 2 for logger", "verbosity", false)

	new, _ := zap.NewDevelopment()
	newLogr := zapr.NewLogger(new)
	rlog.SetLogger(newLogr)

	rlog.V(2).Info("this is info2", "level", 2, "start", false)
	rlog.V(3).Error(errors.New("error2"), "this is err2", "level", 3)
	rlog.SetVerbosity(2)
	rlog.V(3).Info("this is info3", "level", 3, "start", false)

	rlog.SwtichMode(rlog.Production)
	rlog.V(5).Info("stiil should not appear")
	rlog.Info("back to production")

	logger = rlog.NewZaprLogger(rlog.Example)
	logger.Info("this is a example logger")
}
