<script lang="ts">
  import type { BlogImage } from '$lib/types';

  let {
    images,
  }: {
    images: BlogImage[];
  } = $props();

  let lightboxIdx = $state<number | null>(null);

  function open(idx: number) { lightboxIdx = idx; }
  function close() { lightboxIdx = null; }
  function prev() { if (lightboxIdx !== null) lightboxIdx = lightboxIdx > 0 ? lightboxIdx - 1 : images.length - 1; }
  function next() { if (lightboxIdx !== null) lightboxIdx = lightboxIdx < images.length - 1 ? lightboxIdx + 1 : 0; }

  function onKeydown(e: KeyboardEvent) {
    if (lightboxIdx === null) return;
    if (e.key === 'Escape') close();
    if (e.key === 'ArrowLeft') prev();
    if (e.key === 'ArrowRight') next();
  }
</script>

<svelte:window onkeydown={onKeydown} />

<div class="gallery">
  {#each images as img, i}
    <button class="gallery-item" onclick={() => open(i)} aria-label={img.alt || `Image ${i + 1}`}>
      <img src={img.url} alt={img.alt || ''} loading="lazy" />
      {#if img.caption}
        <span class="gallery-caption">{img.caption}</span>
      {/if}
    </button>
  {/each}
</div>

{#if lightboxIdx !== null}
  <div class="lightbox-backdrop" onclick={close} role="presentation"></div>
  <div class="lightbox" role="dialog" aria-label="Image lightbox">
    <button class="lightbox-close" onclick={close} aria-label="Close">&times;</button>
    <button class="lightbox-nav prev" onclick={prev} aria-label="Previous">&lsaquo;</button>
    <button class="lightbox-nav next" onclick={next} aria-label="Next">&rsaquo;</button>
    <div class="lightbox-content">
      <img src={images[lightboxIdx].url} alt={images[lightboxIdx].alt || ''} />
      {#if images[lightboxIdx].caption}
        <p class="lightbox-caption">{images[lightboxIdx].caption}</p>
      {/if}
      <p class="lightbox-counter">{lightboxIdx + 1} / {images.length}</p>
    </div>
  </div>
{/if}

<style>
  .gallery {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 0.75rem;
    margin: 2rem 0;
  }
  .gallery-item {
    position: relative;
    border: 1px solid var(--border-color);
    overflow: hidden;
    cursor: pointer;
    background: var(--surface-dark);
    padding: 0;
    display: block;
    text-align: left;
    transition: border-color 0.15s;
  }
  .gallery-item:hover {
    border-color: var(--accent);
  }
  .gallery-item img {
    width: 100%;
    aspect-ratio: 16 / 10;
    object-fit: cover;
    display: block;
    transition: transform 0.3s ease;
  }
  .gallery-item:hover img {
    transform: scale(1.03);
  }
  .gallery-caption {
    display: block;
    padding: 0.4rem 0.5rem;
    font-size: 0.75rem;
    color: var(--text-secondary);
    line-height: 1.4;
    border-top: 1px solid var(--border-color);
  }

  .lightbox-backdrop {
    position: fixed;
    inset: 0;
    background: rgba(0,0,0,0.85);
    z-index: 200;
  }
  .lightbox {
    position: fixed;
    inset: 0;
    z-index: 201;
    display: flex;
    align-items: center;
    justify-content: center;
    pointer-events: none;
  }
  .lightbox-content {
    pointer-events: auto;
    max-width: 90vw;
    max-height: 90vh;
    display: flex;
    flex-direction: column;
    align-items: center;
  }
  .lightbox-content img {
    max-width: 100%;
    max-height: 80vh;
    object-fit: contain;
    border: 1px solid var(--border-color);
  }
  .lightbox-caption {
    color: var(--text-secondary);
    font-size: 0.85rem;
    margin: 0.75rem 0 0.25rem;
    text-align: center;
  }
  .lightbox-counter {
    color: var(--text-muted);
    font-size: 0.75rem;
    margin: 0;
    text-align: center;
  }
  .lightbox-close {
    position: fixed;
    top: 1rem;
    right: 1.5rem;
    background: none;
    border: none;
    color: #fff;
    font-size: 2rem;
    cursor: pointer;
    z-index: 202;
    pointer-events: auto;
    opacity: 0.7;
    transition: opacity 0.15s;
  }
  .lightbox-close:hover {
    opacity: 1;
  }
  .lightbox-nav {
    position: fixed;
    top: 50%;
    transform: translateY(-50%);
    background: rgba(0,0,0,0.5);
    border: none;
    color: #fff;
    font-size: 2.5rem;
    padding: 0.5rem 0.75rem;
    cursor: pointer;
    z-index: 202;
    pointer-events: auto;
    opacity: 0.5;
    transition: opacity 0.15s;
    line-height: 1;
  }
  .lightbox-nav:hover {
    opacity: 1;
  }
  .lightbox-nav.prev {
    left: 1rem;
  }
  .lightbox-nav.next {
    right: 1rem;
  }

  @media (max-width: 600px) {
    .gallery {
      grid-template-columns: 1fr 1fr;
      gap: 0.5rem;
    }
  }
  @media (max-width: 400px) {
    .gallery {
      grid-template-columns: 1fr;
    }
  }
</style>
