import type { Metadata } from 'next'
import '../index.css'

export const metadata: Metadata = {
  title: 'Calendarr',
  description: 'Calendar feeds from Sonarr/Radarr to Discord and Slack',
}

const antiFlashScript = `
(function() {
  try {
    var stored = localStorage.getItem('calendarr-theme');
    var mode = (stored === 'light' || stored === 'dark' || stored === 'system') ? stored : 'system';
    var active = mode === 'system' 
      ? (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light')
      : mode;
    document.documentElement.setAttribute('data-theme', active);
  } catch (e) {}
})();
`

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
