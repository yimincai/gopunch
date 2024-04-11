default: dev

run:
	@go run cmd/main.go

dev: clean
	@CompileDaemon -build="go build -race -o bin/server cmd/main.go" -command="./bin/server" -color=true -graceful-kill=true

clean:
	@go clean
	@-rm -f bin/server

test:
	@gotestsum --junitfile-hide-empty-pkg --format testname

tidy:
	@go mod tidy
	@go fmt ./...

binary-linux-amd64:
	@rm -rf ./bin/*
	@GOOS=linux GOARCH=amd64 go build -o ./bin/bot-linux-amd64 ./cmd

