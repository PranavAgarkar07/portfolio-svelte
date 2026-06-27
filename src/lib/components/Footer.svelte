<script lang="ts">
    import { onMount } from "svelte";
    import { Icon } from "$lib/components/ui";

    interface Props {
        profile: {
            socials: Array<{ name: string; url: string; label: string }>;
        };
    }

    let { profile }: Props = $props();
    let showBackToTop = $state(false);
    const currentYear = new Date().getFullYear();

    onMount(() => {
        const onScroll = () => {
            showBackToTop = window.scrollY > 400;
        };
        window.addEventListener("scroll", onScroll, { passive: true });
        return () => window.removeEventListener("scroll", onScroll);
    });

    function iconFor(name: string) {
        const map: Record<string, string> = {
            GitHub: "github",
            LinkedIn: "linkedin",
            Email: "envelope",
        };
        return map[name] || "envelope";
    }
</script>

<footer class="aero-footer">
    <div class="footer-content">
        <div class="footer-section brand">
            <h3>Pranav<span class="text-accent">.</span></h3>
            <p>Full-stack developer building secure, high-performance applications with Django, Go, and modern frontends.</p>
            <div class="footer-social-row">
                {#each profile.socials as social}
                    {#if social.label !== "RESUME"}
                        <a
                            href={social.url}
                            target="_blank"
                            aria-label={social.name}
                            rel="noopener"
                            class="footer-social-icon"
                        >
                            <Icon name={iconFor(social.name)} size={18} />
                        </a>
                    {/if}
                {/each}
            </div>
        </div>
        <div class="footer-section links">
            <h4>Navigation</h4>
            <div class="footer-nav-list">
                <a href="#about">About</a>
                <a href="#skills">Skills</a>
                <a href="#certifications">Certifications</a>
                <a href="#projects">Projects</a>
                <a href="#contact">Contact</a>
            </div>
        </div>
        <div class="footer-section info">
            <h4>Contact</h4>
            <p>pranavagarkar8@gmail.com</p>
            <p>Solapur, India</p>
            <p class="footer-status">Available for opportunities</p>
        </div>
    </div>

    <div class="footer-cta">
        <a href="#contact" class="footer-cta-link">
            <span>Get In Touch</span>
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                <line x1="5" y1="12" x2="19" y2="12"></line>
                <polyline points="12 5 19 12 12 19"></polyline>
            </svg>
        </a>
    </div>
    <div class="footer-bottom">
        <p>
            &copy; {currentYear} Pranav Agarkar &mdash; Built with Svelte
        </p>
    </div>
    <button
        class="back-to-top"
        class:visible={showBackToTop}
        onclick={() => window.scrollTo({ top: 0, behavior: 'smooth' })}
        aria-label="Back to top"
        tabindex={showBackToTop ? 0 : -1}
    >
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="18 15 12 9 6 15"></polyline>
        </svg>
    </button>
</footer>

<style>
    .aero-footer {
        background: var(--surface-dark);
        margin-top: 4rem;
        border-top: 2px solid var(--accent);
        padding-top: 3rem;
        padding-bottom: 2rem;
        position: relative;
        z-index: 10;
    }

    .footer-content {
        display: grid;
        grid-template-columns: 1.5fr 1fr 1fr;
        gap: 3rem;
        max-width: 1400px;
        margin: 0 auto;
        padding: 0 2rem;
    }

    .footer-section h3 {
        font-family: 'Clash Display', sans-serif;
        font-size: 1.4rem;
        font-weight: 600;
        margin-bottom: 0.75rem;
    }

    .footer-section h4 {
        font-family: var(--font-heading);
        font-size: 0.75rem;
        font-weight: 700;
        text-transform: uppercase;
        letter-spacing: 0.12em;
        color: var(--accent);
        margin-bottom: 1rem;
    }

    .footer-section p {
        font-size: 0.85rem;
        color: var(--text-secondary);
        line-height: 1.6;
    }

    .footer-section.brand p {
        max-width: 320px;
    }

    .footer-nav-list {
        display: flex;
        flex-direction: column;
        gap: 0.6rem;
    }

    .footer-nav-list a {
        font-size: 0.85rem;
        color: var(--text-secondary);
        text-decoration: none;
        transition: color 0.2s ease;
        width: fit-content;
    }

    .footer-nav-list a:hover {
        color: var(--accent);
    }

    .footer-social-row {
        display: flex;
        gap: 0.5rem;
        margin-top: 1rem;
    }

    .footer-social-icon {
        width: 36px;
        height: 36px;
        display: flex;
        align-items: center;
        justify-content: center;
        border: 1px solid var(--grid-line);
        color: var(--text-secondary);
        transition: all 0.2s ease;
    }

    .footer-social-icon:hover {
        border-color: var(--accent);
        color: var(--accent);
        background: rgba(255, 68, 0, 0.06);
    }

    .footer-section.info p {
        margin-bottom: 0.4rem;
    }

    .footer-status {
        font-size: 0.75rem;
        color: #00cc88;
        margin-top: 0.25rem;
    }

    .footer-bottom {
        text-align: center;
        border-top: 1px solid var(--grid-line);
        padding-top: 2rem;
        margin-top: 3rem;
        font-size: 0.75rem;
        color: var(--text-secondary);
        max-width: 1400px;
        margin-left: auto;
        margin-right: auto;
        padding-left: 2rem;
        padding-right: 2rem;
    }

    .back-to-top {
        position: fixed;
        bottom: 2rem;
        right: 2rem;
        width: 44px;
        height: 44px;
        display: flex;
        align-items: center;
        justify-content: center;
        background: var(--surface-dark);
        border: 1px solid var(--grid-line);
        color: var(--text-secondary);
        cursor: pointer;
        z-index: 999;
        opacity: 0;
        transform: translateY(12px);
        pointer-events: none;
        transition: opacity 0.25s ease, transform 0.25s ease, border-color 0.2s ease, color 0.2s ease;
    }
    .back-to-top.visible {
        opacity: 1;
        transform: translateY(0);
        pointer-events: auto;
    }
    .back-to-top:hover {
        border-color: var(--accent);
        color: var(--accent);
        background: rgba(255,68,0,0.08);
    }
    .back-to-top:focus-visible {
        outline: 2px solid var(--accent);
        outline-offset: 3px;
    }
    :global(body.light-mode) .back-to-top {
        background: var(--surface-dark);
        border-color: var(--wire-color);
    }
    :global(body.light-mode) .back-to-top:hover {
        border-color: var(--accent);
    }

    @media (max-width: 768px) {
        .footer-content {
            grid-template-columns: 1fr;
            gap: 2rem;
            text-align: center;
        }
        .footer-section.brand p {
            max-width: none;
        }
        .footer-social-row {
            justify-content: center;
        }
        .footer-nav-list {
            align-items: center;
        }
    .footer-cta {
        text-align: center;
        padding: 2rem 2rem 0;
        max-width: 1400px;
        margin: 0 auto;
    }
    .footer-cta-link {
        display: inline-flex;
        align-items: center;
        gap: 0.5rem;
        font-family: var(--font-heading);
        font-size: 0.8rem;
        font-weight: 700;
        text-transform: uppercase;
        letter-spacing: 0.12em;
        color: var(--accent);
        text-decoration: none;
        border: 2px solid var(--accent);
        padding: 0.75rem 1.5rem;
        transition: background 0.2s ease, color 0.2s ease;
    }
    .footer-cta-link:hover {
        background: var(--accent);
        color: #050505;
    }
    :global(body.light-mode) .footer-cta-link:hover {
        color: #fff;
    }
    .footer-cta-link svg {
        transition: transform 0.2s ease;
    }
    .footer-cta-link:hover svg {
        transform: translateX(4px);
    }
    .footer-bottom {
            padding-left: 1rem;
            padding-right: 1rem;
        }
        .footer-content {
            padding: 0 1rem;
        }
    }

    @media (prefers-reduced-motion: reduce) {
        .back-to-top {
            transition: none;
        }
    }
</style>
