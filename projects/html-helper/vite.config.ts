import { defineConfig } from 'vite';
import path from 'path';

export default defineConfig({
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  build: {
    target: 'es2020',
    minify: 'esbuild',
    rollupOptions: {
      output: {
        manualChunks: {
          'vendor': ['animejs', 'dompurify', 'localforage'],
          'editor': ['codemirror', '@codemirror/lang-html', '@codemirror/theme-one-dark'],
        },
      },
    },
  },
  server: {
    port: 3000,
    open: true,
  },
});
