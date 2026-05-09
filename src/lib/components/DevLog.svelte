<script lang="ts">
    import { onMount } from "svelte";
    import type { DevLogResponse } from "$lib/types";
    let status = $state("Fetching latest commit...");
    let lastUpdate = $state("");
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
            setTimeout(() => {
                status = data.summary;
                lastUpdate = "Today";
                loading = false;
            }, 600);
        } catch {
            status = "Waking up server...";
            try {
                const data = await fetchStatus(30000, 2);
                status = data.summary;
                lastUpdate = "Today";
                loading = false;
            } catch (e) {
                console.error("[DevLog] Fetch failed:", {
                    attempt: 2,
                    timeout: 30000,
                    error: e instanceof Error ? e.message : String(e),
                });
                status = "System Offline";
                loading = false;
            }
        }
    });

</script>

<div class="devlog-strip">
    <div class="status-block">
        <span class="status-label">SYS::STATUS</span>
        {#if loading}
            <div class="skeleton-bar"></div>
        {:else}
            <span class="status-text">{status}</span>
        {/if}
    </div>
    <div class="pulse-indicator" class:loading class:online={!loading && status !== "System Offline"} class:offline={status === "System Offline"}>
        <span class="pulse-dot"></span>
        <span class="pulse-label">
            {#if loading}
                INITIALIZING
            {:else if status === "System Offline"}
                OFFLINE
            {:else}
                ONLINE
            {/if}
        </span>
    </div>
</div>

<style>
    .devlog-strip {
        display: flex;
        flex-direction: row;
        align-items: center;
        gap: 16px;
        padding: 0.75rem 1rem;
        background: #050505;
        border: 1px solid var(--grid-line);
        border-right: 3px solid var(--accent);
        border-radius: 0;
        width: fit-content;
        box-shadow: 0 4px 20px rgba(0, 0, 0, 0.6);
    }
    :global(body.light-mode) .devlog-strip {
        background: rgba(255, 255, 255, 0.9);
        border-color: rgba(0, 0, 0, 0.15);
        border-right: 3px solid var(--accent);
        box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
    }
    .status-block {
        display: flex;
        gap: 8px;
        align-items: center;
        font-family: var(--font-body);
        font-size: 0.7rem;
    }
    .status-label {
        font-weight: 700;
        color: var(--accent);
        letter-spacing: 0.1em;
    }
    .status-text {
        color: var(--text-secondary);
        white-space: normal;
    }
    :global(body.light-mode) .status-text {
        color: #111;
        font-weight: 500;
    }
    .skeleton-bar {
        width: 180px;
        height: 0.7rem;
        background: linear-gradient(90deg, #1a1a1a 25%, #2a2a2a 50%, #1a1a1a 75%);
        background-size: 200% 100%;
        animation: shimmer 1.5s ease-in-out infinite;
        border-radius: 2px;
    }
    :global(body.light-mode) .skeleton-bar {
        background: linear-gradient(90deg, #e0e0e0 25%, #f0f0f0 50%, #e0e0e0 75%);
        background-size: 200% 100%;
    }
    @keyframes shimmer {
        0% { background-position: 200% 0; }
        100% { background-position: -200% 0; }
    }
    .pulse-indicator {
        display: flex;
        align-items: center;
        gap: 6px;
        padding-left: 16px;
        border-left: 1px solid var(--grid-line);
        font-family: var(--font-body);
        font-size: 0.65rem;
        letter-spacing: 0.1em;
    }
    .pulse-dot {
        width: 6px;
        height: 6px;
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
    :global(body.light-mode) .pulse-indicator {
        border-left-color: rgba(0, 0, 0, 0.15);
    }
    :global(body.light-mode) .pulse-label {
        color: #555;
    }
    @keyframes blink {
        50% { opacity: 0; }
    }
    @media (max-width: 768px) {
        .devlog-strip {
            flex-wrap: wrap;
            gap: 12px;
        }
        .pulse-indicator {
            border-left: none;
            border-top: 1px solid var(--grid-line);
            padding-left: 0;
            padding-top: 8px;
            width: 100%;
        }
    }
</style>
