/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{ts,html}'],
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: '#E8A838',
          hover: '#F0B94D',
          dim: '#C48A2A',
        },
        background: '#0F0F13',
        surface: {
          DEFAULT: '#1A1A24',
          light: '#24243A',
          border: '#2E2E48',
        },
        accent: {
          win: '#4ADE80',
          lose: '#F87171',
          draw: '#A78BFA',
        },
        text: {
          primary: '#F0F0F5',
          secondary: '#8888A0',
          muted: '#55556A',
        },
      },
      fontFamily: {
        display: ['"Space Grotesk"', 'system-ui', 'sans-serif'],
        body: ['Inter', 'system-ui', 'sans-serif'],
      },
      animation: {
        'bounce-in': 'bounceIn 0.5s cubic-bezier(0.34, 1.56, 0.64, 1)',
        'fade-up': 'fadeUp 0.4s ease-out',
        'shake': 'shake 0.5s ease-in-out',
        'pulse-glow': 'pulseGlow 2s ease-in-out infinite',
        'slide-down': 'slideDown 0.3s ease-out',
      },
      keyframes: {
        bounceIn: {
          '0%': { transform: 'scale(0)', opacity: '0' },
          '100%': { transform: 'scale(1)', opacity: '1' },
        },
        fadeUp: {
          '0%': { transform: 'translateY(20px)', opacity: '0' },
          '100%': { transform: 'translateY(0)', opacity: '1' },
        },
        shake: {
          '0%, 100%': { transform: 'translateX(0)' },
          '25%': { transform: 'translateX(-8px)' },
          '75%': { transform: 'translateX(8px)' },
        },
        pulseGlow: {
          '0%, 100%': { boxShadow: '0 0 20px rgba(232, 168, 56, 0.15)' },
          '50%': { boxShadow: '0 0 40px rgba(232, 168, 56, 0.3)' },
        },
        slideDown: {
          '0%': { transform: 'translateY(-10px)', opacity: '0' },
          '100%': { transform: 'translateY(0)', opacity: '1' },
        },
      },
    },
  },
  plugins: [],
}
