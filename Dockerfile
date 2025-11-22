# ----------------------
# 1. Build stage
# ----------------------
FROM golang:1.24-alpine AS build
RUN apk add --no-cache git

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY main.go ./
COPY cmd/ ./cmd/
COPY internal/ ./internal/


RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# ----------------------
# 2. Runtime stage
# ----------------------
FROM alpine:3.20

# Add non-root user
RUN adduser -D appuser

WORKDIR /app

COPY --from=build /app/app .

USER appuser

# App will read env vars automatically
CMD ["./app"]
