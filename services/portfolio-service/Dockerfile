FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY services/portfolio-service /app
RUN go build -o portfolio-service ./cmd

EXPOSE 50052
CMD ["./portfolio-service"]
