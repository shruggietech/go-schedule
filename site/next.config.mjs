import { createMDX } from 'fumadocs-mdx/next';

export default createMDX()({
  output: 'export',
  staticPageGenerationTimeout: 300,
  trailingSlash: true,
  images: { unoptimized: true },
  experimental: { cpus: 2 },
});
