<script lang="ts">
    import { onMount } from "svelte";
    import { Input, Button, Skeleton, Tag } from "$lib/components/ui";

    const BASE = (import.meta.env.VITE_API_URL || "").replace(/\/$/, "");

    let key = $state("");
    let storedKey = $state("");
    let messages = $state<any[]>([]);
    let loading = $state(false);
    let error = $state("");
    onMount(() => {
        const params = new URLSearchParams(window.location.search);
        const urlKey = params.get("key");
        if (urlKey) {
            saveKey(urlKey);
            return;
        }
        const saved = localStorage.getItem("contact_admin_key");
        if (saved) {
            storedKey = saved;
        }
    });

    $effect(() => {
        if (storedKey) {
            loadMessages();
        }
    });

    async function loadMessages() {
        loading = true;
        error = "";
        try {
            const res = await fetch(`${BASE}/api/contact/messages`, { headers: { Authorization: `Bearer ${storedKey}` } });
            if (!res.ok) {
                const text = await res.text();
                throw new Error(text || `HTTP ${res.status}`);
            }
            const data = await res.json();
            messages = data.messages ?? data ?? [];
        } catch (e: any) {
            error = e.message || "Failed to load messages";
            messages = [];
        } finally {
            loading = false;
        }
    }

    async function markRead(id: number) {
        try {
            const res = await fetch(`${BASE}/api/contact/messages/${id}/read`, {
                method: "PATCH",
                headers: { Authorization: `Bearer ${storedKey}` },
            });
            if (!res.ok) {
                const text = await res.text();
                throw new Error(text || `HTTP ${res.status}`);
            }
            messages = messages.map((m) => (m.id === id ? { ...m, is_read: true } : m));
        } catch (e: any) {
            error = e.message || "Failed to mark as read";
        }
    }

    function saveKey(k: string) {
        localStorage.setItem("contact_admin_key", k);
        storedKey = k;
        key = "";
    }

    function clearKey() {
        localStorage.removeItem("contact_admin_key");
        storedKey = "";
    }

    function truncate(s: string, len: number): string {
        return s.length > len ? s.slice(0, len) + "..." : s;
    }

    function formatDate(iso: string): string {
        const d = new Date(iso);
        return d.toLocaleDateString("en-US", { month: "short", day: "numeric", year: "numeric", hour: "2-digit", minute: "2-digit" });
    }
</script>

<div class="terminal-container">
    {#if !storedKey}
        <div class="key-entry">
            <div class="key-entry-inner">
                <div class="key-accent-bar"></div>
                <h1 class="key-title">CONTROL PANEL</h1>
                <p class="key-subtitle">Authenticate to access message terminal</p>
                <div class="key-field">
                    <Input
                        type="password"
                        placeholder="ENTER ADMIN KEY"
                        bind:value={key}
                        onkeydown={(e: KeyboardEvent) => e.key === "Enter" && key && saveKey(key)}
                    />
                </div>
                <div class="key-submit">
                    <Button onclick={() => key && saveKey(key)} disabled={!key}>Access</Button>
                </div>
            </div>
        </div>
    {:else}
        <div class="panel">
            <div class="panel-header">
                <div class="panel-title">
                    <span class="panel-accent"></span>
                    <h1>CONTROL PANEL</h1>
                </div>
                <Button variant="ghost" size="sm" onclick={clearKey}>Change Key</Button>
            </div>

            <div class="stats-bar">
                <div class="stat">
                    <span class="stat-label">Total</span>
                    <span class="stat-value">{messages.length}</span>
                </div>
                <div class="stat">
                    <span class="stat-label">Unread</span>
                    <span class="stat-value accent">{messages.filter(m => !m.is_read).length}</span>
                </div>
            </div>

            {#if loading}
                <div class="loading-grid">
                    <Skeleton variant="card" count={3} />
                </div>
            {:else if error}
                <div class="error-state">
                    <span class="error-icon">&#9888;</span>
                    {error}
                </div>
            {:else if messages.length === 0}
                <div class="empty-state">
                    <span class="empty-text">No incoming transmissions.</span>
                    <span class="empty-cursor"></span>
                </div>
            {:else}
                <div class="table-wrapper">
                    <table>
                        <thead>
                            <tr>
                                <th class="col-num"></th>
                                <th class="col-name">Sender</th>
                                <th class="col-topic">Topic</th>
                                <th class="col-msg">Message</th>
                                <th class="col-date">Date</th>
                                <th class="col-status">Status</th>
                                <th class="col-action"></th>
                            </tr>
                        </thead>
                        <tbody>
                            {#each messages as msg, i}
                                <tr
                                    class="msg-row"
                                    class:unread={!msg.is_read}
                                    style="--i: {i}"
                                >
                                    <td class="row-num">{i + 1}</td>
                                    <td>
                                        <div class="name-cell">{msg.name || msg.sender_name || "Unknown"}</div>
                                        <div class="email-sub">{msg.email || msg.sender_email || ""}</div>
                                    </td>
                                    <td>
                                        <Tag variant="muted">{msg.topic || "general"}</Tag>
                                    </td>
                                    <td>
                                        <div class="msg-cell">
                                            <span class="msg-text">{truncate(msg.message || msg.content || "", 100)}</span>
                                            <span class="msg-fade"></span>
                                        </div>
                                    </td>
                                    <td class="date-cell">{formatDate(msg.created_at || msg.created || msg.date || "")}</td>
                                    <td>
                                        {#if msg.is_read}
                                            <span class="badge badge-read">Read</span>
                                        {:else}
                                            <span class="badge badge-new">New</span>
                                        {/if}
                                    </td>
                                    <td>
                                        {#if !msg.is_read}
                                            <Button variant="ghost" size="sm" onclick={() => markRead(msg.id)}>Mark Read</Button>
                                        {/if}
                                    </td>
                                </tr>
                            {/each}
                        </tbody>
                    </table>
                </div>
            {/if}
        </div>
    {/if}
</div>

<style>
    .terminal-container {
        max-width: 1100px;
        margin: 0 auto;
        padding: 2rem;
        position: relative;
        overflow: hidden;
    }
    .terminal-container::after {
        content: "";
        position: absolute;
        inset: 0;
        pointer-events: none;
        z-index: 0;
        background: repeating-linear-gradient(
            0deg,
            transparent,
            transparent 2px,
            rgba(255,255,255,0.015) 2px,
            rgba(255,255,255,0.015) 4px
        );
    }

    /* ── Key Entry ── */
    .key-entry {
        display: flex;
        align-items: center;
        justify-content: center;
        min-height: 70vh;
        position: relative;
        z-index: 1;
    }
    .key-entry-inner {
        display: flex;
        flex-direction: column;
        align-items: flex-start;
        gap: 0.75rem;
        padding: 3rem 3.5rem;
        background: var(--surface-dark);
        border: 1px solid rgba(255,255,255,0.06);
        position: relative;
        animation: pulseGlow 3s ease-in-out infinite;
    }
    .key-accent-bar {
        width: 3px;
        height: 2.5rem;
        background: var(--accent);
        position: absolute;
        left: 0;
        top: 3rem;
    }
    .key-title {
        font-family: var(--font-heading);
        font-size: 1.6rem;
        letter-spacing: 0.15em;
        color: var(--text-primary);
        font-weight: 600;
        margin: 0;
    }
    .key-subtitle {
        font-size: 0.8rem;
        color: var(--text-secondary);
        letter-spacing: 0.05em;
        margin: 0 0 0.5rem;
    }
    .key-field {
        position: relative;
        width: 320px;
        max-width: 100%;
    }
    .key-submit {
        margin-top: 0.25rem;
    }

    /* ── Panel ── */
    .panel {
        position: relative;
        z-index: 1;
        border: 1px solid rgba(255,255,255,0.06);
        background: var(--surface-dark);
        padding: 1.75rem 2rem;
    }
    .panel-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 1.5rem;
    }
    .panel-title {
        display: flex;
        align-items: center;
        gap: 0.75rem;
    }
    .panel-accent {
        width: 3px;
        height: 1.5rem;
        background: var(--accent);
        flex-shrink: 0;
    }
    .panel-title h1 {
        font-family: var(--font-heading);
        font-size: 1.1rem;
        letter-spacing: 0.15em;
        color: var(--text-primary);
        font-weight: 600;
    }
    /* ── Stats Bar ── */
    .stats-bar {
        display: flex;
        gap: 2rem;
        margin-bottom: 1.75rem;
        padding: 0.75rem 1rem;
        background: rgba(255,255,255,0.02);
        border: 1px solid rgba(255,255,255,0.04);
    }
    .stat {
        display: flex;
        align-items: baseline;
        gap: 0.5rem;
    }
    .stat-label {
        font-size: 0.6rem;
        letter-spacing: 0.15em;
        text-transform: uppercase;
        color: var(--text-secondary);
        font-weight: 500;
    }
    .stat-value {
        font-size: 1rem;
        font-weight: 600;
        color: var(--text-primary);
        font-family: var(--font-body);
        letter-spacing: 0.05em;
    }
    .stat-value.accent {
        color: var(--accent);
    }

    /* ── Table ── */
    .table-wrapper {
        overflow-x: auto;
    }
    table {
        width: 100%;
        border-collapse: collapse;
        font-size: 0.8rem;
    }
    th {
        text-align: left;
        padding: 0.6rem 0.75rem;
        border-bottom: 1px solid var(--grid-line);
        color: var(--text-secondary);
        font-weight: 600;
        font-size: 0.6rem;
        letter-spacing: 0.12em;
        text-transform: uppercase;
    }
    td {
        padding: 0.75rem;
        border-bottom: 1px solid var(--grid-line);
        color: var(--text-primary);
        vertical-align: top;
    }
    .col-num { width: 2.5rem; }
    .col-name { width: 16%; }
    .col-topic { width: 10%; }
    .col-msg {  }
    .col-date { width: 14%; }
    .col-status { width: 8%; }
    .col-action { width: 9%; }

    /* ── Row Animation ── */
    .msg-row {
        transition: background 0.25s;
        animation: slideIn 0.35s ease-out both;
        animation-delay: calc(var(--i) * 0.05s);
    }
    .msg-row:hover {
        background: rgba(255,255,255,0.04);
        box-shadow: 0 0 0 1px rgba(255,255,255,0.04);
    }
    .msg-row.unread {
        border-left: 2px solid var(--accent);
    }
    .msg-row:not(.unread) {
        border-left: 2px solid transparent;
    }

    .row-num {
        color: rgba(255,255,255,0.2);
        font-size: 0.7rem;
        font-variant-numeric: tabular-nums;
    }

    /* ── Name / Email ── */
    .name-cell {
        font-weight: 500;
        font-size: 0.8rem;
        color: var(--text-primary);
    }
    .email-sub {
        color: var(--text-secondary);
        font-size: 0.65rem;
        margin-top: 2px;
        opacity: 0.7;
    }

    /* ── Message Cell ── */
    .msg-cell {
        position: relative;
        max-width: 320px;
    }
    .msg-text {
        font-size: 0.75rem;
        line-height: 1.5;
        color: var(--text-primary);
        opacity: 0.85;
    }
    .msg-fade {
        position: absolute;
        right: 0;
        top: 0;
        bottom: 0;
        width: 2rem;
        background: linear-gradient(to right, transparent, var(--surface-dark));
        pointer-events: none;
    }

    /* ── Date ── */
    .date-cell {
        font-size: 0.7rem;
        font-family: var(--font-body);
        color: var(--text-secondary);
        white-space: nowrap;
        font-variant-numeric: tabular-nums;
    }

    /* ── Badges ── */
    .badge {
        display: inline-block;
        padding: 2px 10px;
        font-size: 0.6rem;
        letter-spacing: 0.1em;
        text-transform: uppercase;
        font-weight: 500;
    }
    .badge-new {
        background: rgba(255,68,0,0.12);
        color: var(--accent);
        border: 1px solid var(--accent);
        box-shadow: 0 0 8px rgba(255,68,0,0.25);
    }
    .badge-read {
        color: rgba(255,255,255,0.3);
        border: 1px solid rgba(255,255,255,0.08);
    }

    /* ── Empty State ── */
    .empty-state {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 0.25rem;
        padding: 4rem 0;
        color: var(--text-secondary);
        font-size: 0.9rem;
    }
    .empty-text {
        letter-spacing: 0.05em;
    }
    .empty-cursor {
        display: inline-block;
        width: 2px;
        height: 1em;
        background: var(--accent);
        animation: blink 1s step-end infinite;
    }
    @keyframes blink {
        0%, 100% { opacity: 1; }
        50% { opacity: 0; }
    }

    /* ── Error State ── */
    .error-state {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        color: #ff4444;
        font-size: 0.8rem;
        padding: 0.75rem 0;
        letter-spacing: 0.03em;
    }
    .error-icon {
        font-size: 1rem;
        flex-shrink: 0;
    }

    .loading-grid {
        display: flex;
        flex-direction: column;
        gap: 2px;
    }

    /* ── Keyframes ── */
    @keyframes slideIn {
        from {
            opacity: 0;
            transform: translateY(10px);
        }
        to {
            opacity: 1;
            transform: translateY(0);
        }
    }
    @keyframes pulseGlow {
        0%, 100% {
            box-shadow: 0 0 0 rgba(255,68,0,0);
        }
        50% {
            box-shadow: 0 0 30px rgba(255,68,0,0.08), inset 0 0 30px rgba(255,68,0,0.02);
        }
    }

    /* ── Reduced Motion ── */
    @media (prefers-reduced-motion: reduce) {
        .key-entry-inner {
            animation: none;
        }
        .msg-row {
            animation: none;
        }
        .msg-row {
            transition: none;
        }
        .empty-cursor {
            animation: none;
            opacity: 1;
        }
        .msg-row:hover {
            box-shadow: none;
        }
        .key-link,
        .msg-row {
            transition: none;
        }
    }

    /* ── Responsive ── */
    @media (max-width: 768px) {
        .terminal-container { padding: 1rem; }
        .panel { padding: 1rem; }
        table { font-size: 0.7rem; }
        th, td { padding: 0.5rem; }
        .key-entry-inner { padding: 2rem 1.5rem; }
        .stats-bar { gap: 1rem; }
        .msg-cell { max-width: 160px; }
    }
</style>
