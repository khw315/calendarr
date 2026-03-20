> [!NOTE]
> This repository is a fork of [jordanlambrecht/calendarr](https://github.com/jordanlambrecht/calendarr)

# Calendarr

A simple, lightweight Docker container that fetches upcoming airings and releases for TV shows and movies directly from your **Sonarr** and **Radarr** calendars, and seamlessly posts them to **Discord** or **Slack** on a schedule.

![Example Discord post](screenshots/discord.png)

## Features

- **Consolidated Feed**: Combines multiple Sonarr and Radarr calendar iCal feeds into one beautiful, easy-to-read summary.
- **Modern Web UI**: A bold, Neobrutalist dashboard to view upcoming events, manually trigger runs, and configure your settings on the fly.
- **Smart Grouping**: Groups upcoming shows and movies intelligently by day of the week to maximize readability and reduce notification fatigue.
- **Flexible Scheduling**: Runs exactly when you want it to. Choose a Daily schedule, a Weekly summary, or provide your own robust Cron expression.
- **Multi-Platform Integrations**: Effortlessly supports Webhooks for both Discord and Slack out-of-the-box.
- **Localization Support**: Native translation support for English (EN), Korean (KO), Japanese (JA), and Indonesian (ID).
- **Dynamic Timezones**: Automatically adapts all outputs and schedules to your explicitly configured timezone without any rigid system dependencies.
- **Extremely Customizable**: Configure headers, footers, timestamp styles, and custom role mentions dynamically.

---

## Web UI Dashboard

Calendarr features a modern Web UI for visualization and configuration without forcing you to restart the container or edit raw files.

### Accessing the Dashboard
Once the container is spun up, access the Web UI at:
```text
http://localhost:5000
```
*(Or replace `localhost` with your remote server IP).*

### UI Capabilities
- **Dashboard Hub**: View analytics including upcoming events count, the next scheduled cron run, and your schedule type.
- **Live Calendar View**: See all upcoming releases cleanly grouped by day, complete with color-coding for TV versus Movies.
- **Manual Automation Trigger**: Force-run the synchronization job instantly with a single button click.
- **Visual Configuration Viewer**: Edit your settings, URLs, scheduling, and timezone dynamically.
- **Auto-Refresh Data**: The event dashboard automatically syncs and refreshes dynamically.

![Web UI Screenshot](screenshots/preview.png)

---

## Getting Started

The official image is available via `ghcr.io/khw315/calendarr:latest`

### Option 1: Docker Compose (Recommended)

1. **Create your `docker-compose.yml`**:
```yaml
services:
  calendarr:
    image: ghcr.io/khw315/calendarr:latest
    container_name: calendarr
    restart: unless-stopped
    ports:
      - "5000:5000"  # Web UI Access
    volumes:
      - ./calendarr/config:/app/config:rw
      - ./calendarr/logs:/app/logs:rw  # Optional but recommended
```

2. **Spin it up**:
```bash
docker compose up -d
```

3. **Configure the App**: Navigate to `http://localhost:5000` to setup your webhook integrations, schedules, and your calendar URLs!

### Option 2: Docker CLI

```bash
docker run -d \
  --name calendarr \
  -p 5000:5000 \
  -v ./calendarr/config:/app/config:rw \
  -v ./calendarr/logs:/app/logs:rw \
  ghcr.io/khw315/calendarr:latest
```

### Manual Trigger (Off-schedule)
If you ever want to bypass the internal scheduler purely through standard CLI commands:
1. Hit the "Run Now" button on the Web UI.
2. POST a request to `http://localhost:5000/api/trigger` with command:
   ```bash
   curl -X POST http://localhost:5000/api/trigger
   ```
3. Execute the primary thread manually:
   `docker exec -it calendarr python src/main.py`

---

## App Configuration

Configuration is managed directly from the **Web UI**. All changes will seamlessly apply instantly without restarting the orchestrator, and are persisted mapped into your `/app/config` volume automatically.

| Settings Panel | Description |
| :--- | :--- |
| **Calendars** | Add endless Sonarr or Radarr `iCal` URLs and declare their media target (`tv` or `movie`). |
| **Integrations** | Enable/Disable Discord and Slack delivery triggers. Define custom Webhook URLs and role mentions. |
| **Format** | Adjust your primary locale translation engine, header verbosity, and dynamic footer injections. |
| **Event Display** | Define logic rules for duplicate event resolutions, toggle 24-hour timestamp clocks, and TZ metadata. |
| **Schedule** | Instruct internal Chron schedulers for Daily logic, Weekly recaps, absolute Run Times, or Raw Cron inputs. |
| **Advanced** | Developer settings for logging velocity controls, max file chunking, debugging, and strict HTTP boundaries. |

---

## Syncing Calendars (Sonarr/Radarr)

Calendarr processes external standard `iCal` (.ics) endpoints. To fetch these:
1. In your media manager (Sonarr or Radarr), navigate to **Settings** > **General**.
2. Copy your static **API Key**.
3. Construct the endpoint URLs manually logic:
   - **Sonarr**: `http://YOUR_HOST:8989/feed/v3/calendar/Sonarr.ics?unmonitored=true&apikey=YOUR_API_KEY`
   - **Radarr**: `http://YOUR_HOST:7878/feed/v3/calendar/Radarr.ics?unmonitored=true&apikey=YOUR_API_KEY`

*Pro-tip: You can quickly utilize the **Calendar > iCal Link** dashboard UI button natively in Sonarr/Radarr.*

---

## Local Development

### Prerequisites & Tech Stack
- **Backend Framework**: Python (Flask)
- **Frontend Framework**: React + TypeScript (Vite + TailwindCSS)
- **Tooling Engine**: [`uv`](https://github.com/astral-sh/uv) (Python Packager)

### 1. Repository Setup

Clone the source matrix:
```bash
git clone https://github.com/khw315/calendarr.git
cd calendarr
```

### 2. Starting the Backend Server
Initialize the Python stack. `uv` will safely scaffold and run standard execution logic directly without boilerplate configurations:
```bash
uv run src/app.py
```
*(The backend logic API sits at `localhost:5000`)*

### 3. Starting the Frontend UI
In a separate terminal, deploy the hot-reload Vite server:
```bash
cd frontend
npm install
npm run dev
```
*(Development GUI is served via `localhost:5173` with proxy linking API queries down to the backend `5000` port).*

### 4. Running Backend Tests
Calendarr holds a robust internal test suite managed via `pytest`. Running the test suite ensures configurations (like translations) natively conform securely to runtime patterns:
```bash
uv pip install -r requirements-test.txt
uv run pytest tests/
```

### 5. Docker Builds
If you desire testing raw container builds directly inside an isolated space:
```bash
docker build -t calendarr .
```

---

## Contributing

Continuous Open-Source participation is heavily valued!
- **Features**: Submit PRs containing architectural backend integrations, frontend tweaks, or general enhancements.
- **Localization**: Improve or append dictionaries directly toward `src/data/locales/`. The JSON system inherently cascades template mappings into global text hooks globally.

## License

[GNU General Public License v3.0](LICENSE)
