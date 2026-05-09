<script lang="ts">
    import { onMount } from "svelte";

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
            const res = await fetch(`${BASE}/api/contact/messages?key=${encodeURIComponent(storedKey)}`);
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
            const res = await fetch(`${BASE}/api/contact/messages/${id}/read?key=${encodeURIComponent(storedKey)}`, {
                method: "PATCH",
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
        editingKey = false;
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

<div class="admin-page">
    {#if !storedKey}
        <div class="key-entry">
            <h1>Admin — Contact Messages</h1>
            <input
                class="key-input"
                type="password"
                placeholder="Enter admin key"
                bind:value={key}
                onkeydown={(e) => e.key === "Enter" && key && saveKey(key)}
            />
            <button class="btn btn-primary key-submit" onclick={() => key && saveKey(key)} disabled={!key}>
                Access
            </button>
        </div>
    {:else}
        <div class="admin-header">
            <h1>Contact Messages</h1>
            <span class="key-link" role="button" tabindex="0" onclick={clearKey} onkeydown={(e) => e.key === "Enter" && clearKey()}>Change Key</span>
        </div>

        {#if loading}
            <div class="skeleton-row"></div>
            <div class="skeleton-row"></div>
            <div class="skeleton-row"></div>
        {:else if error}
            <div class="error-state">{error}</div>
        {:else if messages.length === 0}
            <div class="empty-state">No messages yet.</div>
        {:else}
            <table>
                <thead>
                    <tr>
                        <th>#</th>
                        <th>Name</th>
                        <th>Message</th>
                        <th>Date</th>
                        <th>Status</th>
                        <th>Action</th>
                    </tr>
                </thead>
                <tbody>
                    {#each messages as msg, i}
                        <tr>
                            <td>{i + 1}</td>
                            <td>
                                <div class="name-cell">{msg.name || msg.sender_name || "Unknown"}</div>
                                <div class="email-sub">{msg.email || msg.sender_email || ""}</div>
                            </td>
                            <td>{truncate(msg.message || msg.content || "", 100)}</td>
                            <td>{formatDate(msg.created_at || msg.created || msg.date || "")}</td>
                            <td>
                                {#if msg.is_read}
                                    <span class="badge-read">Read</span>
                                {:else}
                                    <span class="badge-new">New</span>
                                {/if}
                            </td>
                            <td>
                                {#if !msg.is_read}
                                    <button class="mark-read-btn" onclick={() => markRead(msg.id)}>Mark Read</button>
                                {/if}
                            </td>
                        </tr>
                    {/each}
                </tbody>
            </table>
        {/if}
    {/if}
</div>

<style>
    .admin-page { max-width: 1100px; margin: 0 auto; padding: 2rem; }
    .admin-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem; }
    .admin-header h1 { font-size: 1.3rem; text-transform: uppercase; letter-spacing: 0.05em; color: var(--text-primary); border-left: 3px solid var(--accent); padding-left: 1rem; }
    .key-link { color: #555; font-size: 0.75rem; cursor: pointer; text-decoration: underline; }
    .key-link:hover { color: var(--accent); }
    .key-entry { display: flex; flex-direction: column; align-items: center; justify-content: center; min-height: 60vh; gap: 1rem; }
    .key-entry h1 { font-size: 1.1rem; letter-spacing: 0.1em; }
    .key-input { width: 300px; max-width: 90%; background: #0a0a0a; border: 1px solid #222; color: var(--text-primary); padding: 0.75rem 1rem; font-family: var(--font-body); font-size: 0.85rem; outline: none; }
    .key-input:focus { border-color: var(--accent); }
    table { width: 100%; border-collapse: collapse; font-size: 0.8rem; }
    th { text-align: left; padding: 0.75rem 0.5rem; border-bottom: 1px solid var(--grid-line); color: var(--text-secondary); font-weight: 600; font-size: 0.65rem; letter-spacing: 0.1em; text-transform: uppercase; }
    td { padding: 0.75rem 0.5rem; border-bottom: 1px solid var(--grid-line); color: var(--text-primary); }
    tr:hover { background: rgba(255,255,255,0.02); }
    td:first-child { color: #555; font-size: 0.7rem; }
    .name-cell { font-weight: 500; }
    .email-sub { color: #555; font-size: 0.7rem; margin-top: 2px; }
    .badge-new { display: inline-block; padding: 2px 8px; font-size: 0.65rem; letter-spacing: 0.08em; background: rgba(255,68,0,0.15); color: var(--accent); border: 1px solid var(--accent); }
    .badge-read { display: inline-block; padding: 2px 8px; font-size: 0.65rem; letter-spacing: 0.08em; color: #555; border: 1px solid #333; }
    .mark-read-btn { font-size: 0.65rem; font-family: var(--font-body); background: transparent; border: 1px solid #333; color: #555; padding: 4px 10px; cursor: pointer; transition: all 0.2s; text-transform: uppercase; letter-spacing: 0.05em; }
    .mark-read-btn:hover { border-color: var(--accent); color: var(--accent); }
    .empty-state { text-align: center; padding: 4rem 0; color: #555; font-size: 0.9rem; }
    .error-state { color: #ff4444; font-size: 0.85rem; padding: 0.5rem 0; }
    .skeleton-row { height: 3rem; background: linear-gradient(90deg, #1a1a1a 25%, #2a2a2a 50%, #1a1a1a 75%); background-size: 200% 100%; animation: shimmer 1.5s ease-in-out infinite; margin-bottom: 2px; border-radius: 0; }
    @keyframes shimmer { 0% { background-position: 200% 0; } 100% { background-position: -200% 0; } }
    @media (max-width: 768px) { .admin-page { padding: 1rem; } table { font-size: 0.7rem; } th, td { padding: 0.5rem 0.3rem; } }
</style>
