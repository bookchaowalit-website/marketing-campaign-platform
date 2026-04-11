// Simple event-based search store (no external deps)
// Useful for non-React contexts or programmatic control

type SearchListener = (open: boolean) => void
let listeners: SearchListener[] = []

export function openSearch() {
  listeners.forEach((fn) => fn(true))
}

export function closeSearch() {
  listeners.forEach((fn) => fn(false))
}

export function toggleSearch(open?: boolean) {
  listeners.forEach((fn) => fn(open ?? true))
}

export function subscribeSearch(listener: SearchListener) {
  listeners.push(listener)
  return () => {
    listeners = listeners.filter((l) => l !== listener)
  }
}
