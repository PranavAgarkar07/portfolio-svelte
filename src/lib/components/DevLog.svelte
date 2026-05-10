<script lang="ts">
    import { onMount } from "svelte";
    import { Skeleton } from "$lib/components/ui";
    import type { DevLogResponse } from "$lib/types";
    let status = $state("Fetching latest commit...");
    let lastUpdate = $state("");
    let source = $state<"cache" | "live" | "stale-cache" | null>(null);
    let loading = $state(true);
    const API_URL = (import.meta.env.VITE_API_URL || "").replace(/\/$/, "") + "/api/status";

    onMount(async () => {
        const fetchStatus = async (timeoutMs: number, attemptNumber: number) => {
            const controller = new AbortController();
            const timer = setTimeout(() => controller.abort(), timeoutMs);
            try {
                const res = await fetch(API_URL, { signal: controller.signal });
                clearTimeout(timer);
                const data: DevLogResponse = await res.json();
                return data;
            } catch (e) {
                clearTimeout(timer);
                console.error("[DevLog] Fetch failed:", {
                    attempt: attemptNumber,
                    timeout: timeoutMs,
                    error: e instanceof Error ? e.message : String(e),
                });
                throw e;
            }
        };

        try {
            const data = await fetchStatus(8000, 1);
            status = data.summary;
            lastUpdate = data.last_update;
            source = data.source;
            loading = false;
        } catch {
            status = "Waking up server...";
            try {
                const data = await fetchStatus(30000, 2);
                status = data.summary;
                lastUpdate = data.last_update;
                source = data.source;
                loading = false;
            } catch {
                status = "System Offline";
                loading = false;
            }
        }
    });

    function formatSource(src: string) {
        switch (src) {
            case "live": return "LIVE";
            case "cache": return "CACHED";
            case "stale-cache": return "STALE";
            default: return src.toUpperCase();
        }
    }
</script>

<div class="status-panel" class:offline={status === "System Offline"}>
    <div class="panel-header">
        <span class="panel-label">STATUS</span>
        <div class="pulse-indicator" class:loading class:online={!loading && status !== "System Offline"} class:offline={status === "System Offline"}>
            <span class="pulse-dot"></span>
            <span class="pulse-label">
                {#if loading}
                    BOOT
                {:else if status === "System Offline"}
                    OFFLINE
                {:else}
                    ONLINE
                {/if}
            </span>
        </div>
    </div>
    <div class="panel-divider"></div>
    <div class="panel-body">
        {#if loading}
            <div class="skeleton-group">
                <Skeleton width="80%" />
                <Skeleton width="50%" />
            </div>
        {:else}
            <div class="status-message">{status}</div>
            <div class="status-meta">
                {#if lastUpdate}
                    <span class="meta-item">
                        <span class="meta-key">UPDATED</span>
                        <span class="meta-value">{lastUpdate}</span>
                    </span>
                {/if}
                {#if source}
                    <span class="meta-item">
                        <span class="meta-key">SOURCE</span>
                        <span class="meta-value source" class:live={source === "live"} class:cached={source === "cache"}>{formatSource(source)}</span>
                    </span>
                {/if}
            </div>
        {/if}
    </div>
</div>

<style>
    .status-panel {
        display: flex;
        flex-direction: column;
        width: 100%;
        max-width: 640px;
        background: #050505;
        border: 1px solid var(--grid-line);
        border-left: 3px solid var(--accent);
        border-right: 1px solid var(--grid-line);
        font-family: var(--font-body);
        transition: border-color 0.3s ease;
    }
    .status-panel.offline {
        border-left-color: #ff4444;
    }
    :global(body.light-mode) .status-panel {
        background: rgba(255, 255, 255, 0.95);
        border-color: rgba(0, 0, 0, 0.12);
    }
    :global(body.light-mode) .status-panel.offline {
        border-left-color: #ff4444;
    }
    .panel-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 0.5rem 1rem;
    }
    .panel-label {
        font-weight: 700;
        font-size: 0.65rem;
        letter-spacing: 0.15em;
        color: var(--accent);
    }
    :global(body.light-mode) .panel-label {
        color: var(--accent);
    }
    .panel-divider {
        height: 1px;
        background: var(--grid-line);
        margin: 0 1rem;
    }
    :global(body.light-mode) .panel-divider {
        background: rgba(0, 0, 0, 0.1);
    }
    .panel-body {
        padding: 0.75rem 1rem 1rem;
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
    }
    .status-message {
        font-size: 0.85rem;
        line-height: 1.6;
        color: var(--text-secondary);
        word-wrap: break-word;
    }
    :global(body.light-mode) .status-message {
        color: #1a1a1a;
        font-weight: 500;
    }
    .status-meta {
        display: flex;
        gap: 1.25rem;
        flex-wrap: wrap;
    }
    .meta-item {
        display: flex;
        align-items: center;
        gap: 0.4rem;
        font-size: 0.6rem;
        letter-spacing: 0.1em;
    }
    .meta-key {
        color: #555;
        font-weight: 600;
    }
    :global(body.light-mode) .meta-key {
        color: #888;
    }
    .meta-value {
        color: var(--text-secondary);
        font-weight: 600;
    }
    :global(body.light-mode) .meta-value {
        color: #333;
    }
    .meta-value.source.live {
        color: #00ffaa;
    }
    :global(body.light-mode) .meta-value.source.live {
        color: #00aa6e;
    }
    .meta-value.source.cached {
        color: #ffaa00;
    }
    :global(body.light-mode) .meta-value.source.cached {
        color: #cc8800;
    }
    .pulse-indicator {
        display: flex;
        align-items: center;
        gap: 6px;
        font-size: 0.6rem;
        letter-spacing: 0.15em;
    }
    .pulse-dot {
        width: 5px;
        height: 5px;
        border-radius: 50%;
        background: #555;
        transition: background 0.3s ease;
    }
    .pulse-indicator.loading .pulse-dot {
        background: #ffaa00;
        animation: blink 1s step-end infinite;
    }
    .pulse-indicator.online .pulse-dot {
        background: #00ffaa;
        box-shadow: 0 0 6px rgba(0, 255, 170, 0.4);
    }
    .pulse-indicator.offline .pulse-dot {
        background: #ff4444;
    }
    .pulse-label {
        color: var(--text-secondary);
    }
    :global(body.light-mode) .pulse-label {
        color: #555;
    }
    @keyframes blink {
        50% { opacity: 0; }
    }
</style>
