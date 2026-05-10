<script lang="ts">
    import { motion } from '@humanspeak/svelte-motion';
    import type { HTMLAttributes } from 'svelte/elements';

    type Props = {
        title: string;
        count?: number;
        id?: string;
        animate?: boolean;
        class?: string;
        style?: string;
    } & HTMLAttributes<HTMLElement>;

    let {
        title,
        count,
        id,
        animate = false,
        class: className = '',
        style = '',
        ...rest
    }: Props = $props();
</script>

{#if animate}
    <motion.div
        class="section-header{className ? ' ' + className : ''}"
        {style}
        {id}
        initial={{ clipPath: 'polygon(0 0, 0 0, 0 100%, 0 100%)' }}
        whileInView={{ clipPath: 'polygon(0 0, 100% 0, 100% 100%, 0 100%)' }}
        viewport={{ once: true }}
        transition={{ duration: 0.6, ease: [0.16, 1, 0.3, 1] }}
        {...rest}
    >
        <h2>
            {title}
            {#if count !== undefined}
                <span class="count-badge">{String(count).padStart(2, '0')}</span>
            {/if}
        </h2>
    </motion.div>
{:else}
    <div
        class="section-header{className ? ' ' + className : ''}"
        {style}
        {id}
        {...rest}
    >
        <h2>
            {title}
            {#if count !== undefined}
                <span class="count-badge">{String(count).padStart(2, '0')}</span>
            {/if}
        </h2>
    </div>
{/if}

<style lang="scss">
	@use '../../../styles' as *;

	.section-header {
		margin-bottom: 0.5rem;
		position: relative;

		h2 {
			font-family: $font-heading;
			font-size: 1.5rem;
			font-weight: 600;
			text-transform: uppercase;
			letter-spacing: 0.05em;
			color: #fff;
			border-left: 4px solid $accent;
			padding-left: 1rem;
			margin: 0;
		}

		.count-badge {
			display: inline-flex;
			align-items: center;
			justify-content: center;
			font-family: $font-body;
			font-size: 0.65rem;
			color: $accent;
			border: 1px solid $accent;
			padding: 1px 7px;
			margin-left: 8px;
			vertical-align: middle;
			line-height: 1;
		}
	}

	:global(body.light-mode) .section-header h2 {
		color: #111;
		border-color: $light-accent;
	}

	:global(body.light-mode) .section-header .count-badge {
		color: $light-accent;
		border-color: $light-accent;
	}

	.section-header.clip-reveal {
		h2, .count-badge {
			clip-path: polygon(0 0, 0 0, 0 100%, 0 100%);
			transition: clip-path 0.6s cubic-bezier(0.16, 1, 0.3, 1);
		}

		&.visible {
			h2, .count-badge {
				clip-path: polygon(0 0, 100% 0, 100% 100%, 0 100%);
			}
		}
	}
</style>
