<script lang="ts">
    import { theme } from "$lib/stores/theme";
    import { scrollY } from "$lib/stores/scroll";
    import { onMount } from "svelte";
    import { base } from "$app/paths";
    interface Props {
        profile: {
            name: string;
            socials: Array<{ name: string; url: string; label: string }>;
        };
    }

    let { profile }: Props = $props();
    let mobileMenuOpen = $state(false);
    let activeSection = $state("");
    let scrollProgress = $state(0);
    let headerHidden = $state(false);
    let lastScrollY = $state(0);
    let reachedTop = $state(true);
    let detached = $state(false);

    function toggleTheme() {
        theme.toggle();
    }

    onMount(() => {
        const handleScroll = (winScroll: number) => {
            const height =
                document.documentElement.scrollHeight -
                document.documentElement.clientHeight;
            scrollProgress = height > 0 ? (winScroll / height) * 100 : 0;
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
            <li>
                <a href="#about" onclick={() => (mobileMenuOpen = false)}
                    >About</a
                >
            </li>
            <li>
                <a href="#skills" onclick={() => (mobileMenuOpen = false)}
                    >Skills</a
                >
            </li>
            <li>
                <a
                    href="#projects"
                    class:active={activeSection === "projects"}
                    onclick={() => (mobileMenuOpen = false)}>Projects</a
                >
            </li>
            <li>
                <a
                    href="#contact"
                    class:active={activeSection === "contact"}
                    onclick={() => (mobileMenuOpen = false)}>Contact</a
                >
            </li>

            <!-- Mobile Only Theme Switcher -->
            <li class="mobile-only">
                <button
                    class="nav-link-btn"
                    onclick={() => {
                        toggleTheme();
                        mobileMenuOpen = false;
                    }}
                >
                    THEME
                </button>
            </li>
        </ul>

        <button
            class="desktop-only theme-toggle"
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
