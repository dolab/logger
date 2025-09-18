# Bito AI Instructions for `dolab/logger`

This document provides guidance for AI coding agents to effectively contribute to the `dolab/logger` codebase.

## Big Picture

`dolab/logger` is a logging library for Go that prioritizes simplicity, efficiency, and structured logging. It is designed to be a flexible and configurable logging solution.

### Core Components

- **`logger.go`**: The main component, containing the `Logger` struct and methods for creating and configuring loggers. It handles log levels, output destinations, and formatting.
- **`fields.go`**: Implements structured logging by allowing key-value pairs to be added to log entries. It supports both text and JSON formats.
- **`levels.go`**: Defines the available log levels (e.g., `Ldebug`, `Linfo`, `Lwarn`) and provides functions for resolving levels by name.
- **`errors.go`**: Contains custom error types used within the library.

### Key Design Principles

- **Simplicity**: The library is designed to be easy to use, with a straightforward API for creating and using loggers.
- **Efficiency**: The use of `bufio` and `bytes.Buffer` helps to minimize performance overhead.
- **Structured Logging**: The library supports structured logging through the use of fields, which can be output in either text or JSON format.

## Developer Workflows

### Building and Testing

The project uses standard Go tools for building and testing. To run the tests, use the following command:

```sh
go test ./...
```

### Debugging

When debugging, you can use the `Ltrace` log level to get stack traces in your log output. This can be helpful for identifying the source of errors.

## Project-Specific Conventions

### Logging Levels

The library defines a set of custom log levels in `levels.go`. When adding new logging functionality, be sure to use these levels consistently.

### Structured Logging

When adding structured data to your logs, use the `Str`, `Bool`, `Err`, and `Any` functions in `fields.go`. This ensures that your logs are properly formatted.

### Error Handling

The library uses a custom error type, `ErrLevel`, which is defined in `errors.go`. When adding new error types, follow this pattern.

## Integration Points

The library can be integrated with any Go application. It supports logging to `stdout`, `stderr`, and files. The output format can be customized to meet the needs of your application.
