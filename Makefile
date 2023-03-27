.PHONY: all deps test benchmark

all: deps test 

deps:
	go mod download

test:
	go test -race -covermode=atomic

benchmark:
	go test -benchmem -run=^$$ -bench ^Benchmark