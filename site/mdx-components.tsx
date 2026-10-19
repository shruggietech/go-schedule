import defaults from 'fumadocs-ui/mdx';
import type { MDXComponents } from 'mdx/types';
import { Mermaid } from '@/components/mermaid';

export function getMDXComponents(): MDXComponents {
  return { ...defaults, Mermaid };
}
