all: fmt vet lint

vet:
	go vet .
	go vet ./adapters/http
	go vet ./adapters/echo
	go vet ./examples/native
	go vet ./examples/echo

fmt:
	go fmt .
	go fmt ./adapters/http
	go fmt ./adapters/echo
	go fmt ./examples/native
	go fmt ./examples/echo

lint:
	golint .
	golint ./adapters/http
	golint ./adapters/echo
	golint ./examples/native
	golint ./examples/echo

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
