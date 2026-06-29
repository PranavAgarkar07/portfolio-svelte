<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { base } from '$app/paths';
  import MediaUploader from '$lib/components/blog/MediaUploader.svelte';
  import { cacheClear } from '$lib/utils/cache';
  import { renderMarkdown } from '$lib/utils/markdown';
  import type { BlogImage } from '$lib/types';

  const BASE = (import.meta.env.VITE_API_URL || '').replace(/\/$/, '');

  let title = $state('');
  let slug = $state('');
  let content = $state('');
  let excerpt = $state('');
  let coverImage = $state('');
  let galleryImages: BlogImage[] = $state([]);
  let tags = $state('');
  let published = $state(false);
  let loading = $state(true);
  let submitting = $state(false);
  let error = $state('');
  let uploading = $state(false);
  let isDragOver = $state(false);
  let showPreview = $state(false);

  let wordCount = $derived(content ? content.trim().split(/\s+/).length : 0);
  let charCount = $derived(content.length);
  let readingTime = $derived(Math.max(1, Math.round(wordCount / 200)));
  let textareaEl: HTMLTextAreaElement | undefined = $state();

  function captureWheel(node: HTMLElement) {
    const handler = (e: WheelEvent) => {
      const atTop = node.scrollTop === 0;
      const atBottom = node.scrollTop + node.clientHeight >= node.scrollHeight;
      if ((e.deltaY < 0 && atTop) || (e.deltaY > 0 && atBottom)) return;
      e.preventDefault();
      e.stopPropagation();
      node.scrollTop += e.deltaY;
    };
    node.addEventListener('wheel', handler, { passive: false, capture: true });
    return { destroy() { node.removeEventListener('wheel', handler, { capture: true }); } };
  }

  function insertAtCursor(text: string) {
    const ta = textareaEl;
    if (!ta) return;
    const start = ta.selectionStart;
    const end = ta.selectionEnd;
    content = content.substring(0, start) + text + content.substring(end);
    requestAnimationFrame(() => {
      ta.focus();
      ta.setSelectionRange(start + text.length, start + text.length);
    });
  }

  function wrapSelection(before: string, after: string) {
    const ta = textareaEl;
    if (!ta) return;
    const start = ta.selectionStart;
    const end = ta.selectionEnd;
    const selected = content.substring(start, end);
    content = content.substring(0, start) + before + selected + after + content.substring(end);
    requestAnimationFrame(() => {
      ta.focus();
      ta.setSelectionRange(start + before.length, start + before.length + selected.length);
    });
  }

  function insertLinePrefix(prefix: string) {
    const ta = textareaEl;
    if (!ta) return;
    const start = ta.selectionStart;
    const lineStart = content.lastIndexOf('\n', start - 1) + 1;
    content = content.substring(0, lineStart) + prefix + content.substring(lineStart);
    requestAnimationFrame(() => {
      ta.focus();
      ta.setSelectionRange(start + prefix.length, start + prefix.length);
    });
  }

  async function uploadAndInsert(file: File, insertPos?: number) {
    if (!file.type.startsWith('image/')) return;
    uploading = true;
    error = '';
    const token = localStorage.getItem('portfolio_jwt');
    const fd = new FormData();
    fd.append('image', file);
    try {
      const r = await fetch(`${BASE}/api/admin/blog/upload-image`, {
        method: 'POST',
        headers: token ? { 'Authorization': `Bearer ${token}` } : {},
        body: fd,
      });
      if (r.ok) {
        const data = await r.json();
        const md = `![](${data.url})\n`;
        if (insertPos !== undefined) {
          content = content.substring(0, insertPos) + md + content.substring(insertPos);
          requestAnimationFrame(() => {
            textareaEl?.focus();
            const newPos = insertPos + data.url.length + 5;
            textareaEl?.setSelectionRange(newPos, newPos);
          });
        } else {
          const ta = textareaEl;
          if (ta) {
            const start = ta.selectionStart;
            content = content.substring(0, start) + md + content.substring(ta.selectionEnd);
            requestAnimationFrame(() => ta.focus());
          }
        }
      } else {
        error = 'Upload failed';
      }
    } catch {
      error = 'Upload failed';
    } finally {
      uploading = false;
    }
  }

  async function uploadImage(e: Event) {
    const input = e.target as HTMLInputElement;
    const file = input.files?.[0];
    if (!file) return;
    await uploadAndInsert(file);
    input.value = '';
  }

  function handleDragOver(e: DragEvent) {
    e.preventDefault();
    isDragOver = true;
  }

  function handleDragLeave() {
    isDragOver = false;
  }

  function handleDrop(e: DragEvent) {
    e.preventDefault();
    isDragOver = false;
    const files = e.dataTransfer?.files;
    if (!files?.length) return;
    for (const file of Array.from(files)) {
      if (file.type.startsWith('image/')) {
        const ta = textareaEl;
        let insertPos: number | undefined;
        if (ta) {
          const rect = ta.getBoundingClientRect();
          const lineHeight = parseInt(getComputedStyle(ta).lineHeight || '20');
          const relY = e.clientY - rect.top;
          const line = Math.round(relY / lineHeight);
          const lines = content.substring(0, ta.selectionStart).split('\n');
          insertPos = lines.slice(0, line).join('\n').length + (line > 0 ? 1 : 0);
        }
        uploadAndInsert(file, insertPos);
        break;
      }
    }
  }

  function handlePaste(e: ClipboardEvent) {
    const items = e.clipboardData?.items;
    if (!items) return;
    for (const item of Array.from(items)) {
      if (item.type.startsWith('image/')) {
        const file = item.getAsFile();
        if (file) {
          e.preventDefault();
          uploadAndInsert(file);
        }
        break;
      }
    }
  }

  function handleTab(e: KeyboardEvent) {
    if (e.key === 'Tab') {
      e.preventDefault();
      const ta = textareaEl;
      if (!ta) return;
      const start = ta.selectionStart;
      content = content.substring(0, start) + '  ' + content.substring(ta.selectionEnd);
      requestAnimationFrame(() => {
        ta.focus();
        ta.setSelectionRange(start + 2, start + 2);
      });
    }
  }

  onMount(async () => {
    const id = $page.params.id;
    if (!BASE || !id) { loading = false; return; }
    const token = localStorage.getItem('portfolio_jwt');
    try {
      const r = await fetch(`${BASE}/api/admin/blog/${id}`, {
        headers: token ? { 'Authorization': `Bearer ${token}` } : {},
      });
      if (r.ok) {
        const data = await r.json();
        const p = data.post;
        if (p) {
          title = p.title;
          slug = p.slug;
          content = p.content_md;
          excerpt = p.excerpt || '';
          coverImage = p.cover_image || '';
          galleryImages = p.images || [];
          tags = (p.tags || []).join(', ');
          published = p.published;
        }
      }
    } catch {} finally {
      loading = false;
    }
  });

  async function submit() {
    if (!title.trim() || !content.trim()) {
      error = 'Title and content are required';
      return;
    }
    submitting = true;
    error = '';
    const id = $page.params.id;
    const token = localStorage.getItem('portfolio_jwt');
    try {
      const r = await fetch(`${BASE}/api/admin/blog/${id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          ...(token ? { 'Authorization': `Bearer ${token}` } : {}),
        },
        body: JSON.stringify({
          slug,
          title,
          content_md: content,
          excerpt,
          cover_image: coverImage,
          images: galleryImages,
          tags: tags ? tags.split(',').map(t => t.trim()).filter(Boolean) : [],
          published,
        }),
      });
      if (r.ok) {
        cacheClear('blog:');
        goto(`${base}/admin/blog`);
      } else {
        const data = await r.json();
        error = data.error?.message || 'Failed to update post';
      }
    } catch {
      error = 'Network error';
    } finally {
      submitting = false;
    }
  }
</script>

<div class="editor-page">
  <header class="editor-header">
    <div>
      <h1>Edit Post</h1>
      <p class="editor-subtitle">Update your existing blog entry</p>
    </div>
    <a href={`${base}/admin/blog`} class="btn-secondary">
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M19 12H5"/><polyline points="12 19 5 12 12 5"/></svg>
      Back
    </a>
  </header>

  {#if loading}
    <div class="loading-state">
      <svg class="spin" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#71717a" stroke-width="2"><circle cx="12" cy="12" r="10" opacity="0.3"/><path d="M12 2a10 10 0 0 1 10 10" stroke-linecap="round"/></svg>
      Loading post...
    </div>
  {:else}
    {#if error}
      <div class="error-msg">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/></svg>
        {error}
      </div>
    {/if}

    <div class="editor-title-field">
      <input
        id="title"
        bind:value={title}
        placeholder="Post title"
        class="title-input"
      />
    </div>

    <div class="editor-meta-grid">
      <div class="editor-section">
        <div class="editor-field">
          <label for="slug">Slug</label>
          <input id="slug" bind:value={slug} placeholder="post-url-slug" class="meta-input" />
        </div>

        <div class="editor-field">
          <label for="excerpt">Excerpt</label>
          <input id="excerpt" bind:value={excerpt} placeholder="Short description for cards & SEO" class="meta-input" />
        </div>

        <div class="editor-field">
          <label for="tags">Tags</label>
          <input id="tags" bind:value={tags} placeholder="svelte, tutorial, webdev" class="meta-input" />
        </div>
      </div>

      <div class="editor-section">
        <div class="editor-field">
          <label for="coverImage">Cover Image</label>
          <div class="cover-row">
            <input id="coverImage" bind:value={coverImage} placeholder="https://..." class="meta-input" />
            <label class="btn-icon" class:loading={uploading}>
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/></svg>
              <input type="file" accept="image/png,image/jpeg,image/gif,image/webp" onchange={async (e) => {
                const input = e.target as HTMLInputElement;
                const file = input.files?.[0];
                if (!file) return;
                uploading = true;
                const token = localStorage.getItem('portfolio_jwt');
                const fd = new FormData();
                fd.append('image', file);
                try {
                  const r = await fetch(`${BASE}/api/admin/blog/upload-image`, {
                    method: 'POST',
                    headers: token ? { 'Authorization': `Bearer ${token}` } : {},
                    body: fd,
                  });
                  if (r.ok) { const d = await r.json(); coverImage = d.url; }
                } catch {} finally { uploading = false; input.value = ''; }
              }} hidden />
            </label>
          </div>
        </div>

        <div class="editor-field">
          <label>Publish</label>
          <label class="toggle-row">
            <input type="checkbox" bind:checked={published} />
            <span class="toggle-track">
              <span class="toggle-thumb"></span>
            </span>
            <span class="toggle-label">{published ? 'Published' : 'Draft'}</span>
          </label>
        </div>

        <div class="editor-actions">
          <button class="btn-primary" onclick={submit} disabled={submitting}>
            {#if submitting}
              <svg class="spin" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10" opacity="0.3"/><path d="M12 2a10 10 0 0 1 10 10" stroke-linecap="round"/></svg>
              Saving...
            {:else}
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"/><polyline points="17 21 17 13 7 13 7 21"/><polyline points="7 3 7 8 15 8"/></svg>
              Update Post
            {/if}
          </button>
        </div>
      </div>
    </div>

    <div class="editor-section editor-section-full">
      <div class="editor-field">
        <label>Gallery Images</label>
        <MediaUploader bind:images={galleryImages} />
      </div>
    </div>

    <div class="editor-section editor-section-full editor-section-content">
      <div class="editor-field">
        <div class="content-header">
          <label for="content">Content</label>
          <div class="content-stats">
            <span>{wordCount} words</span>
            <span class="stat-sep">·</span>
            <span>{readingTime} min read</span>
            <span class="stat-sep">·</span>
            <span>{charCount} chars</span>
          </div>
        </div>

        <div class="editor-toolbar">
          <button class="tb-btn" onclick={() => wrapSelection('**', '**')} title="Bold" aria-label="Bold (Ctrl+B)">
            <svg width="15" height="15" viewBox="0 0 24 24" fill="currentColor"><path d="M15.6 10.79c.97-.67 1.65-1.77 1.65-2.79 0-2.26-1.75-4-4-4H7v14h7.04c2.09 0 3.71-1.7 3.71-3.79 0-1.52-.86-2.82-2.15-3.42zM10 6.5h3c.83 0 1.5.67 1.5 1.5s-.67 1.5-1.5 1.5h-3v-3zm3.5 9H10v-3h3.5c.83 0 1.5.67 1.5 1.5s-.67 1.5-1.5 1.5z"/></svg>
          </button>
          <button class="tb-btn" onclick={() => wrapSelection('*', '*')} title="Italic" aria-label="Italic (Ctrl+I)">
            <svg width="15" height="15" viewBox="0 0 24 24" fill="currentColor"><path d="M10 4v3h2.21l-3.42 8H6v3h8v-3h-2.21l3.42-8H18V4z"/></svg>
          </button>
          <button class="tb-btn" onclick={() => wrapSelection('~~', '~~')} title="Strikethrough" aria-label="Strikethrough">
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M16 4H9a3 3 0 0 0-2.83 4"/><path d="M14 12a4 4 0 0 1 0 8H6"/><line x1="4" y1="12" x2="20" y2="12"/></svg>
          </button>
          <span class="tb-sep"></span>
          <button class="tb-btn" onclick={() => insertLinePrefix('## ')} title="Heading" aria-label="Heading">
            <svg width="15" height="15" viewBox="0 0 24 24" fill="currentColor"><path d="M3 4h2v7h6V4h2v16h-2v-7H5v7H3V4zm14 0h2v16h-2V4z"/></svg>
          </button>
          <button class="tb-btn" onclick={() => insertLinePrefix('- ')} title="List" aria-label="List">
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="8" y1="6" x2="21" y2="6"/><line x1="8" y1="12" x2="21" y2="12"/><line x1="8" y1="18" x2="21" y2="18"/><line x1="3" y1="6" x2="3.01" y2="6"/><line x1="3" y1="12" x2="3.01" y2="12"/><line x1="3" y1="18" x2="3.01" y2="18"/></svg>
          </button>
          <button class="tb-btn" onclick={() => insertLinePrefix('> ')} title="Blockquote" aria-label="Blockquote">
            <svg width="15" height="15" viewBox="0 0 24 24" fill="currentColor"><path d="M6 17h3l2-4V7H5v6h3zm8 0h3l2-4V7h-6v6h3z"/></svg>
          </button>
          <button class="tb-btn" onclick={() => { insertLinePrefix('```\n'); insertAtCursor('\n```'); }} title="Code block" aria-label="Code block">
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="16 18 22 12 16 6"/><polyline points="8 6 2 12 8 18"/></svg>
          </button>
          <button class="tb-btn" onclick={() => wrapSelection('[', '](url)')} title="Link" aria-label="Link">
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"/><path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"/></svg>
          </button>
          <button class="tb-btn" onclick={() => { const ta = textareaEl; if (ta) { const sel = ta.value.substring(ta.selectionStart, ta.selectionEnd); if (sel) { insertAtCursor(`![${sel}]()`); } else { insertAtCursor('![]()'); } } }} title="Image" aria-label="Insert image">
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="18" height="18" rx="2" ry="2"/><circle cx="8.5" cy="8.5" r="1.5"/><polyline points="21 15 16 10 5 21"/></svg>
          </button>
          <span class="tb-sep"></span>
          <label class="tb-btn" class:loading={uploading}>
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/></svg>
            <input type="file" accept="image/png,image/jpeg,image/gif,image/webp" onchange={uploadImage} hidden />
          </label>
          <span class="tb-sep"></span>
          <button class="tb-btn tb-preview" class:active={showPreview} onclick={() => showPreview = !showPreview} title="Preview" aria-label="Toggle preview">
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/><circle cx="12" cy="12" r="3"/></svg>
            Preview
          </button>
        </div>

        <div class="textarea-wrap">
          {#if showPreview}
            <div class="preview-pane">
              <div class="preview-content">{@html renderMarkdown(content)}</div>
            </div>
          {:else}
              <textarea
                id="content"
                bind:this={textareaEl}
                bind:value={content}
                use:captureWheel
                placeholder="Write your post in markdown...&#10;&#10;Drag & drop images here or paste from clipboard"
              ondragover={handleDragOver}
              ondragleave={handleDragLeave}
              ondrop={handleDrop}
              onpaste={handlePaste}
              onkeydown={handleTab}
            ></textarea>
          {/if}
          {#if isDragOver}
            <div class="drop-overlay">Drop image to upload & insert</div>
          {/if}
        </div>
      </div>
    </div>
  {/if}
</div>

<style>
  .editor-page {
    max-width: 820px;
    margin: 0 auto;
    padding: 2rem 2rem 4rem;
  }
  .editor-header {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    margin-bottom: 1.5rem;
  }
  .editor-header h1 {
    font-family: var(--font-heading);
    font-size: 1.5rem;
    color: var(--text-primary);
    margin: 0;
    line-height: 1.2;
  }
  .editor-subtitle {
    color: var(--text-muted);
    font-size: 0.8rem;
    margin: 0.25rem 0 0;
  }
  .btn-secondary {
    display: inline-flex;
    align-items: center;
    gap: 0.4rem;
    padding: 0.45rem 0.85rem;
    border: 1px solid var(--border-color);
    background: var(--surface-dark);
    color: var(--text-secondary);
    font-size: 0.8rem;
    cursor: pointer;
    font-family: inherit;
    text-decoration: none;
    transition: all 0.15s;
  }
  .btn-secondary:hover {
    background: var(--surface-raised);
    color: var(--text-primary);
    border-color: var(--border-color);
  }
  .btn-secondary svg, .btn-primary svg, .btn-icon svg {
    flex-shrink: 0;
  }

  .loading-state {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    color: var(--text-muted);
    font-size: 0.9rem;
    padding: 3rem;
  }

  .error-msg {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    color: #ef4444;
    font-size: 0.85rem;
    padding: 0.6rem 0.85rem;
    background: rgba(239,68,68,0.1);
    border: 1px solid rgba(239,68,68,0.2);
    margin-bottom: 1rem;
  }

  .editor-title-field {
    margin-bottom: 1.25rem;
  }
  .title-input {
    width: 100%;
    padding: 0.7rem 0.85rem;
    background: var(--surface-dark);
    border: 1px solid var(--border-color);
    color: var(--text-primary);
    font-size: 1.15rem;
    font-family: var(--font-heading);
    box-sizing: border-box;
    transition: border-color 0.15s, box-shadow 0.15s;
  }
  .title-input:focus {
    outline: none;
    border-color: var(--accent);
    box-shadow: 0 0 0 2px rgba(255,68,0,0.15);
  }
  .title-input::placeholder {
    color: var(--text-muted);
  }

  .editor-meta-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
    margin-bottom: 1rem;
  }
  @media (max-width: 640px) {
    .editor-meta-grid {
      grid-template-columns: 1fr;
    }
  }

  .editor-section {
    padding: 1rem;
    border: 1px solid var(--border-color);
    background: var(--surface-dark);
  }
  .editor-section-full {
    margin-bottom: 1rem;
  }

  .editor-field {
    margin-bottom: 0.85rem;
  }
  .editor-field:last-child {
    margin-bottom: 0;
  }
  .editor-field label {
    display: block;
    font-size: 0.75rem;
    color: var(--text-muted);
    margin-bottom: 0.3rem;
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }
  .meta-input {
    width: 100%;
    padding: 0.5rem 0.65rem;
    background: var(--surface-raised);
    border: 1px solid var(--border-color);
    color: var(--text-primary);
    font-size: 0.85rem;
    font-family: inherit;
    box-sizing: border-box;
    transition: border-color 0.15s, box-shadow 0.15s;
  }
  .meta-input:focus {
    outline: none;
    border-color: var(--accent);
    box-shadow: 0 0 0 2px rgba(255,68,0,0.15);
  }
  .meta-input::placeholder {
    color: var(--text-muted);
  }

  .cover-row {
    display: flex;
    gap: 0.4rem;
  }
  .cover-row .meta-input {
    flex: 1;
  }
  .btn-icon {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 34px;
    min-width: 34px;
    padding: 0;
    background: var(--surface-dark);
    border: 1px solid var(--border-color);
    color: var(--text-secondary);
    cursor: pointer;
    transition: all 0.15s;
    position: relative;
  }
  .btn-icon:hover {
    background: var(--surface-raised);
    color: var(--text-primary);
  }
  .btn-icon.loading {
    opacity: 0.5;
  }

  .toggle-row {
    display: flex;
    align-items: center;
    gap: 0.6rem;
    cursor: pointer;
  }
  .toggle-row input {
    display: none;
  }
  .toggle-track {
    width: 36px;
    height: 20px;
    border-radius: 10px;
    background: var(--border-color);
    position: relative;
    transition: background 0.2s;
  }
  .toggle-row input:checked + .toggle-track {
    background: var(--accent);
  }
  .toggle-thumb {
    position: absolute;
    top: 2px;
    left: 2px;
    width: 16px;
    height: 16px;
    border-radius: 50%;
    background: var(--text-primary);
    transition: transform 0.2s;
  }
  .toggle-row input:checked + .toggle-track .toggle-thumb {
    transform: translateX(16px);
  }
  .toggle-label {
    font-size: 0.85rem;
    color: var(--text-secondary);
    text-transform: none;
    letter-spacing: 0;
  }

  .editor-actions {
    display: flex;
    gap: 0.6rem;
    margin-top: 0.85rem;
  }
  .btn-primary {
    display: inline-flex;
    align-items: center;
    gap: 0.4rem;
    padding: 0.55rem 1.25rem;
    background: var(--accent);
    color: #fff;
    border: none;
    cursor: pointer;
    font-size: 0.85rem;
    font-family: inherit;
    transition: opacity 0.15s;
  }
  .btn-primary:hover {
    opacity: 0.9;
  }
  .btn-primary:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .spin {
    animation: spin 0.8s linear infinite;
  }
  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  .content-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 0.3rem;
  }
  .content-header label {
    margin-bottom: 0;
  }
  .content-stats {
    display: flex;
    align-items: center;
    gap: 0.3rem;
    font-size: 0.7rem;
    color: var(--text-muted);
    font-family: monospace;
  }
  .stat-sep {
    color: var(--text-muted);
  }

  .editor-toolbar {
    display: flex;
    align-items: center;
    gap: 2px;
    padding: 0.35rem 0.5rem;
    background: var(--surface-raised);
    border: 1px solid var(--border-color);
    border-bottom: none;
  }
  .tb-btn {
    display: inline-flex;
    align-items: center;
    gap: 0.3rem;
    padding: 0.3rem 0.45rem;
    background: none;
    border: none;
    color: var(--text-muted);
    cursor: pointer;
    font-size: 0.75rem;
    font-family: inherit;
    transition: all 0.12s;
  }
  .tb-btn:hover {
    background: var(--surface-dark);
    color: var(--text-primary);
  }
  .tb-btn:focus-visible {
    outline: 1px solid var(--accent);
    outline-offset: -1px;
  }
  .tb-btn.active {
    color: var(--accent);
    background: var(--accent-glow);
  }
  .tb-btn.loading {
    opacity: 0.5;
  }
  .tb-sep {
    width: 1px;
    height: 18px;
    background: var(--border-color);
    margin: 0 0.2rem;
  }
  .tb-preview {
    margin-left: auto;
  }

  .textarea-wrap {
    position: relative;
  }
  .editor-field textarea {
    width: 100%;
    padding: 0.85rem;
    background: var(--surface-raised);
    border: 1px solid var(--border-color);
    color: var(--text-primary);
    font-family: 'Fira Code', 'JetBrains Mono', 'Cascadia Code', monospace;
    font-size: 0.85rem;
    line-height: 1.6;
    resize: vertical;
    height: 300px;
    min-height: 150px;
    max-height: 80vh;
    box-sizing: border-box;
    transition: border-color 0.15s, box-shadow 0.15s;
    overflow-y: scroll;
  }
  .editor-field textarea::-webkit-scrollbar {
    width: 8px;
  }
  .editor-field textarea::-webkit-scrollbar-track {
    background: transparent;
  }
  .editor-field textarea::-webkit-scrollbar-thumb {
    background: rgba(128,128,128,0.4);
    border-radius: 4px;
  }
  .editor-field textarea::-webkit-scrollbar-thumb:hover {
    background: rgba(128,128,128,0.6);
  }
  .editor-field textarea:focus {
    outline: none;
    border-color: var(--accent);
    box-shadow: 0 0 0 2px rgba(255,68,0,0.15);
  }
  .editor-field textarea::placeholder {
    color: var(--text-muted);
  }

  .preview-pane {
    width: 100%;
    min-height: 400px;
    max-height: 70vh;
    padding: 1.5rem;
    background: var(--surface-raised);
    border: 1px solid var(--border-color);
    box-sizing: border-box;
    color: var(--text-primary);
    font-size: 0.9rem;
    line-height: 1.7;
    overflow-y: auto;
  }
  .preview-pane::-webkit-scrollbar {
    width: 8px;
  }
  .preview-pane::-webkit-scrollbar-track {
    background: transparent;
  }
  .preview-pane::-webkit-scrollbar-thumb {
    background: rgba(128,128,128,0.4);
    border-radius: 4px;
  }
  .preview-pane::-webkit-scrollbar-thumb:hover {
    background: rgba(128,128,128,0.6);
  }
  .preview-content :global(h1), .preview-content :global(h2), .preview-content :global(h3) {
    color: var(--text-primary);
    font-family: var(--font-heading);
    margin: 1.5em 0 0.5em;
  }
  .preview-content :global(p) {
    margin: 0 0 1em;
  }
  .preview-content :global(code) {
    background: var(--code-bg);
    color: var(--code-color);
    padding: 0.15em 0.35em;
    font-size: 0.85em;
  }
  .preview-content :global(pre) {
    background: var(--surface-dark);
    padding: 1rem;
    overflow-x: auto;
  }
  .preview-content :global(pre code) {
    background: none;
    padding: 0;
    color: var(--text-primary);
  }
  .preview-content :global(img) {
    max-width: 100%;
    height: auto;
  }
  .preview-content :global(blockquote) {
    border-left: 3px solid var(--accent);
    padding-left: 1rem;
    color: var(--text-secondary);
    margin: 1em 0;
  }
  .preview-content :global(a) {
    color: var(--accent);
  }

  .drop-overlay {
    position: absolute;
    inset: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--surface-dark);
    color: var(--text-primary);
    font-size: 0.9rem;
    border: 2px dashed var(--accent);
    pointer-events: none;
    z-index: 1;
  }
</style>
