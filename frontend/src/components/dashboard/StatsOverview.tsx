import React from 'react'
import { AllowedRange, ScheduleDTO } from '@/types/event'
import { formatCountdown } from '@/lib/utils'

interface StatsOverviewProps {
  range: AllowedRange
  totalTv: number | string
  totalMovies: number | string
  schedule: ScheduleDTO
  currentTime?: number
  triggering: boolean
  triggerStatus: string | null
  onTrigger: () => void
}

export default function StatsOverview({
  range,
  totalTv,
  totalMovies,
  schedule,
  currentTime,
  triggering,
  triggerStatus,
  onTrigger
}: StatsOverviewProps) {
  const getRangeLabel = () => {
    switch (range) {
      case '1':
        return 'TODAY'
      case '3':
        return 'NEXT 3 DAYS'
      case '7':
        return 'NEXT 7 DAYS'
      case '14':
        return 'NEXT 14 DAYS'
      default:
        return 'UPCOMING'
    }
  }

  const renderStatValue = () => {
    if (totalTv === '-' || totalMovies === '-') return '-'

    const tv = Number(totalTv)
    const movies = Number(totalMovies)

    if (tv === 0 && movies === 0) return '0'

    const tvText = tv > 0 ? `${tv} Show${tv > 1 ? 's' : ''}` : ''
    const movieText = movies > 0 ? `${movies} Movie${movies > 1 ? 's' : ''}` : ''
    const separator = tv > 0 && movies > 0 ? ' and ' : ''

    return `${tvText}${separator}${movieText}`
  }

  const nextRunCountdown = formatCountdown(schedule.nextRun, schedule.timezone, currentTime)
  const nextRunTooltip = schedule.nextRun && schedule.nextRun !== '-' && schedule.nextRun !== 'Not scheduled'
    ? `${schedule.nextRun.replace('Z', '')} (${schedule.timezone || 'UTC'})`
    : undefined

  return (
    <div className="stats-grid">
      {/* 1. Total Shows/Movies */}
      <div className="stat-card brutal-card brutal-interactive">
        <div className="stat-icon stat-icon-primary">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <rect x="3" y="6" width="18" height="15" rx="2" stroke="currentColor" strokeWidth="2" />
            <line x1="3" y1="10" x2="21" y2="10" stroke="currentColor" strokeWidth="2" />
          </svg>
        </div>
        <div className="stat-content">
          <div className="stat-label">{getRangeLabel()}</div>
          <div className="stat-value" style={{ fontSize: '1.2rem' }}>
            {renderStatValue()}
          </div>
        </div>
      </div>

      {/* 2. Next Scheduled Run */}
      <div
        className="stat-card brutal-card brutal-interactive"
        title={nextRunTooltip}
      >
        <div className="stat-icon stat-icon-success">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <circle cx="12" cy="12" r="9" stroke="currentColor" strokeWidth="2" />
            <path d="M12 7V12L15 15" stroke="currentColor" strokeWidth="2" strokeLinecap="round" />
          </svg>
        </div>
        <div className="stat-content">
          <div className="stat-label">Next Run</div>
          <div className="stat-value">{nextRunCountdown}</div>
        </div>
      </div>

      {/* 3. Schedule Type */}
      <div className="stat-card brutal-card brutal-interactive">
        <div className="stat-icon stat-icon-info">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z" stroke="currentColor" strokeWidth="2" strokeLinejoin="round" />
          </svg>
        </div>
        <div className="stat-content">
          <div className="stat-label">Schedule Type</div>
          <div className="stat-value">{schedule.type}</div>
        </div>
      </div>

      {/* 4. Manual Trigger Action */}
      <div className="stat-card brutal-card trigger-card">
        <button
          className="trigger-btn"
          onClick={onTrigger}
          disabled={triggering}
        >
          {triggering ? (
            <div className="btn-spinner" />
          ) : (
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M5 3l14 9-14 9V3z" fill="currentColor" />
            </svg>
          )}
          <span>{triggering ? "Running..." : "Run Now"}</span>
        </button>
        <div className={`trigger-status ${triggerStatus ? (triggerStatus.includes('success') ? 'success' : 'error') : ''}`}>
          {triggerStatus}
        </div>
      </div>
    </div>
  )
}
