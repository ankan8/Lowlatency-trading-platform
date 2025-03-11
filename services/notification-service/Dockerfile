FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY services/notification-service /app
RUN go build -o notification-service ./cmd

EXPOSE 50056
CMD ["./notification-service"]
