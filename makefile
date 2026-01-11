# Name of the output binary
BINARY_NAME := cv_gen

# Default target
all: build

# Build the Go application
build:
	go build -o ./_build/$(BINARY_NAME) ./main.go

# Clean the build
clean:
	rm -f ./_build/$(BINARY_NAME)