<script lang="ts">
    import { page } from "$app/stores";
    import { base } from "$app/paths";
    import { onMount } from "svelte";
    import { isAdmin, isAuthor, init, logout } from "$lib/stores/auth";
    import { goto } from "$app/navigation";

    let { children } = $props();

    const navItems = [
        { label: "Dashboard", path: `${base}/admin/dashboard`, icon: "D" },
        { label: "Contact", path: `${base}/admin/contact`, icon: "C" },
        { label: "Certificates", path: `${base}/admin/certificates`, icon: "A" },
        { label: "Badges", path: `${base}/admin/badges`, icon: "B" },
        { label: "Blog", path: `${base}/admin/blog`, icon: "P" },
        { label: "Users", path: `${base}/admin/users`, icon: "U" },
        { label: "Reviews", path: `${base}/admin/reviews`, icon: "R" },
    ];

    let currentPath = $derived($page.url.pathname.replace(/\/$/, ''));
    let checking = $state(true);

    onMount(async () => {
        await init();
        checking = false;
        if (!$isAdmin && !$isAuthor) {
            goto(base);
        }
    });
</script>

{#if checking}
    <div class="admin-loading">Checking access...</div>
{:else if $isAdmin || $isAuthor}
    <nav class="admin-nav">
        <div class="admin-nav-inner">
            <a href={base} class="admin-nav-brand">
                <span class="admin-nav-icon">◈</span>
                <span class="admin-nav-title">SENTINEL</span>
            </a>
            <div class="admin-nav-links">
                {#each navItems as item}
                    <a
                        href={item.path}
                        class="admin-nav-link"
                        class:active={currentPath === item.path}
                    >
                        <span class="admin-nav-link-icon">{item.icon}</span>
                        <span class="admin-nav-link-label">{item.label}</span>
                    </a>
                {/each}
                <button
                    class="admin-nav-link admin-nav-logout"
                    onclick={() => {
                        logout();
                        window.location.href = base;
                    }}
                >
                    <span class="admin-nav-link-icon">🚪</span>
                    <span class="admin-nav-link-label">Logout</span>
                </button>
            </div>
        </div>
    </nav>

    <div class="admin-page" id="main-content">
        {@render children()}
    </div>
{/if}

<style>
    .admin-loading {
        display: flex;
        align-items: center;
        justify-content: center;
        min-height: 60vh;
        color: var(--text-muted);
    }
    .admin-nav {
        position: sticky;
        top: var(--nav-height);
        z-index: 100;
        background: var(--bg-dark);
        border-bottom: 1px solid var(--border-color);
        padding: 0 2rem;
    }

    .admin-nav-inner {
        max-width: 1100px;
        margin: 0 auto;
        display: flex;
        align-items: center;
        gap: 2rem;
        height: 48px;
    }

    .admin-nav-brand {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        text-decoration: none;
        color: var(--accent);
        flex-shrink: 0;
    }

    .admin-nav-icon {
        font-size: 1rem;
    }

    .admin-nav-title {
        font-family: var(--font-heading);
        font-size: 0.8rem;
        font-weight: 700;
        letter-spacing: 0.2em;
    }

    .admin-nav-links {
        display: flex;
        align-items: center;
        gap: 0.25rem;
        flex: 1;
        overflow-x: auto;
    }

    .admin-nav-link {
        display: flex;
        align-items: center;
        gap: 0.35rem;
        padding: 0.35rem 0.65rem;
        text-decoration: none;
        color: var(--text-muted);
        font-size: 0.75rem;
        font-family: var(--font-body);
        letter-spacing: 0.05em;
        border: 1px solid transparent;
        transition: all 0.15s ease;
        cursor: pointer;
        background: none;
        white-space: nowrap;
    }

    .admin-nav-link:hover {
        color: var(--text-primary);
        border-color: var(--border-color);
    }

    .admin-nav-link.active {
        color: var(--accent);
        border-color: var(--accent-glow);
        background: var(--accent-glow);
    }

    .admin-nav-link-icon {
        font-size: 0.85rem;
        line-height: 1;
    }

    .admin-nav-link-label {
        font-family: 'Space Grotesk', monospace;
        font-weight: 500;
    }

    .admin-nav-logout {
        margin-left: auto;
        flex-shrink: 0;
    }

    .admin-nav-logout:hover {
        color: #ef4444;
        border-color: rgba(239, 68, 68, 0.2);
    }

    .admin-page {
        padding-top: calc(48px + 3rem);
        background: var(--bg-dark);
        min-height: calc(100vh - 48px);
    }

    @media (max-width: 640px) {
        .admin-nav {
            padding: 0 1rem;
        }
        .admin-nav-title {
            display: none;
        }
        .admin-nav-link-label {
            display: none;
        }
        .admin-nav-link {
            padding: 0.35rem 0.5rem;
        }
    }
</style>
