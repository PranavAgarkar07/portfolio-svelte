<script lang="ts">
    import { theme } from "$lib/stores/theme";
    import { scrollY } from "$lib/stores/scroll";
    import { onMount } from "svelte";
    import { base } from "$app/paths";
    import { goto } from "$app/navigation";
    import { page } from "$app/stores";
    import { trackEvent } from "$lib/analytics";

    function navigateToSection(id: string) {
        goto(`${base}/#${id}`);
    }

    interface NavigationItem {
        id: string;
        title: string;
    }

    const navigationItems: NavigationItem[] = [
        { id: "projects", title: "Projects" },
        { id: "skills", title: "Skills" },
        { id: "about", title: "About" },
        { id: "certifications", title: "Certifications" },
        { id: "contact", title: "Contact" },
    ];

    interface Props {
        profile: {
            name: string;
            socials: Array<{ name: string; url: string; label: string }>;
        };
    }

    let { profile }: Props = $props();
    let resumeUrl = $derived(
        profile.socials.find((s) => s.name === "Resume")?.url ||
            "/Pranav_Agarkar_Resume.pdf",
    );
    let mobileMenuOpen = $state(false);
    let activeSection = $state("");
    let scrollProgress = $state(0);
    let headerHidden = $state(false);
    let lastScrollY = $state(0);
    let reachedTop = $state(true);
    let detached = $state(false);
    let cachedDocHeight = $state(0);

    let isDark = $derived($theme === "dark");

    let onBlogPage = $derived($page.url.pathname.startsWith(`${base}/blog`));

    function toggleTheme() {
        theme.toggle();
    }

    function updateDocHeight() {
        cachedDocHeight =
            document.documentElement.scrollHeight -
            document.documentElement.clientHeight;
    }

    function openMenu() {
        mobileMenuOpen = true;
    }

    function closeMenu() {
        mobileMenuOpen = false;
    }

    $effect(() => {
        if (mobileMenuOpen) {
            document.body.style.overflow = "hidden";
            const handleEscape = (e: KeyboardEvent) => {
                if (e.key === "Escape") mobileMenuOpen = false;
            };
            window.addEventListener("keydown", handleEscape);
            return () => {
                document.body.style.overflow = "";
                window.removeEventListener("keydown", handleEscape);
            };
        } else {
            document.body.style.overflow = "";
        }
    });

    onMount(() => {
        updateDocHeight();
        window.addEventListener("resize", updateDocHeight);

        const handleScroll = (winScroll: number) => {
            scrollProgress =
                cachedDocHeight > 0 ? (winScroll / cachedDocHeight) * 100 : 0;
            reachedTop = winScroll < 10;

            const isDesktop = window.innerWidth > 768;

            if (isDesktop) {
                if (winScroll > 150) {
                    detached = true;
                } else if (winScroll < 80) {
                    detached = false;
                }
            } else {
                detached = false;
                const scrolledDown = winScroll > lastScrollY;
                if (scrolledDown && winScroll > 80) {
                    headerHidden = true;
                } else if (!scrolledDown) {
                    headerHidden = false;
                }
            }
            lastScrollY = winScroll;
        };

        const unsubScroll = scrollY.subscribe(handleScroll);

        const onNativeScroll = () => {
            handleScroll(window.scrollY);
        };
        window.addEventListener("scroll", onNativeScroll);

        const observer = new IntersectionObserver(
            (entries) => {
                entries.forEach((entry) => {
                    if (entry.isIntersecting) {
                        activeSection = entry.target.id;
                    }
                });
            },
            { threshold: 0.5 },
        );

        document.querySelectorAll("section").forEach((section) => {
            observer.observe(section);
        });

        return () => {
            unsubScroll();
            window.removeEventListener("scroll", onNativeScroll);
            window.removeEventListener("resize", updateDocHeight);
            observer.disconnect();
        };
    });
</script>

<div class="progress-bar" aria-hidden="true">
    <div class="progress-fill" style="width: {scrollProgress}%"></div>
</div>

<header
    class="aero-header"
    class:header-hidden={headerHidden}
    class:at-top={reachedTop}
    class:detached
>
    <nav class="nav-container">
        <a href="{base}/" class="logo">Pranav<span class="text-accent">.</span></a>

        <ul class="nav-links" class:active={mobileMenuOpen}>
            <button class="nav-close-btn" onclick={closeMenu} aria-label="Close menu">
                <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>
                </svg>
            </button>

            {#each navigationItems as item (item.id)}
                <li>
                    <a
                        href="{base}/#{item.id}"
                        class:active-section={activeSection === item.id}
                        onclick={(e) => { e.preventDefault(); navigateToSection(item.id); closeMenu(); }}>{item.title}</a
                    >
                </li>
            {/each}
            <li class="desktop-only nav-blog-li">
                <a href="{base}/blog" class="nav-blog-link" class:active-section={onBlogPage} onclick={closeMenu}>Blog</a>
            </li>
            <li class="nav-resume-li mobile-only">
                <button
                    class="nav-link-btn"
                    onclick={() => { toggleTheme(); closeMenu(); }}
                >
                    {#if isDark}
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
                            <circle cx="12" cy="12" r="5"></circle>
                            <line x1="12" y1="1" x2="12" y2="3"></line>
                            <line x1="12" y1="21" x2="12" y2="23"></line>
                            <line x1="4.22" y1="4.22" x2="5.64" y2="5.64"></line>
                            <line x1="18.36" y1="18.36" x2="19.78" y2="19.78"></line>
                            <line x1="1" y1="12" x2="3" y2="12"></line>
                            <line x1="21" y1="12" x2="23" y2="12"></line>
                            <line x1="4.22" y1="19.78" x2="5.64" y2="18.36"></line>
                            <line x1="18.36" y1="5.64" x2="19.78" y2="4.22"></line>
                        </svg>
                        Light
                    {:else}
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
                            <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"></path>
                        </svg>
                        Dark
                    {/if}
                </button>
            </li>
            <li class="nav-resume-li mobile-only">
                <a
                    href={resumeUrl}
                    target="_blank"
                    rel="noopener"
                    class="nav-resume-link"
                    onclick={() => {
                        closeMenu();
                        trackEvent("click", "resume_pdf");
                    }}
                >
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
                        <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
                        <polyline points="14 2 14 8 20 8"></polyline>
                        <line x1="16" y1="13" x2="8" y2="13"></line>
                        <line x1="16" y1="17" x2="8" y2="17"></line>
                    </svg>
                    Resume
                </a>
            </li>
        </ul>

        <div class="header-actions">
            <a
                href={resumeUrl}
                target="_blank"
                rel="noopener"
                class="desktop-only header-resume-btn"
                aria-label="Download Resume"
                onclick={() => trackEvent("click", "resume_pdf")}
            >
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
                    <polyline points="14 2 14 8 20 8"></polyline>
                    <line x1="16" y1="13" x2="8" y2="13"></line>
                    <line x1="16" y1="17" x2="8" y2="17"></line>
                </svg>
                Resume
            </a>
            <button
                class="theme-toggle desktop-only"
                onclick={toggleTheme}
                aria-label="Toggle Theme"
            >
                <span class="sun-icon">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <circle cx="12" cy="12" r="5"></circle>
                        <line x1="12" y1="1" x2="12" y2="3"></line>
                        <line x1="12" y1="21" x2="12" y2="23"></line>
                        <line x1="4.22" y1="4.22" x2="5.64" y2="5.64"></line>
                        <line x1="18.36" y1="18.36" x2="19.78" y2="19.78"></line>
                        <line x1="1" y1="12" x2="3" y2="12"></line>
                        <line x1="21" y1="12" x2="23" y2="12"></line>
                        <line x1="4.22" y1="19.78" x2="5.64" y2="18.36"></line>
                        <line x1="18.36" y1="5.64" x2="19.78" y2="4.22"></line>
                    </svg>
                </span>
                <span class="moon-icon">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"></path>
                    </svg>
                </span>
            </button>
        </div>

        <button
            class="mobile-menu-btn"
            class:active={mobileMenuOpen}
            onclick={openMenu}
            aria-label="Open Navigation"
        >
            <span class="hamburger">
                <span class="hamburger-line"></span>
                <span class="hamburger-line"></span>
                <span class="hamburger-line"></span>
            </span>
        </button>
    </nav>
</header>

<style>
    .nav-container {
        width: min(92%, 1200px);
        max-width: 1400px;
        display: flex;
        justify-content: space-between;
        align-items: center;
        gap: 1rem;
    }

    .logo {
        font-family: 'Clash Display', sans-serif;
        font-size: 1.4rem;
        font-weight: 600;
        color: #fff;
        letter-spacing: 0.02em;
        padding: 5px 0;
        transition: opacity 0.2s ease;
        flex-shrink: 0;
    }

    .text-accent {
        color: var(--accent);
    }

    :global(body.light-mode) .logo {
        color: var(--text-primary);
    }

    .nav-links {
        display: flex;
        gap: 2rem;
        list-style: none;
        align-items: center;
    }

    .nav-links a {
        font-family: var(--font-body);
        font-size: 0.85rem;
        text-transform: uppercase;
        color: var(--text-secondary);
        position: relative;
        letter-spacing: 0.04em;
        transition: color 0.2s ease;
        white-space: nowrap;
    }

    .nav-links a::before {
        content: '[';
        margin-right: 5px;
        opacity: 0;
        color: var(--accent);
        transition: opacity 0.2s ease, transform 0.2s ease;
        transform: translateX(4px);
        display: inline-block;
    }

    .nav-links a::after {
        content: ']';
        margin-left: 5px;
        opacity: 0;
        color: var(--accent);
        transition: opacity 0.2s ease, transform 0.2s ease;
        transform: translateX(-4px);
        display: inline-block;
    }

    .nav-links a:hover,
    .nav-links a.active-section {
        color: var(--accent);
    }

    .nav-links a:hover::before,
    .nav-links a:hover::after,
    .nav-links a.active-section::before,
    .nav-links a.active-section::after {
        opacity: 1;
        transform: translateX(0);
    }

    .nav-close-btn {
        display: none;
    }

    .header-actions {
        display: flex;
        align-items: center;
        gap: 0.35rem;
        margin-left: auto;
        flex-shrink: 0;
    }

    .nav-blog-link {
        color: var(--text-muted);
        font-size: 0.8rem;
        padding: 0.3rem 0.5rem;
        font-family: var(--font-body);
        text-transform: uppercase;
        letter-spacing: 0.04em;
        transition: color 0.2s ease;
    }
    .nav-blog-link:hover,
    .nav-blog-link.active-section {
        color: var(--accent);
    }

    .header-resume-btn {
        display: inline-flex;
        align-items: center;
        gap: 6px;
        font-family: var(--font-heading);
        font-size: 0.7rem;
        font-weight: 700;
        text-transform: uppercase;
        letter-spacing: 0.06em;
        color: #000 !important;
        background: var(--accent);
        padding: 7px 14px;
        text-decoration: none;
        transition: opacity 0.2s ease, transform 0.2s ease;
        border: none;
        cursor: pointer;
        flex-shrink: 0;
    }
    .header-resume-btn:hover {
        opacity: 0.85;
    }
    .header-resume-btn svg {
        flex-shrink: 0;
    }

    .theme-toggle {
        background: transparent;
        border: none;
        color: var(--text-secondary);
        cursor: pointer;
        padding: 0;
        display: flex;
        align-items: center;
        justify-content: center;
        width: 36px;
        height: 36px;
        position: relative;
        transition: color 0.2s ease, background 0.2s ease;
        flex-shrink: 0;
        border-radius: 6px;
    }
    .theme-toggle:hover {
        color: var(--accent);
        background: rgba(255, 68, 0, 0.06);
    }
    .theme-toggle:active {
        transform: scale(0.92);
    }
    .theme-toggle svg {
        width: 16px;
        height: 16px;
        stroke-width: 1.5px;
    }
    .theme-toggle .sun-icon,
    .theme-toggle .moon-icon {
        position: absolute;
        inset: 0;
        display: flex;
        align-items: center;
        justify-content: center;
        transition: opacity 0.3s ease, transform 0.3s ease;
    }
    .theme-toggle .sun-icon {
        opacity: 1;
        transform: rotate(0deg) scale(1);
    }
    .theme-toggle .moon-icon {
        opacity: 0;
        transform: rotate(90deg) scale(0.5);
    }
    :global(body.light-mode) .theme-toggle .sun-icon {
        opacity: 0;
        transform: rotate(-90deg) scale(0.5);
    }
    :global(body.light-mode) .theme-toggle .moon-icon {
        opacity: 1;
        transform: rotate(0deg) scale(1);
    }

    .mobile-menu-btn {
        display: none;
    }

    .mobile-only {
        display: none;
    }

    .desktop-only {
        display: flex;
    }

    @media (max-width: 768px) {
        .nav-links {
            position: fixed;
            top: 0;
            left: 0;
            width: 100vw;
            height: 100dvh;
            background: rgba(12, 14, 18, 0.97);
            flex-direction: column;
            justify-content: center;
            align-items: center;
            gap: 2.5rem;
            transform: translateY(-100%);
            transition: transform 0.4s cubic-bezier(0.16, 1, 0.3, 1);
            z-index: 1500;
            pointer-events: none;
            visibility: hidden;
            backdrop-filter: blur(24px);
            -webkit-backdrop-filter: blur(24px);
            padding: 2rem;
            margin: 0;
        }

        :global(body.light-mode) .nav-links {
            background: rgba(240, 242, 245, 0.97);
        }

        .nav-links::after {
            content: '';
            position: absolute;
            top: -30%;
            right: -20%;
            width: 300px;
            height: 300px;
            background: radial-gradient(circle, var(--accent) 0%, transparent 70%);
            opacity: 0.04;
            border-radius: 50%;
            pointer-events: none;
        }

        .nav-links.active {
            transform: translateY(0);
            pointer-events: auto;
            visibility: visible;
        }

        .nav-links.active::after {
            opacity: 0.08;
            transition: opacity 0.6s ease 0.2s;
        }

        .nav-links.active::before {
            content: '';
            position: absolute;
            top: 0;
            left: 0;
            width: 100%;
            height: 2px;
            background: linear-gradient(90deg, var(--accent) 0%, rgba(255,255,255,0.05) 80%, transparent 100%);
        }

        .nav-links.active li {
            width: 100%;
            text-align: center;
            opacity: 0;
            animation: navItemFade 0.5s cubic-bezier(0.16, 1, 0.3, 1) forwards;
            transform: translateY(30px);
        }

        .nav-links.active li:nth-child(1) { animation-delay: 0.05s; }
        .nav-links.active li:nth-child(2) { animation-delay: 0.10s; }
        .nav-links.active li:nth-child(3) { animation-delay: 0.15s; }
        .nav-links.active li:nth-child(4) { animation-delay: 0.20s; }
        .nav-links.active li:nth-child(5) { animation-delay: 0.25s; }
        .nav-links.active li:nth-child(6) { animation-delay: 0.30s; }
        .nav-links.active li:nth-child(7) { animation-delay: 0.35s; }

        .nav-links li {
            width: 100%;
            text-align: center;
            position: relative;
        }

        .nav-links a {
            font-size: 1.65rem;
            font-family: var(--font-heading);
            letter-spacing: 0.08em;
            transition: color 0.25s ease, opacity 0.25s ease;
            opacity: 0.85;
            position: relative;
            display: inline-block;
            color: var(--text-secondary);
        }

        .nav-links a::before {
            content: '[';
            margin-right: 10px;
            opacity: 0;
            color: var(--accent);
            transition: opacity 0.25s ease, transform 0.25s ease;
            transform: translateX(6px);
            display: inline-block;
        }

        .nav-links a::after {
            content: ']';
            margin-left: 10px;
            opacity: 0;
            color: var(--accent);
            transition: opacity 0.25s ease, transform 0.25s ease;
            transform: translateX(-6px);
            display: inline-block;
        }

        .nav-links a:hover,
        .nav-links a:focus-visible,
        .nav-links a.active-section {
            opacity: 1;
            color: var(--accent);
        }

        .nav-links a:hover::before,
        .nav-links a:hover::after,
        .nav-links a:focus-visible::before,
        .nav-links a:focus-visible::after,
        .nav-links a.active-section::before,
        .nav-links a.active-section::after {
            opacity: 1;
            transform: translateX(0);
        }

        .nav-close-btn {
            display: flex;
            position: absolute;
            top: 16px;
            right: 16px;
            z-index: 2001;
            background: none;
            border: none;
            color: var(--text-secondary);
            cursor: pointer;
            width: 44px;
            height: 44px;
            align-items: center;
            justify-content: center;
            opacity: 0;
            transition: opacity 0.3s ease, color 0.2s ease, background 0.2s ease;
            pointer-events: none;
            border-radius: 8px;
        }

        .nav-links.active .nav-close-btn {
            opacity: 0.7;
            pointer-events: auto;
            transition: opacity 0.3s ease 0.3s;
        }

        .nav-close-btn:hover {
            opacity: 1 !important;
            color: var(--accent);
            background: rgba(255,255,255,0.05);
        }

        .nav-link-btn {
            background: none;
            border: none;
            color: var(--text-secondary);
            font-family: var(--font-heading);
            font-size: 1.65rem;
            letter-spacing: 0.08em;
            text-transform: uppercase;
            padding: 6px 12px;
            cursor: pointer;
            width: auto;
            display: inline-flex;
            align-items: center;
            gap: 8px;
            transition: color 0.25s ease, opacity 0.25s ease;
            opacity: 0.85;
            border-radius: 4px;
        }

        .nav-link-btn:hover,
        .nav-link-btn:focus-visible {
            opacity: 1;
            color: var(--accent);
        }

        .nav-link-btn svg {
            width: 16px;
            height: 16px;
            flex-shrink: 0;
        }

        .nav-resume-li {
            display: flex;
            align-items: center;
            justify-content: center;
            width: 100%;
        }

        .nav-resume-li .nav-resume-link {
            font-size: 1.65rem;
            font-family: var(--font-heading);
            color: var(--text-secondary) !important;
            background: none;
            padding: 0;
            text-transform: uppercase;
            letter-spacing: 0.08em;
            border: none;
            transition: color 0.25s ease, opacity 0.25s ease;
            opacity: 0.85;
        }

        .nav-resume-li .nav-resume-link:hover,
        .nav-resume-li .nav-resume-link:focus-visible {
            opacity: 1;
            color: var(--accent) !important;
        }

        .mobile-menu-btn {
            display: flex;
            z-index: 2000;
            position: relative;
            background: none;
            border: none;
            color: var(--accent);
            cursor: pointer;
            padding: 8px;
            transition: opacity 0.2s ease;
            touch-action: manipulation;
            width: 36px;
            height: 36px;
            align-items: center;
            justify-content: center;
            flex-shrink: 0;
        }

        .mobile-menu-btn:hover {
            opacity: 0.7;
        }

        .hamburger {
            display: flex;
            flex-direction: column;
            justify-content: space-between;
            width: 22px;
            height: 16px;
            transition: transform 0.35s cubic-bezier(0.16, 1, 0.3, 1);
        }

        .mobile-menu-btn.active {
            opacity: 0;
            pointer-events: none;
        }

        .hamburger-line {
            display: block;
            width: 100%;
            height: 2px;
            border-radius: 1px;
            background: var(--accent);
            transition: transform 0.35s cubic-bezier(0.16, 1, 0.3, 1),
                        opacity 0.2s ease;
            transform-origin: center;
        }

        .mobile-only {
            display: block;
            width: 100%;
        }

        .desktop-only {
            display: none !important;
        }

        .header-actions {
            gap: 0.25rem;
        }
    }

    @keyframes navItemFade {
        to {
            opacity: 1;
            transform: translateY(0);
        }
    }

    @media (hover: none) {
        .nav-links a:hover,
        .nav-links a:focus-visible,
        .nav-links a:hover::before,
        .nav-links a:hover::after,
        .nav-links a:focus-visible::before,
        .nav-links a:focus-visible::after {
            color: var(--text-secondary);
            opacity: 0.85;
            transform: none;
        }
        .nav-links a.active-section,
        .nav-links a.active-section::before,
        .nav-links a.active-section::after {
            color: var(--accent);
            opacity: 1;
            transform: translateX(0);
        }
        .nav-blog-link:hover,
        .nav-blog-link:focus-visible {
            color: var(--text-muted);
        }
        .nav-blog-link.active-section {
            color: var(--accent);
        }
        .header-resume-btn:hover {
            opacity: 1;
        }
        .theme-toggle:hover {
            color: var(--text-secondary);
            background: transparent;
        }
    }
</style>
