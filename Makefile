all: fmt vet lint test

test:
	go test -cover .

vet:
	go vet .
	go vet ./examples/native
	go vet ./examples/client

fmt:
	go fmt .
	go fmt ./examples/native
	go fmt ./examples/client

lint:
	golint .
	golint ./examples/native
	golint ./examples/client

profile-mem:
	mkdir -p ./bench
	go test -bench=$(NAME) -benchmem -run=None -cpu 1 \
	  -memprofile ./bench/mem.out -test.memprofilerate=1 -o ./bench/bench.a
	go tool pprof -web -alloc_space ./bench/bench.a ./bench/mem.out
	rm -rf ./bench

profile-cpu:
	mkdir -p ./bench
	go test -bench=$(NAME) -run=None -cpu 1 -cpuprofile ./bench/cpu.out \
	  -o ./bench/bench.a
	go tool pprof -web ./bench/bench.a ./bench/cpu.out
	rm -rf ./bench
