import { Providers } from './providers'
import './globals.css'
import type { Metadata } from 'next'
import { Inter, Noto_Sans_Thai } from 'next/font/google'
import ErrorBoundary from '@/components/ErrorBoundary'
import EmotionRegistry from './emotion-registry'

// next/font self-hosts both faces at build time (no runtime CDN). The variables go on <html>, not
// <body>: app-theme.css reads them from :root.
const inter = Inter({ subsets: ['latin'], variable: '--font-inter' })
const notoSansThai = Noto_Sans_Thai({
  subsets: ['thai', 'latin'],
  weight: ['400', '500', '600', '700'],
  variable: '--font-noto-sans-thai',
})

export const metadata: Metadata = {
  title: 'MarketPro - Marketing Dashboard',
  description: 'Marketing SaaS Platform',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en" className={`${inter.variable} ${notoSansThai.variable}`}>
      <body>
        <a className="skip-link" href="#main-content">
          Skip to main content
        </a>
        <EmotionRegistry>
          <ErrorBoundary>
            <Providers>
              {children}
            </Providers>
          </ErrorBoundary>
        </EmotionRegistry>
      </body>
    </html>
  )
}
