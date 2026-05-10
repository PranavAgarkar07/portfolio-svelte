<script lang="ts">
    import { motion } from '@humanspeak/svelte-motion';
    import type { HTMLAttributes } from 'svelte/elements';
    import Icon from '../Icon/Icon.svelte';
    import type { IconName } from '../Icon/Icon.svelte';

    type InputType = 'text' | 'email' | 'tel' | 'number' | 'password' | 'search' | 'url' | 'date';

    type Props = {
        label?: string;
        type?: InputType;
        placeholder?: string;
        value?: string | number;
        error?: string;
        hint?: string;
        prefix?: IconName;
        clearable?: boolean;
        animate?: boolean;
        class?: string;
        style?: string;
        onchange?: (value: string) => void;
    } & HTMLAttributes<HTMLInputElement>;

    let {
        label,
        type = 'text',
        placeholder,
        value = $bindable(''),
        error,
        hint,
        prefix,
        clearable = false,
        animate = false,
        class: className = '',
        style = '',
        onchange,
        ...rest
    }: Props = $props();

    let inputId = $state(`input-${Math.random().toString(36).slice(2, 8)}`);
    let showPassword = $state(false);
    let inputType = $derived(type === 'password' && showPassword ? 'text' : type);
    let hasValue = $derived(value !== '' && value !== undefined && value !== null);

    let autoPrefix = $derived<IconName | undefined>(
        prefix ?? (type === 'email' ? 'envelope' : type === 'search' ? 'search' : type === 'date' ? 'calendar' : undefined)
    );

    function handleInput(e: Event) {
        const target = e.target as HTMLInputElement;
        value = target.value;
        onchange?.(target.value);
    }

    function clear() {
        value = '';
        onchange?.('');
    }
</script>

<div class="form-group{className ? ' ' + className : ''}">
    {#if label}
        <label for={inputId} class="form-label">{label}</label>
    {/if}
    {#if animate}
        <motion.div
            class="input-wrap"
            whileTap={{ scale: 0.998 }}
        >
            {#if autoPrefix}
                <span class="input-prefix">
                    <Icon name={autoPrefix} size={18} />
                </span>
            {/if}
            <input
                id={inputId}
                type={inputType}
                class="form-input{error ? ' has-error' : ''}{autoPrefix ? ' has-prefix' : ''}{(clearable && hasValue) || type === 'password' ? ' has-suffix' : ''}"
                {placeholder}
                {value}
                oninput={handleInput}
                autocomplete={type === 'email' ? 'email' : type === 'tel' ? 'tel' : type === 'url' ? 'url' : undefined}
                {...rest}
            />
            {#if type === 'password'}
                <span class="input-suffix" onclick={() => { showPassword = !showPassword; }} aria-label={showPassword ? 'Hide password' : 'Show password'}>
                    <Icon name={showPassword ? 'eye' : 'eye-off'} size={18} />
                </span>
            {:else if clearable && hasValue}
                <span class="input-suffix" onclick={clear} aria-label="Clear input">
                    <Icon name="close" size={18} />
                </span>
            {/if}
            <div class="input-bar"></div>
        </motion.div>
    {:else}
        <div class="input-wrap">
            {#if autoPrefix}
                <span class="input-prefix">
                    <Icon name={autoPrefix} size={18} />
                </span>
            {/if}
            <input
                id={inputId}
                type={inputType}
                class="form-input{error ? ' has-error' : ''}{autoPrefix ? ' has-prefix' : ''}{(clearable && hasValue) || type === 'password' ? ' has-suffix' : ''}"
                {placeholder}
                {value}
                oninput={handleInput}
                autocomplete={type === 'email' ? 'email' : type === 'tel' ? 'tel' : type === 'url' ? 'url' : undefined}
                {...rest}
            />
            {#if type === 'password'}
                <span class="input-suffix" onclick={() => { showPassword = !showPassword; }} aria-label={showPassword ? 'Hide password' : 'Show password'}>
                    <Icon name={showPassword ? 'eye' : 'eye-off'} size={18} />
                </span>
            {:else if clearable && hasValue}
                <span class="input-suffix" onclick={clear} aria-label="Clear input">
                    <Icon name="close" size={18} />
                </span>
            {/if}
            <div class="input-bar"></div>
        </div>
    {/if}
    {#if hint && !error}
        <span class="form-hint">{hint}</span>
    {/if}
    {#if error}
        <span class="form-error-text">{error}</span>
    {/if}
</div>

<style lang="scss">
	@use '../../../styles' as *;

	.form-group {
		margin: 0;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;

		.form-label {
			font-family: $font-body;
			font-size: 0.7rem;
			font-weight: 700;
			letter-spacing: 0.1em;
			text-transform: uppercase;
			color: $text-primary;
			transition: color 0.1s ease;
			cursor: pointer;
		}

		&:focus-within .form-label {
			color: $accent;
		}
	}

	.input-wrap {
		position: relative;
		display: flex;
		align-items: center;
		width: 100%;
	}

	.input-prefix, .input-suffix {
		position: absolute;
		top: 0;
		bottom: 0;
		display: flex;
		align-items: center;
		justify-content: center;
		width: 44px;
		color: $text-secondary;
		font-size: 0.85rem;
		pointer-events: none;
		z-index: 1;

		svg {
			width: 18px;
			height: 18px;
		}
	}

	.input-prefix {
		left: 0;
	}

	.input-suffix {
		right: 0;
		pointer-events: auto;
		cursor: pointer;

		&:hover {
			color: $accent;
		}
	}

	.input-bar {
		position: absolute;
		bottom: 0;
		left: 0;
		width: 100%;
		height: 2px;
		background: $accent;
		transform: scaleX(0);
		transform-origin: left;
		transition: transform 0.25s cubic-bezier(0.4, 0, 0.2, 1);
		z-index: 2;

		.form-group:focus-within & {
			transform: scaleX(1);
		}
	}

	.form-input {
		width: 100%;
		background: transparent;
		border: 2px solid $text-secondary;
		color: $text-primary;
		padding: 0.85rem 1rem;
		font-family: $font-body;
		font-size: 1rem;
		font-weight: 500;
		outline: none;
		transition: border-color 0.1s ease, box-shadow 0.1s ease, transform 0.1s ease;
		border-radius: 0;
		box-sizing: border-box;

		&:focus {
			border-color: $accent;
			background: rgba(0, 0, 0, 0.2);
			box-shadow: 4px 4px 0px rgba(255, 68, 0, 0.15);
			transform: translate(-2px, -2px);
		}

		&::placeholder {
			color: $text-secondary;
			opacity: 0.5;
		}

		&:disabled {
			opacity: 0.5;
			cursor: not-allowed;
		}

		&.has-error {
			border-color: #ff4444;
			box-shadow: 4px 4px 0px rgba(255, 68, 68, 0.15);

			&:focus {
				border-color: #ff4444;
				box-shadow: 4px 4px 0px rgba(255, 68, 68, 0.25);
			}
		}

		&.has-prefix {
			padding-left: 44px;
		}

		&.has-suffix {
			padding-right: 44px;
		}
	}

	.form-hint {
		font-family: $font-body;
		font-size: 0.7rem;
		color: $text-secondary;
		opacity: 0.7;
	}

	.form-error-text {
		color: #ff4444;
		font-size: 0.75rem;
		font-family: $font-body;
	}

	:global(body.light-mode) {
		.form-input {
			background: transparent;
			border-color: #000;
			color: #000;

			&:focus {
				border-color: $light-accent;
				box-shadow: 4px 4px 0px #000;
				background: #fff;
			}

			&::placeholder {
				color: #555;
			}

			&.has-error {
				border-color: #ff4444;
				box-shadow: 4px 4px 0px rgba(255, 68, 68, 0.15);
			}
		}

		.form-group .form-label {
			color: #111;
		}

		.form-group:focus-within .form-label {
			color: $light-accent;
		}

		.input-prefix, .input-suffix {
			color: #555;
		}

		.form-hint {
			color: #555;
		}
	}

	@media (prefers-reduced-motion: reduce) {
		.form-input { transition: none; }
	}
</style>
