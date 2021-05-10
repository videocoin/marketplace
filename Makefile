NAME=marketplace
VERSION?=$$(git rev-parse HEAD)

REGISTRY_SERVER?=registry.videocoin.net
REGISTRY_PROJECT?=cloud

default: build
version:
	@echo ${VERSION}

build:
	GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=0 \
		go build \
			-mod vendor \
			-ldflags="-w -s -X main.Version=${VERSION}" \
			-o bin/${NAME} \
			./cmd/${NAME}/main.go

.PHONY: vendor
vendor:
	go mod vendor

docker-build:
	docker build -t ${REGISTRY_SERVER}/${REGISTRY_PROJECT}/${NAME}:${VERSION} -f Dockerfile .

docker-push:
	docker push ${REGISTRY_SERVER}/${REGISTRY_PROJECT}/${NAME}:${VERSION}

release: docker-build docker-push
