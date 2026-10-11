<script lang="ts">
    import { onMount } from "svelte";
    import { ensureGsap, motionAllowed } from "$lib/motion";

    import { theme } from "$lib/stores/theme";
    import { portfolioData, fallbackCertificates } from "$lib/data";
    import Hero from "$lib/components/Hero.svelte";
    import About from "$lib/components/About.svelte";
    import ProjectCard from "$lib/components/ProjectCard.svelte";
    import CertificateCard from "$lib/components/CertificateCard.svelte";

    import Skills from "$lib/components/Skills.svelte";
    import Footer from "$lib/components/Footer.svelte";
    import ContactForm from "$lib/components/ContactForm.svelte";
    // P1-6: direct imports — the ui barrel re-exports motion-heavy
    // components (Input/Select/DatePicker) that would join this chunk.
    import Icon from "$lib/components/ui/Icon/Icon.svelte";
    import SectionHeader from "$lib/components/ui/SectionHeader/SectionHeader.svelte";
    import Card from "$lib/components/ui/Card/Card.svelte";
    import Skeleton from "$lib/components/ui/Skeleton/Skeleton.svelte";
    import type { IconName } from "$lib/components/ui/Icon/Icon.svelte";
    import type { Certificate, Badge, MarqueeItem } from "$lib/types";
    import { fetchProjectLikes, toggleProjectLike } from "$lib/analytics";
    import FeedbackMarquee from "$lib/components/feedback/FeedbackMarquee.svelte";
    import { staleWhileRevalidate } from '$lib/utils/cache';

    const { profile, about, skills, projects } = portfolioData;
    const API_BASE = (import.meta.env.VITE_API_URL || "").replace(/\/$/, "");

    let certificates = $state<Certificate[]>([]);
    let certsLoading = $state(true);
    let certsError = $state("");
    let badges = $state<Badge[]>([]);
    let marqueeItems = $state<MarqueeItem[]>([]);
    let likes = $state<Record<string, number>>({});
    let userLikes = $state<Record<string, boolean>>({});
    let prefersReducedMotion = $state(false);

    async function fetchCerts() {
        if (!API_BASE) {
            certificates = fallbackCertificates;
            certsLoading = false;
            return;
        }
        // F-P0-4: critical fetch gets an AbortController timeout (8s).
        const ctrl = new AbortController();
        const timer = setTimeout(() => ctrl.abort(), 8000);
        try {
            const ok = await staleWhileRevalidate<Certificate[]>(
                'home:certs',
                `${API_BASE}/api/certificates`,
                30 * 60 * 1000,
                (data) => {
                    certificates = data.length > 0 ? data : fallbackCertificates;
                    if (data.length === 0) certsError = "No certificates found on server";
                    certsLoading = false;
                },
            (raw) => {
                const d = (raw as { certificates?: Certificate[] } | null);
                if (Array.isArray(raw)) return raw as Certificate[];
                const list = d?.certificates;
                return Array.isArray(list) ? list : [];
            },
                { signal: ctrl.signal },
            );
            if (!ok) {
                certificates = fallbackCertificates;
                certsLoading = false;
            }
        } finally {
            clearTimeout(timer);
        }
    }

    // F-P0-4: below-fold fetches (badges, marquee, likes) — deferred until
    // the projects section nears the viewport. Single-fire guard.
    let belowFoldFired = false;
    function loadBelowFold() {
        if (belowFoldFired || !API_BASE) return;
        belowFoldFired = true;
        void Promise.allSettled([
            staleWhileRevalidate<Badge[]>(
                'home:badges',
                `${API_BASE}/api/badges`,
                30 * 60 * 1000,
                (data) => { badges = data ?? []; },
                (raw) => { const d = (raw as { badges?: Badge[] } | null); const list = d?.badges; if (Array.isArray(list)) return list; if (Array.isArray(raw)) return (raw as Badge[]); return []; },
            ),
            staleWhileRevalidate<MarqueeItem[]>(
                'home:marquee',
                `${API_BASE}/api/feedback/marquee`,
                10 * 60 * 1000,
                (data) => { marqueeItems = data ?? []; },
                (raw) => { const d = (raw as { items?: MarqueeItem[] } | null); const list = d?.items; if (Array.isArray(list)) return list; if (Array.isArray(raw)) return (raw as MarqueeItem[]); return []; },
            ),
            fetchProjectLikes().then(result => {
                likes = result.likes;
                userLikes = result.user_likes;
            }).catch(() => {}),
        ]);
    }

    const siteUrl = "https://pranavagarkar07.github.io/portfolio-svelte/";
    const pageTitle = `${profile.name} | ${profile.role}`;
    const pageDescription = `${profile.tagline} Experienced in building secure, scalable applications with Django, React, Svelte, and Go.`;
    const pageKeywords =
        "Svelte, SvelteKit, Django, Go, Full Stack Developer, Portfolio, Pranav Agarkar, Web Development, TaskVault, BeamSync";
    const coverImage = `${siteUrl}og-image.png`;

    onMount(() => {
        prefersReducedMotion = window.matchMedia("(prefers-reduced-motion: reduce)").matches;

        // F-P0-4: critical fetch first (certs); everything below the fold
        // (badges, marquee, likes) waits for the projects section.
        void fetchCerts();

        let belowFoldObserver: IntersectionObserver | null = null;
        let belowFoldFallback: ReturnType<typeof setTimeout> | undefined;
        if (API_BASE && !belowFoldFired) {
            const fire = () => {
                belowFoldObserver?.disconnect();
                belowFoldObserver = null;
                if (belowFoldFallback) clearTimeout(belowFoldFallback);
                loadBelowFold();
            };
            if ("IntersectionObserver" in window) {
                const target = document.querySelector("#projects");
                if (target) {
                    belowFoldObserver = new IntersectionObserver(
                        (entries) => {
                            if (entries.some((e) => e.isIntersecting)) fire();
                        },
                        { rootMargin: "400px 0px" },
                    );
                    belowFoldObserver.observe(target);
                    // Fallback: never leave below-fold data unfetched.
                    belowFoldFallback = setTimeout(fire, 8000);
                } else {
                    loadBelowFold();
                }
            } else {
                loadBelowFold();
            }
        }

        // F-P0-1/F-P0-2: motion + reveals load dynamically and only when
        // allowed. Content is default-visible; gsap.from() is used so a
        // missing/failed animation can never leave content hidden.
        let destroyed = false;
        let resizeTimer: ReturnType<typeof setTimeout> | undefined;
        let progressBar: HTMLDivElement | null = null;
        let scrollTriggers: Array<{ kill: () => void }> = [];
        const onResize = (refresh: () => void) => () => {
            if (resizeTimer) clearTimeout(resizeTimer);
            resizeTimer = setTimeout(refresh, 200);
        };
        let handleResize: (() => void) | null = null;

        (async () => {
            if (!motionAllowed()) return;
            const mods = await ensureGsap();
            if (destroyed || !mods) return;
            const { gsap, ScrollTrigger } = mods;

            handleResize = onResize(() => ScrollTrigger.refresh());
            window.addEventListener("resize", handleResize);

        // ──────────────────────────────────────────────
        // 1. READING PROGRESS BAR
        // ──────────────────────────────────────────────
        progressBar = document.createElement("div");
        Object.assign(progressBar.style, {
            position: "fixed", top: "0", left: "0", width: "0%", height: "2px",
            background: "linear-gradient(90deg, var(--accent), #ff8800, var(--accent))",
            zIndex: "9999", pointerEvents: "none", transformOrigin: "0 50%",
            transition: "transform 0.1s linear",
        });
        document.body.prepend(progressBar);
        scrollTriggers.push(ScrollTrigger.create({ trigger: "body", start: "top top", end: "bottom bottom", onUpdate: (self: { progress: number }) => {
                if (progressBar) progressBar.style.transform = `scaleX(${self.progress})`;
        }}));

        // ──────────────────────────────────────────────
        // 6. PROJECTS — STAGGER REVEAL (from: default-visible)
        // ──────────────────────────────────────────────
        scrollTriggers.push(ScrollTrigger.create({
            trigger: "#projects", start: "top 85%", once: true,
            onEnter: () => { gsap.fromTo(".project-card",
                { y: 40, opacity: 0 },
                { y: 0, opacity: 1,
                    stagger: 0.08,
                    ease: "power3.out", duration: 0.6,
                    overwrite: "auto",
                });}
        }));

        // About columns
        scrollTriggers.push(ScrollTrigger.create({
            trigger: "#about", start: "top 80%", once: true,
            onEnter: () => {
                gsap.fromTo(".about-visual-col",
                    { x: -30, opacity: 0 },
                    { x: 0, opacity: 1,
                        ease: "power3.out", duration: 0.6, overwrite: "auto",
                    });
                gsap.fromTo(".about-content-col",
                    { x: 30, opacity: 0 },
                    { x: 0, opacity: 1,
                        ease: "power3.out", duration: 0.6, delay: 0.08, overwrite: "auto",
                    });
                // Spec rows sweep up (batched: inside #about trigger)
                gsap.utils.toArray(".about-spec-row").forEach((row: HTMLElement) => {
                    gsap.fromTo(row as HTMLElement,
                        { y: 25, opacity: 0 },
                        { y: 0, opacity: 1, ease: "power3.out", duration: 0.5, overwrite: "auto" });
                });
            }
        }));

        // ──────────────────────────────────────────────
        // 11. ABOUT — METRICS COUNT-UP
        // ──────────────────────────────────────────────
        scrollTriggers.push(ScrollTrigger.create({
            trigger: "#about", start: "top 75%", once: true,
            onEnter: () => {
                // Count-up metrics
                (gsap.utils.toArray(".about-metric-value") as HTMLElement[]).forEach((el: HTMLElement) => {
                    const raw = el.textContent || "0";
                    const num = parseFloat(raw.replace(/[^0-9.]/g, ""));
                    if (isNaN(num)) return;
                    const suffix = raw.replace(/[0-9.]/g, "");
                    const obj = { val: 0 };
                    gsap.to(obj, {
                        val: num, duration: 1.5, ease: "power3.out",
                        onUpdate: () => { el.textContent = Math.round(obj.val) + suffix; },
                    });
                });
                // Metrics fade up
                gsap.fromTo(".about-metric",
                    { y: 20, opacity: 0 },
                    { y: 0, opacity: 1,
                        stagger: 0.08, ease: "power3.out", duration: 0.5, overwrite: "auto",
                    });
            }
        }));

        // ──────────────────────────────────────────────
        // 12-13. CONTACT — ALL (batched: 1 trigger instead of 7)
        // ──────────────────────────────────────────────

        // Contact section single trigger
        scrollTriggers.push(ScrollTrigger.create({
            trigger: "#contact", start: "top 85%", once: true,
            onEnter: () => {
                // Simple fade up
                gsap.fromTo(".contact-info",
                    { y: 15, opacity: 0 },
                    { y: 0, opacity: 1, ease: "power3.out", duration: 0.5, overwrite: "auto" });
                gsap.fromTo(".contact-form-container",
                    { y: 15, opacity: 0 },
                    { y: 0, opacity: 1, ease: "power3.out", duration: 0.5, delay: 0.1, overwrite: "auto" });
                gsap.fromTo(".info-item",
                    { y: 10, opacity: 0 },
                    { y: 0, opacity: 1, ease: "power3.out", duration: 0.4,
                        stagger: 0.05, overwrite: "auto",
                    });
                gsap.fromTo(".social-link",
                    { y: 10, opacity: 0 },
                    { y: 0, opacity: 1, ease: "power3.out", duration: 0.4,
                        stagger: 0.04, overwrite: "auto",
                    });
            }
        }));

        // ──────────────────────────────────────────────
        // 14. SAFETY — 1s only: refresh triggers so stuck reveals recover.
        // No clipPath/opacity forcing (content is default-visible).
        // ──────────────────────────────────────────────
        setTimeout(() => {
            if (destroyed) return;
            try { ScrollTrigger.refresh(); } catch {}
        }, 1000);
        })();

        return () => {
            destroyed = true;
            belowFoldObserver?.disconnect();
            if (belowFoldFallback) clearTimeout(belowFoldFallback);
            if (handleResize) window.removeEventListener("resize", handleResize);
            if (resizeTimer) clearTimeout(resizeTimer);
            scrollTriggers.forEach((t) => { try { t.kill(); } catch {} });
            progressBar?.remove();
        };
    });


</script>

<svelte:head>
    <!-- Primary Meta Tags -->
    <title>{pageTitle}</title>
    <meta name="title" content={pageTitle} />
    <meta name="description" content={pageDescription} />
    <meta name="keywords" content={pageKeywords} />
    <meta name="author" content={profile.name} />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <meta name="robots" content="index, follow" />

    <!-- Open Graph / Facebook -->
    <meta property="og:type" content="website" />
    <meta property="og:url" content={siteUrl} />
    <meta property="og:title" content={pageTitle} />
    <meta property="og:description" content={pageDescription} />
    <meta property="og:image" content={coverImage} />
    <meta property="og:image:width" content="1200" />
    <meta property="og:image:height" content="630" />
    <meta property="og:image:type" content="image/png" />
    <meta property="og:image:alt" content="Pranav Agarkar — Fullstack Developer Portfolio" />
    <meta property="og:site_name" content="Pranav Agarkar Portfolio" />
    <meta property="og:locale" content="en_US" />

    <!-- Twitter -->
    <meta property="twitter:card" content="summary_large_image" />
    <meta property="twitter:url" content={siteUrl} />
    <meta property="twitter:title" content={pageTitle} />
    <meta property="twitter:description" content={pageDescription} />
    <meta property="twitter:image" content={coverImage} />

    <!-- Structured Data (JSON-LD) -->
    {@html `<noscript><style>.blur-segment,.hero-status,.hero-name-text,.hero-tagline-wrap,.cta-group,.tech-marquee,.hero-scroll,.project-card,.about-visual-col,.about-content-col,.about-spec-row,.about-metric,.contact-info,.contact-form-container,.info-item,.social-link,.skill-category-card,.skill-card{opacity:1!important;transform:none!important;filter:none!important;clip-path:none!important}</style></noscript>`}
    <script type="application/ld+json">
        {
            "@context": "https://schema.org",
            "@type": "Person",
            "name": "{profile.name}",
            "url": "{siteUrl}",
            "jobTitle": "{profile.role}",
            "image": "{coverImage}",
            "sameAs": [
                {#each profile.socials as social, i}
                    "{social.url}"{i < profile.socials.length - 1 ? ',' : ''}
                {/each}
            ],
            "description": "{pageDescription}"
        }
    </script>
</svelte:head>


<main id="main-content">
    <Hero {profile} {skills} {about} />

    <FeedbackMarquee items={marqueeItems} />

    <section id="projects" class="section-container snap-section">
        <SectionHeader title="Featured Projects" count={projects.length} animate />
        <div class="projects-grid" style="margin-top: 1.5rem">
            {#each projects as project, i}
                <ProjectCard {project} index={i} eager={i === 0} liked={userLikes[project.name] ?? false} likeCount={likes[project.name] ?? 0} onLike={(liked) => {
                    const name = project.name;
                    const prevLikes = likes;
                    userLikes = { ...userLikes, [name]: liked };
                    likes = { ...likes, [name]: (likes[name] ?? 0) + (liked ? 1 : -1) };
                    toggleProjectLike(name, liked).then(count => {
                        if (count < 0) {
                            userLikes = { ...userLikes, [name]: !liked };
                            likes = { ...likes, [name]: (likes[name] ?? 0) + (liked ? -1 : 1) };
                        } else {
                            likes = { ...likes, [name]: count };
                        }
                    });
                }} />
            {/each}
        </div>
    </section>

    <Skills {skills} />

    <About {profile} {about} {badges} />

    <section id="certifications" class="section-container snap-section">
        <SectionHeader title="Certifications" count={certsLoading ? undefined : certificates.length} animate />
        {#if certsLoading}
            <div class="certs-skeleton-grid">
                {#each Array(3) as _, i}
                    <div class="cert-card-wrap">
                        <div class="cert-skeleton">
                            <div class="skeleton-thumb"></div>
                            <div class="skeleton-body">
                                <Skeleton width="40%" height="10px" />
                                <Skeleton width="85%" height="16px" />
                                <Skeleton width="55%" height="10px" />
                                <div class="skeleton-tags">
                                    <Skeleton width="50px" height="18px" />
                                    <Skeleton width="60px" height="18px" />
                                </div>
                                <Skeleton width="110px" height="24px" />
                            </div>
                        </div>
                    </div>
                {/each}
            </div>
        {:else if certificates.length > 0}
            <div class="certs-grid">
                {#each certificates as cert, i}
                    <div class="cert-card-wrap">
                        <CertificateCard certificate={cert} index={i} />
                    </div>
                {/each}
            </div>
        {:else if certsError}
            <div class="certs-error">
                <span class="certs-error-text">&#9888; Failed to load certificates</span>
                <button class="certs-retry-btn" onclick={() => { certsLoading = true; certsError = ''; fetchCerts(); }}>
                    Retry
                </button>
            </div>
        {/if}
    </section>

    <section id="contact" class="contact-section snap-section">
        <div class="contact-container">
            <div class="contact-header">
                <span class="contact-label">Get in Touch</span>
                <h2>Let's Work Together</h2>
                <p class="contact-description">Have a project in mind? I'd love to hear about it. Send me a message and I'll get back to you within 24 hours.</p>
            </div>

            <div class="contact-grid">
                <div class="contact-info">
                    <Card>
                        <a href="mailto:pranavagarkar8@gmail.com?subject=Opportunity%20from%20your%20portfolio" class="contact-email-cta">
                            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                                <rect x="2" y="4" width="20" height="16" rx="2"></rect>
                                <path d="M22 4L12 13 2 4"></path>
                            </svg>
                            <div class="contact-email-content">
                                <span class="contact-email-label">Email me directly</span>
                                <span class="contact-email-addr">pranavagarkar8@gmail.com</span>
                            </div>
                            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="contact-email-arrow">
                                <line x1="7" y1="17" x2="17" y2="7"></line>
                                <polyline points="7 7 17 7 17 17"></polyline>
                            </svg>
                        </a>
                        <div class="info-item">
                            <div class="info-icon">
                                <Icon name="clock" size={20} />
                            </div>
                            <div class="info-content">
                                <span class="info-title">Response Time</span>
                                <span class="info-text">Within 24 hours</span>
                            </div>
                        </div>
                    </Card>

                    <div class="social-section">
                        <span class="social-label">Connect</span>
                        <div class="social-links">
                            {#each profile.socials as social}
                                <a
                                    href={social.url}
                                    target={social.url.startsWith("mailto") ? undefined : "_blank"}
                                    rel="noopener"
                                    class="social-link"
                                    aria-label={social.label}
                                >
                                    <Icon name={social.icon as IconName} size={18} />
                                </a>
                            {/each}
                        </div>
                    </div>
                </div>

                <div class="contact-form-container">
                    <ContactForm />
                </div>
            </div>
        </div>
    </section>
</main>

<Footer {profile} />

<style>
</style>
