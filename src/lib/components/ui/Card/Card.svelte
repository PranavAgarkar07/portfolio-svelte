<script lang="ts">
    import { motion } from '@humanspeak/svelte-motion';
    import { onMount } from 'svelte';
    import type { Snippet } from 'svelte';
    import type { HTMLAttributes } from 'svelte/elements';

    type Props = {
        interactive?: boolean;
        animate?: boolean;
        badge?: string | number;
        class?: string;
        style?: string;
        children?: Snippet;
    } & HTMLAttributes<HTMLElement>;

    let {
        interactive = false,
        animate = false,
        badge,
        class: className = '',
        style = '',
        children,
        ...rest
    }: Props = $props();

    let badgeDisplay = $derived(badge !== undefined ? String(badge).padStart(2, '0') : '');
    let badgeChars = $derived(badgeDisplay.split(''));

    let cardEl: any = $state();
    let badgeEl: HTMLElement | undefined = $state();
    let rafId: number | undefined;
    let targetX = 0;
    let targetY = 0;
    let currentX = 0;
    let currentY = 0;
    let isHovering = false;
    let isTouch = false;
    let prefersReduced = false;

    let isInside = $derived(isHovering);

    function checkPrefs() {
        if (typeof window === 'undefined') return;
        isTouch = window.matchMedia('(pointer: coarse)').matches;
        prefersReduced = window.matchMedia('(prefers-reduced-motion: reduce)').matches;
    }

    function handleMouseMove(e: MouseEvent) {
        if (!cardEl || isTouch || prefersReduced) return;
        const rect = cardEl.getBoundingClientRect();
        const cx = rect.left + rect.width / 2;
        const cy = rect.top + rect.height / 2;
        const dx = (e.clientX - cx) / (rect.width / 2);
        const dy = (e.clientY - cy) / (rect.height / 2);
        targetX = dx * 24;
        targetY = dy * 24;
        isHovering = true;
    }

    function handleMouseLeave() {
        targetX = 0;
        targetY = 0;
        isHovering = false;
    }

    function tick() {
        if (!badgeEl || isTouch || prefersReduced) return;
        currentX += (targetX - currentX) * 0.1;
        currentY += (targetY - currentY) * 0.1;
        if (Math.abs(currentX) > 0.05 || Math.abs(currentY) > 0.05 || isInside) {
            badgeEl.style.transform = `translate(${currentX}px, ${currentY}px)`;
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

{#if animate}
    <motion.article
        bind:this={cardEl}
        class="aero-card{interactive ? ' card-interactive' : ''}{className ? ' ' + className : ''}"
        {style}
        whileHover={interactive ? { y: -4, transition: { duration: 0.2 } } : undefined}
        onmousemove={handleMouseMove}
        onmouseleave={handleMouseLeave}
        {...rest}
    >
        {#if children}
            {@render children()}
        {/if}
        {#if badge !== undefined}
            <span bind:this={badgeEl} class="card-badge" aria-hidden="true">
                {#each badgeChars as char, i}
                    <span class="badge-char" style="--i: {i}" data-i={i}>{char}</span>
                {/each}
            </span>
        {/if}
    </motion.article>
{:else}
    <article
        bind:this={cardEl}
        class="aero-card{interactive ? ' card-interactive' : ''}{className ? ' ' + className : ''}"
        {style}
        onmousemove={handleMouseMove}
        onmouseleave={handleMouseLeave}
        {...rest}
    >
        {#if children}
            {@render children()}
        {/if}
        {#if badge !== undefined}
            <span bind:this={badgeEl} class="card-badge" aria-hidden="true">
                {#each badgeChars as char, i}
                    <span class="badge-char" style="--i: {i}" data-i={i}>{char}</span>
                {/each}
            </span>
        {/if}
    </article>
{/if}

<style lang="scss">
	@use '../../../styles' as *;

	.aero-card {
		background-color: #171e28;
		background-image: linear-gradient(160deg, #252e3a 0%, #171e28 100%);
		padding: $card-padding;
		border: 1px solid rgba(255, 255, 255, 0.1);
		position: relative;
		clip-path: polygon(20px 0, 100% 0,
				100% calc(100% - 20px), calc(100% - 20px) 100%,
				0 100%, 0 20px);
		transition: box-shadow 0.2s ease, transform 0.2s ease;
		overflow: hidden;
		will-change: clip-path;
		box-shadow: 0 10px 30px -10px rgba(0, 0, 0, 0.5);
		border-left: 3px solid rgba(255, 255, 255, 0.1);

		&.card-interactive {
			cursor: pointer;

			&:hover {
				border-color: $accent;
				border-left-color: $accent;
				box-shadow: 0 0 25px rgba(255, 51, 0, 0.2);
				background-color: #263340;
				transform: translateY(-2px);
			}
		}

		&:focus-visible {
			outline: 3px solid $accent;
			outline-offset: 4px;
		}
	}

	.card-badge {
		position: absolute;
		bottom: -0.5rem;
		right: -0.5rem;
		font-family: $font-heading;
		font-size: clamp(6rem, 10vw, 10rem);
		font-weight: 900;
		line-height: 1;
		letter-spacing: -0.06em;
		pointer-events: none;
		user-select: none;
		color: transparent;
		-webkit-text-stroke: 2px $accent;
		opacity: 0.3;
		z-index: 0;
		filter: drop-shadow(0 4px 6px rgba(0, 0, 0, 0.5));
		white-space: nowrap;
		will-change: transform;

		.card-interactive:hover & {
			opacity: 0.5;
			filter: drop-shadow(0 6px 12px rgba(0, 0, 0, 0.6));
		}
	}

	.badge-char {
		display: inline-block;
		color: transparent;
		-webkit-text-stroke: 2px $accent;
		will-change: transform, filter, clip-path;
		clip-path: polygon(0 0, 0 0, 0 0, 0 0);
		animation: drawChar 0.55s cubic-bezier(0.22, 1, 0.36, 1) forwards;
		animation-delay: calc(0.13s * var(--i) + 0.12s);
	}

	@keyframes drawChar {
		0% {
			clip-path: polygon(0 0, 0 0, 0 100%, 0 100%);
			opacity: 0;
			transform: translateY(8px) scale(1.08);
			filter: drop-shadow(0 0 14px rgba(255, 68, 0, 0.5));
		}
		20% {
			opacity: 1;
		}
		50% {
			transform: translateY(0) scale(1);
			filter: drop-shadow(0 0 5px rgba(255, 68, 0, 0.15));
		}
		100% {
			clip-path: polygon(0 0, 100% 0, 100% 100%, 0 100%);
			opacity: 1;
			transform: translateY(0) scale(1);
			filter: drop-shadow(0 4px 6px rgba(0, 0, 0, 0.5));
		}
	}

	@media (prefers-reduced-motion: reduce) {
		.badge-char {
			animation: none !important;
			clip-path: none !important;
			-webkit-mask: none !important;
			mask: none !important;
			opacity: 1;
			transform: none;
			filter: drop-shadow(0 4px 6px rgba(0, 0, 0, 0.5));
		}
	}

	:global(body.light-mode) .card-badge {
		-webkit-text-stroke: 2px $light-accent;
		filter: drop-shadow(0 4px 6px rgba(0, 0, 0, 0.15));
		opacity: 0.2;
	}

	:global(body.light-mode) .aero-card {
		background-color: var(--surface-dark);
		background-image: linear-gradient(160deg, #f5f5f0 0%, #edede9 100%);
		border: 1px solid var(--wire-color);
		border-left: 3px solid #c8cbcf;
		box-shadow: 0 1px 3px rgba(0,0,0,0.04);

		&.card-interactive:hover {
			border-left-color: $light-accent;
			background: #fff !important;
			background-color: #fff !important;
			box-shadow: 0 4px 16px rgba(0,0,0,0.06);
			transform: translateY(-2px);
		}
	}

	@media (prefers-reduced-motion: reduce) {
		.aero-card {
			transition: none;
		}
	}
</style>
