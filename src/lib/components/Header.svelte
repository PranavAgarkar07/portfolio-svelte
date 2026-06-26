<script lang="ts">
    import { theme } from "$lib/stores/theme";
    import { scrollY } from "$lib/stores/scroll";
    import { onMount } from "svelte";
    import { base } from "$app/paths";
    import { trackEvent } from "$lib/analytics";
    
    interface NavigationItem {
        id: string;
        title: string;
    }

    const navigationItems: NavigationItem[] = [
        { id: "about", title: "About" },
        { id: "skills", title: "Skills" },
        { id: "certifications", title: "Certifications" },
        { id: "projects", title: "Projects" },
        { id: "contact", title: "Contact" }
    ];

    interface Props {
        profile: {
            name: string;
            socials: Array<{ name: string; url: string; label: string }>;
        };
    }

    let { profile }: Props = $props();
    let resumeUrl = $derived(profile.socials.find(s => s.name === "Resume")?.url || "/Pranav_Agarkar_Resume.pdf");
    let mobileMenuOpen = $state(false);
    let activeSection = $state("");
    let scrollProgress = $state(0);
    let headerHidden = $state(false);
    let lastScrollY = $state(0);
    let reachedTop = $state(true);
    let detached = $state(false);
    let cachedDocHeight = $state(0);

    function toggleTheme() {
        theme.toggle();
    }

    function updateDocHeight() {
        cachedDocHeight = document.documentElement.scrollHeight - document.documentElement.clientHeight;
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
            scrollProgress = cachedDocHeight > 0 ? (winScroll / cachedDocHeight) * 100 : 0;
            reachedTop = winScroll < 10;

            const isDesktop = window.innerWidth > 768;

            // Hysteresis: prevent rapid toggling near threshold
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

        // Also keep native scroll as fallback
        const onNativeScroll = () => {
            handleScroll(window.scrollY);
        };
        window.addEventListener("scroll", onNativeScroll);

        // Active Section Observer
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
    class:detached={detached}
>
    <nav class="nav-container">
        <a href="{base}/" class="logo"
            >Pranav<span class="text-accent">.</span></a
        >
        <ul class="nav-links" class:active={mobileMenuOpen}>
            {#each navigationItems as item (item.id)}
                <li>
                    <a href="#{item.id}" class:active-section={activeSection === item.id} onclick={() => (mobileMenuOpen = false)}
                        >{item.title}</a
                    >
                </li>
            {/each}
            {#if theme === "dark"}
                <li class="nav-resume-li mobile-only">
                    <button
                        class="nav-link-btn"
                        onclick={() => {
                            toggleTheme();
                            mobileMenuOpen = false;
                        }}
                    >
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
                        Dark Mode
                    </button>
                </li>
            {:else}
                <li class="nav-resume-li mobile-only">
                    <button
                        class="nav-link-btn"
                        onclick={() => {
                            toggleTheme();
                            mobileMenuOpen = false;
                        }}
                    >
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
                        Light Mode
                    </button>
                </li>
            {/if}
            <li class="nav-resume-li mobile-only">
                <a href={resumeUrl} target="_blank" rel="noopener" class="nav-resume-link" onclick={() => { mobileMenuOpen = false; trackEvent("click", "resume_pdf"); }}>
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
            <a href={resumeUrl} target="_blank" rel="noopener" class="desktop-only header-resume-btn" aria-label="Download Resume" onclick={() => trackEvent("click", "resume_pdf")}>
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
            onclick={() => (mobileMenuOpen = !mobileMenuOpen)}
            aria-label="Toggle Navigation"
        >
            <span class="hamburger">
                <span class="hamburger-line"></span>
                <span class="hamburger-line"></span>
                <span class="hamburger-line"></span>
            </span>
        </button>
    </nav>
</header>
