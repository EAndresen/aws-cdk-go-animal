.PHONY: build clean deploy synth

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/create lambda/create/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/list lambda/list/main.go

clean:
	rm -rf ./bin ./vendor

test: clean build
	go test -v ./...

synth: clean build
	cdk synth

deploy: clean build
	cdk deploy