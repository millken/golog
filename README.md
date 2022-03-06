[![Test status](https://github.com/millken/golog/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/millken/golog/actions?workflow=test)
[![Coverage Status](https://coveralls.io/repos/github/millken/golog/badge.svg?branch=main)](https://coveralls.io/github/millken/golog?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/millken/golog)](https://goreportcard.com/report/github.com/millken/golog)
[![GoDev](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/millken/golog)
[![GitHub release](https://img.shields.io/github/release/millken/golog.svg)](https://github.com/millken/golog/releases)
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
    golog.Infof("hello world")
}

// Output: 2021-05-29 15:48:06 info  hello world
```

> Note: By default log writes to `os.Stderr`, The default log level for is *info*

## Performance on Mac M1 16G
```
BenchmarkFileHandler
BenchmarkFileHandler-8                                   3563302               325.6 ns/op             0 B/op          0 allocs/op
BenchmarkFileHandlerWithFields
BenchmarkFileHandlerWithFields-8                         2612840               461.6 ns/op             0 B/op          0 allocs/op
BenchmarkJSONFormatterFileHandler
BenchmarkJSONFormatterFileHandler-8                      3396434               353.3 ns/op             0 B/op          0 allocs/op
BenchmarkJSONFormatterFileHandlerWithFields
BenchmarkJSONFormatterFileHandlerWithFields-8            2952370               407.1 ns/op             0 B/op          0 allocs/op
BenchmarkGlobalLogger
BenchmarkGlobalLogger-8                                  3920433               305.7 ns/op             0 B/op          0 allocs/op
BenchmarkLoggerNoHandler
BenchmarkLoggerNoHandler-8                              371451189                3.197 ns/op           0 B/op          0 allocs/op
BenchmarkLoggerNoHandlerWithFields
BenchmarkLoggerNoHandlerWithFields-8                    100000000               11.50 ns/op            0 B/op          0 allocs/op
BenchmarkStdlog
BenchmarkStdlog-8                                        3984559               301.1 ns/op             0 B/op          0 allocs/op
BenchmarkStdlogWithFields
BenchmarkStdlogWithFields-8                              2733052               436.1 ns/op             0 B/op          0 allocs/op
```