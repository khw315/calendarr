import type { Metadata } from 'next'
import '../index.css'
import { antiFlashScript } from '@/lib/antiFlashScript'

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
        <script dangerouslySetInnerHTML={{ __html: antiFlashScript }} />
      </head>
      <body>
        <div id="root">
          {children}
        </div>
      </body>
    </html>
  )
}
