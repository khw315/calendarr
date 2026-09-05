import React from 'react'
import { antiFlashScript } from '@/utils'

/**
 * Head component that renders the inline theme anti-flash script.
 * Prevents light/dark theme flickering on initial page load.
 */
export function AntiFlashScript() {
  return (
    <script
      id="anti-flash-theme"
      dangerouslySetInnerHTML={{ __html: antiFlashScript }}
    />
  )
}
