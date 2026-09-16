# Name of the output binary
BINARY_NAME := cv_gen
BUILD_DIR   := ./_build

.PHONY: all build run fmt vet tidy clean

all: build

# Build the CLI
build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) .

# Build and run (pass args with: make run ARGS="generate -i data.json -o out")
run: build
	$(BUILD_DIR)/$(BINARY_NAME) $(ARGS)

fmt:
	go fmt ./...

vet:
	go vet ./...

tidy:
	go mod tidy

# Clean build artifacts
clean:
	go clean
	rm -rf $(BUILD_DIR)
