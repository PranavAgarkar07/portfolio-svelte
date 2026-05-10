<script lang="ts">
    import Radio from './Radio.svelte';
    import type { Snippet } from 'svelte';
    import type { HTMLAttributes } from 'svelte/elements';

    type Props = {
        label?: string;
        value?: string;
        name?: string;
        class?: string;
        style?: string;
        children?: Snippet;
        onchange?: (value: string) => void;
    } & HTMLAttributes<HTMLElement>;

    let {
        label,
        value = $bindable(''),
        name = `radio-group-${Math.random().toString(36).slice(2, 8)}`,
        class: className = '',
        style = '',
        children,
        onchange,
        ...rest
    }: Props = $props();

    function handleChange(e: Event) {
        const target = e.target as HTMLInputElement;
        if (target.type === 'radio' && target.checked) {
            value = target.value;
            onchange?.(target.value);
        }
    }
</script>

<div
    class="radio-group{className ? ' ' + className : ''}"
    {style}
    role="radiogroup"
    aria-label={label}
    onchange={handleChange}
    {...rest}
>
    {#if label}
        <span class="radio-group-label">{label}</span>
    {/if}
    {#if children}
        {@render children()}
    {/if}
</div>

<style lang="scss">
	@use '../../../styles' as *;

	.radio-group {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;

		.radio-group-label {
			font-family: $font-body;
			font-size: 0.7rem;
			font-weight: 700;
			letter-spacing: 0.1em;
			text-transform: uppercase;
			color: $text-primary;
		}
	}

	:global(body.light-mode) .radio-group-label {
		color: #111;
	}
</style>
