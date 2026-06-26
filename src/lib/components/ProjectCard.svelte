<script lang="ts">
    import { onMount } from "svelte";
import { Lightbox } from "lightbox3";
    import { Badge, Tag, Icon } from "$lib/components/ui";

    interface Image {
        src: string;
        alt: string;
    }

    interface Props {
        project: {
            name: string;
            description: string;
            tags: string[];
            isLive: boolean;
            featured?: boolean;
            seriesTag?: string;
            githubStars?: number;
            images?: Image[];
            links: Array<{ label: string; url: string; icon: string }>;
        };
        index: number;
        liked?: boolean;
        likeCount?: number;
        onLike?: (liked: boolean) => void;
    }

    let { project, index, liked = false, likeCount = 0, onLike }: Props = $props();
    let currentSlide = $state(0);
    let carousel = $state<HTMLDivElement | null>(null);
    let cardEl = $state<HTMLElement | null>(null);
    let watermarkEl = $state<HTMLElement | null>(null);
    let reducedMotion = $state(false);

    let galleryName = $derived(`project-${index}`);
    let slideCount = $derived(project.images?.length ?? 0);

    function handleKeyDown(e: KeyboardEvent) {
        if (e.key === 'ArrowLeft') prev();
        if (e.key === 'ArrowRight') next();
    }

    onMount(() => {
        reducedMotion = window.matchMedia(
            "(prefers-reduced-motion: reduce)",
        ).matches;
        window.addEventListener('keydown', handleKeyDown);
        return () => window.removeEventListener('keydown', handleKeyDown);
    });

    function handleMouseMove(e: MouseEvent) {
        if (reducedMotion || !cardEl || !watermarkEl) return;
        watermarkEl.style.transition = "none";
        const rect = cardEl.getBoundingClientRect();
        const x = ((e.clientX - rect.left) / rect.width - 0.5) * 2;
        const y = ((e.clientY - rect.top) / rect.height - 0.5) * 2;
        const factor = 10;
        watermarkEl.style.transform = `translate(${x * factor}px, ${y * factor}px)`;
    }

    function handleMouseLeave() {
        if (!watermarkEl) return;
        watermarkEl.style.transition = "";
        watermarkEl.style.transform = "";
    }

    function next() {
        if (slideCount < 2) return;
        currentSlide = (currentSlide + 1) % slideCount;
        scrollToSlide();
    }

    function prev() {
        if (slideCount < 2) return;
        currentSlide = (currentSlide - 1 + slideCount) % slideCount;
        scrollToSlide();
    }

    function goTo(n: number) {
        currentSlide = n;
        scrollToSlide();
    }

    function projectLinkTarget(name: string, label: string): string {
        const n = name.toLowerCase().replace(/[^a-z0-9]/g, '_').replace(/_+$/, '');
        const l = label.toLowerCase().replace(/\s+/g, '_');
        return `${n}_${l}`;
    }

    function scrollToSlide() {
        if (!carousel) return;
        const child = carousel.children[currentSlide] as HTMLElement;
        if (child)
            carousel.scrollTo({
                left: child.offsetLeft,
                behavior: "smooth",
            });
    }
</script>

<article
    class="project-card"
    class:featured={project.featured}
    bind:this={cardEl}
    onmousemove={handleMouseMove}
    onmouseleave={handleMouseLeave}
>
    {#if project.featured}
        <!-- FEATURED LAYOUT: horizontal split — image left, content right -->
        <div class="featured-inner">
            <div class="card-media featured-media">
                {#if project.images && project.images.length > 0}
                    <div class="carousel-track" bind:this={carousel}>
                        {#each project.images as img, i}
                            <a
                                href={img.src}
                                data-lightbox={galleryName}
                                data-caption={img.alt}
                                class="carousel-slide"
                                class:active={i === currentSlide}
                                aria-label="View {img.alt} full size"
                                onclick={(e) => {
                                    e.preventDefault();
                                    Lightbox.instance.open(img.src, e.currentTarget);
                                }}
                            >
                                <img
                                    src={img.src}
                                    alt={img.alt}
                                    class="carousel-img"
                                    loading="lazy"
                                    decoding="async"
                                />
                            </a>
                        {/each}
                    </div>

                    {#if slideCount > 1}
                        <button class="carousel-btn prev" onclick={prev} aria-label="Previous slide">
                            <Icon name="chevron-left" />
                        </button>
                        <button class="carousel-btn next" onclick={next} aria-label="Next slide">
                            <Icon name="chevron-right" />
                        </button>
                        <div class="carousel-dots">
                            {#each project.images as _, i}
                                <button
                                    class="dot"
                                    class:active={i === currentSlide}
                                    onclick={() => goTo(i)}
                                    aria-label="Go to slide {i + 1}"
                                ></button>
                            {/each}
                        </div>
                    {/if}
                {/if}
            </div>

            <div class="card-body featured-body">
                <div class="card-header">
                    <div class="card-meta">
                        <span class="project-number">PRJ {String(index + 1).padStart(2, "0")}</span>
                        <span class="featured-badge" aria-label="Featured project">FEATURED</span>
                        <span class="badge-align">
                            {#if project.isLive}
                                <Badge variant="live">LIVE</Badge>
                            {:else}
                                <Badge variant="offline">OFFLINE</Badge>
                            {/if}
                        </span>
                    </div>
                    <h3 class="project-title featured-title">{project.name}</h3>
                </div>

                <p class="project-description">{project.description}</p>

                <div class="card-tags">
                    {#each project.tags as tag}
                        <Tag>{tag}</Tag>
                    {/each}
                </div>

                <div class="card-links">
                    {#each project.links as link}
                        <a href={link.url} target="_blank" rel="noopener" class="card-link">
                            <Icon name={link.icon} size={14} />
                            <span>{link.label}</span>
                            <span class="link-arrow-icon"><Icon name="arrow-up-right" size={12} /></span>
                        </a>
                    {/each}
                    <button class="like-btn" class:liked class:interactive={!!onLike} onclick={() => onLike?.(!liked)} aria-label={liked ? "Unlike" : "Like this project"}>
                        <svg width="14" height="14" viewBox="0 0 24 24" fill={liked ? "currentColor" : "none"} stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <path d="M19 14c1.49-1.46 3-3.21 3-5.5A5.5 5.5 0 0 0 16.5 3c-1.76 0-3 .5-4.5 2-1.5-1.5-2.74-2-4.5-2A5.5 5.5 0 0 0 2 8.5c0 2.3 1.5 4.05 3 5.5l7 7Z"/>
                        </svg>
                        <span class="like-count">{likeCount}</span>
                    </button>
                </div>
            </div>
        </div>
    {:else}
        <!-- STANDARD LAYOUT: vertical card -->
        <div class="card-media">
            {#if project.images && project.images.length > 0}
                <div class="carousel-track" bind:this={carousel}>
                    {#each project.images as img, i}
                        <a
                            href={img.src}
                            data-lightbox={galleryName}
                            data-caption={img.alt}
                            class="carousel-slide"
                            class:active={i === currentSlide}
                            aria-label="View {img.alt} full size"
                            onclick={(e) => {
                                e.preventDefault();
                                Lightbox.instance.open(img.src, e.currentTarget);
                            }}
                        >
                            <img
                                src={img.src}
                                alt={img.alt}
                                class="carousel-img"
                                loading={i === 0 ? "eager" : "lazy"}
                                decoding="async"
                            />
                        </a>
                    {/each}
                </div>

                {#if slideCount > 1}
                    <button class="carousel-btn prev" onclick={prev} aria-label="Previous slide">
                        <Icon name="chevron-left" />
                    </button>
                    <button class="carousel-btn next" onclick={next} aria-label="Next slide">
                        <Icon name="chevron-right" />
                    </button>
                    <div class="carousel-dots">
                        {#each project.images as _, i}
                            <button
                                class="dot"
                                class:active={i === currentSlide}
                                onclick={() => goTo(i)}
                                aria-label="Go to slide {i + 1}"
                            ></button>
                        {/each}
                    </div>
                {/if}
            {/if}
        </div>

        <div class="card-body">
            <div class="card-header">
                <div class="card-meta">
                    <span class="project-number">PRJ {String(index + 1).padStart(2, "0")}</span>
                    {#if project.seriesTag}
                        <span class="series-tag" title="Part of the {project.seriesTag}">{project.seriesTag.split(' ')[0].toUpperCase()}</span>
                    {/if}
                    <span class="badge-align">
                        {#if project.isLive}
                            <Badge variant="live">LIVE</Badge>
                        {:else}
                            <Badge variant="offline">OFFLINE</Badge>
                        {/if}
                    </span>
                </div>
                <h3 class="project-title">{project.name}</h3>
            </div>

            <p class="project-description">{project.description}</p>

            <div class="card-tags">
                {#each project.tags as tag}
                    <Tag>{tag}</Tag>
                {/each}
            </div>

            <div class="card-links">
                {#each project.links as link}
                    <a href={link.url} target="_blank" rel="noopener" class="card-link">
                        <Icon name={link.icon} size={14} />
                        <span>{link.label}</span>
                        <span class="link-arrow-icon"><Icon name="arrow-up-right" size={12} /></span>
                    </a>
                {/each}
                <button class="like-btn" class:liked class:interactive={!!onLike} onclick={() => onLike?.(!liked)} aria-label={liked ? "Unlike" : "Like this project"}>
                    <svg width="14" height="14" viewBox="0 0 24 24" fill={liked ? "currentColor" : "none"} stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <path d="M19 14c1.49-1.46 3-3.21 3-5.5A5.5 5.5 0 0 0 16.5 3c-1.76 0-3 .5-4.5 2-1.5-1.5-2.74-2-4.5-2A5.5 5.5 0 0 0 2 8.5c0 2.3 1.5 4.05 3 5.5l7 7Z"/>
                    </svg>
                    <span class="like-count">{likeCount}</span>
                </button>
            </div>
            <span class="card-number-watermark" bind:this={watermarkEl} aria-hidden="true"
                >0{index + 1}</span
            >
        </div>
    {/if}
</article>

<style>
    .project-card {
        position: relative;
        background: #111922;
        padding: 0;
        border: 1.5px solid rgba(255, 255, 255, 0.08);
        border-left: 4px solid rgba(255, 255, 255, 0.12);
        box-shadow: 4px 4px 0px rgba(0, 0, 0, 0.6);
        transition:
            border-color 0.15s ease,
            box-shadow 0.15s ease,
            background 0.15s ease;
        display: flex;
        flex-direction: column;
        height: 100%;
    }

    .project-card:hover {
        border-color: var(--accent);
        border-left-color: var(--accent);
        box-shadow: 4px 4px 0px var(--accent);
        background: #141e2a;
    }

    .project-card:focus-visible {
        outline: 3px solid var(--accent);
        outline-offset: 4px;
    }

    /* Featured card: same size as standard, accent border only */
    .project-card.featured {
        grid-column: auto;
        border-color: rgba(255, 68, 0, 0.2);
        border-left-color: var(--accent);
    }

    .project-card.featured:hover {
        border-color: var(--accent);
    }

    .featured-inner {
        display: grid;
        grid-template-columns: 1.2fr 1fr;
        height: 100%;
    }

    .featured-media {
        border-bottom: none;
        height: 100%;
    }

    .featured-media .carousel-slide {
        aspect-ratio: auto;
        height: 100%;
    }

    .featured-body {
        padding: 1.25rem 1.25rem 1.5rem;
        gap: 0.6rem;
    }

    .featured-title {
        font-size: clamp(1.15rem, 2.5vw, 1.4rem);
    }

    .featured-badge {
        font-family: var(--font-body);
        font-size: 0.55rem;
        font-weight: 600;
        letter-spacing: 0.12em;
        text-transform: uppercase;
        color: var(--accent);
        border: 1px solid var(--accent);
        padding: 2px 7px;
        line-height: 1.6;
        display: inline-flex;
        align-items: center;
    }

    /* Series tag for multi-project ecosystems */
    .series-tag {
        font-family: var(--font-body);
        font-size: 0.55rem;
        font-weight: 600;
        letter-spacing: 0.12em;
        text-transform: uppercase;
        color: var(--text-muted);
        border: 1px dashed rgba(255,255,255,0.15);
        padding: 2px 6px;
        line-height: 1.6;
    }

    :global(body.light-mode) .series-tag {
        border-color: rgba(0,0,0,0.15);
        color: var(--text-muted);
    }

    :global(.project-card:hover .tag) {
        background: rgba(255, 68, 0, 0.15);
        color: var(--accent);
    }

    :global(:global(body.light-mode) .project-card:hover .tag) {
        background: rgba(200, 50, 0, 0.1);
        color: var(--accent);
    }

    .card-media {
        position: relative;
        overflow: hidden;
        background: #070a0f;
        border-bottom: 1.5px solid rgba(255, 255, 255, 0.06);
    }

    .project-card:hover .card-media {
        border-bottom-color: rgba(255, 68, 0, 0.2);
    }

    .carousel-track {
        display: flex;
        overflow-x: auto;
        scroll-snap-type: x mandatory;
        scrollbar-width: none;
        -ms-overflow-style: none;
    }

    .carousel-track::-webkit-scrollbar {
        display: none;
    }

    .carousel-slide {
        flex: 0 0 100%;
        scroll-snap-align: start;
        cursor: pointer;
        overflow: hidden;
        aspect-ratio: 16 / 10;
        height: auto;
        display: block;
        text-decoration: none;
        color: inherit;
    }

    .carousel-img {
        display: block;
        width: 100%;
        height: 100%;
        object-fit: cover;
        transition:
            transform 0.4s ease,
            filter 0.4s ease;
    }

    .project-card:hover .carousel-img {
        transform: scale(1.06);
    }

    @media (max-width: 767px) {
        .carousel-slide {
            aspect-ratio: 16 / 10;
            height: auto;
        }
    }

    .carousel-btn {
        position: absolute;
        top: 50%;
        transform: translateY(-50%);
        background: rgba(0, 0, 0, 0.75);
        border: 1px solid rgba(255, 255, 255, 0.08);
        color: #fff;
        width: 44px;
        height: 44px;
        cursor: pointer;
        z-index: 10;
        display: flex;
        align-items: center;
        justify-content: center;
        transition: all 0.2s ease;
        opacity: 0;
        backdrop-filter: blur(4px);
    }

    @media (hover: hover) {
        .carousel-btn {
            opacity: 0;
        }
        .project-card:hover .carousel-btn {
            opacity: 0.85;
        }
        .carousel-btn:hover {
            opacity: 1 !important;
            background: var(--accent);
            border-color: var(--accent);
            color: #000;
        }
    }

    @media (hover: none) {
        .carousel-btn {
            opacity: 0.85;
        }
    }

    .carousel-btn.prev {
        left: 12px;
    }

    .carousel-btn.next {
        right: 12px;
    }

    .carousel-dots {
        position: absolute;
        bottom: 14px;
        left: 50%;
        transform: translateX(-50%);
        display: flex;
        gap: 8px;
        z-index: 10;
        background: rgba(0, 0, 0, 0.5);
        padding: 6px 10px;
        backdrop-filter: blur(4px);
        border: 1px solid rgba(255, 255, 255, 0.06);
    }

    .dot {
        width: 10px;
        height: 10px;
        border: 1px solid rgba(255, 255, 255, 0.3);
        background: transparent;
        cursor: pointer;
        padding: 0;
        transition: all 0.2s ease;
    }

    .dot.active {
        background: var(--accent);
        border-color: var(--accent);
    }

    .dot:hover {
        border-color: var(--accent);
    }

    @media (min-width: 768px) {
        .carousel-dots {
            bottom: 20px;
            padding: 8px 12px;
        }
        .dot {
            width: 12px;
            height: 12px;
        }
    }

    :global(body.light-mode) .dot {
        border-color: rgba(0, 0, 0, 0.3);
    }

    :global(body.light-mode) .dot.active {
        background: var(--accent);
        border-color: var(--accent);
    }

    :global(body.light-mode) .carousel-btn {
        background: rgba(255, 255, 255, 0.85);
        color: #111;
    }

    .card-body {
        position: relative;
        padding: 1.25rem 1.25rem 1.5rem;
        display: flex;
        flex-direction: column;
        flex: 1;
        gap: 0.6rem;
        min-height: 0;
    }

    .card-number-watermark {
        position: absolute;
        bottom: -0.5rem;
        right: -0.5rem;
        font-family: var(--font-heading);
        font-size: clamp(8rem, 12vw, 12rem);
        font-weight: 900;
        line-height: 1;
        letter-spacing: -0.06em;
        pointer-events: none;
        user-select: none;
        color: transparent;
        -webkit-text-stroke: 2px var(--accent);
        opacity: 0.3;
        z-index: 0;
        filter: drop-shadow(0 4px 6px rgba(0, 0, 0, 0.5));
        transition:
            transform 0.3s ease,
            opacity 0.2s ease,
            filter 0.2s ease;
    }

    .project-card:hover .card-number-watermark {
        opacity: 0.5;
        filter: drop-shadow(0 6px 12px rgba(0, 0, 0, 0.6));
    }

    @media (max-width: 767px) {
        .card-number-watermark {
            bottom: -0.4rem;
            right: -0.3rem;
            font-size: clamp(4rem, 12vw, 8rem);
            -webkit-text-stroke: 1.5px var(--accent);
            opacity: 0.35;
            filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.4));
        }

        .project-card:hover .card-number-watermark {
            filter: drop-shadow(0 4px 8px rgba(0, 0, 0, 0.5));
        }
    }

    :global(body.light-mode) .card-number-watermark {
        -webkit-text-stroke: 2px var(--accent);
        filter: drop-shadow(0 4px 6px rgba(0, 0, 0, 0.15));
        opacity: 0.2;
    }

    :global(body.light-mode) .project-card:hover .card-number-watermark {
        opacity: 0.4;
        filter: drop-shadow(0 6px 12px rgba(0, 0, 0, 0.2));
    }

    .card-header {
        position: relative;
        z-index: 1;
        display: flex;
        flex-direction: column;
        gap: 0.3rem;
    }

    .card-meta {
        display: flex;
        align-items: center;
        gap: 0.5rem;
    }

    .badge-align {
        margin-left: auto;
    }

    .project-number {
        font-family: var(--font-body);
        font-size: 0.65rem;
        color: var(--text-secondary);
        letter-spacing: 3px;
        text-transform: uppercase;
    }

    .project-number::before {
        content: "[";
        color: var(--accent);
        opacity: 0.7;
        margin-right: 1px;
    }

    .project-number::after {
        content: "]";
        color: var(--accent);
        opacity: 0.7;
        margin-left: 1px;
    }

    .project-title {
        font-family: var(--font-heading);
        font-size: clamp(1.05rem, 2vw, 1.25rem);
        font-weight: 700;
        color: #fff;
        line-height: 1.2;
        letter-spacing: 0.02em;
        text-transform: uppercase;
        position: relative;
        z-index: 1;
    }

    .project-description {
        color: var(--text-secondary);
        font-size: 0.8rem;
        line-height: 1.5;
        position: relative;
        z-index: 1;
    }

    .card-tags {
        display: flex;
        flex-wrap: wrap;
        gap: 0.35rem;
        position: relative;
        z-index: 1;
    }

    .card-links {
        display: flex;
        flex-wrap: wrap;
        gap: clamp(0.375rem, 1vw, 0.5rem);
        margin-top: auto;
        padding-top: clamp(0.5rem, 1.25vw, 0.75rem);
        border-top: 1.5px solid rgba(255, 255, 255, 0.06);
        position: relative;
        z-index: 1;
    }

    .card-link {
        display: inline-flex;
        align-items: center;
        gap: 6px;
        font-family: var(--font-heading);
        font-size: 0.65rem;
        font-weight: 700;
        text-transform: uppercase;
        letter-spacing: 0.08em;
        color: #ccc;
        text-decoration: none;
        padding: 5px 10px;
        border: 1.5px solid rgba(255, 255, 255, 0.15);
        background: rgba(255, 255, 255, 0.03);
        box-shadow: 3px 3px 0px rgba(0, 0, 0, 0.5);
        transition:
            color 0.1s ease,
            border-color 0.1s ease,
            background 0.1s ease,
            box-shadow 0.1s ease,
            transform 0.1s ease;
        cursor: pointer;
    }

    .card-link:hover {
        color: #fff;
        border-color: var(--accent);
        background: rgba(255, 68, 0, 0.1);
        box-shadow: 2px 2px 0px var(--accent);
        transform: translate(1px, 1px);
    }

    .card-link:active {
        box-shadow: none;
        transform: translate(3px, 3px);
    }

    .card-link:focus-visible {
        outline: 2px solid var(--accent);
        outline-offset: 2px;
    }

    .link-arrow-icon {
        transition: transform 0.15s ease;
        opacity: 0.6;
    }

    .card-link:hover .link-arrow-icon {
        transform: translate(3px, -3px);
        opacity: 1;
    }

    :global(body.light-mode) .project-card {
        background: #e4e8ed;
        border-color: rgba(0, 0, 0, 0.08);
        border-left-color: rgba(0, 0, 0, 0.2);
        box-shadow: 4px 4px 0px rgba(0, 0, 0, 0.1);
    }

    :global(body.light-mode) .project-card:hover {
        border-color: var(--accent);
        border-left-color: var(--accent);
        box-shadow: 4px 4px 0px var(--accent);
        background: #dce1e8;
    }

    :global(body.light-mode) .project-card.featured {
        border-color: rgba(200, 50, 0, 0.2);
        border-left-color: var(--accent);
    }

    :global(body.light-mode) .project-title {
        color: #111;
    }

    :global(body.light-mode) .card-link {
        color: #333;
        border-color: rgba(0, 0, 0, 0.2);
        background: rgba(0, 0, 0, 0.03);
        box-shadow: 3px 3px 0px rgba(0, 0, 0, 0.2);
    }

    :global(body.light-mode) .card-link:hover {
        color: var(--accent);
        border-color: rgba(200, 50, 0, 0.25);
        background: rgba(200, 50, 0, 0.05);
        box-shadow: 1px 1px 0px rgba(200, 50, 0, 0.25);
    }

    :global(body.light-mode) .card-link:active {
        box-shadow: none;
    }

    /* ── Like Button ── */

    .like-btn {
        display: inline-flex;
        align-items: center;
        gap: 4px;
        font-family: var(--font-heading);
        font-size: 0.65rem;
        font-weight: 700;
        text-transform: uppercase;
        letter-spacing: 0.08em;
        color: rgba(255, 255, 255, 0.25);
        text-decoration: none;
        padding: 5px 10px;
        border: 1.5px solid rgba(255, 255, 255, 0.15);
        background: rgba(255, 255, 255, 0.03);
        box-shadow: 3px 3px 0px rgba(0, 0, 0, 0.5);
        cursor: default;
        transition: all 0.15s ease;
        margin-left: auto;
    }

    .like-btn.interactive {
        cursor: pointer;
    }

    .like-btn.interactive:hover {
        color: #ef4444;
        border-color: #ef4444;
        background: rgba(239, 68, 68, 0.1);
        box-shadow: 2px 2px 0px #ef4444;
        transform: translate(1px, 1px);
    }

    .like-btn.interactive:active {
        box-shadow: none;
        transform: translate(3px, 3px);
    }

    .like-btn.liked {
        color: #ef4444;
        border-color: rgba(239, 68, 68, 0.4);
        background: rgba(239, 68, 68, 0.08);
    }

    .like-btn.liked.interactive:hover {
        border-color: #ef4444;
        background: rgba(239, 68, 68, 0.15);
    }

    .like-count {
        font-size: 0.6rem;
        opacity: 0.8;
    }

    :global(body.light-mode) .like-btn {
        color: rgba(0, 0, 0, 0.25);
        border-color: rgba(0, 0, 0, 0.2);
        background: rgba(0, 0, 0, 0.03);
        box-shadow: 3px 3px 0px rgba(0, 0, 0, 0.2);
    }

    :global(body.light-mode) .like-btn.interactive:hover {
        color: #dc2626;
        border-color: rgba(220, 38, 38, 0.3);
        background: rgba(220, 38, 38, 0.06);
        box-shadow: 1px 1px 0px rgba(220, 38, 38, 0.25);
    }

    :global(body.light-mode) .like-btn.liked {
        color: #dc2626;
        border-color: rgba(220, 38, 38, 0.3);
        background: rgba(220, 38, 38, 0.06);
    }

    :global(body.light-mode) .like-btn.liked.interactive:hover {
        border-color: #dc2626;
        background: rgba(220, 38, 38, 0.12);
    }

    :global(body.light-mode) .like-btn.interactive:active {
        box-shadow: none;
    }

    :global(body.light-mode) .card-links {
        border-top-color: rgba(0, 0, 0, 0.06);
    }

    :global(body.light-mode) .card-media {
        background: #d4d4d8;
    }

    :global(body.light-mode) .carousel-dots {
        background: rgba(255, 255, 255, 0.7);
        border-color: rgba(0, 0, 0, 0.1);
    }
</style>
