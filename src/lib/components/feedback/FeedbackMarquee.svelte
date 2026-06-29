<script lang="ts">
  import type { MarqueeItem } from '$lib/types';
  import StarRating from './StarRating.svelte';

  let { items = [] }: { items: MarqueeItem[] } = $props();

  let duped = $derived([...items, ...items, ...items]);

  function initials(name: string): string {
    return name
      .split(' ')
      .map(w => w[0])
      .filter(Boolean)
      .slice(0, 2)
      .join('')
      .toUpperCase() || '?';
  }

  function handleAvatarError(e: Event) {
    const img = e.currentTarget as HTMLImageElement;
    img.style.display = 'none';
    const wrap = img.parentElement as HTMLElement;
    const fallback = wrap?.querySelector('.marquee-avatar-fallback') as HTMLElement;
    if (fallback) fallback.style.display = 'flex';
  }
</script>

{#if items.length > 0}
  <div class="marquee-section">
    <h3 class="marquee-title">Live Feedback</h3>
    <div class="marquee-track">
      <div class="marquee-content">
        {#each duped as item, i}
          <div class="marquee-item" title={item.comment ? `"${item.comment}"` : ''}>
            <div class="marquee-avatar-wrap">
              {#if item.avatar_url}
                <img
                  src={item.avatar_url}
                  alt=""
                  class="marquee-avatar"
                  onerror={handleAvatarError}
                  loading="lazy"
                />
              {/if}
              <span
                class="marquee-avatar-fallback"
                style:display={item.avatar_url ? 'none' : 'flex'}
              >{initials(item.user_name)}</span>
            </div>
            <span class="marquee-user">{item.user_name}</span>
            <StarRating value={item.rating} size="sm" />
            <span class="marquee-project">on {item.project_name}</span>
            {#if item.comment}
              <span class="marquee-comment">&ldquo;{item.comment.slice(0, 60)}{item.comment.length > 60 ? '...' : ''}&rdquo;</span>
            {/if}
          </div>
        {/each}
      </div>
    </div>
  </div>
{/if}

<style>
  .marquee-section {
    margin: 2rem 0;
    overflow: hidden;
  }
  .marquee-title {
    font-family: var(--font-heading);
    font-size: 0.75rem;
    text-transform: uppercase;
    letter-spacing: 0.15em;
    color: #71717a;
    margin-bottom: 0.75rem;
    text-align: center;
  }
  .marquee-track {
    overflow: hidden;
    mask-image: linear-gradient(to right, transparent, black 5%, black 95%, transparent);
    -webkit-mask-image: linear-gradient(to right, transparent, black 5%, black 95%, transparent);
  }
  .marquee-content {
    display: flex;
    gap: 1rem;
    width: max-content;
    animation: scroll 30s linear infinite;
  }
  @keyframes scroll {
    0% { transform: translateX(0); }
    100% { transform: translateX(-33.33%); }
  }
  .marquee-item {
    display: flex;
    align-items: center;
    gap: 0.35rem;
    padding: 0.4rem 0.75rem;
    background: rgba(255,255,255,0.03);
    border: 1px solid rgba(255,255,255,0.06);
    border-radius: 0;
    white-space: nowrap;
    font-size: 0.75rem;
    cursor: default;
  }
  .marquee-avatar-wrap {
    width: 18px;
    height: 18px;
    flex-shrink: 0;
    position: relative;
    border-radius: 50%;
    overflow: hidden;
    border: 1px solid rgba(255,255,255,0.1);
  }
  .marquee-avatar {
    width: 100%;
    height: 100%;
    object-fit: cover;
    display: block;
  }
  .marquee-avatar-fallback {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--accent);
    color: #000;
    font-size: 0.5rem;
    font-weight: 700;
  }
  .marquee-user {
    color: #e4e4e7;
    font-weight: 500;
  }
  .marquee-project {
    color: var(--accent);
  }
  .marquee-comment {
    color: #71717a;
    font-style: italic;
  }

  :global(body.light-mode) .marquee-item {
    background: rgba(0,0,0,0.02);
    border-color: rgba(0,0,0,0.06);
  }
  :global(body.light-mode) .marquee-user {
    color: #27272a;
  }
  :global(body.light-mode) .marquee-comment {
    color: #52525b;
  }
  :global(body.light-mode) .marquee-avatar-wrap {
    border-color: rgba(0,0,0,0.1);
  }
  :global(body.light-mode) .marquee-avatar-fallback {
    color: #fff;
  }
</style>
