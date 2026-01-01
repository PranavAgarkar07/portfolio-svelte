<script lang="ts">
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
        <span class="signal-text">// END_OF_LINE</span>
        <div class="signal-track-footer"></div>
    </div>

    <div class="footer-content">
        <div class="footer-section brand">
            <h3>Pranav.</h3>
            <p>Crafting digital experiences with precision and passion.</p>
        </div>
        <div class="footer-section links">
            <h4>Navigation</h4>
            <ul>
                <li><a href="#about">About</a></li>
                <li><a href="#skills">Skills</a></li>
                <li><a href="#projects">Projects</a></li>
            </ul>
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
                            <i
                                class={social.name === "GitHub"
                                    ? "devicon-github-original"
                                    : social.name === "LinkedIn"
                                      ? "devicon-linkedin-plain"
                                      : "fas fa-envelope"}
                            ></i>
                        </a>
                    {/if}
                {/each}
            </div>
        </div>
    </div>

    <div class="footer-bottom">
        <p>
            &copy; {currentYear} Pranav Agarkar. All Rights Reserved. // System Operational
        </p>
    </div>
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
        background: rgba(0, 0, 0, 0.2);
    }
</style>
