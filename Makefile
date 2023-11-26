default: windows

windows: windows-amd64 windows-arm64

NAME = go-audit
BUILD_TIME = `date +"%Y-%m-%dT%H:%M:%S%z"`
GIT_COMMIT_HASH = `git rev-parse HEAD`

LDFLAGS = "-s -w -X common.buildTime=${BUILD_TIME} -X common.gitCommitHash=${GIT_COMMIT_HASH}"

help:
	@echo "make [windows]"
	@echo "make [windows-amd64|windows-arm64]"
	@echo "make [test]"
	@echo "make [run]"

up:
	docker-compose up -d
	make run

down:
	docker-compose down --remove-orphans

run:
	go run main.go

windows-amd64:
	GOOS=windows GOARCH=arm64 go build -ldflags ${LDFLAGS} -o bin/${NAME}-windows-amd64.exe main.go

windows-arm64:
	GOOS=windows GOARCH=arm64 go build -ldflags ${LDFLAGS} -o bin/${NAME}-windows-arm64.exe main.go

test:
	go test -v .