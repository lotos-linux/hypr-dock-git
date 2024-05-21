get:
		go mod tidy

build:
		go build -v -o bin/inclayer-go .

run:
		go run .

exec:
		./bin/inclayer-go

push:
		git add .
		git status
		read commit
		git commit -m $commit
		git remote | xargs -L1 git push --all