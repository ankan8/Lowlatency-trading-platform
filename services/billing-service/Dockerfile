FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY services/billing-service /app
RUN go build -o billing-service ./cmd

EXPOSE 50055
CMD ["./billing-service"]
