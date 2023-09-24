# Setup Builder
FROM golang:1.16 as builder

WORKDIR /go/src/github.com/kevingil/blog/app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o /app/app

RUN echo "PORT=\$PORT" > .env \
    && echo "MYSQL_HOST=\$MYSQL_HOST" >> .env \
    && echo "MYSQL_PORT=\$MYSQL_PORT" >> .env \
    && echo "MYSQL_USER=\$MYSQL_USER" >> .env \
    && echo "MYSQL_PASSWORD=\$MYSQL_PASSWORD" >> .env \
    && echo "MYSQL_DATABASE=\$MYSQL_DATABASE" >> .env \
    && echo "MYSQL_ROOT_PASSWORD=\$MYSQL_ROOT_PASSWORD" >> .env


# Small final image
# FROM alpine:latest
# RUN apk add libc6-compat
FROM debian:11-slim
WORKDIR /app
COPY --from=builder /app/app .
COPY --from=builder /go/src/github.com/kevingil/blog/app/.env .

# Run app
EXPOSE 8080
CMD ["./app"]
