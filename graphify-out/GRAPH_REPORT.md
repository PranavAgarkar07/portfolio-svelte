# Graph Report - .  (2026-05-12)

## Corpus Check
- 40 files · ~105,992 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 142 nodes · 190 edges · 20 communities detected
- Extraction: 83% EXTRACTED · 17% INFERRED · 0% AMBIGUOUS · INFERRED: 32 edges (avg confidence: 0.81)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_Core Backend & Config|Core Backend & Config]]
- [[_COMMUNITY_Personal Profile & Skills|Personal Profile & Skills]]
- [[_COMMUNITY_Deployment & Infrastructure|Deployment & Infrastructure]]
- [[_COMMUNITY_TaskVault App|TaskVault App]]
- [[_COMMUNITY_BeamSync App|BeamSync App]]
- [[_COMMUNITY_BeamSync Onboarding|BeamSync Onboarding]]
- [[_COMMUNITY_BeamSync Technical|BeamSync Technical]]
- [[_COMMUNITY_TaskVault Security UI|TaskVault Security UI]]
- [[_COMMUNITY_Resume Skills|Resume Skills]]
- [[_COMMUNITY_Avatar|Avatar]]
- [[_COMMUNITY_Site Branding|Site Branding]]
- [[_COMMUNITY_TaskVault Login|TaskVault Login]]
- [[_COMMUNITY_TaskVault Dashboard|TaskVault Dashboard]]
- [[_COMMUNITY_TaskVault Homepage|TaskVault Homepage]]
- [[_COMMUNITY_Scroll Store|Scroll Store]]
- [[_COMMUNITY_Svelte CLI|Svelte CLI]]
- [[_COMMUNITY_SvelteKit Adapter|SvelteKit Adapter]]
- [[_COMMUNITY_Robots Config|Robots Config]]
- [[_COMMUNITY_Linux Scripts|Linux Scripts]]
- [[_COMMUNITY_Python Skill|Python Skill]]

## God Nodes (most connected - your core abstractions)
1. `Pranav Agarkar` - 24 edges
2. `Community Structure` - 21 edges
3. `God Nodes Analysis` - 10 edges
4. `generateDevLog()` - 7 edges
5. `Knowledge Gap Analysis` - 6 edges
6. `getInitialTheme()` - 5 edges
7. `createThemeStore()` - 5 edges
8. `DevLogResponse` - 5 edges
9. `CachedData` - 5 edges
10. `handleStatus()` - 5 edges

## Surprising Connections (you probably didn't know these)
- `Go Programming Language` --conceptually_related_to--> `Go Sentinel Backend Service`  [INFERRED]
  static/Pranav_Agarkar_Resume.pdf → DEPLOY.md
- `Portfolio Website (portfolio-svelte)` --conceptually_related_to--> `VITE_API_URL Repository Variable`  [INFERRED]
  static/Pranav_Agarkar_Resume.pdf → DEPLOY.md
- `Portfolio Website (portfolio-svelte)` --references--> `Svelte Portfolio Project`  [INFERRED]
  static/Pranav_Agarkar_Resume.pdf → README.md
- `Svelte Frontend Framework` --references--> `Svelte Portfolio Project`  [INFERRED]
  static/Pranav_Agarkar_Resume.pdf → README.md
- `Avatar Portrait (WebP)` --conceptually_related_to--> `Pranav Agarkar`  [INFERRED]
  src/lib/assets/avatar.webp → static/Pranav_Agarkar_Resume.pdf

## Hyperedges (group relationships)
- **Two-Tier Deployment Architecture** — deploy_backend_render, deploy_frontend_github_pages, deploy_go_sentinel, deploy_github_actions, deploy_render_blueprint, deploy_gemini_api_key, deploy_vite_api_url, deploy_twotier_rationale [EXTRACTED 1.00]
- **Backend Technologies** — tech_django, tech_go, tech_postgresql [EXTRACTED 1.00]
- **Frontend Technologies** — tech_react, tech_svelte, tech_typescript [EXTRACTED 1.00]
- **Cloud Platforms** — tech_render, tech_vercel, tech_docker [EXTRACTED 1.00]
- **Security Tools** — tech_jwt, tech_oauth, tech_fernet [EXTRACTED 1.00]
- **Featured Projects** — project_taskvault_lite, project_beamsync [EXTRACTED 1.00]
- **Contact Information** — contact_email, contact_github, contact_linkedin [EXTRACTED 1.00]
- **BeamSync P2P Transfer Flow (UI → Discovery → Pairing → Transfer)** — beamsync_cyberpunk_ui, beamsync_auto_ip_discovery, beamsync_qr_pairing, beamsync_p2p_transfer [INFERRED 0.80]
- **Community Cluster Assignments** — backend_cacheddata, backend_circuitbreaker, backend_contactsubmission, backend_devlogresponse, backend_metrics, main_callgemini, main_fetchgithubevents, main_generatedevlog, main_handlestatus, theme_createthemestore, theme_getinitialtheme, main_handleadmincontact, main_htmlescape, svelte_config_js, vite_config_ts, src_app_d_ts, src_lib_data_ts, src_lib_index_ts, src_lib_types_ts, src_lib_components_ui_index_ts, src_routes_layout_ts [EXTRACTED 1.00]
- **Dev Log Generation Pipeline** — main_generatedevlog, main_fetchgithubevents, main_callgemini [INFERRED 0.90]
- **Backend Core Data Types** — backend_devlogresponse, backend_cacheddata, backend_circuitbreaker, backend_metrics, backend_contactsubmission [INFERRED 0.85]
- **TaskVault Portfolio Showcase** — secureDBwithEncrypetedTasks_png, TaskVault, SecureDatabaseEncryptedTasks, PortfolioSvelteSite [EXTRACTED 1.00]
- **BeamSync Splash Page Composition** — beamsync_app_icon, beamsync_phone_mockup, beamsync_dark_theme, beamsync_orange_accent_color, beamsync_blue_gray_backdrop [INFERRED 0.85]

## Communities

### Community 0 - "Core Backend & Config"
Cohesion: 0.14
Nodes (18): CachedData, CircuitBreaker, ContactSubmission, DevLogResponse, Metrics, Community Structure, Corpus Verdict, God Nodes Analysis (+10 more)

### Community 1 - "Personal Profile & Skills"
Cohesion: 0.09
Nodes (26): Email: pranavagarkar8@gmail.com, GitHub: PranavAgarkar07, LinkedIn: pranavagarkar, 2+ Years Experience, Portfolio OG Image PNG, Portfolio OG Image SVG, Portfolio Website, Pranav Agarkar (+18 more)

### Community 2 - "Deployment & Infrastructure"
Cohesion: 0.12
Nodes (18): Backend Deployment on Render, DevLog Status Check System, Frontend Deployment on GitHub Pages, GEMINI_API_KEY Environment Variable, GitHub Actions Deploy Workflow, Go Sentinel Backend Service, Render Blueprint Config (render.yaml), Two-Tier Deployment Rationale (+10 more)

### Community 3 - "TaskVault App"
Cohesion: 0.2
Nodes (10): Fernet AES-128-CBC Encryption, JWT + OAuth2 Authentication, PostgreSQL SSL Database, API Throughput Metric (1,247 req/min), TaskVault Dashboard UI, TaskVault Lite, TaskVault Preview Image, Security Status Panel (+2 more)

### Community 4 - "BeamSync App"
Cohesion: 0.31
Nodes (10): BeamSync File Transfer Interface, BeamSync File Transfer Interface, BeamSync Peer Discovery and Connection, BeamSync Peer Discovery and Connection, BeamSync QR Sync and Mobile Pairing, BeamSync QR Sync and Mobile Pairing, Auto-IP Detection and Dynamic Port Scouting, Cyberpunk Terminal UI Design (+2 more)

### Community 5 - "BeamSync Onboarding"
Cohesion: 0.31
Nodes (9): BeamSync App Icon, BeamSync App Icon and Starting Page, BeamSync App Icon and Starting Page (WebP), BeamSync Blue-Gray Backdrop Color, BeamSync Dark Theme Design, BeamSync Orange Accent Color (#FF8A00), BeamSync Phone UI Mockup, BeamSync Project (+1 more)

### Community 6 - "BeamSync Technical"
Cohesion: 0.25
Nodes (8): AES-256-GCM Encryption Cipher, BeamSync P2P Transfer, Network Topology Panel, BeamSync Preview Image, BeamSync Terminal Interface, Transfer Stats Panel, Node_A Peer (192.168.1.15), Node_B Peer (192.168.1.22)

### Community 7 - "TaskVault Security UI"
Cohesion: 0.7
Nodes (4): Portfolio Svelte Site, Secure Database with Encrypted Tasks, Task Management User Interface, TaskVault Project

### Community 8 - "Resume Skills"
Cohesion: 0.67
Nodes (3): Django Web Framework, React Frontend Library, TaskVault Lite Project

### Community 9 - "Avatar"
Cohesion: 1.0
Nodes (3): Pranav Agarkar, Avatar Portrait, Avatar Portrait (WebP)

### Community 10 - "Site Branding"
Cohesion: 0.67
Nodes (3): Favicon PNG, Favicon SVG, Svelte Logo

### Community 11 - "TaskVault Login"
Cohesion: 1.0
Nodes (2): TaskVault Lite Login Page Screenshot, TaskVault Lite Login Page Screenshot (WebP)

### Community 12 - "TaskVault Dashboard"
Cohesion: 1.0
Nodes (2): TaskVault Lite Main Dashboard Screenshot, TaskVault Lite Main Dashboard Screenshot (WebP)

### Community 13 - "TaskVault Homepage"
Cohesion: 1.0
Nodes (2): TaskVault Lite Homepage Screenshot, TaskVault Lite Homepage Screenshot (WebP)

### Community 14 - "Scroll Store"
Cohesion: 1.0
Nodes (0): 

### Community 15 - "Svelte CLI"
Cohesion: 1.0
Nodes (1): sv CLI Tool

### Community 16 - "SvelteKit Adapter"
Cohesion: 1.0
Nodes (1): SvelteKit Adapter Concept

### Community 17 - "Robots Config"
Cohesion: 1.0
Nodes (1): Robots.txt Allow All Crawling

### Community 18 - "Linux Scripts"
Cohesion: 1.0
Nodes (1): Linux Optimization Scripts Project

### Community 19 - "Python Skill"
Cohesion: 1.0
Nodes (1): Python Programming Language

## Knowledge Gaps
- **56 isolated node(s):** `Render Blueprint Config (render.yaml)`, `GitHub Actions Deploy Workflow`, `GEMINI_API_KEY Environment Variable`, `DevLog Status Check System`, `sv CLI Tool` (+51 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **Thin community `TaskVault Login`** (2 nodes): `TaskVault Lite Login Page Screenshot`, `TaskVault Lite Login Page Screenshot (WebP)`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `TaskVault Dashboard`** (2 nodes): `TaskVault Lite Main Dashboard Screenshot`, `TaskVault Lite Main Dashboard Screenshot (WebP)`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `TaskVault Homepage`** (2 nodes): `TaskVault Lite Homepage Screenshot`, `TaskVault Lite Homepage Screenshot (WebP)`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Scroll Store`** (1 nodes): `scroll.ts`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Svelte CLI`** (1 nodes): `sv CLI Tool`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `SvelteKit Adapter`** (1 nodes): `SvelteKit Adapter Concept`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Robots Config`** (1 nodes): `Robots.txt Allow All Crawling`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Linux Scripts`** (1 nodes): `Linux Optimization Scripts Project`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.
- **Thin community `Python Skill`** (1 nodes): `Python Programming Language`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Avatar Portrait (WebP)` connect `Avatar` to `Core Backend & Config`?**
  _High betweenness centrality (0.007) - this node is a cross-community bridge._
- **What connects `Render Blueprint Config (render.yaml)`, `GitHub Actions Deploy Workflow`, `GEMINI_API_KEY Environment Variable` to the rest of the system?**
  _56 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Core Backend & Config` be split into smaller, more focused modules?**
  _Cohesion score 0.14 - nodes in this community are weakly interconnected._
- **Should `Personal Profile & Skills` be split into smaller, more focused modules?**
  _Cohesion score 0.09 - nodes in this community are weakly interconnected._
- **Should `Deployment & Infrastructure` be split into smaller, more focused modules?**
  _Cohesion score 0.12 - nodes in this community are weakly interconnected._