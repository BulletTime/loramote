.PHONY: default all build config clean
VERSION := 0.2.0
COMMIT := $(shell git describe --always)
GOOS ?= darwin
GOARCH ?= amd64
BUILD_DATE = `date -u +%Y-%m-%dT%H:%M.%SZ`
BUILD_NAME = loramote
MAIN_FILE = main.go

.SILENT:
default: clean build

all: clean build config

build:
	echo "[===] Build for $(GOOS) $(GOARCH) [===]"
	mkdir -p build
	echo "[GO BUILD] $(MAIN_FILE)"
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=1 go build -a -ldflags "-X main.version=$(VERSION) -X main.build=$(COMMIT) -X main.buildDate=$(BUILD_DATE)" -o build/$(BUILD_NAME) $(MAIN_FILE)

config:
	test -f build/$(BUILD_NAME) || $(MAKE) build
	echo "[===] Writing config to: ~/.$(BUILD_NAME).yaml [===]"
	build/$(BUILD_NAME) config > ~/.$(BUILD_NAME).yaml

clean:
	echo "[===] Cleaning up workspace [===]"
	rm -rf build
	rm -rf $(BUILD_NAME).log
