<script lang="ts">


    const API_URL = (import.meta.env.VITE_API_URL || "").replace(/\/$/, "") + "/api/contact";

    let name = $state("");
    let email = $state("");
    let topic = $state("general");
    let message = $state("");
    let submitting = $state(false);
    let success = $state(false);
    let errorMsg = $state("");

    let nameTouched = $state(false);
    let emailTouched = $state(false);
    let messageTouched = $state(false);

    let nameError = $derived(nameTouched && !name.trim() ? "Name is required" : "");
    let emailError = $derived(emailTouched && !email.trim() ? "Email is required" : emailTouched && !/^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$/.test(email) ? "Invalid email format" : "");
    let messageError = $derived(messageTouched && !message.trim() ? "Message is required" : messageTouched && message.trim().length < 10 ? "At least 10 characters" : "");

    function validate(): string {
        nameTouched = true;
        emailTouched = true;
        messageTouched = true;
        if (nameError || emailError || messageError) return "Please fix the errors above.";
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
                nameTouched = false; emailTouched = false; messageTouched = false;
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
    <div class="form-body">
        <div class="form-row">
            <div class="input-card" class:has-error={nameError}>
                <label class="input-label" for="name-input">Name</label>
                <div class="input-shell">
                    <input
                        id="name-input"
                        type="text"
                        placeholder="Your name"
                        bind:value={name}
                        disabled={submitting}
                        autocomplete="name"
                        onblur={() => nameTouched = true}
                        oninput={() => { nameTouched = true; errorMsg = ""; }}
                    />
                    <span class="input-accent"></span>
                </div>
                {#if nameError}
                    <span class="field-error">{nameError}</span>
                {/if}
            </div>
            <div class="input-card" class:has-error={emailError}>
                <label class="input-label" for="email-input">Email</label>
                <div class="input-shell">
                    <input
                        id="email-input"
                        type="email"
                        placeholder="your@email.com"
                        bind:value={email}
                        disabled={submitting}
                        autocomplete="email"
                        onblur={() => emailTouched = true}
                        oninput={() => { emailTouched = true; errorMsg = ""; }}
                    />
                    <span class="input-accent"></span>
                </div>
                {#if emailError}
                    <span class="field-error">{emailError}</span>
                {/if}
            </div>
        </div>

        <div class="input-card topic-section">
            <label class="input-label">Topic</label>
            <div class="topic-grid">
                <button
                    type="button"
                    class="topic-chip"
                    class:active={topic === 'general'}
                    onclick={() => topic = 'general'}
                    disabled={submitting}
                >
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <circle cx="12" cy="12" r="10"></circle>
                        <path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"></path>
                        <line x1="12" y1="17" x2="12.01" y2="17"></line>
                    </svg>
                    <span class="chip-label">General</span>
                    <span class="chip-dot"></span>
                </button>
                <button
                    type="button"
                    class="topic-chip"
                    class:active={topic === 'project'}
                    onclick={() => topic = 'project'}
                    disabled={submitting}
                >
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <polygon points="12 2 2 7 12 12 22 7 12 2"></polygon>
                        <polyline points="2 17 12 22 22 17"></polyline>
                        <polyline points="2 12 12 17 22 12"></polyline>
                    </svg>
                    <span class="chip-label">Project</span>
                    <span class="chip-dot"></span>
                </button>
                <button
                    type="button"
                    class="topic-chip"
                    class:active={topic === 'freelance'}
                    onclick={() => topic = 'freelance'}
                    disabled={submitting}
                >
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
                        <circle cx="12" cy="7" r="4"></circle>
                    </svg>
                    <span class="chip-label">Freelance</span>
                    <span class="chip-dot"></span>
                </button>
                <button
                    type="button"
                    class="topic-chip"
                    class:active={topic === 'other'}
                    onclick={() => topic = 'other'}
                    disabled={submitting}
                >
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
                    </svg>
                    <span class="chip-label">Other</span>
                    <span class="chip-dot"></span>
                </button>
            </div>
        </div>

        <div class="input-card" class:has-error={messageError}>
            <label class="input-label" for="msg-input">Message</label>
            <div class="input-shell">
                <textarea
                    id="msg-input"
                    placeholder="Tell me about your project, idea, or just say hi..."
                    bind:value={message}
                    rows={4}
                    disabled={submitting}
                    autocomplete="off"
                    onblur={() => messageTouched = true}
                    oninput={() => { messageTouched = true; errorMsg = ""; }}
                ></textarea>
                <span class="input-accent"></span>
            </div>
            {#if messageError}
                <span class="field-error">{messageError}</span>
            {/if}
        </div>

        <button
            class="submit-btn"
            class:sending={submitting}
            class:done={success}
            onclick={handleSubmit}
            disabled={submitting || success}
            type="button"
        >
            {#if submitting}
                <span class="btn-content">
                    <span class="btn-text">Sending</span>
                    <span class="btn-dots"><span>.</span><span>.</span><span>.</span></span>
                </span>
            {:else if success}
                <span class="btn-content">
                    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                        <polyline points="20 6 9 17 4 12"></polyline>
                    </svg>
                    <span class="btn-text">Sent</span>
                </span>
            {:else}
                <span class="btn-content">
                    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <line x1="22" y1="2" x2="11" y2="13"></line>
                        <polygon points="22 2 15 22 11 13 2 9 22 2"></polygon>
                    </svg>
                    <span class="btn-text">Send Message</span>
                </span>
            {/if}
        </button>

        {#if errorMsg}
            <div class="form-error" role="alert">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <circle cx="12" cy="12" r="10"></circle>
                    <line x1="12" y1="8" x2="12" y2="12"></line>
                    <line x1="12" y1="16" x2="12.01" y2="16"></line>
                </svg>
                <span>{errorMsg}</span>
            </div>
        {/if}
    </div>

    <div class="form-footer">
        <span class="form-reply">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="10"></circle>
                <polyline points="12 6 12 12 16 14"></polyline>
            </svg>
            Typically responds within 24 hours
        </span>
    </div>
</div>

<style>
    .contact-form {
        padding: 0;
        display: flex;
        flex-direction: column;
        gap: 1.5rem;
    }
    .form-body {
        display: flex;
        flex-direction: column;
        gap: 1.25rem;
    }

    /* ---- Row ---- */
    .form-row {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 1.25rem;
    }
    @media (max-width: 600px) {
        .form-row { grid-template-columns: 1fr; }
    }

    /* ---- Input Cards ---- */
    .input-card {
        display: flex;
        flex-direction: column;
        gap: 0.35rem;
    }
    .input-label {
        font-family: var(--font-body);
        font-size: 0.65rem;
        font-weight: 700;
        letter-spacing: 0.12em;
        text-transform: uppercase;
        color: var(--text-secondary);
        transition: color 0.2s ease;
        cursor: pointer;
    }
    .input-card:focus-within .input-label {
        color: var(--accent);
    }
    .input-card.has-error .input-label {
        color: #ff4444;
    }

    .input-shell {
        position: relative;
        display: flex;
        align-items: center;
        width: 100%;
        background: transparent;
        border-bottom: 2px solid var(--grid-line);
        transition: border-color 0.2s ease;
    }
    .input-card:focus-within .input-shell {
        border-bottom-color: var(--accent);
    }
    .input-card.has-error .input-shell {
        border-bottom-color: #ff4444;
    }

    .input-shell input,
    .input-shell textarea {
        width: 100%;
        background: transparent;
        border: none;
        outline: none;
        color: var(--text-primary);
        font-family: var(--font-body);
        font-size: 0.95rem;
        font-weight: 500;
        padding: 0.7rem 0;
        resize: vertical;
        min-height: 1.5rem;
    }
    .input-shell input::placeholder,
    .input-shell textarea::placeholder {
        color: var(--text-secondary);
        opacity: 0.5;
        font-weight: 400;
    }
    .input-shell input:disabled,
    .input-shell textarea:disabled {
        opacity: 0.4;
        cursor: not-allowed;
    }

    .input-accent {
        position: absolute;
        bottom: -2px;
        left: 0;
        width: 100%;
        height: 2px;
        background: var(--accent);
        transform: scaleX(0);
        transform-origin: left;
        transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    }
    .input-card:focus-within .input-accent {
        transform: scaleX(1);
    }
    .input-card.has-error .input-accent {
        background: #ff4444;
        transform: scaleX(1);
    }

    .field-error {
        font-size: 0.7rem;
        color: #ff4444;
        font-weight: 500;
        padding-top: 0.15rem;
        animation: errSlide 0.25s ease;
    }

    /* ---- Topic Grid ---- */
    .topic-section {
        gap: 0.5rem;
    }
    .topic-grid {
        display: grid;
        grid-template-columns: repeat(4, 1fr);
        gap: 0.5rem;
    }
    @media (max-width: 700px) {
        .topic-grid { grid-template-columns: repeat(2, 1fr); }
    }
    @media (max-width: 400px) {
        .topic-grid { grid-template-columns: 1fr; }
    }

    .topic-chip {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.65rem 0.75rem;
        background: transparent;
        border: 1.5px solid var(--grid-line);
        clip-path: polygon(6px 0, 100% 0, 100% calc(100% - 6px), calc(100% - 6px) 100%, 0 100%, 0 6px);
        cursor: pointer;
        font-family: var(--font-body);
        transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
        position: relative;
        color: var(--text-secondary);
    }
    .topic-chip svg {
        flex-shrink: 0;
        transition: color 0.2s ease;
    }
    .chip-label {
        font-size: 0.75rem;
        font-weight: 600;
        text-transform: uppercase;
        letter-spacing: 0.06em;
        flex: 1;
        text-align: left;
        transition: color 0.2s ease;
    }
    .chip-dot {
        width: 6px;
        height: 6px;
        border-radius: 50%;
        background: var(--grid-line);
        transition: all 0.25s ease;
        flex-shrink: 0;
    }
    .topic-chip:hover:not(.active):not(:disabled) {
        border-color: var(--accent);
        color: var(--accent);
        background: rgba(255, 68, 0, 0.04);
    }
    .topic-chip:hover:not(.active):not(:disabled) .chip-dot {
        background: var(--accent);
        transform: scale(1.4);
    }
    .topic-chip.active {
        border-color: var(--accent);
        background: linear-gradient(135deg, rgba(255, 68, 0, 0.08), rgba(255, 68, 0, 0.03));
        color: var(--accent);
    }
    .topic-chip.active .chip-label {
        color: var(--accent);
    }
    .topic-chip.active .chip-dot {
        background: var(--accent);
        box-shadow: 0 0 6px var(--accent-glow);
        animation: chipPulse 2s ease infinite;
    }
    .topic-chip:disabled {
        opacity: 0.4;
        cursor: not-allowed;
    }

    /* ---- Submit Button ---- */
    .submit-btn {
        width: 100%;
        padding: 1rem 1.5rem;
        border: 1.5px solid var(--accent);
        background: transparent;
        clip-path: polygon(10px 0, 100% 0, 100% calc(100% - 10px), calc(100% - 10px) 100%, 0 100%, 0 10px);
        cursor: pointer;
        font-family: var(--font-body);
        transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
        position: relative;
        overflow: hidden;
    }
    .submit-btn::before {
        content: '';
        position: absolute;
        inset: 0;
        background: var(--accent);
        opacity: 0;
        transition: opacity 0.25s ease;
    }
    .btn-content {
        position: relative;
        z-index: 1;
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 0.5rem;
    }
    .submit-btn .btn-text {
        font-size: 0.85rem;
        font-weight: 700;
        text-transform: uppercase;
        letter-spacing: 0.1em;
        color: var(--accent);
        transition: color 0.25s ease;
    }
    .submit-btn svg {
        color: var(--accent);
        transition: all 0.25s ease;
    }
    .submit-btn:hover:not(:disabled)::before {
        opacity: 1;
    }
    .submit-btn:hover:not(:disabled) .btn-text,
    .submit-btn:hover:not(:disabled) svg {
        color: #fff;
    }
    .submit-btn:active:not(:disabled) {
        transform: scale(0.98);
    }
    .submit-btn:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }
    .submit-btn.done {
        border-color: #00cc88;
        background: rgba(0, 204, 136, 0.08);
    }
    .submit-btn.done .btn-text {
        color: #00cc88;
    }
    .submit-btn.done svg {
        color: #00cc88;
    }

    .btn-dots span {
        animation: dotPulse 1.4s ease-in-out infinite;
        color: var(--accent);
    }
    .btn-dots span:nth-child(2) { animation-delay: 0.2s; }
    .btn-dots span:nth-child(3) { animation-delay: 0.4s; }

    /* ---- Error Banner ---- */
    .form-error {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        color: #ff4444;
        font-size: 0.8rem;
        font-weight: 500;
        padding: 0.7rem 1rem;
        border: 1px solid rgba(255, 68, 68, 0.25);
        background: rgba(255, 68, 68, 0.05);
        clip-path: polygon(6px 0, 100% 0, 100% calc(100% - 6px), calc(100% - 6px) 100%, 0 100%, 0 6px);
        animation: shake 0.3s ease;
    }
    .form-error svg {
        flex-shrink: 0;
    }

    /* ---- Footer ---- */
    .form-footer {
        display: flex;
        justify-content: center;
        padding-top: 0.5rem;
    }
    .form-reply {
        display: flex;
        align-items: center;
        gap: 0.4rem;
        font-family: var(--font-body);
        font-size: 0.75rem;
        color: var(--text-secondary);
        opacity: 0.65;
    }
    .form-reply svg {
        flex-shrink: 0;
    }

    /* ---- Animations ---- */
    @keyframes errSlide {
        from { opacity: 0; transform: translateY(-4px); }
        to { opacity: 1; transform: translateY(0); }
    }
    @keyframes shake {
        0%, 100% { transform: translateX(0); }
        20% { transform: translateX(-6px); }
        40% { transform: translateX(6px); }
        60% { transform: translateX(-4px); }
        80% { transform: translateX(4px); }
    }
    @keyframes chipPulse {
        0%, 100% { box-shadow: 0 0 4px var(--accent-glow); }
        50% { box-shadow: 0 0 10px var(--accent-glow); }
    }
    @keyframes dotPulse {
        0%, 80%, 100% { opacity: 0; }
        40% { opacity: 1; }
    }

    /* ---- Light Mode ---- */
    :global(body.light-mode) .input-shell {
        border-bottom-color: rgba(0, 0, 0, 0.2);
    }
    :global(body.light-mode) .input-card:focus-within .input-shell {
        border-bottom-color: var(--accent);
    }
    :global(body.light-mode) .input-card.has-error .input-shell {
        border-bottom-color: #ff4444;
    }
    :global(body.light-mode) .input-shell input::placeholder,
    :global(body.light-mode) .input-shell textarea::placeholder {
        color: rgba(0, 0, 0, 0.35);
    }
    :global(body.light-mode) .topic-chip {
        border-color: rgba(0, 0, 0, 0.2);
        color: var(--text-secondary);
    }
    :global(body.light-mode) .topic-chip.active {
        border-color: var(--accent);
        background: rgba(217, 61, 0, 0.06);
    }
    :global(body.light-mode) .topic-chip:hover:not(.active):not(:disabled) {
        border-color: var(--accent);
        background: rgba(217, 61, 0, 0.04);
    }
    :global(body.light-mode) .submit-btn::before {
        background: var(--accent);
    }
    :global(body.light-mode) .submit-btn:hover:not(:disabled) .btn-text,
    :global(body.light-mode) .submit-btn:hover:not(:disabled) svg {
        color: #fff;
    }
    :global(body.light-mode) .form-error {
        border-color: rgba(217, 61, 0, 0.3);
        background: rgba(217, 61, 0, 0.04);
    }
    :global(body.light-mode) .form-reply {
        opacity: 0.5;
    }

    /* ---- Reduced Motion ---- */
    @media (prefers-reduced-motion: reduce) {
        .input-accent { transition: none; }
        .topic-chip { transition: none; }
        .chip-dot { animation: none; }
        .submit-btn { transition: none; }
        .submit-btn::before { transition: none; }
        .btn-dots span { animation: none; }
        .form-error { animation: none; }
        .field-error { animation: none; }
    }

    /* ---- Tiny Screens ---- */
    @media (max-width: 400px) {
        .form-body { gap: 1rem; }
        .form-row { gap: 1rem; }
        .topic-chip { padding: 0.55rem 0.6rem; }
        .topic-chip svg { width: 14px; height: 14px; }
        .chip-label { font-size: 0.65rem; }
        .submit-btn { padding: 0.85rem 1rem; }
        .submit-btn .btn-text { font-size: 0.75rem; }
    }
</style>
