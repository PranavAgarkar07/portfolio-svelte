<script lang="ts">
    import { onMount } from "svelte";
    import { SectionHeader } from "$lib/components/ui";
    import DevLog from "$lib/components/DevLog.svelte";
    import BlurText from "./BlurText.svelte";

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

    let nameParts = $derived(profile.name.split(" "));

    let heroReady = $state(false);

    onMount(() => {
        requestAnimationFrame(() => { heroReady = true; });
    });
</script>

<section id="hero" class="hero-section snap-section">
    <div class="hero-background-pattern" aria-hidden="true"></div>

    <div class="hero-inner">
        <div class="hero-status" class:hero-ready={heroReady}>
            <span class="hero-status-icon" class:available={profile.status.includes("Open")}>
                <svg width="14" height="14" viewBox="0 0 14 14" fill="none">
                    <circle cx="7" cy="7" r="6" stroke="currentColor" stroke-width="1.5"/>
                    <path d="M4 7h6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
                    <path d="M7 4v6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
                </svg>
            </span>
            <span class="hero-status-value" class:available={profile.status.includes("Open")}>{profile.status}</span>
            <span class="hero-status-divider" aria-hidden="true">//</span>
            <span class="hero-status-role">{profile.role}</span>
        </div>

        <h1 class="hero-name" class:hero-ready={heroReady}>
            <span class="hero-name-line hero-name-line--first">
                <span class="hero-name-text">{nameParts[0]}</span>
            </span>
            <span class="hero-name-line hero-name-line--last">
                <span class="hero-name-text">{nameParts.slice(1).join(" ") || nameParts[0]}</span>
            </span>
        </h1>

        <div class="hero-tagline-wrap" class:hero-ready={heroReady}>
            <BlurText
                text={profile.tagline}
                animateBy="words"
                direction="top"
                delay={120}
                class="hero-tagline"
            />
        </div>

        <div class="cta-group" class:hero-ready={heroReady}>
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

    <div class="hero-scroll" class:hero-ready={heroReady}>
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
    
    /* ─── HERO ─── */
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
        position: relative;
    }
    
    .hero-status {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        margin-bottom: 2rem;
        font-family: var(--font-body);
        font-size: 0.75rem;
        text-transform: uppercase;
        letter-spacing: 0.15em;
        opacity: 0;
        transform: translateY(10px);
        transition: opacity 0.8s ease, transform 0.8s ease;
    }
    
    .hero-status.hero-ready {
        opacity: 1;
        transform: translateY(0);
    }
    
    .hero-status-icon {
        width: 14px; height: 14px; flex-shrink: 0;
        color: var(--text-secondary);
    }
    
    .hero-status-icon.available {
        color: #00ffaa;
        animation: statusPulse 2s ease-in-out infinite;
        filter: drop-shadow(0 0 6px rgba(0, 255, 170, 0.3));
    }
    
    .hero-status-value {
        font-weight: 600;
    }
    
    .hero-status-value.available {
        color: #00ffaa;
    }
    
    .hero-status-divider {
        color: var(--accent);
        opacity: 0.4;
        font-size: 0.65rem;
    }
    
    .hero-status-role {
        color: var(--text-secondary);
        font-weight: 400;
    }
    
    .hero-name {
        margin-bottom: 1.5rem;
    }
    
    .hero-name-line {
        display: block;
        overflow: hidden;
    }
    
    .hero-name-line--last {
        display: flex;
        align-items: center;
        gap: 1.5rem;
        flex-wrap: wrap;
    }
    
    .hero-name-text {
        display: inline-block;
        font-family: var(--font-heading);
        font-size: clamp(3rem, 10vw, 6rem);
        font-weight: 800;
        line-height: 1;
        letter-spacing: -0.04em;
        text-transform: uppercase;
        color: var(--text-primary);
        opacity: 0;
        transform: translateY(40px);
        transition: opacity 0.8s cubic-bezier(0.22, 1, 0.36, 1), transform 0.8s cubic-bezier(0.22, 1, 0.36, 1);
        transition-delay: 0.2s;
    }
    
    .hero-name.hero-ready .hero-name-text {
        opacity: 1;
        transform: translateY(0);
    }

    .hero-name.hero-ready .hero-name-line--last .hero-name-text {
        transition-delay: 0.4s;
    }
    
    .hero-tagline-wrap {
        margin-bottom: 2.5rem;
        opacity: 0;
        transition: opacity 0.6s ease;
        transition-delay: 0.8s;
    }
    
    .hero-tagline-wrap.hero-ready {
        opacity: 1;
    }
    
    .hero-tagline {
        font-family: var(--font-body);
        font-size: clamp(0.9rem, 1.2vw, 1.1rem);
        color: var(--text-secondary);
        max-width: 600px;
        line-height: 1.7;
    }
    
    .cta-group {
        display: flex;
        gap: 1rem;
        opacity: 0;
        transform: translateY(20px);
        transition: opacity 0.6s ease, transform 0.6s ease;
        transition-delay: 1s;
    }
    
    .cta-group.hero-ready {
        opacity: 1;
        transform: translateY(0);
    }
    
    .btn {
        padding: 16px 32px; font-size: 0.9rem; letter-spacing: 0.1em;
    }
    
    .btn-icon { width: 20px; height: 20px; }
    
    .hero-scroll {
        position: absolute;
        bottom: 3rem;
        left: 50%;
        transform: translateX(-50%);
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 0.75rem;
        z-index: 1;
        opacity: 0;
        transition: opacity 0.8s ease;
        transition-delay: 1.2s;
    }
    
    .hero-scroll.hero-ready {
        opacity: 1;
    }
    
    .hero-scroll-pulse {
        width: 2px; height: 24px;
        background: var(--accent);
        animation: scrollPulse 2s ease-in-out infinite;
        margin-top: 4px;
    }
    
    /* ─── KEYFRAMES ─── */
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
    }
    
</style>
