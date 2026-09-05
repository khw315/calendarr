import { THEME_STORAGE_KEY, THEME_ATTRIBUTE, THEME_MODES, ThemeMode } from '@/types/theme'

/**
 * Pure client-side initializer function executed before DOM rendering.
 * Written as a typed, linted TypeScript function rather than a raw hardcoded string.
 */
export function themeInitializer(
  storageKey: string,
  attribute: string,
  allowedModes: readonly string[]
): void {
  try {
    const stored = localStorage.getItem(storageKey)
    const mode = stored && allowedModes.includes(stored) ? stored : 'system'
    const active = mode === 'system'
      ? (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light')
      : mode
    document.documentElement.setAttribute(attribute, active)
  } catch {}
}

/**
 * Dynamically generates the self-executing anti-flash script IIFE using centralized constants.
 */
export function createAntiFlashScript(
  storageKey: string = THEME_STORAGE_KEY,
  attribute: string = THEME_ATTRIBUTE,
  modes: readonly ThemeMode[] = THEME_MODES
): string {
  return `(${themeInitializer.toString()})(${JSON.stringify(storageKey)}, ${JSON.stringify(attribute)}, ${JSON.stringify(modes)});`
}

/**
 * Default evaluated anti-flash script ready for injection into HTML <head>.
 */
export const antiFlashScript = createAntiFlashScript()
