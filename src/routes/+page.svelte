<script lang="ts">
    import { onMount } from "svelte";
    import gsap from "gsap";
    import { ScrollTrigger } from "gsap/dist/ScrollTrigger";

    import { theme } from "$lib/stores/theme";
    import { portfolioData } from "$lib/data";
    import Header from "$lib/components/Header.svelte";
    import Hero from "$lib/components/Hero.svelte";
    import ProjectCard from "$lib/components/ProjectCard.svelte";

    import DevLog from "$lib/components/DevLog.svelte";
    import Footer from "$lib/components/Footer.svelte";
    import ContactForm from "$lib/components/ContactForm.svelte";
    import Icon from "$lib/components/Icon.svelte";

    const { profile, about, skills, projects } = portfolioData;

    const siteUrl = "https://pranavagarkar07.github.io/portfolio-svelte/";
    const pageTitle = `${profile.name} | ${profile.role}`;
    const pageDescription = `${profile.tagline} Experienced in building secure, scalable applications with Django, React, Svelte, and Go.`;
    const pageKeywords =
        "Svelte, SvelteKit, Django, Go, Full Stack Developer, Portfolio, Pranav Agarkar, Web Development, TaskVault, BeamSync";
    const coverImage = `${siteUrl}avatar.png`; // Fallback to avatar if no dedicated cover

    onMount(() => {
        const prefersReducedMotion = window.matchMedia("(prefers-reduced-motion: reduce)").matches;
        if (prefersReducedMotion) return;
        gsap.registerPlugin(ScrollTrigger);
        gsap.defaults({ overwrite: "auto" });

        // ──────────────────────────────────────────────
        // 1. READING PROGRESS BAR
        // ──────────────────────────────────────────────
        const progressBar = document.createElement("div");
        Object.assign(progressBar.style, {
            position: "fixed", top: "0", left: "0", width: "0%", height: "2px",
            background: "linear-gradient(90deg, var(--accent), #ff8800, var(--accent))",
            zIndex: "9999", pointerEvents: "none", transition: "width 0.1s linear",
        });
        document.body.prepend(progressBar);
        ScrollTrigger.create({ trigger: "body", start: "top top", end: "bottom bottom", onUpdate: self => {
            progressBar.style.width = self.progress * 100 + "%";
        }});

        // ──────────────────────────────────────────────
        // 2. HERO — SVG CORNER DRAW + MORPH
        // ──────────────────────────────────────────────
        const corners = document.querySelectorAll(".hero-corner svg path");
        corners.forEach((path, i) => {
            const p = path as SVGPathElement;
            const length = p.getTotalLength();
            p.style.strokeDasharray = String(length);
            p.style.strokeDashoffset = String(length);
            gsap.to(p, {
                strokeDashoffset: 0, duration: 2, ease: "elastic.out(1.2, 0.6)", delay: 0.1 + i * 0.08,
            });
        });

        // Corner brackets drift outward on scroll
        gsap.to(".hero-corner--tl", { scrollTrigger: { trigger: "#hero", start: "top top", end: "bottom top", scrub: 1 }, x: -10, y: -10, opacity: 0.3, ease: "power2.out" });
        gsap.to(".hero-corner--tr", { scrollTrigger: { trigger: "#hero", start: "top top", end: "bottom top", scrub: 1 }, x: 10, y: -10, opacity: 0.3, ease: "power2.out" });
        gsap.to(".hero-corner--bl", { scrollTrigger: { trigger: "#hero", start: "top top", end: "bottom top", scrub: 1 }, x: -10, y: 10, opacity: 0.3, ease: "power2.out" });
        gsap.to(".hero-corner--br", { scrollTrigger: { trigger: "#hero", start: "top top", end: "bottom top", scrub: 1 }, x: 10, y: 10, opacity: 0.3, ease: "power2.out" });

        // ──────────────────────────────────────────────
        // 3. HERO — CLIP-PATH TEXT REVEAL
        // ──────────────────────────────────────────────
        gsap.set(".hero-name", { clipPath: "polygon(0 0, 0 0, 0 100%, 0 100%)" });
        gsap.set(".hero-role-block", { clipPath: "polygon(0 0, 0 0, 0 100%, 0 100%)" });
        gsap.set(".hero-tagline", { clipPath: "polygon(0 0, 0 0, 0 100%, 0 100%)" });

        const heroTimeline = gsap.timeline();
        heroTimeline
            .to(".hero-status", { y: 0, opacity: 1, duration: 0.6, ease: "power3.out" })
            .to(".hero-name", { clipPath: "polygon(0 0, 100% 0, 100% 100%, 0 100%)", duration: 1.2, ease: "power4.out" }, "-=0.3")
            .to(".hero-role-block", { clipPath: "polygon(0 0, 100% 0, 100% 100%, 0 100%)", duration: 0.8, ease: "power3.out" }, "-=0.4")
            .to(".hero-tagline", { clipPath: "polygon(0 0, 100% 0, 100% 100%, 0 100%)", duration: 0.8, ease: "power3.out" }, "-=0.4")
            .from(".cta-group .btn", { y: 20, opacity: 0, stagger: 0.15, duration: 0.7, ease: "elastic.out(1, 0.7)", clearProps: "all" }, "-=0.3")
            .to(".hero-scroll", { opacity: 1, duration: 0.8, ease: "power2.out" }, "-=0.3");

        // ──────────────────────────────────────────────
        // 4. HERO — ENTRANCE GLOW SWEEP
        // ──────────────────────────────────────────────
        const glow = document.createElement("div");
        Object.assign(glow.style, {
            position: "absolute", top: "0", left: "-100%", width: "100%", height: "100%",
            background: "linear-gradient(90deg, transparent, rgba(255,68,0,0.04), transparent)",
            pointerEvents: "none", zIndex: "0",
        });
        document.querySelector(".hero-inner")?.appendChild(glow);
        gsap.to(glow, { left: "200%", duration: 2.5, delay: 1, ease: "power2.inOut" });

        // ──────────────────────────────────────────────
        // 5. HERO — PARALLAX ON SCROLL
        // ──────────────────────────────────────────────
        gsap.to(".hero-inner", {
            scrollTrigger: { trigger: "#hero", start: "top top", end: "bottom top", scrub: 1.5 },
            y: 80, scale: 0.97, opacity: 0.5, ease: "power2.out",
        });
        gsap.to(".hero-background-pattern", {
            scrollTrigger: { trigger: "#hero", start: "top top", end: "bottom top", scrub: 1 },
            opacity: 0.01, ease: "power1.out",
        });

        // ──────────────────────────────────────────────
        // 6. PROJECTS — CARD PERSPECTIVE + STAGGER REVEAL
        // ──────────────────────────────────────────────
        gsap.utils.toArray(".project-card").forEach((card: any, i: number) => {
            const dir = i % 2 === 0 ? -1 : 1;
            gsap.set(card, { x: dir * 150, opacity: 0, rotateY: dir * 15, scale: 0.9 });
            gsap.to(card, {
                scrollTrigger: { trigger: card, start: "top 85%", end: "top 50%", scrub: 1.2 },
                x: 0, opacity: 1, rotateY: 0, scale: 1,
                ease: "elastic.out(1, 0.65)",
            });
        });

        // ──────────────────────────────────────────────
        // 7. SKILLS — CASCADE WAVE REVEAL
        // ──────────────────────────────────────────────
        gsap.set(".skill-category-card", { y: 80, opacity: 0, scale: 0.92, rotateX: 10 });
        gsap.to(".skill-category-card", {
            scrollTrigger: { trigger: "#skills", start: "top 85%", end: "top 50%", scrub: 1.5 },
            y: 0, opacity: 1, scale: 1, rotateX: 0,
            stagger: { each: 0.15, from: "edges", ease: "power2.inOut" },
            ease: "elastic.out(1, 0.55)",
        });

        // ──────────────────────────────────────────────
        // 8. SKILL ITEMS — INDIVIDUAL STAGGER POP
        // ──────────────────────────────────────────────
        gsap.utils.toArray(".skill-card").forEach((card: any, i: number) => {
            gsap.set(card, { y: 30, opacity: 0, scale: 0.85 });
            gsap.to(card, {
                scrollTrigger: { trigger: card, start: "top 92%", end: "top 75%", scrub: 1 },
                y: 0, opacity: 1, scale: 1,
                ease: "back.out(2)",
            });
        });

        // ──────────────────────────────────────────────
        // 9. SECTION HEADERS — REVEAL WIPE
        // ──────────────────────────────────────────────
        gsap.utils.toArray(".section-header").forEach((header: any) => {
            gsap.set(header, { clipPath: "polygon(0 0, 0 0, 0 100%, 0 100%)" });
            gsap.to(header, {
                scrollTrigger: { trigger: header, start: "top 85%", end: "top 70%", scrub: 1 },
                clipPath: "polygon(0 0, 100% 0, 100% 100%, 0 100%)",
                ease: "power3.out",
            });
        });

        // ──────────────────────────────────────────────
        // 10. ABOUT — DEPTH LAYER REVEAL
        // ──────────────────────────────────────────────
        gsap.set(".about-visual-col", { x: -80, opacity: 0, rotateY: 8, scale: 0.95 });
        gsap.to(".about-visual-col", {
            scrollTrigger: { trigger: "#about", start: "top 80%", end: "top 45%", scrub: 1.3 },
            x: 0, opacity: 1, rotateY: 0, scale: 1,
            ease: "elastic.out(1, 0.55)",
        });
        gsap.set(".about-content-col", { x: 80, opacity: 0, rotateY: -8, scale: 0.95 });
        gsap.to(".about-content-col", {
            scrollTrigger: { trigger: "#about", start: "top 80%", end: "top 45%", scrub: 1.3 },
            x: 0, opacity: 1, rotateY: 0, scale: 1,
            ease: "elastic.out(1, 0.55)",
        });

        // About specs — sweep up
        gsap.utils.toArray(".about-spec-row").forEach((row: any, i: number) => {
            gsap.set(row, { y: 25, opacity: 0 });
            gsap.to(row, {
                scrollTrigger: { trigger: "#about", start: "top 80%", end: "top 60%", scrub: 1 },
                y: 0, opacity: 1,
                ease: "elastic.out(1, 0.5)",
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
            let obj = { val: 0 };
            gsap.to(obj, {
                val: num, duration: 2, ease: "elastic.out(1, 0.4)",
                scrollTrigger: { trigger: el, start: "top 85%", end: "top 70%", scrub: 1.2 },
                onUpdate: () => { display.textContent = Math.round(obj.val) + suffix; },
            });
        });

        // About metrics — spring up
        gsap.utils.toArray(".about-metric").forEach((metric: any, i: number) => {
            gsap.set(metric, { y: 40, opacity: 0, scale: 0.85 });
            gsap.to(metric, {
                scrollTrigger: { trigger: "#about", start: "top 75%", end: "top 50%", scrub: 1.3 },
                y: 0, opacity: 1, scale: 1,
                ease: "elastic.out(1.2, 0.55)",
            });
        });

        // ──────────────────────────────────────────────
        // 12. CONTACT — DUAL-SIDED REVEAL
        // ──────────────────────────────────────────────
        gsap.set(".contact-info", { x: -50, opacity: 0 });
        gsap.to(".contact-info", {
            scrollTrigger: { trigger: "#contact", start: "top 85%", end: "top 60%", scrub: 1.2 },
            x: 0, opacity: 1,
            ease: "elastic.out(1, 0.65)",
        });
        gsap.set(".contact-form-container", { x: 50, opacity: 0 });
        gsap.to(".contact-form-container", {
            scrollTrigger: { trigger: "#contact", start: "top 85%", end: "top 60%", scrub: 1.2 },
            x: 0, opacity: 1,
            ease: "elastic.out(1, 0.65)",
        });

        // Contact info items — staggered
        gsap.utils.toArray(".info-item").forEach((item: any, i: number) => {
            gsap.set(item, { y: 20, opacity: 0 });
            gsap.to(item, {
                scrollTrigger: { trigger: "#contact", start: "top 85%", end: "top 70%", scrub: 1 },
                y: 0, opacity: 1,
                ease: "back.out(1.7)",
            });
        });

        // ──────────────────────────────────────────────
        // 13. SOCIAL LINKS — STAGGER REVEAL
        // ──────────────────────────────────────────────
        gsap.utils.toArray(".social-link").forEach((link: any, i: number) => {
            gsap.set(link, { y: 15, opacity: 0, scale: 0.8 });
            gsap.to(link, {
                scrollTrigger: { trigger: "#contact", start: "top 85%", end: "top 75%", scrub: 1 },
                y: 0, opacity: 1, scale: 1,
                ease: "elastic.out(1.2, 0.5)",
            });
        });
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
    <Hero {profile} {about} {skills} />

    <section id="projects" class="section-container snap-section">
        <div class="section-header fade-in">
            <h2>
                Featured Projects <span class="count-badge"
                    >{projects.length}</span
                >
            </h2>
        </div>
        <div class="projects-grid">
            {#each projects as project, i}
                <ProjectCard {project} index={i} />
            {/each}
        </div>
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
                    <div class="info-card">
                        <div class="info-item">
                            <div class="info-icon">
                                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
                                    <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
                                </svg>
                            </div>
                            <div class="info-content">
                                <span class="info-title">Chat</span>
                                <span class="info-text">pranavagarkar8@gmail.com</span>
                            </div>
                        </div>
                        <div class="info-item">
                            <div class="info-icon">
                                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
                                    <circle cx="12" cy="12" r="10"></circle>
                                    <polyline points="12 6 12 12 16 14"></polyline>
                                </svg>
                            </div>
                            <div class="info-content">
                                <span class="info-title">Response Time</span>
                                <span class="info-text">Within 24 hours</span>
                            </div>
                        </div>
                        <div class="info-item">
                            <div class="info-icon available">
                                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
                                    <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
                                    <polyline points="22 4 12 14.01 9 11.01"></polyline>
                                </svg>
                            </div>
                            <div class="info-content">
                                <span class="info-title">Status</span>
                                <span class="info-text available-text">Available for new projects</span>
                            </div>
                        </div>
                    </div>

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
