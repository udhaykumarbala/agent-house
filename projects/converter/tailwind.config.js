/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.ts'],
  theme: {
    extend: {
      colors: {
        canvas: '#131110',
        card: {
          DEFAULT: '#1E1B19',
          hover: '#2A2623',
        },
        edge: '#3D3632',
        copper: {
          DEFAULT: '#C17F59',
          light: '#D4956E',
          soft: 'rgba(193, 127, 89, 0.15)',
        },
        content: {
          DEFAULT: '#F2ECE6',
          muted: '#9B918A',
          faint: '#6B5F58',
        },
        ok: '#7BAE7F',
        danger: '#C4654A',
        caution: '#C7882A',
        tint: {
          length: '#D4A053',
          weight: '#C47A7A',
          temp: '#7A9EB2',
        },
      },
      fontFamily: {
        display: ["'Space Grotesk'", 'system-ui', 'monospace'],
        body: ["'Inter'", 'system-ui', 'sans-serif'],
      },
      borderRadius: {
        sm: '6px',
        md: '10px',
        lg: '16px',
        xl: '20px',
      },
      boxShadow: {
        sm: '0 1px 3px rgba(19, 17, 16, 0.12)',
        md: '0 4px 12px rgba(19, 17, 16, 0.16)',
        lg: '0 8px 24px rgba(19, 17, 16, 0.24)',
        glow: '0 0 20px rgba(193, 127, 89, 0.12)',
      },
    },
  },
  plugins: [],
};
