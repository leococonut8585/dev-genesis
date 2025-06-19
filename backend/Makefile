.PHONY: all build run clean test

BINARY_NAME=dev-genesis
BINARY_PATH=../$(BINARY_NAME)

all: build

build:
	@echo "🔨 Building Dev Genesis..."
	cd .. && go build -o $(BINARY_NAME) backend/cmd/server/main.go
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
	cd .. && GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME).exe backend/cmd/server/main.go
	@echo "✅ Windows build complete: $(BINARY_PATH).exe"