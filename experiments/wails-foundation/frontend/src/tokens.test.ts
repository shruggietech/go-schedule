import { describe, expect, it } from 'vitest'
import css from './styles.css?raw'

describe('local design tokens', () => {
  it('defines every approved semantic color role', () => {
    for (const token of ['night', 'panel', 'raised', 'line', 'text', 'muted', 'paper', 'ink', 'interval', 'anchor', 'hold', 'stop']) {
      expect(css).toContain(`--${token}:`)
    }
  })

  it('uses only local font sources and honors reduced motion', () => {
    expect(css).not.toMatch(/url\(['"]?https?:/)
    expect(css).toContain('@media (prefers-reduced-motion: reduce)')
  })
})
