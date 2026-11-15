<script lang="ts">
  import type { BlogPost } from '$lib/types';
  import BlogCard from '$lib/components/blog/BlogCard.svelte';
  import Seo from '$lib/components/Seo.svelte';
  import { SITE } from '$lib/config';
  import { blogData } from '$data/blog';

  const posts: BlogPost[] = blogData;

  let selectedTag = $state('');

  let allTags = $derived.by(() => {
    const counts = new Map<string, number>();
    for (const p of posts) {
      for (const t of p.tags || []) counts.set(t, (counts.get(t) || 0) + 1);
    }
    return [...counts.entries()].sort();
  });

  let filtered = $derived.by(() => {
    if (!selectedTag) return posts;
    return posts.filter(p => (p.tags || []).includes(selectedTag));
  });
</script>

<svelte:head>
  <title>Blog — Pranav Agarkar</title>
</svelte:head>

<Seo
  title="Blog"
  description="Thoughts, tutorials, and updates on fullstack development with Django, React, Svelte, Go, and modern web technologies."
  url="{SITE.url}/blog"
/>

<main id="main-content" class="blog-page">
  <header class="blog-header">
    <p class="blog-kicker" aria-hidden="true">~/blog</p>
    <h1 class="blog-heading">Blog</h1>
    <p class="blog-subtitle">Thoughts, tutorials, and updates &mdash; {posts.length} {posts.length === 1 ? 'entry' : 'entries'}</p>
  </header>

  {#if allTags.length > 0}
    <div class="tag-filters">
      <button
        class="tag-filter"
        class:active={!selectedTag}
        onclick={() => selectedTag = ''}
        aria-pressed={!selectedTag}
      >All<span class="tag-count">{posts.length}</span></button>
      {#each allTags as [tag, count]}
        <button
          class="tag-filter"
          class:active={selectedTag === tag}
          onclick={() => selectedTag = tag}
          aria-pressed={selectedTag === tag}
        >{tag}<span class="tag-count">{count}</span></button>
      {/each}
    </div>
  {/if}

  {#if filtered.length > 0}
    <div class="blog-grid">
      {#each filtered as post, i (post.slug)}
        <BlogCard {post} index={i} featured={i === 0} />
      {/each}
    </div>
  {:else}
    <div class="blog-empty">
      {#if selectedTag}
        No posts tagged "{selectedTag}".
      {:else}
        No posts yet. Check back soon!
      {/if}
    </div>
  {/if}
</main>

<style>
  .blog-page {
    max-width: 720px;
    margin: 0 auto;
    padding: calc(var(--nav-height) + 3rem) 2rem 4rem;
  }
  .blog-header {
    margin-bottom: 2.5rem;
    border-left: 3px solid var(--accent);
    padding-left: 1rem;
  }
  .blog-kicker {
    font-family: var(--font-heading);
    font-size: 0.7rem;
    letter-spacing: 0.18em;
    text-transform: uppercase;
    color: var(--accent);
    margin: 0 0 0.5rem;
  }
  .blog-heading {
    font-family: var(--font-heading);
    font-size: 2.75rem;
    color: var(--text-primary);
    margin: 0 0 0.5rem;
    line-height: 1.1;
    text-transform: uppercase;
  }
  .blog-subtitle {
    color: var(--text-muted);
    margin: 0;
    font-size: 0.9rem;
    letter-spacing: 0.04em;
  }
  .tag-filters {
    display: flex;
    flex-wrap: wrap;
    gap: 0.4rem;
    margin-bottom: 2rem;
  }
  .tag-filter {
    display: inline-flex;
    align-items: center;
    gap: 0.45rem;
    min-height: 42px;
    padding: 0.35rem 0.85rem;
    border: 1px solid var(--border-color);
    background: transparent;
    color: var(--text-secondary);
    font-size: 0.75rem;
    cursor: pointer;
    font-family: var(--font-heading);
    text-transform: uppercase;
    letter-spacing: 0.08em;
    transition: border-color 0.2s ease, color 0.2s ease, background 0.2s ease;
  }
  .tag-filter:hover {
    border-color: var(--accent);
    color: var(--accent);
    background: rgba(255, 68, 0, 0.06);
  }
  .tag-filter.active {
    background: var(--accent);
    color: #fff;
    border-color: var(--accent);
  }
  .tag-count {
    font-size: 0.6rem;
    opacity: 0.55;
  }
  .tag-filter:focus-visible {
    outline: 2px solid var(--accent);
    outline-offset: 2px;
  }
  .blog-grid {
    display: flex;
    flex-direction: column;
    gap: 1.25rem;
  }
  .blog-empty {
    color: var(--text-muted);
    text-align: center;
    padding: 3rem;
    border: 1px dashed var(--border-color);
  }

  @media (max-width: 600px) {
    .blog-page {
      padding: calc(var(--nav-height) + 2rem) 1rem 3rem;
    }
    .blog-heading {
      font-size: 2.1rem;
    }
  }
</style>
