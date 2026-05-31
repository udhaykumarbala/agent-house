import './style.css';
import { initTheme } from './ui/theme';
import { initApp } from './ui/app';

document.addEventListener('DOMContentLoaded', () => {
  initTheme();
  initApp();
});
