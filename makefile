BINARY_NAME=markee
BUILD_DIR=build
VERSION = 0_0_1
OPTIONS = CGO_ENABLED=0
default: build

build:
	# 编译为 macOS 平台 amd64
	GOOS=darwin GOARCH=amd64 $(OPTIONS) go build -trimpath -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)_$(VERSION)_darwin_amd64 main.go
	# 编译为 macOS 平台 arm64
	GOOS=darwin GOARCH=arm64 $(OPTIONS) go build -trimpath -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)_$(VERSION)_darwin_arm64 main.go


	# 编译为 linux 平台 amd64
	GOOS=linux GOARCH=amd64 $(OPTIONS) go build -trimpath -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)_$(VERSION)_linux_amd64 main.go

	# 编译为 linux 平台 arm64
	GOOS=linux GOARCH=arm64 $(OPTIONS)  go build -trimpath -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)_$(VERSION)_linux_arm64 main.go


	# 编译为 BSD 平台 amd64
	GOOS=freebsd GOARCH=amd64 $(OPTIONS) go build -trimpath -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)_$(VERSION)_freebsd_amd64 main.go
	# 编译为 BSD 平台 arm64
	GOOS=freebsd GOARCH=arm64 $(OPTIONS) go build -trimpath -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)_$(VERSION)_freebsd_arm64 main.go


	# 编译为 Windows 平台 arm64
	GOOS=windows GOARCH=amd64 $(OPTIONS) go build -trimpath -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)_$(VERSION)_windows_arm64.exe main.go

	# 编译为 Windows 平台 amd64
	GOOS=windows GOARCH=amd64 $(OPTIONS) go build -trimpath -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)_$(VERSION)_windows_amd64.exe main.go

clean:
	rm -rf $(BUILD_DIR)/$(BINARY_NAME)_$(VERSION)_*

.PHONY: build clean