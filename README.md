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
BenchmarkGlobal-8                       17820811                67.55 ns/op            0 B/op          0 allocs/op
BenchmarkGlobal_WithField-8             10794764               110.7 ns/op            32 B/op          1 allocs/op
BenchmarkLogConsole-8                   21667298                53.65 ns/op            0 B/op          0 allocs/op
BenchmarkLogConsole_WithField-8          4162293               289.5 ns/op             0 B/op          0 allocs/op
BenchmarkLogJSON-8                      11371420               105.0 ns/op            16 B/op          1 allocs/op
BenchmarkLogJSON_WithField-8             3810990               315.7 ns/op            16 B/op          1 allocs/op
```