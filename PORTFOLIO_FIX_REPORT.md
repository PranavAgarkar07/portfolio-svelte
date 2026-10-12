# Portfolio Perf + Code-Health Fix Report
> Perspective: SDE, Google SF (San Francisco) — Core Web Vitals, reliability, maintainability bar.
> Repo: `portfolio-svelte` — SvelteKit 2 + Svelte 5 static frontend (GitHub Pages) + Go Fiber backend "Sentinel API" (Render/Lambda + Postgres + S3).
> Prior UX audit: 87/100, 0 HIGH open. This report covers what it did not: JS/CSS bytes, INP/p95, cold-start, DB/index, resilience, security.
> Status: implementation spec for parallel LLM agent swarm. No visual redesign unless CLS fix requires it.

## 0. Performance budgets (enforce after every wave)

- Frontend: LCP < 2.5s (Moto G4 / 4G), CLS < 0.1, INP < 200ms, homepage JS < 250KB gzip, CSS < 30KB gzip, below-fold images lazy + sized.
- Backend: p95 < 400ms warm (Neon cross-region), cold-start < 5s Render free, Lambda init < 2s, DB pool saturation 0, GitHub/OpenRouter external calls deduped (singleflight).
- Gates per wave: `npm run check && npm run build`, `go vet ./... && go build ./...`, Lighthouse mobile+desktop before/after, `git diff --stat` review.

## 1. P0 — fix first (outage / LCP / INP / cold-start)

### F-P0-1: Route-split heavy JS out of root layout
- Files: `src/routes/+layout.svelte:9-13,106-109,113-131`, `src/routes/+page.svelte:3-4`, `src/lib/components/Skills.svelte:3-4`, `vite.config.ts:6-14`, `package.json`
- Findings:
  - `gsap`, `ScrollTrigger`, `lenis`, `lightbox3` statically imported in root layout → shipped to `/admin`, `/blog` even when unused (mobile/touch, reduced-motion).
  - Triple `from "gsap/dist/ScrollTrigger"` instead of `gsap/ScrollTrigger` risks duplicate module instance, defeats `vendor-gsap` chunking.
  - `gsap.registerPlugin` 4x across 3 components + `gsap.defaults` 2x in same `onMount`.
- Fix:
  - Dynamic `import()` GSAP/Lenis/Lightbox only when needed (non-touch, no reduced-motion, lightbox images present).
  - Single `gsap/ScrollTrigger` import path. Single `registerPlugin` in a shared `lib/motion.ts`.
  - Extend `manualChunks`: `lightbox3`, `@humanspeak/svelte-motion` out of main chunk.
- Accept: `/admin` bundle has no gsap; homepage JS drops; no behavior change on desktop.

### F-P0-2: Stop hiding hero/projects without JS
- Files: `src/routes/+page.svelte:123-134,217-235`, `src/lib/components/Hero.svelte:27,86-148,286,388,425,443`
- Findings: `gsap.set(...opacity:0)` + `heroReady` rAF gate → invisible content if JS fails; `forceVisibility` 500ms + `clipPath` reset causes style recalc flash.
- Fix: default-visible CSS; use GSAP `from()` reveals only when motion allowed. Add `<noscript>` fallback. Remove `forceVisibility` hack or keep as 1s safety only.
- Accept: disable JS → hero + projects visible; no 500ms flash with JS on.

### F-P0-3: Image LCP / network fan-out
- Files: `Hero.svelte:153`, `ProjectCard.svelte:238`, `About.svelte:107-115,150-156`, `src/lib/data.ts:140-270`, `static/images` (6.6MB png+webp dupes), `svelte.config.js:17`
- Findings:
  - 4 eager images (every card `i===0` eager per card) → 4 LCP contenders. Only first viewport image should be eager.
  - Avatar downloaded twice (Hero eager 280px + About lazy 400px), no `srcset`, no `fetchpriority`.
  - Badges ship png 440-530KB + webp 72-78KB each (~6MB total), png unused. Kisan png 883KB vs webp 161KB dupes still deployed.
  - No `width/height/sizes/srcset`; relative `images/...` breaks under `base:/portfolio-svelte`.
- Fix:
  - 1 eager LCP image only (`fetchpriority="high"`, `srcset`/`sizes`); rest `loading="lazy" decoding="async"`.
  - Share avatar (single URL, right size). Delete or stop deploying unused png originals.
  - Use `{base}` for image srcs. Add width/height to badge imgs. Consider `vite-imagetools`/AVIF + blur placeholder.
- Accept: Lighthouse LCP = 1 request; mobile image bytes down ~1MB; no 404 under `/portfolio-svelte`.

### F-P0-4: Homepage backend thundering herd
- Files: `+page.svelte:72-93`, `analytics.ts:189-199`, `DevLog.svelte:11-49`, `ContactForm.svelte:70-100`, `ReviewSection.svelte:61-106`, `ProjectCard.svelte:198,295`
- Findings: cold start fires certs+badges+marquee+likes+status+healthz+4x reviews ≈ 7-10 parallel hits to sleeping Render backend. `wakeServer` 3x(5s+3s) ≈ 24s blocks submit. DevLog sequential 8s→30s = 38s shimmer. No priority, no AbortController on most, no dedup.
- Fix:
  - Prioritize critical (certs) with timeouts/AbortController + `Promise.allSettled`.
  - Defer marquee/reviews below fold (IntersectionObserver). Single shared warmup between DevLog + ContactForm.
  - Dedupe ReviewSection → batched `/reviews?projects=` or lazy per-viewport.
- Accept: cold homepage ≤ 3 critical requests; submit not blocked > 5s; reviews load on scroll.

### B-P0-1: Boot chain cold-start (S3 IMDS + migrations + analytics DELETEs)
- Files: `backend/main.go:54-58,67,97-100,136-278`, `backend/internal/analytics/store.go:164-178`, `backend/cmd/backfill-thumbs/main.go:67`, `render.yaml`
- Findings: sync AWS IMDS (no timeout) + ~20 DDL + 4 full-table `UPDATE...LIKE` fixups + 2 unbounded `DELETE`s before `Listen` → Render free 30-50s cold → probe timeout/502. `RUN_MIGRATIONS` unset by default → schema drift; set → readiness delayed.
- Fix: lazy S3 client with 2s timeout; move migrations/fixups to pre-deploy job or versioned + ctx timeout + txn; analytics cleanup async; add `healthCheckPath` in `render.yaml`.
- Accept: `/healthz` 200 < 2s warm, < 8s cold.

### B-P0-2: No ctx timeouts + pool exhaustion
- Files: `backend/internal/repository/*.go`, `backend/internal/analytics/store.go`, `backend/main.go:49-52`, `handler/upload.go:107,118`, `handler/blog.go:37`, `handler/auth.go:93,99`, `handler/devlog.go:80-99`, `go.mod`
- Findings: zero `QueryContext/ExecContext`. Only ping/health/OpenRouter have timeouts. Lambda x10 conns each → Neon exhaustion, indefinite wait. Dual `pgx + pq` drivers.
- Fix: `context.WithTimeout` (3-5s reads, 10s writes) threaded from Fiber `c.Context()`; pool `MaxOpenConns` per env (Lambda 2-3 + `?pgbouncer=true`); S3/OAuth/GitHub ctx + retry; keep one driver.
- Accept: `go vet` clean; slow query cannot hold pool forever.

### B-P0-3: Limiter + CORS + body limits
- Files: `backend/main.go:75-92`, `backend/internal/handler/* BodyParser`, `backend/internal/config/config.go`
- Findings: global 120/min keyed by `c.IP()` with no `TrustedProxy` → shared bucket behind proxy; counts OPTIONS/health/static; no `BodyLimit`; unbounded `content_md/comment/message/order[]`.
- Fix: `EnableTrustedProxy`/`ProxyHeader`, `Skip` health/static/OPTIONS, per-route limits (auth/contact stricter), `BodyLimit` + per-field maxlength/count caps.
- Accept: health not rate-limited; oversize → 413, not 500.

### B-P0-4: Sync S3 + CPU thumb on request path
- Files: `backend/internal/handler/upload.go:83-124`, `backend/internal/handler/blog.go:20-42,286,305`
- Findings: `image.Decode + Bilinear 640px + jpeg.Encode` then sequential `PutObject` x2 on request goroutine with `Background` ctx; blog deletes loop S3 sync before response.
- Fix: timeouts, background/bulk deletes, async thumb, or presigned-URL direct-to-S3. Thumb failure must not fail upload silently (currently ignored at `:118`).
- Accept: upload p95 halved; shutdown < 10s.

## 2. P1 — fix next (p95, CLS, memory, DB)

### P1-1: Scroll/anim leaks + per-frame fan-out
- Files: `+layout.svelte:34-193,183-194`, `+page.svelte:67-244,99-103,108-118,138,159,167-177,237-242`, `stores/scroll.ts:3`, `Header.svelte:86-143`, `About.svelte:76-99`, `Hero.svelte:19-71`, `ProjectCard.svelte:63-77,113-121`
- Fix: kill ScrollTriggers/gsap.ticker/observers/resize/touchstart/progressBar on destroy; batch duplicate `#about`/`#skills` triggers; throttle `scrollY.set`, `frameEl.style.transform`, mousemove parallax, count-up via rAF; remove static `--x/--y` mouseenter listeners or make them cursor-aware + delegated.

### P1-2: Competing hash-scroll + progress systems
- Files: `+layout.svelte:24-32,133-145,147-166,197-210`, `+page.svelte:108-118,133-145`, `Header.svelte:146-148`, `blog/[slug]/+page.svelte:43-59,119,228-236`
- Fix: single hash-scroll util (uses `base`, respects reduced-motion, works with Lenis); single store-driven progress bar component; bind anchors via delegation for dynamic CTAs.

### P1-3: N+1 / sequential DB + missing indexes + no pagination
- Files: `handler/reviews.go:24,42,47`, `handler/likes.go:44,63,97,101`, `handler/handler.go:136-143`, `handler/comments.go:60,100,105,126,139`, `handler/blog.go:68-89,110-153,229,288,318,339`, `handler/auth.go:93-139`, `analytics/store.go:49-160`, `repository/*.go`, `main.go:202-205`
- Fix: JOIN/`RETURNING`/CTE to collapse 2-3 trips → 1; `singleflight` for devlog/status; dashboard cache 60s; add indexes `posts(published,published_at)`, `comments(post_id,parent_id)`, `likes(visitor_token)`, `reviews(project_name)`, `events(type,target,ts)`, `events(session_id)`; add `LIMIT/OFFSET` + `?limit/cursor`, list-view column projection; check `rows.Err()` after every loop (currently missing everywhere).

### P1-4: CSS / fonts / infinite anims
- Files: `styles.css` (56KB/2721 lines), `light-mode.css` (7.5KB), `+layout.svelte:14-15,264-283`, `app.html:14-32`, `Skills.svelte:51`, `FeedbackMarquee.svelte:63-91`, `DevLog.svelte:226`
- Fix: split critical CSS, code-split admin/blog CSS; subset/self-host devicon (14 icons only); one font provider; drop API/S3 preconnect pre-paint; disable CSS `scroll-behavior:smooth` when Lenis on; pause offscreen infinite anims via IntersectionObserver; merge light-mode truths; remove duplicate `.skip-link` block.

### P1-5: Analytics + SWR cache
- Files: `analytics.ts:69-199`, `utils/cache.ts:54-84`, `CookieConsent.svelte:17`, `blog/+page.svelte:14,22`, `blog/[slug]/+page.svelte:43-59`
- Fix: `sendBeacon`/`keepalive` + queue+flush batching; memoize visitor token; dedup SWR map; pass `signal`; version cache key by `BASE`; fix `allTags`/`filtered` function-call-in-template recompute; throttle blog TOC scroll handler.

### P1-6: Dep / chunk hygiene
- Files: `vite.config.ts:6-14`, `package.json:18,30-34`, `+page.svelte:16`, `ui/Button:2,Input:2,Select:2,DatePicker:2,Card:2`, `showcase/+page.svelte:6,67,102-130`, `Magnetic.svelte:81-87`, `CustomCursor.svelte:45-64`, `service-worker.ts:1-80`, `Makefile:3-14`, `render.yaml:7-9`, `.env.example`
- Fix: prune `perfect-arrows` (zero imports); avoid `ui/` barrel pulling `svelte-motion` into homepage (direct imports); gate showcase rAF loops on `disabled/touch/reduced`; `go build -ldflags="-s -w" -trimpath`; document `ANALYTICS_SECRET/RUN_MIGRATIONS/AWS_*/S3_BUCKET`, pooler guidance; remove or fix dead service-worker.

## 3. P2 / hardening — must not regress

### H-1: Type safety
- `+page.svelte:163 (row:any)`, `:417 (social.icon as any)`, `Button:40/Card:29/Select:46/DatePicker:40 (...:any)`, `raw as {...}` null-unsafe (`+page:50,80,87`, `blog/+page:34`, `[slug]:40`, `CommentSection:55`), `catch(e:any)`, `(entry as any)`, `BlogComment.replies` non-optional vs `?.`, dead `githubStars`/`SkillItem.level`.
- Fix with narrowing + `?.` + `Array.isArray`, import shared `SkillItem`/`IconName` types.

### H-2: Security (do not ship without)
- `handler.go:39-48` `?key==ContactSecret` bypass when empty + key in URL/localStorage (`admin/contact/+page:13-17`). Require non-empty secret, move to `Authorization: Bearer`, expire localStorage.
- `markdown.ts:1-56` XSS (`javascript:`, `" onerror=`) via `{@html}` blog + admin previews. Add DOMPurify + allowlist + link `rel=noopener`.
- OAuth `state` never validated (`auth.go:87-91`), JWT in `#token` + localStorage no expiry (`+layout:37-43`, `auth.ts:76-83`). Validate state, check JWT shape/exp, validate `redirect_after_login` prefix.
- CORS hardcoded (`main.go:77-82`); track `FRONTEND_URL`, tighten methods. Contact/topic/message/comments/reviews/likes validation gaps → maxlength + allowlist.

### H-3: A11y (WCAG 2.2 AA)
- Marquee no reduced-motion guard + clones announced (add `aria-hidden`); article arrows hijack textarea (scope to carousel region); whole-card link no focus-visible; stars no radiogroup + targets < 44px; chips no `aria-pressed`; progress `aria-hidden` no live alt; `type=button` outside `<form>`; menu `aria-expanded`/`inert`/trap; badge modal no dialog semantics; skip-link broken on blog; `/blog` links ignore `base`.

### H-4: SEO / CSP
- Static `connect-src` misses Render backend host → blocked fetch. Sync with `VITE_API_URL`. Replace `unsafe-inline` where possible, add header CSP + `report-uri`.
- Canonical slash mismatch, `SITE_URL` wrong env key, sitemap lists hash sections + omits `/blog/[slug]`. Escape `</script>` in JSON-LD (`\u003c`). Move certs/blog SEO-critical data to `load()` for prerender.

### H-5: Resilience
- Empty vs error conflated (fallback masks outage, dead error branch `+page:352-367`); blog "No posts" on 500; like/contact no cancel/retry/toast; SWR no timeout; auth clears JWT on transient 500. Add `error` state + retry buttons, pending/disable, toasts, AbortController everywhere.

## 4. Swarm execution plan (4 waves, parallel agents)

- Wave 1 — Frontend perf swarm: F-P0-1, F-P0-2, F-P0-3, F-P0-4.
- Wave 2 — Backend reliability swarm: B-P0-1, B-P0-2, B-P0-3, B-P0-4.
- Wave 3 — P1 optimization swarm: P1-1 … P1-6.
- Wave 4 — Hardening swarm: H-1 … H-5.
- Rules for every agent: minimal diffs, no visual redesign, preserve `base:/portfolio-svelte`, run `npm run check` / `go vet`, keep reduced-motion + mobile paths working, update this file's checkbox when done.

## 5. Verification checklist

- [x] `npm run check` clean (0 errors, 75 pre-existing warnings)
- [x] `npm run build` succeeds (homepage 49KB, showcase 54KB, vendor-motion split)
- [x] Go files `gofmt -e` clean (`go vet` blocked by sandbox /tmp quota, not by code)
- [x] Lighthouse: LCP single eager image, CLS icon boxes, INP throttled handlers
- [x] Cold homepage ≤ 3 critical API calls (certs now, badges/marquee/likes on #projects approach); `/healthz` in render.yaml
- [x] No `?key` bypass (empty secret rejected), `{@html}` JSON-LD escaped, CSP includes Render host
