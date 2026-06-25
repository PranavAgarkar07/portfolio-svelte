<script lang="ts">
    import { onMount } from "svelte";

    const BASE = (import.meta.env.VITE_API_URL || "").replace(/\/$/, "");
    const ANALYTICS_URL = BASE + "/api/v1/analytics";

    interface ReferrerCount { referrer: string; count: number; }
    interface TargetCount { target: string; count: number; }
    interface CountryCount { country: string; count: number; }
    interface DashboardStats {
        total_views: number;
        unique_visitors: number;
        top_referrers: ReferrerCount[];
        top_targets: TargetCount[];
        country_breakdown: CountryCount[];
        avg_time_on_site: number;
        resume_downloads: number;
        form_submissions: number;
    }

    let key = $state("");
    let loading = $state(false);
    let error = $state("");
    let stats = $state<DashboardStats | null>(null);
    let since = $state("30");

    onMount(() => {
        const saved = localStorage.getItem("contact_admin_key");
        if (saved) {
            key = saved;
            loadDashboard();
        }
    });

    async function loadDashboard() {
        if (!key) return;
        loading = true;
        error = "";
        stats = null;
        try {
            const days = parseInt(since) || 30;
            const sinceParam = new Date(Date.now() - days * 86400000).toISOString();
            const res = await fetch(`${ANALYTICS_URL}/dashboard?since=${encodeURIComponent(sinceParam)}`, {
                headers: { "X-Analytics-Key": key },
            });
            if (!res.ok) {
                const text = await res.text();
                throw new Error(text || `HTTP ${res.status}`);
            }
            stats = await res.json();
        } catch (e: any) {
            error = e.message || "Failed to load dashboard";
        } finally {
            loading = false;
        }
    }

    function formatTime(seconds: number): string {
        if (!seconds || seconds === 0) return "—";
        if (seconds < 60) return `${Math.round(seconds)}s`;
        const m = Math.floor(seconds / 60);
        const s = Math.round(seconds % 60);
        return `${m}m ${s}s`;
    }

    function shortenUrl(url: string): string {
        try {
            const u = new URL(url);
            return u.hostname + (u.pathname.length > 1 ? u.pathname.substring(0, 30) + "..." : "");
        } catch {
            return url.length > 35 ? url.substring(0, 35) + "..." : url;
        }
    }
</script>

<div class="dashboard">
    <div class="dh">
        <h1 class="dh-title">Analytics Dashboard</h1>
        <div class="dh-controls">
            <select bind:value={since} onchange={loadDashboard} class="dh-select">
                <option value="7">Last 7 days</option>
                <option value="30">Last 30 days</option>
                <option value="90">Last 90 days</option>
                <option value="365">Last year</option>
            </select>
            <button class="dh-refresh" onclick={loadDashboard} disabled={loading}>
                &#x21bb; {loading ? "Loading..." : "Refresh"}
            </button>
        </div>
    </div>

    {#if !key}
        <div class="empty">
            <p>No admin key found. Set your key in localStorage (<code>contact_admin_key</code>) to view analytics.</p>
        </div>
    {:else if loading && !stats}
        <div class="loading">
            <div class="spinner"></div>
            <p>Loading dashboard...</p>
        </div>
    {:else if error}
        <div class="error">
            <p>{error}</p>
            <button class="dh-refresh" onclick={loadDashboard}>Retry</button>
        </div>
    {:else if stats}
        <div class="kc">
            <div class="kc-card">
                <span class="kc-label">Total Views</span>
                <span class="kc-value">{stats.total_views}</span>
            </div>
            <div class="kc-card">
                <span class="kc-label">Unique Visitors</span>
                <span class="kc-value">{stats.unique_visitors}</span>
            </div>
            <div class="kc-card">
                <span class="kc-label">Resume Downloads</span>
                <span class="kc-value">{stats.resume_downloads}</span>
            </div>
            <div class="kc-card">
                <span class="kc-label">Form Submissions</span>
                <span class="kc-value">{stats.form_submissions}</span>
            </div>
            <div class="kc-card">
                <span class="kc-label">Avg Time on Site</span>
                <span class="kc-value">{formatTime(stats.avg_time_on_site)}</span>
            </div>
        </div>

        <div class="grid">
            <div class="panel">
                <h2 class="panel-title">Top Referrers</h2>
                {#if stats.top_referrers.length === 0}
                    <p class="panel-empty">No referrer data</p>
                {:else}
                    <table class="table">
                        <thead>
                            <tr><th>Source</th><th class="num">Visits</th></tr>
                        </thead>
                        <tbody>
                            {#each stats.top_referrers as r}
                                <tr>
                                    <td class="break">{r.referrer || "(direct)"}</td>
                                    <td class="num">{r.count}</td>
                                </tr>
                            {/each}
                        </tbody>
                    </table>
                {/if}
            </div>

            <div class="panel">
                <h2 class="panel-title">Top Targets</h2>
                {#if stats.top_targets.length === 0}
                    <p class="panel-empty">No event data</p>
                {:else}
                    <table class="table">
                        <thead>
                            <tr><th>Target</th><th class="num">Clicks</th></tr>
                        </thead>
                        <tbody>
                            {#each stats.top_targets as t}
                                <tr>
                                    <td class="break">{t.target}</td>
                                    <td class="num">{t.count}</td>
                                </tr>
                            {/each}
                        </tbody>
                    </table>
                {/if}
            </div>

            <div class="panel">
                <h2 class="panel-title">Countries</h2>
                {#if stats.country_breakdown.length === 0}
                    <p class="panel-empty">No geography data</p>
                {:else}
                    <table class="table">
                        <thead>
                            <tr><th>Country</th><th class="num">Visits</th></tr>
                        </thead>
                        <tbody>
                            {#each stats.country_breakdown as c}
                                <tr>
                                    <td>{c.country || "Unknown"}</td>
                                    <td class="num">{c.count}</td>
                                </tr>
                            {/each}
                        </tbody>
                    </table>
                {/if}
            </div>
        </div>
    {/if}
</div>

<style>
    .dashboard {
        max-width: 1100px;
        margin: 0 auto;
        padding: 1.5rem 2rem 4rem;
    }

    .dh {
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 1rem;
        margin-bottom: 1.5rem;
        flex-wrap: wrap;
    }

    .dh-title {
        font-family: var(--font-heading);
        font-size: 1.1rem;
        font-weight: 700;
        letter-spacing: 0.08em;
        text-transform: uppercase;
        color: var(--accent);
        margin: 0;
    }

    .dh-controls {
        display: flex;
        align-items: center;
        gap: 0.5rem;
    }

    .dh-select {
        background: transparent;
        border: 1px solid rgba(255, 255, 255, 0.1);
        color: #e4e4e7;
        padding: 0.4rem 0.65rem;
        font-family: var(--font-body);
        font-size: 0.75rem;
        cursor: pointer;
        outline: none;
    }

    .dh-select:focus {
        border-color: var(--accent);
    }

    .dh-refresh {
        background: transparent;
        border: 1px solid rgba(255, 255, 255, 0.1);
        color: #e4e4e7;
        padding: 0.4rem 0.75rem;
        font-family: var(--font-body);
        font-size: 0.75rem;
        cursor: pointer;
        transition: border-color 0.15s ease, color 0.15s ease;
    }

    .dh-refresh:hover:not(:disabled) {
        border-color: var(--accent);
        color: var(--accent);
    }

    .dh-refresh:disabled {
        opacity: 0.4;
        cursor: not-allowed;
    }

    .kc {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(170px, 1fr));
        gap: 0.75rem;
        margin-bottom: 1.5rem;
    }

    .kc-card {
        border: 1px solid rgba(255, 255, 255, 0.06);
        padding: 1rem 1.25rem;
        display: flex;
        flex-direction: column;
        gap: 0.35rem;
    }

    .kc-label {
        font-family: var(--font-body);
        font-size: 0.65rem;
        font-weight: 600;
        letter-spacing: 0.1em;
        text-transform: uppercase;
        color: #71717a;
    }

    .kc-value {
        font-family: var(--font-heading);
        font-size: 1.65rem;
        font-weight: 700;
        color: #e4e4e7;
        line-height: 1;
    }

    .grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
        gap: 1rem;
    }

    .panel {
        border: 1px solid rgba(255, 255, 255, 0.06);
        padding: 1rem 1.25rem;
    }

    .panel-title {
        font-family: var(--font-heading);
        font-size: 0.7rem;
        font-weight: 700;
        letter-spacing: 0.1em;
        text-transform: uppercase;
        color: #71717a;
        margin: 0 0 0.75rem;
        padding-bottom: 0.5rem;
        border-bottom: 1px solid rgba(255, 255, 255, 0.06);
    }

    .panel-empty {
        color: #52525b;
        font-size: 0.8rem;
        font-family: var(--font-body);
    }

    .table {
        width: 100%;
        border-collapse: collapse;
    }

    .table th {
        font-family: var(--font-body);
        font-size: 0.6rem;
        font-weight: 600;
        letter-spacing: 0.1em;
        text-transform: uppercase;
        color: #52525b;
        text-align: left;
        padding: 0.35rem 0.5rem 0.35rem 0;
        border-bottom: 1px solid rgba(255, 255, 255, 0.04);
    }

    .table td {
        font-family: var(--font-body);
        font-size: 0.78rem;
        color: #d4d4d8;
        padding: 0.35rem 0.5rem 0.35rem 0;
        border-bottom: 1px solid rgba(255, 255, 255, 0.03);
    }

    .table th.num,
    .table td.num {
        text-align: right;
        padding-right: 0;
        width: 60px;
        font-variant-numeric: tabular-nums;
    }

    .table .break {
        word-break: break-all;
        max-width: 240px;
    }

    .empty, .loading, .error {
        text-align: center;
        padding: 3rem 1rem;
        color: #71717a;
        font-family: var(--font-body);
        font-size: 0.85rem;
    }

    .empty code {
        background: rgba(255, 255, 255, 0.05);
        padding: 0.1rem 0.35rem;
        font-size: 0.78rem;
    }

    .error {
        color: #ef4444;
    }

    .error .dh-refresh {
        margin-top: 0.75rem;
        border-color: rgba(239, 68, 68, 0.3);
        color: #ef4444;
    }

    .spinner {
        width: 24px;
        height: 24px;
        border: 2px solid rgba(255, 255, 255, 0.08);
        border-top-color: var(--accent);
        border-radius: 50%;
        animation: spin 0.7s linear infinite;
        margin: 0 auto 0.75rem;
    }

    @keyframes spin {
        to { transform: rotate(360deg); }
    }

    @media (max-width: 640px) {
        .dashboard {
            padding: 1rem 1rem 3rem;
        }
        .kc {
            grid-template-columns: repeat(2, 1fr);
        }
        .kc-value {
            font-size: 1.25rem;
        }
        .grid {
            grid-template-columns: 1fr;
        }
    }
</style>
