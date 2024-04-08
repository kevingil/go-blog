# Go builder
FROM golang:1.22.1 as builder

WORKDIR /go/src/github.com/kevingil/blog/app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /app/app


# Final package
FROM alpine:3.18.4
WORKDIR /app
COPY --from=builder /app/app .
COPY --from=builder /go/src/github.com/kevingil/blog/app/.env .
COPY --from=builder /go/src/github.com/kevingil/blog/app/views ./views
COPY --from=builder /go/src/github.com/kevingil/blog/app/static ./static

RUN apk update && apk add ca-certificates

# Run app
EXPOSE 8080
CMD ["./app"]
