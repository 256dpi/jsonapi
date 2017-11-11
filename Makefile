all: fmt vet lint test

vet:
	go vet ./...

fmt:
	go fmt ./...

lint:
	golint ./...

test:
	go test ./...

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
