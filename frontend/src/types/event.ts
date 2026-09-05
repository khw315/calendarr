export type EventType = 'tv' | 'movie'

export interface EventItem {
  title: string
  type: EventType | string
  start_time: string
  end_time?: string
  date: string
  timestamp?: number
  end_timestamp?: number
  description?: string
  summary?: string
  source_type?: string
  is_bulk?: boolean
  episode_count?: number
  is_series_finale?: boolean
}

export interface DayGroup {
  day_name: string
  date: string
  events: EventItem[]
}

export interface ScheduleDTO {
  type: string
  nextRun: string
  timezone: string
}

export const ALLOWED_RANGES = ['1', '3', '7', '14'] as const
export type AllowedRange = typeof ALLOWED_RANGES[number]

export const ALLOWED_PAST_SORTS = ['newest', 'oldest'] as const
export type AllowedPastSort = typeof ALLOWED_PAST_SORTS[number]
