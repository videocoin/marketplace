FROM golang:1.15 as builder

WORKDIR /go/src/github.com/videocoin/marketplace
COPY . .

ENV GOOS=linux
ENV GOARCH=amd64

RUN make build


FROM registry.videocoin.net/cloud/marketplace:b9b53061d9ea9872beb9f44686e43439121e2682

#RUN apt-get update
#RUN apt-get install -y ca-certificates ffmpeg gpac

COPY --from=builder /go/src/github.com/videocoin/marketplace/api /api
COPY --from=builder /go/src/github.com/videocoin/marketplace/bin/marketplace /marketplace
COPY --from=builder /go/src/github.com/videocoin/marketplace/tools/goose/linux_amd64/goose /goose
COPY --from=builder /go/src/github.com/videocoin/marketplace/migrations /migrations

CMD ["/marketplace"]
