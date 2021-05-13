FROM golang:1.15 as builder

WORKDIR /go/src/github.com/videocoin/marketplace
COPY . .

RUN make build


FROM alpine:3

RUN apk add ca-certificates ffmpeg bash

COPY --from=builder /go/src/github.com/videocoin/marketplace/api /api
COPY --from=builder /go/src/github.com/videocoin/marketplace/bin/marketplace /marketplace
COPY --from=builder /go/src/github.com/videocoin/marketplace/bin/mpek /mpek
COPY --from=builder /go/src/github.com/videocoin/marketplace/tools/goose/linux_amd64/goose /goose
COPY --from=builder /go/src/github.com/videocoin/marketplace/migrations /migrations

CMD ["/marketplace"]
