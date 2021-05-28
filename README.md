# golog
fast structured logger for Go (golang)

## performance
```
BenchmarkFileHandler
BenchmarkFileHandler-4                                   2374411               501.7 ns/op             0 B/op          0 allocs/op
BenchmarkFileHandlerWithFields
BenchmarkFileHandlerWithFields-4                         1629920               712.4 ns/op             0 B/op          0 allocs/op
BenchmarkJSONFormatterFileHandler
BenchmarkJSONFormatterFileHandler-4                      2328375               510.1 ns/op             0 B/op          0 allocs/op
BenchmarkJSONFormatterFileHandlerWithFields
BenchmarkJSONFormatterFileHandlerWithFields-4            2015848               586.1 ns/op             0 B/op          0 allocs/op
BenchmarkLoggerNoHandler
BenchmarkLoggerNoHandler-4                              350686596                3.432 ns/op           0 B/op          0 allocs/op
BenchmarkLoggerNoHandlerWithFields
BenchmarkLoggerNoHandlerWithFields-4                    160966628                7.439 ns/op           0 B/op          0 allocs/op
BenchmarkStdlog
BenchmarkStdlog-4                                        3015333               398.4 ns/op             0 B/op          0 allocs/op
BenchmarkStdlogWithFields
BenchmarkStdlogWithFields-4                              2337860               507.4 ns/op             0 B/op          0 allocs/op
```