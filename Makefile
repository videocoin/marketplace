NAME=marketplace
VERSION?=$$(git rev-parse HEAD)

default: build
version:
	@echo ${VERSION}

build:
	GOOS=${GOOS} GOARCH=${GOARCH} go build \
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

db-status:
	goose -dir migrations postgres ${DB_URI} status

db-up:
	goose -dir migrations postgres ${DB_URI} up

db-down:
	goose -dir migrations postgres ${DB_URI} down
