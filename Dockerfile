# Stage 1: Build Frontend Web UI (Native Build Platform)
FROM --platform="$BUILDPLATFORM" node:20-slim AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci --ignore-scripts
COPY frontend/ ./
RUN npm run build

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

RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S calendarr && adduser -S calendarr -G calendarr

COPY --from=go-builder /app/calendarr /app/calendarr

# Create config and logs directories and set permissions for mounted volumes
RUN mkdir -p /app/config /app/logs && \
    chown -R calendarr:calendarr /app && \
    chmod -R 777 /app/config /app/logs

USER calendarr

EXPOSE 5000

CMD ["/app/calendarr"]
