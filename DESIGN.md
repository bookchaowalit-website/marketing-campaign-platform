---
name: BookMarketing Basic
description: Minimal, system-font design system for the BookMarketing content management platform with MUI integration.
colors:
  text-primary: "#000000"
  background-start: "#f5f5f5"
  background-end: "#ffffff"
  focus-blue: "#2563eb"
  scrollbar-thumb: "#888888"
  scrollbar-hover: "#555555"
typography:
  body:
    fontFamily: "Inter, 'Noto Sans Thai', system-ui, sans-serif"
    fontWeight: 400
    lineHeight: 1.6
rounded:
  sm: "4px"
  md: "8px"
spacing:
  sm: "8px"
  md: "16px"
  lg: "24px"
components:
  button-primary:
    backgroundColor: "{colors.focus-blue}"
    textColor: "#ffffff"
    rounded: "{rounded.sm}"
    padding: "8px 16px"
  input:
    backgroundColor: "{colors.background-end}"
    textColor: "{colors.text-primary}"
    rounded: "{rounded.sm}"
    padding: "8px 12px"
---

# Design System: BookMarketing Basic

## 1. Overview

**Creative North Star: "The Clean Canvas"**

A minimal, functional design system for a content management and marketing platform. This is the most restrained system in the ecosystem — it uses Inter for Latin and Noto Sans Thai for Thai (self-hosted through `next/font`), basic RGB custom properties, and relies heavily on MUI (Material UI) for component rendering. The design intent is "get out of the way" — the content being managed is the focus, not the interface managing it.

The palette is extremely basic: black text on a light gray-to-white gradient background. The only accent color is MUI's default focus blue (#1976d2), which appears on focus rings and skip-to-main accessibility links. There is no custom button system, no card library, no elevation scale — the system defers to MUI's defaults for interactive components.

Dark mode is handled via `prefers-color-scheme` media query, swapping to white text on a dark gray gradient. No class-based dark mode toggle exists.

**Key Characteristics:**
- Inter (Latin) and Noto Sans Thai (Thai), self-hosted by `next/font`; no runtime font CDN.
- MUI-dependent. Interactive components (drawers, buttons, inputs) render through MUI's default theme.
- Minimal token set. Only 6 CSS custom properties for the entire system.
- Accessibility-first. Skip-to-main link, focus-visible outlines (ADR-0007 design floor), responsive table wrappers.
- Auto dark mode. Uses `prefers-color-scheme` — no manual toggle.

## 2. Colors

The palette is basic black-on-white with a single blue accent for focus states and the skip link.

### Primary
- **Accent Blue** (#2563eb): Focus rings, skip-to-main link, accessibility indicators. Same blue as the MUI light theme primary. White text on it is 5.2:1.

### Neutral
- **Black** (#000000): Light mode text.
- **White** (#ffffff): Light mode background end-point.
- **Light Gray** (#f5f5f5): Light mode background start, gradient origin.
- **Dark Gray** (rgb(18, 18, 18)): Dark mode background start.
- **Charcoal** (rgb(24, 24, 24)): Dark mode background end.

### Functional
- **Scrollbar Thumb** (#888888): Custom scrollbar thumb.
- **Scrollbar Hover** (#555555): Scrollbar thumb on hover.

### Named Rules
**The Two-Face Rule.** Only Inter (Latin) and Noto Sans Thai (Thai) are loaded, both self-hosted through `next/font` with the variables on `<html>`. Latin comes first in the stack. Thai text sits in a `lang="th"` element, so the page's own `lang="en"` stays true.

## 3. Typography

**Display/Body Font:** System font stack (-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif)

**Character:** The system font stack is invisible by design. It matches the user's OS typography, creating a native feel. There is no brand personality in the type — the system is a transparent window to the content being managed.

### Hierarchy
- **Headline** (600–700, 1.875rem, 1.35 line-height): Section headers. Thai headings need at least 1.3.
- **Body** (400, 1rem, 1.6 line-height): Reading text. Max line length 65–75ch.
- **Label** (400, 0.875rem): Metadata, secondary labels.

## 4. Elevation

This system has no custom elevation. It relies entirely on MUI's built-in shadow system for any elevated elements (drawers, popovers, dialogs). No custom box-shadow tokens exist.

### Named Rules
**The No-Custom-Shadow Rule.** All elevation is handled by MUI. Don't add custom box-shadow values — use MUI's elevation classes or the default theme.

## 5. Components

### Buttons
- **Shape:** MUI default radius (~4px).
- **Primary:** MUI blue background, white text. Rendered through MUI's Button component.
- **Hover:** MUI's built-in hover darkening.

### Cards
- **Corner Style:** MUI default (4px radius).
- **Background:** White (light mode).
- **Border:** None.
- **Shadow:** MUI default card elevation.
- **Internal Padding:** MUI default (16px).

### Inputs / Fields
- **Style:** White background, system border, 4px radius.
- **Focus:** 3px accent outline, 2px offset, via `:focus-visible`. MUI text inputs get the same ring (the floor rule applies to every focused control).

### Navigation
- **Style:** MUI Drawer component for sidebar navigation. Custom override removes the default right border (`MuiDrawer-paper { border-right: none }`).
- **Active State:** MUI's built-in selected state.

### Signature Components
- **Skip-to-Main Link:** Rendered once in the root layout (`.skip-link`), first Tab stop on every page, targets `#main-content`. Fixed, hidden above the viewport until focused, accent background, white text, z-index 1500.
- **Responsive Table Wrapper:** Horizontal scroll wrapper for tables on mobile (< 768px). Uses `-webkit-overflow-scrolling: touch` for smooth momentum scrolling.
- **Custom Scrollbar:** 8px wide, gray thumb on transparent track. Simple and unobtrusive.

## 6. Do's and Don'ts

### Do:
- **Do** load fonts only through `next/font`. Inter and Noto Sans Thai are the only faces.
- **Do** defer to MUI for interactive components. Buttons, drawers, inputs — use MUI's defaults.
- **Do** maintain the skip-to-main link for accessibility. It's the system's accessibility anchor.
- **Do** keep vendored design-floor files unchanged (`web/src/styles/design-floor/`); put every value in `web/src/styles/app-theme.css`.
- **Do** use `prefers-color-scheme` for dark mode. No class-based toggle exists.
- **Do** keep the token set minimal. If MUI handles it, don't duplicate it in custom CSS.

### Don't:
- **Don't** add custom fonts. The two-face approach is a deliberate performance choice.
- **Don't** add custom shadows. MUI's elevation system handles all shadow needs.
- **Don't** add accent colors beyond MUI blue. The system is deliberately neutral.
- **Don't** remove the skip-to-main link or reduce its visibility.
- **Don't** use `border-left` or `border-right` greater than 1px as colored stripes.
- **Don't** use gradient text (`background-clip: text` with gradient).
- **Don't** add uppercase tracked eyebrows above every section.

## 7. Tokens

ADR-0007 design floor. `web/src/styles/design-floor/tokens.css` (`css-vars` 1.0.0) and `tailwind/v3.config.snippet.js` (`tailwind-v3` 1.0.0) are vendored unchanged from `architecture-docs`. App values are in `web/src/styles/app-theme.css`; `web/tailwind.config.js` merges the snippet. Tailwind is v3, so the floor's v4 `layer(base)` import does not apply.

| Token | Light | Dark (`data-theme="dark"`) | Note |
|---|---|---|---|
| `--color-bg` | `#f5f5f5` | `#121212` | Page background token |
| `--color-surface` | `#ffffff` | `#1d1d1d` | Card surface |
| `--color-fg` | `#000000` | `#f2f2f2` | Primary text |
| `--color-fg-muted` | `#555555` | `#b0b0b0` | Secondary text, 6.5:1 on bg |
| `--color-accent` / `--color-focus` | `#2563eb` | `#8fb8ff` | Skip link, focus ring |
| `--color-accent-fg` | `#ffffff` | `#0b1a33` | Text on accent |
| `--color-border` | `#767676` | `#3a3a3a` | Field and divider edge |
| `--color-danger` | `#d32f2f` | `#ff8a80` | Error |
| `--color-success` | `#2e7d32` | `#6fdc9a` | Success |
| `--font-sans` | Inter, then Noto Sans Thai | same | Latin first; `next/font` variables on `<html>` |
| `--leading-thai` / `--leading-thai-heading` | 1.7 / 1.35 | same | Body and headings; MUI typography uses 1.6 / 1.35 |

**Dark mode: yes.** `ThemeProvider` follows the OS `prefers-color-scheme` until the user picks a mode, then sets `data-theme` on `<html>` so the dark token set applies alongside the MUI dark palette.

**Waivers: none.** The public auth routes and the 404 page pass all 12 automated rules at 320, 375, 390, 768 and 1440 px. Dashboard routes need a login, and the floor checker has no sign-in step, so they are covered by the manual UX review only.
