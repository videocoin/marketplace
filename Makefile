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
	# https://github.com/ethereum/go-ethereum/issues/2738
	cp -r $(GOPATH)/src/github.com/ethereum/go-ethereum/crypto/secp256k1/libsecp256k1 \
	vendor/github.com/ethereum/go-ethereum/crypto/secp256k1/
	cp -r $(GOPATH)/src/github.com/karalabe/usb/hidapi \
	vendor/github.com/karalabe/usb/hidapi/
	cp -r $(GOPATH)/src/github.com/karalabe/usb/libusb \
	vendor/github.com/karalabe/usb/libusb/

docker-build:
	docker build -t ${REGISTRY_SERVER}/${REGISTRY_PROJECT}/${NAME}:${VERSION} -f Dockerfile .

docker-push:
	docker push ${REGISTRY_SERVER}/${REGISTRY_PROJECT}/${NAME}:${VERSION}

release: docker-build docker-push
