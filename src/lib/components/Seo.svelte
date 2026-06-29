<script lang="ts">
  import { SITE } from '$lib/config';

  let {
    title = SITE.title,
    description = SITE.description,
    image = SITE.image,
    url = SITE.url,
    type = 'website',
    canonical,
    publishedTime,
    modifiedTime,
    tags = [] as string[],
    jsonld,
  }: {
    title?: string;
    description?: string;
    image?: string;
    url?: string;
    type?: string;
    canonical?: string;
    publishedTime?: string;
    modifiedTime?: string;
    tags?: string[];
    jsonld?: Record<string, unknown>;
  } = $props();

  const canonicalUrl = canonical || url;
  const siteTitle = title === SITE.title ? title : `${title} — ${SITE.name}`;

  $effect(() => {
    if (typeof document !== 'undefined') {
      document.title = siteTitle;
    }
  });
</script>

<svelte:head>
  <title>{siteTitle}</title>
  <meta name="description" content={description} />
  <meta name="author" content={SITE.author} />
  <link rel="canonical" href={canonicalUrl} />

  <meta property="og:type" content={type} />
  <meta property="og:url" content={url} />
  <meta property="og:title" content={siteTitle} />
  <meta property="og:description" content={description} />
  <meta property="og:image" content={image} />
  <meta property="og:image:width" content="1200" />
  <meta property="og:image:height" content="630" />
  <meta property="og:site_name" content={SITE.name} />
  <meta property="og:locale" content="en_US" />

  <meta name="twitter:card" content="summary_large_image" />
  <meta name="twitter:url" content={url} />
  <meta name="twitter:title" content={siteTitle} />
  <meta name="twitter:description" content={description} />
  <meta name="twitter:image" content={image} />
  <meta name="twitter:creator" content={SITE.twitter} />

  {#if jsonld}
    {@html `<script type="application/ld+json">${JSON.stringify(jsonld)}</script>`}
  {/if}
</svelte:head>
