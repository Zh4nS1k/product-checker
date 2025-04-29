FROM ubuntu:latest
LABEL authors="Zhanserik"

ENTRYPOINT ["top", "-b"]

# Dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o barcode-checker ./cmd/api

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/barcode-checker .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080
CMD ["./barcode-checker"]