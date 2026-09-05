/**
 * Inline head script that runs before DOM render to prevent light/dark theme flashing.
 * Reads 'calendarr-theme' from localStorage or detects system prefers-color-scheme.
 */
export const antiFlashScript = `
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
