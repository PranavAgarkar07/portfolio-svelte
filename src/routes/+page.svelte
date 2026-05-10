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
    import { Icon, SectionHeader, Card } from "$lib/components/ui";

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
        if (window.innerWidth <= 768) {
            ScrollTrigger.normalizeScroll(true);
        }
        window.addEventListener("resize", () => ScrollTrigger.refresh());

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

        // Hero animations are now CSS-transition-driven in Hero.svelte
        // (no GSAP hero animations needed — BlurText + heroReady handles reveal)

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
                ease: "elastic.out(1, 0.55)", duration: 0.9,
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
        // 9. SECTION HEADERS — REVEAL WIPE
        // ──────────────────────────────────────────────
        gsap.utils.toArray(".section-header").forEach((header: any) => {
            gsap.set(header, { clipPath: "polygon(0 0, 0 0, 0 100%, 0 100%)" });
            ScrollTrigger.create({
                trigger: header as Element, start: "top 85%", once: true,
                onEnter: () => { gsap.to(header, {
                    clipPath: "polygon(0 0, 100% 0, 100% 100%, 0 100%)",
                    ease: "power3.out", duration: 0.7,
                });}
            });
        });

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
                    ease: "elastic.out(1, 0.55)", duration: 0.8,
                });
                gsap.to(".about-content-col", {
                    x: 0, opacity: 1, rotateY: 0, scale: 1,
                    ease: "elastic.out(1, 0.55)", duration: 0.8, delay: 0.1,
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
                    ease: "elastic.out(1, 0.5)", duration: 0.6,
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
                        val: num, duration: 2, ease: "elastic.out(1, 0.4)",
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
                    ease: "elastic.out(1.2, 0.55)", duration: 0.7,
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
                    ease: "elastic.out(1, 0.65)", duration: 0.8,
                });
                gsap.to(".contact-form-container", {
                    x: 0, opacity: 1,
                    ease: "elastic.out(1, 0.65)", duration: 0.8, delay: 0.15,
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
                    ease: "elastic.out(1.2, 0.5)", duration: 0.6, delay: i * 0.05,
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
                ".info-item", ".social-link",
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
            setTimeout(forceVisibility, 5000);
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
    <Hero {profile} {about} {skills} />

    <section id="projects" class="section-container snap-section">
        <SectionHeader title="Featured Projects" count={projects.length} animate />
        <div class="projects-grid" style="margin-top: 1.5rem">
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
