<script lang="ts">
    import { onMount } from "svelte";
    import { Lightbox } from "lightbox3";
    import { Badge, Tag, Icon } from "$lib/components/ui";
    import type { Certificate } from "$lib/types";
    import { certPlaceholders } from "$lib/cert-placeholders";

    let { certificate, index }: { certificate: Certificate; index: number } = $props();

    function formatDate(d: string): string {
        if (!d) return d;
        const parts = d.split("-");
        if (parts.length < 2) return d;
        const months = ["January","February","March","April","May","June","July","August","September","October","November","December"];
        const month = months[parseInt(parts[1], 10) - 1];
        if (!month) return d;
        if (parts.length === 3) {
            const day = parseInt(parts[2], 10);
            return `${month} ${day}, ${parts[0]}`;
        }
        return `${month} ${parts[0]}`;
    }

    function thumbUrl(fullUrl: string): string {
        const name = fullUrl.split("/").pop() || "";
        const base = name.replace(/\.\w+$/, "");
        const parts = fullUrl.split("/");
        parts.pop();
        parts.push("thumbs", base + ".webp");
        return parts.join("/");
    }

    function placeholderBg(fullUrl: string): string {
        const name = fullUrl.split("/").pop() || "";
        return certPlaceholders[name] || "";
    }

    let cardEl = $state<HTMLElement | null>(null);
    let watermarkEl = $state<HTMLElement | null>(null);
    let reducedMotion = $state(false);

    let fullUrl = $derived(certificate.image_url);
    let thumb = $derived(thumbUrl(fullUrl));
    let placeholderDataUri = $derived(placeholderBg(fullUrl));

    onMount(() => {
        reducedMotion = window.matchMedia("(prefers-reduced-motion: reduce)").matches;
    });

    function handleMouseMove(e: MouseEvent) {
        if (reducedMotion || !cardEl || !watermarkEl) return;
        watermarkEl.style.transition = "none";
        const rect = cardEl.getBoundingClientRect();
        const x = ((e.clientX - rect.left) / rect.width - 0.5) * 2;
        const y = ((e.clientY - rect.top) / rect.height - 0.5) * 2;
        watermarkEl.style.transform = `translate(${x * 12}px, ${y * 12}px)`;
    }

    function handleMouseLeave() {
        if (!watermarkEl) return;
        watermarkEl.style.transition = "";
        watermarkEl.style.transform = "";
    }
</script>

<article
    class="cert-card"
    class:cert-verified={certificate.is_verified}
    bind:this={cardEl}
    onmousemove={handleMouseMove}
    onmouseleave={handleMouseLeave}
    style="--cert-index: {index}"
>
    {#if certificate.is_verified}
        <div class="cert-ribbon" aria-hidden="true">
            <span>VERIFIED</span>
        </div>
    {/if}

    {#if certificate.image_url}
        <div class="cert-media" style={placeholderDataUri ? `background-image: url('${placeholderDataUri}')` : ""}>
            <div class="cert-media-border" aria-hidden="true"></div>
            <a
                href={certificate.image_url}
                data-lightbox="cert-{certificate.id}"
                data-caption="{certificate.title} — {certificate.issuer}"
                class="cert-image-link"
                onclick={(e) => {
                    e.preventDefault();
                    Lightbox.instance.open(certificate.image_url, e.currentTarget);
                }}
            >
                <img
                    src={thumb}
                    alt={certificate.title}
                    class="cert-img"
                    loading="lazy"
                    decoding="async"
                    fetchpriority="low"
                />
                <div class="cert-image-overlay">
                    <div class="overlay-scanline" aria-hidden="true"></div>
                    <Icon name="search" size={20} />
                    <span class="overlay-text">Inspect</span>
                    <span class="overlay-hint">click to enlarge</span>
                </div>
            </a>
        </div>
    {/if}

    <div class="cert-body">
        <div class="cert-header">
            <div class="cert-meta">
                <span class="cert-number"
                    >CERT {String(index + 1).padStart(2, "0")}</span
                >
            </div>
            <h3 class="cert-title">{certificate.title}</h3>
            <span class="cert-issuer">{certificate.issuer}</span>
        </div>

        <div class="cert-details">
            {#if certificate.date}
                <div class="cert-date-row">
                    <Icon name="calendar" size={12} color="var(--accent)" />
                    <span>{formatDate(certificate.date)}</span>
                </div>
            {/if}
            {#if certificate.id}
                <div class="cert-fingerprint" title="Certificate ID">
                    <span class="fp-dot"></span>
                    <span>CRD-{String(certificate.id).padStart(4, "0")}</span>
                </div>
            {/if}
        </div>

        {#if certificate.tags && certificate.tags.length > 0}
            <div class="cert-tags">
                {#each certificate.tags as tag}
                    <Tag>{tag}</Tag>
                {/each}
            </div>
        {/if}

        {#if certificate.credential_url}
            <div class="cert-links">
                <a
                    href={certificate.credential_url}
                    target="_blank"
                    rel="noopener"
                    class="cert-link"
                >
                    <Icon name="external-link" size={14} />
                    <span>Verify Credential</span>
                    <span class="link-arrow"><Icon name="arrow-up-right" size={12} /></span>
                </a>
            </div>
        {/if}

        <span class="cert-number-watermark" bind:this={watermarkEl} aria-hidden="true">
            {String(index + 1).padStart(2, "0")}
        </span>
    </div>
</article>

<style>
    .cert-card {
        position: relative;
        background: #0f131a;
        padding: 0;
        border: 1px solid rgba(255, 255, 255, 0.06);
        box-shadow: 0 2px 12px rgba(0, 0, 0, 0.4);
        transition:
            border-color 0.25s ease,
            box-shadow 0.25s ease,
            transform 0.25s ease,
            background 0.25s ease;
        display: flex;
        flex-direction: column;
        height: 100%;
        overflow: hidden;
    }

    .cert-card:hover {
        border-color: rgba(255, 68, 0, 0.25);
        box-shadow: 0 4px 24px rgba(255, 68, 0, 0.08), 0 2px 12px rgba(0, 0, 0, 0.4);
        transform: translateY(-3px);
        background: #121720;
    }

    .cert-card:focus-visible {
        outline: 3px solid var(--accent);
        outline-offset: 4px;
    }

    /* ── Ribbon ── */

    .cert-ribbon {
        position: absolute;
        top: 0;
        right: 0;
        width: 120px;
        height: 120px;
        overflow: hidden;
        z-index: 5;
        pointer-events: none;
    }

    .cert-ribbon span {
        position: absolute;
        top: 22px;
        right: -30px;
        width: 160px;
        padding: 4px 0;
        background: var(--accent);
        color: #000;
        font-family: var(--font-heading);
        font-size: 0.6rem;
        font-weight: 800;
        letter-spacing: 0.15em;
        text-transform: uppercase;
        text-align: center;
        transform: rotate(45deg);
        box-shadow: 0 2px 6px rgba(0, 0, 0, 0.4);
        line-height: 1;
    }

    /* ── Media / Document scan ── */

    .cert-media {
        position: relative;
        overflow: hidden;
        background-color: #070a0f;
        background-size: cover;
        background-position: center;
        background-repeat: no-repeat;
        border-bottom: 1px solid rgba(255, 255, 255, 0.04);
        aspect-ratio: 16 / 10;
    }

    .cert-media-border {
        position: absolute;
        inset: 8px;
        border: 1px solid rgba(255, 255, 255, 0.04);
        z-index: 1;
        pointer-events: none;
    }

    .cert-image-link {
        display: block;
        width: 100%;
        height: 100%;
        position: relative;
        cursor: pointer;
        text-decoration: none;
    }

    .cert-img {
        display: block;
        width: 100%;
        height: 100%;
        object-fit: cover;
        transition: transform 0.5s ease, filter 0.35s ease;
    }

    .cert-card:hover .cert-img {
        transform: scale(1.05);
        filter: brightness(0.7);
    }

    .cert-image-overlay {
        position: absolute;
        inset: 0;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        gap: 0.35rem;
        background: rgba(0, 0, 0, 0.65);
        opacity: 0;
        transition: opacity 0.3s ease;
        color: #fff;
        font-family: var(--font-heading);
    }

    .cert-card:hover .cert-image-overlay {
        opacity: 1;
    }

    .overlay-scanline {
        position: absolute;
        inset: 0;
        background: repeating-linear-gradient(
            0deg,
            transparent,
            transparent 2px,
            rgba(255, 68, 0, 0.04) 2px,
            rgba(255, 68, 0, 0.04) 4px
        );
        pointer-events: none;
    }

    .overlay-text {
        font-size: 0.8rem;
        font-weight: 700;
        letter-spacing: 0.15em;
        text-transform: uppercase;
        margin-top: 0.2rem;
    }

    .overlay-hint {
        font-family: var(--font-body);
        font-size: 0.6rem;
        opacity: 0.5;
        letter-spacing: 0.08em;
    }

    /* ── Body ── */

    .cert-body {
        position: relative;
        padding: 1rem 1.15rem 1.15rem;
        display: flex;
        flex-direction: column;
        flex: 1;
        gap: 0.5rem;
        min-height: 0;
    }

    .cert-number-watermark {
        position: absolute;
        bottom: -0.3rem;
        right: -0.2rem;
        font-family: var(--font-heading);
        font-size: clamp(4rem, 8vw, 7rem);
        font-weight: 900;
        line-height: 1;
        letter-spacing: -0.06em;
        pointer-events: none;
        user-select: none;
        color: transparent;
        -webkit-text-stroke: 1.5px var(--accent);
        opacity: 0.18;
        z-index: 0;
        transition: transform 0.3s ease, opacity 0.2s ease;
        will-change: transform;
    }

    .cert-card:hover .cert-number-watermark {
        opacity: 0.35;
    }

    .cert-header {
        position: relative;
        z-index: 1;
        display: flex;
        flex-direction: column;
        gap: 0.2rem;
    }

    .cert-meta {
        display: flex;
        align-items: center;
        gap: 0.6rem;
    }

    .cert-number {
        font-family: var(--font-body);
        font-size: 0.55rem;
        color: var(--text-secondary);
        letter-spacing: 3px;
        text-transform: uppercase;
        opacity: 0.6;
    }

    .cert-number::before {
        content: "[";
        color: var(--accent);
        opacity: 0.7;
        margin-right: 1px;
    }

    .cert-number::after {
        content: "]";
        color: var(--accent);
        opacity: 0.7;
        margin-left: 1px;
    }

    .cert-title {
        font-family: var(--font-heading);
        font-size: 0.95rem;
        font-weight: 600;
        color: #fff;
        line-height: 1.25;
        letter-spacing: 0.02em;
        text-transform: uppercase;
        position: relative;
        z-index: 1;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
    }

    .cert-issuer {
        font-family: var(--font-body);
        font-size: 0.72rem;
        color: var(--text-secondary);
        letter-spacing: 0.05em;
        opacity: 0.7;
    }

    .cert-details {
        display: flex;
        align-items: center;
        gap: 0.65rem;
        flex-wrap: wrap;
        padding: 0.35rem 0;
        border-top: 1px dashed rgba(255, 255, 255, 0.06);
        border-bottom: 1px dashed rgba(255, 255, 255, 0.06);
    }

    .cert-date-row {
        display: flex;
        align-items: center;
        gap: 0.35rem;
        font-size: 0.65rem;
        color: var(--text-secondary);
        font-family: var(--font-body);
    }

    .cert-fingerprint {
        display: flex;
        align-items: center;
        gap: 0.3rem;
        font-size: 0.55rem;
        font-family: var(--font-body);
        color: rgba(255, 255, 255, 0.2);
        letter-spacing: 0.05em;
        font-variant-numeric: tabular-nums;
    }

    .fp-dot {
        width: 3px;
        height: 3px;
        background: var(--accent);
        opacity: 0.5;
        flex-shrink: 0;
    }

    .cert-tags {
        display: flex;
        flex-wrap: wrap;
        gap: 0.25rem;
        position: relative;
        z-index: 1;
    }

    .cert-links {
        margin-top: auto;
        padding-top: 0.6rem;
        border-top: 1px solid rgba(255, 255, 255, 0.04);
        position: relative;
        z-index: 1;
    }

    .cert-link {
        display: inline-flex;
        align-items: center;
        gap: 5px;
        font-family: var(--font-heading);
        font-size: 0.55rem;
        font-weight: 600;
        text-transform: uppercase;
        letter-spacing: 0.1em;
        color: rgba(255, 255, 255, 0.5);
        text-decoration: none;
        padding: 4px 8px;
        border: 1px solid rgba(255, 255, 255, 0.06);
        background: transparent;
        transition:
            color 0.2s ease,
            border-color 0.2s ease,
            background 0.2s ease;
        cursor: pointer;
    }

    .cert-link:hover {
        color: var(--accent);
        border-color: rgba(255, 68, 0, 0.2);
        background: rgba(255, 68, 0, 0.04);
    }

    .cert-link:focus-visible {
        outline: 2px solid var(--accent);
        outline-offset: 2px;
    }

    .link-arrow {
        transition: transform 0.15s ease;
        opacity: 0.4;
    }

    .cert-link:hover .link-arrow {
        transform: translate(3px, -3px);
        opacity: 0.8;
    }

    /* ── Light mode ── */

    :global(body.light-mode) .cert-card {
        background: #fff;
        border-color: var(--wire-color);
        box-shadow: 0 1px 3px rgba(0,0,0,0.04);
    }

    :global(body.light-mode) .cert-card:hover {
        border-color: rgba(217, 65, 0, 0.2);
        box-shadow: 0 8px 24px rgba(0,0,0,0.06);
        background: #fff;
    }

    :global(body.light-mode) .cert-title {
        color: #18181b;
    }

    :global(body.light-mode) .cert-link {
        color: rgba(0, 0, 0, 0.35);
        border-color: rgba(0, 0, 0, 0.06);
    }

    :global(body.light-mode) .cert-link:hover {
        color: var(--accent);
        border-color: rgba(217, 65, 0, 0.15);
        background: rgba(217, 65, 0, 0.04);
    }

    :global(body.light-mode) .cert-links {
        border-top-color: var(--grid-line);
    }

    :global(body.light-mode) .cert-media {
        background: #e8eaed;
    }

    :global(body.light-mode) .cert-number-watermark {
        -webkit-text-stroke: 1.5px var(--accent);
        opacity: 0.08;
    }

    :global(body.light-mode) .cert-card:hover .cert-number-watermark {
        opacity: 0.18;
    }

    :global(body.light-mode) .cert-details {
        border-top-color: var(--grid-line);
        border-bottom-color: var(--grid-line);
    }

    :global(body.light-mode) .cert-fingerprint {
        color: rgba(0, 0, 0, 0.18);
    }

    :global(body.light-mode) .cert-ribbon span {
        color: #18181b;
    }

    :global(body.light-mode) .overlay-scanline {
        background: repeating-linear-gradient(
            0deg,
            transparent,
            transparent 2px,
            rgba(217, 65, 0, 0.05) 2px,
            rgba(217, 65, 0, 0.05) 4px
        );
    }

    :global(body.light-mode) .cert-media-border {
        border-color: var(--grid-line);
    }

    @media (prefers-reduced-motion: reduce) {
        .cert-number-watermark { transition: none; }
        .cert-card { transition: none; }
        .cert-img { transition: none; }
    }
</style>
