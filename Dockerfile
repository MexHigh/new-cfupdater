FROM golang:1.24.2 AS builder
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go install .

FROM alpine:latest
LABEL maintainer="Leon Schmidt"
RUN apk add --no-cache \
    libc6-compat \
    ca-certificates
WORKDIR /app
COPY --from=builder /go/bin/new-cfupdater .

ENTRYPOINT ["/app/new-cfupdater"]
