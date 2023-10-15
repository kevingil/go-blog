# Go builder
FROM golang:1.21 as builder

WORKDIR /go/src/github.com/kevingil/blog/app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /app/app

# Final package
FROM debian:11-slim
WORKDIR /app
COPY --from=builder /app/app .
COPY --from=builder /go/src/github.com/kevingil/blog/app/.env .
COPY --from=builder /go/src/github.com/kevingil/blog/app/templates ./templates 
COPY --from=builder /go/src/github.com/kevingil/blog/app/static ./static

# Run app
EXPOSE 8080
CMD ["./app"]
