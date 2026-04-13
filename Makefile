BIN      := karma-vine
CMD      := ./cmd/karma-vine
OUT      := ./build/$(BIN)

.PHONY: build run test test-verbose lint clean

## build: compile a dev binary to ./build/karma-vine
build:
	go build -o $(OUT) $(CMD)

## run: build and run the explorer
run: build
	$(OUT)

## test: run all tests
test:
	go test ./...

## test-verbose: run all tests with verbose output
test-verbose:
	go test -v ./...

## test-race: run all tests with the race detector
test-race:
	go test -race ./...

## lint: vet all packages
lint:
	go vet ./...

## clean: remove built binaries
clean:
	rm -rf ./bin
