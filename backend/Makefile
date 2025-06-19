.PHONY: all build run clean test

BINARY_NAME=dev-genesis
BINARY_PATH=../$(BINARY_NAME)

all: build

build:
	@echo "🔨 Building Dev Genesis..."
	go build -o $(BINARY_PATH) cmd/server/main.go
	@echo "✅ Build complete: $(BINARY_PATH)"

run: build
	@echo "🚀 Starting Dev Genesis..."
	$(BINARY_PATH)

clean:
	@echo "🧹 Cleaning..."
	go clean
	rm -f $(BINARY_PATH)
	rm -f $(BINARY_PATH).exe

test:
	@echo "🧪 Running tests..."
	go test -v ./...

deps:
	@echo "📦 Installing dependencies..."
	go mod download
	go mod tidy

# Windows specific build
build-windows:
	@echo "🔨 Building for Windows..."
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_PATH).exe cmd/server/main.go
	@echo "✅ Windows build complete: $(BINARY_PATH).exe"