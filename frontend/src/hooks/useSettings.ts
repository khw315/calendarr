import { useState, useEffect, useCallback } from 'react'
import { ConfigState, CalendarUrlItem } from '@/types/config'
import { useTheme, ThemeMode } from './useTheme'

const API_BASE = ''

export function useSettings() {
  const [config, setConfig] = useState<ConfigState>({
    APP_LANGUAGE: "EN",
    USE_DISCORD: false,
    DISCORD_WEBHOOK_URL: "",
    DISCORD_MENTION_ROLE_ID: "",
    DISCORD_HIDE_MENTION_INSTRUCTIONS: false,
    DISCORD_TIMESTAMP_STYLE: "R",
    ENABLE_CUSTOM_DISCORD_FOOTER: false,
    USE_SLACK: false,
    SLACK_WEBHOOK_URL: "",
    ENABLE_CUSTOM_SLACK_FOOTER: false,
    CALENDAR_URLS: [],
    PASSED_EVENT_HANDLING: "DISPLAY",
    DEDUPLICATE_EVENTS: false,
    USE_24_HOUR: false,
    ADD_LEADING_ZERO: false,
    DISPLAY_TIME: true,
    SHOW_DATE_RANGE: false,
    SHOW_TIMEZONE_IN_SUBHEADER: false,
    TZ: "UTC",
    SCHEDULE_TYPE: "DAILY",
    SCHEDULE_DAY: "0",
    RUN_TIME: "09:00",
    CRON_SCHEDULE: "",
    RUN_ON_STARTUP: false,
    DEBUG: false,
    HTTP_TIMEOUT: 30,
    LOG_MAX_SIZE_MB: 10,
    LOG_BACKUP_COUNT: 5
  })
  const [originalConfig, setOriginalConfig] = useState<ConfigState>({})
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [toast, setToast] = useState<{ message: string; type: 'success' | 'error' } | null>(null)
  const [timezones, setTimezones] = useState<{ iana: string; label: string }[]>([])
  const [languages, setLanguages] = useState<{ code: string; name: string }[]>([])
  const { mode: themeMode, setMode: setThemeMode } = useTheme()

  const showToast = useCallback((message: string, type: 'success' | 'error' = 'success') => {
    setToast({ message, type })
    setTimeout(() => setToast(null), 3000)
  }, [])

  const fetchConfig = useCallback(async () => {
    setLoading(true)
    try {
      const res = await fetch(`${API_BASE}/api/config`)
      if (res.ok) {
        const data = await res.json()
        setConfig(data)
        setOriginalConfig(JSON.parse(JSON.stringify(data)))
      }
    } catch {
      showToast("Failed to load settings.", "error")
    } finally {
      setLoading(false)
    }
  }, [showToast])

  const fetchTimezones = useCallback(async () => {
    try {
      const res = await fetch(`${API_BASE}/api/timezones`)
      if (res.ok) {
        const data: Record<string, string> = await res.json()
        const sorted = Object.entries(data)
          .map(([iana, label]) => ({ iana, label }))
          .sort((a, b) => a.iana.localeCompare(b.iana))
        setTimezones(sorted)
      }
    } catch {
      try {
        const tzs = Intl.supportedValuesOf('timeZone')
        setTimezones(tzs.map(iana => ({ iana, label: iana })))
      } catch {
        setTimezones([])
      }
    }
  }, [])

  const fetchLanguages = useCallback(async () => {
    try {
      const res = await fetch(`${API_BASE}/api/languages`)
      if (res.ok) {
        const data = await res.json()
        setLanguages(data)
      }
    } catch {
      // ignore
    }
  }, [])

  useEffect(() => {
    fetchConfig()
    fetchLanguages()
    fetchTimezones()
  }, [fetchConfig, fetchLanguages, fetchTimezones])

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const target = e.target as HTMLInputElement
    const value = target.type === 'checkbox' ? target.checked : target.value
    setConfig(prev => ({
      ...prev,
      [target.name]: value
    }))
  }

  const handleAddCalendarUrl = () => {
    setConfig(prev => ({
      ...prev,
      CALENDAR_URLS: [...(prev.CALENDAR_URLS || []), { url: "", type: "tv" }]
    }))
  }

  const handleRemoveCalendarUrl = (index: number) => {
    setConfig(prev => {
      const newUrls = [...(prev.CALENDAR_URLS || [])]
      newUrls.splice(index, 1)
      return { ...prev, CALENDAR_URLS: newUrls }
    })
  }

  const handleCalendarUrlChange = (index: number, field: 'url' | 'type', value: string) => {
    setConfig(prev => {
      const newUrls = [...(prev.CALENDAR_URLS || [])]
      newUrls[index] = { ...newUrls[index], [field]: value }
      return { ...prev, CALENDAR_URLS: newUrls }
    })
  }

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault()
    setSaving(true)

    const changesToSubmit: Record<string, unknown> = {}
    for (const [key, value] of Object.entries(config)) {
      let isChanged = false
      if (key === 'CALENDAR_URLS') {
        const filteredUrls = ((value as CalendarUrlItem[]) || []).filter(u => u.url && u.url.trim() !== '')
        isChanged = JSON.stringify(filteredUrls) !== JSON.stringify(originalConfig[key] || [])
        if (isChanged) changesToSubmit[key] = filteredUrls
        continue
      } else {
        if (typeof value === 'boolean') {
          const origBool = originalConfig[key] === true || originalConfig[key] === 'true'
          isChanged = value !== origBool
        } else {
          const valStr = String(value || '').trim()
          const origStr = String(originalConfig[key] || '').trim()
          isChanged = valStr !== origStr
        }
      }
      if (isChanged) {
        changesToSubmit[key] = key === 'SCHEDULE_DAY' ? String(value) : value
      }
    }

    if (Object.keys(changesToSubmit).length === 0) {
      showToast("No changes detected to save.")
      setSaving(false)
      return
    }

    try {
      const res = await fetch(`${API_BASE}/api/config`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(changesToSubmit)
      })
      const data = await res.json()
      if (res.ok && (data.success || data.status === 'success')) {
        showToast("Settings saved successfully!", 'success')
        setOriginalConfig(prev => ({ ...prev, ...changesToSubmit }))
      } else {
        showToast(data.error || "Failed to save settings", 'error')
      }
    } catch {
      showToast("Failed to communicate with server.", 'error')
    } finally {
      setSaving(false)
    }
  }

  const handleDiscard = () => {
    fetchConfig()
    showToast("Changes discarded.")
  }

  return {
    config,
    loading,
    saving,
    toast,
    languages,
    timezones,
    themeMode,
    setThemeMode: (mode: ThemeMode) => setThemeMode(mode),
    handleChange,
    handleAddCalendarUrl,
    handleRemoveCalendarUrl,
    handleCalendarUrlChange,
    handleSave,
    handleDiscard
  }
}
