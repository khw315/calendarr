import React from 'react'
import { ConfigState } from '@/types/config'
import { ThemeMode } from '@/lib/useTheme'

interface DisplaySettingsProps {
  config: ConfigState
  languages: { code: string; name: string }[]
  timezones: { iana: string; label: string }[]
  themeMode: ThemeMode
  onThemeChange: (mode: ThemeMode) => void
  onChange: (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => void
}

export default function DisplaySettings({
  config,
  languages,
  timezones,
  themeMode,
  onThemeChange,
  onChange
}: DisplaySettingsProps) {
  return (
    <details className="settings-group">
      <summary>
        <h3>Display & Formatting</h3>
        <span className="details-arrow">▼</span>
      </summary>

      <div className="form-group">
        <label htmlFor="set_THEME_MODE">Theme Mode</label>
        <select
          id="set_THEME_MODE"
          name="THEME_MODE"
          className="brutal-select"
          value={themeMode}
          onChange={(e) => onThemeChange(e.target.value as ThemeMode)}
        >
          <option value="system">System (Follow Browser/OS)</option>
          <option value="light">Light Mode</option>
          <option value="dark">Dark Mode</option>
        </select>
      </div>

      <div className="form-group">
        <label htmlFor="set_APP_LANGUAGE">Language</label>
        <select
          id="set_APP_LANGUAGE"
          name="APP_LANGUAGE"
          className="brutal-select"
          value={config.APP_LANGUAGE || 'EN'}
          onChange={onChange}
        >
          {languages.length > 0 ? (
            languages.map(lang => (
              <option key={lang.code} value={lang.code}>
                {lang.name}
              </option>
            ))
          ) : (
            <option value="EN">English</option>
          )}
        </select>
      </div>

      <div className="form-group">
        <label htmlFor="set_TZ">Timezone</label>
        <select
          id="set_TZ"
          name="TZ"
          className="brutal-select"
          value={config.TZ || 'UTC'}
          onChange={onChange}
        >
          {timezones.map(({ iana, label }) => {
            const isDefault = typeof Intl !== 'undefined' && iana === Intl.DateTimeFormat().resolvedOptions().timeZone
            const displayLabel = label && label !== iana ? `${iana} — ${label}` : iana
            return (
              <option key={iana} value={iana}>
                {isDefault ? `${displayLabel} (System Default)` : displayLabel}
              </option>
            )
          })}
        </select>
      </div>

      <div className="form-group">
        <label htmlFor="set_PASSED_EVENT_HANDLING">Passed Event Handling</label>
        <select
          id="set_PASSED_EVENT_HANDLING"
          name="PASSED_EVENT_HANDLING"
          className="brutal-select"
          value={config.PASSED_EVENT_HANDLING || 'DISPLAY'}
          onChange={onChange}
        >
          <option value="DISPLAY">Display</option>
          <option value="HIDE">Hide</option>
          <option value="STRIKE">Strike</option>
        </select>
      </div>

      <div className="form-group checkbox-group">
        <input
          type="checkbox"
          id="set_DEDUPLICATE_EVENTS"
          name="DEDUPLICATE_EVENTS"
          checked={!!config.DEDUPLICATE_EVENTS}
          onChange={onChange}
        />
        <label htmlFor="set_DEDUPLICATE_EVENTS">Deduplicate Events</label>
      </div>

      <div className="form-group checkbox-group">
        <input
          type="checkbox"
          id="set_USE_24_HOUR"
          name="USE_24_HOUR"
          checked={!!config.USE_24_HOUR}
          onChange={onChange}
        />
        <label htmlFor="set_USE_24_HOUR">Use 24-Hour Time</label>
      </div>

      <div className="form-group checkbox-group">
        <input
          type="checkbox"
          id="set_ADD_LEADING_ZERO"
          name="ADD_LEADING_ZERO"
          checked={!!config.ADD_LEADING_ZERO}
          onChange={onChange}
        />
        <label htmlFor="set_ADD_LEADING_ZERO">Add Leading Zero</label>
      </div>

      <div className="form-group checkbox-group">
        <input
          type="checkbox"
          id="set_DISPLAY_TIME"
          name="DISPLAY_TIME"
          checked={!!config.DISPLAY_TIME}
          onChange={onChange}
        />
        <label htmlFor="set_DISPLAY_TIME">Display Release Time</label>
      </div>

      <div className="form-group checkbox-group">
        <input
          type="checkbox"
          id="set_SHOW_DATE_RANGE"
          name="SHOW_DATE_RANGE"
          checked={!!config.SHOW_DATE_RANGE}
          onChange={onChange}
        />
        <label htmlFor="set_SHOW_DATE_RANGE">Show Date Range</label>
      </div>

      <div className="form-group checkbox-group">
        <input
          type="checkbox"
          id="set_SHOW_TIMEZONE_IN_SUBHEADER"
          name="SHOW_TIMEZONE_IN_SUBHEADER"
          checked={!!config.SHOW_TIMEZONE_IN_SUBHEADER}
          onChange={onChange}
        />
        <label htmlFor="set_SHOW_TIMEZONE_IN_SUBHEADER">Show Timezone in Subheader</label>
      </div>
    </details>
  )
}
