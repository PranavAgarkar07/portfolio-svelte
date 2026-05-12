# Portfolio — Pranav Agarkar

> Fullstack portfolio built with Svelte 5, SvelteKit, and TypeScript. Features a custom aerospace-brutalist design system, live GitHub activity feed, encrypted contact form with admin panel, and interactive project showcase.

## Tech Stack

- **Framework:** Svelte 5 + SvelteKit (SPA mode for GitHub Pages)
- **Language:** TypeScript
- **Styling:** CSS custom properties with dark/light mode
- **Animations:** GSAP, ScrollTrigger, Lenis
- **Backend:** Go (Sentinel API — Render)
- **Database:** PostgreSQL (TaskVault)
- **Deployment:** Frontend → GitHub Pages, Backend → Render

## Features

- **Aerospace Brutalist Design** — Chamfered cards, noise grain overlay, CAD grid backgrounds
- **Dark/Light Mode** — CSS variable-driven theming with smooth toggle
- **Live DevLog** — Real-time GitHub activity status from Go backend
- **Interactive Terminal** — CLI-style navigation widget
- **Project Showcase** — Image carousels with Lightbox3 integration
- **Contact Form** — Topic-based messaging with admin control panel
- **Component Library** — 15+ reusable UI components (Button, Card, Input, Select, DatePicker, etc.)

## Getting Started

```bash
# Install dependencies
npm install

# Start dev server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

## Project Structure

```
src/
├── lib/
│   ├── components/     # Svelte components (Hero, ProjectCard, ContactForm, etc.)
│   │   └── ui/         # Reusable UI components (Button, Card, Input, etc.)
│   ├── stores/         # Svelte stores (theme, scroll)
│   ├── assets/         # Static assets (images)
│   ├── data.ts         # Portfolio content data
│   └── types.ts        # TypeScript type definitions
├── routes/
│   ├── +page.svelte    # Main landing page
│   ├── +layout.svelte  # Root layout (animations, fonts, meta)
│   ├── showcase/       # UI component library showcase
│   └── admin/contact/  # Contact message admin panel
└── styles.css          # Global styles (dark mode)
└── light-mode.css      # Light mode overrides
```

## Deployment

See [DEPLOY.md](./DEPLOY.md) for the two-tier deployment setup.

## License

MIT © Pranav Agarkar
