'use client';

import Link from 'next/link';
import { useTheme } from 'next-themes';
import { useEffect, useState } from 'react';

export function ThemeToggle() {
  const { resolvedTheme, setTheme } = useTheme();
  const [ready, setReady] = useState(false);
  useEffect(() => setReady(true), []);
  const dark = ready && resolvedTheme === 'dark';
  return <button type="button" className="theme-toggle" aria-label={`Switch to ${dark ? 'light' : 'dark'} mode`} onClick={() => setTheme(dark ? 'light' : 'dark')}><span aria-hidden="true">{dark ? '☀' : '◐'}</span><span>{dark ? 'Light' : 'Dark'}</span></button>;
}

export function Header() {
  return <header className="site-header"><Link href="/" className="brand" aria-label="go-schedule home"><img src="/assets/brand/go-schedule-mark-reduced-color.svg" width="34" height="34" alt="" /><span>go-schedule</span></Link><nav aria-label="Main navigation"><Link href="/docs/">Documentation</Link><Link href="/compare/">Compare</Link><a href="https://github.com/shruggietech/go-schedule">GitHub</a><ThemeToggle /></nav></header>;
}

export function Footer() {
  return <footer className="site-footer"><p>A <a href="https://shruggie.tech">ShruggieTech</a> project.</p><nav aria-label="Project resources"><a href="https://github.com/shruggietech/go-schedule">GitHub</a><Link href="/docs/changelog/">Changelog</Link><a href="https://github.com/shruggietech/go-schedule/releases/latest">Latest release</a><Link href="/docs/brand/">Brand</Link></nav><p className="fine-print">Open source. Apache License 2.0.</p></footer>;
}
