PLATFORMS ?= linux/amd64,linux/arm64

IMAGE_TAG = unitecoder/counter-example:latest

build-amd64:
	GOOS=linux GOARCH=amd64 go build -a -o bin/function-linux-amd64 main.go

build-arm64:
	GOOS=linux GOARCH=arm64 go build -a -o bin/function-linux-arm64 main.go

docker-build: build-amd64 build-arm64
	docker buildx build --platform $(PLATFORMS) -t $(IMAGE_TAG) -f Dockerfile bin
