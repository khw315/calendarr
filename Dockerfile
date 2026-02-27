# Build frontend
FROM node:20-slim AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package.json frontend/package-lock.json* ./
RUN npm ci || npm install
COPY frontend/ ./
RUN npm run build

# Build backend
FROM python:3.14.2-slim

WORKDIR /app

LABEL org.opencontainers.image.source=https://github.com/khw315/calendarr
LABEL org.opencontainers.image.description="Calendar feeds from Sonarr/Radarr to Discord and Slack"
LABEL org.opencontainers.image.licenses=GPL-3.0

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application source code
COPY src/ /app/src/

# Copy public directory for web UI from frontend builder
COPY --from=frontend-builder /app/public/ /app/public/

# Copy and set up the entrypoint script
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

# Create logs directory
RUN mkdir -p /app/logs

ENV PYTHONUNBUFFERED=1

EXPOSE 5000

# Set the entrypoint script
ENTRYPOINT ["/app/entrypoint.sh"]

# Set the default command (will be executed by entrypoint's exec "$@")
CMD ["python", "src/app.py"]