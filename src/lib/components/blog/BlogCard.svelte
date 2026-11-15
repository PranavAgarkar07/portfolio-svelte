<script lang="ts">
  import { base } from '$app/paths';
  import type { BlogPost } from '$lib/types';
  import { formatDate } from '$lib/utils/format';
  import { readingTime } from '$lib/utils/markdown';

  let {
    post,
    index = 0,
    featured = false,
  }: { post: BlogPost; index?: number; featured?: boolean } = $props();
</script>

<a href="{base}/blog/{post.slug}" class="blog-card" class:featured>
  {#if post.cover_image}
    <div class="card-image">
      <img src={`${base}${post.cover_image}`} alt={post.title} loading="lazy" decoding="async" />
      <span class="card-chip card-chip-idx" aria-hidden="true">
        {String(index + 1).padStart(2, "0")}
      </span>
      {#if featured}
        <span class="card-chip card-chip-latest">Latest</span>
      {/if}
    </div>
  {/if}
  <div class="card-body">
    <div class="card-toprow">
      {#if post.tags?.length}
        <div class="card-tags">
          {#each post.tags.slice(0, 3) as tag (tag)}
            <span class="tag">{tag}</span>
          {/each}
        </div>
      {/if}
      <time class="card-date">{formatDate(post.published_at || post.created_at)}</time>
    </div>
    <h3 class="card-title">{post.title}</h3>
    {#if post.excerpt}
      <p class="card-excerpt">{post.excerpt}</p>
    {/if}
    <div class="card-meta">
      <span class="card-meta-name">{post.author_name}</span>
      <span class="card-meta-dot" aria-hidden="true">&middot;</span>
      <span class="card-meta-item">{readingTime(post.content_md)}</span>
    </div>
  </div>
</a>

<style>
  .blog-card {
    display: block;
    text-decoration: none;
    position: relative;
    background: var(--card-bg);
    border: 1px solid var(--border-color);
    border-left: 3px solid var(--border-color);
    box-shadow: var(--card-shadow);
    transition:
      transform 0.2s ease,
      box-shadow 0.2s ease,
      border-left-color 0.2s ease;
  }
  .blog-card:hover {
    border-left-color: var(--accent);
    transform: translateY(-4px);
    /* Neobrutalist offset shadow — same language as the certificate cards */
    box-shadow: 6px 6px 0 rgba(0, 0, 0, 0.45);
  }
  .blog-card:focus-visible {
    outline: 2px solid var(--accent);
    outline-offset: 3px;
  }
  .card-image {
    position: relative;
    width: 100%;
    aspect-ratio: 2 / 1;
    overflow: hidden;
    background: var(--surface-dark);
    border-bottom: 1px solid var(--border-color);
  }
  .card-image img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    transition: transform 0.3s ease;
  }
  .blog-card:hover .card-image img {
    transform: scale(1.04);
  }

  /* ── Persistent chips (index / latest) ── */
  .card-chip {
    position: absolute;
    top: 10px;
    z-index: 2;
    font-family: var(--font-body);
    font-size: 0.55rem;
    font-weight: 700;
    letter-spacing: 0.12em;
    text-transform: uppercase;
    line-height: 1;
    padding: 4px 8px;
    color: rgba(255, 255, 255, 0.85);
    background: #0a0d13;
    border: 1px solid rgba(255, 255, 255, 0.14);
    pointer-events: none;
  }
  .card-chip-idx {
    left: 10px;
    color: rgba(255, 255, 255, 0.6);
  }
  .card-chip-latest {
    right: 10px;
    color: var(--accent);
    border-color: rgba(255, 68, 0, 0.55);
    background: rgba(255, 68, 0, 0.12);
  }
  :global(body.light-mode) .card-chip {
    color: rgba(0, 0, 0, 0.65);
    background: rgba(255, 255, 255, 0.92);
    border-color: rgba(0, 0, 0, 0.15);
  }
  :global(body.light-mode) .card-chip-latest {
    color: var(--light-accent, #d94100);
    border-color: rgba(217, 65, 0, 0.45);
    background: rgba(217, 65, 0, 0.06);
  }

  .card-body {
    padding: 1.25rem;
  }
  .card-toprow {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 0.75rem;
    flex-wrap: wrap;
    margin-bottom: 0.6rem;
  }
  .card-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 0.4rem;
  }
  .tag {
    font-size: 0.62rem;
    padding: 0.2rem 0.5rem;
    background: rgba(255, 68, 0, 0.1);
    border: 1px solid rgba(255, 68, 0, 0.3);
    color: var(--accent);
    letter-spacing: 0.06em;
    text-transform: uppercase;
    font-weight: 600;
  }
  .card-date {
    font-family: var(--font-heading);
    font-size: 0.65rem;
    letter-spacing: 0.1em;
    text-transform: uppercase;
    color: var(--text-muted);
  }
  .card-title {
    font-family: var(--font-heading);
    font-size: 1.3rem;
    color: var(--text-primary);
    margin: 0 0 0.5rem;
    line-height: 1.3;
    transition: color 0.2s ease;
  }
  .blog-card:hover .card-title {
    color: var(--accent);
  }
  .card-excerpt {
    color: var(--text-secondary);
    font-size: 0.9rem;
    line-height: 1.6;
    margin: 0 0 0.9rem;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }
  .card-meta {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 0.45rem;
    font-family: var(--font-heading);
    font-size: 0.7rem;
    letter-spacing: 0.06em;
    text-transform: uppercase;
    color: var(--text-muted);
  }
  .card-meta-name {
    color: var(--accent);
  }
  .card-meta-dot {
    opacity: 0.4;
  }

  /* ── Featured (first post): horizontal split on wide screens ── */
  @media (min-width: 700px) {
    .blog-card.featured {
      display: grid;
      grid-template-columns: 5fr 7fr;
    }
    .blog-card.featured .card-image {
      aspect-ratio: auto;
      height: 100%;
      min-height: 240px;
      border-bottom: none;
      border-right: 1px solid var(--border-color);
    }
    .blog-card.featured .card-body {
      padding: 1.75rem;
      display: flex;
      flex-direction: column;
      justify-content: center;
    }
    .blog-card.featured .card-title {
      font-size: 1.5rem;
    }
    .blog-card.featured .card-excerpt {
      -webkit-line-clamp: 3;
    }
  }

  @media (prefers-reduced-motion: reduce) {
    .blog-card,
    .card-image img,
    .card-title {
      transition: none;
    }
    .blog-card:hover {
      transform: none;
    }
    .blog-card:hover .card-image img {
      transform: none;
    }
  }
</style>
