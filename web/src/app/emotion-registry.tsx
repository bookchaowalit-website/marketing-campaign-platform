'use client'

// Collects the Emotion (MUI) styles rendered on the server and emits them in <head> through
// useServerInsertedHTML, so the server HTML and the first client render match (no inline
// <style> tags in the body, which React 19 reports as a hydration mismatch).
// Same pattern as MUI's AppRouterCacheProvider (@mui/material-nextjs), whose v5 build does not
// declare Next.js 16 support.
import createCache from '@emotion/cache'
import { CacheProvider } from '@emotion/react'
import { useServerInsertedHTML } from 'next/navigation'
import { useState, type ReactNode } from 'react'

export default function EmotionRegistry({ children }: { children: ReactNode }) {
  const [registry] = useState(() => {
    const cache = createCache({ key: 'mui' })
    cache.compat = true
    const insert = cache.insert
    let inserted: string[] = []
    cache.insert = (...args) => {
      const serialized = args[1]
      if (cache.inserted[serialized.name] === undefined) inserted.push(serialized.name)
      return insert(...args)
    }
    const flush = () => {
      const names = inserted
      inserted = []
      return names
    }
    return { cache, flush }
  })

  useServerInsertedHTML(() => {
    const names = registry.flush()
    if (names.length === 0) return null
    const styles = names.map((name) => registry.cache.inserted[name]).join('')
    return (
      <style
        key={registry.cache.key}
        data-emotion={`${registry.cache.key} ${names.join(' ')}`}
        dangerouslySetInnerHTML={{ __html: styles }}
      />
    )
  })

  return <CacheProvider value={registry.cache}>{children}</CacheProvider>
}
