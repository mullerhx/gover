APP_NAME = gopilot-bin
BUILD_DIR = bin

GO_FILES := $(shell find . -name '*.go' -type f)

.PHONY: all build install clean

all: build

build: $(GO_FILES)
	@mkdir -p $(BUILD_DIR)
	@echo "🔨 Building $(APP_NAME)..."
	@go build -buildvcs=false -o `pwd`/$(BUILD_DIR)/$(APP_NAME)
	@echo "✅ Build complete: $(BUILD_DIR)/$(APP_NAME)"

install:
	@echo "📦 Installing $(APP_NAME)..."
	@cp `pwd`/$(BUILD_DIR)/$(APP_NAME) /usr/local/bin/
	@cp ./gopilot /usr/local/bin/
	@echo "✅ Installed to: $(shell go env GOPATH)/bin/$(APP_NAME)"

clean:
	@echo "🧹 Cleaning..."
	@rm -rf $(BUILD_DIR)
	@echo "✅ Clean complete."
