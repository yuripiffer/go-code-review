# build stage
FROM golang:1.24.2-alpine3.21 AS builder

RUN apk --no-cache add git gcc libc-dev make

WORKDIR /app

COPY go.* .

RUN go mod download

COPY . .

# Compiles the binary for linux with no dependencies on C libraries
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/coupon_service/

FROM alpine:3.21

WORKDIR /app

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --chown=appuser:appgroup --from=builder /app/main /app/

USER appuser

EXPOSE 8080

ENTRYPOINT ["/app/main"]
