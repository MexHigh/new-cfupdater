FROM go:1.17 AS builder
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go install .

FROM alpine:latest
LABEL maintainer="Leon Schmidt"
WORKDIR /app
COPY --from=builder /go/bin/new-cfupdater .
COPY config.example.json config.json

ENTRYPOINT ["/app/new-cfupdater"]