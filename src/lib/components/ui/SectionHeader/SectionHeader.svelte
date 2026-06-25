<script lang="ts">
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

    let el = $state<HTMLDivElement | null>(null);

    $effect(() => {
        if (!animate || !el) return;
        const obs = new IntersectionObserver(
            ([entry]) => {
                if (entry.isIntersecting) {
                    el.classList.add('in-view');
                    obs.disconnect();
                }
            },
            { threshold: 0.2 }
        );
        obs.observe(el);
        return () => obs.disconnect();
    });
</script>

<div
    class="section-header{animate ? ' section-header-animate' : ''}{className ? ' ' + className : ''}"
    class:section-header-animate={animate}
    {style}
    {id}
    bind:this={el}
    {...rest}
>
    <h2>
        {title}
        {#if count !== undefined}
            <span class="count-badge">{String(count).padStart(2, '0')}</span>
        {/if}
    </h2>
</div>

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

	.section-header-animate {
		h2 {
			clip-path: polygon(0 0, 0 0, 0 100%, 0 100%);
			transition: clip-path 0.6s cubic-bezier(0.16, 1, 0.3, 1);
		}

		.count-badge {
			clip-path: polygon(0 0, 0 0, 0 100%, 0 100%);
			transition: clip-path 0.6s cubic-bezier(0.16, 1, 0.3, 1);
		}

		&.in-view {
			h2, .count-badge {
				clip-path: polygon(0 0, 100% 0, 100% 100%, 0 100%);
			}
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
</style>
