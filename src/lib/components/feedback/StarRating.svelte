<script lang="ts">
  let { value = 0, interactive = false, onRate, size = 'md' }: { value?: number; interactive?: boolean; onRate?: (n: number) => void; size?: 'sm' | 'md' } = $props();

  let hovered = $state(0);

  let starSize = $derived(size === 'sm' ? 16 : 22);
</script>

<div class="star-rating" class:interactive class:sm={size === 'sm'}>
  {#each [1, 2, 3, 4, 5] as star}
    <button
      class="star"
      class:filled={star <= (hovered || value)}
      class:interactive
      onmouseenter={() => interactive && (hovered = star)}
      onmouseleave={() => (hovered = 0)}
      onclick={() => interactive && onRate?.(star)}
      disabled={!interactive}
      aria-label="{star} star{star > 1 ? 's' : ''}"
    >
      <svg width={starSize} height={starSize} viewBox="0 0 24 24" fill={star <= (hovered || value) ? 'var(--accent)' : 'none'} stroke="var(--accent)" stroke-width="2">
        <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2" />
      </svg>
    </button>
  {/each}
</div>

<style>
  .star-rating {
    display: inline-flex;
    gap: 2px;
  }
  .star-rating.interactive {
    gap: 4px;
  }
  .star-rating.sm {
    gap: 0;
  }
  .star {
    background: none;
    border: none;
    padding: 4px;
    cursor: default;
    line-height: 1;
    transition: transform 0.1s;
    min-width: 0;
  }
  .star.interactive {
    cursor: pointer;
    padding: 6px;
  }
  .star-rating.interactive {
    gap: 0;
  }
  .star.interactive:hover {
    transform: scale(1.2);
  }
  .star.filled svg {
    fill: var(--accent);
  }
  .star:not(.filled) svg {
    fill: transparent;
  }
</style>
