<script lang="ts">
    import { onMount } from "svelte";
    import { browser } from "$app/environment";

    let visible = $state(false);

    const CONSENT_KEY = "portfolio_analytics_consent";

    function hasConsent(): boolean | null {
        if (!browser) return null;
        const v = localStorage.getItem(CONSENT_KEY);
        if (v === "true") return true;
        if (v === "false") return false;
        return null;
    }

    function setConsent(val: boolean) {
        localStorage.setItem(CONSENT_KEY, String(val));
        visible = false;
        if (val) {
            import("$lib/analytics").then(m => m.initAnalytics());
        }
    }

    onMount(() => {
        if (hasConsent() === null) {
            setTimeout(() => { visible = true; }, 1000);
        }
    });
</script>

{#if visible}
    <div class="cookie-consent" role="dialog" aria-label="Cookie consent">
        <p class="cookie-text">
            This site uses basic analytics (page views, device type) to improve your experience.
            No personal data is collected or shared.
        </p>
        <div class="cookie-actions">
            <button class="cookie-btn cookie-accept" onclick={() => setConsent(true)}>Accept</button>
            <button class="cookie-btn cookie-decline" onclick={() => setConsent(false)}>Decline</button>
        </div>
    </div>
{/if}

<style>
    .cookie-consent {
        position: fixed;
        bottom: 1.5rem;
        left: 50%;
        transform: translateX(-50%);
        z-index: 10001;
        max-width: 480px;
        width: calc(100% - 2rem);
        background: #1a1a1a;
        border: 1px solid rgba(255,255,255,0.08);
        padding: 1rem 1.25rem;
        display: flex;
        flex-direction: column;
        gap: 0.75rem;
        box-shadow: 0 8px 32px rgba(0,0,0,0.4);
    }

    .cookie-text {
        font-size: 0.8rem;
        line-height: 1.5;
        color: #a1a1aa;
        margin: 0;
    }

    .cookie-actions {
        display: flex;
        gap: 0.5rem;
        justify-content: flex-end;
    }

    .cookie-btn {
        padding: 0.4rem 1rem;
        font-size: 0.75rem;
        font-family: inherit;
        cursor: pointer;
        border: 1px solid rgba(255,255,255,0.1);
        background: none;
        color: #e4e4e7;
        transition: all 0.15s ease;
    }

    .cookie-accept {
        background: var(--accent);
        border-color: var(--accent);
        color: #030405;
        font-weight: 600;
    }

    .cookie-accept:hover {
        opacity: 0.9;
    }

    .cookie-decline:hover {
        border-color: rgba(255,255,255,0.25);
    }

    :global(.light-mode) .cookie-consent {
        background: #f4f4f5;
        border-color: rgba(0,0,0,0.08);
    }

    :global(.light-mode) .cookie-text {
        color: #52525b;
    }

    :global(.light-mode) .cookie-btn {
        color: #18181b;
        border-color: rgba(0,0,0,0.1);
    }

    :global(.light-mode) .cookie-accept {
        color: #f4f4f5;
    }

    @media (max-width: 480px) {
        .cookie-consent {
            bottom: 0;
            left: 0;
            transform: none;
            width: 100%;
            max-width: 100%;
            border-width: 1px 0 0 0;
        }
    }
</style>
