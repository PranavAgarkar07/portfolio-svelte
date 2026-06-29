<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { base } from '$app/paths';
  import { user, isLoggedIn } from '$lib/stores/auth';
  import type { BlogPost } from '$lib/types';
  import { renderMarkdown, readingTime, extractHeadings } from '$lib/utils/markdown';
  import CommentSection from '$lib/components/blog/CommentSection.svelte';
  import ImageGallery from '$lib/components/blog/ImageGallery.svelte';
  import Seo from '$lib/components/Seo.svelte';
  import { staleWhileRevalidate } from '$lib/utils/cache';
  import { SITE } from '$lib/config';

  const BASE = (import.meta.env.VITE_API_URL || '').replace(/\/$/, '');

  let post = $state<BlogPost | null>(null);
  let loading = $state(true);
  let rendered = $state('');
  let headings = $state<Array<{ level: number; text: string; id: string }>>([]);
  let activeHeading = $state('');
  let readProgress = $state(0);

  onMount(() => {
    const slug = $page.params.slug;
    if (!BASE || !slug) { loading = false; return; }

    function applyPost(p: BlogPost) {
      post = p;
      if (p?.content_md) {
        rendered = renderMarkdown(p.content_md);
        headings = extractHeadings(p.content_md);
      }
    }

    staleWhileRevalidate<BlogPost>(
      `blog:post:${slug}`,
      `${BASE}/api/blog/${encodeURIComponent(slug)}`,
      5 * 60 * 1000,
      (data) => { applyPost(data); loading = false; },
      (raw) => (raw as { post: BlogPost }).post ?? null,
    ).then(ok => { if (!ok && !post) loading = false; });

    const onScroll = () => {
      const scrollTop = window.scrollY;
      const docHeight = document.documentElement.scrollHeight - window.innerHeight;
      readProgress = docHeight > 0 ? Math.min(scrollTop / docHeight * 100, 100) : 0;

      for (const h of headings) {
        const el = document.getElementById(h.id);
        if (el) {
          const rect = el.getBoundingClientRect();
          if (rect.top <= 120) {
            activeHeading = h.id;
          }
        }
      }
    };
    window.addEventListener('scroll', onScroll);
    return () => window.removeEventListener('scroll', onScroll);
  });

  function scrollTo(id: string) {
    const el = document.getElementById(id);
    if (el) el.scrollIntoView({ behavior: 'smooth' });
  }

  function shareUrl() {
    return `${SITE.url}/blog/${post?.slug}`;
  }
</script>

<svelte:head>
  <title>{post?.title || 'Blog Post'} — Pranav Agarkar</title>
</svelte:head>

{#if post}
  <Seo
    title={post.title}
    description={post.excerpt}
    image={post.cover_image || SITE.image}
    url={`${SITE.url}/blog/${post.slug}`}
    type="article"
    publishedTime={post.published_at || post.created_at}
    modifiedTime={post.updated_at}
    tags={post.tags || []}
    jsonld={{
      "@context": "https://schema.org",
      "@type": "BlogPosting",
      "headline": post.title,
      "description": post.excerpt,
      "image": post.cover_image || SITE.image,
      "url": `${SITE.url}/blog/${post.slug}`,
      "datePublished": post.published_at || post.created_at,
      "dateModified": post.updated_at,
      "author": {
        "@type": "Person",
        "name": post.author_name || SITE.author,
        "url": SITE.url,
      },
      "publisher": {
        "@type": "Person",
        "name": SITE.author,
        "url": SITE.url,
      },
      "mainEntityOfPage": {
        "@type": "WebPage",
        "@id": `${SITE.url}/blog/${post.slug}`,
      },
      "keywords": (post.tags || []).join(", "),
      "wordCount": post.content_md ? post.content_md.split(/\s+/).length : undefined,
    }}
  />
{/if}

<div class="blog-post-page">
  {#if loading}
    <div class="post-loading">Loading...</div>
  {:else if post}
    <div class="reading-progress" style="width: {readProgress}%"></div>

    {#if post.cover_image}
      <div class="hero-image">
        <img src={post.cover_image} alt={post.title} loading="eager" />
      </div>
    {/if}

      {#if post.images?.length}
        <ImageGallery images={post.images} />
      {/if}

      <article class="post-article">
      <header class="post-header">
        {#if post.tags?.length}
          <div class="post-tags">
            {#each post.tags as tag}
              <span class="tag">{tag}</span>
            {/each}
          </div>
        {/if}
        <h1 class="post-title">{post.title}</h1>
        <div class="post-meta">
          <span class="meta-author">{post.author_name}</span>
          <span class="meta-sep">·</span>
          <time>{post.published_at || post.created_at}</time>
          <span class="meta-sep">·</span>
          <span>{readingTime(post.content_md)}</span>
          {#if $isLoggedIn && $user?.id === post.author_id}
            <span class="meta-sep">·</span>
            <a href={`${base}/admin/blog/${post.id}`} class="meta-edit">Edit</a>
          {/if}
        </div>
      </header>

      {#if headings.length > 1}
        <nav class="toc">
          <strong class="toc-title">On this page</strong>
          {#each headings as h}
            <button
              class="toc-item"
              class:toc-h2={h.level === 2}
              class:toc-h3={h.level === 3}
              class:active={activeHeading === h.id}
              onclick={() => scrollTo(h.id)}
            >{h.text}</button>
          {/each}
        </nav>
      {/if}

      <div class="post-content">{@html rendered}</div>

      <footer class="post-footer">
        <div class="author-bio">
          <div class="author-avatar">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
              <circle cx="12" cy="7" r="4"/>
            </svg>
          </div>
          <div>
            <strong>{post.author_name}</strong>
            <p>Fullstack Developer &mdash; writing about web dev, Svelte, Go, and building things.</p>
          </div>
        </div>

        <div class="share-links">
          <span>Share:</span>
          <a href="https://twitter.com/intent/tweet?text={encodeURIComponent(post.title)}&url={encodeURIComponent(shareUrl())}" target="_blank" rel="noopener" aria-label="Share on X">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor"><path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z"/></svg>
          </a>
          <a href="https://www.linkedin.com/sharing/share-offsite/?url={encodeURIComponent(shareUrl())}" target="_blank" rel="noopener" aria-label="Share on LinkedIn">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor"><path d="M20.447 20.452h-3.554v-5.569c0-1.328-.027-3.037-1.852-3.037-1.853 0-2.136 1.445-2.136 2.939v5.667H9.351V9h3.414v1.561h.046c.477-.9 1.637-1.85 3.37-1.85 3.601 0 4.267 2.37 4.267 5.455v6.286zM5.337 7.433c-1.144 0-2.063-.926-2.063-2.065 0-1.138.92-2.063 2.063-2.063 1.14 0 2.064.925 2.064 2.063 0 1.139-.925 2.065-2.064 2.065zm1.782 13.019H3.555V9h3.564v11.452zM22.225 0H1.771C.792 0 0 .774 0 1.729v20.542C0 23.227.792 24 1.771 24h20.451C23.2 24 24 23.227 24 22.271V1.729C24 .774 23.2 0 22.222 0h.003z"/></svg>
          </a>
        </div>
      </footer>
    </article>

    <CommentSection slug={post.slug} />

    <a href="/blog" class="back-link">&larr; Back to Blog</a>
  {:else}
    <div class="post-error">Post not found.</div>
  {/if}
</div>

<style>
  .blog-post-page {
    max-width: 720px;
    margin: 0 auto;
    padding: calc(var(--nav-height) + 3rem) 2rem 4rem;
    position: relative;
  }

  .post-article {
    background: var(--card-bg);
    padding: 2.5rem;
    border: 1px solid var(--border-color);
    border-left: 3px solid var(--border-color);
    box-shadow: var(--card-shadow);
    clip-path: polygon(12px 0, 100% 0, 100% calc(100% - 12px), calc(100% - 12px) 100%, 0 100%, 0 12px);
    contain: paint;
    will-change: clip-path;
    transition: border-left-color 0.2s ease;
  }
  .post-article:hover {
    border-left-color: var(--accent);
  }

  .reading-progress {
    position: fixed;
    top: 0;
    left: 0;
    height: 3px;
    background: linear-gradient(90deg, var(--accent), #ff7700);
    z-index: 100;
    transition: width 0.05s linear;
  }

  .hero-image {
    width: 100%;
    aspect-ratio: 2 / 1;
    margin-bottom: 2rem;
    background: var(--surface-dark);
    border: 1px solid var(--border-color);
    clip-path: polygon(16px 0, 100% 0, 100% calc(100% - 16px), calc(100% - 16px) 100%, 0 100%, 0 16px);
  }
  .hero-image img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .post-header {
    margin-bottom: 2rem;
  }

  .post-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 0.4rem;
    margin-bottom: 0.75rem;
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

  .post-title {
    font-family: var(--font-heading);
    font-size: 2.25rem;
    color: var(--text-primary);
    margin: 0 0 0.75rem;
    line-height: 1.2;
  }

  .post-meta {
    display: flex;
    flex-wrap: wrap;
    gap: 0.3rem;
    font-size: 0.85rem;
    color: var(--text-muted);
  }
  .meta-author {
    color: var(--accent);
    font-weight: 500;
  }
  .meta-sep {
    opacity: 0.4;
  }
  .meta-edit {
    color: var(--accent);
    text-decoration: none;
    font-size: 0.8rem;
    font-weight: 500;
    padding: 0.1rem 0.35rem;
    border: 1px solid var(--accent);
    transition: background 0.15s;
  }
  .meta-edit:hover {
    background: rgba(255,68,0,0.1);
  }

  .toc {
    background: var(--toc-bg);
    border: 1px solid var(--border-color);
    border-left: 2px solid var(--accent);
    padding: 1rem 1.25rem;
    margin-bottom: 2rem;
    clip-path: polygon(8px 0, 100% 0, 100% calc(100% - 8px), calc(100% - 8px) 100%, 0 100%, 0 8px);
  }
  .toc-title {
    display: block;
    font-size: 0.65rem;
    text-transform: uppercase;
    letter-spacing: 0.1em;
    color: var(--accent);
    font-weight: 600;
    margin-bottom: 0.5rem;
  }
  .toc-item {
    display: block;
    background: none;
    border: none;
    color: var(--text-secondary);
    font-size: 0.85rem;
    padding: 0.2rem 0;
    cursor: pointer;
    text-align: left;
    font-family: inherit;
    width: 100%;
    transition: color 0.15s;
  }
  .toc-item:hover {
    color: var(--text-primary);
  }
  .toc-item.active {
    color: var(--accent);
  }
  .toc-h3 {
    padding-left: 1rem;
    font-size: 0.8rem;
  }

  .post-content {
    color: var(--text-primary);
    font-size: 1.05rem;
    line-height: 1.85;
  }
  .post-content :global(h2) {
    font-family: var(--font-heading);
    font-size: 1.6rem;
    color: var(--text-primary);
    margin: 3rem 0 1rem;
    scroll-margin-top: 5rem;
    padding-bottom: 0.4rem;
    border-bottom: 1px solid var(--border-color);
  }
  .post-content :global(h3) {
    font-family: var(--font-heading);
    font-size: 1.25rem;
    color: var(--text-primary);
    margin: 2.25rem 0 0.75rem;
    scroll-margin-top: 5rem;
  }
  .post-content :global(h4) {
    font-family: var(--font-heading);
    font-size: 1.1rem;
    color: var(--text-primary);
    margin: 1.75rem 0 0.5rem;
    scroll-margin-top: 5rem;
  }
  .post-content :global(p) {
    margin: 0 0 1.25rem;
  }
  .post-content :global(code) {
    background: var(--inline-code-bg);
    padding: 0.2rem 0.45rem;
    font-size: 0.85em;
    font-family: 'JetBrains Mono', 'Fira Code', monospace;
    color: var(--code-color);
  }
  .post-content :global(pre) {
    background: var(--code-bg);
    border: 1px solid var(--border-color);
    padding: 1.25rem;
    overflow-x: auto;
    margin: 1.5rem 0;
    clip-path: polygon(8px 0, 100% 0, 100% calc(100% - 8px), calc(100% - 8px) 100%, 0 100%, 0 8px);
  }
  .post-content :global(pre code) {
    background: none;
    padding: 0;
    font-size: 0.85rem;
    line-height: 1.65;
  }
  .post-content :global(img) {
    max-width: 100%;
    border: 1px solid var(--border-color);
    margin: 2rem auto;
    display: block;
  }
  .post-content :global(a) {
    color: var(--accent);
    text-decoration: underline;
    text-underline-offset: 2px;
    text-decoration-thickness: 1px;
    transition: opacity 0.15s;
  }
  .post-content :global(a:hover) {
    opacity: 0.8;
  }
  .post-content :global(ul) {
    padding-left: 1.5rem;
    margin: 0.5rem 0 1.25rem;
  }
  .post-content :global(ol) {
    padding-left: 1.5rem;
    margin: 0.5rem 0 1.25rem;
  }
  .post-content :global(li) {
    margin: 0.35rem 0;
  }
  .post-content :global(blockquote) {
    border-left: 3px solid var(--accent);
    padding: 0.75rem 1.25rem;
    margin: 1.5rem 0;
    background: rgba(255,68,0,0.05);
    color: var(--text-secondary);
    font-style: italic;
  }
  .post-content :global(blockquote p) {
    margin: 0;
  }
  .post-content :global(hr) {
    border: none;
    border-top: 1px solid var(--border-color);
    margin: 2.5rem 0;
  }
  .post-content :global(table) {
    width: 100%;
    border-collapse: collapse;
    margin: 1.5rem 0;
    font-size: 0.9rem;
    border: 1px solid var(--border-color);
  }
  .post-content :global(th) {
    text-align: left;
    padding: 0.6rem 0.75rem;
    background: rgba(255,255,255,0.04);
    border-bottom: 2px solid var(--border-color);
    color: var(--text-primary);
    font-weight: 600;
    font-size: 0.85rem;
  }
  .post-content :global(td) {
    padding: 0.6rem 0.75rem;
    border-bottom: 1px solid var(--border-color);
  }
  .post-content :global(tr:last-child td) {
    border-bottom: none;
  }
  .post-content :global(tr:hover) {
    background: var(--toc-bg);
  }
  .post-content :global(del) {
    opacity: 0.5;
  }

  .post-footer {
    margin-top: 3rem;
    padding-top: 1.5rem;
    border-top: 1px solid var(--border-color);
  }

  .author-bio {
    display: flex;
    gap: 0.75rem;
    align-items: flex-start;
    margin-bottom: 1.5rem;
    padding: 1rem;
    background: var(--author-bio-bg);
    border: 1px solid var(--border-color);
    clip-path: polygon(6px 0, 100% 0, 100% calc(100% - 6px), calc(100% - 6px) 100%, 0 100%, 0 6px);
  }
  .author-avatar {
    width: 40px;
    height: 40px;
    background: rgba(255,68,0,0.1);
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
    color: var(--accent);
  }
  .author-bio strong {
    color: var(--text-primary);
    font-size: 0.9rem;
  }
  .author-bio p {
    color: var(--text-secondary);
    font-size: 0.8rem;
    margin: 0.2rem 0 0;
    line-height: 1.5;
  }

  .share-links {
    display: flex;
    gap: 0.5rem;
    align-items: center;
    font-size: 0.75rem;
    color: var(--text-muted);
    text-transform: uppercase;
    letter-spacing: 0.06em;
  }
  .share-links a {
    color: var(--text-secondary);
    transition: color 0.15s;
    display: flex;
    padding: 0.4rem;
    border: 1px solid var(--border-color);
  }
  .share-links a:hover {
    color: var(--accent);
    border-color: var(--accent);
    background: rgba(255,68,0,0.06);
  }
  .share-links span {
    margin-right: 0.25rem;
  }

  .back-link {
    display: inline-flex;
    align-items: center;
    gap: 0.3rem;
    margin-top: 2rem;
    color: var(--text-secondary);
    text-decoration: none;
    font-size: 0.85rem;
    padding: 0.5rem 1rem;
    border: 1px solid var(--border-color);
    transition: all 0.15s;
    width: fit-content;
    clip-path: polygon(6px 0, 100% 0, 100% calc(100% - 6px), calc(100% - 6px) 100%, 0 100%, 0 6px);
  }
  .back-link:hover {
    color: var(--accent);
    border-color: var(--accent);
    background: rgba(255,68,0,0.06);
  }

  @media (max-width: 600px) {
    .blog-post-page {
      padding: calc(var(--nav-height) + 2rem) 1rem 3rem;
    }
    .post-article {
      padding: 1.25rem;
    }
    .post-title {
      font-size: 1.5rem;
    }
    .post-content {
      font-size: 0.95rem;
    }
  }

  @media (max-width: 400px) {
    .post-article {
      padding: 1rem;
    }
    .post-title {
      font-size: 1.3rem;
    }
  }

  .post-loading, .post-error {
    color: var(--text-muted);
    text-align: center;
    padding: 3rem;
  }
</style>
