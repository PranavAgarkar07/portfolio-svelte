<script lang="ts">
    import { onMount } from "svelte";
    import { SectionHeader } from "$lib/components/ui";
    import DevLog from "$lib/components/DevLog.svelte";

    interface Props {
        profile: {
            name: string;
            role: string;
            tagline: string;
            status: string;
            avatar: string;
        };
        about: {
            bio: string;
            stats: Array<{ label: string; value: string }>;
        };
        skills: Array<{
            category: string;
            items: Array<{ name: string; icon: string; level: number }>;
        }>;
    }

    let { profile, about, skills }: Props = $props();

    let displayName = $state("");
    let displayRole = $state("");

    const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789@#$%";

    function decryptText(target: string, setter: (val: string) => void) {
        let iteration = 0;
        const interval = setInterval(() => {
            setter(
                target
                    .split("")
                    .map((char, i) => {
                        if (i < iteration) return char;
                        return chars[Math.floor(Math.random() * chars.length)];
                    })
                    .join(""),
            );

            iteration += 1 / 3;
            if (iteration >= target.length) clearInterval(interval);
        }, 30);
    }

    onMount(() => {
        decryptText(profile.name, (v) => (displayName = v));
        setTimeout(
            () => decryptText(profile.role, (v) => (displayRole = v)),
            300,
        );
    });
</script>

<section id="hero" class="hero-section snap-section">
    <div class="hero-corner hero-corner--tl" aria-hidden="true">
        <svg width="48" height="48" viewBox="0 0 48 48" fill="none">
            <path d="M0 48V0h48" stroke="var(--accent)" stroke-width="2" opacity="0.7"/>
        </svg>
    </div>
    <div class="hero-corner hero-corner--tr" aria-hidden="true">
        <svg width="48" height="48" viewBox="0 0 48 48" fill="none">
            <path d="M48 48V0H0" stroke="var(--accent)" stroke-width="2" opacity="0.7"/>
        </svg>
    </div>
    <div class="hero-corner hero-corner--bl" aria-hidden="true">
        <svg width="48" height="48" viewBox="0 0 48 48" fill="none">
            <path d="M0 0v48h48" stroke="var(--accent)" stroke-width="2" opacity="0.7"/>
        </svg>
    </div>
    <div class="hero-corner hero-corner--br" aria-hidden="true">
        <svg width="48" height="48" viewBox="0 0 48 48" fill="none">
            <path d="M48 0v48H0" stroke="var(--accent)" stroke-width="2" opacity="0.7"/>
        </svg>
    </div>

    <!-- Subtle ASCII background pattern -->
    <div class="hero-background-pattern" aria-hidden="true"></div>

    <div class="hero-inner">
        <div class="hero-status">
            <span class="hero-status-icon" class:available={profile.status.includes("Open")}>
                <svg width="14" height="14" viewBox="0 0 14 14" fill="none">
                    <circle cx="7" cy="7" r="6" stroke="currentColor" stroke-width="1.5"/>
                    <path d="M4 7h6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
                    <path d="M7 4v6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
                </svg>
            </span>
            <span class="hero-status-value" class:available={profile.status.includes("Open")}>{profile.status}</span>
        </div>

        <h1 class="hero-name" id="hero-name">{displayName || profile.name}</h1>

        <div class="hero-role-block">
            <h2 class="hero-role" id="hero-role">{displayRole || profile.role}</h2>
        </div>

        <p class="hero-tagline">
            {profile.tagline}
            <span class="hero-tagline-cursor" aria-hidden="true">│</span>
        </p>

        <div class="cta-group">
            <a href="#projects" class="btn btn-primary">
                <span class="btn-label">View Work</span>
                <svg class="btn-icon" width="18" height="18" viewBox="0 0 18 18" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M2 9l6 6 6-6"/>
                </svg>
            </a>
            <a href="#contact" class="btn btn-secondary">
                <span class="btn-label">Contact Me</span>
                <svg class="btn-icon" width="18" height="18" viewBox="0 0 18 18" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <circle cx="9" cy="9" r="7" stroke="currentColor" stroke-width="2"/>
                </svg>
            </a>
        </div>
    </div>

    <div class="hero-scroll">
        <span class="hero-scroll-text">SCROLL</span>
        <div class="hero-scroll-pulse" aria-hidden="true"></div>
    </div>
</section>

<section id="about" class="section-container snap-section-content">
    <div class="about-grid">
        <div class="about-visual-col">
            <div class="about-frame">
                <div class="about-frame-border"></div>
                <img
                    src={profile.avatar}
                    alt="{profile.name}"
                    class="about-img"
                    loading="lazy"
                    decoding="async"
                    width="400"
                    height="400"
                />
                <div class="about-frame-grid"></div>
            </div>
            <div class="about-specs">
                <div class="about-spec-row">
                    <span class="about-spec-key">Role</span>
                    <span class="about-spec-val">{profile.role}</span>
                </div>
                <div class="about-spec-divider"></div>
                <div class="about-spec-row">
                    <span class="about-spec-key">Status</span>
                    <span class="about-spec-val about-spec-val--accent">{profile.status}</span>
                </div>
            </div>
        </div>
        <div class="about-content-col">
            <div class="about-header">
                <h2>About</h2>
            </div>
            <p class="about-bio">{about.bio}</p>
            <div class="about-metrics">
                {#each about.stats as stat}
                    <div class="about-metric">
                        <span class="about-metric-value">{stat.value}</span>
                        <span class="about-metric-label">{stat.label}</span>
                    </div>
                {/each}
            </div>
        </div>
    </div>
</section>

<section id="skills" class="section-container snap-section-content">
    <SectionHeader title="Technical Proficiency" class="fade-in" />
    <div class="skills-wrapper">
        {#each skills as category}
            <div class="aero-card skill-category-card fade-in">
                <h3 class="category-title">{category.category}</h3>
                <div class="skills-grid">
                    {#each category.items as skill}
                        <div class="skill-card">
                            <i class="{skill.icon} skill-icon"></i>
                            <span class="skill-name">{skill.name}</span>
                            <div class="skill-level-bar">
                                <div
                                    class="fill"
                                    style="width: {skill.level}%"
                                ></div>
                            </div>
                        </div>
                    {/each}
                </div>
            </div>
        {/each}
    </div>

    <!-- Live Activity Metric -->
    <div class="skills-footer fade-in">
        <DevLog />
    </div>
</section>

<style>
    .skills-footer {
        margin-top: 3rem;
        display: flex;
        justify-content: center;
        width: 100%;
    }
    
    /* ─── HERO ENHANCEMENTS ─── */
    .hero-background-pattern {
        position: absolute; top: 0; left: 0;
        width: 100%; height: 100%;
        pointer-events: none; z-index: 0;
        opacity: 0.03;
        background-image: radial-gradient(circle at 20px 20px, rgba(255,255,255,0.03) 1px, transparent 1px);
        background-size: 40px 40px;
    }
    
    .hero-inner {
        max-width: 1000px; padding: 6rem 0;
        position: relative; /* for glow sweep */
    }
    
    .hero-status {
        opacity: 0;
    }
    
    .hero-status-icon {
        width: 16px; height: 16px; flex-shrink: 0;
        color: var(--text-secondary);
    }
    
    .hero-status-icon.available {
        color: #00ffaa;
        animation: statusPulse 2s ease-in-out infinite;
        filter: drop-shadow(0 0 6px rgba(0, 255, 170, 0.3));
    }
    
    .hero-status-value {
        font-family: var(--font-heading);
        font-weight: 600; letter-spacing: 0.05em;
    }
    
    .hero-status-value.available {
        color: #00ffaa;
        text-shadow: 0 0 8px rgba(0, 255, 170, 0.3);
    }
    
    .hero-name {
        position: relative;
        margin-left: -12px;
    }
    
    .hero-role {
        font-weight: 500; letter-spacing: 0.1em;
    }
    
    .hero-tagline {
        position: relative; display: inline-block;
    }
    
    .hero-tagline-cursor {
        display: inline-block; width: 1px;
        animation: blinkCursor 1s steps(2, start) infinite;
        color: var(--accent); margin-left: 2px; font-weight: bold;
    }
    
    .hero-scroll {
        opacity: 0;
    }
    
    .hero-scroll-pulse {
        width: 2px; height: 24px;
        background: var(--accent);
        animation: scrollPulse 2s ease-in-out infinite;
        margin-top: 4px;
    }
    
    .btn {
        padding: 16px 32px; font-size: 0.9rem; letter-spacing: 0.1em;
    }
    
    .btn-icon { width: 20px; height: 20px; }
    
    /* ─── KEYFRAMES ─── */
    @keyframes blinkCursor { 50% { opacity: 0; } }
    
    @keyframes statusPulse {
        0%, 100% { opacity: 1; filter: drop-shadow(0 0 6px rgba(0, 255, 170, 0.3)); }
        50% { opacity: 0.7; filter: drop-shadow(0 0 2px rgba(0, 255, 170, 0.2)); }
    }
    
    @keyframes scrollPulse {
        0%, 100% { opacity: 1; transform: scaleY(1); }
        50% { opacity: 0.5; transform: scaleY(0.6); }
    }
    
    /* ─── LIGHT MODE ─── */
    :global(body.light-mode) .hero-status-icon.available {
        color: #008855;
        filter: drop-shadow(0 0 4px rgba(0, 136, 85, 0.3));
    }
    :global(body.light-mode) .hero-status-value.available {
        color: #008855;
        text-shadow: 0 0 8px rgba(0, 136, 85, 0.3);
    }
    :global(body.light-mode) .hero-tagline-cursor { color: #008855; }
    :global(body.light-mode) .hero-scroll-pulse { background: #008855; }
</style>
