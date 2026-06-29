<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { browser } from '$app/environment';
  import { user, isLoggedIn, loginWithGoogle } from '$lib/stores/auth';
  import StarRating from './StarRating.svelte';
  import type { ProjectReview } from '$lib/types';

  let { projectName }: { projectName: string } = $props();

  const BASE = (import.meta.env.VITE_API_URL || '').replace(/\/$/, '');
  const MAX_COMMENT_LENGTH = 1000;

  let reviews = $state<ProjectReview[]>([]);
  let avgRating = $state(0);
  let count = $state(0);
  let userRating = $state(0);
  let userComment = $state('');
  let userHasReview = $state(false);
  let loading = $state(true);
  let loadError = $state('');
  let submitting = $state(false);
  let inlineSubmitting = $state(false);
  let submitError = $state('');
  let showModal = $state(false);
  let modalView = $state<'view' | 'write'>('view');
  let modalUserRating = $state(0);
  let modalUserComment = $state('');
  let charCount = $state(0);
  let textareaEl = $state<HTMLTextAreaElement | null>(null);
  let abortController = $state<AbortController | null>(null);
  let inlineRating = $state(0);

  let isConfigured = $derived(BASE.length > 0);

  function resetModalView() {
    modalUserRating = userRating;
    modalUserComment = userComment;
    charCount = userComment.length;
    submitError = '';
  }

  $effect(() => {
    if (!browser) return;
    if (showModal) {
      document.body.style.overflow = 'hidden';
      resetModalView();
      if (modalView === 'write') {
        requestAnimationFrame(() => textareaEl?.focus());
      }
      const handler = (e: KeyboardEvent) => {
        if (e.key === 'Escape') closeModal();
      };
      window.addEventListener('keydown', handler);
      return () => {
        document.body.style.overflow = '';
        window.removeEventListener('keydown', handler);
      };
    }
  });

  async function loadReviews() {
    if (!isConfigured) {
      loadError = 'API not configured';
      loading = false;
      return;
    }
    loading = true;
    loadError = '';
    abortController?.abort();
    abortController = new AbortController();
    try {
      const headers: Record<string, string> = {};
      const token = localStorage.getItem('portfolio_jwt');
      if (token) headers['Authorization'] = `Bearer ${token}`;

      const r = await fetch(
        `${BASE}/api/projects/${encodeURIComponent(projectName)}/reviews`,
        { headers, signal: abortController.signal }
      );
      if (!r.ok) {
        if (r.status === 404) {
          loadError = 'Project not found';
        } else if (r.status === 401) {
          loadError = 'Session expired — sign in again';
        } else {
          loadError = `Failed to load reviews (${r.status})`;
        }
        return;
      }
      const data = await r.json();
      reviews = data.reviews ?? [];
      avgRating = data.avg_rating ?? 0;
      count = data.count ?? 0;
      if (data.user_review?.rating) {
        userRating = data.user_review.rating;
        userComment = data.user_review.comment || '';
        userHasReview = true;
      }
    } catch (err) {
      if (err instanceof DOMException && err.name === 'AbortError') return;
      loadError = 'Network error — check your connection';
    } finally {
      loading = false;
      abortController = null;
    }
  }

  async function submitReview() {
    if (!isConfigured || modalUserRating < 1) return;
    submitting = true;
    submitError = '';
    try {
      const token = localStorage.getItem('portfolio_jwt');
      const r = await fetch(
        `${BASE}/api/projects/${encodeURIComponent(projectName)}/reviews`,
        {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            ...(token ? { 'Authorization': `Bearer ${token}` } : {}),
          },
          body: JSON.stringify({
            rating: modalUserRating,
            comment: modalUserComment.slice(0, MAX_COMMENT_LENGTH),
          }),
        }
      );
      if (r.ok) {
        userRating = modalUserRating;
        userComment = modalUserComment.slice(0, MAX_COMMENT_LENGTH);
        userHasReview = true;
        showModal = false;
        await loadReviews();
      } else if (r.status === 401) {
        submitError = 'Session expired. Sign in again.';
      } else {
        const msg = await r.text().catch(() => '');
        submitError = msg || `Submission failed (${r.status})`;
      }
    } catch {
      submitError = 'Network error. Please try again.';
    } finally {
      submitting = false;
    }
  }

  async function deleteReview() {
    if (!isConfigured) return;
    try {
      const token = localStorage.getItem('portfolio_jwt');
      const r = await fetch(
        `${BASE}/api/projects/${encodeURIComponent(projectName)}/reviews`,
        {
          method: 'DELETE',
          headers: token ? { 'Authorization': `Bearer ${token}` } : {},
        }
      );
      if (r.ok) {
        userHasReview = false;
        userRating = 0;
        userComment = '';
        await loadReviews();
      }
    } catch {}
  }

  function openModal(mode: 'view' | 'write') {
    modalView = mode;
    showModal = true;
  }

  function closeModal() {
    showModal = false;
  }

  function onBackdropClick(e: MouseEvent) {
    if (e.target === e.currentTarget) closeModal();
  }

  function handleRate(n: number) {
    modalUserRating = n;
  }

  function handleInlineRate(n: number) {
    inlineRating = n;
  }

  async function submitInlineReview() {
    if (!isConfigured || inlineRating < 1) return;
    inlineSubmitting = true;
    try {
      const token = localStorage.getItem('portfolio_jwt');
      const r = await fetch(
        `${BASE}/api/projects/${encodeURIComponent(projectName)}/reviews`,
        {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            ...(token ? { 'Authorization': `Bearer ${token}` } : {}),
          },
          body: JSON.stringify({ rating: inlineRating }),
        }
      );
      if (r.ok) {
        userRating = inlineRating;
        userComment = '';
        userHasReview = true;
        inlineRating = 0;
        await loadReviews();
      }
    } catch {} finally {
      inlineSubmitting = false;
    }
  }

  function handleCommentInput() {
    if (modalUserComment.length > MAX_COMMENT_LENGTH) {
      modalUserComment = modalUserComment.slice(0, MAX_COMMENT_LENGTH);
    }
    charCount = modalUserComment.length;
  }

  function handleAvatarError(e: Event) {
    const img = e.currentTarget as HTMLImageElement;
    img.style.display = 'none';
    const fallback = img.nextElementSibling as HTMLElement | null;
    if (fallback) fallback.style.display = 'flex';
  }

  function reviewAuthorInitials(name: string): string {
    return name
      .split(' ')
      .map(w => w[0])
      .filter(Boolean)
      .slice(0, 2)
      .join('')
      .toUpperCase() || '?';
  }

  function timeAgo(dateStr: string): string {
    const now = Date.now();
    const then = new Date(dateStr).getTime();
    if (!then) return '';
    const diff = Math.max(0, now - then);
    const mins = Math.floor(diff / 60000);
    if (mins < 1) return 'just now';
    if (mins < 60) return `${mins}m ago`;
    const hours = Math.floor(mins / 60);
    if (hours < 24) return `${hours}h ago`;
    const days = Math.floor(hours / 24);
    if (days < 30) return `${days}d ago`;
    const months = Math.floor(days / 30);
    if (months < 12) return `${months}mo ago`;
    return `${Math.floor(months / 12)}y ago`;
  }

  onMount(loadReviews);

  onDestroy(() => {
    if (browser) document.body.style.overflow = '';
    abortController?.abort();
  });
</script>

<div class="review-section">
  <div class="review-summary">
    {#if loading}
      <span class="review-count sm">Loading...</span>
    {:else if loadError}
      <span class="review-load-error">{loadError}</span>
      <button class="review-btn retry" onclick={loadReviews}>
        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="23 4 23 10 17 10"></polyline><path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"></path></svg>
        Retry
      </button>
    {:else if count > 0}
      <span class="review-count">{count} review{count !== 1 ? 's' : ''}</span>
    {:else}
      <span class="review-count">No reviews yet</span>
    {/if}
  </div>

  {#if !loading && !loadError}
    <div class="review-actions-row">
      {#if $isLoggedIn}
        {#if userHasReview}
          <div class="user-review-inline">
            <span class="user-review-label">Your review:</span>
            <StarRating value={userRating} size="sm" />
            <button class="review-btn ghost" onclick={() => openModal('write')}>
              <span class="edit-link">Edit</span>
            </button>
            <button class="review-btn delete-sm" onclick={deleteReview} aria-label="Delete your review">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path></svg>
            </button>
            {#if userComment}
              <button class="user-comment-preview" onclick={() => openModal('view')}>
                &ldquo;{userComment.slice(0, 70)}{userComment.length > 70 ? '...' : ''}&rdquo;
                {#if userComment.length > 70}
                  <span class="read-more">Read more</span>
                {/if}
              </button>
            {/if}
          </div>
        {:else}
          <div class="inline-rate-row">
            <StarRating interactive value={inlineRating} onRate={handleInlineRate} size="sm" />
            <button
              class="review-btn inline-submit"
              onclick={submitInlineReview}
              disabled={inlineSubmitting || inlineRating < 1}
            >
              {#if inlineSubmitting}
                Submitting...
              {:else}
                Submit
              {/if}
            </button>
            <button class="review-btn write-with-comment" onclick={() => { if (inlineRating > 0) modalUserRating = inlineRating; openModal('write'); }}>
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M12 20h9"></path><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"></path></svg>
              Write a Review
            </button>
          </div>
        {/if}
      {:else}
        <button class="review-btn sign-in" onclick={loginWithGoogle}>
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"></path><polyline points="10 17 15 12 10 7"></polyline><line x1="15" y1="12" x2="3" y2="12"></line></svg>
          Sign in to rate
        </button>
      {/if}

      {#if count > 0}
        <button class="review-btn read-all" onclick={() => openModal('view')}>
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path><circle cx="12" cy="12" r="3"></circle></svg>
          Read {count} review{count !== 1 ? 's' : ''}
        </button>
      {/if}
    </div>
  {/if}
</div>

{#if showModal}
  <div class="modal-backdrop" onclick={onBackdropClick} role="dialog" aria-modal="true" aria-label="{modalView === 'write' ? 'Write a review for ' : 'Reviews for '}{projectName}">
    <div class="modal-window">
      <div class="modal-titlebar">
        <div class="window-dots">
          <span class="dot dot-close" onclick={closeModal} role="button" tabindex="0" aria-label="Close"
            onkeydown={(e) => e.key === 'Enter' && closeModal()}></span>
          <span class="dot dot-minimize"></span>
          <span class="dot dot-maximize"></span>
        </div>
        <span class="modal-title" title={projectName}>
          {modalView === 'write' ? 'Review' : 'Reviews'} &mdash; {projectName}
        </span>
      </div>
      <div class="modal-body">
        {#if modalView === 'write'}
          {#if submitError}
            <div class="modal-error">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="15" y1="9" x2="9" y2="15"></line><line x1="9" y1="9" x2="15" y2="15"></line></svg>
              {submitError}
            </div>
          {/if}
          <div class="modal-rating-row">
            <span class="modal-rating-label">Rating</span>
            <StarRating interactive value={modalUserRating} onRate={handleRate} />
            {#if modalUserRating < 1}
              <span class="modal-rating-hint">Select a rating</span>
            {/if}
          </div>
          <div class="modal-field">
            <textarea
              bind:this={textareaEl}
              bind:value={modalUserComment}
              oninput={handleCommentInput}
              placeholder="Optional: share your thoughts..."
              rows="4"
              class="review-textarea"
              maxlength={MAX_COMMENT_LENGTH}
              disabled={submitting}
            ></textarea>
            <span class="char-count" class:near-limit={charCount > MAX_COMMENT_LENGTH * 0.9} class:over={charCount >= MAX_COMMENT_LENGTH}>
              {charCount}/{MAX_COMMENT_LENGTH}
            </span>
          </div>
          <div class="modal-footer-actions">
            <button class="review-btn cancel" onclick={() => modalView = 'view'} disabled={submitting}>Back</button>
            <button
              class="review-btn submit"
              onclick={submitReview}
              disabled={submitting || modalUserRating < 1}
            >
              {#if submitting}
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" class="spin"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
                Submitting...
              {:else}
                Submit Review
              {/if}
            </button>
          </div>
        {:else}
          {#if loading}
            <div class="modal-loading">Loading reviews...</div>
          {:else if reviews.length > 0}
            <div class="reviews-list">
              {#each reviews as review (review.id || review.user_id)}
                <div class="review-item">
                  <div class="review-header">
                    <div class="review-avatar-wrap">
                      {#if review.avatar_url}
                        <img
                          src={review.avatar_url}
                          alt="{review.user_name}'s avatar"
                          class="review-avatar"
                          onerror={handleAvatarError}
                          loading="lazy"
                        />
                      {/if}
                      <span
                        class="review-avatar-fallback"
                        style:display={review.avatar_url ? 'none' : 'flex'}
                      >{reviewAuthorInitials(review.user_name)}</span>
                    </div>
                    <span class="review-author" title={review.user_name}>{review.user_name}</span>
                    <StarRating value={review.rating} size="sm" />
                    {#if review.created_at}
                      <span class="review-time">{timeAgo(review.created_at)}</span>
                    {/if}
                  </div>
                  {#if review.comment}
                    <p class="review-text">{review.comment}</p>
                  {/if}
                </div>
              {/each}
            </div>
          {:else}
            <div class="modal-empty">No reviews yet. Be the first!</div>
          {/if}
          <div class="modal-footer-actions">
            {#if $isLoggedIn && !userHasReview}
              <button class="review-btn write" onclick={() => modalView = 'write'}>
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M12 20h9"></path><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"></path></svg>
                Write a Review
              </button>
            {/if}
          </div>
        {/if}
      </div>
    </div>
  </div>
{/if}

<style>
  .review-section {
    margin-top: 0.75rem;
    padding-top: 0.75rem;
    border-top: 1px solid rgba(255,255,255,0.08);
    position: relative;
  }
  .review-summary {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    margin-bottom: 0.5rem;
    min-height: 28px;
  }
  .review-count {
    font-size: 0.75rem;
    color: #71717a;
    font-weight: 500;
    letter-spacing: 0.03em;
  }
  .review-count.sm {
    font-size: 0.7rem;
  }
  .review-load-error {
    font-size: 0.7rem;
    color: #ef4444;
    font-weight: 500;
  }
  .review-actions-row {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 0.4rem;
  }
  .review-btn {
    font-size: 0.65rem;
    font-weight: 600;
    padding: 0.3rem 0.65rem;
    border: 1.5px solid rgba(255,255,255,0.15);
    background: rgba(255,255,255,0.04);
    color: #e4e4e7;
    cursor: pointer;
    font-family: inherit;
    letter-spacing: 0.04em;
    text-transform: uppercase;
    transition: all 0.15s ease;
    display: inline-flex;
    align-items: center;
    gap: 0.35rem;
    white-space: nowrap;
  }
  .review-btn:hover:not(:disabled) {
    background: rgba(255,255,255,0.1);
  }
  .review-btn:disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }
  .review-btn:focus-visible {
    outline: 2px solid var(--accent);
    outline-offset: 2px;
  }
  .review-btn.ghost {
    border-color: transparent;
    background: none;
    color: var(--accent);
    padding: 0.1rem 0.3rem;
    font-size: 0.6rem;
    text-transform: uppercase;
    font-weight: 600;
    letter-spacing: 0.04em;
  }
  .review-btn.ghost:hover:not(:disabled) {
    opacity: 0.8;
  }
  .review-btn.write {
    border-color: var(--accent);
    color: var(--accent);
    background: rgba(255, 68, 0, 0.06);
  }
  .review-btn.write:hover:not(:disabled) {
    background: rgba(255, 68, 0, 0.15);
  }
  .review-btn.sign-in {
    border-color: rgba(255,255,255,0.1);
    color: #71717a;
    background: none;
  }
  .review-btn.sign-in:hover:not(:disabled) {
    border-color: var(--accent);
    color: var(--accent);
  }
  .review-btn.read-all {
    border: none;
    background: none;
    color: #71717a;
    padding: 0.2rem 0.3rem;
    font-size: 0.6rem;
    font-weight: 500;
    letter-spacing: 0;
    text-transform: none;
    margin-left: auto;
  }
  .review-btn.read-all:hover:not(:disabled) {
    color: var(--accent);
  }
  .inline-rate-row {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    flex-wrap: wrap;
  }
  .review-btn.inline-submit {
    font-size: 0.6rem;
    padding: 0.25rem 0.6rem;
    background: var(--accent);
    color: #000;
    border-color: var(--accent);
    font-weight: 700;
    min-width: 60px;
    justify-content: center;
  }
  .review-btn.inline-submit:disabled {
    opacity: 0.35;
    cursor: not-allowed;
  }
  .review-btn.inline-submit:hover:not(:disabled) {
    opacity: 0.9;
  }
  .review-btn.write-with-comment {
    border-color: rgba(255,255,255,0.1);
    padding: 0.25rem 0.6rem;
    font-size: 0.6rem;
  }
  .review-btn.write-with-comment:hover:not(:disabled) {
    border-color: var(--accent);
    color: var(--accent);
  }
  .review-btn.delete-sm {
    border: none;
    background: none;
    color: rgba(239,68,68,0.4);
    padding: 0.2rem;
    cursor: pointer;
  }
  .review-btn.delete-sm:hover:not(:disabled) {
    color: #ef4444;
  }
  .review-btn.retry {
    display: inline-flex;
    align-items: center;
    gap: 0.3rem;
    border-color: rgba(239,68,68,0.3);
    color: #ef4444;
    font-size: 0.65rem;
    padding: 0.25rem 0.6rem;
    margin-left: auto;
  }
  .review-btn.retry:hover {
    background: rgba(239,68,68,0.1);
    border-color: #ef4444;
  }
  .user-review-inline {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 0.25rem;
    min-width: 0;
    flex: 1;
  }
  .user-review-label {
    font-size: 0.65rem;
    color: #71717a;
    font-weight: 500;
    text-transform: uppercase;
    letter-spacing: 0.04em;
  }
  .user-comment-preview {
    font-size: 0.7rem;
    color: #71717a;
    font-style: italic;
    background: none;
    border: none;
    cursor: pointer;
    font-family: inherit;
    padding: 0.15rem 0.3rem;
    text-align: left;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 180px;
    transition: color 0.15s;
  }
  .user-comment-preview:hover {
    color: #e4e4e7;
  }
  .read-more {
    color: var(--accent);
    font-style: normal;
    font-weight: 600;
    font-size: 0.6rem;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    margin-left: 0.15rem;
  }

  /* ── Modal ── */
  .modal-backdrop {
    position: fixed;
    inset: 0;
    z-index: 10000;
    background: rgba(0, 0, 0, 0.7);
    backdrop-filter: blur(8px);
    -webkit-backdrop-filter: blur(8px);
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 1rem;
  }
  .modal-window {
    width: 100%;
    max-width: 480px;
    max-height: calc(100dvh - 2rem);
    background: #15181e;
    border: 1px solid rgba(255,255,255,0.1);
    box-shadow: 0 24px 64px -12px rgba(0,0,0,0.6);
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }
  .modal-titlebar {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.65rem 1rem;
    background: rgba(255,255,255,0.03);
    border-bottom: 1px solid rgba(255,255,255,0.06);
    user-select: none;
    flex-shrink: 0;
  }
  .window-dots {
    display: flex;
    gap: 6px;
    flex-shrink: 0;
  }
  .dot {
    width: 12px;
    height: 12px;
    border-radius: 50%;
    display: block;
  }
  .dot-close {
    background: #ff5f56;
    cursor: pointer;
    transition: filter 0.1s;
  }
  .dot-close:hover {
    filter: brightness(1.2);
  }
  .dot-minimize {
    background: #ffbd2e;
  }
  .dot-maximize {
    background: #27c93f;
  }
  .modal-title {
    font-size: 0.75rem;
    color: #a1a1aa;
    font-family: var(--font-body);
    letter-spacing: 0.02em;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  .modal-body {
    padding: 1.25rem;
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    overflow-y: auto;
    flex: 1;
  }
  .modal-error {
    display: flex;
    align-items: center;
    gap: 0.4rem;
    font-size: 0.75rem;
    color: #ef4444;
    background: rgba(239,68,68,0.08);
    border: 1px solid rgba(239,68,68,0.2);
    padding: 0.5rem 0.65rem;
  }
  .modal-error svg {
    flex-shrink: 0;
  }
  .modal-rating-row {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    flex-wrap: wrap;
  }
  .modal-rating-label {
    font-size: 0.7rem;
    color: #a1a1aa;
    text-transform: uppercase;
    letter-spacing: 0.06em;
    font-weight: 600;
  }
  .modal-rating-hint {
    font-size: 0.65rem;
    color: #ef4444;
    font-weight: 500;
  }
  .modal-field {
    position: relative;
  }
  .char-count {
    position: absolute;
    bottom: 0.4rem;
    right: 0.5rem;
    font-size: 0.6rem;
    color: #52525b;
    pointer-events: none;
    font-family: monospace;
  }
  .char-count.near-limit {
    color: #f59e0b;
  }
  .char-count.over {
    color: #ef4444;
  }
  .modal-footer-actions {
    display: flex;
    justify-content: flex-end;
    align-items: center;
    gap: 0.5rem;
    padding-top: 0.25rem;
  }
  .modal-loading {
    padding: 2rem 0;
    text-align: center;
    color: #71717a;
    font-size: 0.8rem;
  }
  .modal-empty {
    padding: 2rem 0;
    text-align: center;
    color: #71717a;
    font-size: 0.8rem;
  }
  .review-textarea {
    background: rgba(0,0,0,0.35);
    border: 1px solid rgba(255,255,255,0.1);
    color: #e4e4e7;
    padding: 0.6rem;
    font-size: 0.85rem;
    font-family: inherit;
    resize: vertical;
    transition: border-color 0.15s ease, box-shadow 0.15s ease;
    line-height: 1.5;
    width: 100%;
    box-sizing: border-box;
  }
  .review-textarea:focus {
    outline: none;
    border-color: var(--accent);
    box-shadow: 0 0 0 1px var(--accent);
  }
  .review-textarea::placeholder {
    color: #52525b;
  }
  .review-textarea:disabled {
    opacity: 0.5;
  }

  /* ── Avatar ── */
  .review-avatar-wrap {
    width: 24px;
    height: 24px;
    flex-shrink: 0;
    position: relative;
    border-radius: 50%;
    overflow: hidden;
    border: 1px solid rgba(255,255,255,0.1);
  }
  .review-avatar {
    width: 100%;
    height: 100%;
    object-fit: cover;
    display: block;
  }
  .review-avatar-fallback {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--accent);
    color: #000;
    font-size: 0.6rem;
    font-weight: 700;
    letter-spacing: 0.02em;
  }
  .reviews-list {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
  .review-item {
    padding: 0.5rem 0.65rem;
    background: rgba(255,255,255,0.03);
    border: 1px solid rgba(255,255,255,0.04);
    transition: border-color 0.15s ease;
  }
  .review-item:hover {
    border-color: rgba(255,255,255,0.1);
  }
  .review-header {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    min-width: 0;
  }
  .review-author {
    font-size: 0.8rem;
    font-weight: 600;
    color: #e4e4e7;
    letter-spacing: 0.02em;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  .review-time {
    font-size: 0.6rem;
    color: #52525b;
    margin-left: auto;
    flex-shrink: 0;
    font-family: monospace;
  }
  .review-text {
    font-size: 0.8rem;
    color: #a1a1aa;
    margin: 0.35rem 0 0;
    line-height: 1.5;
    word-break: break-word;
  }

  .review-btn.cancel {
    border-color: rgba(255,255,255,0.15);
  }
  .review-btn.cancel:hover:not(:disabled) {
    border-color: rgba(255,255,255,0.3);
    color: #fff;
  }
  .review-btn.submit {
    background: var(--accent);
    color: #000;
    border-color: var(--accent);
    font-weight: 700;
    min-width: 130px;
    justify-content: center;
  }
  .review-btn.submit:hover:not(:disabled) {
    opacity: 0.9;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }
  .spin {
    animation: spin 0.8s linear infinite;
  }

  /* ── Light Mode Overrides ── */
  :global(body.light-mode) .review-section {
    border-top-color: rgba(0,0,0,0.06);
  }
  :global(body.light-mode) .review-summary {
    color: #52525b;
  }
  :global(body.light-mode) .review-count {
    color: #52525b;
  }
  :global(body.light-mode) .user-review-label {
    color: #52525b;
  }
  :global(body.light-mode) .review-btn {
    color: #27272a;
    border-color: rgba(0,0,0,0.15);
    background: rgba(0,0,0,0.03);
  }
  :global(body.light-mode) .review-btn:hover:not(:disabled) {
    background: rgba(0,0,0,0.06);
  }
  :global(body.light-mode) .review-btn.write {
    color: var(--accent);
    border-color: rgba(200,50,0,0.3);
    background: rgba(200,50,0,0.04);
  }
  :global(body.light-mode) .review-btn.write:hover:not(:disabled) {
    background: rgba(200,50,0,0.1);
  }
  :global(body.light-mode) .review-btn.sign-in {
    color: #71717a;
    border-color: rgba(0,0,0,0.1);
  }
  :global(body.light-mode) .review-btn.sign-in:hover:not(:disabled) {
    border-color: var(--accent);
    color: var(--accent);
  }
  :global(body.light-mode) .review-btn.read-all {
    color: #71717a;
  }
  :global(body.light-mode) .review-btn.read-all:hover:not(:disabled) {
    color: var(--accent);
  }
  :global(body.light-mode) .review-btn.ghost {
    color: #27272a;
  }
  :global(body.light-mode) .review-btn.ghost:hover:not(:disabled) {
    color: #18181b;
  }
  :global(body.light-mode) .review-btn.delete-sm {
    color: rgba(220,38,38,0.4);
  }
  :global(body.light-mode) .review-btn.delete-sm:hover:not(:disabled) {
    color: #dc2626;
  }
  :global(body.light-mode) .review-btn.cancel {
    color: #52525b;
    border-color: rgba(0,0,0,0.12);
  }
  :global(body.light-mode) .review-btn.cancel:hover:not(:disabled) {
    border-color: rgba(0,0,0,0.3);
    color: #18181b;
  }
  :global(body.light-mode) .review-btn.submit {
    color: #fff;
    background: var(--accent);
    border-color: var(--accent);
  }
  :global(body.light-mode) .review-btn.inline-submit {
    color: #fff;
    background: var(--accent);
    border-color: var(--accent);
  }
  :global(body.light-mode) .review-btn.inline-submit:disabled {
    opacity: 0.35;
  }
  :global(body.light-mode) .review-btn.write-with-comment {
    border-color: rgba(0,0,0,0.1);
    color: #52525b;
  }
  :global(body.light-mode) .review-btn.write-with-comment:hover:not(:disabled) {
    border-color: var(--accent);
    color: var(--accent);
  }
  :global(body.light-mode) .review-btn.retry {
    color: #dc2626;
    border-color: rgba(220,38,38,0.3);
  }
  :global(body.light-mode) .review-btn.retry:hover {
    background: rgba(220,38,38,0.06);
    border-color: #dc2626;
  }
  :global(body.light-mode) .review-item {
    background: rgba(0,0,0,0.02);
    border-color: rgba(0,0,0,0.04);
  }
  :global(body.light-mode) .review-item:hover {
    border-color: rgba(0,0,0,0.1);
  }
  :global(body.light-mode) .review-author {
    color: #18181b;
  }
  :global(body.light-mode) .review-time {
    color: #a1a1aa;
  }
  :global(body.light-mode) .review-text {
    color: #52525b;
  }
  :global(body.light-mode) .modal-window {
    background: #f0f2f5;
    border-color: rgba(0,0,0,0.12);
  }
  :global(body.light-mode) .modal-titlebar {
    background: rgba(0,0,0,0.03);
    border-color: rgba(0,0,0,0.06);
  }
  :global(body.light-mode) .modal-title {
    color: #52525b;
  }
  :global(body.light-mode) .review-textarea {
    color: #18181b;
    background: #fff;
    border-color: rgba(0,0,0,0.12);
  }
  :global(body.light-mode) .review-textarea:focus {
    border-color: var(--accent);
    box-shadow: 0 0 0 1px var(--accent);
  }
  :global(body.light-mode) .review-avatar-fallback {
    color: #fff;
  }
  :global(body.light-mode) .review-avatar-wrap {
    border-color: rgba(0,0,0,0.1);
  }
</style>
