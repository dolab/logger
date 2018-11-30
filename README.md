# logger

[![Build Status](https://travis-ci.org/dolab/logger.svg?branch=master)](https://travis-ci.org/dolab/logger)

Efficent and simpe logger for golang with more features supported.

- custom log levels

- custom log tags

- custom log colors (only worked for *NIX)

- custom log output and output format

# Install

```go
go get -u github.com/dolab/logger
```

# Usage
```go
import "github.com/dolab/logger"

func main() {
    log, _ := logger.New("stdout")
    log.SetLevel(logger.Ldebug)

    // normal
    log.Debug("Hello, logger!")
    log.Infof("Hello, %s!", "logger")

    // create new logger with tags based on log
    taggedLog := log.New("X-REQUEST-ID")
    taggedLog.Debug("Receive HTTP request")
    taggedLog.Warnf("Send response with %d.", 200)
}
```

# Output

- stdout = os.Stdout
- stderr = os.Stderr
- null | nil = os.DevNull
- path/to/file = os.OpenFile("path/to/file", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

# Level

- Ldebug = DEBUG
- Linfo = INFO
- Lwarn = WARN
- Lerror = ERROR
- Lfatal = FATAL
- Lpanic = PANIC
- Ltrace = Stack

# License

MIT

# Author

Spring MC
