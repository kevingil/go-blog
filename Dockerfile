# Go builder
FROM golang:1.22.1 as builder

WORKDIR /go/src/github.com/kevingil/blog

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /app


# Final package
FROM alpine:3.18.4
WORKDIR /app
COPY --from=builder /app .
COPY --from=builder /go/src/github.com/kevingil/blog/.env .
COPY --from=builder /go/src/github.com/kevingil/blog/internal/views ./internal/views
COPY --from=builder /go/src/github.com/kevingil/blog/static ./static

RUN apk update && apk add ca-certificates

# Run app
EXPOSE 8080
CMD ["./app"]
