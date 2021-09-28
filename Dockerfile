FROM golang:1.17.1-alpine
RUN apk add git tzdata dep

WORKDIR /go/src/docker-volume-monitor
COPY . .

ENV GO111MODULE="on"

RUN dep ensure
RUN go mod init
RUN go mod vendor
RUN go build ./...
RUN go install -v ./...

ENTRYPOINT ["/go/bin/docker-volume-monitor"]
CMD ["-prune-unused", "-interval", "10"]