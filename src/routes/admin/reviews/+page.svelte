<script lang="ts">
  import { onMount } from 'svelte';

  const BASE = (import.meta.env.VITE_API_URL || '').replace(/\/$/, '');

  interface AdminReview {
    id: number;
    user_name: string;
    project_name: string;
    rating: number;
    comment: string;
    created_at: string;
  }

  let reviews = $state<AdminReview[]>([]);
  let loading = $state(true);

  onMount(async () => {
    if (!BASE) { loading = false; return; }
    const token = localStorage.getItem('portfolio_jwt');
    try {
      const r = await fetch(`${BASE}/api/admin/reviews`, {
        headers: token ? { 'Authorization': `Bearer ${token}` } : {},
      });
      if (r.ok) {
        const data = await r.json();
        reviews = data.reviews ?? [];
      }
    } catch {} finally {
      loading = false;
    }
  });

  async function deleteReview(id: number) {
    if (!confirm('Delete this review?')) return;
    const token = localStorage.getItem('portfolio_jwt');
    try {
      const r = await fetch(`${BASE}/api/admin/reviews/${id}`, {
        method: 'DELETE',
        headers: token ? { 'Authorization': `Bearer ${token}` } : {},
      });
      if (r.ok) {
        reviews = reviews.filter(r => r.id !== id);
      }
    } catch {}
  }
</script>

<div class="admin-page-header">
  <h1>Reviews</h1>
</div>

{#if loading}
  <div class="loading">Loading...</div>
{:else if reviews.length > 0}
  <table class="admin-table">
    <thead>
      <tr>
        <th>User</th>
        <th>Project</th>
        <th>Rating</th>
        <th>Comment</th>
        <th>Date</th>
        <th>Actions</th>
      </tr>
    </thead>
    <tbody>
      {#each reviews as review}
        <tr>
          <td>{review.user_name}</td>
          <td>{review.project_name}</td>
          <td>{'⭐'.repeat(review.rating)}</td>
          <td class="comment-cell">{review.comment || '—'}</td>
          <td>{review.created_at}</td>
          <td>
            <button class="delete-btn" onclick={() => deleteReview(review.id)}>Delete</button>
          </td>
        </tr>
      {/each}
    </tbody>
  </table>
{:else}
  <div class="empty">No reviews yet.</div>
{/if}

<style>
  .admin-page-header {
    margin-bottom: 1.5rem;
    padding: 0 2rem;
  }
  .admin-page-header h1 {
    font-family: var(--font-heading);
    font-size: 1.3rem;
    color: #e4e4e7;
    margin: 0;
  }
  .admin-table {
    width: calc(100% - 4rem);
    border-collapse: collapse;
    margin: 0 2rem;
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
  .comment-cell {
    max-width: 200px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  .delete-btn {
    padding: 0.2rem 0.5rem;
    border: 1px solid rgba(239,68,68,0.3);
    border-radius: 4px;
    font-size: 0.75rem;
    background: none;
    color: #ef4444;
    cursor: pointer;
    font-family: inherit;
  }
  .loading, .empty {
    color: #71717a;
    text-align: center;
    padding: 2rem;
  }
</style>
