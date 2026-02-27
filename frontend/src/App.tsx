import { useEffect, useState } from 'react'
import './index.css'

const API_BASE = ''

interface EventItem {
  title: string
  type: 'tv' | 'movie'
  start_time?: string
}

interface DayGroup {
  day_name: string
  date: string
  events: EventItem[]
}

export default function App() {
  const [days, setDays] = useState<DayGroup[]>([])
  const [pastEvents, setPastEvents] = useState<EventItem[]>([])
  const [totalTv, setTotalTv] = useState<number | string>('-')
  const [totalMovies, setTotalMovies] = useState<number | string>('-')
  const [schedule, setSchedule] = useState({ type: '-', nextRun: '-', timezone: '-' })
  const [loading, setLoading] = useState(true)
  const [range, setRange] = useState('1')
  const [triggering, setTriggering] = useState(false)
  const [triggerStatus, setTriggerStatus] = useState<string | null>(null)
  const [showPast, setShowPast] = useState(false)

  const fetchData = async () => {
    setLoading(true)
    try {
      const schedRes = await fetch(`${API_BASE}/api/schedule`)
      if (schedRes.ok) {
        const schedData = await schedRes.json()
        setSchedule({
          type: schedData.schedule_type || 'N/A',
          nextRun: schedData.next_run ? formatRelativeTime(schedData.next_run) : 'Not scheduled',
          timezone: schedData.timezone || 'UTC'
        })
      }

      const pastRes = await fetch(`${API_BASE}/api/past-events`)
      if (pastRes.ok) {
        const pastData = await pastRes.json()
        setPastEvents(pastData.events || [])
      }

      const eventsRes = await fetch(`${API_BASE}/api/events?days=${range}`)
      if (eventsRes.ok) {
        const eventsData = await eventsRes.json()
        setDays(eventsData.days || [])
        // Calculate splits
        let tvCount = 0;
        let movieCount = 0;
        (eventsData.days || []).forEach((d: any) => {
          (d.events || []).forEach((e: any) => {
            if (e.type === 'tv') tvCount++;
            if (e.type === 'movie') movieCount++;
          })
        })

        setTotalTv(tvCount)
        setTotalMovies(movieCount)
      }
    } catch (e) {
      console.error(e)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchData()
    const intv = setInterval(fetchData, 60000)
    return () => clearInterval(intv)
  }, [range])

  const handleTrigger = async () => {
    setTriggering(true)
    setTriggerStatus(null)
    try {
      const res = await fetch(`${API_BASE}/api/trigger`, { method: 'POST' })
      const data = await res.json()
      setTriggerStatus(data.message || 'Triggered successfully')
      setTimeout(fetchData, 2000)
    } catch (e) {
      setTriggerStatus('Failed to trigger')
    } finally {
      setTriggering(false)
    }
  }

  const formatRelativeTime = (dateString: string) => {
    const date = new Date(dateString)
    const diff = date.getTime() - new Date().getTime()
    if (diff < 0) return 'Overdue'
    const d = Math.floor(diff / 86400000)
    const h = Math.floor((diff % 86400000) / 3600000)
    const m = Math.floor((diff % 3600000) / 60000)
    if (d > 0) return `${d}d ${h}h ${m}m`
    if (h > 0) return `${h}h ${m}m`
    return `${m}m`
  }

  return (
    <>
      <header className="header">
        <div className="container">
          <div className="header-content">
            <div className="logo">
              <svg width="32" height="32" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
                <rect x="4" y="7" width="24" height="20" rx="2" stroke="currentColor" strokeWidth="2" />
                <line x1="4" y1="12" x2="28" y2="12" stroke="currentColor" strokeWidth="2" />
                <line x1="10" y1="4" x2="10" y2="7" stroke="currentColor" strokeWidth="2" strokeLinecap="round" />
                <line x1="22" y1="4" x2="22" y2="7" stroke="currentColor" strokeWidth="2" strokeLinecap="round" />
              </svg>
              <h1>Calendarr</h1>
            </div>
            <nav className="nav">
              <a href="#/" className="nav-link" style={{ background: 'var(--color-secondary)', boxShadow: '0 0 0 #000', transform: 'translate(2px, 2px)' }}>Dashboard</a>
              <a href="#/settings" className="nav-link">Settings</a>
            </nav>
          </div>
        </div>
      </header>

      <main className="main">
        <div className="container" id="dashboard">
          <div className="stats-grid">
            <div className="stat-card brutal-card brutal-interactive">
              <div className="stat-icon stat-icon-primary">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <rect x="3" y="6" width="18" height="15" rx="2" stroke="currentColor" strokeWidth="2" />
                  <line x1="3" y1="10" x2="21" y2="10" stroke="currentColor" strokeWidth="2" />
                </svg>
              </div>
              <div className="stat-content">
                <div className="stat-label">
                  {range === '1' ? 'TODAY' :
                    range === '3' ? 'NEXT 3 DAYS' :
                      range === '7' ? 'NEXT WEEK' :
                        range === '14' ? 'NEXT 2 WEEKS' : 'UPCOMING'}
                </div>
                <div className="stat-value" style={{ fontSize: '1.2rem' }}>
                  {totalTv !== '-' && totalMovies !== '-' ? (
                    <>
                      {Number(totalTv) > 0 ? `${totalTv} Show${Number(totalTv) > 1 ? 's' : ''}` : ''}
                      {Number(totalTv) > 0 && Number(totalMovies) > 0 ? ' and ' : ''}
                      {Number(totalMovies) > 0 ? `${totalMovies} Movie${Number(totalMovies) > 1 ? 's' : ''}` : ''}
                      {Number(totalTv) === 0 && Number(totalMovies) === 0 ? '0' : ''}
                    </>
                  ) : '-'}
                </div>
              </div>
            </div>

            <div className="stat-card brutal-card brutal-interactive">
              <div className="stat-icon stat-icon-success">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <circle cx="12" cy="12" r="9" stroke="currentColor" strokeWidth="2" />
                  <path d="M12 7V12L15 15" stroke="currentColor" strokeWidth="2" strokeLinecap="round" />
                </svg>
              </div>
              <div className="stat-content">
                <div className="stat-label">Next Run</div>
                <div className="stat-value">{schedule.nextRun}</div>
              </div>
            </div>

            <div className="stat-card brutal-card brutal-interactive">
              <div className="stat-icon stat-icon-info">
                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z" stroke="currentColor" strokeWidth="2" strokeLinejoin="round" />
                </svg>
              </div>
              <div className="stat-content">
                <div className="stat-label">Schedule Type</div>
                <div className="stat-value">{schedule.type}</div>
              </div>
            </div>

            <div className="stat-card brutal-card trigger-card">
              <button className="trigger-btn" onClick={handleTrigger} disabled={triggering}>
                {!triggering ? (
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M5 3l14 9-14 9V3z" fill="currentColor" />
                  </svg>
                ) : (
                  <div className="spinner" style={{ width: 16, height: 16, borderWidth: 2 }}></div>
                )}
                <span>{triggering ? "Running..." : "Run Now"}</span>
              </button>
              <div className={`trigger-status ${triggerStatus ? (triggerStatus.includes('success') ? 'success' : 'error') : ''}`}>
                {triggerStatus}
              </div>
            </div>
          </div>

          <section className="section" id="events">
            <div className="section-header">
              <h2>
                {range === '1' ? 'Releases Today' :
                  range === '3' ? 'Releases Next 3 Days' :
                    range === '7' ? 'Releases Next Week' :
                      range === '14' ? 'Releases Next 2 Weeks' : 'Upcoming Releases'}
              </h2>
              <div className="section-actions">
                <select className="range-selector" value={range} onChange={e => setRange(e.target.value)}>
                  <option value="1">Today</option>
                  <option value="3">Next 3 Days</option>
                  <option value="7">Next Week</option>
                  <option value="14">Next 2 Weeks</option>
                </select>
                <button className="refresh-btn" onClick={fetchData}>
                  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M21.5 2v6h-6M2.5 22v-6h6M2 11.5a10 10 0 0 1 18.8-4.3M22 12.5a10 10 0 0 1-18.8 4.2" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
                  </svg>
                  Refresh
                </button>
              </div>
            </div>

            <div className="events-container">
              {loading ? (
                <div className="loading">
                  <div className="spinner"></div>
                  <p>Loading releases...</p>
                </div>
              ) : days.length === 0 ? (
                <div className="empty-state">
                  <svg width="48" height="48" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <rect x="3" y="6" width="18" height="15" rx="2" stroke="currentColor" strokeWidth="2" />
                    <line x1="3" y1="10" x2="21" y2="10" stroke="currentColor" strokeWidth="2" />
                  </svg>
                  <p>No upcoming events found</p>
                </div>
              ) : (
                days.map((day, i) => (
                  <div key={i} className="day-group">
                    <div className="day-header">
                      {day.day_name}
                      <span className="day-date">{day.date}</span>
                    </div>
                    <div className="events-list">
                      {day.events.map((ev, j) => (
                        <div key={j} className={`event-card brutal-card event-${ev.type}`}>
                          <span className="event-type">{ev.type === 'tv' ? 'TV' : 'Movie'}</span>
                          <div className="event-title">{ev.title}</div>
                          {ev.start_time && (
                            <div className="event-time" data-original-time={ev.start_time}>
                              <span className="time-text">{ev.start_time}</span>
                            </div>
                          )}
                        </div>
                      ))}
                    </div>
                  </div>
                ))
              )}
            </div>
          </section>

          <section className="section" id="past-events">
            <div className="section-header">
              <h2>Past Releases</h2>
              <button className="toggle-btn" onClick={() => setShowPast((p) => !p)}>
                <svg className="toggle-icon" style={{ transform: showPast ? 'rotate(180deg)' : 'rotate(0deg)' }} width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M19 9l-7 7-7-7" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
                </svg>
              </button>
            </div>

            {showPast && (
              <div className="past-content">
                <div className="events-container">
                  {pastEvents.length === 0 ? (
                    <div className="empty-state">
                      <svg width="48" height="48" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                        <rect x="3" y="6" width="18" height="15" rx="2" stroke="currentColor" strokeWidth="2" />
                        <line x1="3" y1="10" x2="21" y2="10" stroke="currentColor" strokeWidth="2" />
                      </svg>
                      <p>No past events today</p>
                    </div>
                  ) : (
                    <div className="day-group">
                      <div className="events-list">
                        {pastEvents.map((ev, j) => (
                          <div key={j} className={`event-card brutal-card event-${ev.type} event-past`}>
                            <span className="event-type">{ev.type === 'tv' ? 'TV' : 'Movie'}</span>
                            <div className="event-title">{ev.title}</div>
                            {ev.start_time && (
                              <div className="event-time">
                                <span className="time-text">{ev.start_time}</span>
                              </div>
                            )}
                          </div>
                        ))}
                      </div>
                    </div>
                  )}
                </div>
              </div>
            )}
          </section>
        </div>
      </main>

      <footer className="footer">
        <div className="container">
          <p><a href="https://github.com/khw315/calendarr" target="_blank" rel="noreferrer">GitHub</a></p>
        </div>
      </footer>
    </>
  )
}
