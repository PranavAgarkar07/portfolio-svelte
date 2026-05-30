<script lang="ts">
    import { onMount } from "svelte";
    import { Input, Button, Skeleton } from "$lib/components/ui";
    import type { Badge } from "$lib/types";

    const BASE = (import.meta.env.VITE_API_URL || "").replace(/\/$/, "");

    let key = $state("");
    let storedKey = $state("");
    let badges = $state<Badge[]>([]);
    let loading = $state(false);
    let error = $state("");

    let showForm = $state(false);
    let editingId = $state<number | null>(null);
    let uploading = $state(false);
    let successMsg = $state("");
    let form = $state({
        name: "",
        image_url: "",
        credential_url: "",
        rarity: "common",
        category: "",
        important: false,
        display_order: 0,
    });

    function toNum(v: unknown, fallback = 0): number {
        const n = typeof v === "string" ? parseInt(v, 10) : Number(v);
        return isNaN(n) ? fallback : n;
    }

    function flash(msg: string) {
        successMsg = msg;
        setTimeout(() => successMsg = "", 3000);
    }

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
            loadBadges();
        }
    });

    async function loadBadges() {
        loading = true;
        error = "";
        try {
            const res = await fetch(`${BASE}/api/admin/badges?key=${encodeURIComponent(storedKey)}`);
            if (!res.ok) {
                const text = await res.text();
                throw new Error(text || `HTTP ${res.status}`);
            }
            const data = await res.json();
            badges = data.badges ?? data ?? [];
        } catch (e: any) {
            error = e.message || "Failed to load badges";
            badges = [];
        } finally {
            loading = false;
        }
    }

    async function saveBadge() {
        if (!form.name.trim()) return;

        const body = {
            name: form.name,
            image_url: form.image_url,
            credential_url: form.credential_url,
            rarity: form.rarity,
            category: form.category,
            important: form.important,
            display_order: toNum(form.display_order, badges.length),
        };

        try {
            const url = editingId
                ? `${BASE}/api/admin/badges/${editingId}?key=${encodeURIComponent(storedKey)}`
                : `${BASE}/api/admin/badges?key=${encodeURIComponent(storedKey)}`;
            const method = editingId ? "PUT" : "POST";

            const res = await fetch(url, {
                method,
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(body),
            });

            if (!res.ok) {
                const text = await res.text();
                throw new Error(text || `HTTP ${res.status}`);
            }

            flash(editingId ? "Badge updated." : "Badge created.");
            resetForm();
            await loadBadges();
        } catch (e: any) {
            error = e.message || "Failed to save badge";
        }
    }

    async function uploadImage(e: Event) {
        const input = e.target as HTMLInputElement;
        const file = input.files?.[0];
        if (!file) return;

        uploading = true;
        error = "";
        try {
            const fd = new FormData();
            fd.append("image", file);
            const res = await fetch(`${BASE}/api/admin/certificates/upload?key=${encodeURIComponent(storedKey)}`, {
                method: "POST",
                body: fd,
            });
            if (!res.ok) {
                const txt = await res.text();
                throw new Error(txt || `HTTP ${res.status}`);
            }
            const data = await res.json();
            form.image_url = data.url;
        } catch (e: any) {
            error = e.message || "Upload failed";
        } finally {
            uploading = false;
            input.value = "";
        }
    }

    async function deleteBadge(id: number) {
        if (!confirm("Delete this badge?")) return;

        try {
            const res = await fetch(`${BASE}/api/admin/badges/${id}?key=${encodeURIComponent(storedKey)}`, {
                method: "DELETE",
            });

            if (!res.ok) {
                const text = await res.text();
                throw new Error(text || `HTTP ${res.status}`);
            }

            flash("Badge deleted.");
            await loadBadges();
        } catch (e: any) {
            error = e.message || "Failed to delete badge";
        }
    }

    function editBadge(b: Badge) {
        editingId = b.id;
        form = {
            name: b.name,
            image_url: b.image_url || "",
            credential_url: b.credential_url || "",
            rarity: b.rarity || "common",
            category: b.category || "",
            important: b.important,
            display_order: b.display_order,
        };
        showForm = true;
    }

    function resetForm() {
        showForm = false;
        editingId = null;
        form = { name: "", image_url: "", credential_url: "", rarity: "common", category: "", important: false, display_order: 0 };
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

    function formatDate(iso: string): string {
        const d = new Date(iso);
        return d.toLocaleDateString("en-US", { month: "short", day: "numeric", year: "numeric", hour: "2-digit", minute: "2-digit" });
    }

    const rarityColors: Record<string, string> = {
        common: "var(--accent)",
        uncommon: "#a78bfa",
        rare: "#fbbf24",
    };
</script>

<div class="terminal-container">
    {#if !storedKey}
        <div class="key-entry">
            <div class="key-entry-inner">
                <div class="key-accent-bar"></div>
                <h1 class="key-title">BADGES ADMIN</h1>
                <p class="key-subtitle">Authenticate to manage GSSoC badges</p>
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
                    <h1>BADGE CONTROL</h1>
                </div>
                <div class="panel-actions">
                    <Button variant="ghost" size="sm" onclick={() => { resetForm(); showForm = !showForm; }}>
                        {showForm ? "Cancel" : "Add Badge"}
                    </Button>
                    <Button variant="ghost" size="sm" onclick={clearKey}>Change Key</Button>
                </div>
            </div>

            {#if showForm}
                <div class="form-panel">
                    <h3 class="form-title">{editingId ? "Edit Badge" : "New Badge"}</h3>
                    <div class="form-grid">
                        <div class="form-field">
                            <label class="form-label">Name *</label>
                            <Input bind:value={form.name} placeholder="e.g. Ambassador" />
                        </div>
                        <div class="form-field">
                            <label class="form-label">Rarity</label>
                            <select bind:value={form.rarity} class="rarity-select">
                                <option value="common">Common</option>
                                <option value="uncommon">Uncommon</option>
                                <option value="rare">Rare</option>
                            </select>
                        </div>
                        <div class="form-field">
                            <label class="form-label">Category</label>
                            <Input bind:value={form.category} placeholder="e.g. GSSoC 2026, GitHub" />
                        </div>
                        <div class="form-field">
                            <label class="form-label">Display Order</label>
                            <Input type="number" bind:value={form.display_order} />
                        </div>
                        <div class="form-field">
                            <label class="form-label">Important</label>
                            <label class="important-toggle">
                                <input type="checkbox" bind:checked={form.important} />
                                <span>Show in main grid</span>
                            </label>
                        </div>
                        <div class="form-field form-full">
                            <label class="form-label">Image URL</label>
                            <div class="image-upload-row">
                                <Input bind:value={form.image_url} placeholder="https://..." />
                                <label class="upload-btn" class:uploading disabled={uploading}>
                                    <input type="file" accept="image/png,image/jpeg,image/gif,image/webp" onchange={uploadImage} disabled={uploading} />
                                    {uploading ? "Uploading..." : "Upload"}
                                </label>
                            </div>
                            {#if form.image_url}
                                <div class="image-preview">
                                    <img src={form.image_url} alt="preview" onerror={(e) => (e.target as HTMLElement).style.display = "none"} />
                                </div>
                            {/if}
                        </div>
                        <div class="form-field form-full">
                            <label class="form-label">Credential URL</label>
                            <Input bind:value={form.credential_url} placeholder="https://gssoc.girlscript.org/profile/..." />
                        </div>
                    </div>
                    <div class="form-actions">
                        <Button onclick={saveBadge} disabled={!form.name.trim()}>
                            {editingId ? "Update" : "Create"}
                        </Button>
                        <Button variant="ghost" onclick={resetForm}>Cancel</Button>
                    </div>
                </div>
            {/if}

            {#if error}
                <div class="error-bar">{error}</div>
            {/if}
            {#if successMsg}
                <div class="success-bar">{successMsg}</div>
            {/if}

            <div class="table-wrapper">
                {#if loading}
                    <div class="loading-row">
                        <Skeleton class="w-50 h-6" />
                        <Skeleton class="w-80 h-6" />
                        <Skeleton class="w-30 h-6" />
                    </div>
                {:else if badges.length === 0}
                    <div class="empty-state">
                        <span class="empty-icon">[ ]</span>
                        <p>No badges yet. Add your first one above.</p>
                    </div>
                {:else}
                    <table class="data-table">
                        <thead>
                            <tr>
                                <th>Badge</th>
                                <th>Name</th>
                                <th>Category</th>
                                <th>Rarity</th>
                                <th>Imp</th>
                                <th>Order</th>
                                <th>Created</th>
                                <th>Actions</th>
                            </tr>
                        </thead>
                        <tbody>
                            {#each badges as b}
                                <tr>
                                    <td>
                                        {#if b.image_url}
                                            <img src={b.image_url} alt={b.name} class="badge-thumb" />
                                        {/if}
                                    </td>
                                    <td class="cell-name">{b.name}</td>
                                    <td class="cell-date">{b.category || "—"}</td>
                                    <td>
                                        <span class="rarity-badge" style="color: {rarityColors[b.rarity] || 'var(--accent)'}">
                                            {b.rarity}
                                        </span>
                                    </td>
                                    <td class="cell-imp">{b.important ? "★" : "—"}</td>
                                    <td>{b.display_order}</td>
                                    <td class="cell-date">{formatDate(b.created_at)}</td>
                                    <td class="cell-actions">
                                        <button class="action-btn" onclick={() => editBadge(b)}>Edit</button>
                                        <button class="action-btn action-delete" onclick={() => deleteBadge(b.id)}>Delete</button>
                                    </td>
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
    .terminal-container {
        min-height: 100vh;
        background: #030405;
        color: #e4e4e7;
        font-family: 'Space Grotesk', monospace;
        padding: 2rem;
    }

    .key-entry {
        display: flex;
        align-items: center;
        justify-content: center;
        min-height: 80vh;
    }
    .key-entry-inner {
        width: 100%;
        max-width: 420px;
        text-align: center;
    }
    .key-accent-bar {
        width: 60px; height: 3px;
        background: var(--accent);
        margin: 0 auto 1.5rem;
    }
    .key-title {
        font-family: var(--font-heading);
        font-size: 1.5rem;
        letter-spacing: 0.15em;
        margin-bottom: 0.5rem;
    }
    .key-subtitle {
        font-size: 0.85rem;
        color: #71717a;
        margin-bottom: 2rem;
    }
    .key-field { margin-bottom: 1rem; }
    .key-submit { display: flex; justify-content: center; }

    .panel {
        max-width: 1100px;
        margin: 0 auto;
    }
    .panel-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 2rem;
        padding-bottom: 1rem;
        border-bottom: 1px solid rgba(255,255,255,0.06);
    }
    .panel-title {
        display: flex;
        align-items: center;
        gap: 0.75rem;
    }
    .panel-accent {
        width: 4px; height: 24px;
        background: var(--accent);
    }
    .panel-title h1 {
        font-family: var(--font-heading);
        font-size: 1.2rem;
        letter-spacing: 0.1em;
    }
    .panel-actions {
        display: flex;
        gap: 0.5rem;
    }

    .form-panel {
        background: rgba(255,255,255,0.02);
        border: 1px solid rgba(255,255,255,0.06);
        padding: 1.5rem;
        margin-bottom: 1.5rem;
    }
    .form-title {
        font-family: var(--font-heading);
        font-size: 0.95rem;
        margin-bottom: 1.25rem;
        color: var(--accent);
    }
    .form-grid {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 1rem;
    }
    .form-full { grid-column: 1 / -1; }
    .form-label {
        display: block;
        font-size: 0.75rem;
        color: #71717a;
        margin-bottom: 0.4rem;
        letter-spacing: 0.05em;
    }
    .form-actions {
        display: flex;
        gap: 0.75rem;
        margin-top: 1.25rem;
    }

    .rarity-select {
        width: 100%;
        padding: 0.6rem 0.75rem;
        background: rgba(255,255,255,0.04);
        border: 1px solid rgba(255,255,255,0.08);
        color: #e4e4e7;
        font-family: inherit;
        font-size: 0.9rem;
        border-radius: 0;
    }

    .important-toggle {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        cursor: pointer;
        font-size: 0.85rem;
        color: #a1a1aa;
        padding: 0.5rem 0;
    }
    .important-toggle input[type="checkbox"] {
        width: 16px;
        height: 16px;
        accent-color: var(--accent);
    }

    .cell-imp { text-align: center; font-size: 0.85rem; }

    .image-upload-row {
        display: flex;
        gap: 0.5rem;
    }
    .image-upload-row > :first-child { flex: 1; }
    .upload-btn {
        display: inline-flex;
        align-items: center;
        padding: 0.5rem 1rem;
        background: rgba(255,68,0,0.1);
        border: 1px solid rgba(255,68,0,0.25);
        color: var(--accent);
        font-size: 0.8rem;
        cursor: pointer;
        white-space: nowrap;
    }
    .upload-btn input { display: none; }
    .upload-btn.uploading { opacity: 0.5; pointer-events: none; }
    .image-preview {
        margin-top: 0.75rem;
        border: 1px solid rgba(255,255,255,0.06);
        padding: 0.5rem;
        max-width: 200px;
    }
    .image-preview img {
        width: 100%;
        height: auto;
        display: block;
    }

    .error-bar {
        background: rgba(239,68,68,0.1);
        border: 1px solid rgba(239,68,68,0.25);
        color: #fca5a5;
        padding: 0.75rem 1rem;
        font-size: 0.85rem;
        margin-bottom: 1rem;
    }

    .success-bar {
        background: rgba(34,197,94,0.1);
        border: 1px solid rgba(34,197,94,0.25);
        color: #86efac;
        padding: 0.75rem 1rem;
        font-size: 0.85rem;
        margin-bottom: 1rem;
    }

    .table-wrapper {
        border: 1px solid rgba(255,255,255,0.06);
        overflow-x: auto;
    }
    .data-table {
        width: 100%;
        border-collapse: collapse;
        font-size: 0.85rem;
    }
    .data-table th {
        text-align: left;
        padding: 0.75rem 1rem;
        border-bottom: 1px solid rgba(255,255,255,0.06);
        color: #71717a;
        font-weight: 500;
        font-size: 0.75rem;
        letter-spacing: 0.05em;
    }
    .data-table td {
        padding: 0.75rem 1rem;
        border-bottom: 1px solid rgba(255,255,255,0.03);
    }
    .cell-name { font-weight: 500; }
    .cell-date { color: #71717a; font-size: 0.8rem; }
    .cell-actions {
        display: flex;
        gap: 0.5rem;
    }
    .action-btn {
        background: none;
        border: 1px solid rgba(255,255,255,0.1);
        color: #a1a1aa;
        padding: 0.3rem 0.6rem;
        font-size: 0.75rem;
        cursor: pointer;
        font-family: inherit;
    }
    .action-btn:hover { border-color: var(--accent); color: var(--accent); }
    .action-delete:hover { border-color: #ef4444; color: #ef4444; }

    .badge-thumb {
        width: 36px;
        height: 36px;
        object-fit: contain;
        display: block;
    }

    .loading-row {
        display: flex;
        gap: 1rem;
        padding: 2rem;
    }
    .empty-state {
        text-align: center;
        padding: 3rem 1rem;
        color: #71717a;
    }
    .empty-icon {
        font-size: 2rem;
        display: block;
        margin-bottom: 0.5rem;
    }

    .rarity-badge {
        font-size: 0.75rem;
        font-weight: 600;
        letter-spacing: 0.05em;
        text-transform: uppercase;
    }

    @media (max-width: 768px) {
        .terminal-container { padding: 1rem; }
        .form-grid { grid-template-columns: 1fr; }
        .panel-header { flex-direction: column; gap: 0.75rem; align-items: flex-start; }
    }
</style>
