import { useState, useEffect, useCallback } from 'react'
import {
  DayGroup,
  EventItem,
  ScheduleDTO,
  AllowedRange,
  AllowedPastSort,
  ALLOWED_RANGES,
  ALLOWED_PAST_SORTS
} from '@/types/event'
import { apiService } from '@/services/api'

export function useDashboard() {
  const [days, setDays] = useState<DayGroup[]>([])
  const [pastDays, setPastDays] = useState<DayGroup[]>([])
  const [totalTv, setTotalTv] = useState<number | string>('-')
  const [totalMovies, setTotalMovies] = useState<number | string>('-')
  const [schedule, setSchedule] = useState<ScheduleDTO>({
    type: '-',
    nextRun: '-',
    timezone: '-'
  })
  const [loading, setLoading] = useState(true)
  const [range, setRangeState] = useState<AllowedRange>('1')
  const [pastRange, setPastRangeState] = useState<AllowedRange>('3')
  const [pastSort, setPastSortState] = useState<AllowedPastSort>('newest')
  const [showPast, setShowPast] = useState(false)
  const [triggering, setTriggering] = useState(false)
  const [triggerStatus, setTriggerStatus] = useState<string | null>(null)
  const [currentTime, setCurrentTime] = useState(Math.floor(Date.now() / 1000))

  // Real-time second ticker for countdowns and Airing/Aired badge checks
  useEffect(() => {
    const timer = setInterval(() => {
      setCurrentTime(Math.floor(Date.now() / 1000))
    }, 1000)
    return () => clearInterval(timer)
  }, [])

  const setRange = useCallback((val: string) => {
    if (ALLOWED_RANGES.includes(val as AllowedRange)) {
      setRangeState(val as AllowedRange)
    }
  }, [])

  const setPastRange = useCallback((val: string) => {
    if (ALLOWED_RANGES.includes(val as AllowedRange)) {
      setPastRangeState(val as AllowedRange)
    }
  }, [])

  const setPastSort = useCallback((val: string) => {
    if (ALLOWED_PAST_SORTS.includes(val as AllowedPastSort)) {
      setPastSortState(val as AllowedPastSort)
    }
  }, [])

  const fetchSchedule = useCallback(async () => {
    try {
      const sched = await apiService.getSchedule()
      setSchedule(sched)
    } catch {
      // ignore
    }
  }, [])

  const fetchPastReleases = useCallback(async () => {
    try {
      const pDays = await apiService.getPastReleases(pastRange)
      setPastDays(pDays)
    } catch {
      // ignore
    }
  }, [pastRange])

  const fetchUpcomingReleases = useCallback(async () => {
    setLoading(true)
    try {
      const uDays = await apiService.getReleases(range)
      setDays(uDays)

      // Calculate totals
      let tvCount = 0
      let movieCount = 0
      uDays.forEach((d: DayGroup) => {
        d.events.forEach((e: EventItem) => {
          if (e.type === 'tv') tvCount += (e.episode_count || 1)
          if (e.type === 'movie') movieCount++
        })
      })
      setTotalTv(tvCount)
      setTotalMovies(movieCount)
    } catch {
      setDays([])
      setTotalTv(0)
      setTotalMovies(0)
    } finally {
      setLoading(false)
    }
  }, [range])

  const refreshAll = useCallback(() => {
    fetchUpcomingReleases()
    fetchSchedule()
    if (showPast) {
      fetchPastReleases()
    }
  }, [fetchUpcomingReleases, fetchSchedule, fetchPastReleases, showPast])

  useEffect(() => {
    fetchUpcomingReleases()
  }, [fetchUpcomingReleases])

  useEffect(() => {
    fetchSchedule()
  }, [fetchSchedule])

  useEffect(() => {
    if (showPast) {
      fetchPastReleases()
    }
  }, [showPast, fetchPastReleases])

  const handleTrigger = async () => {
    if (triggering) return
    setTriggering(true)
    setTriggerStatus(null)
    try {
      await apiService.triggerSync()
      setTriggerStatus('Sync triggered successfully')
      setTimeout(() => {
        refreshAll()
      }, 1500)
    } catch (err: unknown) {
      setTriggerStatus(err instanceof Error ? err.message : 'Trigger failed')
    } finally {
      setTriggering(false)
      setTimeout(() => setTriggerStatus(null), 5000)
    }
  }

  return {
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
  }
}
