<script lang="ts">
    import type { Snippet } from 'svelte';
    import type { HTMLAttributes } from 'svelte/elements';

    type Variant = 'default' | 'muted';

    type Props = {
        variant?: Variant;
        removable?: boolean;
        class?: string;
        style?: string;
        children?: Snippet;
    } & HTMLAttributes<HTMLElement>;

    let {
        variant = 'default',
        removable = false,
        class: className = '',
        style = '',
        children,
        ...rest
    }: Props = $props();
</script>

<span
    class="tag tag-{variant}{removable ? ' tag-removable' : ''}{className ? ' ' + className : ''}"
    {style}
    {...rest}
>
    {#if children}
        {@render children()}
    {/if}
    {#if removable}
        <span class="tag-remove">&times;</span>
    {/if}
</span>

<style lang="scss">
	@use '../../../styles' as *;

	.tag {
		font-family: $font-body;
		font-size: 0.68rem;
		font-weight: 600;
		color: $accent;
		border: 1px solid $accent;
		background: transparent;
		padding: 3px 8px;
		text-transform: uppercase;
		letter-spacing: 0.5px;
		transition: background 0.1s ease, color 0.1s ease;
		display: inline-flex;
		align-items: center;
		gap: 4px;

		&.tag-muted {
			color: $text-secondary;
			border-color: $text-secondary;
		}

		&.tag-removable {
			cursor: pointer;

			.tag-remove {
				display: inline-flex;
				align-items: center;
				margin-left: 2px;
				opacity: 0.6;
				transition: opacity 0.1s ease;
			}

			&:hover .tag-remove {
				opacity: 1;
			}

			&:hover {
				background: $accent;
				color: #000;
			}
		}
	}

	:global(body.light-mode) .tag {
		color: $light-accent;
		border-color: $light-accent;

		&.tag-muted {
			color: $light-text-secondary;
			border-color: $light-text-secondary;
		}

		&.tag-removable:hover {
			background: $light-accent;
			color: #fff;
		}
	}

	@media (prefers-reduced-motion: reduce) {
		.tag {
			transition: none;
		}
	}
</style>
