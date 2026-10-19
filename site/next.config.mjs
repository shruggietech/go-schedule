import { createMDX } from 'fumadocs-mdx/next';

export default createMDX()({
  output: 'export',
  trailingSlash: true,
  images: { unoptimized: true },
  experimental: { cpus: 2 },
});
