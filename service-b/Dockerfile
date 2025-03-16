# service-b/Dockerfile
FROM golang:1.23.5-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /service-b

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /service-b .
EXPOSE 8080
CMD ["./service-b"]