<script lang="ts">
    interface Props {
        project: {
            name: string;
            description: string;
            tags: string[];
            isLive: boolean;
            links: Array<{ label: string; url: string; icon: string }>;
        };
        index: number;
    }

    let { project, index }: Props = $props();

    let projectId = $derived(`PRJ // 0${index + 1}`);
</script>

<article class="aero-card project-card">
    <div class="project-content">
        <div class="project-header">
            <div class="project-meta">
                <span class="project-id">{projectId}</span>
                {#if project.isLive}
                    <span class="live-indicator">
                        <span class="pulse-dot"></span>
                        LIVE
                    </span>
                {:else}
                    <span class="status-offline">OFFLINE</span>
                {/if}
            </div>
            <h3 class="title-highlight">{project.name}</h3>
        </div>

        <p class="project-description">{project.description}</p>

        <div class="project-tags">
            {#each project.tags as tag}
                <span class="tag">{tag}</span>
            {/each}
        </div>

        <div class="project-links">
            {#each project.links as link}
                <a
                    href={link.url}
                    target="_blank"
                    rel="noopener"
                    class="project-link"
                >
                    <i class={link.icon}></i>
                    <span class="link-text">{link.label}</span>
                    <i class="fas fa-external-link-alt link-arrow"></i>
                </a>
            {/each}
        </div>
    </div>
</article>

<style>
    .project-meta {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 0.5rem;
        border-bottom: 1px solid var(--grid-line);
        padding-bottom: 0.5rem;
    }

    .project-id {
        font-family: var(--font-body);
        font-size: 0.75rem;
        color: var(--text-secondary);
        letter-spacing: 2px;
    }

    .status-offline {
        font-size: 0.7rem;
        color: var(--text-secondary);
        opacity: 0.7;
    }

    .link-text {
        margin-right: 5px;
    }

    .link-arrow {
        font-size: 0.8em;
        transition: transform 0.3s ease;
    }

    .project-link:hover .link-arrow {
        transform: translate(3px, -3px);
    }
</style>
