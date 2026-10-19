import { defineDocs, defineConfig } from 'fumadocs-mdx/config';

export const docs = defineDocs({ dir: 'content/docs' });

function remarkMermaid() {
  // The compiler's generic Markdown tree exposes different node fields by kind.
  return (tree: any) => {
    const walk = (node: any) => {
      if (!Array.isArray(node?.children)) return;
      node.children = node.children.map((child: any) => {
        if (child.type === 'code' && child.lang === 'mermaid') return {
          type: 'mdxJsxFlowElement', name: 'Mermaid',
          attributes: [{ type: 'mdxJsxAttribute', name: 'chart', value: child.value }], children: [],
        };
        walk(child);
        return child;
      });
    };
    walk(tree);
  };
}

export default defineConfig({ mdxOptions: { remarkPlugins: [remarkMermaid] } });
