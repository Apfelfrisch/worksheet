build-native:
	@echo "Building..."
	@go build -o ./build/timesheet ./cmd/
build-windows:
	@echo "Building for Windows..."
	@env GOOS=windows GOARCH=amd64 go build -o ./build/timesheet.exe ./cmd/
build-mac-arm:
	@echo "Building for Mac Arms..."
	@env GOOS=darwin GOARCH=arm64 go build -o ./build/timesheet-arm ./cmd/
run: build-native
	./build/timesheet
proto:
	@protoc ./game/network/payload/payload.proto --go_out=.
test:
	@go test ./...
