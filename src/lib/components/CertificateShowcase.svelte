<script lang="ts">
    import type { Certificate } from "$lib/types";
    import CertificateCard from "$lib/components/CertificateCard.svelte";

    let { certificates }: { certificates: Certificate[] } = $props();

    /* ── Ordering (fixed) ────────────────────────────────────── */

    /* Folder-structure order: 01_Internships → 02_Certifications →
       03_Trainings → 04_Events, then document order inside each.
       The build pipeline (scripts/build-content.js) already emits
       certs in this order, so no runtime sorting is applied. */

    /* ── Category filter (mirrors the folder structure) ──────── */

    type Category = { key: string; label: string; rank: number };

    const categories: Category[] = $derived.by(() => {
        const map = new Map<string, Category>();
        for (const c of certificates) {
            if (!map.has(c.category)) {
                map.set(c.category, {
                    key: c.category,
                    label: c.category_label,
                    rank: c.category_rank,
                });
            }
        }
        return [...map.values()].sort((a, b) => a.rank - b.rank);
    });

    /* ── Featured grid + explicit "View more" ────────────────
       The data source is already ordered by editorial priority
       (internships → certifications → trainings → events), so
       the first five are the strongest credentials. The sixth
       grid cell becomes a "View more" tile that reveals the
       remaining certificates on demand. */

    const FEATURED_COUNT = 5;
    let expanded = $state(false);
    let revealKey = $state(0);

    let activeCategory = $state<string>("all");

    function setCategory(key: string) {
        if (activeCategory === key) return;
        activeCategory = key;
        expanded = false;
        revealKey += 1; /* re-trigger the stagger reveal */
    }

    const shown = $derived(
        activeCategory === "all"
            ? certificates
            : certificates.filter((c) => c.category === activeCategory)
    );

    const visible = $derived(
        expanded ? shown : shown.slice(0, FEATURED_COUNT)
    );
    const hiddenCount = $derived(shown.length - visible.length);
</script>

<div class="cs-shell">
    {#if certificates.length > 0}
        <!-- ── Toolbar: category filter ── -->
        <div class="cs-toolbar">
            <div
                class="cs-group"
                role="group"
                aria-label="Filter certificates by category"
            >
                <span class="cs-group-label" aria-hidden="true">Category</span>
                <button
                    type="button"
                    class="cs-pill"
                    aria-pressed={activeCategory === "all"}
                    aria-label="Show all categories"
                    onclick={() => setCategory("all")}
                >
                    All
                </button>
                {#each categories as cat (cat.key)}
                    <button
                        type="button"
                        class="cs-pill"
                        aria-pressed={activeCategory === cat.key}
                        aria-label="Show only {cat.label}"
                        onclick={() => setCategory(cat.key)}
                    >
                        {cat.label}
                    </button>
                {/each}
            </div>

            <span class="cs-meta" aria-live="polite">
                {#if hiddenCount > 0}
                    Showing {visible.length} of {shown.length}
                {:else}
                    {shown.length} certificate{shown.length !== 1 ? "s" : ""}{activeCategory !== "all" ? " · " + (categories.find((c) => c.key === activeCategory)?.label || "") : ""}
                {/if}
            </span>
        </div>

        <!-- ── Featured grid + View more tile ── -->
        {#key revealKey}
            <div class="cs-grid">
                {#each visible as cert, i (cert.id)}
                    <div class="cs-item" style="--idx: {i + 1}">
                        <CertificateCard certificate={cert} index={i} />
                    </div>
                {/each}

                {#if !expanded && hiddenCount > 0}
                    <button
                        type="button"
                        class="cs-more"
                        onclick={() => (expanded = true)}
                        aria-label="View all {shown.length} certificates"
                    >
                        <span class="cs-more-count" aria-hidden="true"
                            >+{hiddenCount}</span
                        >
                        <span class="cs-more-title">View more</span>
                        <span class="cs-more-sub">
                            Show all {shown.length}
                        </span>
                        <span class="cs-more-icon" aria-hidden="true">
                            <svg
                                width="18"
                                height="18"
                                viewBox="0 0 24 24"
                                fill="none"
                                stroke="currentColor"
                                stroke-width="2.5"
                                stroke-linecap="round"
                                stroke-linejoin="round"
                            >
                                <path d="M12 5v14" />
                                <path d="m19 12-7 7-7-7" />
                            </svg>
                        </span>
                    </button>
                {/if}
            </div>
        {/key}
    {:else}
        <div class="cs-empty">No certificates yet.</div>
    {/if}
</div>

<style>
    /* ═══════════════════════ Shell ═══════════════════════ */

    .cs-shell {
        margin-top: 1.5rem;
        width: 100%;
    }

    /* ═══════════════════════ Sort toolbar ═══════════════════════ */

    .cs-toolbar {
        display: flex;
        align-items: center;
        justify-content: space-between;
        flex-wrap: wrap;
        gap: 0.6rem;
        margin-bottom: 1.2rem;
        padding-bottom: 0.9rem;
        border-bottom: 1px solid rgba(255, 255, 255, 0.05);
    }

    .cs-group {
        display: inline-flex;
        align-items: center;
        flex-wrap: wrap;
        gap: 0.4rem;
    }

    .cs-group-label {
        font-family: var(--font-heading);
        font-size: 0.55rem;
        font-weight: 700;
        text-transform: uppercase;
        letter-spacing: 0.18em;
        color: var(--text-muted, rgba(255, 255, 255, 0.35));
        margin-right: 0.25rem;
    }

    .cs-pill {
        appearance: none;
        cursor: pointer;
        font-family: var(--font-heading);
        font-size: 0.6rem;
        font-weight: 700;
        text-transform: uppercase;
        letter-spacing: 0.08em;
        color: rgba(255, 255, 255, 0.45);
        background: transparent;
        border: 1px solid rgba(255, 255, 255, 0.08);
        padding: 0.42rem 0.85rem;
        transition:
            color 0.25s ease,
            border-color 0.25s ease,
            background 0.25s ease;
    }

    .cs-pill:hover {
        color: rgba(255, 255, 255, 0.8);
        border-color: rgba(255, 68, 0, 0.35);
    }

    .cs-pill:focus-visible {
        outline: 2px solid var(--accent, #ff4400);
        outline-offset: 2px;
    }

    .cs-pill[aria-pressed="true"] {
        color: var(--accent, #ff4400);
        border-color: rgba(255, 68, 0, 0.5);
        background: rgba(255, 68, 0, 0.08);
    }

    .cs-meta {
        font-family: var(--font-body);
        font-size: 0.55rem;
        letter-spacing: 0.1em;
        text-transform: uppercase;
        color: var(--text-muted, rgba(255, 255, 255, 0.3));
        white-space: nowrap;
    }

    /* ═══════════════════════ Grid ═══════════════════════ */

    .cs-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(min(100%, 340px), 1fr));
        gap: 1.5rem;
        width: 100%;
    }

    .cs-item {
        height: 100%;
        animation: csReveal 0.5s both cubic-bezier(0.16, 1, 0.3, 1);
        animation-delay: calc(var(--idx, 1) * 0.06s);
    }

    @keyframes csReveal {
        from {
            opacity: 0;
            transform: translateY(18px);
        }
        to {
            opacity: 1;
            transform: translateY(0);
        }
    }

    /* ═══════════════════════ View more tile ═══════════════════════ */

    .cs-more {
        appearance: none;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        gap: 0.45rem;
        height: 100%;
        min-height: 320px;
        padding: 2rem 1.5rem;
        background:
            linear-gradient(135deg, rgba(255, 68, 0, 0.07), transparent 60%),
            var(--surface-dark, #0e1116);
        border: 1px solid rgba(255, 255, 255, 0.08);
        border-left: 3px solid rgba(255, 68, 0, 0.4);
        color: var(--text-primary, #f5f5f5);
        font-family: var(--font-body);
        text-align: center;
        cursor: pointer;
        position: relative;
        transition:
            transform 0.35s cubic-bezier(0.16, 1, 0.3, 1),
            box-shadow 0.35s ease,
            border-color 0.35s ease;
        animation: csReveal 0.5s both cubic-bezier(0.16, 1, 0.3, 1);
        animation-delay: 0.36s;
    }

    .cs-more:hover {
        border-color: rgba(255, 68, 0, 0.25);
        border-left-color: var(--accent, #ff4400);
        transform: translateY(-4px);
        box-shadow:
            5px 5px 0 rgba(0, 0, 0, 0.4),
            0 10px 28px rgba(0, 0, 0, 0.25);
    }

    .cs-more:focus-visible {
        outline: 2px solid var(--accent, #ff4400);
        outline-offset: 4px;
    }

    .cs-more-count {
        font-family: var(--font-heading);
        font-size: clamp(2.5rem, 5vw, 3.75rem);
        font-weight: 800;
        line-height: 1;
        letter-spacing: -0.04em;
        color: var(--accent, #ff4400);
    }

    .cs-more-title {
        font-family: var(--font-heading);
        font-size: 0.8rem;
        font-weight: 700;
        text-transform: uppercase;
        letter-spacing: 0.18em;
    }

    .cs-more-sub {
        font-family: var(--font-body);
        font-size: 0.62rem;
        text-transform: uppercase;
        letter-spacing: 0.12em;
        color: var(--text-muted, rgba(255, 255, 255, 0.35));
    }

    .cs-more-icon {
        display: grid;
        place-items: center;
        width: 40px;
        height: 40px;
        margin-top: 1rem;
        border: 1px solid rgba(255, 255, 255, 0.12);
        color: rgba(255, 255, 255, 0.55);
        transition:
            transform 0.3s ease,
            color 0.3s ease,
            border-color 0.3s ease;
    }

    .cs-more:hover .cs-more-icon {
        transform: translateY(4px);
        color: var(--accent, #ff4400);
        border-color: rgba(255, 68, 0, 0.5);
    }

    /* ═══════════════════════ Empty state ═══════════════════════ */

    .cs-empty {
        padding: 2.5rem 0;
        text-align: center;
        color: var(--text-muted, rgba(255, 255, 255, 0.35));
        font-family: var(--font-body);
        font-size: 0.85rem;
    }

    /* ═══════════════════════ Light mode ═══════════════════════ */

    :where(:global(body.light-mode)) .cs-toolbar {
        border-bottom-color: rgba(0, 0, 0, 0.06);
    }

    :where(:global(body.light-mode)) .cs-group-label,
    :where(:global(body.light-mode)) .cs-meta {
        color: rgba(0, 0, 0, 0.35);
    }

    :where(:global(body.light-mode)) .cs-pill {
        color: rgba(0, 0, 0, 0.45);
        border-color: rgba(0, 0, 0, 0.1);
    }

    :where(:global(body.light-mode)) .cs-pill:hover {
        color: rgba(0, 0, 0, 0.8);
        border-color: rgba(217, 65, 0, 0.35);
    }

    :where(:global(body.light-mode)) .cs-pill[aria-pressed="true"] {
        color: var(--light-accent, #d94100);
        border-color: rgba(217, 65, 0, 0.5);
        background: rgba(217, 65, 0, 0.06);
    }

    :where(:global(body.light-mode)) .cs-more {
        background:
            linear-gradient(135deg, rgba(217, 65, 0, 0.05), transparent 60%),
            #fff;
        border-color: rgba(0, 0, 0, 0.08);
        border-left-color: rgba(217, 65, 0, 0.35);
        color: #18181b;
    }

    :where(:global(body.light-mode)) .cs-more:hover {
        border-color: rgba(217, 65, 0, 0.15);
        border-left-color: var(--light-accent, #d94100);
        box-shadow:
            5px 5px 0 rgba(0, 0, 0, 0.14),
            0 10px 28px rgba(217, 65, 0, 0.06);
    }

    :where(:global(body.light-mode)) .cs-more-count {
        color: var(--light-accent, #d94100);
    }

    :where(:global(body.light-mode)) .cs-more-sub {
        color: rgba(0, 0, 0, 0.4);
    }

    :where(:global(body.light-mode)) .cs-more-icon {
        border-color: rgba(0, 0, 0, 0.12);
        color: rgba(0, 0, 0, 0.45);
    }

    :where(:global(body.light-mode)) .cs-more:hover .cs-more-icon {
        color: var(--light-accent, #d94100);
        border-color: rgba(217, 65, 0, 0.5);
    }

    /* ═══════════════════════ Mobile ═══════════════════════ */

    @media (max-width: 768px) {
        .cs-grid {
            grid-template-columns: 1fr;
            gap: 1rem;
        }

        .cs-toolbar {
            flex-direction: column;
            align-items: stretch;
            row-gap: 0.55rem;
            padding-bottom: 0.75rem;
            margin-bottom: 1rem;
        }

        /* Each control row scrolls horizontally instead of wrapping —
           keeps one line per row on 320–420px screens */
        .cs-group {
            display: flex;
            overflow-x: auto;
            overflow-y: hidden;
            padding-bottom: 3px;
            scrollbar-width: none;
            -webkit-overflow-scrolling: touch;
        }

        .cs-group::-webkit-scrollbar {
            display: none;
        }

        .cs-pill {
            flex: 0 0 auto;
            min-height: 40px;
            padding: 0.5rem 0.9rem;
        }

        .cs-meta {
            font-size: 0.5rem;
            text-align: left;
            padding-top: 0.15rem;
        }

        .cs-more {
            min-height: 220px;
        }
    }

    /* Coarse pointers: comfortable tap targets for the pills */
    @media (hover: none) and (pointer: coarse) {
        .cs-pill {
            min-height: 42px;
        }
    }

    /* Very small screens */
    @media (max-width: 380px) {
        .cs-grid {
            gap: 0.85rem;
        }
    }

    /* ═══════════════════════ Reduced motion ═══════════════════════ */

    @media (prefers-reduced-motion: reduce) {
        .cs-item,
        .cs-more {
            animation: none;
        }

        .cs-pill,
        .cs-more,
        .cs-more-icon {
            animation: none;
            transition: none;
        }
    }
</style>
