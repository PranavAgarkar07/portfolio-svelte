import avatar from "$lib/assets/avatar.webp";

// Portfolio Data
export const portfolioData = {
    profile: {
        name: "Pranav Agarkar",
    role: "Fullstack Developer",
    tagline:
        "B.Tech CSE student building production-grade systems — P2P networking tools, encrypted APIs, and ML-powered dashboards. Seeking internships & early-career SDE roles.",
    location: "Solapur, India",
    status: "Seeking Internships & SDE Roles",
    socials: [
      {
        name: "GitHub",
        url: "https://github.com/PranavAgarkar07",
        label: "GH REPO",
        icon: "github",
      },
      {
        name: "LinkedIn",
        url: "https://www.linkedin.com/in/pranavagarkar",
        label: "LinkedIn",
        icon: "linkedin",
      },
      {
        name: "Email",
        url: "mailto:pranavagarkar8@gmail.com",
        label: "MAIL TO",
        icon: "envelope",
      },
      {
        name: "Resume",
        url: "/Pranav_Agarkar_Resume.pdf",
        label: "RESUME",
        icon: "file",
      },
    ],
    avatar: avatar, // Asset path
  },
  about: {
    bio: "3rd Year B.Tech CSE student (9.0 CGPA) at Walchand Institute of Technology, Solapur. I build production-grade systems with Go, Django, and modern frontend frameworks. From zero-knowledge encrypted task APIs to offline P2P file transfer with AES-256-GCM and auto-discovery — I engineer secure, real-world solutions. AICTE internship in Generative AI & Agentic Systems (2026). Open to internships and early-career SDE roles across India.",
    stats: [
      { label: "CGPA", value: "9.0 / 10" },
      { label: "GRADUATION", value: "May 2028" },
      { label: "FOCUS", value: "Systems + Web" },
      { label: "LOCATION", value: "Solapur / Remote" },
    ],
  },
  skills: [
    {
      category: "Languages",
      items: [
        { name: "Python", icon: "devicon-python-plain", level: "Expert", projects: ["TaskVault Lite", "KisanRakshak"] },
        { name: "Go", icon: "devicon-go-original-wordmark colored", level: "Proficient", projects: ["BeamSync", "KisanRakshak"] },
        { name: "JavaScript", icon: "devicon-javascript-plain colored", level: "Proficient", projects: ["TaskVault Lite", "KisanRakshak"] },
        { name: "C", icon: "devicon-c-plain colored", level: "Proficient" },
        { name: "Kotlin", icon: "devicon-kotlin-plain colored", level: "Familiar", projects: ["BeamSync Mobile"] },
        { name: "Bash", icon: "devicon-bash-plain", level: "Proficient" },
        { name: "SQL", icon: "devicon-postgresql-plain colored", level: "Proficient" },
      ],
    },
    {
      category: "Backend & Systems",
      items: [
        { name: "Django", icon: "devicon-django-plain", level: "Expert", projects: ["TaskVault Lite", "KisanRakshak"] },
        { name: "Django REST", icon: "devicon-django-plain", level: "Expert", projects: ["TaskVault Lite"] },
        { name: "Linux", icon: "devicon-linux-plain", level: "Expert" },
        { name: "Git", icon: "devicon-git-plain colored", level: "Proficient" },
        { name: "Docker", icon: "devicon-docker-plain colored", level: "Familiar" },
        { name: "AWS", icon: "devicon-amazonwebservices-plain-wordmark colored", level: "Familiar" },
        { name: "K8s", icon: "devicon-kubernetes-plain colored", level: "Familiar" },
      ],
    },
    {
      category: "Frontend & ML",
      items: [
        { name: "React", icon: "devicon-react-original colored", level: "Proficient", projects: ["TaskVault Lite", "KisanRakshak"] },
        { name: "Svelte", icon: "devicon-svelte-plain colored", level: "Proficient", projects: ["BeamSync", "Portfolio"] },
        { name: "HTML5/CSS3", icon: "devicon-html5-plain colored", level: "Expert" },
        { name: "XGBoost / ML", icon: "devicon-python-plain", level: "Familiar", projects: ["KisanRakshak"] },
        { name: "Rust", icon: "devicon-rust-plain", level: "Familiar" },
      ],
    },
  ],
  projects: [
    {
      name: "BeamSync",
      description:
        "Offline P2P file transfer tool with zero-config auto-discovery and QR pairing. Encrypts data with AES-256-GCM — no cloud, no accounts. Built with Go and Wails. 32 GitHub stars from developers who use it for air-gapped transfers.",
      featured: true,
      seriesTag: "BeamSync Ecosystem",
      githubStars: 32,
      images: [
        {
          src: "images/beamsync/appiconStartingpage.webp",
          alt: "BeamSync app icon and starting page",
        },
        {
          src: "images/beamsync/appSS1.webp",
          alt: "BeamSync file transfer interface",
        },
        {
          src: "images/beamsync/appSS2.webp",
          alt: "BeamSync peer discovery and connection",
        },
        {
          src: "images/beamsync/appSS3.webp",
          alt: "BeamSync QR sync and mobile pairing",
        },
      ],
      tags: ["Go", "Wails", "Svelte", "Networking", "P2P"],
      isLive: true,
      links: [
        {
          label: "Website",
          url: "https://pranavagarkar07.github.io/BeamSync/",
          icon: "external-link",
        },
        {
          label: "GitHub",
          url: "https://github.com/PranavAgarkar07/BeamSync",
          icon: "github",
        },
      ],
    },
    {
      name: "TaskVault Lite",
      description:
        "Encrypted task manager with JWT/OAuth authentication and AES-128-CBC encryption — 1,247 req/min throughput. Replaced manual tracking workflows with a secure API-first system. Built with Django REST Framework and React.",
      images: [
        {
          src: "images/taskvault/HomepageSS.webp",
          alt: "TaskVault Lite homepage",
        },
        {
          src: "images/taskvault/LoginPageSS.webp",
          alt: "TaskVault Lite login page with OAuth options",
        },
        {
          src: "images/taskvault/MainDashboardSS.webp",
          alt: "TaskVault Lite main dashboard with task list",
        },
        {
          src: "images/taskvault/secureDBwithEncrypetedTasks.webp",
          alt: "TaskVault Lite secure database with encrypted tasks",
        },
      ],
      tags: ["Django", "React", "JWT", "OAuth", "Security"],
      githubStars: 14,
      isLive: true,
      links: [
        {
          label: "Live Demo",
          url: "https://taskvault-lite.vercel.app/",
          icon: "external-link",
        },
        {
          label: "GitHub",
          url: "https://github.com/PranavAgarkar07/TaskVault-lite",
          icon: "github",
        },
      ],
    },
    {
      name: "KisanRakshak",
      description:
        "Predicts farmer distress risk across 10 districts using XGBoost on 6 agricultural signals. SHAP-explainable scores + real-time heatmap dashboards for early intervention. Built with React, Django, Go, and TimescaleDB.",
      githubStars: 8,
      images: [
        {
          src: "images/kisanrakshak/01-dashboard.webp",
          alt: "KisanRakshak CDI heatmap dashboard showing all 10 Vidarbha districts ranked by Crop Distress Index",
        },
        {
          src: "images/kisanrakshak/02-district-detail.webp",
          alt: "District detail view with SHAP waterfall chart explaining feature contributions to CDI score",
        },
        {
          src: "images/kisanrakshak/03-notifications.webp",
          alt: "Alert notification panel with PDF distress report generation and CDI threshold-based severity levels",
        },
        {
          src: "images/kisanrakshak/06-dashboard-full.webp",
          alt: "Full dashboard overview with CDI trend chart, district ranking table, and alert timeline",
        },
      ],
      tags: [
        "React",
        "TypeScript",
        "Django",
        "Python",
        "Go",
        "XGBoost",
        "AI/ML",
        "TimescaleDB",
      ],
      isLive: false,
      links: [
        {
          label: "GitHub",
          url: "https://github.com/PranavAgarkar07/KisanRakshak",
          icon: "github",
        },
      ],
    },
    {
      name: "BeamSync Mobile",
      description:
        "Android companion that extends the BeamSync ecosystem to mobile — scan a QR code to pair, transfer files at LAN speed with real-time progress tracking. Built with Kotlin, Jetpack Compose, and Material Design 3.",
      seriesTag: "BeamSync Ecosystem",
      githubStars: 6,
      images: [
        {
          src: "images/beamsync-mobile/permissions.webp",
          alt: "BeamSync Mobile permissions screen with camera and Wi-Fi status indicators",
        },
        {
          src: "images/beamsync-mobile/homescreen.webp",
          alt: "BeamSync Mobile home screen with Receive and Send action cards",
        },
        {
          src: "images/beamsync-mobile/uploadscreen.webp",
          alt: "BeamSync Mobile upload screen showing file selection and transfer progress",
        },
        {
          src: "images/beamsync-mobile/history.webp",
          alt: "BeamSync Mobile transfer history with timestamps and status indicators",
        },
      ],
      tags: [
        "Kotlin",
        "Jetpack Compose",
        "Android",
        "Material Design 3",
        "OkHttp",
        "CameraX",
      ],
      isLive: false,
      links: [
        {
          label: "GitHub",
          url: "https://github.com/PranavAgarkar07/BeamSyncMobile",
          icon: "github",
        },
      ],
    },
  ],
};
