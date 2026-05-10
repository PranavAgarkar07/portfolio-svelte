<script lang="ts">
    import { Button, Input, Textarea } from "$lib/components/ui";

    const API_URL = (import.meta.env.VITE_API_URL || "").replace(/\/$/, "") + "/api/contact";

    let name = $state("");
    let email = $state("");
    let topic = $state("general");
    let message = $state("");
    let submitting = $state(false);
    let success = $state(false);
    let errorMsg = $state("");

    function validate(): string {
        if (!name.trim()) return "Name is required.";
        if (!email.trim()) return "Email is required.";
        if (!/^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$/.test(email)) return "Invalid email format.";
        if (!message.trim()) return "Message is required.";
        if (message.trim().length < 10) return "Message must be at least 10 characters.";
        return "";
    }

    async function handleSubmit() {
        const err = validate();
        if (err) { errorMsg = err; return; }
        submitting = true;
        errorMsg = "";
        try {
            const res = await fetch(API_URL, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ name: name.trim(), email: email.trim(), topic, message: message.trim() }),
            });
            if (res.ok) {
                success = true;
                name = ""; email = ""; topic = "general"; message = "";
                setTimeout(() => success = false, 6000);
            } else if (res.status === 429) {
                errorMsg = "Please wait a moment before sending another message.";
            } else {
                const data = await res.json().catch(() => ({}));
                errorMsg = data.message || "Something went wrong. Try emailing me directly at pranavagarkar8@gmail.com.";
            }
        } catch {
            errorMsg = "Network error. Try emailing me directly at pranavagarkar8@gmail.com.";
        } finally {
            submitting = false;
        }
    }

</script>

<div class="contact-form">
    <div class="form-row">
        <div class="form-group" style="animation-delay: 0.05s">
            <Input label="Name" type="text" placeholder="Your Name" bind:value={name} disabled={submitting} autocomplete="name" />
        </div>
        <div class="form-group" style="animation-delay: 0.1s">
            <Input label="Email" type="email" placeholder="Your Email" bind:value={email} disabled={submitting} autocomplete="email" />
        </div>
    </div>
    <div class="form-group" style="animation-delay: 0.15s">
        <label class="form-label">Topic</label>
        <div class="topic-grid">
            <label class="topic-card" class:active={topic === 'general'}>
                <input type="radio" name="topic" value="general" bind:group={topic} disabled={submitting}>
                <div class="topic-icon">
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <circle cx="12" cy="12" r="10"></circle>
                        <path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"></path>
                        <line x1="12" y1="17" x2="12.01" y2="17"></line>
                    </svg>
                </div>
                <div class="topic-content">
                    <span class="topic-title">General Inquiry</span>
                    <span class="topic-desc">Questions & info</span>
                </div>
                <div class="topic-indicator"></div>
            </label>
            <label class="topic-card" class:active={topic === 'project'}>
                <input type="radio" name="topic" value="project" bind:group={topic} disabled={submitting}>
                <div class="topic-icon">
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <polygon points="12 2 2 7 12 12 22 7 12 2"></polygon>
                        <polyline points="2 17 12 22 22 17"></polyline>
                        <polyline points="2 12 12 17 22 12"></polyline>
                    </svg>
                </div>
                <div class="topic-content">
                    <span class="topic-title">Project</span>
                    <span class="topic-desc">Start a project</span>
                </div>
                <div class="topic-indicator"></div>
            </label>
            <label class="topic-card" class:active={topic === 'freelance'}>
                <input type="radio" name="topic" value="freelance" bind:group={topic} disabled={submitting}>
                <div class="topic-icon">
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
                        <circle cx="12" cy="7" r="4"></circle>
                    </svg>
                </div>
                <div class="topic-content">
                    <span class="topic-title">Freelance</span>
                    <span class="topic-desc">Work together</span>
                </div>
                <div class="topic-indicator"></div>
            </label>
            <label class="topic-card" class:active={topic === 'other'}>
                <input type="radio" name="topic" value="other" bind:group={topic} disabled={submitting}>
                <div class="topic-icon">
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
                    </svg>
                </div>
                <div class="topic-content">
                    <span class="topic-title">Other</span>
                    <span class="topic-desc">Anything else</span>
                </div>
                <div class="topic-indicator"></div>
            </label>
        </div>
    </div>
    <div class="form-group" style="animation-delay: 0.2s">
        <Textarea label="Message" placeholder="Your Message" bind:value={message} rows={4} disabled={submitting} autocomplete="off" />
    </div>
    <Button variant="primary" size="lg" style="width: 100%" disabled={submitting || success} onclick={handleSubmit}>
        {#if submitting}
            Sending<span class="btn-dots"><span>.</span><span>.</span><span>.</span></span>
        {:else if success}
            &#10003; Sent
        {:else}
            Send Message
        {/if}
    </Button>
    {#if errorMsg}
        <div class="form-error" role="alert">{errorMsg}</div>
    {/if}
    <div class="form-footer" style="animation-delay: 0.3s">
        <span class="form-reply">Typically responds within 24 hours</span>
    </div>
</div>

<style>
    .contact-form {
        padding: 0;
        display: flex;
        flex-direction: column;
        gap: 1.5rem;
    }
    .form-row {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 1.5rem;
    }
    @media (max-width: 600px) {
        .form-row {
            grid-template-columns: 1fr;
        }
    }
    .form-group {
        margin: 0;
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
        animation: formSlideUp 0.4s cubic-bezier(0.4, 0, 0.2, 1) both;
    }
    .form-label {
        font-family: var(--font-body);
        font-size: 0.7rem;
        color: var(--text-primary);
        font-weight: 700;
        letter-spacing: 0.1em;
        text-transform: uppercase;
        transition: color 0.1s ease;
        cursor: pointer;
    }


    .topic-grid {
        display: grid;
        grid-template-columns: repeat(2, 1fr);
        gap: 0.75rem;
    }
    @media (max-width: 500px) {
        .topic-grid { grid-template-columns: 1fr; }
    }
    .topic-card {
        display: flex;
        align-items: center;
        gap: 0.875rem;
        padding: 1rem 1.125rem;
        border: 2px solid var(--text-secondary);
        cursor: pointer;
        font-family: var(--font-body);
        transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
        background: transparent;
        position: relative;
        overflow: hidden;
        box-sizing: border-box;
    }
    .topic-card::before {
        content: '';
        position: absolute;
        inset: 0;
        background: rgba(255, 68, 0, 0.08);
        opacity: 0;
        transition: opacity 0.2s ease;
    }
    .topic-card input {
        display: none;
    }
    .topic-icon {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 40px;
        height: 40px;
        border: 1.5px solid currentColor;
        transition: all 0.2s ease;
        position: relative;
        z-index: 1;
        flex-shrink: 0;
    }
    .topic-content {
        display: flex;
        flex-direction: column;
        gap: 0.125rem;
        flex: 1;
        position: relative;
        z-index: 1;
    }
    .topic-title {
        font-size: 0.8rem;
        font-weight: 700;
        text-transform: uppercase;
        letter-spacing: 0.08em;
        transition: all 0.2s ease;
    }
    .topic-desc {
        font-size: 0.65rem;
        font-weight: 500;
        text-transform: uppercase;
        letter-spacing: 0.05em;
        opacity: 0.6;
        transition: all 0.2s ease;
    }
    .topic-indicator {
        width: 8px;
        height: 8px;
        border-radius: 50%;
        background: var(--text-secondary);
        transition: all 0.2s ease;
        position: relative;
        z-index: 1;
        flex-shrink: 0;
    }
    .topic-card:hover:not(.active) {
        border-color: var(--accent);
        transform: translate(-2px, -2px);
        box-shadow: 4px 4px 0px rgba(255, 68, 0, 0.15);
    }
    .topic-card:hover:not(.active) .topic-indicator {
        background: var(--accent);
        transform: scale(1.25);
    }
    .topic-card:hover:not(.active) .topic-icon {
        background: transparent;
        color: var(--accent);
        border-color: var(--accent);
    }
    .topic-card:hover:not(.active) .topic-desc {
        opacity: 1;
    }
    .topic-card.active {
        border-color: var(--accent);
        box-shadow: 4px 4px 0px rgba(255, 68, 0, 0.15);
        transform: translate(-2px, -2px);
    }
    .topic-card.active::before {
        opacity: 1;
    }
    .topic-card.active .topic-icon {
        background: #000;
        border-color: #000;
        color: var(--accent);
    }
    .topic-card.active .topic-title {
        color: var(--accent);
    }
    .topic-card.active .topic-desc {
        color: var(--text-primary);
        opacity: 0.8;
    }
    .topic-card.active .topic-indicator {
        background: var(--accent);
        transform: scale(1.5);
        animation: pulse 1.5s ease infinite;
    }
    @keyframes pulse {
        0%, 100% { transform: scale(1.5); }
        50% { transform: scale(2); }
    }
    :global(body.light-mode) .topic-card {
        border-color: #000;
    }
    :global(body.light-mode) .topic-card:hover:not(.active) {
        box-shadow: 4px 4px 0px #000;
    }
    :global(body.light-mode) .topic-card.active {
        box-shadow: 4px 4px 0px #000;
    }
    :global(body.light-mode) .topic-card.active::before {
        opacity: 1;
    }
    :global(body.light-mode) .topic-card.active .topic-icon {
        background: transparent;
        border-color: var(--accent);
        color: var(--accent);
    }
    :global(body.light-mode) .topic-card.active .topic-indicator {
        background: var(--accent);
    }
    @media (prefers-reduced-motion: reduce) {
        .topic-card,
        .topic-icon,
        .topic-indicator,
        .topic-title,
        .topic-desc {
            transition: none;
        }
        .topic-card.active .topic-indicator {
            animation: none;
        }
    }

    .btn-dots span {
        animation: dotPulse 1.4s ease-in-out infinite;
    }
    .btn-dots span:nth-child(2) { animation-delay: 0.2s; }
    .btn-dots span:nth-child(3) { animation-delay: 0.4s; }

    .form-error {
        color: #ff4444;
        font-size: 0.8rem;
        padding: 0.75rem;
        border: 1px solid rgba(255, 68, 68, 0.3);
        background: rgba(255, 68, 68, 0.06);
        animation: shake 0.3s ease;
    }

    .form-footer {
        display: flex;
        justify-content: space-between;
        align-items: center;
        gap: 1rem;
        padding-top: 0.25rem;
        animation: formSlideUp 0.4s cubic-bezier(0.4, 0, 0.2, 1) both;
    }
.form-reply {
        font-family: var(--font-body);
        font-size: 0.8rem;
        color: var(--text-secondary);
        opacity: 0.7;
    }

    /* REMOVED lock-icon - no longer used */

    @keyframes formSlideUp {
        from { opacity: 0; transform: translateY(12px); }
        to { opacity: 1; transform: translateY(0); }
    }
    @keyframes shake {
        0%, 100% { transform: translateX(0); }
        20% { transform: translateX(-6px); }
        40% { transform: translateX(6px); }
        60% { transform: translateX(-4px); }
        80% { transform: translateX(4px); }
    }
    @keyframes dotPulse {
        0%, 80%, 100% { opacity: 0; }
        40% { opacity: 1; }
    }

    @media (max-width: 400px) {
        .contact-form { gap: 1rem; }
        .form-row { gap: 1rem; }
        .topic-card { padding: 0.75rem 0.875rem; gap: 0.6rem; }
        .topic-icon { width: 32px; height: 32px; }
        .topic-icon :global(svg) { width: 16px; height: 16px; }
        .topic-title { font-size: 0.7rem; }
        .topic-desc { font-size: 0.6rem; }
        .form-label { font-size: 0.6rem; }
        .form-reply { font-size: 0.7rem; }
    }

    @media (prefers-reduced-motion: reduce) {
        .form-group,
        .form-footer {
            animation: none;
        }
        .btn-dots span {
            animation: none;
        }
        .topic-card,
        .topic-icon,
        .topic-indicator {
            transition: none;
        }
    }
</style>
