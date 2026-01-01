import { writable } from 'svelte/store';
import { browser } from '$app/environment';

// Detect OS preference
function getInitialTheme(): 'dark' | 'light' {
    if (!browser) return 'dark';

    const saved = localStorage.getItem('theme');
    if (saved === 'light' || saved === 'dark') return saved;

    return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
}

function createThemeStore() {
    const { subscribe, set, update } = writable<'dark' | 'light'>(getInitialTheme());

    return {
        subscribe,
        toggle: () => {
            update(current => {
                const newTheme = current === 'dark' ? 'light' : 'dark';
                if (browser) {
                    localStorage.setItem('theme', newTheme);
                }
                return newTheme;
            });
        },
        set: (theme: 'dark' | 'light') => {
            if (browser) {
                localStorage.setItem('theme', theme);
            }
            set(theme);
        }
    };
}

export const theme = createThemeStore();

// Listen for OS theme changes
if (browser) {
    window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
        // Only auto-switch if user hasn't set manual preference
        if (!localStorage.getItem('theme')) {
            theme.set(e.matches ? 'dark' : 'light');
        }
    });
}
