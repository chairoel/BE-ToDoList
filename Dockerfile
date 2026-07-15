# Step 1: build
FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/api ./cmd/api

# Step 2: runtime
FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata \
	&& adduser -D -g '' appuser

COPY --from=builder /app/bin/api /app/api

USER appuser

EXPOSE 8080

CMD ["/app/api"]
