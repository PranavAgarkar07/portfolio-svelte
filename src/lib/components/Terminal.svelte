<script lang="ts">
    import { onMount, tick } from "svelte";

    let inputValue = $state("");
    let outputLines = $state<Array<{ type: string; content: string }>>([]);
    let cursorVisible = $state(true);
    let inputElement: HTMLInputElement;
    let terminalBody: HTMLDivElement;

    let history = $state<string[]>([]);
    let historyIndex = $state(-1);
    let tempInput = $state("");

    const CONTACT_INFO = [
        "Email:    pranavagarkar8@gmail.com",
        "GitHub:   github.com/PranavAgarkar07",
        "LinkedIn: linkedin.com/in/pranavagarkar",
        "Resume:   ./Pranav_Agarkar_Resume.pdf",
    ];

    const commands: Record<string, () => string | string[] | void> = {
        help: () => [
            "Available commands: about, projects, skills, contact, clear, cat contact.txt, ls, email, sudo",
        ],
        about: () => {
            const el = document.getElementById("about");
            if (el) el.scrollIntoView({ behavior: "smooth" });
            return "Navigating to [ABOUT] section...";
        },
        projects: () => {
            const el = document.getElementById("projects");
            if (el) el.scrollIntoView({ behavior: "smooth" });
            return "Navigating to [PROJECTS] section...";
        },
        skills: () => {
            const el = document.getElementById("skills");
            if (el) el.scrollIntoView({ behavior: "smooth" });
            return "Navigating to [SKILLS] section...";
        },
        contact: () => {
            const el = document.getElementById("contact");
            if (el) el.scrollIntoView({ behavior: "smooth" });
            return "Navigating to [CONTACT] section...";
        },
        ls: () => "about  skills  projects  contact.txt",
        "cat contact.txt": () => CONTACT_INFO,
        sudo: () => "Permission denied: you are not Pranav 🚫",
        hi: () => openMail(),
        hello: () => openMail(),
        email: () => openMail(),
        clear: () => {
            outputLines = [
                {
                    type: "system",
                    content: `Last login: ${new Date().toLocaleString()} on ttys000`,
                },
            ];
        },
    };

    function openMail() {
        window.location.href = "mailto:pranavagarkar8@gmail.com";
        return "Opening mail client...";
    }

    function handleKeyDown(e: KeyboardEvent) {
        if (e.ctrlKey && e.key === "l") {
            e.preventDefault();
            commands.clear();
            return;
        }

        if (e.ctrlKey && e.key === "c") {
            e.preventDefault();
            outputLines = [
                ...outputLines,
                { type: "input", content: `➜ ~ ${inputValue}^C` },
            ];
            inputValue = "";
            historyIndex = -1;
            scrollToBottom();
            return;
        }

        if (e.key === "Tab") {
            e.preventDefault();
            const input = inputValue.toLowerCase();
            const matches = Object.keys(commands).filter((cmd) =>
                cmd.startsWith(input),
            );

            if (matches.length === 1) {
                inputValue = matches[0];
            } else if (matches.length > 1) {
                // Determine common prefix
                // Simpler: just show matches if not empty
                outputLines = [
                    ...outputLines,
                    { type: "input", content: `➜ ~ ${inputValue}` },
                ];
                outputLines = [
                    ...outputLines,
                    { type: "output", content: matches.join("  ") },
                ];
                scrollToBottom();
            }
            return;
        }

        if (e.key === "ArrowUp") {
            e.preventDefault();
            if (history.length > 0 && historyIndex < history.length - 1) {
                if (historyIndex === -1) tempInput = inputValue;
                historyIndex++;
                inputValue = history[history.length - 1 - historyIndex];
            }
        }

        if (e.key === "ArrowDown") {
            e.preventDefault();
            if (historyIndex > -1) {
                historyIndex--;
                if (historyIndex === -1) {
                    inputValue = tempInput;
                } else {
                    inputValue = history[history.length - 1 - historyIndex];
                }
            }
        }
    }

    async function handleCommand() {
        const cmd = inputValue.trim();
        const cmdKey = cmd.toLowerCase();

        if (cmd) {
            // Add to history if unique or different from last
            if (history.length === 0 || history[history.length - 1] !== cmd) {
                history.push(cmd);
            }
            historyIndex = -1;

            outputLines = [
                ...outputLines,
                { type: "input", content: `➜ ~ ${inputValue}` },
            ];

            if (commands[cmdKey]) {
                const result = commands[cmdKey]();
                if (Array.isArray(result)) {
                    result.forEach((line) => {
                        outputLines = [
                            ...outputLines,
                            { type: "output", content: line },
                        ];
                    });
                } else if (typeof result === "string") {
                    outputLines = [
                        ...outputLines,
                        { type: "output", content: result },
                    ];
                }
            } else {
                outputLines = [
                    ...outputLines,
                    {
                        type: "error",
                        content: `Command not found: ${cmd}. Type "help" for available commands.`,
                    },
                ];
            }
        }

        inputValue = "";
        await tick();
        scrollToBottom();
    }

    function scrollToBottom() {
        if (terminalBody) {
            terminalBody.scrollTop = terminalBody.scrollHeight;
        }
    }

    function focusInput() {
        inputElement?.focus();
    }

    onMount(() => {
        // Boot sequence
        setTimeout(
            () =>
                (outputLines = [
                    ...outputLines,
                    {
                        type: "system-dim",
                        content:
                            "> Initializing restricted shell environment...",
                    },
                ]),
            300,
        );
        setTimeout(
            () =>
                (outputLines = [
                    ...outputLines,
                    {
                        type: "system-dim",
                        content: "> Loading visual modules... [OK]",
                    },
                ]),
            800,
        );
        setTimeout(
            () =>
                (outputLines = [
                    ...outputLines,
                    {
                        type: "system-dim",
                        content:
                            "> Establishing secure connection... connected.",
                    },
                ]),
            1500,
        );
        setTimeout(
            () =>
                (outputLines = [
                    ...outputLines,
                    {
                        type: "system",
                        content:
                            'Type <span class="code-highlight">help</span> for a list of commands.',
                    },
                ]),
            2200,
        );

        // Cursor blink
        const interval = setInterval(() => {
            cursorVisible = !cursorVisible;
        }, 500);

        return () => clearInterval(interval);
    });
</script>

<!-- svelte-ignore a11y_click_events_have_key_events -->
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div class="terminal-wrapper" onclick={focusInput}>
    <div class="terminal-header">
        <span class="sys-status">[SYSTEM :: ONLINE]</span>
        <span class="sys-title">PRANAV_OS // KERNEL_V1.0</span>
    </div>

    <div class="terminal-body" bind:this={terminalBody}>
        {#each outputLines as line}
            <div
                class="output-line"
                class:error={line.type === "error"}
                class:dim={line.type === "system-dim"}
            >
                {#if line.type === "system"}
                    {@html line.content}
                {:else}
                    {line.content}
                {/if}
            </div>
        {/each}

        <div class="input-line">
            <span class="prompt">➜ <span class="path">~</span></span>
            <div class="input-wrapper">
                <span class="typed-text">{inputValue}</span>
                <span class="cursor-block" class:hidden={!cursorVisible}></span>
                <input
                    bind:this={inputElement}
                    bind:value={inputValue}
                    class="cmd-input"
                    type="text"
                    spellcheck="false"
                    autocomplete="off"
                    onkeydown={(e) => {
                        if (e.key === "Enter") handleCommand();
                        else handleKeyDown(e);
                    }}
                />
            </div>
        </div>
    </div>
</div>

<style>
    .hidden {
        opacity: 0;
    }
    .error {
        color: #ff4444;
    }
    .dim {
        color: #666;
    }
    .terminal-body {
        max-height: 300px;
        overflow-y: auto;
        padding-bottom: 2rem;
    }
    /* Ensure prompt colors match original */
    .prompt {
        color: var(--accent);
        margin-right: 8px;
        font-weight: bold;
    }
    .path {
        color: #00bcd4; /* Cyan/Blue path color */
    }
    :global(.code-highlight) {
        color: var(--accent);
        font-weight: bold;
    }
</style>
