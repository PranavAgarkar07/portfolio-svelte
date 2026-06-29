<script lang="ts">
  import { onMount } from 'svelte';
  import type { BlogPost } from '$lib/types';
  import BlogCard from '$lib/components/blog/BlogCard.svelte';
  import Seo from '$lib/components/Seo.svelte';
  import { staleWhileRevalidate } from '$lib/utils/cache';

  const BASE = (import.meta.env.VITE_API_URL || '').replace(/\/$/, '');

  let posts = $state<BlogPost[]>([]);
  let loading = $state(true);
  let selectedTag = $state('');

  let allTags = $derived(() => {
    const tags = new Set<string>();
    for (const p of posts) {
      for (const t of p.tags || []) tags.add(t);
    }
    return [...tags].sort();
  });

  let filtered = $derived(() => {
    if (!selectedTag) return posts;
    return posts.filter(p => (p.tags || []).includes(selectedTag));
  });

  onMount(async () => {
    if (!BASE) { loading = false; return; }
    const ok = await staleWhileRevalidate<BlogPost[]>(
      'blog:list',
      `${BASE}/api/blog`,
      5 * 60 * 1000,
      (data) => { posts = data; loading = false; },
      (raw) => (raw as { posts: BlogPost[] }).posts ?? [],
    );
    if (!ok) loading = false;
  });
</script>

<svelte:head>
  <title>Blog — Pranav Agarkar</title>
</svelte:head>

<Seo
  title="Blog"
  description="Thoughts, tutorials, and updates on fullstack development with Django, React, Svelte, Go, and modern web technologies."
  url="{import.meta.env.SITE_URL || 'https://pranavagarkar07.github.io/portfolio-svelte'}/blog"
/>

<div class="blog-page">
  <header class="blog-header">
    <h1 class="blog-heading">Blog</h1>
    <p class="blog-subtitle">Thoughts, tutorials, and updates</p>
  </header>

  {#if allTags().length > 0}
    <div class="tag-filters">
      <button
        class="tag-filter"
        class:active={!selectedTag}
        onclick={() => selectedTag = ''}
      >All</button>
      {#each allTags() as tag}
        <button
          class="tag-filter"
          class:active={selectedTag === tag}
          onclick={() => selectedTag = tag}
        >{tag}</button>
      {/each}
    </div>
  {/if}

  {#if loading}
    <div class="blog-loading">Loading...</div>
  {:else if filtered().length > 0}
    <div class="blog-grid">
      {#each filtered() as post}
        <BlogCard {post} />
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
</div>

<style>
  .blog-page {
    max-width: 680px;
    margin: 0 auto;
    padding: calc(var(--nav-height) + 3rem) 2rem 4rem;
  }
  .blog-header {
    margin-bottom: 2.5rem;
    border-left: 3px solid var(--accent);
    padding-left: 1rem;
  }
  .blog-heading {
    font-family: var(--font-heading);
    font-size: 2.5rem;
    color: var(--text-primary);
    margin: 0 0 0.35rem;
    line-height: 1.15;
  }
  .blog-subtitle {
    color: var(--text-muted);
    margin: 0;
    font-size: 0.95rem;
  }
  .tag-filters {
    display: flex;
    flex-wrap: wrap;
    gap: 0.4rem;
    margin-bottom: 2rem;
  }
  .tag-filter {
    padding: 0.35rem 0.8rem;
    border: 1px solid var(--border-color);
    background: transparent;
    color: var(--text-secondary);
    font-size: 0.75rem;
    cursor: pointer;
    font-family: inherit;
    transition: all 0.15s;
    letter-spacing: 0.03em;
  }
  .tag-filter:hover {
    border-color: var(--accent);
    color: var(--accent);
    background: rgba(255,68,0,0.06);
  }
  .tag-filter.active {
    background: var(--accent);
    color: #fff;
    border-color: var(--accent);
  }
  .blog-grid {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }
  .blog-loading, .blog-empty {
    color: var(--text-muted);
    text-align: center;
    padding: 3rem;
  }
</style>
