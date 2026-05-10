<script lang="ts">
    import { motion } from '@humanspeak/svelte-motion';
    import { tick } from 'svelte';
    import type { HTMLAttributes } from 'svelte/elements';
    import Icon from '../Icon/Icon.svelte';

    const WEEKDAYS = ['S', 'M', 'T', 'W', 'T', 'F', 'S'];
    const MONTHS = ['January', 'February', 'March', 'April', 'May', 'June', 'July', 'August', 'September', 'October', 'November', 'December'];

    type Props = {
        label?: string;
        value?: string;
        placeholder?: string;
        min?: string;
        max?: string;
        error?: string;
        hint?: string;
        animate?: boolean;
        class?: string;
        style?: string;
        onchange?: (value: string) => void;
    } & HTMLAttributes<HTMLElement>;

    let {
        label,
        value = $bindable(''),
        placeholder = 'Select date',
        min,
        max,
        error,
        hint,
        animate = false,
        class: className = '',
        style = '',
        onchange,
        ...rest
    }: Props = $props();

    let isOpen = $state(false);
    let triggerRef: HTMLButtonElement | undefined = $state();
    let dropdownRef: HTMLDivElement | undefined = $state();

    let viewYear = $state(new Date().getFullYear());
    let viewMonth = $state(new Date().getMonth());

    let selectedDate = $derived(value ? new Date(value + 'T00:00:00') : null);
    let displayText = $derived(value ? new Date(value + 'T00:00:00').toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' }) : '');

    let today = $derived(new Date());

    let daysInMonth = $derived(new Date(viewYear, viewMonth + 1, 0).getDate());
    let firstDayOfWeek = $derived(new Date(viewYear, viewMonth, 1).getDay());
    let calendarDays = $derived.by(() => {
        const days: (number | null)[] = [];
        for (let i = 0; i < firstDayOfWeek; i++) days.push(null);
        for (let d = 1; d <= daysInMonth; d++) days.push(d);
        return days;
    });

    function isPastMin(day: number) {
        if (!min) return false;
        return new Date(viewYear, viewMonth, day) < new Date(min + 'T00:00:00');
    }

    function isPastMax(day: number) {
        if (!max) return false;
        return new Date(viewYear, viewMonth, day) > new Date(max + 'T00:00:00');
    }

    function toggleDropdown() {
        isOpen = !isOpen;
        if (isOpen) {
            if (selectedDate) {
                viewMonth = selectedDate.getMonth();
                viewYear = selectedDate.getFullYear();
            }
        }
    }

    function selectDay(day: number) {
        const m = viewMonth + 1;
        const d = String(day).padStart(2, '0');
        const month = String(m).padStart(2, '0');
        const formatted = `${viewYear}-${month}-${d}`;
        value = formatted;
        isOpen = false;
        onchange?.(formatted);
    }

    function prevMonth() {
        if (viewMonth === 0) { viewMonth = 11; viewYear--; }
        else viewMonth--;
    }

    function nextMonth() {
        if (viewMonth === 11) { viewMonth = 0; viewYear++; }
        else viewMonth++;
    }

    function prevYear() { viewYear--; }
    function nextYear() { viewYear++; }

    function goToday() {
        const d = new Date();
        viewMonth = d.getMonth();
        viewYear = d.getFullYear();
        selectDay(d.getDate());
    }

    function isSelected(day: number) {
        if (!selectedDate) return false;
        return selectedDate.getDate() === day
            && selectedDate.getMonth() === viewMonth
            && selectedDate.getFullYear() === viewYear;
    }

    function isToday(day: number) {
        return today.getDate() === day
            && today.getMonth() === viewMonth
            && today.getFullYear() === viewYear;
    }

    function isDisabled(day: number) {
        return isPastMin(day) || isPastMax(day);
    }

    function handleWrapperKeydown(e: KeyboardEvent) {
        if (e.key === 'Escape' && isOpen) {
            isOpen = false;
            tick().then(() => triggerRef?.focus());
        }
    }
</script>

<div class="datepicker-wrapper{className ? ' ' + className : ''}" {style} onkeydown={handleWrapperKeydown} {...rest}>
    {#if label}
        <label class="datepicker-label">{label}</label>
    {/if}

    {#if animate}
        <motion.button
            bind:this={triggerRef}
            class="datepicker-trigger{isOpen ? ' open' : ''}{error ? ' has-error' : ''}"
            type="button"
            onclick={toggleDropdown}
            whileTap={{ scale: 0.998 }}
            aria-haspopup="dialog"
            aria-expanded={isOpen}
        >
            <span class="datepicker-value{value ? '' : ' placeholder'}">
                {displayText || placeholder}
            </span>
            <Icon class="datepicker-chevron{isOpen ? ' open' : ''}" name="chevron-down" size={14} />
        </motion.button>
    {:else}
        <button
            bind:this={triggerRef}
            class="datepicker-trigger"
            class:open={isOpen}
            class:has-error={!!error}
            type="button"
            onclick={toggleDropdown}
            aria-haspopup="dialog"
            aria-expanded={isOpen}
        >
            <span class="datepicker-value" class:placeholder={!value}>
                {displayText || placeholder}
            </span>
            <Icon class="datepicker-chevron{isOpen ? ' open' : ''}" name="chevron-down" size={14} />
        </button>
    {/if}

    {#if hint && !error}
        <span class="datepicker-hint">{hint}</span>
    {/if}
    {#if error}
        <span class="datepicker-error">{error}</span>
    {/if}

    {#if isOpen}
        <div class="datepicker-dropdown" bind:this={dropdownRef} role="dialog" aria-label="Date picker">
            <div class="datepicker-nav">
                <div class="datepicker-nav-group">
                    <button class="datepicker-nav-btn" type="button" onclick={prevYear} aria-label="Previous year">
                        <Icon name="chevrons-left" size={14} />
                    </button>
                    <button class="datepicker-nav-btn" type="button" onclick={prevMonth} aria-label="Previous month">
                        <Icon name="chevron-left" size={14} />
                    </button>
                </div>
                <span class="datepicker-nav-label">{MONTHS[viewMonth]} {viewYear}</span>
                <div class="datepicker-nav-group">
                    <button class="datepicker-nav-btn" type="button" onclick={nextMonth} aria-label="Next month">
                        <Icon name="chevron-right" size={14} />
                    </button>
                    <button class="datepicker-nav-btn" type="button" onclick={nextYear} aria-label="Next year">
                        <Icon name="chevrons-right" size={14} />
                </div>
            </div>

            <div class="datepicker-weekdays">
                {#each WEEKDAYS as wd}
                    <span class="datepicker-weekday">{wd}</span>
                {/each}
            </div>

            <div class="datepicker-grid">
                {#each calendarDays as day}
                    {#if day}
                        <button
                            class="datepicker-day"
                            class:selected={isSelected(day)}
                            class:today={isToday(day)}
                            class:disabled={isDisabled(day)}
                            type="button"
                            onclick={() => selectDay(day)}
                            disabled={isDisabled(day)}
                        >
                            {day}
                        </button>
                    {:else}
                        <span class="datepicker-day-empty"></span>
                    {/if}
                {/each}
            </div>

            <div class="datepicker-footer">
                <button class="datepicker-today-btn" type="button" onclick={goToday}>Today</button>
                <button class="datepicker-today-btn" type="button" onclick={() => { value = ''; isOpen = false; onchange?.(''); }}>Clear</button>
            </div>
        </div>
    {/if}
</div>

<style lang="scss">
    @use '../../../styles' as *;

    .datepicker-wrapper {
        position: relative;
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
        font-family: $font-body;
    }

    .datepicker-label {
        font-size: 0.7rem;
        font-weight: 700;
        text-transform: uppercase;
        letter-spacing: 0.1em;
        color: var(--text-secondary);
        transition: color 0.1s ease;
    }

    .datepicker-trigger {
        display: flex;
        align-items: center;
        justify-content: space-between;
        width: 100%;
        background: transparent;
        border: 2px solid var(--text-secondary);
        color: var(--text-primary);
        padding: 0.85rem 1rem;
        font-family: $font-body;
        font-size: 1rem;
        font-weight: 500;
        cursor: pointer;
        outline: none;
        transition: border-color 0.1s ease, box-shadow 0.1s ease, transform 0.1s ease;
        box-sizing: border-box;
        gap: 0.5rem;
        border-radius: 0;

        &:focus {
            border-color: $accent;
            box-shadow: 4px 4px 0px rgba(255, 68, 0, 0.15);
            transform: translate(-2px, -2px);
        }

        &.open {
            border-color: $accent;
            box-shadow: 4px 4px 0px rgba(255, 68, 0, 0.15);
            transform: translate(-2px, -2px);

            .datepicker-chevron {
                transform: rotate(180deg);
            }
        }

        &.has-error {
            border-color: #ff4444;
        }
    }

    .datepicker-value {
        flex: 1;
        text-align: left;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;

        &.placeholder {
            color: var(--text-secondary);
            opacity: 0.5;
        }
    }

    .datepicker-chevron {
        width: 14px;
        height: 14px;
        flex-shrink: 0;
        transition: transform 0.2s ease;
        color: var(--text-secondary);
    }

    .datepicker-hint,
    .datepicker-error {
        font-size: 0.7rem;
        letter-spacing: 0.05em;
    }

    .datepicker-hint {
        color: var(--text-secondary);
        opacity: 0.7;
    }

    .datepicker-error {
        color: #ff4444;
    }

    .datepicker-dropdown {
        position: absolute;
        top: calc(100% + 4px);
        left: -2px;
        right: -2px;
        z-index: 100;
        background: var(--surface-dark, $surface-dark);
        border: 2px solid $accent;
        box-shadow: 6px 6px 0px rgba(255, 68, 0, 0.2);
        animation: dropdownIn 0.12s ease-out;
    }

    @keyframes dropdownIn {
        from { opacity: 0; transform: translateY(-4px); }
        to { opacity: 1; transform: translateY(0); }
    }

    .datepicker-nav {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 0.6rem 0.5rem;
        border-bottom: 1px solid var(--grid-line, rgba(255,255,255,0.08));
    }

    .datepicker-nav-group {
        display: flex;
        gap: 2px;
    }

    .datepicker-nav-btn {
        background: transparent;
        border: 1px solid transparent;
        color: var(--text-secondary);
        width: 28px;
        height: 28px;
        display: flex;
        align-items: center;
        justify-content: center;
        cursor: pointer;
        transition: all 0.1s ease;
        border-radius: 0;
        font-size: 0;

        &:hover {
            color: $accent;
            border-color: $accent;
            background: rgba(255, 68, 0, 0.08);
        }

        svg {
            width: 14px;
            height: 14px;
        }
    }

    .datepicker-nav-label {
        font-family: $font-heading;
        font-size: 0.85rem;
        font-weight: 600;
        text-transform: uppercase;
        letter-spacing: 0.05em;
        color: var(--text-primary);
    }

    .datepicker-weekdays {
        display: grid;
        grid-template-columns: repeat(7, 1fr);
        padding: 0.4rem 0.25rem 0.2rem;
        border-bottom: 1px solid var(--grid-line, rgba(255,255,255,0.08));
    }

    .datepicker-weekday {
        text-align: center;
        font-size: 0.65rem;
        font-weight: 700;
        text-transform: uppercase;
        letter-spacing: 0.05em;
        color: var(--text-secondary);
        opacity: 0.6;
        padding: 0.25rem 0;
    }

    .datepicker-grid {
        display: grid;
        grid-template-columns: repeat(7, 1fr);
        padding: 0.25rem;
        gap: 1px;
    }

    .datepicker-day {
        display: flex;
        align-items: center;
        justify-content: center;
        height: 34px;
        background: transparent;
        border: 1px solid transparent;
        color: var(--text-primary);
        font-family: $font-body;
        font-size: 0.8rem;
        font-weight: 500;
        cursor: pointer;
        transition: all 0.1s ease;
        border-radius: 0;

        &:hover:not(:disabled) {
            border-color: $accent;
            background: rgba(255, 68, 0, 0.08);
            color: $accent;
        }

        &.today {
            border-color: var(--text-secondary);
            font-weight: 700;
        }

        &.selected {
            background: $accent;
            color: #000;
            border-color: $accent;
            font-weight: 700;
        }

        &:disabled {
            opacity: 0.25;
            cursor: not-allowed;
        }
    }

    .datepicker-day-empty {
        height: 34px;
    }

    .datepicker-footer {
        display: flex;
        justify-content: space-between;
        padding: 0.4rem 0.5rem 0.5rem;
        border-top: 1px solid var(--grid-line, rgba(255,255,255,0.08));
    }

    .datepicker-today-btn {
        background: transparent;
        border: none;
        color: var(--text-secondary);
        font-family: $font-body;
        font-size: 0.7rem;
        font-weight: 600;
        text-transform: uppercase;
        letter-spacing: 0.1em;
        cursor: pointer;
        padding: 0.25rem 0.5rem;
        transition: color 0.1s ease;
        border-radius: 0;

        &:hover {
            color: $accent;
        }
    }

    :global(body.light-mode) {
        .datepicker-trigger {
            border-color: #000;
            color: #000;

            &:focus,
            &.open {
                border-color: $light-accent;
                box-shadow: 4px 4px 0px #000;
                background: #fff;
            }
        }

        .datepicker-value.placeholder {
            color: #555;
            opacity: 0.6;
        }

        .datepicker-dropdown {
            background: #fff;
            box-shadow: 6px 6px 0px rgba(0, 0, 0, 0.15);
        }

        .datepicker-day {
            color: #111;

            &:hover:not(:disabled) {
                color: $light-accent;
                border-color: $light-accent;
            }

            &.today {
                border-color: #111;
            }

            &.selected {
                background: $light-accent;
                color: #fff;
                border-color: $light-accent;
            }
        }

        .datepicker-nav-label {
            color: #111;
        }

        .datepicker-today-btn:hover {
            color: $light-accent;
        }
    }

    @media (prefers-reduced-motion: reduce) {
        .datepicker-dropdown {
            animation: none;
        }

        .datepicker-trigger,
        .datepicker-chevron {
            transition: none;
        }
    }

    @media (max-width: 480px) {
        .datepicker-dropdown {
            position: fixed;
            top: 50%;
            left: 50%;
            transform: translate(-50%, -50%);
            width: calc(100% - 2rem);
            max-width: 340px;
            box-shadow: 0 0 0 9999px rgba(0, 0, 0, 0.5);
        }
    }
</style>
