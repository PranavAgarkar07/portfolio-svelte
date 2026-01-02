<script lang="ts">
    import { theme } from "$lib/stores/theme";
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

    function toggleTheme() {
        theme.toggle();
    }

    onMount(() => {
        // Scroll Progress
        const handleScroll = () => {
            const winScroll =
                document.body.scrollTop || document.documentElement.scrollTop;
            const height =
                document.documentElement.scrollHeight -
                document.documentElement.clientHeight;
            scrollProgress = (winScroll / height) * 100;
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

<header class="aero-header">
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
                    [ SWITCH THEME ]
                </button>
            </li>
        </ul>

        <button
            id="theme-toggle"
            class="theme-btn desktop-only"
            onclick={toggleTheme}
            aria-label="Toggle Theme"
        >
            <svg
                class="sun-icon"
                xmlns="http://www.w3.org/2000/svg"
                width="24"
                height="24"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
            >
                <circle cx="12" cy="12" r="5"></circle>
                <path d="M12 1v2"></path>
                <path d="M12 21v2"></path>
                <path d="M4.22 4.22l1.42 1.42"></path>
                <path d="M18.36 18.36l1.42 1.42"></path>
                <path d="M1 12h2"></path>
                <path d="M21 12h2"></path>
                <path d="M4.22 19.78l1.42-1.42"></path>
                <path d="M18.36 5.64l1.42-1.42"></path>
            </svg>
            <svg
                class="moon-icon"
                xmlns="http://www.w3.org/2000/svg"
                width="24"
                height="24"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
            >
                <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"
                ></path>
            </svg>
        </button>

        <button
            class="mobile-menu-btn"
            onclick={() => (mobileMenuOpen = !mobileMenuOpen)}
            aria-label="Toggle Navigation"
        >
            <span class="btn-text"
                >{mobileMenuOpen ? "[ CLOSE ]" : "[ MENU ]"}</span
            >
        </button>
    </nav>

    <!-- Scroll Progress Signal -->
    <div class="signal-line">
        <div class="signal-track"></div>
        <div class="signal-fill" style="width: {scrollProgress}%"></div>
    </div>
</header>
