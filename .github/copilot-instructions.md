# GitHub Copilot Instructions

This document provides guidance for AI coding agents on how to effectively contribute to the `logger` codebase.

## Project Overview

`logger` is a flexible and efficient logging library for Go. It supports both simple, `printf`-style logging and structured, key-value based logging. The library is designed to be extensible, allowing for custom log levels, output formats, and colorization.

## Core Components

The main components of the `logger` library are:

- **`logger.go`**: This file contains the primary `Logger` struct, which is the main entry point for all logging operations. It manages the log output, level, and formatting. When creating a new logger, you can specify the output destination, such as `stdout`, `stderr`, or a file path.

- **`struct.go`**: This file defines the `StructLogger` interface, which provides a fluent API for creating structured logs. You can create a `StructLogger` instance using `logger.NewTextLogger()` or `logger.NewJsonLogger()`. The fluent API allows you to chain calls to add fields to your log entries, like so:

  ```go
  log.NewTextLogger().Str("component", "api").Int("status", 200).Info("Request handled successfully")
  ```

- **`slog.go`**: This file provides a compatibility layer that allows `logger` to be used as a handler for Go's standard `log/slog` package. This is useful for integrating `logger` into projects that already use `slog`.

- **`levels.go` and `format.go`**: These files define the available log levels (e.g., `Ldebug`, `Linfo`) and output formats (e.g., `TextFormat`, `JSONFormat`).

## Development Workflow

### Running Tests

The project has a comprehensive test suite. To run all tests, use the following command:

```sh
go test ./...
```

When adding new features or fixing bugs, please ensure that you also add or update the corresponding tests.

### Dependencies

The project uses Go modules to manage dependencies. The primary external dependencies are:

- `github.com/dolab/colorize`: Used for adding color to log output in the terminal.
- `github.com/golib/assert`: Used in the test suite for making assertions.

## Key Patterns and Conventions

- **Structured Logging**: When adding new logging statements, prefer using the `StructLogger` for structured data. This makes the logs more machine-readable and easier to parse.

- **Immutability**: When creating a new logger from an existing one (e.g., with `log.New(...)`), a new `Logger` instance is returned. This ensures that loggers can be safely shared across different parts of an application without unintended side effects.

- **Error Handling**: When a logging operation fails, it should return an error. This allows the calling code to decide how to handle the failure.
