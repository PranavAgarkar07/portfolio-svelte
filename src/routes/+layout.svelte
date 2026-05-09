<script lang="ts">
    import { theme } from "$lib/stores/theme";
    import { onMount } from "svelte";
    import "../styles.css";
    import "../light-mode.css";

    let { children } = $props();

    onMount(() => {
        const prefersReducedMotion = window.matchMedia(
            "(prefers-reduced-motion: reduce)",
        ).matches;

        const handleMove = (e: MouseEvent) => {
            if (prefersReducedMotion) return;

            const cards = document.querySelectorAll(
                ".aero-card, .btn, .terminal-wrapper, .nav-links a, .theme-btn, .mobile-menu-btn",
            ) as NodeListOf<HTMLElement>;
            cards.forEach((card) => {
                const rect = card.getBoundingClientRect();
                const x = e.clientX - rect.left;
                const y = e.clientY - rect.top;
                card.style.setProperty("--x", `${x}px`);
                card.style.setProperty("--y", `${y}px`);
            });
        };

        if (!prefersReducedMotion) {
            window.addEventListener("mousemove", handleMove);
        }

        // Sync body class on mount (hydration fix)
        const unsubscribe = theme.subscribe((value) => {
            if (value === "light") {
                document.body.classList.add("light-mode");
            } else {
                document.body.classList.remove("light-mode");
            }
        });

        return () => {
            window.removeEventListener("mousemove", handleMove);
            unsubscribe();
        };
    });
</script>

<svelte:head>
    <!-- Primary Meta -->
    <title>Pranav Agarkar | Fullstack Developer</title>
    <meta name="description" content="Portfolio of Pranav Agarkar — Fullstack Developer specializing in Django, React, Svelte, and Go. Building fast, modern web apps." />
    <meta name="keywords" content="Pranav Agarkar, Fullstack Developer, Django, React, Svelte, Go, SvelteKit, portfolio, web developer, India" />
    <meta name="author" content="Pranav Agarkar" />
    <meta name="robots" content="index, follow" />
    <meta name="theme-color" content="#030405" />
    <link rel="canonical" href="https://pranavagarkar07.github.io/portfolio-svelte/" />

    <!-- Open Graph (Facebook, LinkedIn, WhatsApp previews) -->
    <meta property="og:type" content="website" />
    <meta property="og:url" content="https://pranavagarkar07.github.io/portfolio-svelte/" />
    <meta property="og:title" content="Pranav Agarkar | Fullstack Developer" />
    <meta property="og:description" content="Fullstack Developer specializing in Django, React, Svelte, and Go. Explore my projects, skills, and contact info." />
    <meta property="og:image" content="https://pranavagarkar07.github.io/portfolio-svelte/og-image.png" />
    <meta property="og:image:width" content="1200" />
    <meta property="og:image:height" content="630" />
    <meta property="og:site_name" content="Pranav Agarkar Portfolio" />
    <meta property="og:locale" content="en_US" />

    <!-- Twitter / X Card -->
    <meta name="twitter:card" content="summary_large_image" />
    <meta name="twitter:url" content="https://pranavagarkar07.github.io/portfolio-svelte/" />
    <meta name="twitter:title" content="Pranav Agarkar | Fullstack Developer" />
    <meta name="twitter:description" content="Fullstack Developer specializing in Django, React, Svelte, and Go. Explore my projects, skills, and contact info." />
    <meta name="twitter:image" content="https://pranavagarkar07.github.io/portfolio-svelte/og-image.png" />

    <!-- JSON-LD Structured Data (Person Schema) -->
    {@html `<script type="application/ld+json">${JSON.stringify({
        "@context": "https://schema.org",
        "@type": "Person",
        "name": "Pranav Agarkar",
        "url": "https://pranavagarkar07.github.io/portfolio-svelte/",
        "jobTitle": "Fullstack Developer",
        "description": "Fullstack Developer specializing in Django, React, Svelte, and Go.",
        "knowsAbout": ["Django", "React", "Svelte", "SvelteKit", "Go", "Python", "JavaScript", "TypeScript", "PostgreSQL"],
        "sameAs": [
            "https://github.com/PranavAgarkar07"
        ]
    })}</script>`}

    <!-- Preconnects -->
    <link rel="preconnect" href="https://fonts.googleapis.com" />
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous" />
    <link rel="preconnect" href="https://cdn.jsdelivr.net" crossorigin="anonymous" />
    <link rel="preconnect" href="https://cdnjs.cloudflare.com" crossorigin="anonymous" />
    <link rel="preconnect" href="https://api.fontshare.com" crossorigin="anonymous" />

    <!-- Fonts -->
    <link
        href="https://fonts.googleapis.com/css2?family=Space+Grotesk:wght@300;400;500;600;700&display=swap"
        rel="stylesheet"
    />

    <!-- Icon CDNs -->
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/gh/devicons/devicon@latest/devicon.min.css" />

    <!-- Prevent FOIT on icon fonts -->
    <style>
        @font-face {
            font-display: swap;
        }
    </style>
</svelte:head>


<div class="noise-overlay"></div>
<div class="background-gradient"></div>

<svelte:body class:light-mode={$theme === "light"} />

{@render children()}
