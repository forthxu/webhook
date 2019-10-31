# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=release
BINARY_MAC=$(BINARY_NAME).app
BINARY_LINUX=$(BINARY_NAME).bin
BINARY_WINDOWS=$(BINARY_NAME).exe

all: mac linux windows
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_MAC)
	rm -f $(BINARY_LINUX)
	rm -f $(BINARY_WINDOWS)
	rm -f *.log


# Cross compilation
mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_MAC) -v
linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_LINUX) -v
windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_WINDOWS) -v