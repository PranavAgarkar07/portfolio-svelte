<script lang="ts">
    import type { HTMLAttributes } from 'svelte/elements';

    type Props = {
        value: string;
        group?: string;
        label?: string;
        disabled?: boolean;
        class?: string;
        style?: string;
    } & HTMLAttributes<HTMLElement>;

    let {
        value,
        group,
        label,
        disabled = false,
        class: className = '',
        style = '',
        ...rest
    }: Props = $props();

    let inputId = $state(`radio-${Math.random().toString(36).slice(2, 8)}`);
</script>

<label
    class="radio-wrapper{className ? ' ' + className : ''}"
    class:disabled
    {style}
    for={inputId}
    data-radio-value={value}
    data-radio-group={group}
    {...rest}
>
    <input
        id={inputId}
        type="radio"
        {value}
        {disabled}
        aria-label={label || value}
    />
    <span class="radio-circle" tabindex="0">
        <span class="radio-dot"></span>
    </span>
    {#if label}
        <span class="radio-label">{label}</span>
    {/if}
</label>

<style lang="scss">
	@use '../../../styles' as *;

	.radio-wrapper {
		display: inline-flex;
		align-items: center;
		gap: 0.75rem;
		cursor: pointer;
		font-family: $font-body;
		font-size: 0.9rem;
		color: $text-primary;
		user-select: none;
		position: relative;

		&.disabled {
			opacity: 0.5;
			cursor: not-allowed;
		}

		input {
			position: absolute;
			opacity: 0;
			width: 0;
			height: 0;
			pointer-events: none;
		}
	}

	.radio-circle {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 20px;
		height: 20px;
		border: 2px solid $text-secondary;
		background: transparent;
		flex-shrink: 0;
		transition: all 0.1s ease;
		position: relative;

		.radio-dot {
			width: 10px;
			height: 10px;
			background: $accent;
			transform: scale(0);
			transition: transform 0.15s cubic-bezier(0.16, 1, 0.3, 1);
		}

		&.selected {
			border-color: $accent;
			box-shadow: 2px 2px 0px rgba(255, 68, 0, 0.2);

			.radio-dot {
				transform: scale(1);
			}
		}

		&:focus-visible {
			outline: 3px solid $accent;
			outline-offset: 4px;
		}
	}

	.radio-wrapper:hover .radio-circle:not(.selected) {
		border-color: $accent;
	}

	:global(body.light-mode) {
		.radio-wrapper {
			color: #111;
		}

		.radio-circle {
			border-color: #000;

			&.selected {
				border-color: $light-accent;
				box-shadow: 2px 2px 0px rgba(217, 61, 0, 0.2);

				.radio-dot {
					background: $light-accent;
				}
			}
		}

		.radio-wrapper:hover .radio-circle:not(.selected) {
			border-color: $light-accent;
		}
	}

	@media (prefers-reduced-motion: reduce) {
		.radio-circle {
			transition: none;
			.radio-dot { transition: none; }
		}
	}
</style>
