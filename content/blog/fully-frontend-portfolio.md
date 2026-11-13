---
title: "Building a Fully Frontend Portfolio"
slug: fully-frontend-portfolio
date: 2026-11-09
summary: "How I migrated my portfolio to a static-only stack on GitHub Pages — markdown + JSON content, no backend server."
tags: ["svelte", "static", "portfolio"]
published: true
cover: "/images/blog/fullyStatic.png"
---

# Building a Fully Frontend Portfolio

This post walks through the migration from a Go + Neon + S3 backend to a fully static frontend hosted on GitHub Pages.

## Why go static?

The previous stack had a Go Fiber backend on Render (sleeping after 15 minutes), Neon Postgres, and AWS S3 for image storage. Each added cost, complexity, and cold-start latency.

## The new stack

- **Hosting:** GitHub Pages (`adapter-static`)
- **Content:** Markdown files with YAML frontmatter, committed to Git
- **Images:** Committed to `static/images/` with lazy loading
- **Contact:** `mailto:` redirect
- **Comments:** giscus (GitHub Discussions)
- **Analytics:** Client-side `localStorage` only

## Build pipeline

A small Node script (`scripts/build-content.js`) parses markdown frontmatter and emits `src/data/blog.json` and `src/data/certs.json`. Svelte components import these JSON files directly — no API calls, no backend.
