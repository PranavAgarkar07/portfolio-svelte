<script lang="ts">
    import { onMount } from "svelte";
    import gsap from "gsap";
    import { ScrollTrigger } from "gsap/dist/ScrollTrigger";

    import { theme } from "$lib/stores/theme";
    import { portfolioData } from "$lib/data";
    import Header from "$lib/components/Header.svelte";
    import Hero from "$lib/components/Hero.svelte";
    import ProjectCard from "$lib/components/ProjectCard.svelte";
    import Terminal from "$lib/components/Terminal.svelte";
    import Footer from "$lib/components/Footer.svelte";

    const { profile, about, skills, projects } = portfolioData;

    onMount(() => {
        gsap.registerPlugin(ScrollTrigger);

        // Hero Timeline
        const tl = gsap.timeline();
        tl.from(".hero-content", {
            y: 50,
            opacity: 0,
            duration: 1,
            ease: "power3.out",
        }).from(
            ".cta-group .btn",
            {
                y: 20,
                opacity: 0,
                stagger: 0.2,
                duration: 0.8,
                ease: "power2.out",
            },
            "-=0.5",
        );

        // Featured Projects Stagger
        gsap.from(".project-card", {
            scrollTrigger: {
                trigger: "#projects",
                start: "top 80%",
            },
            scale: 0.98,
            opacity: 0,
            stagger: 0.1,
            duration: 0.8,
            ease: "power2.out",
        });

        // Skills Pop-in
        gsap.from(".skills-grid .skill-tag", {
            scrollTrigger: {
                trigger: "#skills",
                start: "top 85%",
            },
            scale: 0.8,
            opacity: 0,
            stagger: 0.05,
            duration: 0.5,
            ease: "back.out(1.7)",
        });

        // Section Headers Glitch-in (Simple Slide)
        gsap.utils.toArray(".section-header").forEach((header: any) => {
            gsap.from(header, {
                scrollTrigger: {
                    trigger: header,
                    start: "top 85%",
                },
                x: -50,
                opacity: 0,
                duration: 0.8,
                ease: "power2.out",
            });
        });

        // Contact Terminal
        gsap.from(".contact-card", {
            scrollTrigger: {
                trigger: "#contact",
                start: "top 80%",
            },
            y: 50,
            opacity: 0,
            duration: 1,
            ease: "power2.out",
        });
    });
</script>

<Header {profile} />

<main>
    <Hero {profile} {about} {skills} />

    <section id="projects" class="section-container">
        <div class="section-header fade-in">
            <h2>Featured Projects</h2>
        </div>
        <div class="projects-grid">
            {#each projects as project, i}
                <ProjectCard {project} index={i} />
            {/each}
        </div>
    </section>

    <section id="contact" class="section-container">
        <div class="aero-card contact-card fade-in">
            <div class="contact-header">
                <h2>Quick Connect Interface</h2>
                <p>
                    Type <span class="code-highlight">hi</span> or
                    <span class="code-highlight">email</span> to initiate protocol.
                </p>
            </div>

            <Terminal />

            <div class="quick-socials">
                {#each profile.socials as social}
                    <a
                        href={social.url}
                        target="_blank"
                        rel="noopener"
                        class="btn"
                    >
                        <i class={social.icon}></i>
                        {social.label}</a
                    >
                {/each}
            </div>
        </div>
    </section>
</main>

<Footer {profile} />
