// API Base URL
const API_BASE = '';

// State
let configOpen = false;
let autoRefreshInterval = null;

// Initialize app on DOM load
document.addEventListener('DOMContentLoaded', () => {
    initializeApp();
    setupEventListeners();
    startAutoRefresh();
});

// Initialize application
async function initializeApp() {
    await Promise.all([
        loadEvents(),
        loadTimezone(),
        loadScheduleInfo()
    ]);
}

// Setup event listeners
function setupEventListeners() {
    // Trigger button
    const triggerBtn = document.getElementById('triggerBtn');
    triggerBtn.addEventListener('click', handleTrigger);

    // Refresh button
    const refreshBtn = document.getElementById('refreshBtn');
    refreshBtn.addEventListener('click', () => {
        loadEvents();
        loadScheduleInfo();
    });

    // Range selector
    const rangeSelector = document.getElementById('rangeSelector');
    rangeSelector.addEventListener('change', () => {
        loadEvents();
    });

    // Past events toggle
    const pastToggle = document.getElementById('pastToggle');
    pastToggle.addEventListener('click', togglePastEvents);

    // Navigation
    document.querySelectorAll('.nav-link').forEach(link => {
        link.addEventListener('click', (e) => {
            e.preventDefault();
            const target = link.getAttribute('href').substring(1);
            scrollToSection(target);
        });
    });
}

// Scroll to section
function scrollToSection(sectionId) {
    const section = document.getElementById(sectionId);
    if (section) {
        section.scrollIntoView({ behavior: 'smooth', block: 'start' });
    }
}

// Load events
async function loadEvents() {
    const container = document.getElementById('eventsContainer');

    try {
        showLoading(container);

        // Get selected range from dropdown
        const rangeSelector = document.getElementById('rangeSelector');
        const days = rangeSelector.value;

        // Update the events label based on selected range
        updateEventsLabel(days);

        const response = await fetch(`${API_BASE}/api/events?days=${days}`);
        if (!response.ok) throw new Error('Failed to fetch events');

        const data = await response.json();

        // Get past events count to subtract from total
        let pastCount = 0;
        try {
            const pastResponse = await fetch(`${API_BASE}/api/past-events`);
            if (pastResponse.ok) {
                const pastData = await pastResponse.json();
                pastCount = pastData.count || 0;
            }
        } catch (error) {
            console.error('Error fetching past events count:', error);
        }

        // Update total events stat (subtract past events)
        const totalEvents = (data.total_events || 0) - pastCount;
        document.getElementById('totalEvents').textContent = Math.max(0, totalEvents);

        // Render events
        if (data.days && data.days.length > 0) {
            renderEvents(container, data.days);
        } else {
            showEmptyState(container, 'No upcoming events found');
        }
    } catch (error) {
        console.error('Error loading events:', error);
        showError(container, 'Failed to load events. Please try again.');
    }
}

// Update events label based on selected range
function updateEventsLabel(days) {
    const label = document.getElementById('eventsLabel');
    const labelMap = {
        '1': 'Today Releases',
        '3': 'Next 3 Days',
        '7': 'This Week',
        '14': 'Next 2 Weeks'
    };

    label.textContent = labelMap[days] || 'Upcoming Events';
}

// Render events
function renderEvents(container, days) {
    container.innerHTML = '';

    days.forEach(day => {
        const dayGroup = document.createElement('div');
        dayGroup.className = 'day-group';

        const dayHeader = document.createElement('div');
        dayHeader.className = 'day-header';
        dayHeader.innerHTML = `
            ${day.day_name}
            <span class="day-date">${day.date}</span>
        `;

        const eventsList = document.createElement('div');
        eventsList.className = 'events-list';

        day.events.forEach(event => {
            const eventCard = createEventCard(event);
            eventsList.appendChild(eventCard);
        });

        dayGroup.appendChild(dayHeader);
        dayGroup.appendChild(eventsList);
        container.appendChild(dayGroup);
    });
}

// Create event card
function createEventCard(event) {
    const card = document.createElement('div');
    const pastClass = event.is_past ? ' event-past' : '';
    card.className = `event-card glass-card event-${event.type}${pastClass}`;

    card.innerHTML = `
        <span class="event-type">${event.type}</span>
        <div class="event-title">${event.title}</div>
        ${event.time ? `
            <div class="event-time">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <circle cx="12" cy="12" r="9" stroke="currentColor" stroke-width="2"/>
                    <path d="M12 7V12L15 15" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                </svg>
                ${event.time}
            </div>
        ` : ''}
    `;

    return card;
}

// Load timezone
async function loadTimezone() {
    try {
        const response = await fetch(`${API_BASE}/api/schedule`);
        if (!response.ok) throw new Error('Failed to fetch timezone');

        const data = await response.json();
        document.getElementById('timezone').textContent = data.timezone || 'UTC';
    } catch (error) {
        console.error('Error loading timezone:', error);
        document.getElementById('timezone').textContent = 'Error';
    }
}

// Load past events
async function loadPastEvents() {
    const container = document.getElementById('pastEventsContainer');

    try {
        showLoading(container);

        const response = await fetch(`${API_BASE}/api/past-events`);
        if (!response.ok) throw new Error('Failed to fetch past events');

        const data = await response.json();

        // Render past events
        if (data.events && data.events.length > 0) {
            container.innerHTML = '';
            const dayDiv = document.createElement('div');
            dayDiv.className = 'day-group';

            const eventsGrid = document.createElement('div');
            eventsGrid.className = 'events-list';

            data.events.forEach(event => {
                const card = createEventCard(event);
                card.classList.add('event-past');
                eventsGrid.appendChild(card);
            });

            dayDiv.appendChild(eventsGrid);
            container.appendChild(dayDiv);
        } else {
            showEmptyState(container, 'No past events today');
        }
    } catch (error) {
        console.error('Error loading past events:', error);
        showError(container, 'Failed to load past events');
    }
}

// Toggle past events visibility
function togglePastEvents() {
    const content = document.getElementById('pastContent');
    const button = document.getElementById('pastToggle');
    const icon = button.querySelector('.toggle-icon');

    if (content.style.display === 'none' || content.style.display === '') { // Check for empty string too
        content.style.display = 'block';
        icon.style.transform = 'rotate(180deg)';
        // Load past events when opening for the first time
        if (!content.dataset.loaded) {
            loadPastEvents();
            content.dataset.loaded = 'true';
        }
    } else {
        content.style.display = 'none';
        icon.style.transform = 'rotate(0deg)';
    }
}

// Load schedule information
async function loadScheduleInfo() {
    try {
        const response = await fetch(`${API_BASE}/api/schedule`);
        if (!response.ok) throw new Error('Failed to fetch schedule');

        const data = await response.json();

        // Update schedule type
        document.getElementById('scheduleType').textContent = data.schedule_type || 'N/A';

        // Update next run time
        const nextRunElement = document.getElementById('nextRun');
        if (data.next_run) {
            nextRunElement.textContent = formatRelativeTime(data.next_run);
        } else {
            nextRunElement.textContent = 'Not scheduled';
        }
    } catch (error) {
        console.error('Error loading schedule:', error);
        document.getElementById('scheduleType').textContent = 'Error';
        document.getElementById('nextRun').textContent = 'Error';
    }
}

// Handle manual trigger
async function handleTrigger() {
    const btn = document.getElementById('triggerBtn');
    const status = document.getElementById('triggerStatus');

    btn.disabled = true;
    btn.innerHTML = `
        <div class="spinner" style="width: 16px; height: 16px; border-width: 2px;"></div>
        <span>Running...</span>
    `;
    status.textContent = '';
    status.className = 'trigger-status';

    try {
        const response = await fetch(`${API_BASE}/api/trigger`, {
            method: 'POST'
        });

        if (!response.ok) throw new Error('Trigger failed');

        const result = await response.json();

        status.textContent = result.message || 'Job triggered successfully!';
        status.className = 'trigger-status success';

        // Reload events and schedule after a short delay
        setTimeout(() => {
            loadEvents();
            loadScheduleInfo();
        }, 2000);

    } catch (error) {
        console.error('Error triggering job:', error);
        status.textContent = 'Failed to trigger job. Please try again.';
        status.className = 'trigger-status error';
    } finally {
        btn.disabled = false;
        btn.innerHTML = `
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M5 3l14 9-14 9V3z" fill="currentColor"/>
            </svg>
            <span>Run Now</span>
        `;
    }
}

// Auto-refresh functionality
function startAutoRefresh() {
    // Refresh data every 60 seconds
    autoRefreshInterval = setInterval(() => {
        loadEvents();
        loadScheduleInfo();
    }, 60000);
}

// Utility: Show loading state
function showLoading(container) {
    container.innerHTML = `
        <div class="loading">
            <div class="spinner"></div>
            <p>Loading...</p>
        </div>
    `;
}

// Utility: Show error state
function showError(container, message) {
    container.innerHTML = `
        <div class="empty-state">
            <svg width="48" height="48" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" style="color: var(--color-danger);">
                <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/>
                <line x1="12" y1="8" x2="12" y2="12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                <circle cx="12" cy="16" r="1" fill="currentColor"/>
            </svg>
            <p>${message}</p>
        </div>
    `;
}

// Utility: Show empty state
function showEmptyState(container, message) {
    container.innerHTML = `
        <div class="empty-state">
            <svg width="48" height="48" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <rect x="3" y="6" width="18" height="15" rx="2" stroke="currentColor" stroke-width="2"/>
                <line x1="3" y1="10" x2="21" y2="10" stroke="currentColor" stroke-width="2"/>
            </svg>
            <p>${message}</p>
        </div>
    `;
}

// Utility: Format relative time
function formatRelativeTime(dateString) {
    const date = new Date(dateString);
    const now = new Date();
    const diff = date - now;

    if (diff < 0) return 'Overdue';

    const minutes = Math.floor(diff / 60000);
    const hours = Math.floor(minutes / 60);
    const days = Math.floor(hours / 24);

    if (days > 0) return `${days}d ${hours % 24}h`;
    if (hours > 0) return `${hours}h ${minutes % 60}m`;
    return `${minutes}m`;
}
