FROM golang:1.10.3-alpine

WORKDIR /go/src/docker-volume-monitor
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["docker-volume-monitor"]