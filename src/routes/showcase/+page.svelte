<script lang="ts">
	import { theme } from '$lib/stores/theme';
	import {
		Button, Card, Icon, Badge, Tag, SectionHeader, Skeleton,
		Input, Textarea, Select, DatePicker, Checkbox, Radio, RadioGroup,
		Magnetic, CustomCursor
	} from '$lib/components/ui';

	let enableAnimations = $state(false);
	let themeLabel = $derived($theme === 'dark' ? 'DARK' : 'LIGHT');

	let countryOptions = [
		{ value: 'us', label: 'United States' },
		{ value: 'ca', label: 'Canada' },
		{ value: 'uk', label: 'United Kingdom' },
		{ value: 'de', label: 'Germany' },
		{ value: 'fr', label: 'France' },
		{ value: 'jp', label: 'Japan' },
		{ value: 'in', label: 'India' },
		{ value: 'au', label: 'Australia' },
		{ value: 'br', label: 'Brazil' },
		{ value: 'sg', label: 'Singapore' },
	];

	let techOptions = [
		{ value: 'svelte', label: 'Svelte' },
		{ value: 'react', label: 'React' },
		{ value: 'vue', label: 'Vue' },
		{ value: 'angular', label: 'Angular' },
		{ value: 'solid', label: 'Solid.js' },
	];

	let selectedCountry = $state('');
	let selectedTech = $state('svelte');
	let selectedDate = $state('');

	let checkbox1 = $state(true);
	let checkbox2 = $state(false);
	let checkbox3 = $state(false);

	let radioValue = $state('option-a');
</script>

<svelte:head>
	<title>Component Library — Showcase</title>
</svelte:head>

<div class="showcase">
	<CustomCursor color="var(--accent)" />

	<header class="showcase-header">
		<div class="showcase-header-inner">
			<div>
				<h1>SHOWCASE</h1>
				<p class="showcase-subtitle">Component Library &middot; {themeLabel} MODE</p>
			</div>
			<div class="showcase-controls">
				<Checkbox bind:checked={enableAnimations} label="Animations" />
	
				<Button variant="ghost" size="sm" onclick={() => theme.toggle()}>
					<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="margin-right: 4px;">
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
					{themeLabel}
				</Button>
			</div>
		</div>
	</header>

	<section class="showcase-section" id="section-buttons">
		<SectionHeader title="Buttons" count={1} />

		<div class="demo-grid">
			<div class="demo-item">
				<span class="demo-label">Primary</span>
				<Magnetic><Button variant="primary" animate={enableAnimations}>Primary</Button></Magnetic>
			</div>
			<div class="demo-item">
				<span class="demo-label">Secondary</span>
				<Magnetic><Button variant="secondary" animate={enableAnimations}>Secondary</Button></Magnetic>
			</div>
			<div class="demo-item">
				<span class="demo-label">Ghost</span>
				<Magnetic><Button variant="ghost" animate={enableAnimations}>Ghost</Button></Magnetic>
			</div>
			<div class="demo-item">
				<span class="demo-label">Icon</span>
				<Magnetic><Button variant="icon" animate={enableAnimations} aria-label="Icon button"><Icon name="github" size={18} /></Button></Magnetic>
			</div>
			<div class="demo-item">
				<span class="demo-label">Disabled</span>
				<Magnetic><Button variant="primary" disabled>Disabled</Button></Magnetic>
			</div>
			<div class="demo-item">
				<span class="demo-label">As Link</span>
				<Magnetic><Button variant="primary" href="#section-buttons">Link</Button></Magnetic>
			</div>
			<div class="demo-item">
				<span class="demo-label">Small</span>
				<Magnetic><Button variant="primary" size="sm">Small</Button></Magnetic>
			</div>
			<div class="demo-item">
				<span class="demo-label">Large</span>
				<Magnetic><Button variant="primary" size="lg">Large</Button></Magnetic>
			</div>
		</div>
	</section>

	<section class="showcase-section" id="section-cards">
		<SectionHeader title="Cards" count={2} />

		<div class="demo-grid cols-2">
			<div class="demo-item">
				<span class="demo-label">Default Card</span>
					<Card badge={1}>
						<p style="font-family: 'Space Grotesk', sans-serif; font-size: 0.9rem; color: var(--text-secondary);">
							Default card with chamfered corners. This is the base aero-card style.
						</p>
					</Card>
			</div>
			<div class="demo-item">
				<span class="demo-label">Interactive Card</span>
					<Card interactive badge={2} animate={enableAnimations}>
						<p style="font-family: 'Space Grotesk', sans-serif; font-size: 0.9rem; color: var(--text-secondary);">
							Hover me! Interactive cards lift and highlight on hover with accent border.
						</p>
					</Card>
			</div>
		</div>
	</section>

	<section class="showcase-section" id="section-badges">
		<SectionHeader title="Badges" count={3} />

		<div class="demo-grid">
			<div class="demo-item">
				<span class="demo-label">Live</span>
				<Badge variant="live">LIVE</Badge>
			</div>
			<div class="demo-item">
				<span class="demo-label">Offline</span>
				<Badge variant="offline">OFFLINE</Badge>
			</div>
			<div class="demo-item">
				<span class="demo-label">Default</span>
				<Badge variant="default">STATUS</Badge>
			</div>
		</div>
	</section>

	<section class="showcase-section" id="section-tags">
		<SectionHeader title="Tags" count={4} />

		<div class="demo-grid">
			<div class="demo-item">
				<span class="demo-label">Default</span>
				<div class="demo-tags">
					<Tag>Svelte</Tag>
					<Tag>TypeScript</Tag>
					<Tag>SCSS</Tag>
				</div>
			</div>
			<div class="demo-item">
				<span class="demo-label">Muted</span>
				<div class="demo-tags">
					<Tag variant="muted">Rust</Tag>
					<Tag variant="muted">Go</Tag>
					<Tag variant="muted">Python</Tag>
				</div>
			</div>
			<div class="demo-item">
				<span class="demo-label">Removable</span>
				<div class="demo-tags">
					<Tag removable onclick={() => alert('Remove')}>Svelte</Tag>
					<Tag removable onclick={() => alert('Remove')}>SCSS</Tag>
				</div>
			</div>
		</div>
	</section>

	<section class="showcase-section" id="section-headers">
		<SectionHeader title="Section Headers" count={5} />

		<div class="demo-grid">
			<div class="demo-item">
				<span class="demo-label">Basic</span>
				<SectionHeader title="Projects" />
			</div>
			<div class="demo-item">
				<span class="demo-label">With Count</span>
				<SectionHeader title="Skills" count={8} />
			</div>
			<div class="demo-item">
				<span class="demo-label">Animated</span>
				<SectionHeader title="About" count={3} animate={enableAnimations} />
			</div>
		</div>
	</section>

	<section class="showcase-section" id="section-skeletons">
		<SectionHeader title="Skeletons" count={6} />

		<div class="demo-grid cols-2">
			<div class="demo-item">
				<span class="demo-label">Text (3 lines)</span>
				<Skeleton variant="text" count={3} />
			</div>
			<div class="demo-item">
				<span class="demo-label">Card</span>
				<Skeleton variant="card" />
			</div>
			<div class="demo-item">
				<span class="demo-label">Custom Width</span>
				<Skeleton variant="text" width="75%" />
			</div>
			<div class="demo-item">
				<span class="demo-label">Custom Height</span>
				<Skeleton variant="text" height="1rem" count={2} />
			</div>
		</div>
	</section>

	<section class="showcase-section" id="section-inputs">
		<SectionHeader title="Input" count={7} />

		<div class="demo-grid cols-2">
			<div class="demo-item">
				<span class="demo-label">Default</span>
				<Input label="Name" placeholder="Your Name" />
			</div>
			<div class="demo-item">
				<span class="demo-label">Email Type</span>
				<Input label="Email" type="email" placeholder="user@example.com" hint="We'll never share your email" />
			</div>
			<div class="demo-item">
				<span class="demo-label">Password (with toggle)</span>
				<Input label="Password" type="password" placeholder="Enter password" />
			</div>
			<div class="demo-item">
				<span class="demo-label">Search (clearable)</span>
				<Input label="Search" type="search" placeholder="Search..." clearable />
			</div>
			<div class="demo-item">
				<span class="demo-label">With Prefix Icon</span>
				<Input label="Username" placeholder="Your username" prefix="user" />
			</div>
			<div class="demo-item">
				<span class="demo-label">With Error</span>
				<Input label="Email" type="email" placeholder="Your Email" error="Invalid email format" />
			</div>
			<div class="demo-item">
				<span class="demo-label">Disabled</span>
				<Input label="Username" placeholder="Not editable" disabled />
			</div>
			<div class="demo-item">
				<span class="demo-label">Date</span>
				<Input label="Date" type="date" />
			</div>
		</div>
	</section>

	<section class="showcase-section" id="section-textareas">
		<SectionHeader title="Textarea" count={8} />

		<div class="demo-grid cols-2">
			<div class="demo-item">
				<span class="demo-label">Default</span>
				<Textarea label="Message" placeholder="Your Message" rows={3} />
			</div>
			<div class="demo-item">
				<span class="demo-label">With Error</span>
				<Textarea label="Message" placeholder="Your Message" error="Message is required" />
			</div>
			<div class="demo-item">
				<span class="demo-label">With Hint</span>
				<Textarea label="Feedback" placeholder="Share your thoughts" hint="We value your input" rows={3} />
			</div>
			<div class="demo-item">
				<span class="demo-label">Auto-grow + Char Count</span>
				<Textarea label="Bio" placeholder="Write about yourself..." autoGrow showCount maxlength={200} rows={2} />
			</div>
		</div>
	</section>

	<section class="showcase-section" id="section-select">
		<SectionHeader title="Select" count={9} />

		<div class="demo-grid cols-2">
			<div class="demo-item">
				<span class="demo-label">Default</span>
				<Select label="Country" options={countryOptions} />
			</div>
			<div class="demo-item">
				<span class="demo-label">With Default Value</span>
				<Select label="Technology" options={techOptions} value={selectedTech} />
			</div>
			<div class="demo-item">
				<span class="demo-label">Searchable</span>
				<Select label="Country (search)" options={countryOptions} searchable placeholder="Type to search..." />
			</div>
			<div class="demo-item">
				<span class="demo-label">With Error</span>
				<Select label="Country" options={countryOptions} error="Please select a country" />
			</div>
			<div class="demo-item">
				<span class="demo-label">Animated</span>
				<Select label="Technology" options={techOptions} animate={enableAnimations} searchable />
			</div>
			<div class="demo-item">
				<span class="demo-label">Disabled</span>
				<Select label="Country" options={countryOptions} disabled />
			</div>
		</div>
	</section>

	<section class="showcase-section" id="section-datepicker">
		<SectionHeader title="DatePicker" count={10} />

		<div class="demo-grid cols-2">
			<div class="demo-item">
				<span class="demo-label">Default</span>
				<DatePicker label="Pick a date" placeholder="Select a date" />
			</div>
			<div class="demo-item">
				<span class="demo-label">With Default Value</span>
				<DatePicker label="Appointment" bind:value={selectedDate} />
			</div>
			<div class="demo-item">
				<span class="demo-label">With Constraints</span>
				<DatePicker label="Travel Date" min="2026-05-01" max="2026-12-31" />
			</div>
			<div class="demo-item">
				<span class="demo-label">With Hint</span>
				<DatePicker label="Birthday" hint="We'll send you a surprise" />
			</div>
			<div class="demo-item">
				<span class="demo-label">With Error</span>
				<DatePicker label="Deadline" error="Date is required" />
			</div>
			<div class="demo-item">
				<span class="demo-label">Animated</span>
				<DatePicker label="Event Date" animate={enableAnimations} />
			</div>
		</div>
	</section>

	<section class="showcase-section" id="section-checkbox">
		<SectionHeader title="Checkbox" count={11} />

		<div class="demo-grid">
			<div class="demo-item">
				<span class="demo-label">Checked</span>
				<Checkbox label="Remember me" bind:checked={checkbox1} />
			</div>
			<div class="demo-item">
				<span class="demo-label">Unchecked</span>
				<Checkbox label="Subscribe to newsletter" bind:checked={checkbox2} />
			</div>
			<div class="demo-item">
				<span class="demo-label">Indeterminate</span>
				<Checkbox label="Select all" indeterminate />
			</div>
			<div class="demo-item">
				<span class="demo-label">Disabled (checked)</span>
				<Checkbox label="Read only" checked disabled />
			</div>
			<div class="demo-item">
				<span class="demo-label">Disabled</span>
				<Checkbox label="Not available" disabled />
			</div>
		</div>
	</section>

	<section class="showcase-section" id="section-radio">
		<SectionHeader title="Radio" count={12} />

		<div class="demo-grid cols-2">
			<div class="demo-item">
				<span class="demo-label">RadioGroup Default</span>
				<RadioGroup label="Choose option" bind:value={radioValue}>
					<Radio value="option-a" label="Option A" group="demo-1" />
					<Radio value="option-b" label="Option B" group="demo-1" />
					<Radio value="option-c" label="Option C" group="demo-1" />
				</RadioGroup>
			</div>
			<div class="demo-item">
				<span class="demo-label">With Disabled</span>
				<RadioGroup label="Delivery method">
					<Radio value="standard" label="Standard (5-7 days)" group="delivery" />
					<Radio value="express" label="Express (2-3 days)" group="delivery" />
					<Radio value="overnight" label="Overnight" group="delivery" disabled />
				</RadioGroup>
			</div>
		</div>
	</section>

	<section class="showcase-section" id="section-icons">
		<SectionHeader title="Icons" count={13} />

		<div class="demo-grid icon-gallery">
			<div class="demo-item icon-demo">
				<Icon name="github" size={24} />
				<span class="demo-label">github</span>
			</div>
			<div class="demo-item icon-demo">
				<Icon name="linkedin" size={24} />
				<span class="demo-label">linkedin</span>
			</div>
			<div class="demo-item icon-demo">
				<Icon name="envelope" size={24} />
				<span class="demo-label">envelope</span>
			</div>
			<div class="demo-item icon-demo">
				<Icon name="external-link" size={24} />
				<span class="demo-label">external-link</span>
			</div>
			<div class="demo-item icon-demo">
				<Icon name="chevron-left" size={24} />
				<span class="demo-label">chevron-left</span>
			</div>
			<div class="demo-item icon-demo">
				<Icon name="chevron-right" size={24} />
				<span class="demo-label">chevron-right</span>
			</div>
			<div class="demo-item icon-demo">
				<Icon name="close" size={24} />
				<span class="demo-label">close</span>
			</div>
			<div class="demo-item icon-demo">
				<Icon name="file" size={24} />
				<span class="demo-label">file</span>
			</div>
			<div class="demo-item icon-demo">
				<Icon name="arrow-up-right" size={24} />
				<span class="demo-label">arrow-up-right</span>
			</div>
			<div class="demo-item icon-demo">
				<Icon name="github" size={16} color="var(--accent)" />
				<span class="demo-label">size=16 accent</span>
			</div>
			<div class="demo-item icon-demo">
				<Icon name="linkedin" size={32} color="var(--text-secondary)" />
				<span class="demo-label">size=32 muted</span>
			</div>
		</div>
	</section>

	<section class="showcase-section" id="section-composition">
		<SectionHeader title="Composition Example" count={14} />

		<div class="composition-example">
				<Card interactive badge={3} animate={enableAnimations}>
					<div style="display: flex; align-items: center; gap: 1rem; margin-bottom: 1rem;">
						<Badge variant="live">LIVE</Badge>
						<Badge variant="default">PRJ 01</Badge>
				</div>
				<h3 style="font-family: 'Clash Display', sans-serif; font-size: 1.5rem; text-transform: uppercase; letter-spacing: 0.02em; margin-bottom: 0.5rem;">TaskVault Lite</h3>
				<p style="font-family: 'Space Grotesk', sans-serif; font-size: 0.875rem; color: var(--text-secondary); line-height: 1.65; margin-bottom: 1rem;">
					A lightweight task management application with real-time sync.
				</p>
				<div style="display: flex; flex-wrap: wrap; gap: 0.4rem; margin-bottom: 1.5rem;">
					{#each ['Svelte', 'TypeScript', 'Go', 'PostgreSQL'] as tag}
						<Tag>{tag}</Tag>
					{/each}
				</div>
				<div style="display: flex; gap: 0.75rem;">
					<Magnetic><Button variant="primary" size="sm" animate={enableAnimations}>View Project</Button></Magnetic>
					<Magnetic><Button variant="ghost" size="sm" animate={enableAnimations}>Source Code</Button></Magnetic>
				</div>
			</Card>
		</div>
	</section>
</div>

<style lang="scss">
	.showcase {
		max-width: 1100px;
		margin: 0 auto;
		padding: calc(var(--nav-height, 80px) + 3rem) 2rem 4rem;
	}

	.showcase-header {
		margin-bottom: 4rem;
		padding-bottom: 2rem;
		border-bottom: 2px solid var(--accent);
	}

	.showcase-header-inner {
		display: flex;
		justify-content: space-between;
		align-items: flex-end;
		gap: 2rem;
		flex-wrap: wrap;
	}

	.showcase-header h1 {
		font-family: 'Clash Display', sans-serif;
		font-size: 2rem;
		font-weight: 700;
		letter-spacing: 0.05em;
		text-transform: uppercase;
		color: var(--text-primary);
		margin: 0 0 0.25rem;
	}

	.showcase-subtitle {
		font-family: 'Space Grotesk', sans-serif;
		font-size: 0.8rem;
		color: var(--text-secondary);
		letter-spacing: 0.15em;
	}

	.showcase-controls {
		display: flex;
		align-items: center;
		gap: 1.5rem;
	}



	.theme-toggle-btn:hover {
		background: rgba(255, 68, 0, 0.15);
	}

	.showcase-section {
		margin-bottom: 4rem;
	}

	.demo-grid {
		display: flex;
		flex-wrap: wrap;
		gap: 1.5rem;
		align-items: flex-start;
		margin-top: 1.5rem;

		&.cols-2 {
			display: grid;
			grid-template-columns: repeat(2, 1fr);
		}

		&.icon-gallery {
			display: grid;
			grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
			gap: 1rem;
		}
	}

	.demo-item {
		display: flex;
		flex-direction: column;
		gap: 0.75rem;
	}

	.demo-label {
		font-family: 'Space Grotesk', sans-serif;
		font-size: 0.65rem;
		color: var(--text-secondary);
		text-transform: uppercase;
		letter-spacing: 0.15em;
		opacity: 0.7;
	}

	.demo-tags {
		display: flex;
		flex-wrap: wrap;
		gap: 0.4rem;
	}

	.icon-demo {
		align-items: center;
		text-align: center;
		padding: 1rem;
		border: 1px solid var(--grid-line);
		gap: 0.5rem;
	}

	.composition-example {
		max-width: 500px;
		margin-top: 1.5rem;
	}

	@media (max-width: 768px) {
		.showcase {
			padding: calc(var(--nav-height, 80px) + 2rem) 1rem 3rem;
		}

		.showcase-header-inner {
			flex-direction: column;
			align-items: flex-start;
		}

		.demo-grid.cols-2 {
			grid-template-columns: 1fr;
		}
	}


</style>
