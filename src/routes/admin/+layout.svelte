<script lang="ts">
    import { page } from "$app/stores";
    import { base } from "$app/paths";

    let { children } = $props();

    const navItems = [
        { label: "Dashboard", path: `${base}/admin/dashboard`, icon: "📊" },
        { label: "Contact", path: `${base}/admin/contact`, icon: "📬" },
        { label: "Certificates", path: `${base}/admin/certificates`, icon: "📜" },
        { label: "Badges", path: `${base}/admin/badges`, icon: "🏅" },
    ];

    let currentPath = $derived($page.url.pathname);
</script>

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
                    localStorage.removeItem("contact_admin_key");
                    window.location.reload();
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

<style>
    .admin-nav {
        position: sticky;
        top: 0;
        z-index: 100;
        background: #030405;
        border-bottom: 1px solid rgba(255, 255, 255, 0.06);
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
    }

    .admin-nav-link {
        display: flex;
        align-items: center;
        gap: 0.35rem;
        padding: 0.35rem 0.65rem;
        text-decoration: none;
        color: #71717a;
        font-size: 0.75rem;
        font-family: var(--font-body);
        letter-spacing: 0.05em;
        border: 1px solid transparent;
        transition: all 0.15s ease;
        cursor: pointer;
        background: none;
    }

    .admin-nav-link:hover {
        color: #e4e4e7;
        border-color: rgba(255, 255, 255, 0.08);
    }

    .admin-nav-link.active {
        color: var(--accent);
        border-color: rgba(255, 68, 0, 0.2);
        background: rgba(255, 68, 0, 0.06);
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
    }

    .admin-nav-logout:hover {
        color: #ef4444;
        border-color: rgba(239, 68, 68, 0.2);
    }

    .admin-page {
        padding-top: 1rem;
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
