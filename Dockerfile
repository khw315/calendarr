# Stage 1: Build Frontend Web UI
FROM node:20-slim AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package.json frontend/package-lock.json* ./
RUN npm ci || npm install
COPY frontend/ ./
RUN npm run build

# Stage 2: Build Go Static Binary
FROM golang:1.24-alpine AS go-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend-builder /app/public /app/public
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o calendarr main.go

# Stage 3: Minimal Final Image
FROM alpine:3.21

WORKDIR /app

LABEL org.opencontainers.image.source=https://github.com/khw315/calendarr
LABEL org.opencontainers.image.description="Calendar feeds from Sonarr/Radarr to Discord and Slack"
LABEL org.opencontainers.image.licenses=GPL-3.0

RUN apk add --no-cache ca-certificates tzdata

COPY --from=go-builder /app/calendarr /app/calendarr

# Create config and logs directories
RUN mkdir -p /app/config /app/logs

EXPOSE 5000

CMD ["/app/calendarr"]