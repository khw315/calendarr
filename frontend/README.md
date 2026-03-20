# Calendarr Web UI

This directory contains the frontend Web UI for Calendarr.

It is a modern, responsive dashboard built with **React**, **TypeScript**, **TailwindCSS**, and **Vite**. The UI provides a beautiful way to visualize upcoming TV and movie releases, trigger manual calendar syncs, and configure all application settings dynamically without restarting the Docker container.

## Features

- **Calendar Dashboard**: Visual representation of upcoming releases grouped by day.
- **Settings Manager**: Full control over integrations (Discord/Slack), format settings, iCal URLs, and scheduling.
- **Manual Trigger**: Quick button to manually run the sync job.
- **Auto-Refresh**: Dashboard updates automatically every 60 seconds (configurable).

## Development Setup

### Prerequisites

- Node.js (v18+ recommended)
- npm or yarn

### Getting Started

1. **Install dependencies**:
   ```bash
   npm install
   ```

2. **Start the development server**:
   ```bash
   npm run dev
   ```
   The development server will start at `http://localhost:5173`. 

3. **Backend API Connection**:
   During development, the frontend expects the Python backend to be running simultaneously (usually on port 5000). The Vite development server automatically proxies API requests (`/api/*`) to the backend. Ensure the backend is running to load actual calendar data.

## Building for Production

To build the frontend for production:

```bash
npm run build
```

This will compile and minify the assets into the `dist/` directory, which the Flask backend will serve statically when running in the Docker container.
