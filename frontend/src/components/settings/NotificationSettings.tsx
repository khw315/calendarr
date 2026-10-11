import React, { useState } from 'react'
import { ConfigState } from '@/types/config'
import { apiService } from '@/services/api'

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
  const [testingPlatform, setTestingPlatform] = useState<string | null>(null)
  const [testStatus, setTestStatus] = useState<Record<string, { type: 'success' | 'error'; message: string } | null>>({})

  const handleTestWebhook = async (platform: string, url: string, key = platform) => {
    if (!url) return
    setTestingPlatform(key)
    setTestStatus((prev) => ({ ...prev, [key]: null }))
    try {
      const res = await apiService.testWebhook(platform, url)
      setTestStatus((prev) => ({
        ...prev,
        [key]: { type: 'success', message: res.message || 'Test notification sent!' }
      }))
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Test failed'
      setTestStatus((prev) => ({
        ...prev,
        [key]: { type: 'error', message: msg }
      }))
    } finally {
      setTestingPlatform(null)
    }
  }

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
              <label htmlFor="set_DISCORD_WEBHOOK_URL">Default Discord Webhook URL</label>
              <input
                type="password"
                id="set_DISCORD_WEBHOOK_URL"
                name="DISCORD_WEBHOOK_URL"
                className="brutal-input"
                value={config.DISCORD_WEBHOOK_URL || ''}
                onChange={onChange}
                placeholder="https://discord.com/api/webhooks/..."
              />
              <div style={{ display: 'flex', gap: '8px', alignItems: 'center', marginTop: '6px', flexWrap: 'wrap' }}>
                <button
                  type="button"
                  className="brutal-btn"
                  style={{ padding: '4px 12px', fontSize: '0.85rem', cursor: config.DISCORD_WEBHOOK_URL ? 'pointer' : 'not-allowed' }}
                  disabled={!config.DISCORD_WEBHOOK_URL || testingPlatform === 'discord'}
                  onClick={() => handleTestWebhook('discord', config.DISCORD_WEBHOOK_URL || '', 'discord')}
                >
                  {testingPlatform === 'discord' ? 'Testing...' : 'Test Discord Webhook'}
                </button>
                {testStatus['discord'] && (
                  <span
                    style={{
                      fontSize: '0.85rem',
                      fontWeight: 600,
                      color: testStatus['discord'].type === 'success' ? '#16a34a' : '#dc2626'
                    }}
                  >
                    {testStatus['discord'].type === 'success' ? '✅ ' : '❌ '}
                    {testStatus['discord'].message}
                  </span>
                )}
              </div>
            </div>

            {/* Separate Discord Channels */}
            <details
              style={{ marginTop: '12px', marginBottom: '12px', border: '1px dashed var(--color-border)', borderRadius: '4px', padding: '8px 12px' }}
              open={!!(config.DISCORD_TV_WEBHOOK_URL || config.DISCORD_MOVIE_WEBHOOK_URL)}
            >
              <summary style={{ cursor: 'pointer', fontWeight: 600, fontSize: '0.9rem' }}>
                📺 🎬 Separate TV and Movie Channels (Optional)
              </summary>
              <p style={{ fontSize: '0.8rem', color: 'var(--color-text-muted)', margin: '6px 0 12px 0' }}>
                Route TV shows (Sonarr) and movies (Radarr) to dedicated Discord channels. If not specified, notifications fall back to the default webhook above.
              </p>
              <div className="form-group">
                <label htmlFor="set_DISCORD_TV_WEBHOOK_URL">Discord TV Webhook URL</label>
                <input
                  type="password"
                  id="set_DISCORD_TV_WEBHOOK_URL"
                  name="DISCORD_TV_WEBHOOK_URL"
                  className="brutal-input"
                  value={config.DISCORD_TV_WEBHOOK_URL || ''}
                  onChange={onChange}
                  placeholder="https://discord.com/api/webhooks/..."
                />
                <div style={{ display: 'flex', gap: '8px', alignItems: 'center', marginTop: '6px', flexWrap: 'wrap' }}>
                  <button
                    type="button"
                    className="brutal-btn"
                    style={{ padding: '4px 12px', fontSize: '0.85rem', cursor: config.DISCORD_TV_WEBHOOK_URL ? 'pointer' : 'not-allowed' }}
                    disabled={!config.DISCORD_TV_WEBHOOK_URL || testingPlatform === 'discord_tv'}
                    onClick={() => handleTestWebhook('discord', config.DISCORD_TV_WEBHOOK_URL || '', 'discord_tv')}
                  >
                    {testingPlatform === 'discord_tv' ? 'Testing...' : 'Test TV Webhook'}
                  </button>
                  {testStatus['discord_tv'] && (
                    <span
                      style={{
                        fontSize: '0.85rem',
                        fontWeight: 600,
                        color: testStatus['discord_tv'].type === 'success' ? '#16a34a' : '#dc2626'
                      }}
                    >
                      {testStatus['discord_tv'].type === 'success' ? '✅ ' : '❌ '}
                      {testStatus['discord_tv'].message}
                    </span>
                  )}
                </div>
              </div>

              <div className="form-group" style={{ marginTop: '8px' }}>
                <label htmlFor="set_DISCORD_MOVIE_WEBHOOK_URL">Discord Movie Webhook URL</label>
                <input
                  type="password"
                  id="set_DISCORD_MOVIE_WEBHOOK_URL"
                  name="DISCORD_MOVIE_WEBHOOK_URL"
                  className="brutal-input"
                  value={config.DISCORD_MOVIE_WEBHOOK_URL || ''}
                  onChange={onChange}
                  placeholder="https://discord.com/api/webhooks/..."
                />
                <div style={{ display: 'flex', gap: '8px', alignItems: 'center', marginTop: '6px', flexWrap: 'wrap' }}>
                  <button
                    type="button"
                    className="brutal-btn"
                    style={{ padding: '4px 12px', fontSize: '0.85rem', cursor: config.DISCORD_MOVIE_WEBHOOK_URL ? 'pointer' : 'not-allowed' }}
                    disabled={!config.DISCORD_MOVIE_WEBHOOK_URL || testingPlatform === 'discord_movie'}
                    onClick={() => handleTestWebhook('discord', config.DISCORD_MOVIE_WEBHOOK_URL || '', 'discord_movie')}
                  >
                    {testingPlatform === 'discord_movie' ? 'Testing...' : 'Test Movie Webhook'}
                  </button>
                  {testStatus['discord_movie'] && (
                    <span
                      style={{
                        fontSize: '0.85rem',
                        fontWeight: 600,
                        color: testStatus['discord_movie'].type === 'success' ? '#16a34a' : '#dc2626'
                      }}
                    >
                      {testStatus['discord_movie'].type === 'success' ? '✅ ' : '❌ '}
                      {testStatus['discord_movie'].message}
                    </span>
                  )}
                </div>
              </div>
            </details>

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
              <label htmlFor="set_SLACK_WEBHOOK_URL">Default Slack Webhook URL</label>
              <input
                type="password"
                id="set_SLACK_WEBHOOK_URL"
                name="SLACK_WEBHOOK_URL"
                className="brutal-input"
                value={config.SLACK_WEBHOOK_URL || ''}
                onChange={onChange}
                placeholder="https://hooks.slack.com/services/..."
              />
              <div style={{ display: 'flex', gap: '8px', alignItems: 'center', marginTop: '6px', flexWrap: 'wrap' }}>
                <button
                  type="button"
                  className="brutal-btn"
                  style={{ padding: '4px 12px', fontSize: '0.85rem', cursor: config.SLACK_WEBHOOK_URL ? 'pointer' : 'not-allowed' }}
                  disabled={!config.SLACK_WEBHOOK_URL || testingPlatform === 'slack'}
                  onClick={() => handleTestWebhook('slack', config.SLACK_WEBHOOK_URL || '', 'slack')}
                >
                  {testingPlatform === 'slack' ? 'Testing...' : 'Test Slack Webhook'}
                </button>
                {testStatus['slack'] && (
                  <span
                    style={{
                      fontSize: '0.85rem',
                      fontWeight: 600,
                      color: testStatus['slack'].type === 'success' ? '#16a34a' : '#dc2626'
                    }}
                  >
                    {testStatus['slack'].type === 'success' ? '✅ ' : '❌ '}
                    {testStatus['slack'].message}
                  </span>
                )}
              </div>
            </div>

            {/* Separate Slack Channels */}
            <details
              style={{ marginTop: '12px', marginBottom: '12px', border: '1px dashed var(--color-border)', borderRadius: '4px', padding: '8px 12px' }}
              open={!!(config.SLACK_TV_WEBHOOK_URL || config.SLACK_MOVIE_WEBHOOK_URL)}
            >
              <summary style={{ cursor: 'pointer', fontWeight: 600, fontSize: '0.9rem' }}>
                📺 🎬 Separate TV and Movie Channels (Optional)
              </summary>
              <p style={{ fontSize: '0.8rem', color: 'var(--color-text-muted)', margin: '6px 0 12px 0' }}>
                Route TV shows (Sonarr) and movies (Radarr) to dedicated Slack channels. If not specified, notifications fall back to the default webhook above.
              </p>
              <div className="form-group">
                <label htmlFor="set_SLACK_TV_WEBHOOK_URL">Slack TV Webhook URL</label>
                <input
                  type="password"
                  id="set_SLACK_TV_WEBHOOK_URL"
                  name="SLACK_TV_WEBHOOK_URL"
                  className="brutal-input"
                  value={config.SLACK_TV_WEBHOOK_URL || ''}
                  onChange={onChange}
                  placeholder="https://hooks.slack.com/services/..."
                />
                <div style={{ display: 'flex', gap: '8px', alignItems: 'center', marginTop: '6px', flexWrap: 'wrap' }}>
                  <button
                    type="button"
                    className="brutal-btn"
                    style={{ padding: '4px 12px', fontSize: '0.85rem', cursor: config.SLACK_TV_WEBHOOK_URL ? 'pointer' : 'not-allowed' }}
                    disabled={!config.SLACK_TV_WEBHOOK_URL || testingPlatform === 'slack_tv'}
                    onClick={() => handleTestWebhook('slack', config.SLACK_TV_WEBHOOK_URL || '', 'slack_tv')}
                  >
                    {testingPlatform === 'slack_tv' ? 'Testing...' : 'Test TV Webhook'}
                  </button>
                  {testStatus['slack_tv'] && (
                    <span
                      style={{
                        fontSize: '0.85rem',
                        fontWeight: 600,
                        color: testStatus['slack_tv'].type === 'success' ? '#16a34a' : '#dc2626'
                      }}
                    >
                      {testStatus['slack_tv'].type === 'success' ? '✅ ' : '❌ '}
                      {testStatus['slack_tv'].message}
                    </span>
                  )}
                </div>
              </div>

              <div className="form-group" style={{ marginTop: '8px' }}>
                <label htmlFor="set_SLACK_MOVIE_WEBHOOK_URL">Slack Movie Webhook URL</label>
                <input
                  type="password"
                  id="set_SLACK_MOVIE_WEBHOOK_URL"
                  name="SLACK_MOVIE_WEBHOOK_URL"
                  className="brutal-input"
                  value={config.SLACK_MOVIE_WEBHOOK_URL || ''}
                  onChange={onChange}
                  placeholder="https://hooks.slack.com/services/..."
                />
                <div style={{ display: 'flex', gap: '8px', alignItems: 'center', marginTop: '6px', flexWrap: 'wrap' }}>
                  <button
                    type="button"
                    className="brutal-btn"
                    style={{ padding: '4px 12px', fontSize: '0.85rem', cursor: config.SLACK_MOVIE_WEBHOOK_URL ? 'pointer' : 'not-allowed' }}
                    disabled={!config.SLACK_MOVIE_WEBHOOK_URL || testingPlatform === 'slack_movie'}
                    onClick={() => handleTestWebhook('slack', config.SLACK_MOVIE_WEBHOOK_URL || '', 'slack_movie')}
                  >
                    {testingPlatform === 'slack_movie' ? 'Testing...' : 'Test Movie Webhook'}
                  </button>
                  {testStatus['slack_movie'] && (
                    <span
                      style={{
                        fontSize: '0.85rem',
                        fontWeight: 600,
                        color: testStatus['slack_movie'].type === 'success' ? '#16a34a' : '#dc2626'
                      }}
                    >
                      {testStatus['slack_movie'].type === 'success' ? '✅ ' : '❌ '}
                      {testStatus['slack_movie'].message}
                    </span>
                  )}
                </div>
              </div>
            </details>

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
