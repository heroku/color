PROJECT_ROOT := $(shell pwd)

.PHONY: test bench

test:
	go test -race -coverprofile c.out $(PROJECT_ROOT)/...

bench:
	go test -bench=. -count=5  -cpu=1,2,4,8 -cpuprofile=cpu.prof github.com/murphybytes/color