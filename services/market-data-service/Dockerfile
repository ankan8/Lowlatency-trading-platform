FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY services/market-data-service /app
RUN go build -o market-data-service ./cmd

EXPOSE 50054
CMD ["./market-data-service"]
