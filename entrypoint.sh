#!/bin/sh
set -e

# Fix permissions on host-mounted volume directories at container startup
if [ "$(id -u)" = '0' ]; then
    chown -R calendarr:calendarr /app/config /app/logs 2>/dev/null || chmod -R 777 /app/config /app/logs 2>/dev/null || true
    exec su-exec calendarr "$@"
fi

exec "$@"
