FROM golang:1.21-alpine

WORKDIR /app

# 1) Copy go.mod and go.sum from the root build context
COPY go.mod go.sum ./
RUN go mod download

# 2) Copy the auth-service code into /app
COPY services/auth-service /app

# 3) Build the binary from cmd/main.go
RUN go build -o auth-service ./cmd

# 4) Expose the gRPC port (e.g., 50051)
EXPOSE 50051

# 5) Run the compiled binary
CMD ["./auth-service"]
