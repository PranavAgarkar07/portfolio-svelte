<script lang="ts">
    import { onMount } from "svelte";
    import gsap from "gsap";
    import { ScrollTrigger } from "gsap/dist/ScrollTrigger";

    import { theme } from "$lib/stores/theme";
    import { portfolioData } from "$lib/data";
    import Header from "$lib/components/Header.svelte";
    import Hero from "$lib/components/Hero.svelte";
    import ProjectCard from "$lib/components/ProjectCard.svelte";
    import CertificateCard from "$lib/components/CertificateCard.svelte";

    import DevLog from "$lib/components/DevLog.svelte";
    import Footer from "$lib/components/Footer.svelte";
    import ContactForm from "$lib/components/ContactForm.svelte";
    import { Icon, SectionHeader, Card, Skeleton } from "$lib/components/ui";
    import type { Certificate, Badge } from "$lib/types";

    const { profile, about, skills, projects } = portfolioData;
    const API_BASE = (import.meta.env.VITE_API_URL || "").replace(/\/$/, "");

    let certificates = $state<Certificate[]>([]);
    let certsLoading = $state(true);
    let certsError = $state("");
    let badges = $state<Badge[]>([]);

    const siteUrl = "https://pranavagarkar07.github.io/portfolio-svelte/";
    const pageTitle = `${profile.name} | ${profile.role}`;
    const pageDescription = `${profile.tagline} Experienced in building secure, scalable applications with Django, React, Svelte, and Go.`;
    const pageKeywords =
        "Svelte, SvelteKit, Django, Go, Full Stack Developer, Portfolio, Pranav Agarkar, Web Development, TaskVault, BeamSync";
    const coverImage = `${siteUrl}avatar.png`; // Fallback to avatar if no dedicated cover

    onMount(() => {
        const prefersReducedMotion = window.matchMedia("(prefers-reduced-motion: reduce)").matches;

        if (API_BASE) {
            fetch(`${API_BASE}/api/certificates`)
                .then(r => r.ok ? r.json() : Promise.reject("HTTP " + r.status))
                .then(data => { certificates = data.certificates ?? data ?? []; certsLoading = false; })
                .catch(e => { certsError = String(e); certsLoading = false; });

            fetch(`${API_BASE}/api/badges`)
                .then(r => r.ok ? r.json() : Promise.reject("HTTP " + r.status))
                .then(data => { badges = data.badges ?? data ?? []; })
                .catch(() => {});
        } else {
            certsLoading = false;
        }

        if (prefersReducedMotion) return;
        gsap.registerPlugin(ScrollTrigger);
        gsap.defaults({ overwrite: "auto" });

        let resizeTimer: ReturnType<typeof setTimeout>;
        window.addEventListener("resize", () => {
            clearTimeout(resizeTimer);
            resizeTimer = setTimeout(() => ScrollTrigger.refresh(), 200);
        });

        // ──────────────────────────────────────────────
        // 1. READING PROGRESS BAR
        // ──────────────────────────────────────────────
        const progressBar = document.createElement("div");
        Object.assign(progressBar.style, {
            position: "fixed", top: "0", left: "0", width: "0%", height: "2px",
            background: "linear-gradient(90deg, var(--accent), #ff8800, var(--accent))",
            zIndex: "9999", pointerEvents: "none", transformOrigin: "0 50%",
            transition: "transform 0.1s linear",
        });
        document.body.prepend(progressBar);
        ScrollTrigger.create({ trigger: "body", start: "top top", end: "bottom bottom", onUpdate: self => {
                progressBar.style.transform = `scaleX(${self.progress})`;
        }});

        // ──────────────────────────────────────────────
        // 6. PROJECTS — CARD PERSPECTIVE + STAGGER REVEAL
        // ──────────────────────────────────────────────
        gsap.utils.toArray(".project-card").forEach((card: any, i: number) => {
            const dir = i % 2 === 0 ? -1 : 1;
            gsap.set(card, { x: dir * 150, opacity: 0, rotateY: dir * 15, scale: 0.9 });
            ScrollTrigger.create({
                trigger: card as Element, start: "top 85%", once: true,
                onEnter: () => { gsap.to(card, {
                    x: 0, opacity: 1, rotateY: 0, scale: 1,
                    ease: "elastic.out(1, 0.65)", duration: 0.8,
                });}
            });
        });

        // ──────────────────────────────────────────────
        // 7. SKILLS — CASCADE WAVE REVEAL
        // ──────────────────────────────────────────────
        gsap.set(".skill-category-card", { y: 80, opacity: 0, scale: 0.92, rotateX: 10 });
        ScrollTrigger.create({
            trigger: "#skills", start: "top 85%", once: true,
            onEnter: () => { gsap.to(".skill-category-card", {
                y: 0, opacity: 1, scale: 1, rotateX: 0,
                stagger: { each: 0.15, from: "edges", ease: "power2.inOut" },
                ease: "power3.out", duration: 0.7,
            });}
        });

        // ──────────────────────────────────────────────
        // 8. SKILL ITEMS — INDIVIDUAL STAGGER POP
        // ──────────────────────────────────────────────
        gsap.utils.toArray(".skill-card").forEach((card: any, i: number) => {
            gsap.set(card, { y: 30, opacity: 0, scale: 0.85 });
            ScrollTrigger.create({
                trigger: card as Element, start: "top 92%", once: true,
                onEnter: () => { gsap.to(card, {
                    y: 0, opacity: 1, scale: 1,
                    ease: "back.out(2)", duration: 0.5,
                });}
            });
        });

        // ──────────────────────────────────────────────
        // CERTIFICATES — FILE DRAWER REVEAL (batched: 1 trigger instead of 13)
        // ──────────────────────────────────────────────
        const certCards = gsap.utils.toArray(".cert-card-wrap");
        if (certCards.length > 0) {
            const dirs = certCards.map((_: any, i: number) => i % 2 === 0 ? -1 : 1);
            gsap.set(certCards, { x: (i: number) => dirs[i] * 80, y: 60, opacity: 0, rotateX: 5, scale: 0.95 });
            ScrollTrigger.create({
                trigger: "#certifications", start: "top 85%", once: true,
                onEnter: () => { gsap.to(certCards, {
                    x: 0, y: 0, opacity: 1, rotateX: 0, scale: 1,
                    stagger: 0.06, ease: "power3.out", duration: 0.65,
                });}
            });
        }

        // Handled by motion.div whileInView in SectionHeader.svelte
        // (GSAP version removed to avoid duplicate clip-path animation)

        // ──────────────────────────────────────────────
        // 10. ABOUT — DEPTH LAYER REVEAL
        // ──────────────────────────────────────────────
        gsap.set(".about-visual-col", { x: -80, opacity: 0, rotateY: 8, scale: 0.95 });
        gsap.set(".about-content-col", { x: 80, opacity: 0, rotateY: -8, scale: 0.95 });
        ScrollTrigger.create({
            trigger: "#about", start: "top 80%", once: true,
            onEnter: () => {
                gsap.to(".about-visual-col", {
                    x: 0, opacity: 1, rotateY: 0, scale: 1,
                    ease: "power3.out", duration: 0.7,
                });
                gsap.to(".about-content-col", {
                    x: 0, opacity: 1, rotateY: 0, scale: 1,
                    ease: "power3.out", duration: 0.7, delay: 0.1,
                });
            }
        });

        // About specs — sweep up
        gsap.utils.toArray(".about-spec-row").forEach((row: any, i: number) => {
            gsap.set(row, { y: 25, opacity: 0 });
            ScrollTrigger.create({
                trigger: row as Element, start: "top 80%", once: true,
                onEnter: () => { gsap.to(row, {
                    y: 0, opacity: 1,
                    ease: "power3.out", duration: 0.5,
                });}
            });
        });

        // ──────────────────────────────────────────────
        // 11. ABOUT — COUNT-UP METRICS
        // ──────────────────────────────────────────────
        gsap.utils.toArray<HTMLElement>(".about-metric-value").forEach((el) => {
            const raw = el.textContent || "0";
            const num = parseFloat(raw.replace(/[^0-9.]/g, ""));
            if (isNaN(num)) return;
            const suffix = raw.replace(/[0-9.]/g, "");
            const display = el;
            ScrollTrigger.create({
                trigger: el, start: "top 85%", once: true,
                onEnter: () => {
                    let obj = { val: 0 };
                    gsap.to(obj, {
                        val: num, duration: 1.5, ease: "power3.out",
                        onUpdate: () => { display.textContent = Math.round(obj.val) + suffix; },
                    });
                }
            });
        });

        // About metrics — spring up
        gsap.utils.toArray(".about-metric").forEach((metric: any, i: number) => {
            gsap.set(metric, { y: 40, opacity: 0, scale: 0.85 });
            ScrollTrigger.create({
                trigger: "#about", start: "top 75%", once: true,
                onEnter: () => { gsap.to(metric, {
                    y: 0, opacity: 1, scale: 1,
                    ease: "power3.out", duration: 0.6,
                    delay: i * 0.1,
                });}
            });
        });

        // ──────────────────────────────────────────────
        // 12. CONTACT — DUAL-SIDED REVEAL
        // ──────────────────────────────────────────────
        gsap.set(".contact-info", { x: -50, opacity: 0 });
        gsap.set(".contact-form-container", { x: 50, opacity: 0 });
        ScrollTrigger.create({
            trigger: "#contact", start: "top 85%", once: true,
            onEnter: () => {
                gsap.to(".contact-info", {
                    x: 0, opacity: 1,
                    ease: "power3.out", duration: 0.7,
                });
                gsap.to(".contact-form-container", {
                    x: 0, opacity: 1,
                    ease: "power3.out", duration: 0.7, delay: 0.15,
                });
            }
        });

        // Contact info items — staggered
        gsap.utils.toArray(".info-item").forEach((item: any, i: number) => {
            gsap.set(item, { y: 20, opacity: 0 });
            ScrollTrigger.create({
                trigger: "#contact", start: "top 85%", once: true,
                onEnter: () => { gsap.to(item, {
                    y: 0, opacity: 1,
                    ease: "back.out(1.7)", duration: 0.5, delay: i * 0.06,
                });}
            });
        });

        // ──────────────────────────────────────────────
        // 13. SOCIAL LINKS — STAGGER REVEAL
        // ──────────────────────────────────────────────
        gsap.utils.toArray(".social-link").forEach((link: any, i: number) => {
            gsap.set(link, { y: 15, opacity: 0, scale: 0.8 });
            ScrollTrigger.create({
                trigger: "#contact", start: "top 85%", once: true,
                onEnter: () => { gsap.to(link, {
                    y: 0, opacity: 1, scale: 1,
                    ease: "power3.out", duration: 0.5, delay: i * 0.05,
                });}
            });
        });

        // ──────────────────────────────────────────────
        // 14. MOBILE SAFETY — recover stuck hidden elements
        // ──────────────────────────────────────────────
        function forceVisibility() {
            ScrollTrigger.refresh();
            const selectors = [
                ".section-header", ".skill-category-card", ".skill-card", ".project-card",
                ".about-visual-col", ".about-content-col", ".about-spec-row",
                ".about-metric", ".contact-info", ".contact-form-container",
                ".info-item", ".social-link", ".cert-card-wrap",
                ".hero-status", ".hero-scroll"
            ];
            document.querySelectorAll(selectors.join(",")).forEach(el => {
                const computed = getComputedStyle(el);
                if (parseFloat(computed.opacity) < 0.5) {
                    gsap.set(el, {
                        opacity: 1, y: 0, x: 0, scale: 1,
                        rotateY: 0, rotateX: 0,
                        clipPath: "polygon(0% 0%, 100% 0%, 100% 100%, 0% 100%)"
                    });
                }
});
        gsap.set(".hero-name-text", {
            clipPath: "polygon(0% 0%, 100% 0%, 100% 100%, 0% 100%)"
        });
        }

        const isMobile = window.innerWidth <= 768 || "ontouchstart" in window;
        if (isMobile) {
            setTimeout(forceVisibility, 1500);
            window.addEventListener("touchstart", () => ScrollTrigger.refresh(), { once: true });
        }
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

    <!-- Twitter -->
    <meta property="twitter:card" content="summary_large_image" />
    <meta property="twitter:url" content={siteUrl} />
    <meta property="twitter:title" content={pageTitle} />
    <meta property="twitter:description" content={pageDescription} />
    <meta property="twitter:image" content={coverImage} />

    <!-- Structured Data (JSON-LD) -->
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

<Header {profile} />

<main>
    <Hero {profile} {about} {skills} {badges} />

    <section id="projects" class="section-container snap-section">
        <SectionHeader title="Featured Projects" count={projects.length} animate />
        <div class="projects-grid" style="margin-top: 1.5rem">
            {#each projects as project, i}
                <ProjectCard {project} index={i} />
            {/each}
        </div>
    </section>

    <section id="certifications" class="section-container snap-section">
        <SectionHeader title="Certifications" count={certificates.length} animate />
        {#if certsLoading}
            <div class="certs-loading">
                {#each { length: 3 } as _}
                    <div class="cert-skeleton">
                        <div class="cert-skeleton-media shimmer" />
                        <div class="cert-skeleton-body">
                            <div class="shimmer w-30" />
                            <div class="shimmer w-80" />
                            <div class="shimmer w-50" />
                            <div class="cert-skeleton-tags">
                                <div class="shimmer w-15" />
                                <div class="shimmer w-20" />
                                <div class="shimmer w-18" />
                            </div>
                            <div class="shimmer w-40" />
                        </div>
                    </div>
                {/each}
            </div>
        {:else if certsError}
            <div class="certs-error">
                <span class="certs-error-text">&#9888; Failed to load certificates</span>
            </div>
        {:else if certificates.length > 0}
            <div class="certs-grid">
                {#each certificates as cert, i}
                    <div class="cert-card-wrap">
                        <CertificateCard certificate={cert} index={i} />
                    </div>
                {/each}
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
                        <div class="info-item">
                            <div class="info-icon">
                                <Icon name="chat" size={20} />
                            </div>
                            <div class="info-content">
                                <span class="info-title">Chat</span>
                                <span class="info-text">pranavagarkar8@gmail.com</span>
                            </div>
                        </div>
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
                                    <Icon name={social.icon} size={18} />
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
    :global(.project-card),
    :global(.skill-category-card),
    :global(.skill-card),
    :global(.cert-card-wrap),
    :global(.section-header),
    :global(.about-visual-col),
    :global(.about-content-col),
    :global(.about-spec-row),
    :global(.about-metric),
    :global(.contact-info),
    :global(.contact-form-container),
    :global(.info-item),
    :global(.social-link) {
        will-change: transform, opacity;
    }
</style>
