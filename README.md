[![Test status](https://github.com/millken/golog/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/millken/golog/actions?workflow=test)
[![Coverage Status](https://coveralls.io/repos/github/millken/golog/badge.svg?branch=main)](https://coveralls.io/github/millken/golog?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/millken/golog)](https://goreportcard.com/report/github.com/millken/golog)
[![GoDev](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/millken/golog)
[![GitHub release](https://img.shields.io/github/release/millken/golog.svg)](https://github.com/millken/golog/releases)
# golog
Fast structure logger for Golang

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
    golog.Infof("hello %s", "world")
    golog.Info("hello world with fields", "a", 1, "b", true, "c", "string")
}

// Output: 
2022-08-07T21:48:21+08:00 WARN hello world
2022-08-07T21:48:21+08:00 INFO hello world with fields a=1 b=true c=string
```

> Note: By default log writes to `os.Stdout`, The default log level for is *info*

## Performance 
> Note: disabled time and colors

```
BenchmarkGlobal-8                       22448018                53.41 ns/op            0 B/op          0 allocs/op
BenchmarkGlobal_WithField-8             10720359               112.4 ns/op            32 B/op          1 allocs/op
BenchmarkLogConsole-8                   22745239                51.94 ns/op            0 B/op          0 allocs/op
BenchmarkLogConsole_WithField-8          4971008               216.8 ns/op             0 B/op          0 allocs/op
BenchmarkLogJSON-8                      11603274               106.6 ns/op            16 B/op          1 allocs/op
BenchmarkLogJSON_WithField-8             4482342               266.6 ns/op            16 B/op          1 allocs/op
```
