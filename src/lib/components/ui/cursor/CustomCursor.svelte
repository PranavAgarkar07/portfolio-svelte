<script lang="ts">
    import { onMount } from 'svelte';

    type Props = {
        size?: number;
        hoverSize?: number;
        color?: string;
        lag?: number;
        ringLag?: number;
        mixBlend?: boolean;
        magnetElements?: string;
        disabled?: boolean;
    };

    let {
        size = 8,
        hoverSize = 48,
        color = 'var(--accent)',
        lag = 0.12,
        ringLag = 0.06,
        mixBlend = false,
        magnetElements = 'a, button, .card-interactive, .tag, .badge, .demo-item, .theme-toggle-btn, .toggle-label',
        disabled = false,
    }: Props = $props();

    let dot: HTMLElement | undefined;
    let ring: HTMLElement | undefined;
    let rafId: number | undefined;
    let isTouch = false;
    let prefersReduced = false;
    let visible = false;

    let mouse = { x: -100, y: -100 };
    let dotPos = { x: -100, y: -100 };
    let ringPos = { x: -100, y: -100 };
    let ringScale = 1;
    let targetScale = 1;
    let isHovering = false;

    function checkPrefs() {
        isTouch = window.matchMedia('(pointer: coarse)').matches;
        prefersReduced = window.matchMedia('(prefers-reduced-motion: reduce)').matches;
    }

    function tick() {
        if (isTouch || disabled || prefersReduced) return;

        dotPos.x += (mouse.x - dotPos.x) * lag;
        dotPos.y += (mouse.y - dotPos.y) * lag;

        ringPos.x += (mouse.x - ringPos.x) * ringLag;
        ringPos.y += (mouse.y - ringPos.y) * ringLag;

        ringScale += (targetScale - ringScale) * 0.1;

        if (dot) {
            dot.style.transform = `translate(${dotPos.x}px, ${dotPos.y}px)`;
        }
        if (ring) {
            ring.style.transform = `translate(${ringPos.x}px, ${ringPos.y}px) scale(${ringScale})`;
        }

        rafId = requestAnimationFrame(tick);
    }

    function handleMove(e: MouseEvent) {
        mouse.x = e.clientX;
        mouse.y = e.clientY;
        if (!visible) {
            visible = true;
            if (dot) dot.style.opacity = '1';
            if (ring) ring.style.opacity = '1';
        }
    }

    function handleLeave() {
        visible = false;
        if (dot) dot.style.opacity = '0';
        if (ring) ring.style.opacity = '0';
    }

    function handleEnter(e: MouseEvent) {
        const target = e.target as HTMLElement;
        if (target.closest(magnetElements)) {
            isHovering = true;
            targetScale = 1.5;
        }
    }

    function handleExit(e: MouseEvent) {
        const target = e.target as HTMLElement;
        if (target.closest(magnetElements)) {
            isHovering = false;
            targetScale = 1;
        }
    }

    onMount(() => {
        checkPrefs();
        if (isTouch || disabled || prefersReduced) return;

        rafId = requestAnimationFrame(tick);

        window.addEventListener('mousemove', handleMove);
        document.addEventListener('mouseleave', handleLeave);
        document.addEventListener('mouseenter', handleEnter, true);
        document.addEventListener('mouseleave', handleExit, true);

        return () => {
            if (rafId) cancelAnimationFrame(rafId);
            window.removeEventListener('mousemove', handleMove);
            document.removeEventListener('mouseleave', handleLeave);
            document.removeEventListener('mouseenter', handleEnter, true);
            document.removeEventListener('mouseleave', handleExit, true);
        };
    });
</script>

<svelte:head>
    <style>
        .cursor-dot,
        .cursor-ring {
            position: fixed;
            pointer-events: none;
            z-index: 99999;
            top: 0;
            left: 0;
            will-change: transform;
            transition: opacity 0.15s ease;
        }

        .cursor-dot {
            width: {size}px;
            height: {size}px;
            background: {color};
            border-radius: 50%;
            margin: -{size / 2}px 0 0 -{size / 2}px;
            opacity: 0;
        }

        .cursor-ring {
            width: {hoverSize}px;
            height: {hoverSize}px;
            border: 2px solid {color};
            border-radius: 50%;
            margin: -{hoverSize / 2}px 0 0 -{hoverSize / 2}px;
            opacity: 0;
            {mixBlend ? 'mix-blend-mode: difference;' : ''}
        }

        @media (hover: none), (prefers-reduced-motion: reduce) {
            .cursor-dot,
            .cursor-ring {
                display: none;
            }
        }
    </style>
</svelte:head>

<div class="cursor-dot" bind:this={dot}></div>
<div class="cursor-ring" bind:this={ring}></div>
