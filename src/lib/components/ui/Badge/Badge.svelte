<script lang="ts">
    import type { Snippet } from 'svelte';
    import type { HTMLAttributes } from 'svelte/elements';

    type Variant = 'live' | 'offline' | 'default';

    type Props = {
        variant?: Variant;
        class?: string;
        style?: string;
        children?: Snippet;
    } & HTMLAttributes<HTMLElement>;

    let {
        variant = 'default',
        class: className = '',
        style = '',
        children,
        ...rest
    }: Props = $props();
</script>

<span
    class="badge badge-{variant}{className ? ' ' + className : ''}"
    {style}
    {...rest}
>
    <span class="badge-dot"></span>
    {#if children}
        {@render children()}
    {/if}
</span>

<style lang="scss">
	@use '../../../styles' as *;

	.badge {
		font-family: $font-body;
		font-size: 0.65rem;
		letter-spacing: 1.5px;
		display: inline-flex;
		align-items: center;
		gap: 5px;
		padding: 2px 7px;
		text-transform: uppercase;
		transition: all 0.2s ease;

		&.badge-live {
			color: $accent;
			border: 1px solid rgba(255, 68, 0, 0.4);
			background: rgba(255, 68, 0, 0.08);

			.badge-dot {
				width: 5px;
				height: 5px;
				background: $accent;
				border-radius: 50%;
				box-shadow: 0 0 6px $accent;
				animation: blink 1.2s step-end infinite;
			}
		}

		&.badge-offline {
			color: $text-secondary;
			opacity: 0.5;
			border: 1px solid rgba(255, 255, 255, 0.08);

			.badge-dot {
				display: none;
			}
		}

		&.badge-default {
			color: $text-secondary;
			border: 1px solid $grid-line;

			.badge-dot {
				display: none;
			}
		}
	}

	@keyframes blink {
		50% { opacity: 0; }
	}

	:global(body.light-mode) .badge {
		&.badge-live {
			color: #00aa77;
			border-color: #00aa77;
			background: rgba(0, 170, 119, 0.1);

			.badge-dot {
				background: #00aa77;
				box-shadow: 0 0 6px rgba(0, 170, 119, 0.4);
			}
		}

		&.badge-offline {
			color: $light-text-secondary;
			opacity: 0.6;
			border-color: rgba(0, 0, 0, 0.12);
		}

		&.badge-default {
			color: $light-text-secondary;
			border-color: $light-grid-line;
		}
	}

	@media (prefers-reduced-motion: reduce) {
		.badge.badge-live .badge-dot {
			animation: none;
		}
	}
</style>
