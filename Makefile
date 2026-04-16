#@author Fred Brooker <git@gscloud.cz>

.PHONY: all build wr install clean

all:
	@echo "build | buildwr | wr | install"

buildwr:
	@echo "Building Docker container ..."
	@bash ./build.sh
	@echo "✅ Done."

build:
	@echo "🐹 Building Go toolchains ..."
	@mkdir -p dist
	@cd go && go mod tidy
	
	@echo "  -> Linux amd64"
	@cd go && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ../dist/cf-linux-amd64 main.go
	
	@echo "  -> Windows amd64"
	@cd go && CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ../dist/cf-windows-amd64.exe main.go
	
	@echo "  -> macOS amd64 (Intel)"
	@cd go && CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o ../dist/cf-darwin-amd64 main.go
	
	@echo "  -> macOS arm64 (Apple Silicon)"
	@cd go && CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o ../dist/cf-darwin-arm64 main.go

	@cp dist/cf-linux-amd64 cf
	@chmod +x cf
	@echo "✅ Done. Binaries are in ./dist/ and local 'cf' is updated."

wr:
	@bash ./run.sh

install:
	@echo "🚚 Installing 'cf' to /usr/local/bin..."
	@sudo cp go/cf /usr/local/bin/cf
	@sudo chmod +x /usr/local/bin/cf
	@echo "🎉 'cf' is ready to use."

everything: build install
