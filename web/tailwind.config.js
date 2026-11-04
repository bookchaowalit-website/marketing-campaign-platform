const floor = require('./src/styles/design-floor/tailwind/v3.config.snippet.js')

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './src/pages/**/*.{js,ts,jsx,tsx,mdx}',
    './src/components/**/*.{js,ts,jsx,tsx,mdx}',
    './src/app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  // ADR-0007 design floor: token names come from the vendored snippet (tokens.css).
  darkMode: floor.darkMode,
  theme: {
    extend: {
      ...floor.theme.extend,
      colors: {
        ...floor.theme.extend.colors,
        primary: {
          50: '#eff6ff',
          100: '#dbeafe',
          200: '#bfdbfe',
          300: '#93c5fd',
          400: '#60a5fa',
          500: '#3b82f6',
          600: '#2563eb',
          700: '#1d4ed8',
          800: '#1e40af',
          900: '#1e3a8a',
        },
      },
      // Thai line-heights: body >= 1.6, headings >= 1.3 (design floor thai-typography rule).
      fontSize: {
        xs: ['0.75rem', { lineHeight: 'var(--leading-thai)' }],
        sm: ['0.875rem', { lineHeight: 'var(--leading-thai)' }],
        base: ['1rem', { lineHeight: 'var(--leading-thai)' }],
        lg: ['1.125rem', { lineHeight: 'var(--leading-thai)' }],
        xl: ['1.25rem', { lineHeight: 'var(--leading-thai)' }],
        '2xl': ['1.5rem', { lineHeight: 'var(--leading-thai-heading)' }],
        '3xl': ['1.875rem', { lineHeight: 'var(--leading-thai-heading)' }],
        '4xl': ['2.25rem', { lineHeight: 'var(--leading-thai-heading)' }],
        '5xl': ['3rem', { lineHeight: 'var(--leading-thai-heading)' }],
      },
    },
  },
  plugins: [],
}
