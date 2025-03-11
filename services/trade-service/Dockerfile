FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY services/trade-service /app
RUN go build -o trade-service ./cmd

EXPOSE 50053
CMD ["./trade-service"]
