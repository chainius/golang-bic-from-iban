# Build Geth in a stock Go builder container
FROM golang:1.11-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers

ADD . /go-bic
RUN cd /go-bic && make bic

# Pull Geth into a second stage deploy alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /go-bic/build/bin/bic /usr/local/bin/

WORKDIR /var/www
COPY data /var/www/data

RUN addgroup -g 1000 geth && \
    adduser -h /root -D -u 1000 -G geth geth && \
    chown geth:geth /root

USER geth

EXPOSE 8080
ENTRYPOINT ["bic"]
