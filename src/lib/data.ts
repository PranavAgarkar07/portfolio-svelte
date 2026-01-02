import avatar from '$lib/assets/avatar.png';

// Portfolio Data
export const portfolioData = {
    profile: {
        name: "Pranav Agarkar",
        role: "Fullstack Developer",
        tagline: "Building secure, scalable applications from frontend to backend.",
        location: "Solapur, India",
        status: "Open to Opportunities",
        socials: [
            { name: "GitHub", url: "https://github.com/PranavAgarkar07", label: "GH_REPO", icon: "fab fa-github" },
            { name: "LinkedIn", url: "https://www.linkedin.com/in/pranavagarkar", label: "LINK_NET", icon: "fab fa-linkedin" },
            { name: "Email", url: "mailto:pranavagarkar8@gmail.com", label: "MAIL_TO", icon: "fas fa-envelope" },
            { name: "Resume", url: "./Pranav_Agarkar_Resume.pdf", label: "RESUME", icon: "fas fa-file-alt" }
        ],
        avatar: avatar // Asset path
    },
    about: {
        bio: "2nd Year CSE student at Walchand Institute of Technology, Solapur. I build full-stack applications with Django and React/Svelte, focusing on clean architecture, secure data handling, and modern UI/UX. From encrypted task managers to offline P2P file transfer systems—I enjoy solving real problems with code.",
        stats: [
            { label: "EXPERIENCE", value: "2+ Years" },
            { label: "EDUCATION", value: "B.Tech CSE" },
            { label: "TECH STACK", value: "Full Stack" },
            { label: "STATUS", value: "Available" }
        ]
    },
    skills: [
        { name: "Python", icon: "devicon-python-plain" },
        { name: "Django", icon: "devicon-django-plain" },
        { name: "React", icon: "devicon-react-original colored" },
        { name: "Svelte", icon: "devicon-svelte-plain colored" },
        { name: "Go", icon: "devicon-go-original-wordmark colored" },
        { name: "C", icon: "devicon-c-plain colored" },
        { name: "JavaScript", icon: "devicon-javascript-plain colored" },
        { name: "HTML5", icon: "devicon-html5-plain colored" },
        { name: "CSS3", icon: "devicon-css3-plain colored" },
        { name: "Docker", icon: "devicon-docker-plain colored" },
        { name: "Linux", icon: "devicon-linux-plain" },
        { name: "Git", icon: "devicon-git-plain colored" }
    ],
    projects: [
        {
            name: "TaskVault Lite",
            description: "Secure full-stack task manager with JWT authentication, Google/GitHub OAuth, and Fernet encryption for sensitive data. Built with Django REST Framework and React.",
            tags: ["Django", "React", "JWT", "OAuth", "Security"],
            isLive: true,
            links: [
                { label: "Live Demo", url: "https://taskvault-lite.vercel.app/", icon: "fas fa-external-link-alt" },
                { label: "GitHub", url: "https://github.com/PranavAgarkar07/TaskVault-lite", icon: "fab fa-github" }
            ]
        },
        {
            name: "BeamSync",
            description: "High-performance offline P2P file transfer system with a cyberpunk terminal UI. Features auto-IP detection, dynamic port scouting, and QR-based mobile sync—no internet required.",
            tags: ["Go", "Wails", "Svelte", "Networking", "P2P"],
            isLive: false,
            links: [
                { label: "GitHub", url: "https://github.com/PranavAgarkar07/BeamSync", icon: "fab fa-github" }
            ]
        }
    ]
};
