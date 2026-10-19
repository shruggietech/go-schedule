import { DocsLayout } from 'fumadocs-ui/layouts/docs';
import { source } from '@/lib/source';

export default function Layout({ children }: { children: React.ReactNode }) {
  return <DocsLayout tree={source.pageTree} nav={{ title: 'go-schedule', url: '/' }} githubUrl="https://github.com/shruggietech/go-schedule" links={[{ text: 'Compare', url: '/compare/' }, { text: 'Latest release', url: 'https://github.com/shruggietech/go-schedule/releases/latest', external: true }]}>{children}</DocsLayout>;
}
