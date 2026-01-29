/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: '#FF2E63',
          hover: '#E6175A',
          light: '#FF5C85',
        },
        secondary: {
          DEFAULT: '#00D9FF',
          hover: '#00BFE6',
          light: '#33E2FF',
        },
        background: {
          dark: '#1A1A2E',
          light: '#EEEEF7',
        },
        surface: {
          dark: '#16213E',
          elevated: '#1F2D52',
          light: '#FFFFFF',
        },
        border: {
          dark: '#2A3B62',
          light: '#D4D4E8',
        },
        text: {
          'primary-dark': '#FFFFFF',
          'secondary-dark': '#B8B8D1',
          'tertiary-dark': '#8585A3',
          'primary-light': '#1A1A2E',
          'secondary-light': '#5F5F7A',
        },
        success: {
          DEFAULT: '#00FF9D',
          dark: '#00CC7E',
        },
        error: {
          DEFAULT: '#FF3366',
          dark: '#E61A47',
        },
        warning: '#FFB627',
      },
      fontFamily: {
        display: ['Space Grotesk', 'Inter', 'Segoe UI', 'system-ui', 'sans-serif'],
        body: ['Inter', 'Segoe UI', 'system-ui', '-apple-system', 'sans-serif'],
        code: ['JetBrains Mono', 'Fira Code', 'Monaco', 'Courier New', 'monospace'],
      },
      spacing: {
        'xs': '4px',
        'sm': '8px',
        'md': '16px',
        'lg': '24px',
        'xl': '32px',
        '2xl': '48px',
        '3xl': '64px',
        '4xl': '96px',
      },
      borderRadius: {
        'sm': '6px',
        'md': '8px',
        'lg': '12px',
        'xl': '16px',
        '2xl': '24px',
      },
      boxShadow: {
        'sm': '0 1px 3px rgba(0, 0, 0, 0.12), 0 1px 2px rgba(0, 0, 0, 0.08)',
        'md': '0 4px 6px rgba(0, 0, 0, 0.12), 0 2px 4px rgba(0, 0, 0, 0.08)',
        'lg': '0 10px 15px rgba(0, 0, 0, 0.15), 0 4px 6px rgba(0, 0, 0, 0.08)',
        'xl': '0 20px 25px rgba(0, 0, 0, 0.15), 0 10px 10px rgba(0, 0, 0, 0.04)',
        'glow-primary': '0 0 20px rgba(255, 46, 99, 0.4), 0 0 40px rgba(255, 46, 99, 0.2)',
        'glow-secondary': '0 0 20px rgba(0, 217, 255, 0.4), 0 0 40px rgba(0, 217, 255, 0.2)',
        'glow-success': '0 0 20px rgba(0, 255, 157, 0.5)',
        'glow-error': '0 0 20px rgba(255, 51, 102, 0.5)',
      },
      animation: {
        'bounce-in': 'bounceIn 0.6s cubic-bezier(0.68, -0.55, 0.265, 1.55)',
        'shake': 'shake 0.4s',
        'pulse-glow': 'pulseGlow 2s infinite',
        'float': 'float 3s ease-in-out infinite',
        'confetti-fall': 'confettiFall 2s ease-in-out forwards',
      },
      keyframes: {
        bounceIn: {
          '0%': { transform: 'scale(0.3)', opacity: '0' },
          '50%': { transform: 'scale(1.05)' },
          '70%': { transform: 'scale(0.9)' },
          '100%': { transform: 'scale(1)', opacity: '1' },
        },
        shake: {
          '0%, 100%': { transform: 'translateX(0)' },
          '10%, 30%, 50%, 70%, 90%': { transform: 'translateX(-4px)' },
          '20%, 40%, 60%, 80%': { transform: 'translateX(4px)' },
        },
        pulseGlow: {
          '0%, 100%': { boxShadow: '0 0 10px rgba(255, 46, 99, 0.5)' },
          '50%': { boxShadow: '0 0 20px rgba(255, 46, 99, 0.8), 0 0 30px rgba(255, 46, 99, 0.6)' },
        },
        float: {
          '0%, 100%': { transform: 'translateY(0px)' },
          '50%': { transform: 'translateY(-10px)' },
        },
        confettiFall: {
          '0%': { transform: 'translateY(-100%) rotate(0deg)', opacity: '1' },
          '100%': { transform: 'translateY(100vh) rotate(360deg)', opacity: '0' },
        },
      },
    },
  },
  plugins: [],
}
