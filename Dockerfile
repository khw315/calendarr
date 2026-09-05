# Stage 1: Build Frontend Web UI (Native Build Platform)
FROM --platform="$BUILDPLATFORM" oven/bun:1-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package.json frontend/bun.lock* ./
RUN bun install --frozen-lockfile
COPY frontend/ ./
RUN bun run build

# Stage 2: Build Go Static Binary (Native Build Platform + Go Cross Compilation)
FROM --platform="$BUILDPLATFORM" golang:1.24-alpine AS go-builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend-builder /app/public /app/public
RUN CGO_ENABLED=0 GOOS="${TARGETOS:-linux}" GOARCH="${TARGETARCH:-amd64}" go build -ldflags="-w -s" -o calendarr main.go

# Stage 3: Minimal Final Image
FROM alpine:3.21

WORKDIR /app

LABEL org.opencontainers.image.title="Calendarr"
LABEL org.opencontainers.image.description="Calendar feeds from Sonarr/Radarr to Discord and Slack"
LABEL org.opencontainers.image.source="https://github.com/khw315/calendarr"
LABEL org.opencontainers.image.documentation="https://github.com/khw315/calendarr/blob/main/README.md"
LABEL org.opencontainers.image.licenses="GPL-3.0"

RUN apk add --no-cache ca-certificates su-exec tzdata && \
    addgroup -g 1000 -S calendarr && \
    adduser -u 1000 -S calendarr -G calendarr

COPY --from=go-builder /app/calendarr /app/calendarr
COPY entrypoint.sh /app/entrypoint.sh

# Fix CRLF line endings for entrypoint.sh on Windows hosts and set permissions
RUN tr -d '\r' < /app/entrypoint.sh > /tmp/entrypoint.sh && \
    mv /tmp/entrypoint.sh /app/entrypoint.sh && \
    chmod +x /app/entrypoint.sh && \
    mkdir -p /app/config /app/logs && \
    chown -R calendarr:calendarr /app && \
    chmod -R 775 /app/config /app/logs

EXPOSE 5000

ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["/app/calendarr"]
