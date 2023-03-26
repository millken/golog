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
$ go test -benchmem -run=^$ -bench ^Benchmark
goos: darwin
goarch: arm64
pkg: github.com/millken/golog
BenchmarkGlobal-8               25937742                45.86 ns/op            0 B/op          0 allocs/op
BenchmarkGlobal_WithField-8     12905550                92.81 ns/op           32 B/op          1 allocs/op
BenchmarkLogText-8              20500350                49.62 ns/op            0 B/op          0 allocs/op
BenchmarkLogText_WithField-8     7224022               165.2 ns/op             0 B/op          0 allocs/op
BenchmarkLogJSON-8              22373744                54.55 ns/op            0 B/op          0 allocs/op
BenchmarkLogJSON_WithField-8     6632361               179.6 ns/op             0 B/op          0 allocs/op
PASS
ok      github.com/millken/golog        9.307s
```
