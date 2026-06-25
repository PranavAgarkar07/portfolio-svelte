<script lang="ts">
    import type { Certificate } from "$lib/types";
    import { Lightbox } from "lightbox3";

    let { certificate, index }: { certificate: Certificate; index: number } = $props();

    let thumbFailed = $state(false);
    let imageLink: HTMLAnchorElement | undefined = $state();

    function openLightbox(e: MouseEvent) {
        e.preventDefault();
        try {
            Lightbox.instance.open(certificate.image_url, imageLink!);
        } catch {
            // fallback: navigate directly to the image
            window.open(certificate.image_url, "_blank");
        }
    }

    function formatDate(d: string): string {
        if (!d) return "";
        const parts = d.split("-");
        if (parts.length < 2) return parts[0];
        const months = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
        const monthIndex = parseInt(parts[1], 10) - 1;
        const month = months[monthIndex] || "";
        if (!month) return parts[0];
        if (parts.length >= 3 && parts[2]) {
            const day = parseInt(parts[2], 10);
            return isNaN(day) ? `${month} ${parts[0]}` : `${month} ${day}, ${parts[0]}`;
        }
        return `${month} ${parts[0]}`;
    }

    let showThumb = $derived(!!certificate.thumb_url && !thumbFailed);
</script>

<article class="cert-card" style="--idx: {index + 1}">
    {#if certificate.image_url}
        <a
            href={certificate.image_url}
            class="cert-image-link"
            class:has-thumb={showThumb}
            data-lightbox="certificates"
            data-caption="{certificate.title} — {certificate.issuer}"
            aria-label="View {certificate.title} certificate"
            tabindex="0"
            bind:this={imageLink}
            onclick={openLightbox}
        >
            {#if showThumb}
                <img
                    src={certificate.thumb_url}
                    alt="{certificate.title} certificate"
                    class="cert-thumb"
                    loading="lazy"
                    decoding="async"
                    onerror={() => { thumbFailed = true; }}
                />
                <span class="cert-thumb-overlay">
                    <span class="cert-thumb-idx">{String(index + 1).padStart(2, "0")}</span>
                    <span class="cert-thumb-label">VIEW</span>
                </span>
            {:else}
                <span class="cert-placeholder-inner">
                    <span class="cert-placeholder-idx">{String(index + 1).padStart(2, "0")}</span>
                    <span class="cert-placeholder-label">VIEW</span>
                </span>
            {/if}
            <span class="cert-image-hint">
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <circle cx="11" cy="11" r="8"/><path d="m21 21-4.35-4.35"/>
                </svg>
                Click to view
            </span>
            <img
                alt="" aria-hidden="true"
                class="cert-hidden-img"
            />
        </a>
    {/if}

    <div class="cert-content">
        <div class="cert-meta">
            {#if certificate.date}
                <time class="cert-date">{formatDate(certificate.date)}</time>
            {/if}
            {#if certificate.is_verified}
                <span class="cert-verified">
                    <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round">
                        <polyline points="20 6 9 17 4 12"/>
                    </svg>
                    Verified
                </span>
            {/if}
        </div>

        <h3 class="cert-title">{certificate.title}</h3>
        <p class="cert-issuer">{certificate.issuer}</p>

        {#if certificate.tags && certificate.tags.length > 0}
            <div class="cert-tags">
                {#each certificate.tags as tag}
                    <span class="cert-tag">{tag}</span>
                {/each}
            </div>
        {/if}

        {#if certificate.credential_url}
            <a href={certificate.credential_url} target="_blank" rel="noopener" class="cert-link">
                <span>Verify Credential</span>
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <line x1="7" y1="17" x2="17" y2="7"/><polyline points="7 7 17 7 17 17"/>
                </svg>
            </a>
        {/if}
    </div>
</article>

<style>
    .cert-card {
        display: flex;
        flex-direction: column;
        height: 100%;
        background: var(--surface-dark);
        border: 1px solid rgba(255, 255, 255, 0.06);
        border-left: 3px solid var(--accent);
        position: relative;
        transition:
            border-color 0.2s ease,
            transform 0.2s ease,
            box-shadow 0.2s ease;
    }

    .cert-card:hover {
        border-color: rgba(255, 68, 0, 0.3);
        border-left-color: var(--accent);
        transform: translateY(-2px);
        box-shadow: 0 4px 20px rgba(255, 68, 0, 0.06);
    }

    .cert-card:focus-within {
        outline: 2px solid var(--accent);
        outline-offset: 4px;
    }

    /* ── Image link (shared placeholder + thumb container) ── */

    .cert-image-link {
        display: flex;
        align-items: center;
        justify-content: center;
        aspect-ratio: 16 / 10;
        background: #070a0f;
        border-bottom: 1px solid rgba(255, 255, 255, 0.03);
        text-decoration: none;
        position: relative;
        overflow: hidden;
    }

    .cert-image-link::before {
        content: "";
        position: absolute;
        inset: 0;
        background:
            linear-gradient(135deg, transparent 40%, rgba(255, 68, 0, 0.03) 50%, transparent 60%),
            repeating-linear-gradient(
                0deg,
                transparent,
                transparent 2px,
                rgba(255, 255, 255, 0.015) 2px,
                rgba(255, 255, 255, 0.015) 4px
            );
        pointer-events: none;
        z-index: 1;
    }

    .cert-card:hover .cert-image-link::before {
        opacity: 0;
    }

    .cert-image-link.has-thumb {
        background: #000;
    }

    .cert-image-link.has-thumb::before {
        background: linear-gradient(0deg, rgba(0,0,0,0.6) 0%, rgba(0,0,0,0.1) 50%, transparent 100%);
        z-index: 2;
        transition: opacity 0.25s ease;
    }

    .cert-card:hover .cert-image-link.has-thumb::before {
        background: linear-gradient(0deg, rgba(0,0,0,0.8) 0%, rgba(0,0,0,0.2) 60%, transparent 100%);
        opacity: 1;
    }

    /* ── Thumbnail ── */

    .cert-thumb {
        position: absolute;
        inset: 0;
        width: 100%;
        height: 100%;
        object-fit: cover;
        z-index: 0;
    }

    /* ── Overlay (index + VIEW label on thumb) ── */

    .cert-thumb-overlay {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 2px;
        position: relative;
        z-index: 3;
        opacity: 0;
        transition: opacity 0.2s ease;
        pointer-events: none;
    }

    .cert-card:hover .cert-thumb-overlay {
        opacity: 1;
    }

    .cert-thumb-idx {
        font-family: var(--font-heading);
        font-size: clamp(2rem, 5vw, 3rem);
        font-weight: 900;
        color: rgba(255, 255, 255, 0.7);
        letter-spacing: -0.04em;
        line-height: 1;
        text-shadow: 0 2px 8px rgba(0,0,0,0.5);
    }

    .cert-thumb-label {
        font-family: var(--font-heading);
        font-size: 0.5rem;
        font-weight: 700;
        text-transform: uppercase;
        letter-spacing: 0.25em;
        color: rgba(255, 255, 255, 0.5);
        text-shadow: 0 1px 4px rgba(0,0,0,0.5);
    }

    /* ── Placeholder (no thumb or thumb failed) ── */

    .cert-placeholder-inner {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 2px;
        position: relative;
        z-index: 1;
        transition: opacity 0.2s ease;
    }

    .cert-placeholder-idx {
        font-family: var(--font-heading);
        font-size: clamp(2rem, 5vw, 3rem);
        font-weight: 900;
        color: rgba(255, 68, 0, 0.3);
        letter-spacing: -0.04em;
        line-height: 1;
    }

    .cert-placeholder-label {
        font-family: var(--font-heading);
        font-size: 0.5rem;
        font-weight: 700;
        text-transform: uppercase;
        letter-spacing: 0.25em;
        color: rgba(255, 255, 255, 0.15);
    }

    .cert-card:hover .cert-placeholder-idx {
        color: rgba(255, 68, 0, 0.5);
    }

    .cert-card:hover .cert-placeholder-label {
        color: rgba(255, 255, 255, 0.3);
    }

    /* ── Hint (bottom-right "Click to view") ── */

    .cert-image-hint {
        position: absolute;
        bottom: 8px;
        right: 8px;
        display: inline-flex;
        align-items: center;
        gap: 4px;
        font-family: var(--font-body);
        font-size: 0.5rem;
        color: rgba(255, 255, 255, 0.25);
        opacity: 0;
        transition: opacity 0.2s ease;
        pointer-events: none;
        z-index: 4;
    }

    .cert-card:hover .cert-image-hint {
        opacity: 1;
    }

    .cert-hidden-img {
        position: absolute;
        inset: 0;
        width: 100%;
        height: 100%;
        visibility: hidden;
        pointer-events: none;
        object-fit: cover;
    }

    /* ── Content ── */

    .cert-content {
        display: flex;
        flex-direction: column;
        gap: 0.3rem;
        padding: 1rem 1rem 1.1rem;
        flex: 1;
    }

    .cert-meta {
        display: flex;
        align-items: center;
        gap: 0.6rem;
        margin-bottom: 0.15rem;
    }

    .cert-date {
        font-family: var(--font-body);
        font-size: 0.6rem;
        color: var(--text-secondary);
        letter-spacing: 0.08em;
        text-transform: uppercase;
        opacity: 0.6;
    }

    .cert-verified {
        display: inline-flex;
        align-items: center;
        gap: 3px;
        font-family: var(--font-body);
        font-size: 0.55rem;
        font-weight: 600;
        text-transform: uppercase;
        letter-spacing: 0.08em;
        color: var(--accent);
        padding: 2px 6px;
        border: 1px solid rgba(255, 68, 0, 0.2);
        background: rgba(255, 68, 0, 0.06);
    }

    .cert-title {
        font-family: var(--font-heading);
        font-size: 0.95rem;
        font-weight: 600;
        color: var(--text-primary);
        line-height: 1.25;
        letter-spacing: 0.02em;
        text-transform: uppercase;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
    }

    .cert-issuer {
        font-family: var(--font-body);
        font-size: 0.7rem;
        color: var(--text-secondary);
        opacity: 0.7;
        letter-spacing: 0.04em;
        line-height: 1.3;
        margin-bottom: 0.1rem;
    }

    .cert-tags {
        display: flex;
        flex-wrap: wrap;
        gap: 4px;
    }

    .cert-tag {
        font-family: var(--font-body);
        font-size: 0.55rem;
        font-weight: 500;
        letter-spacing: 0.06em;
        text-transform: uppercase;
        color: var(--text-secondary);
        padding: 2px 7px;
        border: 1px solid rgba(255, 255, 255, 0.06);
        background: rgba(255, 255, 255, 0.02);
    }

    .cert-link {
        display: inline-flex;
        align-items: center;
        gap: 5px;
        margin-top: auto;
        padding-top: 0.75rem;
        font-family: var(--font-heading);
        font-size: 0.55rem;
        font-weight: 600;
        text-transform: uppercase;
        letter-spacing: 0.1em;
        color: rgba(255, 255, 255, 0.35);
        text-decoration: none;
        border-top: 1px solid rgba(255, 255, 255, 0.04);
        transition: color 0.2s ease;
    }

    .cert-link svg {
        transition: transform 0.15s ease;
        opacity: 0.4;
    }

    .cert-link:hover {
        color: var(--accent);
    }

    .cert-link:hover svg {
        transform: translate(2px, -2px);
        opacity: 0.8;
    }

    /* ── Light mode ── */

    :where(:global(body.light-mode)) .cert-card {
        background: var(--surface-dark, #fff);
        border-color: var(--grid-line, rgba(0,0,0,0.06));
        border-left-color: var(--light-accent, #d94100);
    }

    :where(:global(body.light-mode)) .cert-card:hover {
        border-color: rgba(217, 65, 0, 0.15);
        border-left-color: var(--light-accent, #d94100);
    }

    :where(:global(body.light-mode)) .cert-title {
        color: #18181b;
    }

    :where(:global(body.light-mode)) .cert-link {
        color: rgba(0, 0, 0, 0.25);
        border-top-color: var(--grid-line, rgba(0,0,0,0.06));
    }

    :where(:global(body.light-mode)) .cert-link:hover {
        color: var(--light-accent, #d94100);
    }

    :where(:global(body.light-mode)) .cert-verified {
        border-color: rgba(217, 65, 0, 0.15);
        background: rgba(217, 65, 0, 0.05);
        color: var(--light-accent, #d94100);
    }

    :where(:global(body.light-mode)) .cert-tag {
        border-color: var(--grid-line, rgba(0,0,0,0.06));
        background: rgba(0, 0, 0, 0.02);
    }

    @media (prefers-reduced-motion: reduce) {
        .cert-card,
        .cert-image-link,
        .cert-thumb-overlay,
        .cert-placeholder-inner,
        .cert-placeholder-idx,
        .cert-placeholder-label,
        .cert-image-hint,
        .cert-link svg { transition: none; }
    }
</style>
