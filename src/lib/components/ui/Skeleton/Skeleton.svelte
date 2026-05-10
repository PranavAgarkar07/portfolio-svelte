<script lang="ts">
    import type { HTMLAttributes } from 'svelte/elements';

    type Variant = 'text' | 'card';

    type Props = {
        variant?: Variant;
        width?: string;
        height?: string;
        count?: number;
        class?: string;
        style?: string;
    } & HTMLAttributes<HTMLElement>;

    let {
        variant = 'text',
        width,
        height,
        count = 1,
        class: className = '',
        style = '',
        ...rest
    }: Props = $props();

    let items = $derived(Array.from({ length: count }));
</script>

<div class="skeleton-group{className ? ' ' + className : ''}" {...rest}>
    {#each items as _, i}
        <div
            class="skeleton skeleton-{variant}"
            style={[
                width ? `width: ${width}` : (count > 1 && i === count - 1 ? 'width: 60%' : 'width: 100%'),
                height ? `height: ${height}` : '',
                style
            ].filter(Boolean).join('; ')}
        ></div>
    {/each}
</div>

<style lang="scss">
	@use '../../../styles' as *;

	.skeleton-group {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.skeleton {
		height: 0.65rem;
		background: linear-gradient(90deg, #1a1a1a 25%, #2a2a2a 50%, #1a1a1a 75%);
		background-size: 200% 100%;
		animation: shimmer 1.5s ease-in-out infinite;
		border-radius: 2px;

		&.skeleton-text {
			height: 0.65rem;
		}

		&.skeleton-card {
			height: 120px;
			border-radius: 0;
		}

		&.w-25 { width: 25%; }
		&.w-50 { width: 50%; }
		&.w-75 { width: 75%; }
		&.w-80 { width: 80%; }
		&.w-100 { width: 100%; }
	}

	@keyframes shimmer {
		0% { background-position: 200% 0; }
		100% { background-position: -200% 0; }
	}

	:global(body.light-mode) .skeleton {
		background: linear-gradient(90deg, #e0e0e0 25%, #f0f0f0 50%, #e0e0e0 75%);
		background-size: 200% 100%;
	}

	@media (prefers-reduced-motion: reduce) {
		.skeleton {
			animation: none;
		}
	}
</style>
