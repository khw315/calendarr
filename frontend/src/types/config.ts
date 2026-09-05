export interface CalendarUrlItem {
  url: string
  type: string
}

export interface ConfigState {
  APP_LANGUAGE?: string
  USE_DISCORD?: boolean
  DISCORD_WEBHOOK_URL?: string
  DISCORD_MENTION_ROLE_ID?: string
  DISCORD_HIDE_MENTION_INSTRUCTIONS?: boolean
  DISCORD_TIMESTAMP_STYLE?: string
  ENABLE_CUSTOM_DISCORD_FOOTER?: boolean
  USE_SLACK?: boolean
  SLACK_WEBHOOK_URL?: string
  ENABLE_CUSTOM_SLACK_FOOTER?: boolean
  CALENDAR_URLS?: CalendarUrlItem[]
  PASSED_EVENT_HANDLING?: string
  DEDUPLICATE_EVENTS?: boolean
  USE_24_HOUR?: boolean
  ADD_LEADING_ZERO?: boolean
  DISPLAY_TIME?: boolean
  SHOW_DATE_RANGE?: boolean
  SHOW_TIMEZONE_IN_SUBHEADER?: boolean
  TZ?: string
  SCHEDULE_TYPE?: string
  SCHEDULE_DAY?: string
  RUN_TIME?: string
  CRON_SCHEDULE?: string
  RUN_ON_STARTUP?: boolean
  DEBUG?: boolean
  HTTP_TIMEOUT?: number
  LOG_MAX_SIZE_MB?: number
  LOG_BACKUP_COUNT?: number
  [key: string]: unknown
}
