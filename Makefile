get:
		go mod tidy

build:
		go build -v -o bin/hypr-dock .

run:
		go run .

exec:
		./bin/hypr-dock