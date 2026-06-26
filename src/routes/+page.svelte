<script lang="ts">
    import { onMount } from "svelte";
    import gsap from "gsap";
    import { ScrollTrigger } from "gsap/dist/ScrollTrigger";

    import { theme } from "$lib/stores/theme";
    import { portfolioData, fallbackCertificates } from "$lib/data";
    import Header from "$lib/components/Header.svelte";
    import Hero from "$lib/components/Hero.svelte";
    import ProjectCard from "$lib/components/ProjectCard.svelte";
    import CertificateCard from "$lib/components/CertificateCard.svelte";

    import DevLog from "$lib/components/DevLog.svelte";
    import Footer from "$lib/components/Footer.svelte";
    import ContactForm from "$lib/components/ContactForm.svelte";
    import { Icon, SectionHeader, Card, Skeleton } from "$lib/components/ui";
    import type { Certificate, Badge } from "$lib/types";
    import { fetchProjectLikes, toggleProjectLike } from "$lib/analytics";

    const { profile, about, skills, projects } = portfolioData;
    const API_BASE = (import.meta.env.VITE_API_URL || "").replace(/\/$/, "");

    let certificates = $state<Certificate[]>([]);
    let certsLoading = $state(true);
    let certsError = $state("");
    let badges = $state<Badge[]>([]);
    let likes = $state<Record<string, number>>({});
    let userLikes = $state<Record<string, boolean>>({});
    let prefersReducedMotion = $state(false);

    async function fetchCerts() {
        if (!API_BASE) {
            certificates = fallbackCertificates;
            certsLoading = false;
            return;
        }
        try {
            const r = await fetch(`${API_BASE}/api/certificates`);
            if (!r.ok) throw new Error("HTTP " + r.status);
            const data = await r.json();
            certificates = Array.isArray(data.certificates) ? data.certificates : (Array.isArray(data) ? data : []);
            certsLoading = false;
        } catch (e) {
            certsError = String(e);
            certificates = fallbackCertificates;
            certsLoading = false;
        }
    }

    const siteUrl = "https://pranavagarkar07.github.io/portfolio-svelte/";
    const pageTitle = `${profile.name} | ${profile.role}`;
    const pageDescription = `${profile.tagline} Experienced in building secure, scalable applications with Django, React, Svelte, and Go.`;
    const pageKeywords =
        "Svelte, SvelteKit, Django, Go, Full Stack Developer, Portfolio, Pranav Agarkar, Web Development, TaskVault, BeamSync";
    const coverImage = `${siteUrl}og-image.png`;

    onMount(() => {
        prefersReducedMotion = window.matchMedia("(prefers-reduced-motion: reduce)").matches;
        gsap.registerPlugin(ScrollTrigger);
        gsap.defaults({ overwrite: "auto" });

        fetchCerts();

        if (API_BASE) {
            fetch(`${API_BASE}/api/badges`)
                .then(r => r.ok ? r.json() : Promise.reject("HTTP " + r.status))
                .then(data => { badges = data.badges ?? data ?? []; })
                .catch(() => {});
            fetchProjectLikes().then(result => {
                likes = result.likes;
                userLikes = result.user_likes;
            }).catch(() => {});
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
        // 6. PROJECTS — STAGGER REVEAL
        // ──────────────────────────────────────────────
        gsap.set(".project-card", {
            y: 40,
            opacity: 0,
        });
        ScrollTrigger.create({
            trigger: "#projects", start: "top 85%", once: true,
            onEnter: () => { gsap.to(".project-card", {
                y: 0, opacity: 1,
                stagger: 0.08,
                ease: "power3.out", duration: 0.6,
            });}
        });

        // ──────────────────────────────────────────────
        // 7. SKILLS — CASCADE REVEAL
        // ──────────────────────────────────────────────
        gsap.set(".skill-category-card", { y: 40, opacity: 0 });
        ScrollTrigger.create({
            trigger: "#skills", start: "top 85%", once: true,
            onEnter: () => { gsap.to(".skill-category-card", {
                y: 0, opacity: 1,
                stagger: 0.1,
                ease: "power3.out", duration: 0.6,
            });}
        });

        // ──────────────────────────────────────────────
        // 8. SKILL ITEMS — BATCH STAGGER
        // ──────────────────────────────────────────────
        const skillCards = gsap.utils.toArray(".skill-card");
        if (skillCards.length > 0) {
            gsap.set(skillCards, { y: 20, opacity: 0 });
            ScrollTrigger.create({
                trigger: "#skills", start: "top 75%", once: true,
                onEnter: () => { gsap.to(skillCards, {
                    y: 0, opacity: 1,
                    stagger: 0.03,
                    ease: "power3.out", duration: 0.4,
                });}
            });
        }

        // Handled by motion.div whileInView in SectionHeader.svelte
        // (GSAP version removed to avoid duplicate clip-path animation)

        // ──────────────────────────────────────────────
        gsap.set(".about-visual-col", { x: -30, opacity: 0 });
        gsap.set(".about-content-col", { x: 30, opacity: 0 });
        ScrollTrigger.create({
            trigger: "#about", start: "top 80%", once: true,
            onEnter: () => {
                gsap.to(".about-visual-col", {
                    x: 0, opacity: 1,
                    ease: "power3.out", duration: 0.6,
                });
                gsap.to(".about-content-col", {
                    x: 0, opacity: 1,
                    ease: "power3.out", duration: 0.6, delay: 0.08,
                });
            }
        });

        // About specs — sweep up (batched: inside #about trigger)
        gsap.set(".about-spec-row", { y: 25, opacity: 0 });

        // ──────────────────────────────────────────────
        // 11. ABOUT — METRICS COUNT-UP
        // ──────────────────────────────────────────────
        gsap.set(".about-metric", { y: 20, opacity: 0 });
        ScrollTrigger.create({
            trigger: "#about", start: "top 75%", once: true,
            onEnter: () => {
                // Spec rows
                gsap.utils.toArray(".about-spec-row").forEach((row: any) => {
                    gsap.to(row, { y: 0, opacity: 1, ease: "power3.out", duration: 0.5 });
                });
                // Count-up metrics
                gsap.utils.toArray<HTMLElement>(".about-metric-value").forEach((el) => {
                    const raw = el.textContent || "0";
                    const num = parseFloat(raw.replace(/[^0-9.]/g, ""));
                    if (isNaN(num)) return;
                    const suffix = raw.replace(/[0-9.]/g, "");
                    let obj = { val: 0 };
                    gsap.to(obj, {
                        val: num, duration: 1.5, ease: "power3.out",
                        onUpdate: () => { el.textContent = Math.round(obj.val) + suffix; },
                    });
                });
                // Metrics fade up
                gsap.to(".about-metric", {
                    y: 0, opacity: 1,
                    stagger: 0.08, ease: "power3.out", duration: 0.5,
                });
            }
        });

        // ──────────────────────────────────────────────
        // 12-13. CONTACT — ALL (batched: 1 trigger instead of 7)
        // ──────────────────────────────────────────────
        gsap.set(".contact-info", { y: 15, opacity: 0 });
        gsap.set(".contact-form-container", { y: 15, opacity: 0 });
        gsap.set(".info-item", { y: 10, opacity: 0 });
        gsap.set(".social-link", { y: 10, opacity: 0 });

        // Contact section single trigger
        ScrollTrigger.create({
            trigger: "#contact", start: "top 85%", once: true,
            onEnter: () => {
                // Simple fade up
                gsap.to(".contact-info", { y: 0, opacity: 1, ease: "power3.out", duration: 0.5 });
                gsap.to(".contact-form-container", {
                    y: 0, opacity: 1, ease: "power3.out", duration: 0.5, delay: 0.1,
                });
                gsap.to(".info-item", {
                    y: 0, opacity: 1, ease: "power3.out", duration: 0.4,
                    stagger: 0.05,
                });
                gsap.to(".social-link", {
                    y: 0, opacity: 1, ease: "power3.out", duration: 0.4,
                    stagger: 0.04,
                });
            }
        });

        // ──────────────────────────────────────────────
        // 14. SAFETY — recover stuck hidden elements (mobile + reduced motion)
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
                gsap.set(el, {
                    opacity: 1, y: 0, x: 0, scale: 1,
                    clipPath: "polygon(0% 0%, 100% 0%, 100% 100%, 0% 100%)"
                });
            });
            gsap.set(".hero-name-text", {
                clipPath: "polygon(0% 0%, 100% 0%, 100% 100%, 0% 100%)"
            });
        }

        const isMobile = window.innerWidth <= 768 || "ontouchstart" in window;
        if (isMobile || prefersReducedMotion) {
            setTimeout(forceVisibility, 500);
            if (isMobile) {
                window.addEventListener("touchstart", () => ScrollTrigger.refresh(), { once: true });
            }
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

<main id="main-content">
    <Hero {profile} {about} {skills} {badges} />

    <section id="projects" class="section-container snap-section">
        <SectionHeader title="Featured Projects" count={projects.length} animate />
        <div class="projects-grid" style="margin-top: 1.5rem">
            {#each projects as project, i}
                <ProjectCard {project} index={i} liked={userLikes[project.name] ?? false} likeCount={likes[project.name] ?? 0} onLike={(liked) => {
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
</style>
