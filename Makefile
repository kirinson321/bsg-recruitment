.PHONY: default
default: build run

.PHONY: build
build: # Build the app binary.
	go build -o nbp main.go

.PHONY: run
run: # Run the app on localhost, port 8000.
	./nbp

.PHONY: format
format: # Run the code formatter.
	gofumpt -s -w .
	gofumports -w .

.PHONY: dockerize
dockerize: # Build the app binary and build the Docker image.
	build
	docker build --tag nbp .