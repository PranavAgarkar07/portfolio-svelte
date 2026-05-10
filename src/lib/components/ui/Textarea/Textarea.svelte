<script lang="ts">
    import type { HTMLAttributes } from 'svelte/elements';

    type Props = {
        label?: string;
        placeholder?: string;
        value?: string;
        error?: string;
        hint?: string;
        maxlength?: number;
        rows?: number;
        autoGrow?: boolean;
        showCount?: boolean;
        class?: string;
        style?: string;
        onchange?: (value: string) => void;
    } & HTMLAttributes<HTMLTextAreaElement>;

    let {
        label,
        placeholder,
        value = $bindable(''),
        error,
        hint,
        maxlength,
        rows = 4,
        autoGrow = false,
        showCount = false,
        class: className = '',
        style = '',
        onchange,
        ...rest
    }: Props = $props();

    let textareaId = $state(`textarea-${Math.random().toString(36).slice(2, 8)}`);
    let textareaRef: HTMLTextAreaElement | undefined = $state();

    let charLength = $derived(typeof value === 'string' ? value.length : 0);
    let charCountClass = $derived(
        maxlength
            ? charLength > maxlength
                ? 'over-limit'
                : charLength > maxlength * 0.9
                    ? 'near-limit'
                    : ''
            : ''
    );

    function handleInput(e: Event) {
        const target = e.target as HTMLTextAreaElement;
        value = target.value;
        onchange?.(target.value);
        if (autoGrow) autoGrowTextarea(target);
    }

    function autoGrowTextarea(el: HTMLTextAreaElement) {
        el.style.height = 'auto';
        el.style.height = el.scrollHeight + 'px';
    }

    function onTextareaMount(el: HTMLTextAreaElement) {
        if (autoGrow) autoGrowTextarea(el);
    }
</script>

<div class="form-group{className ? ' ' + className : ''}">
    {#if label}
        <label for={textareaId} class="form-label">{label}</label>
    {/if}
    <div class="textarea-wrap">
        <textarea
            bind:this={textareaRef}
            id={textareaId}
            class="form-textarea{error ? ' has-error' : ''}{autoGrow ? ' auto-grow' : ''}"
            {placeholder}
            {value}
            {rows}
            {maxlength}
            oninput={handleInput}
            use:onTextareaMount
            {...rest}
        ></textarea>
        <div class="textarea-bar"></div>
    </div>
    <div class="textarea-footer">
        <div>
            {#if hint && !error}
                <span class="form-hint">{hint}</span>
            {/if}
            {#if error}
                <span class="form-error-text">{error}</span>
            {/if}
        </div>
        {#if showCount && maxlength}
            <span class="char-count {charCountClass}">
                {charLength}/{maxlength}
            </span>
        {/if}
    </div>
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

	.textarea-wrap {
		position: relative;
		display: flex;
		flex-direction: column;
	}

	.textarea-bar {
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

	.form-textarea {
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
		resize: vertical;
		min-height: 120px;
		line-height: 1.5;

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

		&.auto-grow {
			resize: none;
			overflow: hidden;
		}
	}

	.textarea-footer {
		display: flex;
		justify-content: space-between;
		align-items: center;
		gap: 0.5rem;
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

	.char-count {
		font-family: $font-body;
		font-size: 0.65rem;
		color: $text-secondary;
		font-weight: 500;
		letter-spacing: 0.05em;

		&.near-limit {
			color: #ffaa00;
		}

		&.over-limit {
			color: #ff4444;
		}
	}

	:global(body.light-mode) {
		.form-textarea {
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
		}

		.form-group .form-label {
			color: #111;
		}

		.form-group:focus-within .form-label {
			color: $light-accent;
		}

		.form-hint {
			color: #555;
		}

		.char-count {
			color: #555;
		}
	}

	@media (prefers-reduced-motion: reduce) {
		.form-textarea { transition: none; }
	}
</style>
