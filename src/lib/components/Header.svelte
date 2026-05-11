<script lang="ts">
    import { theme } from "$lib/stores/theme";
    import { onMount } from "svelte";
    import { base } from "$app/paths";
    import { Button, Icon } from "$lib/components/ui";
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
        const handleScroll = () => {
            const winScroll =
                document.body.scrollTop || document.documentElement.scrollTop;
            const height =
                document.documentElement.scrollHeight -
                document.documentElement.clientHeight;
            scrollProgress = height > 0 ? (winScroll / height) * 100 : 0;
            reachedTop = winScroll < 10;

            // Desktop: smooth morph to floating past 120px
            // Mobile: hide/show on scroll direction
            const isDesktop = window.innerWidth > 768;
            detached = isDesktop && winScroll > 120;
            if (!isDesktop) {
                const scrolledDown = winScroll > lastScrollY;
                if (scrolledDown && winScroll > 80) {
                    headerHidden = true;
                } else if (!scrolledDown) {
                    headerHidden = false;
                }
            }
            lastScrollY = winScroll;
        };

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

        window.addEventListener("scroll", handleScroll);

        return () => {
            window.removeEventListener("scroll", handleScroll);
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

        <Button
            variant="icon"
            class="desktop-only"
            onclick={toggleTheme}
            aria-label="Toggle Theme"
        >
            <Icon name="sun" size={20} class="sun-icon" />
            <Icon name="moon" size={20} class="moon-icon" />
        </Button>

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
