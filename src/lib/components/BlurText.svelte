<script lang="ts">
    import { onMount } from "svelte";

    interface Props {
        text: string;
        delay?: number;
        animateBy?: "words" | "letters";
        direction?: "top" | "bottom";
        class?: string;
    }

    let { text, delay = 50, animateBy = "words", direction = "top", class: className = "" }: Props = $props();

    let inView = $state(false);
    let el: HTMLParagraphElement | undefined = $state();

    onMount(() => {
        const observer = new IntersectionObserver(
            ([entry]) => {
                if (entry.isIntersecting) {
                    inView = true;
                    observer.disconnect();
                }
            },
            { threshold: 0.1 },
        );
        if (el) observer.observe(el);
        return () => observer.disconnect();
    });

    const segments = $derived(
        animateBy === "words" ? text.split(" ") : text.split(""),
    );

    const yOffset = $derived(direction === "top" ? "-30px" : "30px");
</script>

<p bind:this={el} class={className}>
    {#each segments as segment, i}
        <span
            class="blur-segment"
            style="transition: all 0.6s cubic-bezier(0.22, 1, 0.36, 1) {i * delay}ms; filter: blur({inView ? 0 : 10}px); opacity: {inView ? 1 : 0}; transform: translateY({inView ? 0 : yOffset});"
        >{segment}</span>{animateBy === "words" && i < segments.length - 1 ? " " : ""}
    {/each}
</p>

<style>
    .blur-segment {
        display: inline-block;
    }
</style>
