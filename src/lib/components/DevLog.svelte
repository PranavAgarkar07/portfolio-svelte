<script lang="ts">
    import { onMount } from "svelte";
    let status = $state("Fetching latest commit...");
    let lastUpdate = $state("");
    let loading = $state(true);

    // Use deployed URL if valid, otherwise fallback to localhost
    const API_URL = import.meta.env.VITE_API_URL
        ? `${import.meta.env.VITE_API_URL}/api/status`
        : "http://localhost:8080/api/status";

    // Industrial/Aerospace Theme Colors (Orange Scale)
    const LEVEL_COLORS = [
        "rgba(255, 255, 255, 0.03)", // Level 0: Subtle surface
        "rgba(255, 68, 0, 0.2)", // Level 1: Faint Orange
        "rgba(255, 68, 0, 0.4)", // Level 2: Visible Orange
        "rgba(255, 68, 0, 0.7)", // Level 3: Strong Orange
        "#ff4400", // Level 4: Full Accent
    ];
    // Light Mode equivalents handled via CSS variables if needed, but these usually work well on dark.

    // Mock data for the "Past Days" to create the graph look
    const history = [1, 3, 0, 2, 4, 1, 3]; // 0-4 intensity levels

    onMount(async () => {
        try {
            const res = await fetch(API_URL);
            if (!res.ok) throw new Error("Network response was not ok");
            const data = await res.json();

            setTimeout(() => {
                status = data.summary;
                lastUpdate = "Today"; // GitHub style
                loading = false;
            }, 600);
        } catch (e) {
            console.error(e);
            status = "System Offline";
            loading = false;
        }
    });

    function getIntensityColor(level: number) {
        return LEVEL_COLORS[level] || LEVEL_COLORS[0];
    }
</script>

<div class="github-activity-strip">
    <div class="graph-label">Contribution Activity</div>

    <div class="graph-grid">
        <!-- Render past days (static visual) -->
        {#each history as level, i}
            <div
                class="contribution-cell"
                style="background-color: {getIntensityColor(level)}"
                title="Activity level: {level}"
            ></div>
        {/each}

        <!-- Today's Live Cell (Static) -->
        <div class="contribution-cell live"></div>
    </div>

    <div class="live-status">
        <span class="status-label">STATUS:</span>
        <span class="status-text">{status}</span>
    </div>

    <!-- AI Watermark -->
    <div class="ai-watermark">AI GEN</div>
</div>

<style>
    /* GitHub Grid Container */
    .github-activity-strip {
        position: relative; /* For watermark */
        /* overflow: visible; Default */
        display: flex;
        flex-direction: row; /* Horizontal for Skills Footer */
        align-items: center;
        gap: 16px;
        padding: 0.75rem 1rem;
        background: #050505; /* Deep Jet Black (Solid) */
        border: 1px solid var(--grid-line);
        border-right: 3px solid var(--accent); /* Technical Accent Edge */
        border-radius: 0px; /* Brutalist */
        width: fit-content;
        margin-top: 0;
        backdrop-filter: none; /* Removed blur for solid bg */
        box-shadow: 0 4px 20px rgba(0, 0, 0, 0.6); /* Deeper shadow */
    }

    /* AI Watermark Sticker */
    .ai-watermark {
        position: absolute;
        top: -6px;
        right: -6px;
        transform: rotate(25deg); /* Natural sticker angle */

        background: var(--accent);
        color: #000;

        font-family: var(--font-body);
        font-size: 0.45rem;
        font-weight: 800;
        letter-spacing: 0.05em;

        padding: 2px 6px;
        border: 2px solid rgba(255, 255, 255, 0.8); /* Die-cut effect */
        border-radius: 2px;

        box-shadow: 2px 2px 5px rgba(0, 0, 0, 0.5); /* Texture depth */

        pointer-events: none;
        z-index: 20;
        white-space: nowrap;
        text-transform: uppercase;
    }

    /* Light Mode Sticker: Darker border for contrast */
    :global(body.light-mode) .ai-watermark {
        border-color: #000;
        color: #000;
    }

    /* Light Mode Overrides for Container */
    :global(body.light-mode) .github-activity-strip {
        background: rgba(255, 255, 255, 0.9);
        border-color: rgba(0, 0, 0, 0.15);
        border-right: 3px solid var(--accent);
        box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
    }

    .graph-label {
        font-family: var(--font-body); /* JetBrains Mono */
        font-size: 0.65rem;
        color: var(--text-secondary);
        letter-spacing: 0.1em;
        font-weight: 700;
        opacity: 0.8;
    }

    :global(body.light-mode) .graph-label {
        color: #333; /* Darker clean text */
        opacity: 1;
    }

    .graph-grid {
        display: flex;
        gap: 2px; /* Tight technical gap */
    }

    .contribution-cell {
        width: 10px; /* Precise small blocks */
        height: 10px;
        border-radius: 0px; /* Sharp */
        background-color: rgba(255, 255, 255, 0.03);
    }

    :global(body.light-mode) .contribution-cell {
        background-color: #e1e4e8; /* Visible grey on white */
    }

    /* Today's Cell (Static Industrial) */
    .contribution-cell.live {
        background-color: var(--accent); /* International Orange */
        box-shadow: 0 0 6px var(--accent-glow);
    }

    /* Live Status Text */
    .live-status {
        display: flex;
        gap: 6px;
        align-items: center;
        font-family: var(--font-body);
        font-size: 0.7rem;
        color: var(--text-primary);
        border-left: 1px solid var(--grid-line);
        padding-left: 16px;
    }

    .status-label {
        font-weight: 700;
        color: var(--accent); /* Industrial Orange */
    }

    .status-text {
        color: var(--text-secondary);
        white-space: normal; /* Allow wrap */
        max-width: none; /* Show full text */
        overflow: visible;
    }

    :global(body.light-mode) .status-text {
        color: #111;
        font-weight: 500;
    }

    @media (max-width: 768px) {
        .github-activity-strip {
            flex-wrap: wrap;
            gap: 12px;
        }

        .live-status {
            display: flex; /* Show on mobile */
            border-left: none; /* Remove separator */
            border-top: 1px solid var(--grid-line); /* Stack with separator */
            padding-left: 0;
            padding-top: 8px;
            width: 100%;
        }

        .status-text {
            max-width: none; /* Full width on mobile */
        }
    }
</style>
