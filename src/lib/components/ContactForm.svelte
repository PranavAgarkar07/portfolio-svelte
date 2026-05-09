<script lang="ts">
    import { onMount } from "svelte";
    const API_URL = (import.meta.env.VITE_API_URL || "").replace(/\/$/, "") + "/api/contact";

    let name = $state("");
    let email = $state("");
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
                body: JSON.stringify({ name: name.trim(), email: email.trim(), message: message.trim() }),
            });
            if (res.ok) {
                success = true;
                name = ""; email = ""; message = "";
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
    <div class="form-group">
        <input type="text" class="form-input" placeholder="Name" bind:value={name} disabled={submitting} />
    </div>
    <div class="form-group">
        <input type="email" class="form-input" placeholder="Email" bind:value={email} disabled={submitting} />
    </div>
    <div class="form-group">
        <textarea class="form-textarea" placeholder="Message" bind:value={message} rows="4" disabled={submitting}></textarea>
    </div>
    <button class="btn btn-primary form-submit" onclick={handleSubmit} disabled={submitting}>
        {submitting ? "Sending..." : "Send Message"}
    </button>
    {#if success}
        <div class="form-success">Message sent. I'll get back to you soon.</div>
    {/if}
    {#if errorMsg}
        <div class="form-error">{errorMsg}</div>
    {/if}
</div>

<style>
    .contact-form {
        padding: 1.5rem;
        display: flex;
        flex-direction: column;
        gap: 1rem;
    }
    .form-input,
    .form-textarea {
        width: 100%;
        background: #0a0a0a;
        border: 1px solid #222;
        color: var(--text-primary);
        padding: 0.75rem 1rem;
        font-family: var(--font-body);
        font-size: 0.85rem;
        outline: none;
        transition: border-color 0.2s;
        border-radius: 0;
    }
    .form-input:focus,
    .form-textarea:focus {
        border-color: var(--accent);
        box-shadow: 0 0 0 1px var(--accent);
    }
    .form-textarea {
        resize: vertical;
        min-height: 100px;
        background-image: repeating-linear-gradient(0deg, transparent, transparent 2px, rgba(255, 255, 255, 0.008) 2px, rgba(255, 255, 255, 0.008) 4px);
    }
    :global(body.light-mode) .form-input,
    :global(body.light-mode) .form-textarea {
        background: #f5f5f5;
        border-color: #ddd;
        color: #111;
    }
    .form-submit {
        margin-top: 0.5rem;
        width: 100%;
    }
    .form-success {
        color: #00ffaa;
        font-size: 0.8rem;
        padding: 0.5rem 0;
        animation: fadeSlideIn 0.4s ease-out;
    }
    .form-error {
        color: #ff4444;
        font-size: 0.8rem;
        padding: 0.5rem 0;
        animation: shake 0.3s ease-out;
    }
    :global(body.light-mode) .form-success {
        color: #00aa6e;
    }

    @keyframes fadeSlideIn {
        from { opacity: 0; transform: translateY(8px); }
        to { opacity: 1; transform: translateY(0); }
    }
    @keyframes shake {
        0%, 100% { transform: translateX(0); }
        25% { transform: translateX(-4px); }
        75% { transform: translateX(4px); }
    }

    @media (prefers-reduced-motion: reduce) {
        .form-success,
        .form-error,
        .form-input,
        .form-textarea {
            animation: none !important;
            transition: none !important;
        }
    }
</style>
