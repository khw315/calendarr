'use client'

import React from 'react'
import Header from '@/components/Header'
import Footer from '@/components/Footer'
import { useDashboard } from '@/hooks/useDashboard'
import StatsOverview from '@/components/dashboard/StatsOverview'
import DayGroupList from '@/components/events/DayGroupList'
import PastReleasesSection from '@/components/events/PastReleasesSection'
import LoadingSpinner from '@/components/common/LoadingSpinner'
import { ALLOWED_RANGES, AllowedRange } from '@/types/event'

export default function DashboardPage() {
  const {
    days,
    pastDays,
    totalTv,
    totalMovies,
    schedule,
    loading,
    range,
    pastRange,
    pastSort,
    showPast,
    triggering,
    triggerStatus,
    currentTime,
    setRange,
    setPastRange,
    setPastSort,
    setShowPast,
    handleTrigger,
    refreshAll
  } = useDashboard()

  const getSectionTitle = () => {
    switch (range) {
      case '1':
        return 'Releases Today'
      case '3':
        return 'Releases Next 3 Days'
      case '7':
        return 'Releases Next 7 Days'
      case '14':
        return 'Releases Next 14 Days'
      default:
        return 'Upcoming Releases'
    }
  }

  return (
    <div id="root">
      <Header activePage={'dashboard'} />

      <main className="main">
        <div className="container" id="dashboard">
          {/* Top Metrics & Action Row */}
          <StatsOverview
            range={range}
            totalTv={totalTv}
            totalMovies={totalMovies}
            schedule={schedule}
            currentTime={currentTime}
            triggering={triggering}
            triggerStatus={triggerStatus}
            onTrigger={handleTrigger}
          />

          {/* Upcoming Releases Section */}
          <section className="section" id="events">
            <div className="section-header">
              <h2>{getSectionTitle()}</h2>
              <div className="section-actions">
                <select
                  className="range-selector"
                  value={range}
                  onChange={e => {
                    if (ALLOWED_RANGES.includes(e.target.value as AllowedRange)) {
                      setRange(e.target.value)
                    }
                  }}
                >
                  <option value="1">Today</option>
                  <option value="3">Next 3 Days</option>
                  <option value="7">Next 7 Days</option>
                  <option value="14">Next 14 Days</option>
                </select>
                <button className="refresh-btn" onClick={refreshAll}>
                  <svg width="18" height="18" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path
                      d="M21.5 2v6h-6M2.5 22v-6h6M2 11.5a10 10 0 0 1 18.8-4.3M22 12.5a10 10 0 0 1-18.8 4.2"
                      stroke="currentColor"
                      strokeWidth="2"
                      strokeLinecap="round"
                      strokeLinejoin="round"
                    />
                  </svg>
                  Refresh
                </button>
              </div>
            </div>

            <div className="events-container">
              {loading ? (
                <LoadingSpinner label="Loading releases..." />
              ) : (
                <DayGroupList days={days} currentTime={currentTime} />
              )}
            </div>
          </section>

          {/* Past Releases Collapsible Section */}
          <PastReleasesSection
            pastDays={pastDays}
            showPast={showPast}
            pastRange={pastRange}
            pastSort={pastSort}
            onToggleShowPast={() => setShowPast(!showPast)}
            onSelectPastRange={setPastRange}
            onSelectPastSort={setPastSort}
          />
        </div>
      </main>

      <Footer />
    </div>
  )
}
