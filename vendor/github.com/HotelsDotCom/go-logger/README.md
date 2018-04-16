# logger

Enhanced standard [go logger](https://golang.org/pkg/log/) with different log levels:

* `DEBUG`
* `INFO`
* `ERROR`
* `FATAL`

Level is set using env. variable `LOGLEVEL` e.g.: `LOGLEVEL=INFO`. Setting lower level automatically enables the ones that
follow, for example setting log level to `INFO` will log `ERROR` and `FATAL` messages as well.

Logging using fatal `Fatal(...)` or `Fatalf(..., ...)` will by default call `os.Exit(1)`.

Logger can optionally log to a file (instead of stdout) by setting `LOGFILE` env. variable, e.g.: `LOGFILE=/var/log/test`.
Specified log file will be created if it does not exist, but the parent directory must exist.

## Usage

Import `import "github.com/HotelsDotCom/go-logger"` or use alias to use it same way as [go logger](https://golang.org/pkg/log/)
`import log "github.com/HotelsDotCom/go-logger"`.

Example:

```go
import log "github.com/HotelsDotCom/go-logger"

...

log.Debug("some debug")
log.Debugf("some formatted %s", "debug")

log.Fatal("logging fatal and exiting")
log.Fatalf("some formatted %s and exiting", "fatal")
```

## Output

Output format:

`[LogLevel] Date Time FileName:LineNumber: LogMessage`

Example:

```
[INFO] 2017/04/28 11:48:14 main.go:8: Some information
[ERROR] 2017/04/28 11:48:14 main.go:9: Some error
```

## Package loggertest

Project contains `loggertest` package for testing, to check logged messages and change `os.Exit(1)` behaviour of
`FATAL` logger, to `panic()`. This can be either recovered from or tested with `assert.Panics(...)`

### code under test

```go
package main

import "github.com/HotelsDotCom/go-logger"

func foo(val string) {
	// some code...
	logger.Infof("You passed %q", val)
	// some more code...
}


func bar(val string) {
	// some code...
	logger.Fatalf("Oops %q", val)
	// some more code...
}
```

### test code

```go
package main

import (
	"testing"
	"github.com/HotelsDotCom/go-logger/loggertest"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
)

func TestLogEntry(t *testing.T) {
	// setup loggertest
	loggertest.Init(loggertest.LogLevelInfo)
	defer loggertest.Reset()

	// execute test code
	foo("hello")

	// assert logging behaviour - should have logged e.g. `[INFO] 2017/10/03 14:13:31 foo.go:6: You passed "hello"`
	logMessages := loggertest.GetLogMessages()
	require.Len(t, logMessages, 1)
	assert.Equal(t, `You passed "hello"`, logMessages[0].Message)
}


func TestFatalLogEntry(t *testing.T) {
	// setup loggertest
	loggertest.Init(loggertest.LogLevelInfo)
	defer loggertest.Reset()

	// fatal log will now cause a panic so that we can verify the output
	defer func() {
		if r := recover(); r != nil {
			// assert logging behaviour - should have logged e.g. `[FATAL] 2017/10/03 14:26:48 foo.go:14: Oops "bad times"`
			logMessages := loggertest.GetLogMessages()
			require.Len(t, logMessages, 1)
			assert.Regexp(t, `^\[FATAL\] .+ Oops "bad times"$`, logMessages[0].RawMessage)
		}
	}()

	// execute test code
	bar("bad times")
}
```
