BIN      := karma-vine
CMD      := ./cmd/karma-vine
OUT      := ./build/$(BIN)
COV_OUT  := ./build/coverage.out
COV_HTML := ./build/coverage.html

.PHONY: build run test test-verbose test-race test-coverage test-coverage-html lint clean

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

## test-coverage: run tests and print a coverage summary
test-coverage:
	@mkdir -p ./build
	go test -coverprofile=$(COV_OUT) -covermode=atomic ./...
	go tool cover -func=$(COV_OUT)

## test-coverage-html: run tests and open an HTML coverage report
test-coverage-html: test-coverage
	go tool cover -html=$(COV_OUT) -o $(COV_HTML)
	@echo "Coverage report written to $(COV_HTML)"

## lint: vet all packages
lint:
	go vet ./...

## clean: remove built binaries
clean:
	rm -rf ./bin
