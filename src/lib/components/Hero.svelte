<script lang="ts">
    import { onMount } from "svelte";
    import BlurText from "./BlurText.svelte";
    import { base } from "$app/paths";

    interface Props {
        profile: {
            name: string;
            role: string;
            tagline: string;
            status: string;
            avatar: string;
        };
        skills: Array<{ category: string; items: Array<{ name: string }> }>;
        about: { stats: Array<{ label: string; value: string }> };
    }

    let { profile, skills, about }: Props = $props();

    let marqueeTechList = $derived([...new Set(skills.flatMap(cat => cat.items.map(s => s.name)))]);

    let nameParts = $derived(profile.name.split(" "));

    let heroReady = $state(false);

    onMount(() => {
        // F-P0-2: mark JS so hidden pre-ready states only apply with JS on.
        // Without JS the content below stays default-visible.
        document.documentElement.classList.add("js");
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
        <div class="hero-left-col">
            <h1 class="hero-name" class:hero-ready={heroReady}>
                <span class="hero-name-line hero-name-line--first">
                    <span class="hero-name-text">{nameParts[0]}</span>
                </span>
                <span class="hero-name-line hero-name-line--last">
                    <span class="hero-name-text">{nameParts.slice(1).join(" ") || nameParts[0]}</span>
                </span>
            </h1>

            <div class="hero-availability" class:hero-ready={heroReady}>
                <span class="avail-dot"></span>
                <span class="avail-text">Seeking Internships &amp; SDE Roles · Open to Remote &amp; Relocation (Pune / Bengaluru)</span>
            </div>

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
                <a href="{base}/Pranav_Agarkar_Resume.pdf" target="_blank" rel="noopener" download class="btn btn-resume">
                    <svg class="btn-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                        <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
                        <polyline points="14 2 14 8 20 8"></polyline>
                        <line x1="12" y1="18" x2="12" y2="12"></line>
                        <polyline points="9 15 12 18 15 15"></polyline>
                    </svg>
                    <span class="btn-label">Resume PDF</span>
                </a>
                <a href="#contact" class="btn btn-secondary">
                    <span class="btn-label">Contact Me</span>
                    <svg class="btn-icon" width="18" height="18" viewBox="0 0 18 18" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <polyline points="10 3 16 9 10 15"/>
                    <line x1="2" y1="9" x2="16" y2="9"/>
                    </svg>
                </a>
            </div>

            <div class="tech-marquee" class:hero-ready={heroReady}>
                <span class="marquee-label">TECH STACK</span>
                <div class="marquee-track">
                    <div class="marquee-content">
                        {#each [...marqueeTechList, ...marqueeTechList] as tech, i}
                            <span class="marquee-item" aria-hidden={i >= marqueeTechList.length ? "true" : undefined}>
                                <span class="marquee-dot"></span>
                                {tech}
                            </span>
                        {/each}
                    </div>
                </div>
            </div>
        </div>

        <div class="hero-right-col">
            <div class="hero-avatar-card">
                <img src={profile.avatar} alt={profile.name} width="280" height="280" sizes="(max-width: 600px) 280px, 280px" loading="eager" fetchpriority="high" decoding="async" />
                <div class="hero-stats-mini" aria-label="Quick stats">
                    {#each about.stats.slice(0, 3) as stat}
                        <div class="hero-stat-row">
                            <span class="hero-stat-label">{stat.label}</span>
                            <span class="hero-stat-value">{stat.value}</span>
                        </div>
                    {/each}
                </div>
            </div>
        </div>
    </div>

    <div class="hero-scroll" class:hero-ready={heroReady}>
        <span class="hero-scroll-text">SCROLL</span>
        <div class="hero-scroll-pulse" aria-hidden="true"></div>
    </div>
</section>

<style>
    :global(.hero-ambient-glow) {
        position: absolute;
        top: 0; left: 0;
        width: 100%; height: 100%;
        pointer-events: none;
        z-index: 0;
        background-size: 200% 200%;
        will-change: background-position;
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
        width: 100%;
        max-width: 1000px; padding: 6rem 0;
        position: relative;
    }

    .hero-left-col {
        width: 100%;
        display: flex;
        flex-direction: column;
    }

    .hero-right-col {
        display: none;
    }

    @media (max-width: 768px) {
        .hero-inner {
            padding: 3rem 0;
            display: flex;
            flex-direction: column;
            gap: 2rem;
        }
    }

    @media (min-width: 1100px) {
        .hero-inner {
            max-width: 1300px;
            display: grid;
            grid-template-columns: 1fr 320px;
            align-items: center;
            gap: 4rem;
        }
        .hero-right-col {
            display: flex;
            flex-direction: column;
            align-items: flex-end;
            gap: 1.5rem;
        }
        .hero-avatar-card {
            border: 2px solid var(--grid-line);
            border-left: 4px solid var(--accent);
            padding: 1.5rem;
            display: flex;
            flex-direction: column;
            gap: 1rem;
            width: 100%;
        }
        .hero-avatar-card img {
            width: 100%;
            aspect-ratio: 1;
            object-fit: cover;
            filter: grayscale(20%);
        }
        .hero-stats-mini {
            display: flex;
            flex-direction: column;
            gap: 0.5rem;
            border-top: 1px solid var(--grid-line);
            padding-top: 0.875rem;
            margin-top: 0.125rem;
        }
        .hero-stat-row {
            display: flex;
            justify-content: space-between;
            align-items: baseline;
        }
        .hero-stat-label {
            font-family: var(--font-body);
            font-size: 0.6rem;
            font-weight: 600;
            letter-spacing: 0.12em;
            text-transform: uppercase;
            color: var(--text-secondary);
        }
        .hero-stat-value {
            font-family: var(--font-body);
            font-size: 0.8rem;
            font-weight: 700;
            letter-spacing: 0.04em;
            color: var(--text-primary);
        }
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
        /* F-P0-2: default-visible; hidden pre-ready state only under html.js */
        opacity: 1;
        transform: none;
        transition: opacity 0.8s cubic-bezier(0.22, 1, 0.36, 1), transform 0.8s cubic-bezier(0.22, 1, 0.36, 1);
        transition-delay: 0.2s;
    }
    :global(html.js) .hero-name:not(.hero-ready) .hero-name-text {
        opacity: 0;
        transform: translateY(40px);
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
        /* F-P0-2: default-visible; hidden pre-ready state only under html.js */
        opacity: 1;
        transition: opacity 0.6s ease;
        transition-delay: 0.8s;
    }
    :global(html.js) .hero-tagline-wrap:not(.hero-ready) {
        opacity: 0;
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
        /* F-P0-2: default-visible; hidden pre-ready state only under html.js */
        opacity: 1;
        transform: none;
        transition: opacity 0.6s ease, transform 0.6s ease;
        transition-delay: 1s;
    }
    :global(html.js) .cta-group:not(.hero-ready) {
        opacity: 0;
        transform: translateY(20px);
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
        /* F-P0-2: default-visible; hidden pre-ready state only under html.js */
        opacity: 1;
        transform: none;
        transition: opacity 0.6s ease, transform 0.6s ease;
        transition-delay: 2.5s;
    }
    :global(html.js) .tech-marquee:not(.hero-ready) {
        opacity: 0;
        transform: translateY(10px);
    }
    .tech-marquee.hero-ready {
        opacity: 1;
        transform: translateY(0);
    }

    .marquee-track:hover .marquee-content {
        animation-play-state: paused;
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
        animation-delay: 2.5s;
        animation-fill-mode: backwards;
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
        /* F-P0-2: default-visible; hidden pre-ready state only under html.js */
        opacity: 1;
        transition: opacity 0.8s ease;
        transition-delay: 1.2s;
    }
    :global(html.js) .hero-scroll:not(.hero-ready) {
        opacity: 0;
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

    /* ── RESUME CTA BUTTON ── */
    .btn-resume {
        display: inline-flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.6rem 1.1rem;
        font-family: var(--font-body);
        font-size: 0.75rem;
        font-weight: 700;
        letter-spacing: 0.08em;
        text-transform: uppercase;
        text-decoration: none;
        border: 1.5px solid var(--accent);
        color: var(--accent);
        background: transparent;
        cursor: pointer;
        transition: background 0.2s, color 0.2s;
    }
    .btn-resume:hover {
        background: var(--accent);
        color: #000;
    }

    /* ── AVAILABILITY BADGE ── */
    .hero-availability {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        margin-bottom: 1rem;
        opacity: 1;
        transform: none;
        transition: opacity 0.6s ease, transform 0.6s ease;
        transition-delay: 0.1s;
    }
    :global(html.js) .hero-availability:not(.hero-ready) {
        opacity: 0;
        transform: translateY(10px);
    }
    .hero-availability.hero-ready {
        opacity: 1;
        transform: none;
    }
    .avail-dot {
        display: inline-block;
        width: 8px;
        height: 8px;
        border-radius: 50%;
        background: #22c55e;
        box-shadow: 0 0 6px #22c55e88;
        flex-shrink: 0;
        animation: avail-pulse 2s ease-in-out infinite;
    }
    @keyframes avail-pulse {
        0%, 100% { opacity: 1; box-shadow: 0 0 6px #22c55e88; }
        50% { opacity: 0.7; box-shadow: 0 0 12px #22c55ecc; }
    }
    .avail-text {
        font-family: var(--font-body);
        font-size: 0.7rem;
        font-weight: 600;
        letter-spacing: 0.06em;
        color: var(--text-secondary);
        text-transform: uppercase;
    }
</style>
