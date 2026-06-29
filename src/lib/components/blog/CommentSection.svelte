<script lang="ts">
  import { onMount } from 'svelte';
  import { user, token, isLoggedIn } from '$lib/stores/auth';
  import { loginWithGoogle } from '$lib/stores/auth';
  import type { BlogComment } from '$lib/types';
  import { staleWhileRevalidate, cacheDel } from '$lib/utils/cache';

  let { slug }: { slug: string } = $props();

  const BASE = (import.meta.env.VITE_API_URL || '').replace(/\/$/, '');

  let comments = $state<BlogComment[]>([]);
  let loading = $state(true);
  let replyToId = $state<string | null>(null);
  let replyToName = $state('');
  let replyText = $state('');
  let showLoginPopup = $state(false);
  let toasts = $state<Array<{ id: number; message: string; type: 'error' | 'success' }>>([]);
  let toastId = $state(0);

  let currentUserId = $derived($user?.id);
  let loggedIn = $derived($isLoggedIn);
  let cacheKey = $derived(`blog:comments:${slug}`);

  let flatList = $derived.by<Array<{ comment: BlogComment; depth: number }>>(() => {
    function flatten(items: BlogComment[], depth: number): Array<{ comment: BlogComment; depth: number }> {
      const result: Array<{ comment: BlogComment; depth: number }> = [];
      for (const c of items) {
        result.push({ comment: c, depth });
        if (c.replies?.length) {
          result.push(...flatten(c.replies, depth + 1));
        }
      }
      return result;
    }
    return flatten(comments, 0);
  });

  function toast(message: string, type: 'error' | 'success' = 'error') {
    const id = ++toastId;
    toasts = [...toasts, { id, message, type }];
    setTimeout(() => { toasts = toasts.filter(t => t.id !== id); }, 4000);
  }

  onMount(fetchComments);

  async function fetchComments() {
    if (!BASE || !slug) { loading = false; return; }
    try {
      const ok = await staleWhileRevalidate<BlogComment[]>(
        cacheKey,
        `${BASE}/api/blog/${slug}/comments`,
        2 * 60 * 1000,
        (data) => { comments = data; loading = false; },
        (raw) => (raw as { comments: BlogComment[] }).comments ?? [],
      );
      if (!ok) throw new Error('Failed to load comments');
    } catch (e) {
      if (!comments.length) loading = false;
      toast(e instanceof Error ? e.message : 'Could not load comments');
    }
  }

  function startReply(parentId: string, parentName: string) {
    if (!loggedIn) {
      showLoginPopup = true;
      return;
    }
    replyToId = parentId;
    replyToName = parentName;
    replyText = '';
  }

  function cancelReply() {
    replyToId = null;
    replyToName = '';
    replyText = '';
  }

  async function submitComment(parentId: string | null = null) {
    if (!replyText.trim()) { toast('Comment cannot be empty'); return; }
    const t = $token;
    if (!t) { toast('You must be signed in to comment'); showLoginPopup = true; return; }

    try {
      const body: Record<string, unknown> = { content: replyText.trim() };
      if (parentId) body.parent_id = parentId;

      const r = await fetch(`${BASE}/api/blog/${slug}/comments`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${t}`,
        },
        body: JSON.stringify(body),
      });

      if (r.ok) {
        replyText = '';
        replyToId = null;
        replyToName = '';
        cacheDel(cacheKey);
        await fetchComments();
        toast('Comment posted', 'success');
      } else if (r.status === 401) {
        toast('Session expired. Please sign in again.');
        showLoginPopup = true;
      } else {
        const data = await r.json().catch(() => ({}));
        toast(data.error?.message || data.error || 'Failed to post comment');
      }
    } catch {
      toast('Network error. Check your connection and try again.');
    }
  }

  async function deleteComment(id: string) {
    const t = $token;
    if (!t) { toast('You must be signed in'); return; }

    try {
      const r = await fetch(`${BASE}/api/blog/${slug}/comments/${id}`, {
        method: 'DELETE',
        headers: { Authorization: `Bearer ${t}` },
      });
      if (r.ok) {
        cacheDel(cacheKey);
        await fetchComments();
        toast('Comment deleted', 'success');
      } else if (r.status === 401) {
        toast('Session expired. Please sign in again.');
        showLoginPopup = true;
      } else if (r.status === 403 || r.status === 404) {
        toast('You can only delete your own comments');
      } else {
        toast('Failed to delete comment');
      }
    } catch {
      toast('Network error. Check your connection and try again.');
    }
  }
</script>

<div class="comments-section">
  <h2 class="comments-heading">Comments ({comments.length})</h2>

  {#if !loggedIn}
    <div class="login-prompt">
      <button class="login-btn" onclick={loginWithGoogle}>Sign in with Google to comment</button>
    </div>
  {:else}
    <div class="comment-form">
      <textarea
        bind:value={replyText}
        placeholder="Write a comment..."
        rows="3"
      ></textarea>
      <div class="form-actions">
        <button class="submit-btn" onclick={() => submitComment(null)} disabled={!replyText.trim()}>Post Comment</button>
      </div>
    </div>
  {/if}

  {#if loading}
    <div class="loading-text">Loading comments...</div>
  {:else if flatList.length > 0}
    <div class="comments-list">
      {#each flatList as { comment, depth }}
        <div class="comment" style="padding-left: {depth * 16}px">
          <div class="comment-avatar">
            <img
              src={comment.avatar_url || ''}
              alt={comment.user_name}
              loading="lazy"
              onerror={(e) => { (e.currentTarget as HTMLElement).style.display = 'none'; }}
              onload={(e) => { (e.currentTarget.nextElementSibling as HTMLElement)?.style.setProperty('display', 'none'); }}
            />
            <div class="avatar-placeholder">{(comment.user_name || '?')[0]}</div>
          </div>
          <div class="comment-body">
            <div class="comment-header">
              <span class="comment-author">{comment.user_name}</span>
              <span class="comment-time">{comment.created_at}</span>
            </div>
            <div class="comment-text">{comment.content}</div>
            <div class="comment-actions">
              <button class="action-btn" onclick={() => startReply(comment.id, comment.user_name)}>Reply</button>
              {#if currentUserId === comment.user_id}
                <button class="action-btn delete" onclick={() => deleteComment(comment.id)}>Delete</button>
              {/if}
            </div>
          </div>
        </div>
      {/each}
    </div>
  {:else}
    <div class="empty-text">No comments yet. Be the first to share your thoughts!</div>
  {/if}
</div>

{#if replyToId}
  <div class="reply-form-overlay" onclick={cancelReply} role="presentation"></div>
  <div class="reply-form">
    <div class="reply-header">
      <span>Replying to <strong>{replyToName}</strong></span>
      <button class="cancel-btn" onclick={cancelReply}>Cancel</button>
    </div>
    <textarea
      bind:value={replyText}
      placeholder="Write your reply..."
      rows="2"
    ></textarea>
    <button class="submit-btn" onclick={() => submitComment(replyToId)} disabled={!replyText.trim()}>Post Reply</button>
  </div>
{/if}

{#if showLoginPopup}
  <div class="popup-overlay" onclick={() => showLoginPopup = false} role="presentation"></div>
  <div class="popup" role="dialog" aria-label="Sign in required">
    <div class="popup-icon">
      <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="12" cy="8" r="4"/>
        <path d="M20 21a8 8 0 0 0-16 0"/>
      </svg>
    </div>
    <h3 class="popup-title">Sign in Required</h3>
    <p class="popup-text">You need to sign in with Google to comment or reply.</p>
    <div class="popup-actions">
      <button class="popup-btn primary" onclick={() => { showLoginPopup = false; loginWithGoogle(); }}>Sign in with Google</button>
      <button class="popup-btn secondary" onclick={() => showLoginPopup = false}>Cancel</button>
    </div>
  </div>
{/if}

{#if toasts.length > 0}
  <div class="toast-container">
    {#each toasts as t (t.id)}
      <div class="toast" class:toast-error={t.type === 'error'} class:toast-success={t.type === 'success'}>
        <span>{t.message}</span>
        <button class="toast-close" onclick={() => { toasts = toasts.filter(x => x.id !== t.id); }}>×</button>
      </div>
    {/each}
  </div>
{/if}

<style>
  .comments-section {
    margin-top: 3rem;
    padding-top: 1.5rem;
    border-top: 1px solid var(--border-color);
  }
  .comments-heading {
    font-family: var(--font-heading);
    font-size: 1.2rem;
    color: var(--text-primary);
    margin: 0 0 1rem;
  }
  .login-prompt {
    margin-bottom: 1.5rem;
  }
  .login-btn {
    padding: 0.5rem 1rem;
    background: var(--accent);
    color: #fff;
    border: none;
    cursor: pointer;
    font-family: inherit;
    font-size: 0.85rem;
  }
  .login-btn:hover {
    opacity: 0.9;
  }
  .comment-form {
    margin-bottom: 1.5rem;
  }
  .comment-form textarea,
  .reply-form textarea {
    width: 100%;
    padding: 0.5rem;
    background: var(--surface-dark);
    border: 1px solid var(--border-color);
    color: var(--text-primary);
    font-family: inherit;
    font-size: 0.85rem;
    resize: vertical;
    box-sizing: border-box;
  }
  .comment-form textarea:focus,
  .reply-form textarea:focus {
    outline: none;
    border-color: var(--accent);
  }
  .form-actions {
    display: flex;
    justify-content: flex-end;
    margin-top: 0.4rem;
  }
  .submit-btn {
    padding: 0.35rem 0.9rem;
    background: var(--accent);
    color: #fff;
    border: none;
    cursor: pointer;
    font-family: inherit;
    font-size: 0.8rem;
  }
  .submit-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
  .submit-btn:hover:not(:disabled) {
    opacity: 0.9;
  }
  .loading-text, .empty-text {
    color: var(--text-muted);
    text-align: center;
    padding: 1.5rem;
    font-size: 0.85rem;
  }
  .comments-list {
    display: flex;
    flex-direction: column;
  }
  .comment {
    display: flex;
    gap: 0.6rem;
    padding: 0.75rem 0;
  }
  .comment-avatar {
    flex-shrink: 0;
    width: 32px;
    height: 32px;
    overflow: hidden;
    position: relative;
  }
  .comment-avatar img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    position: absolute;
    inset: 0;
    z-index: 1;
  }
  .comment-avatar img[src=""],
  .comment-avatar img[style*="display: none"] {
    display: none;
  }
  .avatar-placeholder {
    width: 100%;
    height: 100%;
    background: var(--border-color);
    color: var(--text-muted);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 0.75rem;
    font-weight: 600;
  }
  .comment-body {
    flex: 1;
    min-width: 0;
  }
  .comment-header {
    display: flex;
    gap: 0.5rem;
    align-items: center;
    margin-bottom: 0.2rem;
  }
  .comment-author {
    font-size: 0.85rem;
    font-weight: 600;
    color: var(--text-primary);
  }
  .comment-time {
    font-size: 0.7rem;
    color: var(--text-muted);
  }
  .comment-text {
    font-size: 0.85rem;
    color: var(--text-secondary);
    line-height: 1.5;
    white-space: pre-wrap;
    word-break: break-word;
  }
  .comment-actions {
    display: flex;
    gap: 0.5rem;
    margin-top: 0.3rem;
  }
  .action-btn {
    background: none;
    border: none;
    color: var(--text-muted);
    font-size: 0.75rem;
    cursor: pointer;
    padding: 0.15rem 0.3rem;
    font-family: inherit;
    transition: color 0.15s;
  }
  .action-btn:hover {
    color: var(--accent);
  }
  .action-btn.delete:hover {
    color: #ef4444;
  }
  .reply-form-overlay {
    position: fixed;
    inset: 0;
    background: rgba(0,0,0,0.4);
    z-index: 90;
  }
  .reply-form {
    position: fixed;
    bottom: 0;
    left: 50%;
    transform: translateX(-50%);
    width: min(640px, 100%);
    background: var(--bg-dark);
    border: 1px solid var(--border-color);
    padding: 1rem;
    z-index: 91;
    box-sizing: border-box;
    clip-path: polygon(12px 0, 100% 0, 100% 100%, 0 100%, 0 12px);
  }
  .reply-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 0.5rem;
    font-size: 0.85rem;
    color: var(--text-secondary);
  }
  .reply-header strong {
    color: var(--text-primary);
  }
  .cancel-btn {
    background: none;
    border: none;
    color: var(--text-muted);
    cursor: pointer;
    font-family: inherit;
    font-size: 0.8rem;
  }
  .cancel-btn:hover {
    color: var(--text-primary);
  }
  .reply-form .submit-btn {
    margin-top: 0.4rem;
  }

  .popup-overlay {
    position: fixed;
    inset: 0;
    background: rgba(0,0,0,0.55);
    z-index: 99;
  }
  .popup {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    background: var(--bg-dark);
    border: 1px solid var(--border-color);
    padding: 1.75rem;
    z-index: 100;
    width: min(360px, 90vw);
    text-align: center;
    box-sizing: border-box;
    clip-path: polygon(12px 0, 100% 0, 100% calc(100% - 12px), calc(100% - 12px) 100%, 0 100%, 0 12px);
  }
  .popup-icon {
    width: 48px;
    height: 48px;
    background: rgba(255,68,0,0.1);
    color: var(--accent);
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0 auto 0.75rem;
  }
  .popup-title {
    font-family: var(--font-heading);
    font-size: 1.1rem;
    color: var(--text-primary);
    margin: 0 0 0.4rem;
  }
  .popup-text {
    font-size: 0.85rem;
    color: var(--text-muted);
    margin: 0 0 1.25rem;
    line-height: 1.5;
  }
  .popup-actions {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
  .popup-btn {
    padding: 0.5rem;
    border: none;
    font-family: inherit;
    font-size: 0.85rem;
    cursor: pointer;
  }
  .popup-btn.primary {
    background: var(--accent);
    color: #fff;
  }
  .popup-btn.primary:hover {
    opacity: 0.9;
  }
  .popup-btn.secondary {
    background: transparent;
    color: var(--text-muted);
    border: 1px solid var(--border-color);
  }
  .popup-btn.secondary:hover {
    color: var(--text-primary);
  }

  .toast-container {
    position: fixed;
    bottom: 1.5rem;
    right: 1.5rem;
    z-index: 200;
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    max-width: 380px;
  }
  .toast {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.65rem 1rem;
    font-size: 0.8rem;
    animation: slideIn 0.2s ease-out;
    box-shadow: 0 4px 12px rgba(0,0,0,0.3);
  }
  .toast-error {
    background: #dc2626;
    color: #fff;
  }
  .toast-success {
    background: #16a34a;
    color: #fff;
  }
  .toast-close {
    background: none;
    border: none;
    color: inherit;
    cursor: pointer;
    font-size: 1rem;
    opacity: 0.7;
    padding: 0;
    line-height: 1;
    margin-left: auto;
    flex-shrink: 0;
  }
  .toast-close:hover {
    opacity: 1;
  }
  @keyframes slideIn {
    from { transform: translateX(100%); opacity: 0; }
    to { transform: translateX(0); opacity: 1; }
  }
</style>
