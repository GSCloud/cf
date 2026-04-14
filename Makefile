#@author Fred Brooker <git@gscloud.cz>

.PHONY: all build wr install clean

all:
	@echo "build | buildwr | wr | install"

buildwr:
	@echo "Building Docker container ..."
	@bash ./build.sh

build:
	@echo "🐹 Building Go ..."
	@cd go && go mod tidy && go build -ldflags="-s -w" -o cf main.go
	@cp go/cf .

wr:
	@bash ./run.sh

install:
	@echo "🚚 Installing 'cf' to /usr/local/bin..."
	@sudo cp go/cf /usr/local/bin/cf
	@sudo chmod +x /usr/local/bin/cf
	@echo "🎉 'cf' is ready to use."

everything: build install
