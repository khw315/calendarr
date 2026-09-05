import React from 'react'
import { ConfigState } from '@/types/config'

interface NotificationSettingsProps {
  config: ConfigState
  onChange: (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => void
}

export function getDiscordTimestampExample(style: string): string {
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

export default function NotificationSettings({ config, onChange }: NotificationSettingsProps) {
  return (
    <>
      {/* Discord Integration */}
      <details className="settings-group" name="settings-accordion">
        <summary>
          <h3>Discord Integration</h3>
          <span className="details-arrow">▼</span>
        </summary>
        <div className="form-group checkbox-group">
          <input
            type="checkbox"
            id="set_USE_DISCORD"
            name="USE_DISCORD"
            checked={!!config.USE_DISCORD}
            onChange={onChange}
          />
          <label htmlFor="set_USE_DISCORD">Enable Discord Support</label>
        </div>
        {config.USE_DISCORD && (
          <>
            <div className="form-group">
              <label htmlFor="set_DISCORD_WEBHOOK_URL">Discord Webhook URL</label>
              <input
                type="password"
                id="set_DISCORD_WEBHOOK_URL"
                name="DISCORD_WEBHOOK_URL"
                className="brutal-input"
                value={config.DISCORD_WEBHOOK_URL || ''}
                onChange={onChange}
              />
            </div>
            <div className="form-group">
              <label htmlFor="set_DISCORD_MENTION_ROLE_ID">Discord Mention Role ID</label>
              <input
                type="text"
                id="set_DISCORD_MENTION_ROLE_ID"
                name="DISCORD_MENTION_ROLE_ID"
                className="brutal-input"
                value={config.DISCORD_MENTION_ROLE_ID || ''}
                onChange={onChange}
              />
            </div>
            <div className="form-group checkbox-group">
              <input
                type="checkbox"
                id="set_DISCORD_HIDE_MENTION_INSTRUCTIONS"
                name="DISCORD_HIDE_MENTION_INSTRUCTIONS"
                checked={!!config.DISCORD_HIDE_MENTION_INSTRUCTIONS}
                onChange={onChange}
              />
              <label htmlFor="set_DISCORD_HIDE_MENTION_INSTRUCTIONS">
                Hide Discord Mention Instructions
              </label>
            </div>
            <div className="form-group">
              <label htmlFor="set_DISCORD_TIMESTAMP_STYLE">Discord Timestamp Style</label>
              <select
                id="set_DISCORD_TIMESTAMP_STYLE"
                name="DISCORD_TIMESTAMP_STYLE"
                className="brutal-select"
                value={config.DISCORD_TIMESTAMP_STYLE || 'R'}
                onChange={onChange}
              >
                <option value="t">Short Time</option>
                <option value="T">Long Time</option>
                <option value="d">Short Date</option>
                <option value="D">Long Date</option>
                <option value="f">Short Date/Time</option>
                <option value="F">Long Date/Time</option>
                <option value="R">Relative Time</option>
              </select>
              <div
                className="timestamp-preview mt-sm"
                style={{
                  fontFamily: 'monospace',
                  background: 'var(--color-surface-hover)',
                  padding: '2px 4px',
                  borderRadius: '3px',
                  border: '1px solid var(--color-border)',
                  fontSize: '0.9em',
                  marginTop: 'var(--spacing-sm)'
                }}
              >
                <strong>Example:</strong>{' '}
                <span style={{ background: 'rgba(0,0,0,0.2)', padding: '2px 4px', borderRadius: '3px' }}>
                  {getDiscordTimestampExample(config.DISCORD_TIMESTAMP_STYLE || 'R')}
                </span>
              </div>
            </div>
            <div className="form-group checkbox-group">
              <input
                type="checkbox"
                id="set_ENABLE_CUSTOM_DISCORD_FOOTER"
                name="ENABLE_CUSTOM_DISCORD_FOOTER"
                checked={!!config.ENABLE_CUSTOM_DISCORD_FOOTER}
                onChange={onChange}
              />
              <label htmlFor="set_ENABLE_CUSTOM_DISCORD_FOOTER">
                Enable Custom Discord Footer
              </label>
            </div>
          </>
        )}
      </details>

      {/* Slack Integration */}
      <details className="settings-group" name="settings-accordion">
        <summary>
          <h3>Slack Integration</h3>
          <span className="details-arrow">▼</span>
        </summary>
        <div className="form-group checkbox-group">
          <input
            type="checkbox"
            id="set_USE_SLACK"
            name="USE_SLACK"
            checked={!!config.USE_SLACK}
            onChange={onChange}
          />
          <label htmlFor="set_USE_SLACK">Enable Slack Support</label>
        </div>
        {config.USE_SLACK && (
          <>
            <div className="form-group">
              <label htmlFor="set_SLACK_WEBHOOK_URL">Slack Webhook URL</label>
              <input
                type="password"
                id="set_SLACK_WEBHOOK_URL"
                name="SLACK_WEBHOOK_URL"
                className="brutal-input"
                value={config.SLACK_WEBHOOK_URL || ''}
                onChange={onChange}
              />
            </div>
            <div className="form-group checkbox-group">
              <input
                type="checkbox"
                id="set_ENABLE_CUSTOM_SLACK_FOOTER"
                name="ENABLE_CUSTOM_SLACK_FOOTER"
                checked={!!config.ENABLE_CUSTOM_SLACK_FOOTER}
                onChange={onChange}
              />
              <label htmlFor="set_ENABLE_CUSTOM_SLACK_FOOTER">
                Enable Custom Slack Footer
              </label>
            </div>
          </>
        )}
      </details>
    </>
  )
}
