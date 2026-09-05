/**
 * Formats a target timestamp into a dynamic human-readable countdown string (e.g. "1h 33m 12s").
 * Accurately parses and evaluates target wall-clock dates within the given IANA timezone.
 */
export function formatCountdown(
    target: string | Date | number | undefined,
    timezoneOrCurrentSeconds?: string | number,
    currentUnixSeconds?: number
): string {
    if (!target || target === '-' || target === 'Not scheduled') return '-'

    let timezone: string | undefined
    let currentSec = currentUnixSeconds

    if (typeof timezoneOrCurrentSeconds === 'string') {
        timezone = timezoneOrCurrentSeconds
    } else if (typeof timezoneOrCurrentSeconds === 'number') {
        currentSec = timezoneOrCurrentSeconds
    }

    let targetMs: number | null = null

    // If target is an ISO date string and timezone is specified,
    // evaluate the wall-clock schedule time within that timezone.
    if (typeof target === 'string') {
        const match = target.match(/^(\d{4})-(\d{2})-(\d{2})T(\d{2}):(\d{2}):(\d{2})/)
        if (match) {
            const [, yStr, mStr, dStr, hStr, minStr, sStr] = match
            const y = parseInt(yStr, 10)
            const m = parseInt(mStr, 10)
            const d = parseInt(dStr, 10)
            const h = parseInt(hStr, 10)
            const min = parseInt(minStr, 10)
            const s = parseInt(sStr, 10)

            if (y <= 1) return 'Not scheduled'

            if (timezone && typeof Intl !== 'undefined') {
                try {
                    const guessUtc = Date.UTC(y, m - 1, d, h, min, s)
                    const formatter = new Intl.DateTimeFormat('en-US', {
                        timeZone: timezone,
                        year: 'numeric',
                        month: '2-digit',
                        day: '2-digit',
                        hour: '2-digit',
                        minute: '2-digit',
                        second: '2-digit',
                        hour12: false
                    })

                    const parts = formatter.formatToParts(new Date(guessUtc))
                    const getPart = (type: string) => parts.find(p => p.type === type)?.value || '0'
                    const inTzYear = parseInt(getPart('year'), 10)
                    const inTzMonth = parseInt(getPart('month'), 10)
                    const inTzDay = parseInt(getPart('day'), 10)
                    let inTzHour = parseInt(getPart('hour'), 10)
                    if (inTzHour === 24) inTzHour = 0
                    const inTzMin = parseInt(getPart('minute'), 10)
                    const inTzSec = parseInt(getPart('second'), 10)

                    const inTzUtc = Date.UTC(inTzYear, inTzMonth - 1, inTzDay, inTzHour, inTzMin, inTzSec)
                    const offset = inTzUtc - guessUtc
                    targetMs = guessUtc - offset
                } catch {
                    targetMs = Date.UTC(y, m - 1, d, h, min, s)
                }
            } else {
                targetMs = new Date(target).getTime()
            }
        }
    }

    if (targetMs === null) {
        targetMs = typeof target === 'number'
            ? (target < 1e12 ? target * 1000 : target)
            : new Date(target).getTime()
    }

    if (isNaN(targetMs) || targetMs <= 0) return '-'
    if (new Date(targetMs).getFullYear() <= 1) return 'Not scheduled'

    const currentMs = currentSec ? currentSec * 1000 : Date.now()
    const diffSec = Math.floor((targetMs - currentMs) / 1000)

    if (diffSec <= 0) return 'Due now'

    const days = Math.floor(diffSec / 86400)
    const hours = Math.floor((diffSec % 86400) / 3600)
    const minutes = Math.floor((diffSec % 3600) / 60)
    const seconds = diffSec % 60

    if (days > 0) {
        return `${days}d ${hours}h ${minutes}m ${seconds}s`
    }
    if (hours > 0) {
        return `${hours}h ${minutes}m ${seconds}s`
    }
    if (minutes > 0) {
        return `${minutes}m ${seconds}s`
    }
    return `${seconds}s`
}
