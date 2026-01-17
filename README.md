![Docker Pulls](https://img.shields.io/docker/pulls/hann315/calendarr)
![GitHub Release](https://img.shields.io/github/v/release/khw315/calendarr)
![GitHub last commit](https://img.shields.io/github/last-commit/khw315/calendarr)

# 📆 Calendarr

A simple Docker container that fetches upcoming airings/releases for TV shows and movies from Sonarr and Radarr calendars and posts them to Discord on a schedule.

![Example Discord post](https://i.imgur.com/abGhhg4.png)

## ✨ Features

- **Consolidated Feed**: Combines multiple Sonarr and Radarr calendar feeds into one summary.
- **Smart Grouping**: Groups shows and movies by day of the week for easy reading.
- **Flexible Scheduling**: Runs on a customizable schedule (Daily or Weekly) or Cron expression.
- **Multi-Platform**: Supports both Discord and Slack notifications.
- **Localization**: Native support for English (EN), Korean (KO), Japanese (JA), and Indonesian (ID).
- **Dynamic Timezones**: Automatically adapts to your configured timezone.
- **Highly Customizable**: Configure headers, footers, timestamp styles, and more.

## 🚀 Getting Started

Available via `ghcr.io/khw315/calendarr:latest`.

### 🐳 With Docker Compose (Recommended)

1.  **Create a `.env` file** with your configuration (see [Configuration](#-configuration) below).
2.  **Create a `docker-compose.yml`**:

```yaml
services:
  calendarr:
    image: ghcr.io/khw315/calendarr:latest
    container_name: calendarr
    restart: unless-stopped
    volumes:
      - ./logs:/app/logs:rw
    env_file: .env
    environment:
      - TZ=Asia/Seoul
```

3.  **Run it**:
    ```bash
    docker compose up -d
    ```

### ⌨️ With Docker CLI

```bash
docker run -d \
  --name calendarr \
  -e DISCORD_WEBHOOK_URL="https://discord.com/api/webhooks/..." \
  -e CALENDAR_URLS='[{"url":"https://sonarr.example.com/feed/calendar/api.ics","type":"tv"},{"url":"https://radarr.example.com/feed/calendar/api.ics","type":"movie"}]' \
  -e TZ="Asia/Seoul" \
  ghcr.io/khw315/calendarr:latest
```

## 🛠️ Configuration

Configure the application using environment variables in your `.env` file or Docker compose.

| Variable | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `CALENDAR_URLS` * | JSON | `[]` | JSON array of calendar URLs and types. Example: `[{"url":"http://...","type":"tv"}, {"url":"http://...","type":"movie"}]` |
| `DISCORD_WEBHOOK_URL` ** | String | `""` | Discord webhook URL. |
| `SLACK_WEBHOOK_URL` *** | String | `""` | Slack webhook URL. |
| `TZ` | String | `UTC` | Timezone (e.g., `America/New_York`). |
| `APP_LANGUAGE` | String | `EN` | Language for generated messages. Options: `EN`, `KO`, `JA`, `ID`. |
| `SCHEDULE_TYPE` | String | `WEEKLY` | `DAILY` or `WEEKLY`. |
| `RUN_TIME` | Time | `09:00` | Time to run the job (HH:MM). |
| `DISCORD_TIMESTAMP_STYLE` | String | `Relative Time` | *Discord only* Style for timestamps. Options: `Short Time`, `Long Time`, `Short Date`, `Long Date`, `Short/Long Date/Time`, `Relative Time`. |
| `CUSTOM_HEADER` | String | `New Releases` | Custom header text for the notification. |
| `DISPLAY_TIME` | Boolean | `true` | Display the release time next to events. |
| `SHOW_DATE_RANGE` | Boolean | `true` | Show the date range in the header. |
| `SHOW_TIMEZONE_IN_SUBHEADER`| Boolean | `false`| Show the configured timezone in the subheader. |
| `USE_24_HOUR` | Boolean | `true` | Use 24-hour time format. |
| `ADD_LEADING_ZERO` | Boolean | `true` | Add leading zero to single-digit hours. |
| `DEDUPLICATE_EVENTS` | Boolean | `true` | Remove duplicate events from multiple sources. |
| `PASSED_EVENT_HANDLING` | String | `DISPLAY`| How to handle past events: `DISPLAY`, `HIDE`, `STRIKE`. |
| `CRON_SCHEDULE` | Cron | `None` | Custom CRON expression (Overrides simple scheduling). |
| `ENABLE_CUSTOM_DISCORD_FOOTER`| Boolean| `false`| Enable custom footer for Discord messages. |
| `ENABLE_CUSTOM_SLACK_FOOTER` | Boolean | `false`| Enable custom footer for Slack messages. |
| `DEBUG` | Boolean | `false` | Enable debug logging. |

\* Required.
\** Required if `USE_DISCORD` is `true` (default).
\*** Required if `USE_SLACK` is `true`.

## 🤝 Obtaining Calendar URLs

### Sonarr / Radarr
1.  Go to **Settings** > **General**.
2.  Copy your **API Key**.
3.  Construct your URL:
    *   **Sonarr**: `http://YOUR_HOST:8989/feed/v3/calendar/Sonarr.ics?apikey=YOUR_API_KEY`
    *   **Radarr**: `http://YOUR_HOST:7878/feed/v3/calendar/Radarr.ics?apikey=YOUR_API_KEY`

*Alternatively, use the "Calendar > iCal Link" button in the UI (ensure no tags/filters are selected).*

## ✍️ Custom Footers

You can inject custom markdown into the footer of your messages.

1.  Set `ENABLE_CUSTOM_DISCORD_FOOTER=true` (or Slack equivalent).
2.  Mount a volume to `/app/custom_footers`:
    ```yaml
    volumes:
      - ./custom_footers:/app/custom_footers:rw
    ```
3.  Edit the generated `discord_footer.md` or `slack_footer.md` in that directory.

## 🚧 Development

To build the container locally:

```bash
git clone https://github.com/khw315/calendarr.git
cd calendarr
docker build -t calendarr .
```

## 🧑‍💻 Contributing

Contributions are welcome!
- **Localization**: Help us translate Calendarr into more languages by adding to `src/data/locales.json`.
- **Features**: Submit PRs for new platform integrations or improvements.

## 🧑‍⚖️ License

GNU GENERAL PUBLIC LICENSE
