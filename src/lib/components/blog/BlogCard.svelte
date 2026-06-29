<script lang="ts">
  import type { BlogPost } from '$lib/types';

  let { post }: { post: BlogPost } = $props();

  function readingTime(md: string): string {
    const words = md.split(/\s+/).length;
    const min = Math.max(1, Math.round(words / 200));
    return `${min} min read`;
  }
</script>

<a href="/blog/{post.slug}" class="blog-card">
  {#if post.cover_image}
    <div class="card-image">
      <img src={post.cover_image} alt={post.title} loading="lazy" />
    </div>
  {/if}
  <div class="card-body">
    {#if post.tags?.length}
      <div class="card-tags">
        {#each post.tags.slice(0, 3) as tag}
          <span class="tag">{tag}</span>
        {/each}
      </div>
    {/if}
    <h3 class="card-title">{post.title}</h3>
    {#if post.excerpt}
      <p class="card-excerpt">{post.excerpt}</p>
    {/if}
    <div class="card-meta">
      <span class="card-meta-name">{post.author_name}</span>
      <span class="card-meta-divider">·</span>
      <span class="card-meta-item">{post.published_at || post.created_at}</span>
      <span class="card-meta-divider">·</span>
      <span class="card-meta-item">{readingTime(post.content_md)}</span>
    </div>
  </div>
</a>

<style>
  .blog-card {
    display: block;
    text-decoration: none;
    position: relative;
    contain: paint;
    will-change: clip-path;

    background: var(--card-bg);
    border: 1px solid var(--border-color);
    border-left: 3px solid var(--border-color);
    box-shadow: var(--card-shadow);
    clip-path: polygon(12px 0, 100% 0, 100% calc(100% - 12px), calc(100% - 12px) 100%, 0 100%, 0 12px);

    transition: transform 0.2s ease, box-shadow 0.2s ease;
  }
  .blog-card:hover {
    border-left-color: var(--accent);
    transform: translateY(-3px);
    box-shadow: var(--card-shadow-hover);
  }
  .card-image {
    width: 100%;
    aspect-ratio: 2 / 1;
    overflow: hidden;
    background: var(--surface-dark);
  }
  .card-image img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    transition: transform 0.4s ease;
  }
  .blog-card:hover .card-image img {
    transform: scale(1.05);
  }
  .card-body {
    padding: 1.25rem;
  }
  .card-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 0.4rem;
    margin-bottom: 0.6rem;
  }
  .tag {
    font-size: 0.65rem;
    padding: 0.2rem 0.5rem;
    background: rgba(255,68,0,0.1);
    color: var(--accent);
    letter-spacing: 0.04em;
    text-transform: uppercase;
    font-weight: 600;
  }
  .card-title {
    font-family: var(--font-heading);
    font-size: 1.15rem;
    color: var(--text-primary);
    margin: 0 0 0.5rem;
    line-height: 1.35;
  }
  .card-excerpt {
    color: var(--text-secondary);
    font-size: 0.85rem;
    line-height: 1.55;
    margin: 0 0 0.75rem;
  }
  .card-meta {
    display: flex;
    flex-wrap: wrap;
    gap: 0.3rem;
    font-size: 0.75rem;
    color: var(--text-muted);
  }
  .card-meta-name {
    color: var(--accent);
    font-weight: 500;
  }
  .card-meta-divider {
    opacity: 0.4;
  }
</style>
