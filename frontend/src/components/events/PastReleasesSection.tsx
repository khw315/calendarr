import React from 'react'
import { DayGroup, AllowedRange, AllowedPastSort, ALLOWED_RANGES, ALLOWED_PAST_SORTS } from '@/types/event'
import EventCard from './EventCard'
import EmptyState from '@/components/common/EmptyState'

interface PastReleasesSectionProps {
  pastDays: DayGroup[]
  showPast: boolean
  pastRange: AllowedRange
  pastSort: AllowedPastSort
  onToggleShowPast: () => void
  onSelectPastRange: (range: string) => void
  onSelectPastSort: (sort: string) => void
}

export default function PastReleasesSection({
  pastDays,
  showPast,
  pastRange,
  pastSort,
  onToggleShowPast,
  onSelectPastRange,
  onSelectPastSort
}: PastReleasesSectionProps) {
  const sortedDays = [...pastDays].sort((a, b) => {
    return pastSort === 'newest' ? b.date.localeCompare(a.date) : a.date.localeCompare(b.date)
  })

  return (
    <section className="section" id="past-events">
      <div className="section-header">
        <h2>Past Releases</h2>
        <div className="section-actions">
          {showPast && (
            <>
              <select
                className="range-selector"
                value={pastSort}
                onChange={e => {
                  if (ALLOWED_PAST_SORTS.includes(e.target.value as AllowedPastSort)) {
                    onSelectPastSort(e.target.value)
                  }
                }}
              >
                <option value="newest">Newest to Oldest</option>
                <option value="oldest">Oldest to Newest</option>
              </select>
              <select
                className="range-selector"
                value={pastRange}
                onChange={e => {
                  if (ALLOWED_RANGES.includes(e.target.value as AllowedRange)) {
                    onSelectPastRange(e.target.value)
                  }
                }}
              >
                <option value="1">Last 1 Day</option>
                <option value="3">Last 3 Days</option>
                <option value="7">Last 7 Days</option>
                <option value="14">Last 14 Days</option>
              </select>
            </>
          )}
          <button
            className="toggle-btn"
            onClick={onToggleShowPast}
            aria-label="Toggle Past Releases"
          >
            <svg
              className="toggle-icon"
              style={{ transform: showPast ? 'rotate(180deg)' : 'rotate(0deg)' }}
              width="20"
              height="20"
              viewBox="0 0 24 24"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                d="M19 9l-7 7-7-7"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              />
            </svg>
          </button>
        </div>
      </div>

      {showPast && (
        <div className="past-content">
          <div className="events-container">
            {sortedDays.length === 0 ? (
              <EmptyState
                message={`No past releases in the last ${pastRange} ${pastRange === '1' ? 'day' : 'days'}`}
              />
            ) : (
              sortedDays.map((day, i) => (
                <div key={day.date || i} className="day-group">
                  <div className="day-header">{day.day_name}</div>
                  <div className="events-list">
                    {[...day.events]
                      .sort((a, b) => {
                        const tA = a.timestamp || 0
                        const tB = b.timestamp || 0
                        return pastSort === 'newest' ? tB - tA : tA - tB
                      })
                      .map((ev, j) => (
                        <EventCard
                          key={`${ev.title}-${ev.start_time}-${j}`}
                          event={ev}
                          isPast={true}
                        />
                      ))}
                  </div>
                </div>
              ))
            )}
          </div>
        </div>
      )}
    </section>
  )
}
