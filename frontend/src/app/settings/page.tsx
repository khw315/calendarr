'use client'

import React from 'react'
import Header from '@/components/Header'
import Footer from '@/components/Footer'
import LoadingSpinner from '@/components/common/LoadingSpinner'
import CalendarFeedManager from '@/components/settings/CalendarFeedManager'
import NotificationSettings from '@/components/settings/NotificationSettings'
import DisplaySettings from '@/components/settings/DisplaySettings'
import ScheduleSettings from '@/components/settings/ScheduleSettings'
import AdvancedSettings from '@/components/settings/AdvancedSettings'
import { useSettings } from '@/hooks/useSettings'

export default function SettingsPage() {
  const {
    config,
    loading,
    saving,
    toast,
    languages,
    timezones,
    themeMode,
    setThemeMode,
    handleChange,
    handleAddCalendarUrl,
    handleRemoveCalendarUrl,
    handleCalendarUrlChange,
    handleSave,
    handleDiscard
  } = useSettings()

  if (loading) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', marginTop: '50px' }}>
        <LoadingSpinner label="Loading settings..." />
      </div>
    )
  }

  return (
    <>
      <Header activePage="settings" />

      <main className="main">
        <div className="container" id="settingsView">
          <div className="section-header">
            <h2>Configuration Settings</h2>
          </div>

          <div className="section">
            <form id="settingsForm" onSubmit={handleSave}>
              {/* 1. Calendar Feeds */}
              <CalendarFeedManager
                feeds={config.CALENDAR_URLS || []}
                onAddFeed={handleAddCalendarUrl}
                onRemoveFeed={handleRemoveCalendarUrl}
                onUpdateFeed={handleCalendarUrlChange}
              />

              {/* 2. Notification Integrations (Discord & Slack) */}
              <NotificationSettings
                config={config}
                onChange={handleChange}
              />

              {/* 3. Display & Formatting */}
              <DisplaySettings
                config={config}
                languages={languages}
                timezones={timezones}
                themeMode={themeMode}
                onThemeChange={setThemeMode}
                onChange={handleChange}
              />

              {/* 4. Scheduling */}
              <ScheduleSettings
                config={config}
                onChange={handleChange}
              />

              {/* 5. Advanced Settings */}
              <AdvancedSettings
                config={config}
                onChange={handleChange}
              />

              {/* Actions */}
              <div
                className="settings-actions"
                style={{
                  display: 'flex',
                  justifyContent: 'flex-end',
                  gap: 'var(--spacing-sm)',
                  marginTop: 'var(--spacing-lg)'
                }}
              >
                <button
                  type="button"
                  className="brutal-btn brutal-btn-outline"
                  onClick={handleDiscard}
                >
                  Discard Changes
                </button>
                <button
                  type="submit"
                  className="brutal-btn brutal-btn-primary"
                  disabled={saving}
                >
                  {saving ? 'Saving...' : 'Save Changes'}
                </button>
              </div>
            </form>
          </div>
        </div>
      </main>

      <Footer />

      {/* Toast Notification */}
      {toast && (
        <div className="toast-container">
          <div className={`toast ${toast.type}`}>
            {toast.message}
          </div>
        </div>
      )}
    </>
  )
}
