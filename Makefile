build:
	@go build -o bin/tcp-server

run: build
	@./bin/tcp-server