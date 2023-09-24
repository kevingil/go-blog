# Setup Builder
FROM golang:1.16 as builder

WORKDIR /go/src/github.com/kevingil/blog/app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o /app/app

# Small final image
# FROM alpine:latest
# RUN apk add libc6-compat
FROM debian:11-slim
WORKDIR /app
COPY --from=builder /app/app .
COPY --from=builder /go/src/github.com/kevingil/blog/app/entrypoint.sh /entrypoint.sh
COPY --from=builder /go/src/github.com/kevingil/blog/app/.env .
COPY --from=builder /go/src/github.com/kevingil/blog/app/views ./views 
COPY --from=builder /go/src/github.com/kevingil/blog/app/assets ./assets

ADD https://raw.githubusercontent.com/eficode/wait-for/v2.1.0/wait-for /usr/local/bin/wait-for
RUN chmod +rx /usr/local/bin/wait-for /entrypoint.sh

RUN apt-get update && apt-get install -y netcat

# Run app
EXPOSE 8080
ENTRYPOINT [ "sh", "/entrypoint.sh" ]
