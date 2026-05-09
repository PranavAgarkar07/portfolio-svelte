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
        const prefersReducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches;
        if (prefersReducedMotion) return; // Skip animations if user prefers reduced motion

        gsap.registerPlugin(ScrollTrigger);

        // Hero Timeline
        const tl = gsap.timeline();
        tl.from(".hero-content", {
            y: 50,
            opacity: 0,
            duration: 1,
            ease: "power3.out",
        }).from(
            ".cta-group .btn",
            {
                y: 20,
                opacity: 0,
                stagger: 0.2,
                duration: 0.8,
                ease: "power2.out",
                clearProps: "all", // Ensure props are cleared after animation
            },
            "-=0.5",
        );

        // Featured Projects Stagger
        gsap.utils.toArray(".project-card").forEach((card: any, i: number) => {
            gsap.from(card, {
                scrollTrigger: {
                    trigger: card,
                    start: "top 85%",
                },
                x: i === 0 ? -120 : 120,
                opacity: 0,
                duration: 1,
                ease: "power3.out",
            });
        });

        // Skills Pop-in
        gsap.from(".skill-category-card", {
            scrollTrigger: {
                trigger: "#skills",
                start: "top 85%",
            },
            y: 30,
            opacity: 0,
            stagger: 0.1,
            duration: 0.6,
            ease: "power2.out",
            clearProps: "all", // Ensure no residual transforms affect layout
        });

        // Section Headers Glitch-in (Simple Slide)
        gsap.utils.toArray(".section-header").forEach((header: any) => {
            gsap.from(header, {
                scrollTrigger: {
                    trigger: header,
                    start: "top 85%",
                },
                x: -50,
                opacity: 0,
                duration: 0.8,
                ease: "power2.out",
            });
        });

        // Contact Section
        gsap.from(".contact-card", {
            scrollTrigger: {
                trigger: "#contact",
                start: "top 80%",
            },
            y: 50,
            duration: 1,
            ease: "power2.out",
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
            <h2>Featured Projects <span class="count-badge">{projects.length}</span></h2>
        </div>
        <div class="projects-grid">
            {#each projects as project, i}
                <ProjectCard {project} index={i} />
            {/each}
        </div>
    </section>

    <section id="contact" class="section-container snap-section">
        <div class="aero-card contact-card fade-in">
            <div class="contact-header">
                <h2>Connect</h2>
                <p>
                    Send me a message below or reach out directly on social platforms.
                </p>
            </div>

            <ContactForm />

            <div class="contact-socials">
                {#each profile.socials as social}
                    <a
                        href={social.url}
                        target={social.url.startsWith("mailto") ? undefined : "_blank"}
                        rel="noopener"
                        class="contact-social-link"
                    >
                        <Icon name={social.icon} size={14} />
                        <span>{social.label}</span>
                    </a>
                {/each}
            </div>
        </div>
    </section>
</main>

<Footer {profile} />
