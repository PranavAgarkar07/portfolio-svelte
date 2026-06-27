# UX/UI Re-Audit — Pranav Agarkar Portfolio

**Site:** https://pranavagarkar07.github.io/portfolio-svelte/
**Stack:** SvelteKit + Vanilla CSS, deployed on GitHub Pages
**Audit Date:** 2026-06-26
**Type:** Re-audit (comparing against prior audit of same date)
**Aesthetic Intent:** Neo-Brutalist (Aerospace / Industrial Luxury variant)

---

## Executive Summary

### Overall UX Score: **87 / 100** (▲ +13 from 74)

Nearly every high-priority issue from the first audit has been addressed. The portfolio has materially improved — especially in visual balance, form robustness, accessibility, and mobile interaction design. The remaining issues are refinements, not blockers.

### ✅ Top 3 Improvements Since Last Audit

1. **Hero two-column layout at 1100px+** — The avatar + stats card fills the right viewport. The dead-space issue at 1440px is completely resolved.
2. **Contact form production hardening** — Autocomplete attributes, submit disabled during wakeup, server-down banner with email fallback, field-level validation. This form is now best-in-class for a student portfolio.
3. **All a11y quick wins landed** — Like button `aria-label`, cert verify link touch target, mobile carousel swipe, skip-link deduplication, global `prefers-reduced-motion` support.

### ❌ Top 3 Remaining Issues

1. **About section lives inside Hero.svelte** — Architectural smell. The `#about` scroll-spy target shares a component with `#hero`, making nav highlighting unreliable and creating a confusing component boundary.
2. **Contact form progress counter still subtle** — The `0/3` bar is a thin 2px line. Progress is not visually noisy enough to leverage Zeigarnik Effect effectively.
3. **Section order: About blocks Projects** — IA still routes recruiters through a bio before showing proof (projects). Changing order would improve time-to-value.

---

## Previously Flagged — Status Check

| # | Issue | Priority | Status | Current State |
|---|-------|----------|--------|---------------|
| 1 | Hero right-side whitespace at 1440px | HIGH | ✅ FIXED | Avatar + stats card fills right column |
| 2 | Marquee clips first word on load | HIGH | ✅ FIXED | Duplicated array + `animation-delay: 2.5s` |
| 3 | Serial Position: BeamSync not first | HIGH | ✅ FIXED | BeamSync is now index 0 |
| 4 | Contact form: submit not disabled during wakeup | HIGH | ✅ FIXED | `disabled={submitting \|\| waking \|\| success}` |
| 5 | Missing autocomplete on form fields | HIGH | ✅ FIXED | `autocomplete="name"`, `autocomplete="email"` |
| 6 | Like button missing aria-label | HIGH | ✅ FIXED | `aria-label={liked ? "Unlike" : "Like this project"}` |
| 7 | No mobile swipe on carousel | HIGH | ✅ FIXED | Touch start/end handlers with 50px threshold |
| 8 | Circe icon on "Contact Me" CTA | MEDIUM | ✅ FIXED | Arrow-right icon with line + polyline |
| 9 | "In Orbit" label opaque to recruiters | MEDIUM | ✅ FIXED | Renamed to "Exploring" in data.ts |
| 10 | GitHub stars buried in description | MEDIUM | ✅ FIXED | Visual `★ 32` badge in card header |
| 11 | Tech tags overflow on KisanRakshak | MEDIUM | ✅ FIXED | `slice(0,5)` + `+N overflow` pill |
| 12 | Global arrow key conflict across cards | MEDIUM | ✅ FIXED | Keydown moved to `onkeydown` on `<article>` |
| 13 | "Verify Credential" touch target < 44px | MEDIUM (a11y) | ✅ FIXED | `min-height: 44px` on cert link |
| 14 | Form progress counter too subtle | MEDIUM | ⚠️ IMPROVED | Now a progress bar, but still low-visibility |
| 15 | Duplicate skip-link in layout | LOW | ✅ FIXED | No duplicate skip-link visible in current code |
| 16 | Footer missing strong CTA | MEDIUM | ✅ FIXED | "Get In Touch" CTA + Back-to-top button added |
| 17 | Hero→About→Projects IA order | MEDIUM | ❌ NOT CHANGED | Still Hero→About→Projects in page flow |

---

## Law-by-Law Verdict Table (Updated)

| # | Law | Status | Severity | Component | Finding |
|---|-----|--------|----------|-----------|---------|
| 1 | Aesthetic-Usability Effect | ✅ | — | Global | Intentional Neo-Brutalism, light mode exceptional, border/shadow weights consistent |
| 2 | Choice Overload | ✅ | — | Header | 5 nav links + Resume + theme toggle — now well-balanced since hero layout is filled |
| 3 | Chunking | ✅ | — | Skills | 3 card columns with clear category headers |
| 4 | Cognitive Bias | ✅ | — | Hero, Cards | "32 stars", "1,247 req/min", "2+ Years" — strong specificity anchors |
| 5 | Cognitive Load | ⚠️ | L | Hero | Now 6 primary elements (was 9). Avatar card consolidates stats. **Improved** |
| 6 | Doherty Threshold | ✅ | — | ContactForm | Server-wake spinner visible immediately; submit disabled during wait |
| 7 | Fitts's Law | ✅ | — | Multiple | CTAs ≥44px, verify link 44px min-height, carousel buttons 44px. **All fixed** |
| 8 | Flow | ✅ | — | layout | Lenis smooth scroll + GSAP reveals with no pop-in interruptions |
| 9 | Goal-Gradient | ⚠️ | M | ContactForm | \(0/3\) progress bar exists but visual weight still low |
| 10 | Hick's Law | ✅ | — | ProjectCard | Tags capped at 5 + overflow pill. **Fixed** |
| 11 | Jakob's Law | ✅ | — | Global | Nav top, logo→top, resume→new tab, socials in footer |
| 12 | Common Region | ✅ | — | All | Card containers, skill groups, section headers all have clear boundaries |
| 13 | Proximity | ⚠️ | L | Skills | Project refs (`TaskVault Lite · KisanRakshak`) still clip on narrow cards |
| 14 | Prägnanz | ✅ | — | Skills | "Exploring" replaces "In Orbit". **Fixed** |
| 15 | Similarity | ✅ | — | Global | All cards/tags styled identically per type |
| 16 | Uniform Connectedness | ✅ | — | layout | Sections flow visually via shared grid and accent borders |
| 17 | Mental Model | ✅ | — | Skills | "Core Stack", "Proficient", "Exploring" — recruiter-compatible. **Fixed** |
| 18 | Miller's Law | ✅ | — | ProjectCard | Card items capped (5 tags, 2 CTAs). Nav = 7 items at limit but manageable |
| 19 | Occam's Razor | ⚠️ | L | Hero | Social links in 4 places (hero→about→contact→footer) — intentional but review |
| 20 | Paradox of Active User | ✅ | — | Hero | Name + role + "View Work" CTA immediately obvious |
| 21 | Pareto Principle | ✅ | — | Hero | Avatar card fills right 30% at 1100px+. **Fixed** |
| 22 | Parkinson's Law | ✅ | — | ContactForm | Autocomplete + draft save + sensible defaults. **Fixed** |
| 23 | Peak-End Rule | ✅ | — | Footer | "Get In Touch" CTA + Back-to-top button. **Fixed** |
| 24 | Postel's Law | ✅ | — | ContactForm | Email regex, default topic, autofill, draft persistence. **Fixed** |
| 25 | Selective Attention | ⚠️ | L | Hero | Marquee is now well-delayed (2.5s) and visually polished |
| 26 | Serial Position Effect | ✅ | — | data.ts | BeamSync (32★, featured) is now PRJ 01. **Fixed** |
| 27 | Tesler's Law | ✅ | — | layout | FOUC-free, dark/light handled by designer, simple for user |
| 28 | Von Restorff Effect | ✅ | — | ProjectCard | "★ 32" badge is now visually isolated. **Fixed** |
| 29 | Zeigarnik Effect | ⚠️ | L | ContactForm | \(0/3\) counter exists but barely visible — needs bold treatment |
| 30 | WCAG 2.2 Accessibility | ✅ | — | Global | skip-link ✅, focus-visible ✅, reduced-motion ✅, aria-labels ✅, semantic HTML ✅. **All fixed** |
| 31 | F/Z-Pattern | ✅ | — | Hero | Left-aligned name + CTAs — correct F-pattern |
| 32 | Visual Hierarchy | ✅ | — | Global | H1→H2→H3 scale consistent, weight + color unambiguous |
| 33 | Gestalt Principles | ✅ | — | Hero | Avatar contrast good, marquee has mask + delayed start. **Improved** |
| 34 | Mobile-First Responsive | ✅ | — | ProjectCard | Swipe + buttons + dots. Hamburger menu. Hero wraps cleanly. **All fixed** |
| 35 | Information Architecture | ❌ | M | page | Hero→About→Projects. Should be Hero→Projects. **NOT FIXED** |
| 36 | Typography UX | ✅ | — | styles.css | 14px body, lh 1.6, ~56 chars/line — within ideal |
| 37 | Performance UX | ✅ | — | Multiple | .webp, lazy loading, explicit dimensions, preconnects. Strong |
| 38 | Interaction Feedback | ✅ | — | ContactForm | Hover states, success animation, disabled states, shake on error. **Improved** |
| 39 | Trust & Credibility | ✅ | — | Multiple | Real photo, working links, specific metrics, live DevLog status |
| 40 | Dark/Light Mode | ✅ | — | app.html | localStorage, OS preference, FOUC-free, fully styled both modes |

---

## New Issues Identified

### 1. About Section Lives Inside Hero.svelte (ARCHITECTURE — MEDIUM)

`Hero.svelte` renders both `<section id="hero">` and `<section id="about">` as siblings. This means:
- The `#about` scroll-margin (`scroll-margin-top: 80px`) defined in `styles.css:585` may conflict with the hero's padding
- Scroll-spy for "about" fires while the hero/avatar/content of the previous section is still visible
- The component name misleads future developers

**Fix:** Extract `#about` into its own `+page.svelte` block or a dedicated `About.svelte` component. Currently rendered at `Hero.svelte:258-344`.

### 2. Duplicate Avatar HTTP Request (PERFORMANCE — LOW)

The same `profile.avatar` is loaded twice:
- `Hero.svelte:239` — `loading="eager"` (in hero right card)
- `Hero.svelte:263-270` — `loading="lazy"` (in about frame)

Browser may cache, but it's an unnecessary second `<img>` element for the same 280×280 asset.

**Fix:** Use a shared `<img>` with CSS `object-position` tricks, or lazy-load the about version only.

### 3. Form Progress Bar Still Subtle (ZEIGARNIK — MEDIUM)

The \(0/3\) progress at `ContactForm.svelte:151-161` uses a 2px-tall bar and tiny gray text. Despite the previous fix adding the bar, it still lacks visual weight. The checkmark appears only at 3/3 — no partial celebration.

**Fix:** Increase bar to 3px, use accent orange, add per-field glow effect when each field is completed.

### 4. Body Font Size = 14px (TYPEOGRAPHY — LOW)

`font-size: 14px` on body (`styles.css:54`) is below the recommended 16px minimum for body text per WCAG and UX convention. At 14px, the line length (~56 chars) is still within range, but readability suffers on high-DPI mobile screens.

**Fix:** `font-size: 15px` or `16px` on body, adjust spacing accordingly.

### 5. Marquee Tech List Hardcoded in Template (MAINTAINABILITY — LOW)

The tech ticker array is duplicated inline at `Hero.svelte:226`:
```js
['Python','Django','React','Go','Svelte','TypeScript','PostgreSQL','Docker','Linux',
 'Python','Django','React','Go','Svelte','TypeScript','PostgreSQL','Docker','Linux']
```
This is duplicated from the skills data in `data.ts`. If skills change, the marquee must be manually kept in sync.

**Fix:** Derive the marquee list from `portfolioData.skills` or define it once in data.

### 6. Light Mode: Hero Status Icon Colors (VISUAL — LOW)

In light mode, the hero status icon (`hero-status-icon.available`) and value turn green (`#008f5e`) in `Hero.svelte:911-917`. Green status dot is fine, but the green `[+]` crosshair icon in the status bar looks disconnected from the orange accent system.

**Fix:** Keep the status dot/value green but tint the SVG with accent orange. Currently the entire icon turns green.

---

## Updated Quick Wins (< 1 hour each)

| # | Fix | Est. Time | Impact |
|---|-----|-----------|--------|
| 1 | Derive marquee tech list from data source | 10 min | 🟢 LOW (maintainability) |
| 2 | Bump body font-size to 15px | 5 min | 🟡 MEDIUM (readability) |
| 3 | Increase form progress bar to 3px + accent glow per field | 20 min | 🟡 MEDIUM (Zeigarnik) |
| 4 | Fix light-mode hero status icon to use accent orange | 5 min | 🟢 LOW (visual consistency) |

---

## Architecture Suggestions (Still Open From Prior Audit)

### A. Section Order Restructure (Medium Priority)

**Current:** Hero → About → Projects → Skills → Certifications → Contact
**Recommended:** Hero → Projects → Skills → About → Certifications → Contact

Recruiters scan for proof (projects) before story (bio). Moving Projects to immediately follow Hero maximizes the 20% of content that delivers 80% of value (Pareto Principle). This requires:
1. Extracting `#about` from `Hero.svelte` into its own section block
2. Reordering sections in `+page.svelte`
3. Updating nav order in `Header.svelte:13-18`

### B. About Extraction from Hero.svelte (Medium Priority)

Before reordering sections, the `#about` content must be extracted from `Hero.svelte`. Currently both `<section id="hero">` and `<section id="about">` are rendered by the same component. This causes:
- Scroll-spy to trigger "about" nav highlight while hero content is still visible
- Misleading component architecture
- Blocks section reordering

**Suggested approach:** Create `About.svelte` component, move lines 258-344 from `Hero.svelte` into it, render in `+page.svelte`.

### C. Server Warmup: Passive vs Proactive (Low Priority)

The current wakeup approach (triggered on submit click) is now well-handled with clear loading states. But consider: calling `wakeServer()` when `#contact` scrolls into view (IntersectionObserver in `+page.svelte`) would pre-warm the server 30-60s before submission, eliminating the wait entirely.

---

## Summary

| Metric | Prior Audit | This Audit | Δ |
|--------|-------------|------------|---|
| UX Score | 74 / 100 | **87 / 100** | **+13** |
| HIGH issues | 7 | 0 | **-7** |
| MEDIUM issues | 14 | 4 | **-10** |
| LOW issues | 5 | 6 | +1 (new edge cases found) |
| Blockers | 3 | 0 | **Cleared** |

The portfolio has moved from "good with notable gaps" to "polished with refinement opportunities." Every high-severity issue from the first audit has been resolved. The remaining work is architectural cleanup (About extraction, section reorder) and visual polish (form progress weight, body font size, light-mode icon consistency).

---

*Re-audit generated: 2026-06-26 | Methodology: Full code re-read + comparison against prior audit report + Laws of UX + WCAG 2.2 + supplemental UX principles*
