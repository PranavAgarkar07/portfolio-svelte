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
            images?: Image[];
            links: Array<{ label: string; url: string; icon: string }>;
        };
        index: number;
    }

    let { project, index }: Props = $props();
    let currentSlide = $state(0);
    let carousel = $state<HTMLDivElement | null>(null);
    let cardEl = $state<HTMLElement | null>(null);
    let watermarkEl = $state<HTMLElement | null>(null);
    let reducedMotion = $state(false);

    let galleryName = $derived(`project-${index}`);
    let slideCount = $derived(project.images?.length ?? 0);

    onMount(() => {
        reducedMotion = window.matchMedia(
            "(prefers-reduced-motion: reduce)",
        ).matches;
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
    bind:this={cardEl}
    onmousemove={handleMouseMove}
    onmouseleave={handleMouseLeave}
>
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
                            Lightbox.open(img.src, e.currentTarget);
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
                <button
                    class="carousel-btn prev"
                    onclick={prev}
                    aria-label="Previous slide"
                >
                    <Icon name="chevron-left" />
                </button>
                <button
                    class="carousel-btn next"
                    onclick={next}
                    aria-label="Next slide"
                >
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
                <span class="project-number"
                    >PRJ {String(index + 1).padStart(2, "0")}</span
                >
                {#if project.isLive}
                    <Badge variant="live">LIVE</Badge>
                {:else}
                    <Badge variant="offline">OFFLINE</Badge>
                {/if}
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
                <a
                    href={link.url}
                    target="_blank"
                    rel="noopener"
                    class="card-link"
                >
                    <Icon name={link.icon} size={14} />
                    <span>{link.label}</span>
                    <span class="link-arrow-icon"
                        ><Icon name="arrow-up-right" size={12} /></span
                    >
                </a>
            {/each}
        </div>
        <span class="card-number-watermark" bind:this={watermarkEl} aria-hidden="true"
            >0{index + 1}</span
        >
    </div>
</article>

<style>
    .project-card {
        position: relative;
        background: #111922;
        padding: 0;
        /* Hard brutalist borders — no chamfer, raw edges */
        border: 2px solid rgba(255, 255, 255, 0.08);
        border-left: 4px solid rgba(255, 255, 255, 0.12);
        box-shadow: 6px 6px 0px rgba(0, 0, 0, 0.6);
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
        /* Brutalist: hard offset shadow in accent color */
        box-shadow: 6px 6px 0px var(--accent);
        background: #141e2a;
    }

    .card-media {
        position: relative;
        overflow: hidden;
        background: #070a0f;
        border-bottom: 2px solid rgba(255, 255, 255, 0.06);
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
        aspect-ratio: 4 / 3;
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
            aspect-ratio: 4 / 3;
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
        padding: 1.75rem 1.75rem 2rem;
        display: flex;
        flex-direction: column;
        flex: 1;
        gap: 0.9rem;
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
        gap: 0.4rem;
    }

    .card-meta {
        display: flex;
        align-items: center;
        gap: 0.6rem;
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

    :global(.project-card:hover .tag) {
        background: var(--accent);
        color: #000;
    }

    :global(body.light-mode .project-card:hover .tag) {
        background: rgba(200, 50, 0, 0.12);
        color: var(--accent);
    }

    .project-title {
        font-family: var(--font-heading);
        /* Larger, heavier — brutalist punch */
        font-size: 1.5rem;
        font-weight: 700;
        color: #fff;
        line-height: 1.15;
        letter-spacing: 0.02em;
        text-transform: uppercase;
        position: relative;
        z-index: 1;
    }

    .project-description {
        color: var(--text-secondary);
        font-size: 0.875rem;
        line-height: 1.65;
        position: relative;
        z-index: 1;
    }

    .card-tags {
        display: flex;
        flex-wrap: wrap;
        gap: 0.4rem;
        position: relative;
        z-index: 1;
    }

    .card-links {
        display: flex;
        flex-wrap: wrap;
        gap: 0.6rem;
        margin-top: auto;
        padding-top: 1rem;
        border-top: 2px solid rgba(255, 255, 255, 0.06);
        position: relative;
        z-index: 1;
    }

    /* Brutalist link button: hard offset shadow, visible border */
    .card-link {
        display: inline-flex;
        align-items: center;
        gap: 6px;
        font-family: var(--font-heading);
        font-size: 0.7rem;
        font-weight: 700;
        text-transform: uppercase;
        letter-spacing: 0.1em;
        color: #ccc;
        text-decoration: none;
        padding: 7px 12px;
        border: 1.5px solid rgba(255, 255, 255, 0.15);
        background: rgba(255, 255, 255, 0.03);
        /* Hard offset — brutalism signature */
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
        /* Push into shadow on hover */
        box-shadow: 2px 2px 0px var(--accent);
        transform: translate(1px, 1px);
    }

    .card-link:active {
        box-shadow: none;
        transform: translate(3px, 3px);
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
        border-color: rgba(0, 0, 0, 0.1);
        border-left-color: rgba(0, 0, 0, 0.2);
        box-shadow: 6px 6px 0px rgba(0, 0, 0, 0.15);
    }

    :global(body.light-mode) .project-card:hover {
        border-color: var(--accent);
        border-left-color: var(--accent);
        box-shadow: 6px 6px 0px var(--accent);
        background: #dce1e8;
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

    :global(body.light-mode) .card-link {
        box-shadow: 2px 2px 0px rgba(0, 0, 0, 0.12);
    }

    :global(body.light-mode) .card-link:hover {
        box-shadow: 1px 1px 0px rgba(200, 50, 0, 0.25);
    }

    :global(body.light-mode) .card-link:active {
        box-shadow: none;
    }

    @media (min-width: 768px) {
        .card-body {
            padding: 1.5rem 1.5rem 1.75rem;
            gap: 0.6rem;
        }

        .project-title {
            font-size: 1.35rem;
        }

        .project-description {
            font-size: 0.85rem;
            line-height: 1.6;
        }

        .card-tags {
            gap: 0.3rem;
        }

        .card-links {
            padding-top: 0.75rem;
            gap: 0.5rem;
        }

        .card-link {
            font-size: 0.65rem;
            padding: 6px 11px;
        }
    }
</style>
