<script lang="ts">
  import { user, isLoggedIn, isAdmin, isAuthor, logout } from '$lib/stores/auth';
  import { base } from '$app/paths';
  import GoogleLoginButton from './GoogleLoginButton.svelte';

  let showMenu = $state(false);
</script>

{#if $isLoggedIn}
  <div class="user-badge" onclick={() => showMenu = !showMenu} onkeydown={(e) => e.key === 'Enter' && (showMenu = !showMenu)} role="button" tabindex="0" aria-haspopup="true" aria-expanded={showMenu}>
    <div class="user-badge-trigger">
      <img src={$user?.avatar_url || ''} alt="" class="user-avatar" />
      <span class="user-name">{$user?.name}</span>
    </div>

    {#if showMenu}
      <div class="user-menu" role="menu" use:clickOutside={() => showMenu = false}>
        <div class="user-menu-header">
          <span class="user-menu-email">{$user?.email}</span>
          <span class="user-menu-role">{$user?.role}</span>
        </div>
        {#if $isAdmin || $isAuthor}
          <a href={`${base}/admin/dashboard`} class="user-menu-item" role="menuitem">Admin Panel</a>
        {/if}
        {#if $isAdmin || $isAuthor}
          <a href={`${base}/admin/blog`} class="user-menu-item" role="menuitem">My Posts</a>
        {/if}
        <button class="user-menu-item user-menu-logout" onclick={() => { logout(); showMenu = false; }} role="menuitem">
          Sign Out
        </button>
      </div>
    {/if}
  </div>
{:else}
  <GoogleLoginButton />
{/if}

<script context="module" lang="ts">
  function clickOutside(node: HTMLElement, cb: () => void) {
    function handler(e: MouseEvent) {
      if (!node.contains(e.target as Node)) cb();
    }
    document.addEventListener('click', handler, true);
    return {
      destroy() { document.removeEventListener('click', handler, true); }
    };
  }
</script>

<style>
  .user-badge {
    position: relative;
    cursor: pointer;
  }
  .user-badge-trigger {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.3rem 0.5rem;
    border-radius: 6px;
    transition: background 0.15s;
  }
  .user-badge-trigger:hover {
    background: rgba(255,255,255,0.05);
  }
  .user-avatar {
    width: 28px;
    height: 28px;
    border-radius: 50%;
    object-fit: cover;
  }
  .user-name {
    font-size: 0.8rem;
    color: #e4e4e7;
  }
  .user-menu {
    position: absolute;
    top: calc(100% + 4px);
    right: 0;
    background: #1a1a2e;
    border: 1px solid rgba(255,255,255,0.1);
    border-radius: 8px;
    padding: 0.5rem;
    min-width: 180px;
    z-index: 100;
    box-shadow: 0 8px 24px rgba(0,0,0,0.4);
  }
  .user-menu-header {
    padding: 0.5rem;
    border-bottom: 1px solid rgba(255,255,255,0.06);
    margin-bottom: 0.25rem;
  }
  .user-menu-email {
    display: block;
    font-size: 0.75rem;
    color: #71717a;
  }
  .user-menu-role {
    display: block;
    font-size: 0.7rem;
    color: var(--accent);
    text-transform: uppercase;
    letter-spacing: 0.05em;
    margin-top: 2px;
  }
  .user-menu-item {
    display: block;
    width: 100%;
    padding: 0.5rem;
    color: #e4e4e7;
    text-decoration: none;
    font-size: 0.8rem;
    border: none;
    background: none;
    text-align: left;
    cursor: pointer;
    border-radius: 4px;
    font-family: inherit;
  }
  .user-menu-item:hover {
    background: rgba(255,255,255,0.05);
  }
  .user-menu-logout {
    color: #ef4444;
    margin-top: 0.25rem;
    border-top: 1px solid rgba(255,255,255,0.06);
    padding-top: 0.5rem;
  }
</style>
