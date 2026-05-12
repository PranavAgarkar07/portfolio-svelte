<script lang="ts">
    import { Icon } from "$lib/components/ui";

    interface Props {
        profile: {
            socials: Array<{ name: string; url: string; label: string }>;
        };
    }

    let { profile }: Props = $props();
    const currentYear = new Date().getFullYear();
</script>

<footer class="aero-footer">
    <div class="scroll-signal">
        <div class="signal-track-footer"></div>
    </div>

    <div class="footer-content">
        <div class="footer-section brand">
            <h3>Pranav.</h3>
            <p>Crafting digital experiences with precision and passion.</p>
        </div>
        <div class="footer-section links">
            <h4>Navigation</h4>
            <a href="#about">About</a><br />
            <a href="#skills">Skills</a><br />
            <a href="#projects">Projects</a><br />
        </div>
        <div class="footer-section social">
            <h4>Connect</h4>
            <div class="social-icons">
                {#each profile.socials as social}
                    {#if social.label !== "RESUME"}
                        <a
                            href={social.url}
                            target="_blank"
                            aria-label={social.name}
                            rel="noopener"
                        >
                            {#if social.name === "GitHub"}
                                <Icon name="github" size={20} />
                            {:else if social.name === "LinkedIn"}
                                <Icon name="linkedin" size={20} />
                            {:else}
                                <Icon name="envelope" size={20} />
                            {/if}
                        </a>
                    {/if}
                {/each}
            </div>
        </div>
    </div>

    <div class="footer-bottom">
        <p>
            &copy; {currentYear} Pranav Agarkar. All Rights Reserved.
        </p>
    </div>
    <button
        class="back-to-top"
        onclick={() => window.scrollTo({ top: 0, behavior: 'smooth' })}
        aria-label="Back to top"
    >
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="18 15 12 9 6 15"></polyline>
        </svg>
    </button>
</footer>

<style>
    /* Add specific styles for the signal if not in global CSS */
    .scroll-signal {
        position: relative;
        display: flex;
        align-items: center;
        gap: 1rem;
        margin-bottom: 3rem;
        opacity: 0.7;
    }

    .signal-text {
        font-family: var(--font-heading);
        font-size: 0.8rem;
        letter-spacing: 2px;
        color: var(--accent);
        white-space: nowrap;
    }

    .signal-track-footer {
        height: 1px;
        flex-grow: 1;
        background: linear-gradient(90deg, var(--accent) 0%, transparent 100%);
    }

    /* Override or ensure footer styles match */
    .aero-footer {
        margin-top: 4rem;
        border-top: 1px solid var(--grid-line);
        padding-top: 3rem;
        padding-bottom: 2rem;
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
        border: 1px solid rgba(255,255,255,0.1);
        color: var(--text-secondary);
        cursor: pointer;
        z-index: 999;
        transition: all 0.2s ease;
        clip-path: polygon(8px 0, 100% 0, 100% calc(100% - 8px), calc(100% - 8px) 100%, 0 100%, 0 8px);
    }
    .back-to-top:hover {
        border-color: var(--accent);
        color: var(--accent);
        background: rgba(255,68,0,0.08);
    }
    :global(body.light-mode) .back-to-top {
        background: #f0f0f3;
        border-color: rgba(0,0,0,0.12);
    }
    :global(body.light-mode) .back-to-top:hover {
        border-color: var(--accent);
    }
    @media (prefers-reduced-motion: reduce) {
        .back-to-top { display: none; }
    }
</style>
