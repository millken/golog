# golog
Fast logger for Golang

## Feature
  - Zero Allocation
  - Multiple Handlers, Formatters
  - Customize 

## Installation

```bash
go get -u github.com/millken/golog
```

## Getting Started

### Simple Logging Example

For simple logging, import the global logger package **github.com/millken/golog**

```go
package main

import (
    "github.com/millken/golog"
)

func main() {
    golog.Info("hello world")
}

// Output: 2021-05-29 15:48:06 info  hello world
```
> Note: By default log writes to `os.Stderr`
> Note: The default log level for is *info*

## Performance on Apple M1 16G
```
BenchmarkFileHandler
BenchmarkFileHandler-8                                   3512186               332.7 ns/op             0 B/op          0 allocs/op
BenchmarkFileHandlerWithFields
BenchmarkFileHandlerWithFields-8                         2621605               428.4 ns/op             0 B/op          0 allocs/op
BenchmarkJSONFormatterFileHandler
BenchmarkJSONFormatterFileHandler-8                      3326818               360.3 ns/op             0 B/op          0 allocs/op
BenchmarkJSONFormatterFileHandlerWithFields
BenchmarkJSONFormatterFileHandlerWithFields-8            2949661               368.9 ns/op             0 B/op          0 allocs/op
BenchmarkLoggerNoHandler
BenchmarkLoggerNoHandler-8                              570979944                2.073 ns/op           0 B/op          0 allocs/op
BenchmarkLoggerNoHandlerWithFields
BenchmarkLoggerNoHandlerWithFields-8                    253720461                4.723 ns/op           0 B/op          0 allocs/op
BenchmarkStdlog
BenchmarkStdlog-8                                        4189542               284.7 ns/op             0 B/op          0 allocs/op
BenchmarkStdlogWithFields
BenchmarkStdlogWithFields-8                              3464176               345.8 ns/op             0 B/op          0 allocs/op
```