FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/harpocrates
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /tmp/harpocrates

FROM alpine
RUN apk add --no-cache bash
COPY --from=builder /tmp/harpocrates /harpocrates
COPY docker-entrypoint.sh /

ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["/harpocrates"]