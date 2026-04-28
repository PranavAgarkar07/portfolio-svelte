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
    <title>Pranav Agarkar | Fullstack Developer</title>
    <meta
        name="description"
        content="Portfolio of Pranav Agarkar - Fullstack Developer specializing in Django, React, and Svelte"
    />
    <link rel="preconnect" href="https://fonts.googleapis.com" />
    <link
        rel="preconnect"
        href="https://fonts.gstatic.com"
        crossorigin="anonymous"
    />
    <link
        rel="preconnect"
        href="https://cdn.jsdelivr.net"
        crossorigin="anonymous"
    />
    <link
        rel="preconnect"
        href="https://cdnjs.cloudflare.com"
        crossorigin="anonymous"
    />
    <link
        rel="stylesheet"
        href="https://cdn.jsdelivr.net/gh/devicons/devicon@latest/devicon.min.css"
    />
    <link
        rel="stylesheet"
        href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css"
    />
    <link
        href="https://fonts.googleapis.com/css2?family=Space+Grotesk:wght@300;400;500;600;700&family=Syncopate:wght@400;700&display=swap"
        rel="stylesheet"
    />
    <!-- Override icon font-display to prevent invisible text (FOIT) during CDN load -->
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
