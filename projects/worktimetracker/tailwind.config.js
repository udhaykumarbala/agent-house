/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{html,ts}'],
  theme: {
    extend: {
      spacing: {
        xs: '4px',
        sm: '8px',
        md: '16px',
        lg: '24px',
        xl: '32px',
        '2xl': '48px',
        '3xl': '64px',
      },
      colors: {
        primary: {
          DEFAULT: '#5746B2',
          hover: '#483A96',
          muted: 'rgba(87, 70, 178, 0.15)',
        },
        accent: {
          DEFAULT: '#CF8A2E',
          hover: '#B87724',
          muted: 'rgba(207, 138, 46, 0.12)',
        },
        background: '#0D0A1A',
        surface: {
          DEFAULT: '#17132B',
          raised: '#201B38',
        },
        divider: {
          DEFAULT: 'rgba(91, 70, 178, 0.15)',
          strong: 'rgba(91, 70, 178, 0.30)',
        },
        content: {
          DEFAULT: '#ECE9F5',
          secondary: '#8C83A8',
          disabled: '#564F6D',
        },
        success: '#3DB87A',
        error: '#D94B4B',
        warning: '#DFA23E',
        info: '#5B8AD4',
        chart: {
          1: '#CF8A2E',
          2: '#5746B2',
          3: '#E07B5F',
          4: '#2E9E8F',
          5: '#C75A8A',
          6: '#5B8AD4',
          7: '#8BA055',
          8: '#9B6DB0',
        },
      },
      fontFamily: {
        heading: ['Manrope', 'sans-serif'],
        body: ['Inter', 'sans-serif'],
      },
      fontSize: {
        timer: ['48px', { lineHeight: '1', fontWeight: '700', letterSpacing: '0.02em' }],
        'timer-compact': ['20px', { lineHeight: '1', fontWeight: '600', letterSpacing: '0.02em' }],
      },
      borderRadius: {
        sm: '4px',
        md: '8px',
        lg: '12px',
        xl: '16px',
      },
      boxShadow: {
        sm: '0 1px 2px rgba(13, 10, 26, 0.15)',
        md: '0 4px 8px rgba(13, 10, 26, 0.20)',
        lg: '0 12px 24px rgba(13, 10, 26, 0.30)',
        'glow-primary': '0 0 20px rgba(87, 70, 178, 0.25)',
        'glow-accent': '0 0 20px rgba(207, 138, 46, 0.25)',
      },
    },
  },
  plugins: [],
};
