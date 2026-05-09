import avatar from "$lib/assets/avatar.webp";

// Portfolio Data
export const portfolioData = {
  profile: {
    name: "Pranav Agarkar",
    role: "Fullstack Developer",
    tagline: "Building secure, scalable applications from frontend to backend.",
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
        { name: "Python", icon: "devicon-python-plain", level: 85 },
        { name: "Django", icon: "devicon-django-plain", level: 90 },
        { name: "Linux", icon: "devicon-linux-plain", level: 95 },
        { name: "React", icon: "devicon-react-original colored", level: 70 },
      ],
    },
    {
      category: "Proficient",
      items: [
        { name: "C", icon: "devicon-c-plain colored", level: 80 },
        {
          name: "JavaScript",
          icon: "devicon-javascript-plain colored",
          level: 85,
        },
        { name: "Git", icon: "devicon-git-plain colored", level: 75 },
        { name: "HTML5/CSS3", icon: "devicon-html5-plain colored", level: 95 },
      ],
    },
    {
      category: "In Orbit",
      items: [
        { name: "Go", icon: "devicon-go-original-wordmark colored", level: 65 },
        { name: "Svelte", icon: "devicon-svelte-plain colored", level: 55 },
        { name: "Docker", icon: "devicon-docker-plain colored", level: 50 },
      ],
    },
  ],
  projects: [
    {
      name: "TaskVault Lite",
      description:
        "Secure full-stack task manager with JWT authentication, Google/GitHub OAuth, and Fernet encryption for sensitive data. Built with Django REST Framework and React.",
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
        "High-performance offline P2P file transfer system with a cyberpunk terminal UI. Features auto-IP detection, dynamic port scouting, and QR-based mobile sync—no internet required.",
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
  ],
};
