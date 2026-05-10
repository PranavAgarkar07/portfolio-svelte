<script lang="ts">
    import { onMount } from 'svelte';
    import type { Snippet } from 'svelte';

    type Props = {
        strength?: number;
        range?: number;
        speed?: number;
        disabled?: boolean;
        class?: string;
        children?: Snippet;
    };

    let {
        strength = 0.25,
        range = 120,
        speed = 0.12,
        disabled = false,
        class: className = '',
        children,
    }: Props = $props();

    let el: HTMLElement | undefined = $state();
    let rafId: number | undefined;
    let mouseX = 0;
    let mouseY = 0;
    let currentX = 0;
    let currentY = 0;
    let isHovering = false;
    let isTouch = false;
    let prefersReduced = false;

    function checkPrefs() {
        isTouch = window.matchMedia('(pointer: coarse)').matches;
        prefersReduced = window.matchMedia('(prefers-reduced-motion: reduce)').matches;
    }

    function handleMove(e: MouseEvent) {
        if (!el || isTouch || disabled || prefersReduced) return;
        const rect = el.getBoundingClientRect();
        const cx = rect.left + rect.width / 2;
        const cy = rect.top + rect.height / 2;
        const dx = e.clientX - cx;
        const dy = e.clientY - cy;
        const dist = Math.hypot(dx, dy);

        if (dist < range) {
            const factor = 1 - dist / range;
            mouseX = dx * factor * strength;
            mouseY = dy * factor * strength;
            if (!isHovering) {
                isHovering = true;
                el.style.transition = 'none';
            }
        } else {
            mouseX = 0;
            mouseY = 0;
        }
    }

    function handleLeave() {
        mouseX = 0;
        mouseY = 0;
        isHovering = false;
        if (el) {
            el.style.transition = `transform ${0.3}s cubic-bezier(0, 0, 0.2, 1)`;
            el.style.transform = 'translate(0, 0)';
        }
    }

    function tick() {
        if (!el || isTouch || disabled || prefersReduced) return;
        currentX += (mouseX - currentX) * speed;
        currentY += (mouseY - currentY) * speed;
        if (Math.abs(currentX) > 0.01 || Math.abs(currentY) > 0.01 || isHovering) {
            el.style.transform = `translate(${currentX}px, ${currentY}px)`;
        }
        rafId = requestAnimationFrame(tick);
    }

    onMount(() => {
        checkPrefs();
        rafId = requestAnimationFrame(tick);
        return () => {
            if (rafId) cancelAnimationFrame(rafId);
        };
    });
</script>

<div
    bind:this={el}
    class="magnetic{className ? ' ' + className : ''}"
    onmousemove={handleMove}
    onmouseleave={handleLeave}
>
    {#if children}
        {@render children()}
    {/if}
</div>

<style lang="scss">
    .magnetic {
        display: inline-flex;
        will-change: transform;
    }
</style>
