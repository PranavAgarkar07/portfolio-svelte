<script lang="ts">
    import type { HTMLAttributes } from 'svelte/elements';

    type Props = {
        checked?: boolean;
        indeterminate?: boolean;
        disabled?: boolean;
        label?: string;
        class?: string;
        style?: string;
        onchange?: (checked: boolean) => void;
    } & HTMLAttributes<HTMLElement>;

    let {
        checked = $bindable(false),
        indeterminate = false,
        disabled = false,
        label,
        class: className = '',
        style = '',
        onchange,
        ...rest
    }: Props = $props();

    let inputId = $state(`checkbox-${Math.random().toString(36).slice(2, 8)}`);

    function handleChange() {
        if (disabled) return;
        checked = !checked;
        indeterminate = false;
        onchange?.(checked);
    }

    function onKeydown(e: KeyboardEvent) {
        if (disabled) return;
        if (e.key === 'Enter' || e.key === ' ') {
            e.preventDefault();
            handleChange();
        }
    }
</script>

<label
    class="checkbox-wrapper{className ? ' ' + className : ''}"
    class:disabled
    {style}
    for={inputId}
    {...rest}
>
    <input
        id={inputId}
        type="checkbox"
        {checked}
        {disabled}
        onchange={handleChange}
        aria-checked={indeterminate ? 'mixed' : checked}
    />
    <span
        class="checkbox-box"
        class:checked={checked && !indeterminate}
        class:indeterminate={indeterminate}
        tabindex="0"
        onkeydown={onKeydown}
        role="presentation"
    >
        {#if checked && !indeterminate}
            <svg class="checkbox-check" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="20 6 9 17 4 12" />
            </svg>
        {/if}
    </span>
    {#if label}
        <span class="checkbox-label">{label}</span>
    {/if}
</label>

<style lang="scss">
	@use '../../../styles' as *;

	.checkbox-wrapper {
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

	.checkbox-box {
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

		.checkbox-check {
			opacity: 0;
			transform: scale(0);
			transition: all 0.15s cubic-bezier(0.16, 1, 0.3, 1);
			width: 14px;
			height: 14px;
			color: $accent;
		}

		&.checked {
			border-color: $accent;
			background: $accent;

			.checkbox-check {
				opacity: 1;
				transform: scale(1);
				color: #000;
			}
		}

		&.indeterminate {
			border-color: $accent;

			&::after {
				content: '';
				width: 10px;
				height: 2px;
				background: $accent;
				position: absolute;
			}

			.checkbox-check {
				display: none;
			}
		}

		&:focus-visible {
			outline: 3px solid $accent;
			outline-offset: 4px;
		}
	}

	.checkbox-wrapper:hover .checkbox-box:not(.checked):not(.indeterminate) {
		border-color: $accent;
	}

	:global(body.light-mode) {
		.checkbox-wrapper {
			color: #111;
		}

		.checkbox-box {
			border-color: #000;

			&.checked {
				border-color: $light-accent;
				background: $light-accent;

				.checkbox-check {
					color: #fff;
				}
			}

			&.indeterminate {
				border-color: $light-accent;

				&::after {
					background: $light-accent;
				}
			}
		}

		.checkbox-wrapper:hover .checkbox-box:not(.checked):not(.indeterminate) {
			border-color: $light-accent;
		}
	}

	@media (prefers-reduced-motion: reduce) {
		.checkbox-box {
			transition: none;

			.checkbox-check {
				transition: none;
			}
		}
	}
</style>
