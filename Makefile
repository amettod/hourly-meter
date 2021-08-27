## help: print this help message
help:
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## build_cli: build the cmd/cli application
build_cli:
	go build -ldflags="-s" -o=./bin/ ./cmd/cli

## build_cli_win: build the cmd/cli application for windows os
build_cli_win:
	GOOS=windows go build -ldflags="-s" -o=./bin/ ./cmd/cli

## build_web: build the cmd/web application
build_web:
	go build -ldflags="-s" -o=./bin/ ./cmd/web

## build_web_win: build the cmd/web application for windows os
build_web_web:
	GOOS=windows go build -ldflags="-s" -o=./bin/ ./cmd/web

## test: run all tests
test:
	go test ./...

## test_cli_manual: сreate files to view
test_cli_manual:
	go run ./cmd/cli -filename="internal/app/testdata/23456789_feb.html" \
		-contract="98765432" \
		-name="OOO STAR" \
		-coefficient=4000

## test_web_manual: create file to view
test_web_manual:
	go run ./cmd/web