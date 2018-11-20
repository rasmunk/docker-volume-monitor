FROM golang:1.10.3-alpine
RUN apk add git tzdata dep

WORKDIR /go/src/docker-volume-monitor
COPY . .

RUN dep ensure
RUN go build ./...
RUN go install -v ./...

ENTRYPOINT ["/go/bin/docker-volume-monitor"]
CMD ["-pruneUnused", "-interval", "10"]