# rlog

A minimal wrapper for [logr].

## Usage

```go

package main

import (
    "errors"

    "github.com/rnrch/rlog"
)

func main() {
    rlog.Info("hello!", "visible", true)
    rlog.SetVerbosity(3)
    rlog.V(5).Error(errors.New("err1"), "this should not be printed out", "visible", false)
    rlog.V(3).Error(errors.New("err2"), "this should be printed out", "visible", true, "level", 3)

    logger := rlog.NewLogger(rlog.Development).WithName("development").WithValues("mode", "dev")
    logger.Info("this is a new logger!", "visible", true)
    logger.Error(errors.New("err"), "error from new logger", "visible", true)
}
```

The default logr implementation is [zapr]. You can use `SetLogger` to change it.

The default zapr mode is production mode. Use `SwtichMode` to change it.

[logr]: https://github.com/go-logr/logr
[zapr]: https://github.com/go-logr/zapr
