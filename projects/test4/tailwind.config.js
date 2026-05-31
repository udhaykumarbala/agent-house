/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        // === Primary palette ===
        primary: {
          DEFAULT: 'var(--color-primary)',
          hover:   'var(--color-primary-hover)',
          light:   'var(--color-primary-light)',
        },
        accent: {
          DEFAULT: 'var(--color-accent)',
          light:   'var(--color-accent-light)',
        },

        // === Backgrounds ===
        background: 'var(--color-background)',
        surface:          'var(--color-surface)',
        'surface-hover':  'var(--color-surface-hover)',

        // === Borders ===
        border:       'var(--color-border)',
        'border-focus': 'var(--color-border-focus)',

        // === Text ===
        'text-primary':   'var(--color-text-primary)',
        'text-secondary': 'var(--color-text-secondary)',
        'text-tertiary':  'var(--color-text-tertiary)',

        // === Semantic ===
        success:       'var(--color-success)',
        'success-light': 'var(--color-success-light)',
        error:         'var(--color-error)',
        warning:       'var(--color-warning)',

        // === Category colors (auto-switched via CSS vars) ===
        category: {
          research:  'var(--color-category-research)',
          design:    'var(--color-category-design)',
          testing:   'var(--color-category-testing)',
          review:    'var(--color-category-review)',
          discovery: 'var(--color-category-discovery)',
          handoff:   'var(--color-category-handoff)',
          admin:     'var(--color-category-admin)',
        },
      },

      fontFamily: {
        heading: ['"Cabinet Grotesk"', 'Inter', 'system-ui', '-apple-system', 'sans-serif'],
        body:    ['Inter', 'system-ui', '-apple-system', 'sans-serif'],
        mono:    ['"JetBrains Mono"', '"SF Mono"', 'Consolas', 'monospace'],
      },

      fontSize: {
        'chip':     ['11px',  { lineHeight: '1.2', letterSpacing: '0.02em', fontWeight: '500' }],
        'meta':     ['12px',  { lineHeight: '1.4' }],
        'body':     ['14px',  { lineHeight: '1.5' }],
        'task':     ['15px',  { lineHeight: '1.4', fontWeight: '500' }],
        'section':  ['14px',  { lineHeight: '1.4', fontWeight: '600', letterSpacing: '0' }],
        'h3':       ['14px',  { lineHeight: '1.4', fontWeight: '600' }],
        'h2':       ['20px',  { lineHeight: '1.3', fontWeight: '600', letterSpacing: '-0.01em' }],
        'h1':       ['28px',  { lineHeight: '1.2', fontWeight: '700', letterSpacing: '-0.02em' }],
        'logo':     ['20px',  { lineHeight: '1.2', fontWeight: '700', letterSpacing: '-0.02em' }],
        'estimate': ['12px',  { lineHeight: '1.4', fontFamily: 'var(--font-mono)' }],
      },

      spacing: {
        '0.5': '2px',
        '1':   '4px',
        '2':   '8px',
        '3':   '12px',
        '4':   '16px',
        '5':   '20px',
        '6':   '24px',
        '8':   '32px',
        '10':  '40px',
        '12':  '48px',
        '16':  '64px',
      },

      borderRadius: {
        'none':   '0px',
        'sm':     '4px',
        'DEFAULT': '8px',
        'md':     '8px',
        'lg':     '12px',
        'xl':     '16px',
        'full':   '9999px',
      },

      boxShadow: {
        // Warm-tinted shadows — auto-switched via CSS vars
        'warm-sm':  'var(--shadow-sm)',
        'warm-md':  'var(--shadow-md)',
        'warm-lg':  'var(--shadow-lg)',
        'warm-xl':  'var(--shadow-xl)',
      },

      transitionDuration: {
        'fast': '150ms',
        'base': '200ms',
        'slow': '300ms',
      },

      transitionTimingFunction: {
        'spring': 'cubic-bezier(0.34, 1.56, 0.64, 1)',
        'smooth': 'cubic-bezier(0.4, 0, 0.2, 1)',
      },

      animation: {
        'spring':    'spring 200ms cubic-bezier(0.34, 1.56, 0.64, 1)',
        'slide-in':  'slideIn 250ms cubic-bezier(0.4, 0, 0.2, 1)',
        'fade-in':   'fadeIn 200ms ease-out',
      },

      keyframes: {
        spring: {
          '0%':   { transform: 'scale(0.95)' },
          '100%': { transform: 'scale(1)' },
        },
        slideIn: {
          '0%':   { transform: 'translateX(100%)' },
          '100%': { transform: 'translateX(0)' },
        },
        fadeIn: {
          '0%':   { opacity: '0' },
          '100%': { opacity: '1' },
        },
      },
    },
  },
  plugins: [],
}
