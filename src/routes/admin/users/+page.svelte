<script lang="ts">
  import { onMount } from 'svelte';

  const BASE = (import.meta.env.VITE_API_URL || '').replace(/\/$/, '');

  interface AdminUser {
    id: string;
    email: string;
    name: string;
    avatar_url: string;
    role: string;
    created_at: string;
  }

  let users = $state<AdminUser[]>([]);
  let loading = $state(true);

  onMount(async () => {
    if (!BASE) { loading = false; return; }
    const token = localStorage.getItem('portfolio_jwt');
    try {
      const r = await fetch(`${BASE}/api/admin/users`, {
        headers: token ? { 'Authorization': `Bearer ${token}` } : {},
      });
      if (r.ok) {
        const data = await r.json();
        users = data.users ?? [];
      }
    } catch {} finally {
      loading = false;
    }
  });

  async function changeRole(userId: string, role: string) {
    const token = localStorage.getItem('portfolio_jwt');
    try {
      const r = await fetch(`${BASE}/api/admin/users/${userId}/role`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          ...(token ? { 'Authorization': `Bearer ${token}` } : {}),
        },
        body: JSON.stringify({ role }),
      });
      if (r.ok) {
        users = users.map(u => u.id === userId ? { ...u, role } : u);
      }
    } catch {}
  }
</script>

<div class="admin-page-header">
  <h1>Users</h1>
</div>

{#if loading}
  <div class="loading">Loading...</div>
{:else if users.length > 0}
  <table class="admin-table">
    <thead>
      <tr>
        <th>User</th>
        <th>Email</th>
        <th>Role</th>
        <th>Joined</th>
      </tr>
    </thead>
    <tbody>
      {#each users as user}
        <tr>
          <td class="user-cell">
            {#if user.avatar_url}
              <img src={user.avatar_url} alt="" class="user-avatar" />
            {/if}
            {user.name}
          </td>
          <td>{user.email}</td>
          <td>
            <select
              value={user.role}
              onchange={(e: Event) => changeRole(user.id, (e.target as HTMLSelectElement).value)}
              class="role-select"
            >
              <option value="user">User</option>
              <option value="author">Author</option>
              <option value="admin">Admin</option>
            </select>
          </td>
          <td>{user.created_at}</td>
        </tr>
      {/each}
    </tbody>
  </table>
{:else}
  <div class="empty">No users yet.</div>
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
  .user-cell {
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }
  .user-avatar {
    width: 24px;
    height: 24px;
    border-radius: 50%;
    object-fit: cover;
  }
  .role-select {
    background: rgba(0,0,0,0.3);
    border: 1px solid rgba(255,255,255,0.1);
    border-radius: 4px;
    color: #e4e4e7;
    padding: 0.2rem 0.4rem;
    font-size: 0.8rem;
    font-family: inherit;
  }
  .loading, .empty {
    color: #71717a;
    text-align: center;
    padding: 2rem;
  }
</style>
