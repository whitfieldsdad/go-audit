default: windows

windows: windows-amd64

NAME = go-audit
BUILD_TIME = `date +"%Y-%m-%dT%H:%M:%S%z"`
GIT_COMMIT_HASH = `git rev-parse HEAD`

LDFLAGS = "-s -w"

help:
	@echo "make [windows]"
	@echo "make [windows-amd64|windows-arm64]"
	@echo "make [test]"
	@echo "make [run]"

run:
	go run main.go

windows-amd64:
	GOOS=windows GOARCH=amd64 go build -o bin/${NAME}-windows-amd64.exe main.go

test:
	go test -v .
