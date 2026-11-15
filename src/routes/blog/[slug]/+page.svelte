<script lang="ts">
  import { page } from '$app/stores';
  import { base } from '$app/paths';
  import type { BlogPost } from '$lib/types';
  import { renderMarkdown, readingTime, extractHeadings } from '$lib/utils/markdown';
  import { formatDate } from '$lib/utils/format';
  import ImageGallery from '$lib/components/blog/ImageGallery.svelte';
  import Seo from '$lib/components/Seo.svelte';
  import { blogData } from '$data/blog';
  import { SITE } from '$lib/config';

  const allPosts: BlogPost[] = blogData;

  const slug = $page.params.slug;
  let post = $state<BlogPost | null>(null);
  let rendered = $state('');
  let headings = $state<Array<{ level: number; text: string; id: string }>>([]);
  let activeHeading = $state('');
  let readProgress = $state(0);

  function escapeJsonLd(str: string): string {
    return str.replace(/<\/script>/g, '<\\/script>');
  }

  // Adjacent posts for the bottom navigation (newest → oldest ordering)
  const sortedPosts = [...allPosts].sort((a, b) =>
    (b.published_at || b.created_at).localeCompare(a.published_at || a.created_at)
  );
  const postIndex = slug ? sortedPosts.findIndex(p => p.slug === slug) : -1;
  const newerPost = postIndex > 0 ? sortedPosts[postIndex - 1] : null;
  const olderPost = postIndex >= 0 && postIndex < sortedPosts.length - 1 ? sortedPosts[postIndex + 1] : null;

  if (slug) {
    const found = allPosts.find((p: BlogPost) => p.slug === slug);
    if (found) {
      post = found;
      if (found.content_md) {
        rendered = renderMarkdown(found.content_md);
        headings = extractHeadings(found.content_md);
      }
    }
  }

  function scrollTo(id: string) {
    const el = document.getElementById(id);
    if (!el) return;
    const reduce = window.matchMedia('(prefers-reduced-motion: reduce)').matches;
    el.scrollIntoView({ behavior: reduce ? 'auto' : 'smooth', block: 'start' });
  }

  function shareUrl() {
    return `${SITE.url}/blog/${post?.slug}`;
  }

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

  // Inject a copy button into every code block after the content renders
  $effect(() => {
    if (!rendered) return;
    const pres = document.querySelectorAll<HTMLElement>('.post-content pre');
    for (const pre of pres) {
      if (pre.querySelector('.code-copy')) continue;
      const btn = document.createElement('button');
      btn.type = 'button';
      btn.className = 'code-copy';
      btn.textContent = 'COPY';
      btn.setAttribute('aria-label', 'Copy code block');
      btn.addEventListener('click', () => {
        const code = pre.querySelector('code')?.textContent ?? '';
        const flash = (label: string) => {
          btn.textContent = label;
          setTimeout(() => (btn.textContent = 'COPY'), 1400);
        };
        if (navigator.clipboard) {
          navigator.clipboard.writeText(code).then(
            () => flash('COPIED'),
            () => flash('ERROR')
          );
        } else {
          flash('ERROR');
        }
      });
      pre.appendChild(btn);
    }
  });
</script>

<svelte:head>
  <title>{post?.title || 'Blog Post'} — Pranav Agarkar</title>
</svelte:head>

<svelte:window onscroll={onScroll} />

{#if post}
  <Seo
    title={post.title}
    description={post.excerpt}
    image={post.cover_image ? `${SITE.url}${post.cover_image}` : SITE.image}
    url={`${SITE.url}/blog/${post.slug}`}
    type="article"
    publishedTime={post.published_at || post.created_at}
    modifiedTime={post.updated_at}
    tags={post.tags || []}
    jsonld={{
      "@context": "https://schema.org",
      "@type": "BlogPosting",
      "headline": escapeJsonLd(post.title),
      "description": escapeJsonLd(post.excerpt),
      "image": escapeJsonLd(post.cover_image || SITE.image),
      "url": escapeJsonLd(`${SITE.url}/blog/${post.slug}`),
      "datePublished": post.published_at || post.created_at,
      "dateModified": post.updated_at,
      "author": {
        "@type": "Person",
        "name": escapeJsonLd(post.author_name || SITE.author),
        "url": SITE.url,
      },
      "publisher": {
        "@type": "Person",
        "name": escapeJsonLd(SITE.author),
        "url": SITE.url,
      },
      "mainEntityOfPage": {
        "@type": "WebPage",
        "@id": escapeJsonLd(`${SITE.url}/blog/${post.slug}`),
      },
      "keywords": escapeJsonLd((post.tags || []).join(", ")),
      "wordCount": post.content_md ? post.content_md.split(/\s+/).length : undefined,
    }}
  />
{/if}

<main id="main-content" class="blog-post-page">
  {#if post}
    <div class="reading-progress" style="width: {readProgress}%"></div>

    <a href="{base}/blog" class="crumb">&larr; All posts</a>

    <div class="post-shell">
      <article class="post-article">
        <header class="post-header">
          {#if post.tags?.length}
            <div class="post-tags">
              {#each post.tags as tag (tag)}
                <span class="tag">{tag}</span>
              {/each}
            </div>
          {/if}
          <h1 class="post-title">{post.title}</h1>
          <div class="post-meta">
            <span class="meta-author">{post.author_name}</span>
            <span class="meta-sep" aria-hidden="true">&middot;</span>
            <time>{formatDate(post.published_at || post.created_at)}</time>
            <span class="meta-sep" aria-hidden="true">&middot;</span>
            <span>{readingTime(post.content_md)}</span>
          </div>
        </header>

        {#if post.cover_image}
          <figure class="hero-image">
            <img
              src={`${base}${post.cover_image}`}
              alt={post.title}
              loading="lazy"
              width="800"
              height="400"
              decoding="async"
            />
          </figure>
        {/if}

        {#if post.images?.length}
          <ImageGallery images={post.images} />
        {/if}

        {#if headings.length > 1}
          <nav class="toc toc-flow" aria-label="Table of contents">
            <strong class="toc-title">On this page</strong>
            {#each headings as h (h.id)}
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
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
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
              <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true"><path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z"/></svg>
            </a>
            <a href="https://www.linkedin.com/sharing/share-offsite/?url={encodeURIComponent(shareUrl())}" target="_blank" rel="noopener" aria-label="Share on LinkedIn">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true"><path d="M20.447 20.452h-3.554v-5.569c0-1.328-.027-3.037-1.852-3.037-1.853 0-2.136 1.445-2.136 2.939v5.667H9.351V9h3.414v1.561h.046c.477-.9 1.637-1.85 3.37-1.85 3.601 0 4.267 2.37 4.267 5.455v6.286zM5.337 7.433c-1.144 0-2.063-.926-2.063-2.065 0-1.138.92-2.063 2.063-2.063 1.14 0 2.064.925 2.064 2.063 0 1.139-.925 2.065-2.064 2.065zm1.782 13.019H3.555V9h3.564v11.452zM22.225 0H1.771C.792 0 0 .774 0 1.729v20.542C0 23.227.792 24 1.771 24h20.451C23.2 24 24 23.227 24 22.271V1.729C24 .774 23.2 0 22.222 0h.003z"/></svg>
            </a>
          </div>
        </footer>
      </article>

      {#if headings.length > 1}
        <aside class="toc toc-rail" aria-label="On this page">
          <strong class="toc-title">On this page</strong>
          {#each headings as h (h.id)}
            <button
              class="toc-item"
              class:toc-h2={h.level === 2}
              class:toc-h3={h.level === 3}
              class:active={activeHeading === h.id}
              onclick={() => scrollTo(h.id)}
            >{h.text}</button>
          {/each}
        </aside>
      {/if}
    </div>

    {#if newerPost || olderPost}
      <nav class="post-nav" aria-label="Adjacent posts">
        {#if newerPost}
          <a class="pn-item" href="{base}/blog/{newerPost.slug}">
            <span class="pn-label">&larr; Newer</span>
            <strong class="pn-title">{newerPost.title}</strong>
          </a>
        {:else}
          <span class="pn-spacer" aria-hidden="true"></span>
        {/if}
        {#if olderPost}
          <a class="pn-item pn-next" href="{base}/blog/{olderPost.slug}">
            <span class="pn-label">Older &rarr;</span>
            <strong class="pn-title">{olderPost.title}</strong>
          </a>
        {/if}
      </nav>
    {/if}
  {:else}
    <div class="post-error" role="alert">Post not found.</div>
  {/if}
</main>

<style>
  .blog-post-page {
    max-width: 720px;
    margin: 0 auto;
    padding: calc(var(--nav-height) + 2.5rem) 2rem 4rem;
    position: relative;
  }

  .crumb {
    display: inline-flex;
    align-items: center;
    gap: 0.4rem;
    font-family: var(--font-heading);
    font-size: 0.7rem;
    letter-spacing: 0.14em;
    text-transform: uppercase;
    color: var(--text-secondary);
    text-decoration: none;
    margin-bottom: 1.5rem;
    transition: color 0.2s ease;
  }
  .crumb:hover {
    color: var(--accent);
  }

  .reading-progress {
    position: fixed;
    top: 0;
    left: 0;
    height: 3px;
    background: var(--accent);
    z-index: 100;
    transition: width 0.05s linear;
  }

  .post-article {
    background: var(--card-bg);
    padding: 2.5rem;
    border: 1px solid var(--border-color);
    border-left: 3px solid var(--border-color);
    box-shadow: var(--card-shadow);
  }

  .post-header {
    margin-bottom: 2rem;
  }

  .post-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 0.4rem;
    margin-bottom: 0.85rem;
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

  .post-title {
    font-family: var(--font-heading);
    font-size: 2.5rem;
    color: var(--text-primary);
    margin: 0 0 0.9rem;
    line-height: 1.15;
    text-transform: uppercase;
  }

  .post-meta {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 0.45rem;
    font-family: var(--font-heading);
    font-size: 0.75rem;
    letter-spacing: 0.06em;
    text-transform: uppercase;
    color: var(--text-muted);
  }
  .meta-author {
    color: var(--accent);
  }
  .meta-sep {
    opacity: 0.4;
  }

  .hero-image {
    margin: 0 0 2rem;
    width: 100%;
    aspect-ratio: 2 / 1;
    background: var(--surface-dark);
    border: 1px solid var(--border-color);
    overflow: hidden;
  }
  .hero-image img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  /* ── Table of contents: in-flow on narrow screens ── */
  .toc {
    background: var(--toc-bg);
    border: 1px solid var(--border-color);
    border-left: 2px solid var(--accent);
    padding: 1rem 1.25rem;
  }
  .toc-flow {
    margin-bottom: 2rem;
  }
  .toc-title {
    display: block;
    font-family: var(--font-heading);
    font-size: 0.65rem;
    text-transform: uppercase;
    letter-spacing: 0.14em;
    color: var(--accent);
    margin-bottom: 0.5rem;
  }
  .toc-item {
    display: block;
    background: none;
    border: none;
    color: var(--text-secondary);
    font-size: 0.85rem;
    padding: 0.25rem 0;
    cursor: pointer;
    text-align: left;
    font-family: inherit;
    width: 100%;
    transition: color 0.2s ease;
  }
  .toc-item:hover {
    color: var(--text-primary);
  }
  .toc-item.active {
    color: var(--accent);
  }
  .toc-item:focus-visible {
    outline: 2px solid var(--accent);
    outline-offset: 1px;
  }
  .toc-h3 {
    padding-left: 1rem;
    font-size: 0.8rem;
  }
  .toc-rail {
    display: none;
    border-left: none;
    border-top: 2px solid var(--accent);
  }

  .post-shell {
    display: block;
  }
  /* ── Wide screens: article + sticky TOC rail (dev.to / docs pattern) ── */
  @media (min-width: 1140px) {
    .blog-post-page {
      max-width: 1080px;
    }
    .post-shell {
      display: grid;
      grid-template-columns: minmax(0, 720px) 240px;
      gap: 3rem;
      align-items: start;
    }
    .toc-flow {
      display: none;
    }
    .toc-rail {
      display: block;
      position: sticky;
      top: calc(var(--nav-height) + 2rem);
      max-height: calc(100vh - var(--nav-height) - 4rem);
      overflow-y: auto;
    }
  }

  /* ── Prose: 68ch measure, generous leading, rhythm per heading ── */
  .post-content {
    color: var(--text-primary);
    font-size: 1.1rem;
    line-height: 1.8;
    max-width: 68ch;
  }
  .post-content :global(h2) {
    font-family: var(--font-heading);
    font-size: 1.6rem;
    color: var(--text-primary);
    margin: 2.6em 0 0.6em;
    scroll-margin-top: 5rem;
    padding-bottom: 0.4rem;
    border-bottom: 1px solid var(--border-color);
  }
  .post-content :global(h3) {
    font-family: var(--font-heading);
    font-size: 1.25rem;
    color: var(--text-primary);
    margin: 2em 0 0.5em;
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
    padding: 0.15rem 0.4rem;
    font-size: 0.85em;
    font-family: 'JetBrains Mono', 'Fira Code', monospace;
    color: var(--code-color);
  }
  .post-content :global(pre) {
    position: relative;
    background: var(--code-bg);
    border: 1px solid var(--border-color);
    padding: 1.25rem;
    overflow-x: auto;
    margin: 1.5rem 0;
    /* code blocks run wider than the 68ch prose column */
    max-width: none;
  }
  .post-content :global(pre code) {
    background: none;
    padding: 0;
    font-size: 0.85rem;
    line-height: 1.65;
  }
  .post-content :global(.code-copy) {
    position: absolute;
    top: 0.75rem;
    right: 0.75rem;
    font-family: var(--font-heading);
    font-size: 0.58rem;
    font-weight: 700;
    letter-spacing: 0.14em;
    padding: 0.35rem 0.6rem;
    color: var(--text-secondary);
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid var(--border-color);
    cursor: pointer;
    transition: color 0.2s ease, border-color 0.2s ease;
  }
  .post-content :global(.code-copy:hover) {
    color: var(--accent);
    border-color: var(--accent);
  }
  .post-content :global(.code-copy:focus-visible) {
    outline: 2px solid var(--accent);
    outline-offset: 1px;
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
    transition: opacity 0.2s ease;
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
    background: rgba(255, 68, 0, 0.05);
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
  /* Tables scroll horizontally instead of squashing on narrow screens */
  .post-content :global(table) {
    width: 100%;
    border-collapse: collapse;
    margin: 1.5rem 0;
    font-size: 0.9rem;
    border: 1px solid var(--border-color);
    display: block;
    overflow-x: auto;
  }
  .post-content :global(th) {
    text-align: left;
    padding: 0.6rem 0.75rem;
    background: rgba(255, 255, 255, 0.04);
    border-bottom: 2px solid var(--border-color);
    color: var(--text-primary);
    font-weight: 600;
    font-size: 0.85rem;
    white-space: nowrap;
  }
  .post-content :global(td) {
    padding: 0.6rem 0.75rem;
    border-bottom: 1px solid var(--border-color);
  }
  .post-content :global(tr:last-child td) {
    border-bottom: none;
  }
  .post-content :global(del) {
    opacity: 0.5;
  }
  :global(body.light-mode) .post-content :global(.code-copy) {
    background: rgba(0, 0, 0, 0.04);
    color: rgba(0, 0, 0, 0.6);
  }
  :global(body.light-mode) .post-content :global(.code-copy:hover) {
    color: var(--light-accent, #d94100);
    border-color: rgba(217, 65, 0, 0.5);
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
  }
  .author-avatar {
    width: 40px;
    height: 40px;
    background: rgba(255, 68, 0, 0.1);
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
    font-family: var(--font-heading);
    font-size: 0.7rem;
    color: var(--text-muted);
    text-transform: uppercase;
    letter-spacing: 0.1em;
  }
  .share-links a {
    color: var(--text-secondary);
    transition: color 0.2s ease, border-color 0.2s ease;
    display: flex;
    padding: 0.55rem;
    border: 1px solid var(--border-color);
  }
  .share-links a:hover {
    color: var(--accent);
    border-color: var(--accent);
    background: rgba(255, 68, 0, 0.06);
  }
  .share-links a:focus-visible {
    outline: 2px solid var(--accent);
    outline-offset: 2px;
  }
  .share-links span {
    margin-right: 0.25rem;
  }

  /* ── Prev / next navigation (Ghost / dev.to pattern) ── */
  .post-nav {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
    margin-top: 2rem;
  }
  .pn-item {
    display: flex;
    flex-direction: column;
    gap: 0.35rem;
    padding: 1rem 1.25rem;
    background: var(--card-bg);
    border: 1px solid var(--border-color);
    border-left: 3px solid var(--border-color);
    text-decoration: none;
    transition: border-left-color 0.2s ease, transform 0.2s ease, box-shadow 0.2s ease;
  }
  .pn-item:hover {
    border-left-color: var(--accent);
    transform: translateY(-3px);
    box-shadow: 5px 5px 0 rgba(0, 0, 0, 0.4);
  }
  .pn-item:focus-visible {
    outline: 2px solid var(--accent);
    outline-offset: 3px;
  }
  .pn-item.pn-next {
    text-align: right;
  }
  .pn-label {
    font-family: var(--font-heading);
    font-size: 0.65rem;
    letter-spacing: 0.14em;
    text-transform: uppercase;
    color: var(--text-muted);
  }
  .pn-item:hover .pn-label {
    color: var(--accent);
  }
  .pn-title {
    font-family: var(--font-heading);
    font-size: 0.95rem;
    color: var(--text-primary);
    line-height: 1.35;
    overflow-wrap: anywhere;
  }
  :global(body.light-mode) .pn-item:hover {
    box-shadow: 5px 5px 0 rgba(0, 0, 0, 0.14);
  }

  @media (max-width: 600px) {
    .blog-post-page {
      padding: calc(var(--nav-height) + 1.5rem) 1rem 3rem;
    }
    .post-article {
      padding: 1.25rem;
    }
    .post-title {
      font-size: 1.6rem;
    }
    .post-content {
      font-size: 1rem;
    }
    .post-nav {
      grid-template-columns: 1fr;
    }
  }

  @media (max-width: 400px) {
    .post-article {
      padding: 1rem;
    }
    .post-title {
      font-size: 1.4rem;
    }
  }

  @media (prefers-reduced-motion: reduce) {
    .reading-progress {
      transition: none;
    }
    .post-content :global(.code-copy),
    .toc-item,
    .crumb,
    .share-links a,
    .pn-item {
      transition: none;
    }
    .pn-item:hover {
      transform: none;
    }
  }

  .post-error {
    color: var(--text-muted);
    text-align: center;
    padding: 3rem;
  }
</style>
