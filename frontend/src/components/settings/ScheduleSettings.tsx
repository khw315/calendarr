import React from 'react'
import { ConfigState } from '@/types/config'

interface ScheduleSettingsProps {
  config: ConfigState
  onChange: (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => void
}

export default function ScheduleSettings({ config, onChange }: ScheduleSettingsProps) {
  return (
    <details className="settings-group">
      <summary>
        <h3>Scheduling</h3>
        <span className="details-arrow">▼</span>
      </summary>
      <div className="form-group">
        <label htmlFor="set_SCHEDULE_TYPE">Schedule Type</label>
        <select
          id="set_SCHEDULE_TYPE"
          name="SCHEDULE_TYPE"
          className="brutal-select"
          value={config.SCHEDULE_TYPE || 'DAILY'}
          onChange={onChange}
        >
          <option value="DAILY">Daily</option>
          <option value="WEEKLY">Weekly</option>
        </select>
      </div>

      {config.SCHEDULE_TYPE === 'WEEKLY' && (
        <div className="form-group">
          <label htmlFor="set_SCHEDULE_DAY">Schedule Day</label>
          <select
            id="set_SCHEDULE_DAY"
            name="SCHEDULE_DAY"
            className="brutal-select"
            value={config.SCHEDULE_DAY || '0'}
            onChange={onChange}
          >
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
        <input
          type="time"
          id="set_RUN_TIME"
          name="RUN_TIME"
          className="brutal-input"
          value={config.RUN_TIME || ''}
          onChange={onChange}
        />
      </div>

      <div className="form-group">
        <label htmlFor="set_CRON_SCHEDULE">Cron Schedule (Overrides others)</label>
        <input
          type="text"
          id="set_CRON_SCHEDULE"
          name="CRON_SCHEDULE"
          className="brutal-input"
          placeholder="e.g. 0 10 * * 1"
          value={config.CRON_SCHEDULE || ''}
          onChange={onChange}
        />
      </div>

      <div className="form-group checkbox-group">
        <input
          type="checkbox"
          id="set_RUN_ON_STARTUP"
          name="RUN_ON_STARTUP"
          checked={!!config.RUN_ON_STARTUP}
          onChange={onChange}
        />
        <label htmlFor="set_RUN_ON_STARTUP">Run on Startup</label>
      </div>
    </details>
  )
}
