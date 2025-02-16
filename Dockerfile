FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o avito-shop ./cmd/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/avito-shop .
EXPOSE 8080
CMD ["./avito-shop"]
