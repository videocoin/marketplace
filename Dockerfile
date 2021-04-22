FROM golang:1.14 as builder

WORKDIR /go/src/github.com/videocoin/marketplace
COPY . .

RUN make build


FROM alpine:3

RUN apk add ca-certificates ffmpeg bash

COPY --from=builder /go/src/github.com/videocoin/marketplace/bin/marketplace /marketplace

CMD ["/marketplace"]
