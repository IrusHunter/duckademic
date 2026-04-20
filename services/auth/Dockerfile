FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o auth

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/auth .
COPY --from=builder /app/data ./data

EXPOSE 10009

CMD ["./auth"]