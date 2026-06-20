import avatar from "$lib/assets/avatar.webp";

// Portfolio Data
export const portfolioData = {
  profile: {
    name: "Pranav Agarkar",
    role: "Fullstack Developer",
    tagline: "Fullstack developer specializing in Django, React, and Go. I build secure SaaS backends, encrypted systems, and modern web UIs.",
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
        { name: "React", icon: "devicon-react-original colored", level: "Proficient" },
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
        { name: "HTML5/CSS3", icon: "devicon-html5-plain colored", level: "Expert" },
      ],
    },
    {
      category: "In Orbit",
      items: [
        { name: "Go", icon: "devicon-go-original-wordmark colored", level: "Proficient" },
        { name: "Svelte", icon: "devicon-svelte-plain colored", level: "Proficient" },
        { name: "Docker", icon: "devicon-docker-plain colored", level: "Familiar" },
      ],
    },
  ],
  projects: [
    {
      name: "TaskVault Lite",
      description:
        "Secure full-stack task manager with JWT auth, OAuth, and AES-128-CBC encryption for sensitive data. Handles 1,247 req/min API throughput with PostgreSQL. Built with Django REST Framework and React.",
      brandIcon: "lock",
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
        "Offline P2P file transfer system with a cyberpunk terminal UI. Auto-IP discovery, dynamic port scouting, QR-based mobile pairing, and AES-256-GCM encrypted transfers — no internet or server required. Built with Go and Wails.",
      brandIcon: "zap",
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
        "Farmer Distress Early Warning System — computes a Crop Distress Index (CDI) across 10 Vidarbha districts with ensemble ML (Random Forest, XGBoost, Logistic Regression) and SHAP explainability. Features a real-time Leaflet heatmap, PDF distress reports, WebSocket alerts, 5-layer architecture (React + Django + Go + TimescaleDB + Bhashini API).",
      brandIcon: "leaf",
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
      tags: ["React", "TypeScript", "Django", "Python", "Go", "XGBoost", "AI/ML", "TimescaleDB"],
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
        "Wireless file transfers between Android and BeamSync desktop — no cables, no cloud, no sign-up. Scan a QR code to connect over Wi-Fi, then browse, upload, and download files at full LAN speed. Material Design 3, CameraX + ML Kit QR pairing, SAF file picker. Built with Kotlin and Jetpack Compose.",
      brandIcon: "smartphone",
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
      tags: ["Kotlin", "Jetpack Compose", "Android", "Material Design 3", "OkHttp", "CameraX"],
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
