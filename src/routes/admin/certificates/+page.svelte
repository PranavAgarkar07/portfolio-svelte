<script lang="ts">
    import { onMount } from "svelte";
    import { Input, Button, Skeleton, Tag } from "$lib/components/ui";
    import type { Certificate } from "$lib/types";

    const BASE = (import.meta.env.VITE_API_URL || "").replace(/\/$/, "");

    let key = $state("");
    let storedKey = $state("");
    let certs = $state<Certificate[]>([]);
    let loading = $state(false);
    let error = $state("");

    let showForm = $state(false);
    let editingId = $state<number | null>(null);
    let uploading = $state(false);
    let form = $state({
        title: "",
        issuer: "",
        date: "",
        credential_url: "",
        image_url: "",
        tags: "",
        is_verified: false,
        display_order: 0,
    });

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
            loadCerts();
        }
    });

    async function loadCerts() {
        loading = true;
        error = "";
        try {
            const res = await fetch(`${BASE}/api/admin/certificates`, { headers: { Authorization: `Bearer ${storedKey}` } });
            if (!res.ok) {
                const text = await res.text();
                throw new Error(text || `HTTP ${res.status}`);
            }
            const data = await res.json();
            certs = data.certificates ?? data ?? [];
        } catch (e: any) {
            error = e.message || "Failed to load certificates";
            certs = [];
        } finally {
            loading = false;
        }
    }

    function toNum(v: unknown, fallback = 0): number {
        const n = typeof v === "string" ? parseInt(v, 10) : Number(v);
        return isNaN(n) ? fallback : n;
    }

    async function saveCert() {
        if (!form.title.trim() || !form.issuer.trim()) return;

        const body = {
            title: form.title,
            issuer: form.issuer,
            date: form.date,
            credential_url: form.credential_url,
            image_url: form.image_url,
            tags: form.tags ? form.tags.split(",").map((t: string) => t.trim()).filter(Boolean) : [],
            is_verified: form.is_verified,
            display_order: toNum(form.display_order, certs.length),
        };

        try {
            const url = editingId
                ? `${BASE}/api/admin/certificates/${editingId}`
                : `${BASE}/api/admin/certificates`;
            const method = editingId ? "PUT" : "POST";

            const res = await fetch(url, {
                method,
                headers: { "Content-Type": "application/json", Authorization: `Bearer ${storedKey}` },
                body: JSON.stringify(body),
            });

            if (!res.ok) {
                const text = await res.text();
                throw new Error(text || `HTTP ${res.status}`);
            }

            resetForm();
            await loadCerts();
        } catch (e: any) {
            error = e.message || "Failed to save certificate";
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
            const res = await fetch(`${BASE}/api/admin/certificates/upload`, {
                method: "POST",
                headers: { Authorization: `Bearer ${storedKey}` },
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

    async function deleteCert(id: number) {
        if (!confirm("Delete this certificate?")) return;

        try {
            const res = await fetch(`${BASE}/api/admin/certificates/${id}`, {
                method: "DELETE",
                headers: { Authorization: `Bearer ${storedKey}` },
            });

            if (!res.ok) {
                const text = await res.text();
                throw new Error(text || `HTTP ${res.status}`);
            }

            await loadCerts();
        } catch (e: any) {
            error = e.message || "Failed to delete certificate";
        }
    }

    function editCert(cert: Certificate) {
        editingId = cert.id;
        form = {
            title: cert.title,
            issuer: cert.issuer,
            date: cert.date || "",
            credential_url: cert.credential_url || "",
            image_url: cert.image_url || "",
            tags: (cert.tags || []).join(", "),
            is_verified: cert.is_verified,
            display_order: cert.display_order,
        };
        showForm = true;
    }

    function resetForm() {
        showForm = false;
        editingId = null;
        form = { title: "", issuer: "", date: "", credential_url: "", image_url: "", tags: "", is_verified: false, display_order: 0 };
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
                <h1 class="key-title">CERTIFICATES ADMIN</h1>
                <p class="key-subtitle">Authenticate to manage certificates</p>
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
                    <h1>CERTIFICATE CONTROL</h1>
                </div>
                <div class="panel-actions">
                    <Button variant="ghost" size="sm" onclick={() => { resetForm(); showForm = !showForm; }}>
                        {showForm ? "Cancel" : "Add Certificate"}
                    </Button>
                    <Button variant="ghost" size="sm" onclick={clearKey}>Change Key</Button>
                </div>
            </div>

            {#if showForm}
                <div class="form-panel">
                    <h3 class="form-title">{editingId ? "Edit Certificate" : "New Certificate"}</h3>
                    <div class="form-grid">
                        <div class="form-field">
                            <label class="form-label">Title *</label>
                            <Input bind:value={form.title} placeholder="e.g. AWS Certified Solutions Architect" />
                        </div>
                        <div class="form-field">
                            <label class="form-label">Issuer *</label>
                            <Input bind:value={form.issuer} placeholder="e.g. Amazon Web Services" />
                        </div>
                        <div class="form-field">
                            <label class="form-label">Date</label>
                            <Input type="date" bind:value={form.date} />
                        </div>
                        <div class="form-field">
                            <label class="form-label">Display Order</label>
                            <Input type="number" bind:value={form.display_order} />
                        </div>
                        <div class="form-field form-full">
                            <label class="form-label">Image URL</label>
                            <div class="image-upload-row">
                                <Input bind:value={form.image_url} placeholder="https://imgur.com/..." />
                                <label class="upload-btn" class:uploading aria-disabled={uploading}>
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
                            <Input bind:value={form.credential_url} placeholder="https://credential.example.com/..." />
                        </div>
                        <div class="form-field form-full">
                            <label class="form-label">Tags (comma-separated)</label>
                            <Input bind:value={form.tags} placeholder="Cloud, AWS, Architecture" />
                        </div>
                        <div class="form-field form-checkbox">
                            <label>
                                <input type="checkbox" bind:checked={form.is_verified} />
                                <span>Verified / Badge</span>
                            </label>
                        </div>
                    </div>
                    <div class="form-actions">
                        <Button onclick={saveCert} disabled={!form.title.trim() || !form.issuer.trim()}>
                            {editingId ? "Update" : "Create"}
                        </Button>
                        <Button variant="ghost" onclick={resetForm}>Cancel</Button>
                    </div>
                </div>
            {/if}

            <div class="stats-bar">
                <div class="stat">
                    <span class="stat-label">Total</span>
                    <span class="stat-value">{certs.length}</span>
                </div>
                <div class="stat">
                    <span class="stat-label">Verified</span>
                    <span class="stat-value accent">{certs.filter(c => c.is_verified).length}</span>
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
            {:else if certs.length === 0}
                <div class="empty-state">
                    <span class="empty-text">No certificates added yet.</span>
                    <span class="empty-cursor"></span>
                </div>
            {:else}
                <div class="table-wrapper">
                    <table>
                        <thead>
                            <tr>
                                <th class="col-num">#</th>
                                <th class="col-title">Title</th>
                                <th class="col-issuer">Issuer</th>
                                <th class="col-date">Date</th>
                                <th class="col-status">Status</th>
                                <th class="col-action"></th>
                            </tr>
                        </thead>
                        <tbody>
                            {#each certs as cert, i}
                                <tr class="cert-row" style="--i: {i}">
                                    <td class="row-num">{cert.display_order + 1}</td>
                                    <td>
                                        <div class="title-cell">{cert.title}</div>
                                    </td>
                                    <td><Tag variant="muted">{cert.issuer}</Tag></td>
                                    <td class="date-cell">{cert.date || "-"}</td>
                                    <td>
                                        {#if cert.is_verified}
                                            <span class="badge badge-new">Verified</span>
                                        {:else}
                                            <span class="badge badge-read">Standard</span>
                                        {/if}
                                    </td>
                                    <td>
                                        <div class="row-actions">
                                            <Button variant="ghost" size="sm" onclick={() => editCert(cert)}>Edit</Button>
                                            <Button variant="ghost" size="sm" onclick={() => deleteCert(cert.id)}>Del</Button>
                                        </div>
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
    .key-field { width: 320px; max-width: 100%; }
    .key-submit { margin-top: 0.25rem; }

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
    .panel-actions {
        display: flex;
        gap: 0.5rem;
        align-items: center;
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

    .form-panel {
        background: rgba(255,255,255,0.02);
        border: 1px solid rgba(255,255,255,0.06);
        padding: 1.5rem;
        margin-bottom: 1.5rem;
    }
    .form-title {
        font-family: var(--font-heading);
        font-size: 0.85rem;
        letter-spacing: 0.1em;
        color: var(--accent);
        margin-bottom: 1rem;
        text-transform: uppercase;
    }
    .form-grid {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 1rem;
    }
    .form-full { grid-column: 1 / -1; }
    .form-field {
        display: flex;
        flex-direction: column;
        gap: 0.35rem;
    }
    .form-label {
        font-size: 0.65rem;
        text-transform: uppercase;
        letter-spacing: 0.1em;
        color: var(--text-secondary);
    }
    .form-checkbox label {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        font-size: 0.8rem;
        color: var(--text-primary);
        cursor: pointer;
    }
    .form-checkbox input[type="checkbox"] {
        accent-color: var(--accent);
        width: 16px;
        height: 16px;
    }
    .form-actions {
        display: flex;
        gap: 0.75rem;
        margin-top: 1rem;
    }

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
    .stat-value.accent { color: var(--accent); }

    .table-wrapper { overflow-x: auto; }
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
    .col-num { width: 3rem; }
    .col-title { }
    .col-issuer { width: 18%; }
    .col-date { width: 12%; }
    .col-status { width: 10%; }
    .col-action { width: 12%; }

    .cert-row {
        transition: background 0.25s;
        animation: slideIn 0.35s ease-out both;
        animation-delay: calc(var(--i) * 0.05s);
    }
    .cert-row:hover {
        background: rgba(255,255,255,0.04);
        box-shadow: 0 0 0 1px rgba(255,255,255,0.04);
    }
    .row-num {
        color: rgba(255,255,255,0.2);
        font-size: 0.7rem;
        font-variant-numeric: tabular-nums;
    }
    .title-cell {
        font-weight: 500;
        font-size: 0.8rem;
        color: var(--text-primary);
    }
    .date-cell {
        font-size: 0.7rem;
        font-family: var(--font-body);
        color: var(--text-secondary);
        white-space: nowrap;
        font-variant-numeric: tabular-nums;
    }
    .row-actions {
        display: flex;
        gap: 0.35rem;
    }
    @media (max-width: 768px) {
        .row-actions { flex-direction: column; }
    }

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

    .empty-state {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 0.25rem;
        padding: 4rem 0;
        color: var(--text-secondary);
        font-size: 0.9rem;
    }
    .empty-text { letter-spacing: 0.05em; }
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

    .error-state {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        color: #ff4444;
        font-size: 0.8rem;
        padding: 0.75rem 0;
        letter-spacing: 0.03em;
    }
    .error-icon { font-size: 1rem; flex-shrink: 0; }
    .loading-grid {
        display: flex;
        flex-direction: column;
        gap: 2px;
    }

    @keyframes slideIn {
        from { opacity: 0; transform: translateY(10px); }
        to { opacity: 1; transform: translateY(0); }
    }
    @keyframes pulseGlow {
        0%, 100% { box-shadow: 0 0 0 rgba(255,68,0,0); }
        50% { box-shadow: 0 0 30px rgba(255,68,0,0.08), inset 0 0 30px rgba(255,68,0,0.02); }
    }

    @media (prefers-reduced-motion: reduce) {
        .key-entry-inner { animation: none; }
        .cert-row { animation: none; transition: none; }
        .empty-cursor { animation: none; opacity: 1; }
    }

    .image-upload-row {
        display: flex;
        gap: 0.5rem;
        align-items: flex-start;
    }
    .image-upload-row > :first-child {
        flex: 1;
        min-width: 0;
    }
    .upload-btn {
        flex-shrink: 0;
        display: inline-flex;
        align-items: center;
        justify-content: center;
        padding: 0.85rem 1.25rem;
        background: var(--accent);
        color: #fff;
        font-family: var(--font-heading);
        font-size: 0.65rem;
        font-weight: 700;
        letter-spacing: 0.1em;
        text-transform: uppercase;
        cursor: pointer;
        border: none;
        transition: opacity 0.15s ease;
        white-space: nowrap;
    }
    .upload-btn:hover { opacity: 0.85; }
    .upload-btn.uploading { opacity: 0.5; pointer-events: none; }
    .upload-btn input[type="file"] {
        display: none;
    }
    .image-preview {
        margin-top: 0.5rem;
        border: 1px solid rgba(255,255,255,0.06);
        max-height: 200px;
        overflow: hidden;
    }
    .image-preview img {
        display: block;
        max-width: 100%;
        max-height: 200px;
        object-fit: contain;
        background: #070a0f;
    }
    body.light-mode .upload-btn { color: #000; }
    body.light-mode .image-preview { border-color: rgba(0,0,0,0.1); }
    body.light-mode .image-preview img { background: #d4d4d8; }

    @media (max-width: 768px) {
        .terminal-container { padding: 1rem; }
        .panel { padding: 1rem; }
        .key-entry-inner { padding: 2rem 1.5rem; }
        .form-grid { grid-template-columns: 1fr; }
        .stats-bar { gap: 1rem; }
        table { font-size: 0.7rem; }
        th, td { padding: 0.5rem; }
        .image-upload-row { flex-direction: column; }
    }
</style>
