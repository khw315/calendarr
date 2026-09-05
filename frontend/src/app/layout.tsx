import type { Metadata } from 'next'
import '../index.css'
import { AntiFlashScript } from '@/components/common/AntiFlashScript'

export const metadata: Metadata = {
  title: 'Calendarr',
  description: 'Calendar feeds from Sonarr/Radarr to Discord and Slack',
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <head>
        <AntiFlashScript />
      </head>
      <body>
        <div id="root">
          {children}
        </div>
      </body>
    </html>
  )
}
