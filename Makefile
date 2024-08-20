build:
	go build -o bin/app ./cmd/web/

test:
	go test ./...

run: build
	./bin/app