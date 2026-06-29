<script lang="ts">
    import { motion } from '@humanspeak/svelte-motion';
    import { tick } from 'svelte';
    import type { HTMLAttributes } from 'svelte/elements';
    import Icon from '../Icon/Icon.svelte';

    export interface SelectOption {
        value: string;
        label: string;
    }

    type Props = {
        label?: string;
        options: SelectOption[];
        value?: string;
        placeholder?: string;
        searchable?: boolean;
        clearable?: boolean;
        disabled?: boolean;
        error?: string;
        animate?: boolean;
        class?: string;
        style?: string;
        onchange?: (value: string) => void;
    } & HTMLAttributes<HTMLElement>;

    let {
        label,
        options,
        value = $bindable(''),
        placeholder = 'Select...',
        searchable = false,
        clearable = false,
        disabled = false,
        error,
        animate = false,
        class: className = '',
        style = '',
        onchange,
        ...rest
    }: Props = $props();

    let isOpen = $state(false);
    let searchQuery = $state('');
    let highlightedIndex = $state(0);
    let triggerRef: any = $state();
    let dropdownRef: HTMLDivElement | undefined = $state();
    let searchInputRef: HTMLInputElement | undefined = $state();

    let selectedOption = $derived(options.find(o => o.value === value) || null);
    let filteredOptions = $derived(
        searchable && searchQuery
            ? options.filter(o =>
                  o.label.toLowerCase().includes(searchQuery.toLowerCase())
              )
            : options
    );

    let selectId = $state(`select-${Math.random().toString(36).slice(2, 8)}`);

    function toggle() {
        if (disabled) return;
        isOpen = !isOpen;
        if (isOpen) {
            highlightedIndex = Math.max(0, options.findIndex(o => o.value === value));
            tick().then(() => {
                searchInputRef?.focus();
            });
        } else {
            searchQuery = '';
        }
    }

    function open() {
        if (disabled || isOpen) return;
        isOpen = true;
        highlightedIndex = Math.max(0, options.findIndex(o => o.value === value));
        tick().then(() => searchInputRef?.focus());
    }

    function close() {
        isOpen = false;
        searchQuery = '';
        triggerRef?.focus();
    }

    function select(option: SelectOption) {
        value = option.value;
        onchange?.(option.value);
        close();
    }

    function clear() {
        value = '';
        onchange?.('');
        close();
    }

    function onKeydown(e: KeyboardEvent) {
        if (!isOpen) {
            if (e.key === 'Enter' || e.key === ' ' || e.key === 'ArrowDown') {
                e.preventDefault();
                open();
            }
            return;
        }

        switch (e.key) {
            case 'ArrowDown':
                e.preventDefault();
                highlightedIndex = Math.min(highlightedIndex + 1, filteredOptions.length - 1);
                break;
            case 'ArrowUp':
                e.preventDefault();
                highlightedIndex = Math.max(highlightedIndex - 1, 0);
                break;
            case 'Enter':
                e.preventDefault();
                if (filteredOptions[highlightedIndex]) {
                    select(filteredOptions[highlightedIndex]);
                }
                break;
            case 'Escape':
                e.preventDefault();
                close();
                break;
            case 'Tab':
                close();
                break;
        }
    }

    function onTriggerKeydown(e: KeyboardEvent) {
        if (disabled) return;
        if (e.key === 'Enter' || e.key === ' ' || e.key === 'ArrowDown') {
            e.preventDefault();
            toggle();
        }
    }

    function onBlur(e: FocusEvent) {
        if (!dropdownRef?.contains(e.relatedTarget as Node) &&
            e.relatedTarget !== triggerRef) {
            close();
        }
    }
</script>

<div
    class="select-wrapper{className ? ' ' + className : ''}"
    {style}
    onkeydown={onKeydown}
    onfocusout={onBlur}
    role="combobox"
    aria-expanded={isOpen}
    aria-haspopup="listbox"
    aria-controls={`${selectId}-listbox`}
    {...rest}
>
    {#if label}
        <label for={selectId} class="select-label">{label}</label>
    {/if}

    {#if animate}
        <motion.button
            bind:this={triggerRef}
            id={selectId}
            class="select-trigger{isOpen ? ' open' : ''}{value ? ' has-value' : ''}{error ? ' has-error' : ''}"
            onclick={toggle}
            onkeydown={onTriggerKeydown}
            {disabled}
            type="button"
            whileTap={{ scale: disabled ? 1 : 0.99 }}
            aria-labelledby={label ? selectId : undefined}
        >
            <span class={value ? 'select-value' : 'select-placeholder'}>
                {selectedOption ? selectedOption.label : placeholder}
            </span>
            <Icon class="select-chevron{isOpen ? ' rotated' : ''}" name="chevron-down" size={14} />
        </motion.button>
    {:else}
        <button
            bind:this={triggerRef}
            id={selectId}
            class="select-trigger{isOpen ? ' open' : ''}{value ? ' has-value' : ''}{error ? ' has-error' : ''}"
            onclick={toggle}
            onkeydown={onTriggerKeydown}
            {disabled}
            type="button"
            aria-labelledby={label ? selectId : undefined}
        >
            <span class={value ? 'select-value' : 'select-placeholder'}>
                {selectedOption ? selectedOption.label : placeholder}
            </span>
            <Icon class="select-chevron{isOpen ? ' rotated' : ''}" name="chevron-down" size={14} />
        </button>
    {/if}

    {#if isOpen}
        <!-- svelte-ignore a11y_no_static_element_interactions -->
        <div
            bind:this={dropdownRef}
            class="select-dropdown"
            role="listbox"
            id={`${selectId}-listbox`}
            onmouseleave={() => { highlightedIndex = -1; }}
        >
            {#if searchable}
                <div class="select-search">
                    <input
                        bind:this={searchInputRef}
                        type="text"
                        bind:value={searchQuery}
                        placeholder="Search..."
                        onkeydown={(e) => e.stopPropagation()}
                        aria-label="Search options"
                    />
                </div>
            {/if}
            <div class="select-options">
                {#if filteredOptions.length === 0}
                    <div class="select-no-results">No results found</div>
                {:else}
                    {#each filteredOptions as option, i}
                        <!-- svelte-ignore a11y_no_static_element_interactions -->
                        <div
                            class="select-option"
                            class:selected={option.value === value}
                            class:highlighted={i === highlightedIndex}
                            role="option"
                            aria-selected={option.value === value}
                            onclick={() => select(option)}
                            onmouseenter={() => { highlightedIndex = i; }}
                            onkeydown={(e) => { if (e.key === 'Enter') select(option); }}
                        >
                            <span>{option.label}</span>
                            {#if option.value === value}
                                <Icon class="select-check" name="check" size={14} />
                            {/if}
                        </div>
                    {/each}
                {/if}
            </div>
        </div>
    {/if}

    {#if error}
        <span class="select-error">{error}</span>
    {/if}
</div>

<style lang="scss">
	@use '../../../styles' as *;

	.select-wrapper {
		position: relative;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		font-family: $font-body;
	}

	.select-label {
		font-size: 0.7rem;
		font-weight: 700;
		letter-spacing: 0.1em;
		text-transform: uppercase;
		color: $text-primary;
		transition: color 0.1s ease;
		cursor: pointer;
	}

	.select-trigger {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 0.75rem;
		width: 100%;
		padding: 0.85rem 1rem;
		background: transparent;
		border: 2px solid $text-secondary;
		color: $text-primary;
		font-family: $font-body;
		font-size: 1rem;
		font-weight: 500;
		cursor: pointer;
		transition: border-color 0.1s ease, box-shadow 0.1s ease, transform 0.1s ease;
		position: relative;
		text-align: left;

		&:focus, &.open {
			border-color: $accent;
			box-shadow: 4px 4px 0px rgba(255, 68, 0, 0.15);
			transform: translate(-2px, -2px);
		}

		&.has-value {
			color: $text-primary;
		}

		&.has-error {
			border-color: #ff4444;
			box-shadow: 4px 4px 0px rgba(255, 68, 68, 0.15);

			&:focus {
				border-color: #ff4444;
				box-shadow: 4px 4px 0px rgba(255, 68, 68, 0.25);
			}
		}

		&:disabled {
			opacity: 0.5;
			cursor: not-allowed;
		}

		.select-placeholder {
			color: $text-secondary;
			opacity: 0.5;
		}

		.select-value {
			display: flex;
			align-items: center;
			gap: 0.5rem;
		}

		.select-chevron {
			width: 16px;
			height: 16px;
			flex-shrink: 0;
			color: $text-secondary;
			transition: transform 0.2s ease;

			&.rotated {
				transform: rotate(180deg);
			}
		}
	}

	.select-dropdown {
		position: absolute;
		top: calc(100% + 4px);
		left: -2px;
		right: -2px;
		z-index: 100;
		background: $surface-dark;
		border: 2px solid $accent;
		box-shadow: 6px 6px 0px rgba(255, 68, 0, 0.2);
		max-height: 250px;
		display: flex;
		flex-direction: column;
		overflow: hidden;
	}

	.select-search {
		padding: 0.5rem;
		border-bottom: 1px solid $grid-line;

		input {
			width: 100%;
			padding: 0.5rem;
			background: transparent;
			border: 1px solid $grid-line;
			color: $text-primary;
			font-family: $font-body;
			font-size: 0.85rem;
			outline: none;
			transition: border-color 0.1s ease;

			&:focus {
				border-color: $accent;
			}

			&::placeholder {
				color: $text-secondary;
				opacity: 0.5;
			}
		}
	}

	.select-options {
		overflow-y: auto;
		padding: 0.25rem 0;
		flex: 1;

		&::-webkit-scrollbar {
			width: 4px;
		}

		&::-webkit-scrollbar-thumb {
			background: $accent;
		}
	}

	.select-option {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 0.5rem;
		padding: 0.75rem 1rem;
		cursor: pointer;
		font-family: $font-body;
		font-size: 0.9rem;
		color: $text-secondary;
		transition: all 0.1s ease;
		position: relative;

		&:hover, &.highlighted {
			background: rgba(255, 68, 0, 0.1);
			color: $text-primary;
		}

		&.selected {
			color: $accent;
			font-weight: 600;
		}

		&.active {
			background: $accent;
			color: #000;
		}

		.select-check {
			width: 14px;
			height: 14px;
			flex-shrink: 0;
			color: $accent;
		}
	}

	.select-no-results {
		padding: 1.5rem 1rem;
		text-align: center;
		color: $text-secondary;
		font-size: 0.8rem;
		opacity: 0.6;
	}

	.select-error {
		color: #ff4444;
		font-size: 0.75rem;
		margin-top: 0.25rem;
	}

	:global(body.light-mode) {
		.select-trigger {
			border-color: #000;
			color: #000;

			.select-placeholder {
				color: #555;
				opacity: 0.6;
			}

			&.open, &:focus {
				border-color: $light-accent;
				box-shadow: 4px 4px 0px #000;
			}

			&.has-error {
				border-color: #ff4444;
			}
		}

		.select-label {
			color: #111;
		}

		.select-dropdown {
			background: #fff;
			border-color: #111;
			box-shadow: 6px 6px 0px rgba(0, 0, 0, 0.1);
		}

		.select-search input {
			border-color: rgba(0, 0, 0, 0.15);
			color: #000;

			&:focus {
				border-color: $light-accent;
			}

			&::placeholder {
				color: #888;
			}
		}

		.select-option {
			color: #555;

			&:hover, &.highlighted {
				background: rgba(217, 61, 0, 0.08);
				color: #000;
			}

			&.selected {
				color: $light-accent;
			}

			&.active {
				background: $light-accent;
				color: #fff;
			}
		}

		.select-no-results {
			color: #888;
		}
	}

	@media (prefers-reduced-motion: reduce) {
		.select-trigger,
		.select-option,
		.select-chevron {
			transition: none;
		}
	}
</style>
