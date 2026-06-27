<script lang="ts">
    import { onMount } from "svelte";
    import gsap from "gsap";
    import { ScrollTrigger } from "gsap/dist/ScrollTrigger";
    import { SectionHeader } from "$lib/components/ui";
    import DevLog from "./DevLog.svelte";

    let { skills }: { skills: Array<{ category: string; items: Array<{ name: string; icon: string; level: string; projects?: string[] }> }> } = $props();
    let prefersReducedMotion = $state(false);

    onMount(() => {
        prefersReducedMotion = window.matchMedia("(prefers-reduced-motion: reduce)").matches;
        if (prefersReducedMotion) return;
        gsap.registerPlugin(ScrollTrigger);
        gsap.defaults({ overwrite: "auto" });

        gsap.set(".skill-category-card", { y: 40, opacity: 0 });
        ScrollTrigger.create({
            trigger: "#skills", start: "top 85%", once: true,
            onEnter: () => { gsap.to(".skill-category-card", {
                y: 0, opacity: 1,
                stagger: 0.1,
                ease: "power3.out", duration: 0.6,
            });}
        });

        const skillCards = gsap.utils.toArray(".skill-card");
        if (skillCards.length > 0) {
            gsap.set(skillCards, { y: 20, opacity: 0 });
            ScrollTrigger.create({
                trigger: "#skills", start: "top 75%", once: true,
                onEnter: () => { gsap.to(skillCards, {
                    y: 0, opacity: 1,
                    stagger: 0.03,
                    ease: "power3.out", duration: 0.4,
                });}
            });
        }
    });
</script>

<section id="skills" class="section-container snap-section-content">
    <SectionHeader title="Technical Proficiency" animate />
    <div class="skills-wrapper">
        {#each skills as category}
            <div class="aero-card skill-category-card">
                <h3 class="category-title">{category.category}</h3>
                <div class="skills-grid">
                    {#each category.items as skill}
                        <div class="skill-card">
                            <i class="{skill.icon} skill-icon"></i>
                            <span class="skill-name">{skill.name}</span>
                            <span class="skill-level-label">{skill.level}</span>
                            {#if skill.projects && skill.projects.length > 0}
                                <span class="skill-projects">{skill.projects.slice(0, 2).join(" · ")}</span>
                            {/if}
                        </div>
                    {/each}
                </div>
            </div>
        {/each}
    </div>
    <div class="skills-footer">
        <DevLog />
    </div>
</section>
