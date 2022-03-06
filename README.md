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

// Output: 2021-05-29 15:48:06 INFO hello world
```

> Note: By default log writes to `os.Stderr`, The default log level for is *info*

## Performance 
> Note: disabled time and colors

```
goos: darwin
goarch: arm64
pkg: github.com/millken/golog
BenchmarkGlobalLogger
BenchmarkGlobalLogger-8                                  9920839               116.0 ns/op             0 B/op          0 allocs/op
BenchmarkLoggerNoHandler
BenchmarkLoggerNoHandler-8                              371435955                3.210 ns/op           0 B/op          0 allocs/op
BenchmarkLoggerNoHandlerWithFields
BenchmarkLoggerNoHandlerWithFields-8                    100000000               11.89 ns/op            0 B/op          0 allocs/op
BenchmarkStdlog
BenchmarkStdlog-8                                       11016960               109.1 ns/op             0 B/op          0 allocs/op
BenchmarkStdlogWithFields
BenchmarkStdlogWithFields-8                              6322047               192.4 ns/op             0 B/op          0 allocs/op
BenchmarkWriterHandler
BenchmarkWriterHandler-8                                 8436931               142.1 ns/op             0 B/op          0 allocs/op
BenchmarkWriterHandlerWithFields
BenchmarkWriterHandlerWithFields-8                       6333751               189.5 ns/op             0 B/op          0 allocs/op
BenchmarkJSONFormatterWriterHandler
BenchmarkJSONFormatterWriterHandler-8                    8484501               141.4 ns/op             0 B/op          0 allocs/op
BenchmarkJSONFormatterWriterHandlerWithFields
BenchmarkJSONFormatterWriterHandlerWithFields-8          5928202               202.2 ns/op             0 B/op          0 allocs/op
PASS
ok      github.com/millken/golog        12.724s
```