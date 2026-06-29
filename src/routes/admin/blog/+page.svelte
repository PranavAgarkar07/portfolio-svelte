<script lang="ts">
  import { onMount } from 'svelte';
  import { base } from '$app/paths';
  import type { BlogPost } from '$lib/types';

  const BASE = (import.meta.env.VITE_API_URL || '').replace(/\/$/, '');

  let posts = $state<BlogPost[]>([]);
  let loading = $state(true);

  onMount(async () => {
    if (!BASE) { loading = false; return; }
    const token = localStorage.getItem('portfolio_jwt');
    try {
      const r = await fetch(`${BASE}/api/admin/blog`, {
        headers: token ? { 'Authorization': `Bearer ${token}` } : {},
      });
      if (r.ok) {
        const data = await r.json();
        posts = data.posts ?? [];
      }
    } catch {} finally {
      loading = false;
    }
  });

  async function deletePost(id: string) {
    if (!confirm('Delete this post?')) return;
    const token = localStorage.getItem('portfolio_jwt');
    try {
      const r = await fetch(`${BASE}/api/admin/blog/${id}`, {
        method: 'DELETE',
        headers: token ? { 'Authorization': `Bearer ${token}` } : {},
      });
      if (r.ok) {
        posts = posts.filter(p => p.id !== id);
      }
    } catch {}
  }
</script>

<div class="admin-page-header">
  <h1>Blog Posts</h1>
  <a href={`${base}/admin/blog/new`} class="create-btn">+ New Post</a>
</div>

{#if loading}
  <div class="loading">Loading...</div>
{:else if posts.length > 0}
  <table class="admin-table">
    <thead>
      <tr>
        <th>Title</th>
        <th>Author</th>
        <th>Status</th>
        <th>Date</th>
        <th>Actions</th>
      </tr>
    </thead>
    <tbody>
      {#each posts as post}
        <tr>
          <td>{post.title}</td>
          <td>{post.author_name}</td>
          <td><span class="status-badge" class:published={post.published}>{post.published ? 'Published' : 'Draft'}</span></td>
          <td>{post.published_at || post.created_at}</td>
          <td class="actions-cell">
            <a href={`${base}/admin/blog/${post.id}`} class="action-btn">Edit</a>
            <button class="action-btn delete" onclick={() => deletePost(post.id)}>Delete</button>
          </td>
        </tr>
      {/each}
    </tbody>
  </table>
{:else}
  <div class="empty">No posts yet. Create your first post!</div>
{/if}

<style>
  .admin-page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
    padding: 0 2rem;
  }
  .admin-page-header h1 {
    font-family: var(--font-heading);
    font-size: 1.3rem;
    color: #e4e4e7;
    margin: 0;
  }
  .create-btn {
    padding: 0.4rem 0.8rem;
    background: var(--accent);
    color: #fff;
    text-decoration: none;
    border-radius: 6px;
    font-size: 0.8rem;
  }
  .admin-table {
    width: 100%;
    border-collapse: collapse;
    margin: 0 2rem;
    width: calc(100% - 4rem);
  }
  .admin-table th, .admin-table td {
    padding: 0.6rem 0.8rem;
    text-align: left;
    border-bottom: 1px solid rgba(255,255,255,0.06);
    font-size: 0.8rem;
    color: #e4e4e7;
  }
  .admin-table th {
    color: #71717a;
    font-weight: 500;
  }
  .status-badge {
    padding: 0.15rem 0.5rem;
    border-radius: 4px;
    font-size: 0.7rem;
    background: rgba(255,255,255,0.05);
    color: #71717a;
  }
  .status-badge.published {
    background: rgba(34, 197, 94, 0.1);
    color: #22c55e;
  }
  .actions-cell {
    display: flex;
    gap: 0.5rem;
  }
  .action-btn {
    padding: 0.2rem 0.5rem;
    border: 1px solid rgba(255,255,255,0.15);
    border-radius: 4px;
    font-size: 0.75rem;
    background: none;
    color: #e4e4e7;
    cursor: pointer;
    text-decoration: none;
    font-family: inherit;
  }
  .action-btn.delete {
    color: #ef4444;
    border-color: rgba(239,68,68,0.3);
  }
  .loading, .empty {
    color: #71717a;
    text-align: center;
    padding: 2rem;
  }
</style>
