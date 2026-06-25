import avatar from "$lib/assets/avatar.webp";

// Portfolio Data
export const portfolioData = {
  profile: {
    name: "Pranav Agarkar",
    role: "Fullstack Developer",
    tagline:
      "Fullstack developer specializing in Django, React, and Go. I build secure SaaS backends, encrypted systems, and modern web UIs.",
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
        { name: "Python", icon: "devicon-python-plain", level: "Expert" },
        { name: "Django", icon: "devicon-django-plain", level: "Expert" },
        { name: "Linux", icon: "devicon-linux-plain", level: "Expert" },
        {
          name: "React",
          icon: "devicon-react-original colored",
          level: "Proficient",
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
        },
        { name: "Git", icon: "devicon-git-plain colored", level: "Proficient" },
        {
          name: "HTML5/CSS3",
          icon: "devicon-html5-plain colored",
          level: "Expert",
        },
        {
          name: "Rust",
          icon: "devicon-rust-plain",
          level: "Familiar",
        },
      ],
    },
    {
      category: "In Orbit",
      items: [
        {
          name: "Go",
          icon: "devicon-go-original-wordmark colored",
          level: "Proficient",
        },
        {
          name: "Svelte",
          icon: "devicon-svelte-plain colored",
          level: "Proficient",
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
      ],
    },
  ],
  projects: [
    {
      name: "TaskVault Lite",
      description:
        "Full-stack task manager with JWT/OAuth auth, AES-128-CBC encryption, and 1,247 req/min throughput. Built with Django REST Framework and React.",
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
      name: "BeamSync",
      description:
        "LAN file transfer with auto-IP discovery, QR pairing, and AES-256-GCM encryption. No cloud, no sign-up. Built with Go and Wails.",
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
      name: "KisanRakshak",
      description:
        "Predicts farmer distress using XGBoost on 6 agri signals with SHAP-explainable scores and real-time heatmaps. Built with React, Django, Go, and TimescaleDB.",
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
        "Android companion for wireless file transfer over local Wi-Fi. Scan a QR code, transfer at LAN speed with real-time progress. Built with Kotlin and Jetpack Compose.",
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
    image_url: "https://portfolio-uploads-sentinel.s3.ap-south-1.amazonaws.com/static/cert-nptl.png",
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
    credential_url: "https://www.redhat.com/en/services/training-and-certification",
    image_url: "https://portfolio-uploads-sentinel.s3.ap-south-1.amazonaws.com/static/cert-redhat-python.png",
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
    credential_url: "https://www.redhat.com/en/services/training-and-certification",
    image_url: "https://portfolio-uploads-sentinel.s3.ap-south-1.amazonaws.com/static/cert-redhat-linux.png",
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
    image_url: "https://portfolio-uploads-sentinel.s3.ap-south-1.amazonaws.com/static/cert-joomla.jpeg",
    tags: ["Joomla", "CMS", "Web Development"],
    is_verified: false,
    display_order: 4,
    created_at: "2024-04-01T00:00:00Z",
  },
];
