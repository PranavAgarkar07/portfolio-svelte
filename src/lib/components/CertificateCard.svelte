<script lang="ts">
    import type { Certificate } from "$lib/types";
    import { Lightbox } from "lightbox3";

    let { certificate, index }: { certificate: Certificate; index: number } =
        $props();

    let thumbFailed = $state(false);
    let imageLink: HTMLAnchorElement | undefined = $state();
    let cardEl: HTMLElement | undefined = $state();

    function openLightbox(e: MouseEvent) {
        e.preventDefault();
        if (certificate.kind === "pdf") {
            window.open(certificate.file_url, "_blank");
            return;
        }
        try {
            (Lightbox as any).instance.open(certificate.image_url, imageLink!);
        } catch {
            window.open(certificate.image_url, "_blank");
        }
    }

    function formatDate(d: string): string {
        if (!d) return "";
        const parts = d.split("-");
        if (parts.length < 2) return parts[0];
        const months = [
            "Jan",
            "Feb",
            "Mar",
            "Apr",
            "May",
            "Jun",
            "Jul",
            "Aug",
            "Sep",
            "Oct",
            "Nov",
            "Dec",
        ];
        const monthIndex = parseInt(parts[1], 10) - 1;
        const month = months[monthIndex] || "";
        if (!month) return parts[0];
        if (parts.length >= 3 && parts[2]) {
            const day = parseInt(parts[2], 10);
            return isNaN(day)
                ? `${month} ${parts[0]}`
                : `${month} ${day}, ${parts[0]}`;
        }
        return `${month} ${parts[0]}`;
    }

    let showThumb = $derived(!!certificate.image_url && !thumbFailed);
    const isPdf = $derived(certificate.kind === "pdf");
    const openLabel = $derived(isPdf ? "PDF" : "View");
</script>

<article
    class="cert-card"
    class:has-thumb={showThumb}
    style="--idx: {index + 1}"
    bind:this={cardEl}
>
    {#if certificate.image_url}
        <a
            href={isPdf ? certificate.file_url : certificate.image_url}
            class="cert-image-link"
            data-lightbox={isPdf ? undefined : "certificates"}
            data-caption="{certificate.title} — {certificate.issuer}"
            aria-label={isPdf
                ? `Open ${certificate.title} certificate (PDF) in a new tab`
                : `View ${certificate.title} certificate`}
            tabindex="0"
            bind:this={imageLink}
            onclick={openLightbox}
        >
            {#if showThumb}
                <img
                    src={certificate.image_url}
                    alt="{certificate.title} certificate"
                    class="cert-thumb"
                    loading="lazy"
                    decoding="async"
                    onerror={() => {
                        thumbFailed = true;
                    }}
                />
                <!-- Persistent chips: always visible (incl. touch, where the
                     hover-only overlay stays hidden) -->
                <span class="cert-chip cert-chip-type" class:accent={isPdf}>
                    {isPdf ? "PDF" : "Image"}
                </span>
                <span class="cert-chip cert-chip-idx" aria-hidden="true">
                    {String(index + 1).padStart(2, "0")}
                </span>
                <span class="cert-image-overlay">
                    <span class="cert-image-action">
                        <svg
                            width="14"
                            height="14"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="2"
                            stroke-linecap="round"
                            stroke-linejoin="round"
                        >
                            <circle cx="11" cy="11" r="8" /><path
                                d="m21 21-4.35-4.35"
                            />
                        </svg>
                        <span>{openLabel}</span>
                    </span>
                </span>
            {:else}
                <span class="cert-image-index"
                    >{String(index + 1).padStart(2, "0")}</span
                >
            {/if}
            <img alt="" aria-hidden="true" class="cert-hidden-img" />
        </a>
    {/if}

    <div class="cert-content">
        <div class="cert-meta">
            <span class="cert-cat">{certificate.category_label}</span>
            {#if certificate.date}
                <time class="cert-date">{formatDate(certificate.date)}</time>
            {/if}
            {#if certificate.is_verified}
                <span class="cert-verified">
                    <svg
                        class="cert-verified-check"
                        width="10"
                        height="10"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="3.5"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        aria-hidden="true"
                    >
                        <path d="M20 6 9 17l-5-5" />
                    </svg>
                    Verified
                </span>
            {/if}
        </div>

        <h3 class="cert-title">{certificate.title}</h3>
        <p class="cert-issuer">{certificate.issuer}</p>

        {#if certificate.note}
            <p class="cert-note">{certificate.note}</p>
        {/if}

        {#if certificate.tags && certificate.tags.length > 0}
            <div class="cert-tags">
                {#each certificate.tags as tag}
                    <span class="cert-tag">{tag}</span>
                {/each}
            </div>
        {/if}

        {#if certificate.credential_url}
            <a
                href={certificate.credential_url}
                target="_blank"
                rel="noopener"
                class="cert-credential"
            >
                <span>Verify Credential</span>
                <svg
                    width="14"
                    height="14"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    class="cert-credential-icon"
                >
                    <path d="M7 7h10v10" />
                    <path d="M7 17 21 3" />
                </svg>
            </a>
        {/if}
    </div>
</article>

<style>
    /* ═══════════════════════════ Card Container ═══════════════════════════ */

    .cert-card {
        display: flex;
        flex-direction: column;
        height: 100%;
        background: var(--surface-dark, #0d1117);
        border: 1px solid rgba(255, 255, 255, 0.06);
        border-left: 3px solid rgba(255, 68, 0, 0.4);
        position: relative;
        transition:
            transform 0.35s cubic-bezier(0.16, 1, 0.3, 1),
            box-shadow 0.35s ease,
            border-color 0.35s ease;
        animation: cardEnter 0.5s both cubic-bezier(0.16, 1, 0.3, 1);
        animation-delay: calc(var(--idx, 0) * 0.06s);
    }

    @keyframes cardEnter {
        from {
            opacity: 0;
            transform: translateY(16px);
        }
        to {
            opacity: 1;
            transform: translateY(0);
        }
    }

    .cert-card::after {
        content: "";
        position: absolute;
        inset: -1px;
        pointer-events: none;
        opacity: 0;
        transition: opacity 0.35s ease;
        box-shadow: 0 0 0 1px rgba(255, 68, 0, 0.08);
    }

    .cert-card:hover {
        border-color: rgba(255, 68, 0, 0.25);
        border-left-color: var(--accent, #ff4400);
        transform: translateY(-4px);
        /* Neobrutalist offset shadow — matches the site's hero language */
        box-shadow:
            5px 5px 0 rgba(0, 0, 0, 0.4),
            0 10px 28px rgba(0, 0, 0, 0.25);
    }

    .cert-card:hover::after {
        opacity: 1;
    }

    .cert-card:focus-within {
        outline: 2px solid var(--accent, #ff4400);
        outline-offset: 4px;
    }

    /* ═══════════════════════════ Image Area ═══════════════════════════ */

    .cert-image-link {
        display: flex;
        align-items: center;
        justify-content: center;
        /* Landscape 4:3 frame (w:h) — every certificate in the set is a
           landscape document (A4-landscape ≈1.41, letter-landscape ≈1.29),
           so the horizontal frame shows each one essentially complete */
        aspect-ratio: 4 / 3;
        background: #070a0f;
        border-bottom: 1px solid rgba(255, 255, 255, 0.03);
        text-decoration: none;
        position: relative;
        overflow: hidden;
        color: inherit;
    }

    .cert-image-link::before {
        content: "";
        position: absolute;
        inset: 0;
        background: linear-gradient(
            135deg,
            transparent 30%,
            rgba(255, 68, 0, 0.04) 50%,
            transparent 70%
        );
        pointer-events: none;
        z-index: 1;
    }

    .cert-image-link.has-thumb::before {
        background: linear-gradient(
            0deg,
            rgba(0, 0, 0, 0.55) 0%,
            rgba(0, 0, 0, 0.05) 50%,
            transparent 100%
        );
        z-index: 2;
        transition: opacity 0.35s ease;
    }

    .cert-card:hover .cert-image-link.has-thumb::before {
        background: linear-gradient(
            0deg,
            rgba(0, 0, 0, 0.75) 0%,
            rgba(0, 0, 0, 0.15) 55%,
            transparent 100%
        );
        opacity: 1;
    }

    /* ── Thumbnail ── */

    .cert-thumb {
        position: absolute;
        inset: 0;
        width: 100%;
        height: 100%;
        /* contain, not cover — the complete certificate stays visible;
           certs deviate ≤ ~6% from 4:3, so letterbox bars stay minimal */
        object-fit: contain;
        z-index: 0;
        transition: transform 0.3s cubic-bezier(0.16, 1, 0.3, 1);
    }

    .cert-card:hover .cert-thumb {
        transform: scale(1.04);
    }

    /* ── Persistent chips (type + index) ── */

    .cert-chip {
        position: absolute;
        top: 10px;
        z-index: 4;
        font-family: var(--font-body);
        font-size: 0.52rem;
        font-weight: 700;
        text-transform: uppercase;
        letter-spacing: 0.14em;
        line-height: 1;
        padding: 4px 8px;
        color: rgba(255, 255, 255, 0.8);
        background: #0a0d13;
        border: 1px solid rgba(255, 255, 255, 0.12);
        pointer-events: none;
    }

    .cert-chip-type {
        left: 10px;
    }

    .cert-chip-type.accent {
        color: var(--accent, #ff4400);
        border-color: rgba(255, 68, 0, 0.55);
        background: rgba(255, 68, 0, 0.1);
    }

    .cert-chip-idx {
        right: 10px;
        font-family: var(--font-heading);
        color: rgba(255, 255, 255, 0.55);
        border-color: rgba(255, 255, 255, 0.08);
    }

    .cert-card:hover .cert-chip-idx {
        color: rgba(255, 255, 255, 0.85);
        border-color: rgba(255, 68, 0, 0.4);
    }

    /* ── Image overlay (hover action) ── */

    .cert-image-overlay {
        display: flex;
        justify-content: center;
        position: relative;
        z-index: 3;
        transform: translateY(6px);
        opacity: 0;
        transition:
            transform 0.3s cubic-bezier(0.16, 1, 0.3, 1),
            opacity 0.3s ease;
        pointer-events: none;
    }

    .cert-card:hover .cert-image-overlay {
        transform: translateY(0);
        opacity: 1;
    }

    /* Hover not needed on touch — chips stay visible instead */
    @media (hover: none) and (pointer: coarse) {
        .cert-image-overlay {
            display: none;
        }
    }

    /* ═══════════════════════ Mobile ═══════════════════════ */

    @media (max-width: 768px) {
        .cert-content {
            padding: 0.9rem 0.95rem 1rem;
        }

        .cert-title {
            font-size: 0.875rem;
        }

        .cert-credential {
            padding-top: 0.6rem;
            padding-bottom: 0.4rem;
        }

        .cert-meta {
            row-gap: 0.35rem;
        }

        .cert-image-index {
            font-size: clamp(1.9rem, 7vw, 2.5rem);
        }
    }

    /* Very small screens (320–380px) */
    @media (max-width: 380px) {
        .cert-content {
            padding: 0.8rem 0.85rem 0.9rem;
        }

        .cert-title {
            font-size: 0.82rem;
        }

        .cert-tags {
            gap: 2px;
        }
    }

    .cert-image-index {
        position: relative;
        z-index: 1;
        font-family: var(--font-heading);
        font-size: clamp(2rem, 4vw, 2.8rem);
        font-weight: 900;
        color: rgba(255, 68, 0, 0.35);
        letter-spacing: -0.04em;
        line-height: 1;
        text-shadow: 0 2px 12px rgba(0, 0, 0, 0.4);
        transition: color 0.25s ease;
    }

    .cert-card:hover .cert-image-index {
        color: rgba(255, 68, 0, 0.6);
    }

    .cert-card:not(.has-thumb) .cert-image-index {
        position: relative;
        z-index: 1;
        font-size: clamp(2.5rem, 5vw, 3.5rem);
        color: rgba(255, 68, 0, 0.3);
    }

    .cert-card:not(.has-thumb):hover .cert-image-index {
        color: rgba(255, 68, 0, 0.55);
    }

    .cert-image-action {
        display: inline-flex;
        align-items: center;
        gap: 4px;
        font-family: var(--font-heading);
        font-size: 0.55rem;
        font-weight: 700;
        text-transform: uppercase;
        letter-spacing: 0.2em;
        color: rgba(255, 255, 255, 0.5);
        text-shadow: 0 1px 6px rgba(0, 0, 0, 0.5);
    }

    .cert-image-action svg {
        width: 12px;
        height: 12px;
        opacity: 0.6;
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

    /* ═══════════════════════════ Content ═══════════════════════════ */

    .cert-content {
        display: flex;
        flex-direction: column;
        gap: 0.3rem;
        padding: 1rem 1.1rem 1.1rem;
        flex: 1;
    }

    /* ── Meta row: date + verified badge ── */

    .cert-meta {
        display: flex;
        align-items: center;
        gap: 0.6rem;
        margin-bottom: 0.1rem;
        flex-wrap: wrap;
    }

    .cert-date {
        font-family: var(--font-body);
        font-size: 0.6rem;
        font-weight: 500;
        color: var(--text-secondary);
        letter-spacing: 0.08em;
        text-transform: uppercase;
        opacity: 0.55;
    }

    .cert-cat {
        display: inline-block;
        font-family: var(--font-heading);
        font-size: 0.5rem;
        font-weight: 700;
        text-transform: uppercase;
        letter-spacing: 0.08em;
        color: var(--text-secondary);
        padding: 2px 7px;
        border: 1px solid rgba(255, 255, 255, 0.08);
        background: rgba(255, 255, 255, 0.02);
    }

    .cert-note {
        font-family: var(--font-body);
        font-size: 0.55rem;
        letter-spacing: 0.04em;
        color: var(--text-muted, rgba(255, 255, 255, 0.3));
        line-height: 1.4;
    }

    .cert-verified {
        display: inline-flex;
        align-items: center;
        gap: 4px;
        font-family: var(--font-body);
        font-size: 0.5rem;
        font-weight: 700;
        text-transform: uppercase;
        letter-spacing: 0.08em;
        color: var(--accent, #ff4400);
        padding: 2px 7px;
        border: 1px solid rgba(255, 68, 0, 0.15);
        background: rgba(255, 68, 0, 0.05);
    }

    .cert-verified-check {
        width: 10px;
        height: 10px;
        flex-shrink: 0;
    }

    /* ── Title ── */

    .cert-title {
        font-family: var(--font-heading);
        font-size: 0.9rem;
        font-weight: 700;
        color: var(--text-primary);
        line-height: 1.3;
        letter-spacing: 0.015em;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
    }

    /* ── Issuer ── */

    .cert-issuer {
        font-family: var(--font-body);
        font-size: 0.7rem;
        color: var(--text-secondary);
        opacity: 0.65;
        letter-spacing: 0.03em;
        line-height: 1.35;
        margin-bottom: 0.05rem;
    }

    /* ── Tags ── */

    .cert-tags {
        display: flex;
        flex-wrap: wrap;
        gap: 3px;
        margin-top: 0.1rem;
    }

    .cert-tag {
        font-family: var(--font-body);
        font-size: 0.55rem;
        font-weight: 500;
        letter-spacing: 0.06em;
        text-transform: uppercase;
        color: rgba(255, 255, 255, 0.38);
        padding: 2px 6px;
        border: 1px solid rgba(255, 255, 255, 0.05);
        background: rgba(255, 255, 255, 0.015);
    }

    /* ── Credential link ── */

    .cert-credential {
        display: inline-flex;
        align-items: center;
        gap: 5px;
        min-height: 44px;
        margin-top: auto;
        padding-top: 0.7rem;
        padding-bottom: 0.5rem;
        font-family: var(--font-heading);
        font-size: 0.55rem;
        font-weight: 700;
        text-transform: uppercase;
        letter-spacing: 0.1em;
        color: rgba(255, 255, 255, 0.42);
        text-decoration: none;
        border-top: 1px solid rgba(255, 255, 255, 0.04);
        transition:
            color 0.25s ease,
            gap 0.25s ease,
            border-color 0.25s ease;
    }

    .cert-credential-icon {
        transition: transform 0.25s cubic-bezier(0.16, 1, 0.3, 1);
        opacity: 0.45;
    }

    .cert-credential:hover {
        color: var(--accent, #ff4400);
        gap: 7px;
        border-top-color: rgba(255, 68, 0, 0.3);
    }

    .cert-credential:hover .cert-credential-icon {
        transform: translate(2px, -2px);
        opacity: 0.8;
    }

    /* ═══════════════════════════ Light Mode ═══════════════════════════ */

    :where(:global(body.light-mode)) .cert-card {
        background: #fff;
        border-color: rgba(0, 0, 0, 0.06);
        border-left-color: rgba(217, 65, 0, 0.35);
    }

    :where(:global(body.light-mode)) .cert-card:hover {
        border-color: rgba(217, 65, 0, 0.15);
        border-left-color: var(--light-accent, #d94100);
        box-shadow:
            5px 5px 0 rgba(0, 0, 0, 0.14),
            0 10px 28px rgba(217, 65, 0, 0.06);
    }

    :where(:global(body.light-mode)) .cert-card::after {
        box-shadow: 0 0 0 1px rgba(217, 65, 0, 0.04);
    }

    :where(:global(body.light-mode)) .cert-image-link {
        background: #f4f4f5;
        border-bottom-color: rgba(0, 0, 0, 0.04);
    }

    :where(:global(body.light-mode)) .cert-image-link::before {
        background: linear-gradient(
            135deg,
            transparent 30%,
            rgba(217, 65, 0, 0.03) 50%,
            transparent 70%
        );
    }

    :where(:global(body.light-mode)) .cert-title {
        color: #18181b;
    }

    :where(:global(body.light-mode)) .cert-verified {
        border-color: rgba(217, 65, 0, 0.12);
        background: rgba(217, 65, 0, 0.04);
        color: var(--light-accent, #d94100);
    }

    :where(:global(body.light-mode)) .cert-chip {
        color: rgba(0, 0, 0, 0.6);
        background: rgba(255, 255, 255, 0.88);
        border-color: rgba(0, 0, 0, 0.1);
    }

    :where(:global(body.light-mode)) .cert-chip-type.accent {
        color: var(--light-accent, #d94100);
        border-color: rgba(217, 65, 0, 0.4);
        background: rgba(217, 65, 0, 0.06);
    }

    :where(:global(body.light-mode)) .cert-chip-idx {
        color: rgba(0, 0, 0, 0.45);
        border-color: rgba(0, 0, 0, 0.08);
    }

    :where(:global(body.light-mode)) .cert-card:hover .cert-chip-idx {
        color: rgba(0, 0, 0, 0.75);
        border-color: rgba(217, 65, 0, 0.4);
    }

    :where(:global(body.light-mode)) .cert-tag {
        border-color: rgba(0, 0, 0, 0.05);
        background: rgba(0, 0, 0, 0.015);
        color: rgba(0, 0, 0, 0.45);
    }

    :where(:global(body.light-mode)) .cert-cat {
        border-color: rgba(0, 0, 0, 0.08);
        background: rgba(0, 0, 0, 0.02);
        color: rgba(0, 0, 0, 0.55);
    }

    :where(:global(body.light-mode)) .cert-note {
        color: rgba(0, 0, 0, 0.45);
    }

    :where(:global(body.light-mode)) .cert-credential {
        color: rgba(0, 0, 0, 0.37);
        border-top-color: rgba(0, 0, 0, 0.04);
    }

    :where(:global(body.light-mode)) .cert-credential:hover {
        color: var(--light-accent, #d94100);
        border-top-color: rgba(217, 65, 0, 0.3);
    }

    :where(:global(body.light-mode)) .cert-image-index {
        color: rgba(217, 65, 0, 0.3);
    }

    :where(:global(body.light-mode)) .cert-card:hover .cert-image-index {
        color: rgba(217, 65, 0, 0.55);
    }

    :where(:global(body.light-mode)) .cert-image-action {
        color: rgba(0, 0, 0, 0.35);
    }

    /* ═══════════════════════════ Reduced Motion ═══════════════════════════ */

    @media (prefers-reduced-motion: reduce) {
        .cert-card {
            animation: none;
        }

        .cert-card,
        .cert-card::after,
        .cert-thumb,
        .cert-chip,
        .cert-image-overlay,
        .cert-image-index,
        .cert-image-action,
        .cert-credential,
        .cert-credential-icon {
            transition: none;
        }
    }
</style>
