<script lang="ts">
    import { page } from "$app/stores";
</script>

<svelte:head>
    <title>SYSTEM::ERROR {$page.status}</title>
</svelte:head>

<div class="error-page">
    <div class="error-container">
        <h1 class="error-code">{$page.status}</h1>
        <h2 class="error-heading">SYSTEM::ERROR</h2>
        <p class="error-message">{$page.error?.message}</p>
        <a href="/" class="btn">try again</a>
    </div>
</div>

<style>
    .error-page {
        display: flex;
        align-items: center;
        justify-content: center;
        min-height: 100vh;
        padding: 2rem;
        background: var(--bg-dark);
        font-family: var(--font-body);
    }
    :global(body.light-mode) .error-page {
        background: var(--bg-dark);
    }
    .error-container {
        text-align: center;
        max-width: 500px;
    }
    .error-code {
        font-family: var(--font-heading);
        font-size: 6rem;
        color: var(--accent);
        line-height: 1;
        margin-bottom: 1rem;
        opacity: 0.6;
    }
    .error-heading {
        font-family: var(--font-heading);
        font-size: 1.5rem;
        color: var(--text-primary);
        margin-bottom: 1rem;
        text-transform: uppercase;
        letter-spacing: 0.15em;
    }
    .error-message {
        font-family: var(--font-body);
        color: var(--text-secondary);
        margin-bottom: 2.5rem;
        font-size: 0.9rem;
        line-height: 1.6;
    }
    .btn {
        display: inline-block;
        font-family: var(--font-heading);
        text-transform: uppercase;
        letter-spacing: 0.08em;
        font-size: 0.85rem;
        padding: 15px 30px;
        font-weight: 600;
        border: 2px solid var(--accent);
        color: var(--accent);
        background: transparent;
        cursor: pointer;
        box-shadow: 4px 4px 0px var(--accent);
        transition: all 0.1s cubic-bezier(0, 0, 0.2, 1);
        text-decoration: none;
    }
    .btn:hover {
        transform: translate(2px, 2px);
        box-shadow: 2px 2px 0px var(--accent);
        background: rgba(255, 68, 0, 0.15);
        color: var(--accent);
        border-color: var(--accent);
    }
    .btn:active {
        transform: translate(4px, 4px);
        box-shadow: none;
    }
</style>
