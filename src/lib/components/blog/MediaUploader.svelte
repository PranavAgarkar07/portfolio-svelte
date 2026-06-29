<script lang="ts">
  import type { BlogImage } from '$lib/types';

  let {
    images = $bindable([] as BlogImage[]),
    disabled = false,
  }: {
    images: BlogImage[];
    disabled?: boolean;
  } = $props();

  const BASE = (import.meta.env.VITE_API_URL || '').replace(/\/$/, '');
  let uploading = $state(false);
  let dragIdx = $state<number | null>(null);

  function addImage(url: string) {
    if (!url.trim()) return;
    images = [...images, { url: url.trim() }];
  }

  function removeImage(idx: number) {
    const removed = images[idx];
    images = images.filter((_, i) => i !== idx);
    if (removed?.url) {
      const token = localStorage.getItem('portfolio_jwt');
      fetch(`${BASE}/api/admin/blog/delete-image`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          ...(token ? { 'Authorization': `Bearer ${token}` } : {}),
        },
        body: JSON.stringify({ url: removed.url }),
      }).catch(() => {});
    }
  }

  function moveImage(from: number, to: number) {
    if (to < 0 || to >= images.length) return;
    const arr = [...images];
    const [item] = arr.splice(from, 1);
    arr.splice(to, 0, item);
    images = arr;
  }

  function updateImage(idx: number, field: 'alt' | 'caption', value: string) {
    images = images.map((img, i) => (i === idx ? { ...img, [field]: value } : img));
  }

  function handleFileUpload(e: Event) {
    const input = e.target as HTMLInputElement;
    const file = input.files?.[0];
    if (!file) return;
    uploading = true;
    const token = localStorage.getItem('portfolio_jwt');
    const fd = new FormData();
    fd.append('image', file);
    fetch(`${BASE}/api/admin/blog/upload-image`, {
      method: 'POST',
      headers: token ? { 'Authorization': `Bearer ${token}` } : {},
      body: fd,
    })
      .then((r) => (r.ok ? r.json() : Promise.reject()))
      .then((data) => {
        images = [...images, { url: data.url }];
      })
      .catch(() => {})
      .finally(() => {
        uploading = false;
        input.value = '';
      });
  }
</script>

<div class="media-uploader">
  <div class="mu-header">
    <span class="mu-label">Gallery Images ({images.length})</span>
    <div class="mu-actions">
      <label class="mu-upload-btn" class:loading={uploading} tabindex="0">
        {uploading ? 'Uploading...' : '+ Upload'}
        <input type="file" accept="image/png,image/jpeg,image/gif,image/webp" onchange={handleFileUpload} hidden disabled={disabled || uploading} />
      </label>
      <button class="mu-add-btn" onclick={() => {
        const url = prompt('Paste image URL:');
        if (url) addImage(url);
      }} disabled={disabled}>+ URL</button>
    </div>
  </div>

  {#if images.length > 0}
    <div class="mu-grid" role="list">
      {#each images as img, i (img.url + i)}
        <div
          class="mu-item"
          class:mu-dragging={dragIdx === i}
          draggable={!disabled}
          ondragstart={() => dragIdx = i}
          ondragover={(e) => { e.preventDefault(); }}
          ondrop={() => { if (dragIdx !== null && dragIdx !== i) moveImage(dragIdx, i); dragIdx = null; }}
          ondragend={() => dragIdx = null}
          role="listitem"
        >
          <img src={img.url} alt={img.alt || ''} loading="lazy" />
          <div class="mu-overlay">
            <button class="mu-remove" onclick={() => removeImage(i)} disabled={disabled} aria-label="Remove image">&times;</button>
            <span class="mu-order">{i + 1}</span>
          </div>
          <div class="mu-fields">
            <input
              class="mu-input"
              placeholder="Alt text"
              value={img.alt || ''}
              oninput={(e) => updateImage(i, 'alt', (e.target as HTMLInputElement).value)}
              disabled={disabled}
            />
            <input
              class="mu-input"
              placeholder="Caption"
              value={img.caption || ''}
              oninput={(e) => updateImage(i, 'caption', (e.target as HTMLInputElement).value)}
              disabled={disabled}
            />
          </div>
        </div>
      {/each}
    </div>
  {:else}
    <div class="mu-empty">No gallery images yet.</div>
  {/if}
</div>

<style>
  .media-uploader {
    border: 1px solid rgba(255,255,255,0.12);
    padding: 0.75rem;
  }
  .mu-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 0.75rem;
  }
  .mu-label {
    font-size: 0.8rem;
    color: #a1a1aa;
  }
  .mu-actions {
    display: flex;
    gap: 0.4rem;
  }
  .mu-upload-btn, .mu-add-btn {
    padding: 0.3rem 0.6rem;
    background: rgba(255,255,255,0.08);
    border: 1px solid rgba(255,255,255,0.12);
    color: #a1a1aa;
    font-size: 0.75rem;
    cursor: pointer;
    font-family: inherit;
    transition: background 0.15s, color 0.15s;
  }
  .mu-upload-btn:hover, .mu-add-btn:hover {
    background: rgba(255,255,255,0.14);
    color: #e4e4e7;
  }
  .mu-upload-btn.loading {
    opacity: 0.5;
    cursor: not-allowed;
  }
  .mu-add-btn:disabled, .mu-upload-btn:disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }
  .mu-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
    gap: 0.75rem;
  }
  .mu-item {
    border: 1px solid rgba(255,255,255,0.1);
    overflow: hidden;
    transition: opacity 0.15s;
  }
  .mu-item.mu-dragging {
    opacity: 0.4;
  }
  .mu-item img {
    width: 100%;
    aspect-ratio: 16 / 9;
    object-fit: cover;
    display: block;
    background: rgba(0,0,0,0.3);
  }
  .mu-overlay {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0.25rem 0.4rem;
    background: rgba(0,0,0,0.4);
  }
  .mu-remove {
    background: none;
    border: none;
    color: #ef4444;
    font-size: 1rem;
    cursor: pointer;
    padding: 0;
    line-height: 1;
  }
  .mu-remove:disabled {
    opacity: 0.4;
  }
  .mu-order {
    font-size: 0.65rem;
    color: #a1a1aa;
    font-family: monospace;
  }
  .mu-fields {
    display: flex;
    flex-direction: column;
    gap: 0.2rem;
    padding: 0.35rem;
  }
  .mu-input {
    width: 100%;
    padding: 0.25rem 0.4rem;
    background: rgba(0,0,0,0.3);
    border: 1px solid rgba(255,255,255,0.08);
    color: #e4e4e7;
    font-size: 0.7rem;
    font-family: inherit;
    box-sizing: border-box;
  }
  .mu-input:focus {
    outline: none;
    border-color: var(--accent);
  }
  .mu-input:disabled {
    opacity: 0.5;
  }
  .mu-empty {
    color: #52525b;
    font-size: 0.8rem;
    text-align: center;
    padding: 1rem;
  }
</style>
