import avatar from "$lib/assets/avatar.webp";

// Portfolio Data
export const portfolioData = {
  profile: {
    name: "Pranav Agarkar",
    role: "Fullstack Developer",
    tagline:
      "Building secure, high-throughput backends with Django and Go. 2+ years of full-stack experience turning complex problems into production-ready systems.",
    location: "Solapur, India",
    status: "Open to Opportunities",
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
        url: "./Pranav_Agarkar_Resume.pdf",
        label: "RESUME",
        icon: "file",
      },
    ],
    avatar: avatar, // Asset path
  },
  about: {
    bio: "2nd Year CSE student at Walchand Institute of Technology, Solapur. I build full-stack applications with Django and React/Svelte, focusing on clean architecture, secure data handling, and modern UI/UX. From encrypted task managers to offline P2P file transfer systems—I enjoy solving real problems with code.",
    stats: [
      { label: "EXPERIENCE", value: "2+ Years" },
      { label: "EDUCATION", value: "B.Tech CSE" },
      { label: "TECH STACK", value: "Full Stack" },
      { label: "STATUS", value: "Available" },
    ],
  },
  skills: [
    {
      category: "Core Stack",
      items: [
        {
          name: "Python",
          icon: "devicon-python-plain",
          level: "Expert",
          projects: ["TaskVault Lite", "KisanRakshak"],
        },
        {
          name: "Django",
          icon: "devicon-django-plain",
          level: "Expert",
          projects: ["TaskVault Lite", "KisanRakshak"],
        },
        { name: "Linux", icon: "devicon-linux-plain", level: "Expert" },
        {
          name: "React",
          icon: "devicon-react-original colored",
          level: "Proficient",
          projects: ["TaskVault Lite", "KisanRakshak"],
        },
      ],
    },
    {
      category: "Proficient",
      items: [
        { name: "C", icon: "devicon-c-plain colored", level: "Proficient" },
        {
          name: "JavaScript",
          icon: "devicon-javascript-plain colored",
          level: "Proficient",
          projects: ["TaskVault Lite", "KisanRakshak"],
        },
        { name: "Git", icon: "devicon-git-plain colored", level: "Proficient" },
        {
          name: "HTML5/CSS3",
          icon: "devicon-html5-plain colored",
          level: "Expert",
        },
      ],
    },
    {
      category: "Exploring",
      items: [
        {
          name: "Go",
          icon: "devicon-go-original-wordmark colored",
          level: "Proficient",
          projects: ["BeamSync", "KisanRakshak"],
        },
        {
          name: "Svelte",
          icon: "devicon-svelte-plain colored",
          level: "Proficient",
          projects: ["BeamSync", "Portfolio"],
        },
        {
          name: "Docker",
          icon: "devicon-docker-plain colored",
          level: "Familiar",
        },
        {
          name: "AWS",
          icon: "devicon-amazonwebservices-plain-wordmark colored",
          level: "Familiar",
        },
        {
          name: "K8s",
          icon: "devicon-kubernetes-plain colored",
          level: "Familiar",
        },
        {
          name: "Rust",
          icon: "devicon-rust-plain",
          level: "Familiar",
        },
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

export const fallbackCertificates = [
  {
    id: 1,
    title: "NPTEL — Privacy and Security in Online Social Media",
    issuer: "NPTEL (IIT Kharagpur)",
    date: "2024-10",
    credential_url: "https://archive.nptel.ac.in/noc/Ecertificate/",
    image_url:
      "https://portfolio-uploads-sentinel.s3.ap-south-1.amazonaws.com/static/cert-nptl.png",
    tags: ["Security", "Social Media", "Privacy"],
    is_verified: true,
    display_order: 1,
    created_at: "2024-10-01T00:00:00Z",
  },
  {
    id: 2,
    title: "Red Hat — Python Programming",
    issuer: "Red Hat Academy",
    date: "2024-08",
    credential_url:
      "https://www.redhat.com/en/services/training-and-certification",
    image_url:
      "https://portfolio-uploads-sentinel.s3.ap-south-1.amazonaws.com/static/cert-redhat-python.png",
    tags: ["Python", "Red Hat"],
    is_verified: true,
    display_order: 2,
    created_at: "2024-08-01T00:00:00Z",
  },
  {
    id: 3,
    title: "Red Hat — Linux Essentials",
    issuer: "Red Hat Academy",
    date: "2024-06",
    credential_url:
      "https://www.redhat.com/en/services/training-and-certification",
    image_url:
      "https://portfolio-uploads-sentinel.s3.ap-south-1.amazonaws.com/static/cert-redhat-linux.png",
    tags: ["Linux", "Red Hat"],
    is_verified: true,
    display_order: 3,
    created_at: "2024-06-01T00:00:00Z",
  },
  {
    id: 4,
    title: "Joomla! — Web Development Basics",
    issuer: "Joomla! Official",
    date: "2024-04",
    credential_url: "https://www.joomla.org/",
    image_url:
      "https://portfolio-uploads-sentinel.s3.ap-south-1.amazonaws.com/static/cert-joomla.jpeg",
    tags: ["Joomla", "CMS", "Web Development"],
    is_verified: false,
    display_order: 4,
    created_at: "2024-04-01T00:00:00Z",
  },
];
