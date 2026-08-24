import { useTheme } from '../lib/useTheme'

interface HeaderProps {
    activePage: 'dashboard' | 'settings';
}

export default function Header({ activePage }: Readonly<HeaderProps>) {
    const { mode, theme, toggleTheme } = useTheme()

    const activeStyle = { background: 'var(--color-success)', boxShadow: '0 0 0 #000', transform: 'translate(2px, 2px)' }

    const getTooltip = () => {
        if (mode === 'system') return `Theme: System (${theme === 'dark' ? 'Dark' : 'Light'})`
        if (mode === 'light') return 'Theme: Light'
        return 'Theme: Dark'
    }

    return (
        <header className="header">
            <div className="container">
                <div className="header-content">
                    <div className="logo">
                        <svg width="32" height="32" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
                            <rect x="4" y="7" width="24" height="20" rx="2" stroke="currentColor" strokeWidth="2" />
                            <line x1="4" y1="12" x2="28" y2="12" stroke="currentColor" strokeWidth="2" strokeLinecap="round" />
                            <line x1="10" y1="4" x2="10" y2="7" stroke="currentColor" strokeWidth="2" strokeLinecap="round" />
                            <line x1="22" y1="4" x2="22" y2="7" stroke="currentColor" strokeWidth="2" strokeLinecap="round" />
                        </svg>
                        <h1>Calendarr</h1>
                    </div>
                    <nav className="nav">
                        <a href="#/" className="nav-link" style={activePage === 'dashboard' ? activeStyle : undefined}>Dashboard</a>
                        <a href="#/settings" className="nav-link" style={activePage === 'settings' ? activeStyle : undefined}>Settings</a>
                        <button className="theme-toggle" onClick={toggleTheme} title={getTooltip()} aria-label={getTooltip()}>
                            {mode === 'system' ? (
                                <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                                    <rect x="2" y="3" width="20" height="14" rx="2" stroke="currentColor" strokeWidth="2" />
                                    <path d="M8 21h8M12 17v4" stroke="currentColor" strokeWidth="2" strokeLinecap="round" />
                                </svg>
                            ) : mode === 'light' ? (
                                <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                                    <circle cx="12" cy="12" r="5" stroke="currentColor" strokeWidth="2" />
                                    <path d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42" stroke="currentColor" strokeWidth="2" strokeLinecap="round" />
                                </svg>
                            ) : (
                                <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                                    <path d="M21 12.79A9 9 0 1111.21 3a7 7 0 009.79 9.79z" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
                                </svg>
                            )}
                            <span className="theme-label" style={{ marginLeft: '6px', fontSize: '0.75rem', fontWeight: 'bold', textTransform: 'uppercase' }}>
                                {mode === 'system' ? 'SYSTEM' : mode}
                            </span>
                        </button>
                    </nav>
                </div>
            </div>
        </header>
    )
}
