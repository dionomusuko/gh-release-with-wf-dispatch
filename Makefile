DOCKER_IMAGE_NAME := ghcr.io/dionomusuko/release-helper
DOCKER_SCOPE := docker-latest

.PHONY: test
test:
	go test -v ./...

.PHONY: build
build:
	CGO_ENABLED=0 go build

.PHONY: docker-setup-buildx
docker-setup-buildx:
	docker buildx create --use --driver docker-container

.PHONY: docker-build
docker-build: docker-setup-buildx
	docker buildx build . \
		--tag "${DOCKER_IMAGE_NAME}:latest" \
		--platform "linux/amd64" \
		--output "type=docker" \
		--cache-from "type=gha,scope=${DOCKER_SCOPE}" \
		--cache-to "type=gha,mode=max,scope=${DOCKER_SCOPE}"
