import React from 'react'
import { EventItem } from '@/types/event'

interface EventCardProps {
  event: EventItem
  currentTime?: number
  isPast?: boolean
}

export default function EventCard({ event, currentTime = 0, isPast = false }: EventCardProps) {
  const title = event.title || event.summary || 'Untitled'
  const rawType = event.type || event.source_type || 'tv'
  const isTv = String(rawType).toLowerCase() === 'tv'

  // Time display formatting (support both HH:mm and ISO strings)
  let timeDisplay = event.start_time || ''
  if (timeDisplay.includes('T')) {
    try {
      const d = new Date(timeDisplay)
      const hours = String(d.getHours()).padStart(2, '0')
      const mins = String(d.getMinutes()).padStart(2, '0')
      timeDisplay = `${hours}:${mins}`
    } catch {
      // keep fallback timeDisplay
    }
  }

  const isAiring = !isPast && Boolean(
    event.timestamp &&
    event.end_timestamp &&
    event.end_timestamp > event.timestamp &&
    currentTime >= event.timestamp &&
    currentTime < event.end_timestamp
  )

  const isAired = isPast || Boolean(
    (event.end_timestamp && event.timestamp && event.end_timestamp > event.timestamp)
      ? currentTime >= event.end_timestamp
      : (event.timestamp && currentTime >= event.timestamp)
  )

  const startsIn = event.timestamp ? event.timestamp - currentTime : -1
  const isStartingSoon = !isPast && !isAiring && !isAired && startsIn > 0 && startsIn <= 3600

  const cardClasses = [
    'event-card',
    'brutal-card',
    `event-${isTv ? 'tv' : 'movie'}`,
    isPast ? 'event-past' : '',
    isAiring ? 'airing' : '',
    isAired && !isPast ? 'aired' : ''
  ].filter(Boolean).join(' ')

  return (
    <div className={cardClasses}>
      {isAiring && <span className="airing-badge">ON AIR</span>}
      {isAired && <span className="aired-badge">AIRED</span>}

      <span className="event-type">{isTv ? 'TV' : 'Movie'}</span>

      <div className="event-title">
        {title}
        {event.is_bulk && (
          <span className="bulk-badge">{event.episode_count} Episodes</span>
        )}
        {event.is_series_finale && (
          <span className="finale-badge">Series End</span>
        )}
      </div>

      {timeDisplay && (
        <div className="event-time" data-original-time={event.start_time}>
          {isStartingSoon ? (
            <span className="time-text countdown-active">
              Starts in {Math.ceil(startsIn / 60)}m
            </span>
          ) : isAiring && event.end_time ? (
            <span className="time-text">{event.end_time}</span>
          ) : (
            <span className="time-text">{timeDisplay}</span>
          )}
        </div>
      )}
    </div>
  )
}
