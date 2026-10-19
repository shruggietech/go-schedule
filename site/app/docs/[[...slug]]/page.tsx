import { source } from '@/lib/source';
import { DocsPage, DocsBody } from 'fumadocs-ui/page';
import { notFound } from 'next/navigation';
import { getMDXComponents } from '@/mdx-components';
import routes from '@/.source/routes.json';
import { Footer } from '@/components/chrome';

export default async function Page({ params }: { params: Promise<{ slug?: string[] }> }) {
  const { slug } = await params;
  const page = source.getPage(slug);
  if (!page) notFound();
  const Body = page.data.body;
  const original = routes.find(p => p.route.replace(/\/$/, '') === page.url.replace(/\/$/, ''));
  return <DocsPage id="main-content" tabIndex={-1} toc={page.data.toc}><DocsBody><Body components={getMDXComponents()} /></DocsBody>{original && <a className="edit-link" href={`https://github.com/shruggietech/go-schedule/edit/main/docs/${original.source}`}>Edit this page on GitHub</a>}<Footer /></DocsPage>;
}

export const generateStaticParams = () => source.generateParams();

export async function generateMetadata({ params }: { params: Promise<{ slug?: string[] }> }) {
  const { slug } = await params;
  const page = source.getPage(slug);
  if (!page) notFound();
  return { title: page.data.title, alternates: { canonical: `${page.url.replace(/\/$/, '')}/` } };
}
