# service-a/Dockerfile
FROM golang:1.23.5-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /service-a

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /service-a .
EXPOSE 8000
CMD ["./service-a"]