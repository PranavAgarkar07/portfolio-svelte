<script lang="ts">
    import { onMount } from "svelte";

    interface Props {
        profile: {
            name: string;
            role: string;
            tagline: string;
            status: string;
            avatar: string;
        };
        about: {
            bio: string;
            stats: Array<{ label: string; value: string }>;
        };
        skills: Array<{
            category: string;
            items: Array<{ name: string; icon: string; level: number }>;
        }>;
    }

    let { profile, about, skills }: Props = $props();

    let displayName = $state("");
    let displayRole = $state("");

    const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789@#$%";

    function decryptText(target: string, setter: (val: string) => void) {
        let iteration = 0;
        const interval = setInterval(() => {
            setter(
                target
                    .split("")
                    .map((char, i) => {
                        if (i < iteration) return char;
                        return chars[Math.floor(Math.random() * chars.length)];
                    })
                    .join(""),
            );

            iteration += 1 / 3;
            if (iteration >= target.length) clearInterval(interval);
        }, 30);
    }

    onMount(() => {
        decryptText(profile.name, (v) => (displayName = v));
        setTimeout(
            () => decryptText(profile.role, (v) => (displayRole = v)),
            300,
        );
    });
</script>

<section id="hero" class="hero-section snap-section">
    <div class="aero-card hero-content fade-in-up">
        <span
            class="badge"
            id="hero-status"
            class:available={profile.status.includes("Open")}
            >{profile.status}</span
        >
        <h1 id="hero-name">{displayName || profile.name}</h1>
        <h2 id="hero-role">{displayRole || profile.role}</h2>
        <p class="tagline" id="hero-tagline">{profile.tagline}</p>

        <div class="cta-group">
            <a href="#projects" class="btn btn-primary">View Work</a>
            <a href="#contact" class="btn btn-secondary">Contact Me</a>
        </div>
    </div>
</section>

<section id="about" class="section-container snap-section">
    <div class="section-header fade-in">
        <h2>About Me</h2>
    </div>
    <div class="aero-card content-card fade-in">
        <div class="visual-profile">
            <!-- Left Colon: Holographic Frame -->
            <div class="profile-frame">
                <img
                    src={profile.avatar}
                    alt="Pilot Avatar"
                    class="profile-img"
                />
                <div class="frame-corners"></div>
            </div>

            <!-- Right Column: Intel -->
            <div class="profile-intel">
                <p id="about-bio" class="bio-text">{about.bio}</p>

                <div class="mission-stats">
                    {#each about.stats as stat}
                        <div class="stat-item">
                            <span class="stat-label">{stat.label}</span>
                            <span class="stat-value">{stat.value}</span>
                        </div>
                    {/each}
                </div>
            </div>
        </div>
    </div>
</section>

<section id="skills" class="section-container snap-section">
    <div class="section-header fade-in">
        <h2>Technical Proficiency</h2>
    </div>
    <div class="skills-wrapper">
        {#each skills as category}
            <div class="aero-card skill-category-card fade-in">
                <h3 class="category-title">{category.category}</h3>
                <div class="skills-grid">
                    {#each category.items as skill}
                        <div class="skill-card">
                            <i class="{skill.icon} skill-icon"></i>
                            <span class="skill-name">{skill.name}</span>
                            <div class="skill-level-bar">
                                <div
                                    class="fill"
                                    style="width: {skill.level}%"
                                ></div>
                            </div>
                        </div>
                    {/each}
                </div>
            </div>
        {/each}
    </div>
</section>
