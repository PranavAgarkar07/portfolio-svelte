<script lang="ts">
    import { base } from "$app/paths";
    import { theme } from "$lib/stores/theme";
    import { scrollY } from "$lib/stores/scroll";
    import { onMount } from "svelte";
    import { initAnalytics, trackEvent } from "$lib/analytics";
    import gsap from "gsap";
    import { ScrollTrigger } from "gsap/dist/ScrollTrigger";
    import Lenis from "lenis";
    import { Lightbox } from "lightbox3";
    import "lightbox3/style.css";
    import "../styles.css";
    import "../light-mode.css";

    let { children } = $props();

    onMount(() => {
        initAnalytics();
        trackEvent("pageview", window.location.pathname);

        function applyTheme(value: 'dark' | 'light') {
            document.getElementById('theme-guard')?.remove();
            const style = document.createElement('style');
            style.id = 'theme-guard';
            style.textContent = '*,*::before,*::after{transition-duration:0s!important;transition-property:none!important;animation-duration:0s!important}';
            document.head.appendChild(style);
            document.body.offsetHeight;
            if (value === 'light') {
                document.body.classList.add('light-mode');
            } else {
                document.body.classList.remove('light-mode');
            }
            requestAnimationFrame(() => {
                document.getElementById('theme-guard')?.remove();
            });
        }
        const unsub = theme.subscribe(applyTheme);
        const prefersReducedMotion = window.matchMedia(
            "(prefers-reduced-motion: reduce)",
        ).matches;

        if ("serviceWorker" in navigator) {
            navigator.serviceWorker.register(`${base}/service-worker.js`);
        }

        if ("web-vital" in navigator === false && "PerformanceObserver" in window) {
            try {
                new PerformanceObserver((list) => {
                    for (const entry of list.getEntries()) {
                        if (entry.entryType === "largest-contentful-paint") {
                            console.log("LCP:", entry.startTime);
                        }
                        if (entry.entryType === "layout-shift" && !entry.hadRecentInput) {
                            console.log("CLS:", entry.value);
                        }
                        if (entry.entryType === "first-input") {
                            console.log("FID:", entry.processingStart - entry.startTime);
                        }
                    }
                }).observe({ type: "largest-contentful-paint", buffered: true });
                new PerformanceObserver((list) => {
                    let cls = 0;
                    for (const entry of list.getEntries()) {
                        if (!entry.hadRecentInput) cls += entry.value;
                    }
                    console.log("CLS total:", cls);
                }).observe({ type: "layout-shift", buffered: true });
            } catch {}
        }

        Lightbox.init({
            springOpen: { stiffness: 500, damping: 40, mass: 0.8 },
            springClose: { stiffness: 400, damping: 35, mass: 0.8 },
        });

        let lenis: Lenis | null = null;

        const isTouchDevice = "ontouchstart" in window || navigator.maxTouchPoints > 0;

        if (!prefersReducedMotion && !isTouchDevice) {
            gsap.registerPlugin(ScrollTrigger);

            lenis = new Lenis({
                duration: 1.2,
                easing: (t) => Math.min(1, 1.001 - Math.pow(2, -10 * t)),
                orientation: "vertical",
                gestureOrientation: "vertical",
                smoothWheel: true,
            });

            lenis.on("scroll", ScrollTrigger.update);

            gsap.ticker.add((time) => {
                lenis!.raf(time * 1000);
                scrollY.set(lenis!.scroll);
            });
            gsap.ticker.lagSmoothing(33, 16);
        }

        function getDocumentTop(el: HTMLElement): number {
            let top = 0;
            let current: HTMLElement | null = el;
            while (current) {
                top += current.offsetTop;
                current = current.offsetParent as HTMLElement | null;
            }
            return top;
        }

        function handleAnchorClick(e: MouseEvent) {
            if (prefersReducedMotion) return;
            const target = e.currentTarget as HTMLAnchorElement;
            const href = target.getAttribute("href");
            if (!href || !href.startsWith("#")) return;
            const id = href.slice(1);
            const el = document.getElementById(id);
            if (!el) return;
            e.preventDefault();
            if (lenis) {
                const targetY = getDocumentTop(el);
                lenis.scrollTo(targetY, { duration: 1.2 });
            } else {
                el.scrollIntoView({ behavior: "smooth" });
            }
        }

        document.querySelectorAll('a[href^="#"]').forEach((a) => {
            a.addEventListener("click", handleAnchorClick);
        });

        const cardListeners: Array<{ el: Element; handler: (e: Event) => void }> = [];
        const cards = document.querySelectorAll(
            ".aero-card, .btn, .terminal-wrapper, .nav-links a, .mobile-menu-btn",
        );
        cards.forEach((card) => {
            const handler = (e: Event) => {
                if (prefersReducedMotion) return;
                const rect = (e.currentTarget as HTMLElement).getBoundingClientRect();
                (card as HTMLElement).style.setProperty("--x", `${rect.width / 2}px`);
                (card as HTMLElement).style.setProperty("--y", `${rect.height / 2}px`);
            };
            card.addEventListener("mouseenter", handler);
            cardListeners.push({ el: card, handler });
        });

        return () => {
            unsub();
            document.querySelectorAll('a[href^="#"]').forEach((a) => {
                a.removeEventListener("click", handleAnchorClick);
            });
            cardListeners.forEach(({ el, handler }) => {
                el.removeEventListener("mouseenter", handler);
            });
            lenis?.destroy();
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
    <link rel="preconnect" href="https://api.fontshare.com" crossorigin="anonymous" />
    <link rel="preconnect" href="https://syci1ayb7l.execute-api.ap-south-1.amazonaws.com" crossorigin="anonymous" />
    <link rel="preconnect" href="https://portfolio-uploads-sentinel.s3.ap-south-1.amazonaws.com" crossorigin="anonymous" />

    <!-- Fonts -->
    <link
        href="https://fonts.googleapis.com/css2?family=Space+Grotesk:wght@300;400;500;600;700&display=swap"
        rel="stylesheet"
    />

    <!-- Icon CDNs -->
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/gh/devicons/devicon@latest/devicon.min.css" />

    <!-- Prevent FOIT on icon fonts -->
    <style>
        @font-face { font-family: "devicon"; font-display: swap; }
    </style>
    <style>
        .lightbox3-overlay { --lb-image-border-radius: 0px; }
    </style>
</svelte:head>


<div class="background-gradient"></div>

<svelte:body />

{@render children()}
