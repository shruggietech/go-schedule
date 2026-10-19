import './global.css';
import type { Metadata } from 'next';
import { RootProvider } from 'fumadocs-ui/provider/next';

export const metadata: Metadata = {
  metadataBase: new URL('https://go-schedule.io'),
  title: { default: 'go-schedule | Put your workflows on schedule', template: '%s | go-schedule' },
  description: 'Readable schedules, event-driven workflows and a desktop control center. Open-source workload automation for Windows, Linux and macOS.',
  manifest: '/assets/favicons/site.webmanifest',
  icons: { icon: '/assets/favicons/favicon.svg', apple: '/assets/favicons/apple-touch-icon.png' },
  openGraph: { type: 'website', siteName: 'go-schedule', images: ['/assets/brand/go-schedule-social-preview-1280.png'] },
  twitter: { card: 'summary_large_image' },
};

export default function Layout({ children }: { children: React.ReactNode }) {
  return <html lang="en" suppressHydrationWarning><body><a className="skip-link" href="#main-content">Skip to content</a><RootProvider theme={{ defaultTheme: 'system', enableSystem: true }} search={{ options: { type: 'static', api: '/static.json' } }}>{children}</RootProvider></body></html>;
}
