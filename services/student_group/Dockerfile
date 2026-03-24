FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o student_group

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/student_group .
COPY --from=builder /app/data ./data

EXPOSE 10004

CMD ["./student_group"]