import { useState, useEffect } from 'react'
import {
    ThemeMode,
    THEME_STORAGE_KEY,
    THEME_ATTRIBUTE,
    THEME_MODES
} from '@/types/theme'

export type { ThemeMode }
export { THEME_STORAGE_KEY, THEME_ATTRIBUTE, THEME_MODES }

export function useTheme() {
    const [mode, setMode] = useState<ThemeMode>(() => {
        if (typeof window !== 'undefined') {
            const stored = localStorage.getItem(THEME_STORAGE_KEY) as ThemeMode
            if (THEME_MODES.includes(stored)) return stored
        }
        return 'system'
    })

    const [theme, setTheme] = useState<'light' | 'dark'>(() => {
        if (mode === 'dark') return 'dark'
        if (mode === 'light') return 'light'
        if (typeof window !== 'undefined') {
            return window.matchMedia?.('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
        }
        return 'light'
    })

    useEffect(() => {
        const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')

        const applyTheme = () => {
            const active: 'light' | 'dark' = mode === 'system'
                ? (mediaQuery.matches ? 'dark' : 'light')
                : mode
            setTheme(active)
            document.documentElement.setAttribute(THEME_ATTRIBUTE, active)
        }

        applyTheme()
        if (THEME_MODES.includes(mode)) {
            localStorage.setItem(THEME_STORAGE_KEY, mode)
        }

        const handleSystemChange = (e: MediaQueryListEvent) => {
            if (mode === 'system') {
                const active = e.matches ? 'dark' : 'light'
                setTheme(active)
                document.documentElement.setAttribute(THEME_ATTRIBUTE, active)
            }
        }

        mediaQuery.addEventListener('change', handleSystemChange)
        return () => mediaQuery.removeEventListener('change', handleSystemChange)
    }, [mode])

    const toggleTheme = () => {
        setMode(prev => {
            if (prev === 'system') return 'light'
            if (prev === 'light') return 'dark'
            return 'system'
        })
    }

    return { mode, theme, toggleTheme, setMode }
}
