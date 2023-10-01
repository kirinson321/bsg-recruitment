.PHONY: default
default: build run

.PHONY: build
build: # Build the app binary.
	go build -o nbp main.go

.PHONY: run
run: # Run the app.
	./nbp

.PHONY: format
format: # Run the code formatter.
	gofumpt -s -w .
	gofumports -w .
