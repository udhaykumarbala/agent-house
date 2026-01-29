import type { Metadata, Viewport } from 'next';
import { Providers } from './providers';
import './globals.css';

export const metadata: Metadata = {
  title: 'SubTrack - Subscription Manager',
  description: 'See all your subscriptions in one beautiful place. No bank connection needed.',
  manifest: '/manifest.json',
  appleWebApp: {
    capable: true,
    statusBarStyle: 'default',
    title: 'SubTrack',
  },
};

export const viewport: Viewport = {
  width: 'device-width',
  initialScale: 1,
  maximumScale: 1,
  userScalable: false,
  themeColor: [
    { media: '(prefers-color-scheme: light)', color: '#FAFAFA' },
    { media: '(prefers-color-scheme: dark)', color: '#0F0F12' },
  ],
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" suppressHydrationWarning>
      <head>
        <link
          rel="preload"
          href="https://cdn.jsdelivr.net/npm/@fontsource/geist-sans@5.0.0/files/geist-sans-latin-400-normal.woff2"
          as="font"
          type="font/woff2"
          crossOrigin="anonymous"
        />
        <style dangerouslySetInnerHTML={{ __html: `
          @font-face {
            font-family: 'Geist';
            font-style: normal;
            font-weight: 400;
            font-display: swap;
            src: url('https://cdn.jsdelivr.net/npm/@fontsource/geist-sans@5.0.0/files/geist-sans-latin-400-normal.woff2') format('woff2');
          }
          @font-face {
            font-family: 'Geist';
            font-style: normal;
            font-weight: 500;
            font-display: swap;
            src: url('https://cdn.jsdelivr.net/npm/@fontsource/geist-sans@5.0.0/files/geist-sans-latin-500-normal.woff2') format('woff2');
          }
          @font-face {
            font-family: 'Geist';
            font-style: normal;
            font-weight: 600;
            font-display: swap;
            src: url('https://cdn.jsdelivr.net/npm/@fontsource/geist-sans@5.0.0/files/geist-sans-latin-600-normal.woff2') format('woff2');
          }
          @font-face {
            font-family: 'Geist';
            font-style: normal;
            font-weight: 700;
            font-display: swap;
            src: url('https://cdn.jsdelivr.net/npm/@fontsource/geist-sans@5.0.0/files/geist-sans-latin-700-normal.woff2') format('woff2');
          }
        ` }} />
      </head>
      <body className="font-sans">
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}
