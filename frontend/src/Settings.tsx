import { useState, useEffect } from 'react'
import { useTheme } from './lib/useTheme'

const API_BASE = ''

// Reuse common type
export default function Settings() {
    const [config, setConfig] = useState<any>({
        APP_LANGUAGE: "EN",
        USE_DISCORD: false,
        DISCORD_WEBHOOK_URL: "",
        DISCORD_MENTION_ROLE_ID: "",
        DISCORD_HIDE_MENTION_INSTRUCTIONS: false,
        DISCORD_TIMESTAMP_STYLE: "R",
        ENABLE_CUSTOM_DISCORD_FOOTER: false,
        USE_SLACK: false,
        SLACK_WEBHOOK_URL: "",
        ENABLE_CUSTOM_SLACK_FOOTER: false,
        CALENDAR_URLS: [],
        PASSED_EVENT_HANDLING: "DISPLAY",
        DEDUPLICATE_EVENTS: false,
        USE_24_HOUR: false,
        ADD_LEADING_ZERO: false,
        DISPLAY_TIME: true,
        SHOW_DATE_RANGE: false,
        SHOW_TIMEZONE_IN_SUBHEADER: false,
        TZ: "UTC",
        SCHEDULE_TYPE: "DAILY",
        SCHEDULE_DAY: "0",
        RUN_TIME: "09:00",
        CRON_SCHEDULE: "",
        RUN_ON_STARTUP: false,
        DEBUG: false,
        HTTP_TIMEOUT: 30,
        LOG_MAX_SIZE_MB: 10,
        LOG_BACKUP_COUNT: 5
    })
    const [originalConfig, setOriginalConfig] = useState<any>({})

    const [loading, setLoading] = useState(true)
    const [saving, setSaving] = useState(false)

    // Custom Toast State (replacing showToast)
    const [toast, setToast] = useState<{ message: string, type: 'success' | 'error' } | null>(null)
    const [timezones, setTimezones] = useState<string[]>([])
    const { theme, toggleTheme } = useTheme()
    const [languages, setLanguages] = useState<{ code: string, name: string }[]>([])

    useEffect(() => {
        fetchConfig()
        fetchLanguages()
        try {
            const tzs = Intl.supportedValuesOf('timeZone')
            setTimezones(tzs)
        } catch (e) {
            console.warn('Browser does not support Intl.supportedValuesOf API, using fallback...')
            setTimezones(['UTC'])
        }
    }, [])

    const fetchConfig = async () => {
        setLoading(true)
        try {
            const res = await fetch(`${API_BASE}/api/config`)
            if (res.ok) {
                const data = await res.json()
                setConfig(data)
                setOriginalConfig(JSON.parse(JSON.stringify(data)))
            }
        } catch (e) {
            console.error(e)
        } finally {
            setLoading(false)
        }
    }

    const fetchLanguages = async () => {
        try {
            const res = await fetch(`${API_BASE}/api/languages`)
            if (res.ok) {
                const data = await res.json()
                setLanguages(data)
            }
        } catch (e) {
            console.error(e)
        }
    }

    const showToast = (message: string, type: 'success' | 'error' = 'success') => {
        setToast({ message, type })
        setTimeout(() => setToast(null), 3000)
    }

    const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
        const target = e.target as HTMLInputElement;
        const value = target.type === 'checkbox' ? target.checked : target.value;
        setConfig((prev: any) => ({
            ...prev,
            [target.name]: value
        }));
    }

    const getDiscordTimestampExample = (style: string) => {
        const now = new Date()
        const futureDate = new Date(now)
        futureDate.setDate(futureDate.getDate() + 1)
        futureDate.setHours(15, 30, 0, 0)

        const fTime = futureDate.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', hour12: false })
        const fTimeLong = futureDate.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit', hour12: false })
        const fDateShort = futureDate.toLocaleDateString([], { month: '2-digit', day: '2-digit', year: 'numeric' })
        const fDateLong = futureDate.toLocaleDateString([], { year: 'numeric', month: 'long', day: 'numeric' })

        switch (style) {
            case 't': return fTime
            case 'T': return fTimeLong
            case 'd': return fDateShort
            case 'D': return fDateLong
            case 'f': return `${fDateLong} ${fTime}`
            case 'F': return `${futureDate.toLocaleDateString([], { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })} ${fTime}`
            case 'R': return "in 1 day"
            default: return "Select a style"
        }
    }

    const handleSave = async (e: React.FormEvent) => {
        e.preventDefault()
        setSaving(true)

        // Calculate diff against originalConfig
        const changesToSubmit: any = {};
        for (const [key, value] of Object.entries(config)) {
            let isChanged = false;
            if (key === 'CALENDAR_URLS') {
                const filteredUrls = (value as any[] || []).filter(u => u.url && u.url.trim() !== '')
                isChanged = JSON.stringify(filteredUrls) !== JSON.stringify(originalConfig[key] || []);
                if (isChanged) changesToSubmit[key] = filteredUrls;
                continue;
            } else {
                if (typeof value === 'boolean') {
                    const origBool = originalConfig[key] === true || originalConfig[key] === 'true';
                    isChanged = value !== origBool;
                } else {
                    const valStr = String(value || '').trim();
                    const origStr = String(originalConfig[key] || '').trim();
                    isChanged = valStr !== origStr;
                }
            }
            if (isChanged) {
                changesToSubmit[key] = key === 'SCHEDULE_DAY' ? String(value) : value;
            }
        }

        if (Object.keys(changesToSubmit).length === 0) {
            showToast("No changes detected to save.")
            setSaving(false)
            return
        }

        try {
            const res = await fetch(`${API_BASE}/api/config`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(changesToSubmit)
            })
            const data = await res.json()
            if (res.ok && data.success) {
                showToast("Settings saved successfully!", 'success')
                setOriginalConfig((prev: any) => ({ ...prev, ...changesToSubmit }))
            } else {
                showToast(data.error || "Failed to save settings", 'error')
            }
        } catch (e) {
            showToast("Failed to communicate with server.", 'error')
        } finally {
            setSaving(false)
        }
    }

    const handleDiscard = () => {
        fetchConfig()
        showToast("Changes discarded.")
    }

    // Calendar URLs handlers
    const handleAddCalendarUrl = () => {
        setConfig((prev: any) => ({
            ...prev,
            CALENDAR_URLS: [...(prev.CALENDAR_URLS || []), { url: '', type: 'tv' }]
        }))
    }

    const handleRemoveCalendarUrl = (index: number) => {
        setConfig((prev: any) => {
            const newUrls = [...(prev.CALENDAR_URLS || [])]
            newUrls.splice(index, 1)
            return { ...prev, CALENDAR_URLS: newUrls }
        })
    }

    const handleCalendarUrlChange = (index: number, field: 'url' | 'type', value: string) => {
        setConfig((prev: any) => {
            const newUrls = [...(prev.CALENDAR_URLS || [])]
            newUrls[index][field] = value
            return { ...prev, CALENDAR_URLS: newUrls }
        })
    }

    if (loading) return (
        <div style={{ display: 'flex', justifyContent: 'center', marginTop: '50px' }}>
            <div className="loading">
                <div className="spinner"></div>
                <p>Loading settings...</p>
            </div>
        </div>
    )

    return (
        <>
            <header className="header">
                <div className="container">
                    <div className="header-content">
                        <div className="logo">
                            <svg width="32" height="32" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
                                <rect x="4" y="7" width="24" height="20" rx="2" stroke="currentColor" strokeWidth="2" />
                                <line x1="4" y1="12" x2="28" y2="12" stroke="currentColor" strokeWidth="2" strokeLinecap="round" />
                                <line x1="10" y1="4" x2="10" y2="7" stroke="currentColor" strokeWidth="2" strokeLinecap="round" />
                                <line x1="22" y1="4" x2="22" y2="7" stroke="currentColor" strokeWidth="2" strokeLinecap="round" />
                            </svg>
                            <h1>Calendarr</h1>
                        </div>
                        <nav className="nav">
                            <a href="#/" className="nav-link">Dashboard</a>
                            <a href="#/settings" className="nav-link" style={{ background: 'var(--color-success)', boxShadow: '0 0 0 #000', transform: 'translate(2px, 2px)' }}>Settings</a>
                            <button className="theme-toggle" onClick={toggleTheme} title={theme === 'light' ? 'Switch to dark mode' : 'Switch to light mode'}>
                                {theme === 'light' ? (
                                    <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                                        <path d="M21 12.79A9 9 0 1111.21 3a7 7 0 009.79 9.79z" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
                                    </svg>
                                ) : (
                                    <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                                        <circle cx="12" cy="12" r="5" stroke="currentColor" strokeWidth="2" />
                                        <path d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42" stroke="currentColor" strokeWidth="2" strokeLinecap="round" />
                                    </svg>
                                )}
                            </button>
                        </nav>
                    </div>
                </div>
            </header>

            <main className="main">
                <div className="container" id="settingsView">
                    <div className="section-header">
                        <h2>Configuration Settings</h2>
                    </div>

                    <div className="section">
                        <form id="settingsForm" onSubmit={handleSave}>
                            {/* 1. Source Configuration */}
                            <details className="settings-group">
                                <summary>
                                    <h3>Source Configuration</h3>
                                    <span className="details-arrow">▼</span>
                                </summary>
                                <div className="form-group">
                                    <label>Calendar URLs (<span id="calendarUrlCount">{(config.CALENDAR_URLS || []).length}</span>)</label>
                                    <div className="calendar-url-container">
                                        {(config.CALENDAR_URLS || []).length === 0 ? (
                                            <p style={{ color: 'var(--color-text-dim)', fontSize: '0.9em', marginBottom: '5px' }}>No calendar URLs added. Add one to get started.</p>
                                        ) : (
                                            (config.CALENDAR_URLS || []).map((urlObj: any, index: number) => (
                                                <div key={index} className="calendar-url-item">
                                                    <div className="calendar-url-inputs">
                                                        <input
                                                            type="url"
                                                            className="brutal-input url-val"
                                                            placeholder="https://..."
                                                            value={urlObj.url || ''}
                                                            onChange={(e) => handleCalendarUrlChange(index, 'url', e.target.value)}
                                                            required
                                                        />
                                                        <select
                                                            className="brutal-select type-val"
                                                            value={urlObj.type || 'tv'}
                                                            onChange={(e) => handleCalendarUrlChange(index, 'type', e.target.value)}
                                                        >
                                                            <option value="tv">TV Show (Sonarr)</option>
                                                            <option value="movie">Movie (Radarr)</option>
                                                        </select>
                                                    </div>
                                                    <button type="button" className="brutal-btn brutal-btn-danger remove-url-btn" onClick={() => handleRemoveCalendarUrl(index)} title="Remove">&times;</button>
                                                </div>
                                            ))
                                        )}
                                    </div>
                                    <button type="button" className="brutal-btn mt-none" onClick={handleAddCalendarUrl} style={{ marginTop: '0' }}>+ Add Calendar</button>
                                </div>
                            </details>

                            {/* 2. Discord Integration */}
                            <details className="settings-group">
                                <summary>
                                    <h3>Discord Integration</h3>
                                    <span className="details-arrow">▼</span>
                                </summary>
                                <div className="form-group checkbox-group">
                                    <input type="checkbox" id="set_USE_DISCORD" name="USE_DISCORD" checked={!!config.USE_DISCORD} onChange={handleChange} />
                                    <label htmlFor="set_USE_DISCORD">Enable Discord Support</label>
                                </div>
                                {config.USE_DISCORD && (
                                    <>
                                        <div className="form-group">
                                            <label htmlFor="set_DISCORD_WEBHOOK_URL">Discord Webhook URL</label>
                                            <input type="password" id="set_DISCORD_WEBHOOK_URL" name="DISCORD_WEBHOOK_URL" className="brutal-input" value={config.DISCORD_WEBHOOK_URL || ''} onChange={handleChange} />
                                        </div>
                                        <div className="form-group">
                                            <label htmlFor="set_DISCORD_MENTION_ROLE_ID">Discord Mention Role ID</label>
                                            <input type="text" id="set_DISCORD_MENTION_ROLE_ID" name="DISCORD_MENTION_ROLE_ID" className="brutal-input" value={config.DISCORD_MENTION_ROLE_ID || ''} onChange={handleChange} />
                                        </div>
                                        <div className="form-group checkbox-group">
                                            <input type="checkbox" id="set_DISCORD_HIDE_MENTION_INSTRUCTIONS" name="DISCORD_HIDE_MENTION_INSTRUCTIONS" checked={!!config.DISCORD_HIDE_MENTION_INSTRUCTIONS} onChange={handleChange} />
                                            <label htmlFor="set_DISCORD_HIDE_MENTION_INSTRUCTIONS">Hide Discord Mention Instructions</label>
                                        </div>
                                        <div className="form-group">
                                            <label htmlFor="set_DISCORD_TIMESTAMP_STYLE">Discord Timestamp Style</label>
                                            <select id="set_DISCORD_TIMESTAMP_STYLE" name="DISCORD_TIMESTAMP_STYLE" className="brutal-select" value={config.DISCORD_TIMESTAMP_STYLE || 'R'} onChange={handleChange}>
                                                <option value="t">Short Time</option>
                                                <option value="T">Long Time</option>
                                                <option value="d">Short Date</option>
                                                <option value="D">Long Date</option>
                                                <option value="f">Short Date/Time</option>
                                                <option value="F">Long Date/Time</option>
                                                <option value="R">Relative Time</option>
                                            </select>
                                            <div className="timestamp-preview mt-sm" style={{ fontFamily: 'monospace', background: 'var(--color-surface-hover)', padding: '2px 4px', borderRadius: '3px', border: '1px solid var(--color-border)', fontSize: '0.9em', marginTop: 'var(--spacing-sm)' }}>
                                                <strong>Example:</strong> <span style={{ background: 'rgba(0,0,0,0.2)', padding: '2px 4px', borderRadius: '3px' }}>{getDiscordTimestampExample(config.DISCORD_TIMESTAMP_STYLE || 'R')}</span>
                                            </div>
                                        </div>
                                        <div className="form-group checkbox-group">
                                            <input type="checkbox" id="set_ENABLE_CUSTOM_DISCORD_FOOTER" name="ENABLE_CUSTOM_DISCORD_FOOTER" checked={!!config.ENABLE_CUSTOM_DISCORD_FOOTER} onChange={handleChange} />
                                            <label htmlFor="set_ENABLE_CUSTOM_DISCORD_FOOTER">Enable Custom Discord Footer</label>
                                        </div>
                                    </>
                                )}
                            </details>

                            {/* 3. Slack Integration */}
                            <details className="settings-group">
                                <summary>
                                    <h3>Slack Integration</h3>
                                    <span className="details-arrow">▼</span>
                                </summary>
                                <div className="form-group checkbox-group">
                                    <input type="checkbox" id="set_USE_SLACK" name="USE_SLACK" checked={!!config.USE_SLACK} onChange={handleChange} />
                                    <label htmlFor="set_USE_SLACK">Enable Slack Support</label>
                                </div>
                                {config.USE_SLACK && (
                                    <>
                                        <div className="form-group">
                                            <label htmlFor="set_SLACK_WEBHOOK_URL">Slack Webhook URL</label>
                                            <input type="password" id="set_SLACK_WEBHOOK_URL" name="SLACK_WEBHOOK_URL" className="brutal-input" value={config.SLACK_WEBHOOK_URL || ''} onChange={handleChange} />
                                        </div>
                                        <div className="form-group checkbox-group">
                                            <input type="checkbox" id="set_ENABLE_CUSTOM_SLACK_FOOTER" name="ENABLE_CUSTOM_SLACK_FOOTER" checked={!!config.ENABLE_CUSTOM_SLACK_FOOTER} onChange={handleChange} />
                                            <label htmlFor="set_ENABLE_CUSTOM_SLACK_FOOTER">Enable Custom Slack Footer</label>
                                        </div>
                                    </>
                                )}
                            </details>

                            {/* 4. Display & Formatting */}
                            <details className="settings-group">
                                <summary>
                                    <h3>Display & Formatting</h3>
                                    <span className="details-arrow">▼</span>
                                </summary>
                                <div className="form-group">
                                    <label htmlFor="set_APP_LANGUAGE">Language</label>
                                    <select id="set_APP_LANGUAGE" name="APP_LANGUAGE" className="brutal-select" value={config.APP_LANGUAGE || 'EN'} onChange={handleChange}>
                                        {languages.length > 0 ? languages.map(lang => (
                                            <option key={lang.code} value={lang.code}>{lang.name}</option>
                                        )) : (
                                            <option value="EN">English</option>
                                        )}
                                    </select>
                                </div>
                                <div className="form-group">
                                    <label htmlFor="set_TZ">Timezone</label>
                                    <select id="set_TZ" name="TZ" className="brutal-select" value={config.TZ || 'UTC'} onChange={handleChange}>
                                        {timezones.map(tz => {
                                            const isDefault = tz === Intl.DateTimeFormat().resolvedOptions().timeZone;
                                            return (
                                                <option key={tz} value={tz}>
                                                    {isDefault ? `${tz} (System Default)` : tz}
                                                </option>
                                            )
                                        })}
                                    </select>
                                </div>
                                <div className="form-group">
                                    <label htmlFor="set_PASSED_EVENT_HANDLING">Passed Event Handling</label>
                                    <select id="set_PASSED_EVENT_HANDLING" name="PASSED_EVENT_HANDLING" className="brutal-select" value={config.PASSED_EVENT_HANDLING || 'DISPLAY'} onChange={handleChange}>
                                        <option value="DISPLAY">Display</option>
                                        <option value="HIDE">Hide</option>
                                        <option value="STRIKE">Strike</option>
                                    </select>
                                </div>

                                <div className="form-group checkbox-group">
                                    <input type="checkbox" id="set_DEDUPLICATE_EVENTS" name="DEDUPLICATE_EVENTS" checked={!!config.DEDUPLICATE_EVENTS} onChange={handleChange} />
                                    <label htmlFor="set_DEDUPLICATE_EVENTS">Deduplicate Events</label>
                                </div>
                                <div className="form-group checkbox-group">
                                    <input type="checkbox" id="set_USE_24_HOUR" name="USE_24_HOUR" checked={!!config.USE_24_HOUR} onChange={handleChange} />
                                    <label htmlFor="set_USE_24_HOUR">Use 24-Hour Time</label>
                                </div>
                                <div className="form-group checkbox-group">
                                    <input type="checkbox" id="set_ADD_LEADING_ZERO" name="ADD_LEADING_ZERO" checked={!!config.ADD_LEADING_ZERO} onChange={handleChange} />
                                    <label htmlFor="set_ADD_LEADING_ZERO">Add Leading Zero</label>
                                </div>
                                <div className="form-group checkbox-group">
                                    <input type="checkbox" id="set_DISPLAY_TIME" name="DISPLAY_TIME" checked={!!config.DISPLAY_TIME} onChange={handleChange} />
                                    <label htmlFor="set_DISPLAY_TIME">Display Release Time</label>
                                </div>
                                <div className="form-group checkbox-group">
                                    <input type="checkbox" id="set_SHOW_DATE_RANGE" name="SHOW_DATE_RANGE" checked={!!config.SHOW_DATE_RANGE} onChange={handleChange} />
                                    <label htmlFor="set_SHOW_DATE_RANGE">Show Date Range</label>
                                </div>
                                <div className="form-group checkbox-group">
                                    <input type="checkbox" id="set_SHOW_TIMEZONE_IN_SUBHEADER" name="SHOW_TIMEZONE_IN_SUBHEADER" checked={!!config.SHOW_TIMEZONE_IN_SUBHEADER} onChange={handleChange} />
                                    <label htmlFor="set_SHOW_TIMEZONE_IN_SUBHEADER">Show Timezone in Subheader</label>
                                </div>
                            </details>

                            {/* 5. Scheduling */}
                            <details className="settings-group">
                                <summary>
                                    <h3>Scheduling</h3>
                                    <span className="details-arrow">▼</span>
                                </summary>
                                <div className="form-group">
                                    <label htmlFor="set_SCHEDULE_TYPE">Schedule Type</label>
                                    <select id="set_SCHEDULE_TYPE" name="SCHEDULE_TYPE" className="brutal-select" value={config.SCHEDULE_TYPE || 'DAILY'} onChange={handleChange}>
                                        <option value="DAILY">Daily</option>
                                        <option value="WEEKLY">Weekly</option>
                                    </select>
                                </div>
                                {config.SCHEDULE_TYPE === 'WEEKLY' && (
                                    <div className="form-group">
                                        <label htmlFor="set_SCHEDULE_DAY">Schedule Day</label>
                                        <select id="set_SCHEDULE_DAY" name="SCHEDULE_DAY" className="brutal-select" value={config.SCHEDULE_DAY || '0'} onChange={handleChange}>
                                            <option value="0">Sunday</option>
                                            <option value="1">Monday</option>
                                            <option value="2">Tuesday</option>
                                            <option value="3">Wednesday</option>
                                            <option value="4">Thursday</option>
                                            <option value="5">Friday</option>
                                            <option value="6">Saturday</option>
                                        </select>
                                    </div>
                                )}
                                <div className="form-group">
                                    <label htmlFor="set_RUN_TIME">Run Time (HH:MM)</label>
                                    <input type="time" id="set_RUN_TIME" name="RUN_TIME" className="brutal-input" value={config.RUN_TIME || ''} onChange={handleChange} />
                                </div>
                                <div className="form-group">
                                    <label htmlFor="set_CRON_SCHEDULE">Cron Schedule (Overrides others)</label>
                                    <input type="text" id="set_CRON_SCHEDULE" name="CRON_SCHEDULE" className="brutal-input" placeholder="e.g. 0 10 * * 1" value={config.CRON_SCHEDULE || ''} onChange={handleChange} />
                                </div>
                                <div className="form-group checkbox-group">
                                    <input type="checkbox" id="set_RUN_ON_STARTUP" name="RUN_ON_STARTUP" checked={!!config.RUN_ON_STARTUP} onChange={handleChange} />
                                    <label htmlFor="set_RUN_ON_STARTUP">Run on Startup</label>
                                </div>
                            </details>

                            {/* 6. Advanced Settings */}
                            <details className="settings-group">
                                <summary>
                                    <h3>Advanced Settings</h3>
                                    <span className="details-arrow">▼</span>
                                </summary>
                                <div className="form-group checkbox-group">
                                    <input type="checkbox" id="set_DEBUG" name="DEBUG" checked={!!config.DEBUG} onChange={handleChange} />
                                    <label htmlFor="set_DEBUG">Enable Debug Logging</label>
                                </div>
                                <div className="form-group">
                                    <label htmlFor="set_HTTP_TIMEOUT">HTTP Timeout (seconds)</label>
                                    <input type="number" id="set_HTTP_TIMEOUT" name="HTTP_TIMEOUT" min="1" className="brutal-input" value={config.HTTP_TIMEOUT || 30} onChange={handleChange} />
                                </div>
                                <div className="form-group">
                                    <label htmlFor="set_LOG_MAX_SIZE_MB">Log Max Size (MB)</label>
                                    <input type="number" id="set_LOG_MAX_SIZE_MB" name="LOG_MAX_SIZE_MB" min="1" className="brutal-input" value={config.LOG_MAX_SIZE_MB || 10} onChange={handleChange} />
                                </div>
                                <div className="form-group">
                                    <label htmlFor="set_LOG_BACKUP_COUNT">Log Backup Count</label>
                                    <input type="number" id="set_LOG_BACKUP_COUNT" name="LOG_BACKUP_COUNT" min="1" className="brutal-input" value={config.LOG_BACKUP_COUNT || 5} onChange={handleChange} />
                                </div>
                            </details>

                            <div className="settings-actions" style={{ display: 'flex', justifyContent: 'flex-end', gap: 'var(--spacing-sm)', marginTop: 'var(--spacing-lg)' }}>
                                <button type="button" className="brutal-btn brutal-btn-outline" onClick={handleDiscard}>Discard Changes</button>
                                <button type="submit" className="brutal-btn brutal-btn-primary" disabled={saving}>{saving ? 'Saving...' : 'Save Changes'}</button>
                            </div>
                        </form>
                    </div>
                </div>
            </main>

            <footer className="footer">
                <div className="container">
                    <p><a href="https://github.com/khw315/calendarr" target="_blank" rel="noreferrer">GitHub</a></p>
                </div>
            </footer>

            {/* Toast Notification Container */}
            {toast && (
                <div className="toast-container" style={{ position: 'fixed', bottom: '20px', right: '20px', zIndex: 1000, display: 'flex', flexDirection: 'column', gap: '10px' }}>
                    <div className="toast" style={{
                        background: toast.type === 'success' ? 'var(--color-success)' : 'var(--color-danger)',
                        color: toast.type === 'success' ? 'var(--color-text)' : '#fff',
                        padding: '12px 24px',
                        border: '2px solid var(--color-border)',
                        boxShadow: '4px 4px 0 var(--color-border)',
                        fontWeight: 800,
                        animation: 'slideIn 0.3s ease-out forwards',
                        textTransform: 'uppercase'
                    }}>
                        {toast.message}
                    </div>
                </div>
            )}
        </>
    )
}
