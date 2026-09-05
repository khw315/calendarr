import React from 'react'
import { CalendarUrlItem } from '@/types/config'

interface CalendarFeedManagerProps {
  feeds: CalendarUrlItem[]
  onAddFeed: () => void
  onRemoveFeed: (index: number) => void
  onUpdateFeed: (index: number, field: keyof CalendarUrlItem, value: string) => void
}

export default function CalendarFeedManager({
  feeds,
  onAddFeed,
  onRemoveFeed,
  onUpdateFeed
}: CalendarFeedManagerProps) {
  return (
    <details className="settings-group" name="settings-accordion">
      <summary>
        <h3>Source Configuration</h3>
        <span className="details-arrow">▼</span>
      </summary>
      <div className="form-group">
        <label>Calendar URLs (<span id="calendarUrlCount">{feeds.length}</span>)</label>
        <div className="calendar-url-container">
          {feeds.length === 0 ? (
            <p style={{ color: 'var(--color-text-dim)', fontSize: '0.9em', marginBottom: '5px' }}>
              No calendar URLs added. Add one to get started.
            </p>
          ) : (
            feeds.map((urlObj, index) => (
              <div key={index} className="calendar-url-item">
                <div className="calendar-url-inputs">
                  <input
                    type="url"
                    className="brutal-input url-val"
                    placeholder="https://..."
                    value={urlObj.url || ''}
                    onChange={(e) => onUpdateFeed(index, 'url', e.target.value)}
                    required
                  />
                  <select
                    className="brutal-select type-val"
                    value={urlObj.type || 'tv'}
                    onChange={(e) => onUpdateFeed(index, 'type', e.target.value)}
                  >
                    <option value="tv">TV Show (Sonarr)</option>
                    <option value="movie">Movie (Radarr)</option>
                  </select>
                </div>
                <button
                  type="button"
                  className="brutal-btn brutal-btn-danger remove-url-btn"
                  onClick={() => onRemoveFeed(index)}
                  title="Remove"
                >
                  &times;
                </button>
              </div>
            ))
          )}
        </div>
        <button
          type="button"
          className="brutal-btn mt-none"
          onClick={onAddFeed}
          style={{ marginTop: '0' }}
        >
          + Add Calendar
        </button>
      </div>
    </details>
  )
}
