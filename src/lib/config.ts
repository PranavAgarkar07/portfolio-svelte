const BASE = (import.meta.env.VITE_API_URL || '').replace(/\/$/, '');
const PROD = import.meta.env.PROD;

export const SITE = {
  url: PROD
    ? 'https://pranavagarkar07.github.io/portfolio-svelte'
    : 'http://localhost:5173/portfolio-svelte',
  name: 'Pranav Agarkar',
  title: 'Pranav Agarkar | Fullstack Developer',
  description: 'Fullstack Developer specializing in Django, React, Svelte, and Go. Building fast, modern web apps.',
  author: 'Pranav Agarkar',
  twitter: '@PranavAgarkar18',
  image: PROD
    ? 'https://pranavagarkar07.github.io/portfolio-svelte/og-image.png'
    : 'http://localhost:5173/portfolio-svelte/og-image.png',
  apiBase: BASE,
};

export const STATIC_ROUTES = ['', '/about', '/projects', '/blog', '/contact'];
