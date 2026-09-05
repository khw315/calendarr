import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
    return twMerge(clsx(inputs))
}

export function formatCountdown(target: string | Date | number | undefined, currentUnixSeconds?: number): string {
    if (!target || target === '-' || target === 'Not scheduled') return '-'

    const targetMs = typeof target === 'number'
        ? (target < 1e12 ? target * 1000 : target)
        : new Date(target).getTime()

    if (isNaN(targetMs) || targetMs <= 0) return '-'
    if (new Date(targetMs).getFullYear() <= 1) return 'Not scheduled'

    const currentMs = currentUnixSeconds ? currentUnixSeconds * 1000 : Date.now()
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
