<script lang="ts">
    import type { HTMLAttributes } from 'svelte/elements';

    export type IconName =
        | 'chevron-left'
        | 'chevron-right'
        | 'chevron-down'
        | 'chevrons-left'
        | 'chevrons-right'
        | 'external-link'
        | 'close'
        | 'envelope'
        | 'file'
        | 'github'
        | 'linkedin'
        | 'arrow-up-right'
        | 'eye'
        | 'eye-off'
        | 'search'
        | 'calendar'
        | 'check'
        | 'user'
        | 'chat'
        | 'clock'
        | 'sun'
        | 'moon'
        | 'check-circle'
        | 'verified'
        | 'lock'
        | 'zap'
        | 'leaf'
        | 'smartphone';

    type Props = {
        name: IconName;
        size?: number;
        color?: string;
        class?: string;
        style?: string;
    } & HTMLAttributes<HTMLElement>;

    let {
        name,
        size = 16,
        color,
        class: className = '',
        style = '',
        ...rest
    }: Props = $props();

    const iconStyle = () => {
        let s = style || '';
        if (color) s += `; color: ${color}`;
        return s;
    };
</script>

<span
    class="ui-icon {className}"
    style={iconStyle()}
    {...rest}
>
    <svg
        xmlns="http://www.w3.org/2000/svg"
        width={size}
        height={size}
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
    >
        {#if name === 'calendar'}
            <rect x="3" y="4" width="18" height="18" rx="2" ry="2" />
            <line x1="16" y1="2" x2="16" y2="6" />
            <line x1="8" y1="2" x2="8" y2="6" />
            <line x1="3" y1="10" x2="21" y2="10" />
        {:else if name === 'check'}
            <polyline points="20 6 9 17 4 12" />
        {:else if name === 'check-circle'}
            <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14" />
            <polyline points="22 4 12 14.01 9 11.01" />
        {:else if name === 'chevron-down'}
            <polyline points="6 9 12 15 18 9" />
        {:else if name === 'chevron-left'}
            <path d="M15 18l-6-6 6-6" />
        {:else if name === 'chevron-right'}
            <path d="M9 18l6-6-6-6" />
        {:else if name === 'chevrons-left'}
            <polyline points="11 17 6 12 11 7" />
            <polyline points="18 17 13 12 18 7" />
        {:else if name === 'chevrons-right'}
            <polyline points="13 17 18 12 13 7" />
            <polyline points="6 17 11 12 6 7" />
        {:else if name === 'eye'}
            <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" />
            <circle cx="12" cy="12" r="3" />
        {:else if name === 'eye-off'}
            <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94" />
            <path d="M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19" />
            <line x1="1" y1="1" x2="23" y2="23" />
        {:else if name === 'external-link'}
            <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6" />
            <polyline points="15 3 21 3 21 9" />
            <line x1="10" y1="14" x2="21" y2="3" />
        {:else if name === 'close'}
            <line x1="18" y1="6" x2="6" y2="18" />
            <line x1="6" y1="6" x2="18" y2="18" />
        {:else if name === 'envelope'}
            <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z" />
            <polyline points="22,6 12,13 2,6" />
        {:else if name === 'file'}
            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z" />
            <polyline points="14 2 14 8 20 8" />
            <line x1="16" y1="13" x2="8" y2="13" />
            <line x1="16" y1="17" x2="8" y2="17" />
            <polyline points="10 9 9 9 8 9" />
        {:else if name === 'search'}
            <circle cx="11" cy="11" r="8" />
            <line x1="21" y1="21" x2="16.65" y2="16.65" />
        {:else if name === 'user'}
            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
            <circle cx="12" cy="7" r="4" />
        {:else if name === 'github'}
            <path d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0 1 12 6.844a9.59 9.59 0 0 1 2.504.338c1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.02 10.02 0 0 0 22 12.017C22 6.484 17.522 2 12 2z" />
        {:else if name === 'linkedin'}
            <path d="M16 8a6 6 0 0 1 6 6v7h-4v-7a2 2 0 0 0-2-2 2 2 0 0 0-2 2v7h-4v-7a6 6 0 0 1 6-6z" />
            <rect x="2" y="9" width="4" height="12" />
            <circle cx="4" cy="4" r="2" />
        {:else if name === 'arrow-up-right'}
            <line x1="7" y1="17" x2="17" y2="7" />
            <polyline points="7 7 17 7 17 17" />
        {:else if name === 'chat'}
            <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
        {:else if name === 'clock'}
            <circle cx="12" cy="12" r="10" />
            <polyline points="12 6 12 12 16 14" />
        {:else if name === 'sun'}
            <circle cx="12" cy="12" r="5" />
            <line x1="12" y1="1" x2="12" y2="3" />
            <line x1="12" y1="21" x2="12" y2="23" />
            <line x1="4.22" y1="4.22" x2="5.64" y2="5.64" />
            <line x1="18.36" y1="18.36" x2="19.78" y2="19.78" />
            <line x1="1" y1="12" x2="3" y2="12" />
            <line x1="21" y1="12" x2="23" y2="12" />
            <line x1="4.22" y1="19.78" x2="5.64" y2="18.36" />
            <line x1="18.36" y1="5.64" x2="19.78" y2="4.22" />
        {:else if name === 'verified'}
            <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" />
            <polyline points="9 12 11 14 15 10" style="stroke:currentColor;stroke-width:3;stroke-linecap:round;stroke-linejoin:round" />
        {:else if name === 'moon'}
            <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z" />
        {:else if name === 'lock'}
            <rect x="3" y="11" width="18" height="11" rx="2" ry="2" />
            <path d="M7 11V7a5 5 0 0 1 10 0v4" />
        {:else if name === 'zap'}
            <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2" />
        {:else if name === 'leaf'}
            <path d="M11 20A7 7 0 0 1 9.8 6.9C15.5 4.9 17 3.5 19 2c1 2 2 4.5 2 8 0 5.5-4.78 10-10 10Z" />
            <path d="M2 21c0-3 1.85-5.36 5.08-6C9.5 14.52 12 13 13 12" />
        {:else if name === 'smartphone'}
            <rect x="5" y="2" width="14" height="20" rx="2" ry="2" />
            <line x1="12" y1="18" x2="12.01" y2="18" />
        {/if}
    </svg>
</span>

<style lang="scss">
	.ui-icon {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		flex-shrink: 0;
		transition: color 0.2s ease;

		svg {
			display: block;
		}
	}
</style>
