# Build stage
FROM golang:1.21-alpine

WORKDIR /go/src/auth-service

# Copy the source code
COPY . .

# Download dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

EXPOSE 8080
EXPOSE 50051

CMD ["./main"] 