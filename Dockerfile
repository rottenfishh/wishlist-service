FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/swaggo/swag/cmd/swag@v1.16.4

COPY . .

RUN /go/bin/swag init -g ./cmd/main.go -o ./docs
RUN go build -o main ./cmd

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/database ./database

EXPOSE 8080

CMD ["./main"]