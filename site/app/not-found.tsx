import Link from 'next/link';
import { Header, Footer } from '@/components/chrome';

export default function NotFound() {
  return <><Header /><main id="main-content" className="not-found"><h1>This page is no longer here.</h1><p>Find the current guide in the documentation, or return to the project home.</p><Link className="button primary" href="/docs/">Open documentation</Link><Link className="button secondary" href="/">Project home</Link></main><Footer /></>;
}
