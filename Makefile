all: fmt vet lint

vet:
	go vet .

fmt:
	go fmt .

lint:
	golint .
