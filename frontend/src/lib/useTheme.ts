import { useState, useEffect } from 'react'

type Theme = 'light' | 'dark'

export function useTheme() {
    const [theme, setTheme] = useState<Theme>(() => {
        const stored = localStorage.getItem('calendarr-theme')
        if (stored === 'light' || stored === 'dark') return stored
        // Respect system preference
        if (window.matchMedia?.('(prefers-color-scheme: dark)').matches) return 'dark'
        return 'light'
    })

    useEffect(() => {
        document.documentElement.setAttribute('data-theme', theme)
        localStorage.setItem('calendarr-theme', theme)
    }, [theme])

    const toggleTheme = () => setTheme(prev => prev === 'light' ? 'dark' : 'light')

    return { theme, toggleTheme }
}
