<script lang="ts">
    import { onMount } from "svelte";
    import { fade } from "svelte/transition";
    import { SectionHeader } from "$lib/components/ui";
    import DevLog from "$lib/components/DevLog.svelte";
    import BlurText from "./BlurText.svelte";
    import type { Badge } from "$lib/types";

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
            items: Array<{ name: string; icon: string; level: string }>;
        }>;
        badges?: Badge[];
    }

    let { profile, about, skills, badges = [] }: Props = $props();

    let nameParts = $derived(profile.name.split(" "));

    let sortedBadges = $derived([...badges].sort((a, b) => a.display_order - b.display_order));

    let importantBadges = $derived(sortedBadges.filter(b => b.important));
    let nonImportantBadges = $derived(sortedBadges.filter(b => !b.important));

    let showAllBadges = $state(false);

    let displayBadges = $derived(showAllBadges ? sortedBadges : importantBadges);

    let badgeGroups = $derived.by(() => {
        const groups: Record<string, Badge[]> = {};
        for (const b of displayBadges) {
            const cat = b.category || "Achievements";
            if (!groups[cat]) groups[cat] = [];
            groups[cat].push(b);
        }
        return groups;
    });

    let heroReady = $state(false);

    let previewBadge: Badge | null = $state(null);

    function togglePreview(badge: Badge) {
        previewBadge = previewBadge?.id === badge.id ? null : badge;
    }

    function closePreview() {
        previewBadge = null;
    }

    let allDisplayBadges = $derived.by(() => {
        const flat: Badge[] = [];
        for (const [, group] of Object.entries(badgeGroups)) {
            flat.push(...group);
        }
        return flat;
    });

    let currentIndex = $derived(previewBadge ? allDisplayBadges.findIndex(b => b.id === previewBadge.id) : -1);

    function prevBadge() {
        if (currentIndex > 0) previewBadge = allDisplayBadges[currentIndex - 1];
    }

    function nextBadge() {
        if (currentIndex < allDisplayBadges.length - 1) previewBadge = allDisplayBadges[currentIndex + 1];
    }

    onMount(() => {
        const prefersReducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches;
        requestAnimationFrame(() => { heroReady = true; });

        // Ambient cursor glow
        let glowCleanup: (() => void) | null = null;
        const heroSection = document.querySelector('.hero-section');
        if (heroSection && !prefersReducedMotion) {
            const glow = document.createElement('div');
            glow.className = 'hero-ambient-glow';
            glow.style.background = `radial-gradient(600px circle at 50% 50%, rgba(255,68,0,0.08), transparent 60%)`;
            heroSection.appendChild(glow);

            let rafId: number;
            const handleMove = (e: MouseEvent) => {
                cancelAnimationFrame(rafId);
                rafId = requestAnimationFrame(() => {
                    const rect = heroSection.getBoundingClientRect();
                    const x = ((e.clientX - rect.left) / rect.width) * 100;
                    const y = ((e.clientY - rect.top) / rect.height) * 100;
                    glow.style.backgroundPosition = `${x}% ${y}%`;
                });
            };

            window.addEventListener('mousemove', handleMove, { passive: true });
            glowCleanup = () => {
                window.removeEventListener('mousemove', handleMove);
                glow.remove();
            };
        }

        // Parallax on about avatar — pre-calculate rect, use passive scroll
        const aboutFrame = document.querySelector('.about-frame');
        if (aboutFrame && !prefersReducedMotion) {
            const frameEl = aboutFrame as HTMLElement;
            let rect = frameEl.getBoundingClientRect();
            const updateRect = () => { rect = frameEl.getBoundingClientRect(); };
            window.addEventListener('resize', updateRect, { passive: true });

            const handleParallax = () => {
                const center = rect.top + rect.height / 2;
                const viewportCenter = window.innerHeight / 2;
                const offset = (center - viewportCenter) / viewportCenter;
                frameEl.style.transform = `translateY(${offset * -8}px)`;
            };
            window.addEventListener('scroll', handleParallax, { passive: true });

            const parallaxCleanup = () => {
                window.removeEventListener('scroll', handleParallax);
                window.removeEventListener('resize', updateRect);
            };

            // Store for cleanup
            const existingCleanup = glowCleanup;
            glowCleanup = () => {
                existingCleanup?.();
                parallaxCleanup();
            };
        }

        // Pulse CTA after idle
        let pulseTimer: ReturnType<typeof setTimeout>;
        function resetPulseTimer() {
            clearTimeout(pulseTimer);
            const btn = document.querySelector('.cta-group .btn-primary');
            if (btn) btn.classList.remove('idle-pulse');
            pulseTimer = setTimeout(() => {
                const btn = document.querySelector('.cta-group .btn-primary');
                if (btn) btn.classList.add('idle-pulse');
            }, 5000);
        }
        resetPulseTimer();
        const events = ['scroll', 'mousemove', 'keydown'];
        events.forEach(ev =>
            window.addEventListener(ev, resetPulseTimer, { passive: true })
        );

        return () => {
            clearTimeout(pulseTimer);
            events.forEach(ev => window.removeEventListener(ev, resetPulseTimer));
            glowCleanup?.();
        };
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

        <div class="tech-marquee" class:hero-ready={heroReady}>
            <span class="marquee-label">TECH STACK</span>
            <div class="marquee-track">
                <div class="marquee-content">
                    {#each ['Python', 'Django', 'React', 'Go', 'Svelte', 'TypeScript', 'PostgreSQL', 'Docker', 'Linux'] as tech}
                        <span class="marquee-item">
                            <span class="marquee-dot"></span>
                            {tech}
                        </span>
                    {/each}
                </div>
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
            {#if badges.length > 0}
                {#each Object.entries(badgeGroups) as [category, groupBadges]}
                    <div class="about-badges">
                        <div class="about-badges-header">
                            <span class="about-badges-label">{category.toUpperCase()}</span>
                            <span class="about-badges-count">{groupBadges.length}</span>
                        </div>
                        <div class="about-badges-grid">
                            {#each groupBadges as badge, i}
                                <button
                                    class="badge-cell"
                                    onclick={() => togglePreview(badge)}
                                    style="--i: {i}"
                                >
                                    <span class="badge-cell-frame">
                                        <img
                                            src={badge.image_url}
                                            alt={badge.name}
                                            loading="lazy"
                                            decoding="async"
                                        />
                                        <span class="badge-cell-dot" class:uncommon={badge.rarity === "uncommon"} class:rare={badge.rarity === "rare"}></span>
                                    </span>
                                    <span class="badge-cell-name">{badge.name}</span>
                                </button>
                            {/each}
                        </div>
                    </div>
                {/each}
                {#if nonImportantBadges.length > 0 && !showAllBadges}
                    <button class="about-badges-more" onclick={() => showAllBadges = true}>
                        <span class="more-text">SHOW MORE BADGES</span>
                        <span class="more-count">+{nonImportantBadges.length}</span>
                    </button>
                {/if}
                {#if showAllBadges && nonImportantBadges.length > 0}
                    <button class="about-badges-more about-badges-less" onclick={() => showAllBadges = false}>
                        <span class="more-text">SHOW LESS</span>
                    </button>
                {/if}
            {/if}
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

{#if previewBadge}
    {#key previewBadge.id}
        <!-- svelte-ignore a11y_no_static_element_interactions -->
        <div class="badge-preview-backdrop" onclick={closePreview} transition:fade={{ duration: 150 }}>
            <!-- svelte-ignore a11y_no_static_element_interactions -->
            <div class="badge-preview-card" onclick={(e) => e.stopPropagation()}>
                <button class="badge-preview-close" onclick={closePreview}>×</button>
                <div class="badge-preview-nav">
                    <button class="badge-prev" onclick={prevBadge} disabled={currentIndex <= 0} aria-label="Previous badge">‹</button>
                    <div class="badge-preview-image">
                        <img src={previewBadge.image_url} alt={previewBadge.name} />
                    </div>
                    <button class="badge-next" onclick={nextBadge} disabled={currentIndex >= allDisplayBadges.length - 1} aria-label="Next badge">›</button>
                </div>
                <div class="badge-preview-name">{previewBadge.name}</div>
                <div class="badge-preview-meta">
                    <span class="badge-preview-rarity" class:uncommon={previewBadge.rarity === 'uncommon'} class:rare={previewBadge.rarity === 'rare'}>
                        {previewBadge.rarity}
                    </span>
                    {#if previewBadge.credential_url}
                        <span class="badge-preview-dot"></span>
                    {/if}
                    {#if previewBadge.credential_url}
                        <a href={previewBadge.credential_url} target="_blank" rel="noopener noreferrer">credential</a>
                    {/if}
                </div>
                <div class="badge-preview-counter">{currentIndex + 1} / {allDisplayBadges.length}</div>
            </div>
        </div>
    {/key}
{/if}

<svelte:window onkeydown={(e) => { if (e.key === 'Escape') closePreview(); if (e.key === 'ArrowLeft') prevBadge(); if (e.key === 'ArrowRight') nextBadge(); }} />

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
                            <span class="skill-level-label">{skill.level}</span>
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
    .about-frame {
        transform-style: preserve-3d;
        will-change: transform;
    }
    
    :global(.hero-ambient-glow) {
        position: absolute;
        top: 0; left: 0;
        width: 100%; height: 100%;
        pointer-events: none;
        z-index: 0;
        background-size: 200% 200%;
        will-change: background-position;
    }
    
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
    
    .tech-marquee {
        display: flex;
        align-items: center;
        gap: 1rem;
        margin-top: 3rem;
        width: 100%;
        max-width: 600px;
        opacity: 0;
        transform: translateY(10px);
        transition: opacity 0.6s ease, transform 0.6s ease;
        transition-delay: 1.4s;
    }
    .tech-marquee.hero-ready {
        opacity: 1;
        transform: translateY(0);
    }
    .marquee-label {
        font-family: var(--font-body);
        font-size: 0.6rem;
        letter-spacing: 0.2em;
        color: var(--text-secondary);
        white-space: nowrap;
        flex-shrink: 0;
        opacity: 0.5;
    }
    .marquee-track {
        flex: 1;
        overflow: hidden;
        mask-image: linear-gradient(to right, transparent 0%, black 5%, black 95%, transparent 100%);
        -webkit-mask-image: linear-gradient(to right, transparent 0%, black 5%, black 95%, transparent 100%);
    }
    .marquee-content {
        display: flex;
        gap: 0.5rem;
        animation: marqueeScroll 20s linear infinite;
        width: max-content;
    }
    .marquee-item {
        display: flex;
        align-items: center;
        gap: 0.35rem;
        font-family: var(--font-body);
        font-size: 0.75rem;
        color: var(--text-secondary);
        white-space: nowrap;
        text-transform: uppercase;
        letter-spacing: 0.08em;
    }
    .marquee-dot {
        width: 4px;
        height: 4px;
        border-radius: 50%;
        background: var(--accent);
        flex-shrink: 0;
    }
    @keyframes marqueeScroll {
        0% { transform: translateX(0); }
        100% { transform: translateX(-50%); }
    }
    @media (prefers-reduced-motion: reduce) {
        .marquee-content { animation: none; }
    }
    
    .btn {
        padding: 16px 32px; font-size: 0.9rem; letter-spacing: 0.1em;
    }
    
    .btn-icon { width: 20px; height: 20px; }
    
    .hero-scroll {
        position: absolute;
        bottom: 2.5rem;
        right: 2rem;
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 0.5rem;
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
    
    @keyframes idlePulse {
        0%, 100% { box-shadow: 4px 4px 0px var(--accent); }
        50% { box-shadow: 4px 4px 20px var(--accent-glow), 4px 4px 0px var(--accent); }
    }
    :global(.btn-primary.idle-pulse) {
        animation: idlePulse 2s ease-in-out 3;
    }
    
    /* ─── LIGHT MODE ─── */
    :global(body.light-mode) .hero-status-icon.available {
        color: #008f5e;
        filter: drop-shadow(0 0 4px rgba(0, 143, 94, 0.25));
    }
    :global(body.light-mode) .hero-status-value.available {
        color: #008f5e;
    }
    
</style>
