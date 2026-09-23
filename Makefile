#@author Fred Brooker <git@gscloud.cz>

BINARY_NAME=cf
DIST_DIR=dist
GO_DIR=go
LDFLAGS=-s -w

.PHONY: all build buildwr wr install everything

all:
	@echo "build | buildwr | wr | install"

buildwr:
	@echo "Building Docker container ..."
	@bash ./build.sh
	@echo "✅ Done."

build:
	@echo "🐹 Building Go toolchains ..."
	@mkdir -p $(DIST_DIR)
	@cd $(GO_DIR) && go mod tidy
	
	@echo "  -> Linux amd64"
	@cd $(GO_DIR) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o ../$(DIST_DIR)/$(BINARY_NAME)-linux-amd64 main.go
	
	@echo "  -> Windows amd64"
	@cd $(GO_DIR) && CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o ../$(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe main.go
	
	@echo "  -> macOS amd64 (Intel)"
	@cd $(GO_DIR) && CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o ../$(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 main.go
	
	@echo "  -> macOS arm64 (Apple Silicon)"
	@cd $(GO_DIR) && CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o ../$(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 main.go

	@cp $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 $(BINARY_NAME)
	@chmod +x $(BINARY_NAME)
	@echo "✅ Done. Binaries are in ./$(DIST_DIR)/ + local linux '$(BINARY_NAME)' is updated."

wr:
	@bash ./run.sh

install:
	@echo "🚚 Installing '$(BINARY_NAME)' to /usr/local/bin..."
	@sudo cp $(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	@echo "🎉 '$(BINARY_NAME)' is ready to use."

everything: build install
