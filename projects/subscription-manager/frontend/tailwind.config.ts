import type { Config } from 'tailwindcss'

const config: Config = {
  darkMode: 'class',
  content: [
    './src/pages/**/*.{js,ts,jsx,tsx,mdx}',
    './src/components/**/*.{js,ts,jsx,tsx,mdx}',
    './src/app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        // Light mode colors (from UI spec)
        primary: {
          DEFAULT: '#FF7A5C',
          hover: '#FF6347',
          pressed: '#E85C3A',
        },
        background: '#FAFAFA',
        surface: {
          DEFAULT: '#FFFFFF',
          elevated: '#FFFFFF',
          hover: '#F5F5F5',
        },
        text: {
          primary: '#1A1A1E',
          secondary: '#6B6B73',
          tertiary: '#9CA3AF',
        },
        border: {
          DEFAULT: '#E5E5E8',
          focus: '#FF7A5C',
        },
        // Semantic colors
        success: {
          DEFAULT: '#059669',
          bg: '#ECFDF5',
        },
        error: {
          DEFAULT: '#DC2626',
          bg: '#FEF2F2',
        },
        warning: {
          DEFAULT: '#D97706',
          bg: '#FFFBEB',
        },
        info: {
          DEFAULT: '#0891B2',
          bg: '#ECFEFF',
        },
        // Dark mode specific
        dark: {
          primary: '#FF8A6C',
          'primary-hover': '#FF9A7C',
          'primary-pressed': '#FF7A5C',
          background: '#0F0F12',
          surface: '#1A1A1E',
          'surface-elevated': '#252529',
          'surface-hover': '#2A2A2F',
          'text-primary': '#FAFAFA',
          'text-secondary': '#A1A1AA',
          'text-tertiary': '#6B6B73',
          border: '#2A2A2F',
          success: '#10B981',
          error: '#EF4444',
          warning: '#F59E0B',
        },
      },
      fontFamily: {
        sans: ['Geist', 'system-ui', '-apple-system', 'BlinkMacSystemFont', 'SF Pro', 'Inter', 'sans-serif'],
      },
      fontSize: {
        display: ['40px', { lineHeight: '1.1', letterSpacing: '-0.02em', fontWeight: '700' }],
        h1: ['32px', { lineHeight: '1.2', letterSpacing: '-0.02em', fontWeight: '700' }],
        h2: ['24px', { lineHeight: '1.3', letterSpacing: '-0.01em', fontWeight: '600' }],
        h3: ['20px', { lineHeight: '1.4', letterSpacing: '-0.01em', fontWeight: '600' }],
        'body-lg': ['18px', { lineHeight: '1.5', fontWeight: '400' }],
        body: ['16px', { lineHeight: '1.5', fontWeight: '400' }],
        'body-sm': ['14px', { lineHeight: '1.5', fontWeight: '400' }],
        label: ['12px', { lineHeight: '1.4', letterSpacing: '0.02em', fontWeight: '500' }],
        caption: ['11px', { lineHeight: '1.4', letterSpacing: '0.02em', fontWeight: '400' }],
      },
      spacing: {
        '1': '4px',
        '2': '8px',
        '3': '12px',
        '4': '16px',
        '5': '20px',
        '6': '24px',
        '8': '32px',
        '10': '40px',
        '12': '48px',
        '16': '64px',
      },
      borderRadius: {
        none: '0px',
        sm: '4px',
        md: '8px',
        lg: '12px',
        xl: '16px',
        '2xl': '24px',
        full: '9999px',
      },
      boxShadow: {
        xs: '0 1px 2px rgba(0, 0, 0, 0.04)',
        sm: '0 2px 4px rgba(0, 0, 0, 0.06)',
        md: '0 4px 12px rgba(0, 0, 0, 0.08)',
        lg: '0 8px 24px rgba(0, 0, 0, 0.12)',
        xl: '0 12px 32px rgba(0, 0, 0, 0.16)',
        fab: '0 4px 16px rgba(255, 122, 92, 0.3)',
      },
      transitionDuration: {
        DEFAULT: '200ms',
        fast: '100ms',
        normal: '200ms',
        slow: '300ms',
      },
      animation: {
        'pulse-slow': 'pulse 1.5s ease-in-out infinite',
        'slide-up': 'slideUp 0.3s ease-out',
        'slide-down': 'slideDown 0.25s ease-in',
        'fade-in': 'fadeIn 0.3s ease-out',
      },
      keyframes: {
        slideUp: {
          '0%': { transform: 'translateY(100%)', opacity: '0' },
          '100%': { transform: 'translateY(0)', opacity: '1' },
        },
        slideDown: {
          '0%': { transform: 'translateY(0)', opacity: '1' },
          '100%': { transform: 'translateY(100%)', opacity: '0' },
        },
        fadeIn: {
          '0%': { opacity: '0', transform: 'translateY(10px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
      },
    },
  },
  plugins: [],
}
export default config
