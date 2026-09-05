import { DayGroup, ScheduleDTO } from '@/types/event'
import { ConfigState } from '@/types/config'

const API_BASE = ''

export const apiService = {
  async getSchedule(): Promise<ScheduleDTO> {
    const res = await fetch(`${API_BASE}/api/schedule`)
    if (!res.ok) {
      throw new Error(`Failed to fetch schedule: ${res.statusText}`)
    }
    const data = await res.json()
    return {
      type: data.schedule_type || 'N/A',
      nextRun: data.next_run ? data.next_run : 'Not scheduled',
      timezone: data.timezone || 'UTC'
    }
  },

  async getReleases(days: string): Promise<DayGroup[]> {
    const res = await fetch(`${API_BASE}/api/releases?days=${encodeURIComponent(days)}`)
    if (!res.ok) {
      throw new Error(`Failed to fetch releases: ${res.statusText}`)
    }
    const data = await res.json()
    return data.days || []
  },

  async getPastReleases(days: string): Promise<DayGroup[]> {
    const res = await fetch(`${API_BASE}/api/past-releases?days=${encodeURIComponent(days)}`)
    if (!res.ok) {
      throw new Error(`Failed to fetch past releases: ${res.statusText}`)
    }
    const data = await res.json()
    return data.days || []
  },

  async getConfig(): Promise<ConfigState> {
    const res = await fetch(`${API_BASE}/api/config`)
    if (!res.ok) {
      throw new Error(`Failed to fetch config: ${res.statusText}`)
    }
    return res.json()
  },

  async saveConfig(config: ConfigState): Promise<{ success: boolean; message?: string }> {
    const res = await fetch(`${API_BASE}/api/config`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(config)
    })
    const data = await res.json()
    if (!res.ok) {
      throw new Error(data.error || data.message || `Failed to save config: ${res.statusText}`)
    }
    return data
  },

  async triggerSync(): Promise<{ success: boolean; message?: string }> {
    const res = await fetch(`${API_BASE}/api/trigger`, {
      method: 'POST'
    })
    const data = await res.json()
    if (!res.ok) {
      throw new Error(data.error || data.message || 'Trigger failed')
    }
    return data
  },

  async getLanguages(): Promise<string[]> {
    const res = await fetch(`${API_BASE}/api/languages`)
    if (!res.ok) {
      return ["EN", "FR", "ID", "JA", "KO"]
    }
    const data = await res.json()
    return data.languages || ["EN", "FR", "ID", "JA", "KO"]
  },

  async getTimezones(): Promise<string[]> {
    const res = await fetch(`${API_BASE}/api/timezones`)
    if (!res.ok) {
      return ["UTC"]
    }
    const data = await res.json()
    return data.timezones || ["UTC"]
  }
}
