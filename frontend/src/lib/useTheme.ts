import { useState, useEffect } from 'react'

export type ThemeMode = 'system' | 'light' | 'dark'

const ALLOWED_THEMES: readonly ThemeMode[] = ['system', 'light', 'dark']

export function useTheme() {
    const [mode, setMode] = useState<ThemeMode>(() => {
        const stored = localStorage.getItem('calendarr-theme') as ThemeMode
        if (ALLOWED_THEMES.includes(stored)) return stored
        return 'system'
    })

    const [theme, setTheme] = useState<'light' | 'dark'>(() => {
        if (mode === 'dark') return 'dark'
        if (mode === 'light') return 'light'
        return window.matchMedia?.('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
    })

    useEffect(() => {
        const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')

        const applyTheme = () => {
            const active: 'light' | 'dark' = mode === 'system'
                ? (mediaQuery.matches ? 'dark' : 'light')
                : mode
            setTheme(active)
            document.documentElement.setAttribute('data-theme', active)
        }

        applyTheme()
        if (ALLOWED_THEMES.includes(mode)) {
            localStorage.setItem('calendarr-theme', mode)
        }

        const handleSystemChange = (e: MediaQueryListEvent) => {
            if (mode === 'system') {
                const active = e.matches ? 'dark' : 'light'
                setTheme(active)
                document.documentElement.setAttribute('data-theme', active)
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
