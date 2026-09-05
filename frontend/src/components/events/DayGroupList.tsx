import React from 'react'
import { DayGroup } from '@/types/event'
import EventCard from './EventCard'
import EmptyState from '@/components/common/EmptyState'

interface DayGroupListProps {
  days: DayGroup[]
  currentTime: number
  emptyMessage?: string
}

export default function DayGroupList({
  days,
  currentTime,
  emptyMessage = 'No upcoming events found'
}: DayGroupListProps) {
  if (days.length === 0) {
    return <EmptyState message={emptyMessage} />
  }

  return (
    <>
      {days.map((day, i) => (
        <div key={day.date || i} className="day-group">
          <div className="day-header">{day.day_name}</div>
          <div className="events-list">
            {day.events.map((ev, j) => (
              <EventCard
                key={`${ev.title}-${ev.start_time}-${j}`}
                event={ev}
                currentTime={currentTime}
                isPast={false}
              />
            ))}
          </div>
        </div>
      ))}
    </>
  )
}
