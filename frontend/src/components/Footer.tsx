import { useEffect, useState } from 'react'

const API_BASE = ''

export default function Footer() {
    const [version, setVersion] = useState<string>('v2.1.0')

    useEffect(() => {
        const fetchVersion = async () => {
            try {
                const res = await fetch(`${API_BASE}/api/version`)
                if (res.ok) {
                    const data = await res.json()
                    if (data.version) {
                        setVersion(data.version.startsWith('v') ? data.version : `v${data.version}`)
                    }
                }
            } catch {
                // fallback to default
            }
        }
        fetchVersion()
    }, [])

    return (
        <footer className="footer">
            <div className="container footer-container">
                <p>
                    <a href="https://github.com/khw315/calendarr" target="_blank" rel="noreferrer">
                        GitHub
                    </a>
                </p>
                <div className="footer-badges">
                    <span className="docker-badge" title="Docker Container Version">
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg" style={{ marginRight: '6px', display: 'inline-block', verticalAlign: 'middle' }}>
                            <path d="M4 11h3v3H4v-3zm4 0h3v3H8v-3zm4 0h3v3h-3v-3zm4 0h3v3h-3v-3zm-8-4h3v3H8V7zm4 0h3v3h-3V7zm4 0h3v3h-3V7zm4 4h3v3h-3v-3zM2 17c0 2.2 3.6 4 10 4s10-1.8 10-4c0-2-3-3.6-8.5-3.9l-.5 1.9C18 15.3 20 16.3 20 17c0 .8-3.2 2-8 2s-8-1.2-8-2c0-.7 2-1.7 7-2l-.5-1.9C5 13.4 2 15 2 17z" fill="currentColor"/>
                        </svg>
                        docker: {version}
                    </span>
                </div>
            </div>
        </footer>
    )
}
