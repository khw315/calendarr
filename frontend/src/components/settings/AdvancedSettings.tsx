import React from 'react'
import { ConfigState } from '@/types/config'

interface AdvancedSettingsProps {
  config: ConfigState
  onChange: (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => void
}

export default function AdvancedSettings({ config, onChange }: AdvancedSettingsProps) {
  return (
    <details className="settings-group">
      <summary>
        <h3>Advanced Settings</h3>
        <span className="details-arrow">▼</span>
      </summary>
      <div className="form-group checkbox-group">
        <input
          type="checkbox"
          id="set_DEBUG"
          name="DEBUG"
          checked={!!config.DEBUG}
          onChange={onChange}
        />
        <label htmlFor="set_DEBUG">Enable Debug Logging</label>
      </div>

      <div className="form-group">
        <label htmlFor="set_HTTP_TIMEOUT">HTTP Timeout (seconds)</label>
        <input
          type="number"
          id="set_HTTP_TIMEOUT"
          name="HTTP_TIMEOUT"
          min="1"
          className="brutal-input"
          value={config.HTTP_TIMEOUT || 30}
          onChange={onChange}
        />
      </div>

      <div className="form-group">
        <label htmlFor="set_LOG_MAX_SIZE_MB">Log Max Size (MB)</label>
        <input
          type="number"
          id="set_LOG_MAX_SIZE_MB"
          name="LOG_MAX_SIZE_MB"
          min="1"
          className="brutal-input"
          value={config.LOG_MAX_SIZE_MB || 10}
          onChange={onChange}
        />
      </div>

      <div className="form-group">
        <label htmlFor="set_LOG_BACKUP_COUNT">Log Backup Count</label>
        <input
          type="number"
          id="set_LOG_BACKUP_COUNT"
          name="LOG_BACKUP_COUNT"
          min="1"
          className="brutal-input"
          value={config.LOG_BACKUP_COUNT || 5}
          onChange={onChange}
        />
      </div>
    </details>
  )
}
