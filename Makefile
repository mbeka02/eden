build:
	@go build -o bin/eden

run: build
	@./bin/eden

test:
	@go test -v -cover ./...
