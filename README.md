# logger

[![CircleCI](https://circleci.com/gh/dolab/logger.svg?style=svg)](https://circleci.com/gh/dolab/logger) [![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/dolab/logger) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/dolab/logger/master/LICENSE)

Efficient and simple logger for golang with more features supported.

- custom log levels
- custom log tags
- custom log colors (only worked for *NIX)
- custom log output and output format
- *structured log*

# Install

```go
go get -u github.com/dolab/logger
```

# Usage
```go
package main

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
    
    // or use struct log
    textLog := log.NewTextLogger()
    textLog.Str("key", "value").Err(err, true).Error("it's for demo")
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
