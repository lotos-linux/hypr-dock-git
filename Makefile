get:
		go mod tidy

build:
		go build -v -o bin/inclayer-go .

run:
		go run .

exec:
		./bin/inclayer-go