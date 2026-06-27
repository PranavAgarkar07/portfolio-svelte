<script lang="ts">
    import { onMount } from "svelte";
    import { fade } from "svelte/transition";
    import { base } from "$app/paths";
    import type { Badge } from "$lib/types";

    interface Props {
        profile: {
            name: string;
            role: string;
            status: string;
            avatar: string;
        };
        about: {
            bio: string;
            stats: Array<{ label: string; value: string }>;
        };
        badges?: Badge[];
    }

    let { profile, about, badges = [] }: Props = $props();

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

    let previewBadge: Badge | null = $state(null);

    function badgeImgUrl(url: string): string {
        if (url.startsWith("http")) return url;
        const webp = url.replace(/\.png$/, ".webp");
        return base + webp;
    }

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

            return () => {
                window.removeEventListener('scroll', handleParallax);
                window.removeEventListener('resize', updateRect);
            };
        }
    });
</script>

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
                        <div class="about-badges-category">
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
                                                src={badgeImgUrl(badge.image_url)}
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
    <div class="badge-preview-backdrop" onclick={closePreview} transition:fade={{ duration: 150 }}>
        <div class="badge-preview-card" onclick={(e) => e.stopPropagation()}>
            <button class="badge-preview-close" onclick={closePreview}>×</button>
            <div class="badge-preview-nav">
                <button class="badge-prev" onclick={prevBadge} disabled={currentIndex <= 0} aria-label="Previous badge">‹</button>
                <div class="badge-preview-image">
                    {#key previewBadge.id}
                        <img src={badgeImgUrl(previewBadge.image_url)} alt={previewBadge.name} />
                    {/key}
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
{/if}

<svelte:window onkeydown={(e) => { if (e.key === 'Escape') closePreview(); if (e.key === 'ArrowLeft') prevBadge(); if (e.key === 'ArrowRight') nextBadge(); }} />

<style>
    .about-frame {
        transform-style: preserve-3d;
        will-change: transform;
    }

    /* ─── BADGE PREVIEW ─── */

    .badge-preview-backdrop {
        position: fixed;
        inset: 0;
        z-index: 9999;
        background: rgba(0, 0, 0, 0.8);
        display: flex;
        align-items: center;
        justify-content: center;
        backdrop-filter: blur(4px);
    }

    .badge-preview-card {
        background: #0d1117;
        border: 1px solid rgba(255, 255, 255, 0.08);
        border-radius: 0;
        padding: 2rem;
        max-width: 480px;
        width: 90vw;
        position: relative;
    }

    .badge-preview-close {
        position: absolute;
        top: 8px;
        right: 8px;
        width: 32px;
        height: 32px;
        border-radius: 50%;
        border: none;
        background: rgba(255, 255, 255, 0.06);
        color: rgba(255, 255, 255, 0.5);
        font-size: 1.2rem;
        cursor: pointer;
        display: flex;
        align-items: center;
        justify-content: center;
        transition: background 0.15s ease, color 0.15s ease;
        z-index: 1;
    }

    .badge-preview-close:hover {
        background: rgba(255, 68, 0, 0.15);
        color: var(--accent);
    }

    .badge-preview-nav {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        margin-bottom: 1rem;
    }

    .badge-prev,
    .badge-next {
        flex-shrink: 0;
        width: 40px;
        height: 40px;
        border-radius: 50%;
        border: 1px solid rgba(255, 255, 255, 0.1);
        background: rgba(255, 255, 255, 0.04);
        color: rgba(255, 255, 255, 0.6);
        font-size: 1.4rem;
        cursor: pointer;
        display: flex;
        align-items: center;
        justify-content: center;
        transition: background 0.15s ease, color 0.15s ease, border-color 0.15s ease;
    }

    .badge-prev:hover:not(:disabled),
    .badge-next:hover:not(:disabled) {
        background: rgba(255, 68, 0, 0.1);
        color: var(--accent);
        border-color: rgba(255, 68, 0, 0.2);
    }

    .badge-prev:disabled,
    .badge-next:disabled {
        opacity: 0.25;
        cursor: not-allowed;
    }

    .badge-preview-image {
        flex: 1;
        display: flex;
        align-items: center;
        justify-content: center;
        aspect-ratio: 1;
        background: #070a0f;
        border-radius: 8px;
        overflow: hidden;
    }

    .badge-preview-image img {
        max-width: 100%;
        max-height: 100%;
        object-fit: contain;
        display: block;
    }

    .badge-preview-name {
        font-family: var(--font-heading);
        font-size: 0.85rem;
        font-weight: 600;
        color: var(--text-primary);
        text-transform: uppercase;
        letter-spacing: 0.04em;
        text-align: center;
        margin-bottom: 0.4rem;
    }

    .badge-preview-meta {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 0.5rem;
        margin-bottom: 0.5rem;
    }

    .badge-preview-rarity {
        font-family: var(--font-body);
        font-size: 0.55rem;
        font-weight: 600;
        text-transform: uppercase;
        letter-spacing: 0.1em;
        padding: 2px 8px;
        border-radius: 4px;
        background: rgba(255, 255, 255, 0.04);
        color: rgba(255, 255, 255, 0.4);
        border: 1px solid rgba(255, 255, 255, 0.06);
    }

    .badge-preview-rarity.uncommon {
        color: #60a5fa;
        border-color: rgba(96, 165, 250, 0.2);
        background: rgba(96, 165, 250, 0.06);
    }

    .badge-preview-rarity.rare {
        color: #f59e0b;
        border-color: rgba(245, 158, 11, 0.2);
        background: rgba(245, 158, 11, 0.06);
    }

    .badge-preview-dot {
        width: 4px;
        height: 4px;
        border-radius: 50%;
        background: rgba(255, 255, 255, 0.2);
    }

    .badge-preview-meta a {
        font-family: var(--font-body);
        font-size: 0.55rem;
        font-weight: 600;
        text-transform: uppercase;
        letter-spacing: 0.1em;
        color: var(--accent);
        text-decoration: none;
        transition: opacity 0.15s ease;
    }

    .badge-preview-meta a:hover {
        opacity: 0.7;
    }

    .badge-preview-counter {
        font-family: var(--font-body);
        font-size: 0.5rem;
        color: rgba(255, 255, 255, 0.25);
        text-align: center;
        letter-spacing: 0.15em;
    }

    /* ─── LIGHT MODE ─── */
    :global(body.light-mode) .badge-preview-card {
        background: #fff;
        border-color: rgba(0, 0, 0, 0.08);
    }

    :global(body.light-mode) .badge-preview-close {
        color: rgba(0, 0, 0, 0.4);
        background: rgba(0, 0, 0, 0.04);
    }

    :global(body.light-mode) .badge-preview-close:hover {
        background: rgba(217, 65, 0, 0.08);
        color: var(--light-accent, #d94100);
    }

    :global(body.light-mode) .badge-prev,
    :global(body.light-mode) .badge-next {
        border-color: rgba(0, 0, 0, 0.1);
        background: rgba(0, 0, 0, 0.03);
        color: rgba(0, 0, 0, 0.5);
    }

    :global(body.light-mode) .badge-prev:hover:not(:disabled),
    :global(body.light-mode) .badge-next:hover:not(:disabled) {
        background: rgba(217, 65, 0, 0.06);
        color: var(--light-accent, #d94100);
        border-color: rgba(217, 65, 0, 0.15);
    }

    :global(body.light-mode) .badge-preview-image {
        background: #f4f4f5;
    }

    :global(body.light-mode) .badge-preview-name {
        color: #18181b;
    }

    :global(body.light-mode) .badge-preview-rarity {
        color: rgba(0, 0, 0, 0.4);
        border-color: rgba(0, 0, 0, 0.06);
        background: rgba(0, 0, 0, 0.02);
    }

    :global(body.light-mode) .badge-preview-rarity.uncommon {
        color: #2563eb;
        border-color: rgba(37, 99, 235, 0.15);
        background: rgba(37, 99, 235, 0.04);
    }

    :global(body.light-mode) .badge-preview-rarity.rare {
        color: #d97706;
        border-color: rgba(217, 119, 6, 0.15);
        background: rgba(217, 119, 6, 0.04);
    }
</style>
