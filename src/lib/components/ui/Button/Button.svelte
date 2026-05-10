<script lang="ts">
    import { motion } from '@humanspeak/svelte-motion';
    import type { Snippet } from 'svelte';
    import type { HTMLButtonAttributes, HTMLAnchorAttributes } from 'svelte/elements';

    type Variant = 'primary' | 'secondary' | 'ghost' | 'icon';
    type Size = 'sm' | 'md' | 'lg';

    type Props = {
        variant?: Variant;
        size?: Size;
        href?: string;
        target?: string;
        rel?: string;
        disabled?: boolean;
        loading?: boolean;
        ripple?: boolean;
        animate?: boolean;
        class?: string;
        style?: string;
        children?: Snippet;
    } & HTMLButtonAttributes & HTMLAnchorAttributes;

    let {
        variant = 'primary',
        size = 'md',
        href,
        target,
        rel,
        disabled = false,
        loading = false,
        ripple = false,
        animate = false,
        class: className = '',
        style = '',
        children,
        ...rest
    }: Props = $props();

    let btnRef: HTMLElement | undefined = $state();

    function handleRipple(e: MouseEvent) {
        if (!ripple || disabled) return;
        const btn = e.currentTarget as HTMLElement;
        const rect = btn.getBoundingClientRect();
        const x = ((e.clientX - rect.left) / rect.width) * 100;
        const y = ((e.clientY - rect.top) / rect.height) * 100;
        btn.style.setProperty('--ripple-x', `${x}%`);
        btn.style.setProperty('--ripple-y', `${y}%`);
        const rippleEl = btn.querySelector('.btn-ripple');
        if (rippleEl) {
            rippleEl.classList.remove('active');
            void (rippleEl as HTMLElement).offsetWidth;
            rippleEl.classList.add('active');
        }
    }

    const baseClass = () => {
        let c = 'btn';
        c += ` btn-${variant}`;
        c += size !== 'md' ? ` btn-${size}` : '';
        if (loading) c += ' btn-loading';
        if (className) c += ` ${className}`;
        return c;
    };
</script>

{#if href}
    <motion.a
        {href}
        {target}
        {rel}
        class={baseClass()}
        {style}
        bind:this={btnRef}
        {...rest}
    >
        {#if children}
            {@render children()}
        {/if}
    </motion.a>
{:else if animate}
    <motion.button
        class={baseClass()}
        {style}
        {disabled}
        bind:this={btnRef}
        onclick={handleRipple}
        whileHover={{ scale: disabled ? 1 : 1.02 }}
        whileTap={{ scale: disabled ? 1 : 0.98 }}
        {...rest}
    >
        {#if ripple}
            <span class="btn-ripple"></span>
        {/if}
        {#if children}
            {@render children()}
        {/if}
    </motion.button>
{:else}
    <button
        class={baseClass()}
        {style}
        {disabled}
        bind:this={btnRef}
        onclick={handleRipple}
        {...rest}
    >
        {#if ripple}
            <span class="btn-ripple"></span>
        {/if}
        {#if children}
            {@render children()}
        {/if}
    </button>
{/if}

<style lang="scss">
	@use '../../../styles' as *;

	.btn {
		display: inline-flex;
		align-items: center;
		gap: 0.75rem;
		font-family: $font-heading;
		text-transform: uppercase;
		letter-spacing: 0.08em;
		font-size: 0.8rem;
		padding: 10px 24px;
		font-weight: 600;
		border: 2px solid $accent;
		color: $accent;
		background: transparent;
		cursor: pointer;
		position: relative;
		transition: all 0.1s cubic-bezier(0, 0, 0.2, 1);
		box-shadow: 4px 4px 0px $accent;
		line-height: 1;
		text-decoration: none;

		&:hover {
			transform: translate(2px, 2px);
			box-shadow: 2px 2px 0px $accent;
			background: rgba(255, 68, 0, 0.15);
			color: $accent;
			border-color: $accent;
		}

		&:active {
			transform: translate(4px, 4px);
			box-shadow: none;
		}

		&:disabled {
			opacity: 0.6;
			cursor: not-allowed;
			transform: translate(1px, 1px);
			box-shadow: 4px 4px 0px rgba(255, 68, 0, 0.1);
		}

		&.btn-sm {
			padding: 7px 16px;
			font-size: 0.7rem;
		}

		&.btn-lg {
			padding: 14px 30px;
			font-size: 0.9rem;
		}

		&.btn-primary {
			background: $accent;
			color: #000;
			box-shadow: 4px 4px 0px #fff;
			border-color: $accent;

			&:hover {
				background: #fff;
				color: #000;
				border-color: #fff;
				box-shadow: 2px 2px 0px $accent;
			}

			&:active {
				transform: translate(4px, 4px);
				box-shadow: none;
			}
		}

		&.btn-secondary {
			border-color: #fff;
			color: #fff;
			box-shadow: 4px 4px 0px #fff;

			&:hover {
				background: #fff;
				color: #000;
				box-shadow: 2px 2px 0px #fff;
				transform: translate(2px, 2px);
			}

			&:active {
				transform: translate(4px, 4px);
				box-shadow: none;
			}
		}

		&.btn-ghost {
			border-color: transparent;
			color: $text-secondary;
			box-shadow: none;

			&:hover {
				color: $accent;
				background: rgba(255, 68, 0, 0.08);
				transform: none;
				box-shadow: none;
			}

			&:active {
				transform: none;
				box-shadow: none;
			}
		}

		&.btn-icon {
			border: 2px solid $accent;
			color: $accent;
			cursor: pointer;
			padding: 8px;
			width: 40px;
			height: 40px;
			box-shadow: 4px 4px 0px $accent;
			background: transparent;

			&:hover {
				background: rgba(255, 68, 0, 0.15);
				color: $accent;
				transform: translate(2px, 2px);
				box-shadow: 2px 2px 0px $accent;
			}

			&:active {
				transform: translate(4px, 4px);
				box-shadow: none;
			}
		}

		&.btn-loading {
			pointer-events: none;
		}

		.btn-ripple {
			position: absolute;
			inset: 0;
			pointer-events: none;
			overflow: hidden;

			&::after {
				content: '';
				position: absolute;
				top: var(--ripple-y, 50%);
				left: var(--ripple-x, 50%);
				width: 100px;
				height: 100px;
				background: radial-gradient(circle, rgba(255,255,255,0.6) 0%, transparent 70%);
				transform: translate(-50%, -50%) scale(0);
				opacity: 0;
			}

			&.active::after {
				animation: rippleAnim 0.4s ease-out forwards;
			}
		}

		@keyframes rippleAnim {
			0% { transform: translate(-50%, -50%) scale(0); opacity: 1; }
			100% { transform: translate(-50%, -50%) scale(2.5); opacity: 0; }
		}

		.btn-icon-slot {
			width: 16px;
			height: 16px;
			transition: transform 0.2s ease;
			flex-shrink: 0;
			display: inline-flex;
			align-items: center;
			justify-content: center;
		}

		&:hover .btn-icon-slot {
			transform: translateX(4px);
		}
	}

	:global(body.light-mode) .btn {
		border-color: #111;
		color: #111;
		box-shadow: 4px 4px 0px #111;

		&:hover {
			background: #111;
			color: #fff;
			box-shadow: 2px 2px 0px #111;
			border-color: $light-accent;
		}

		&:active {
			transform: translate(4px, 4px);
			box-shadow: none;
		}

		&.btn-primary {
			background: $light-accent;
			border-color: $light-accent;
			color: #000;

			&:hover {
				background: #111;
				color: #fff;
				border-color: #111;
			}
		}

		&.btn-secondary {
			border-color: #111;
			color: #111;
			box-shadow: 4px 4px 0px #111;

			&:hover {
				background: #111;
				color: #fff;
				box-shadow: 2px 2px 0px #111;
			}
		}

		&.btn-ghost {
			box-shadow: none;
			color: $light-text-secondary;

			&:hover {
				color: $light-accent;
				background: rgba(217, 61, 0, 0.08);
				box-shadow: none;
			}
		}

		&.btn-icon {
			border-color: #111;
			color: #111;
			box-shadow: 4px 4px 0px #111;

			&:hover {
				background: #111;
				color: #fff;
				box-shadow: 2px 2px 0px #111;
				border-color: $light-accent;
			}

			&:active {
				transform: translate(4px, 4px);
				box-shadow: none;
			}
		}
	}

	@media (prefers-reduced-motion: reduce) {
		.btn {
			transition: none;
			animation: none;
		}
	}
</style>
