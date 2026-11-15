<script lang="ts">
  let { slug }: { slug: string } = $props();
</script>

<div class="comments-section">
  <h2 class="comments-heading">Comments</h2>

  <script
    src="https://giscus.app/client.js"
    data-repo="PranavAgarkar07/portfolio-svelte"
    data-repo-id="R_kgDOPx7Q4A"
    data-category="Comments"
    data-category-id="DIC_kwDOPx7Q4C4C6X9m"
    data-mapping="pathname"
    data-term="{slug}"
    data-reactions-enabled="1"
    data-emit-metadata="0"
    data-input-position="top"
    data-theme="preferred_color_scheme"
    data-lang="en"
    data-loading="lazy"
    crossorigin="anonymous"
    async>
  </script>

  <noscript>
    <p class="empty-text">Comments require JavaScript. You can also discuss this post on <a href="https://github.com/PranavAgarkar07/portfolio-svelte/discussions" target="_blank" rel="noopener">GitHub Discussions</a>.</p>
  </noscript>
</div>

<style>
  .comments-section {
    margin-top: 3rem;
    padding-top: 1.5rem;
    border-top: 1px solid var(--border-color);
  }
  .comments-heading {
    font-family: var(--font-heading);
    font-size: 1.2rem;
    color: var(--text-primary);
    margin: 0 0 1rem;
  }
  .empty-text {
    color: var(--text-muted);
    text-align: center;
    padding: 1.5rem;
    font-size: 0.85rem;
  }
</style>
