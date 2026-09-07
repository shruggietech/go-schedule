import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import axe from 'axe-core'
import { describe, expect, it, vi } from 'vitest'
import { App } from './App'
import type { DesktopBridge } from './connection/model'

const bridge: DesktopBridge = { snapshot: vi.fn().mockResolvedValue({ generation: 1, state: 'connected', target: { id: 'local', displayName: 'This computer', platform: 'linux', capabilities: ['tasks'], permissions: ['read'] }, message: 'Available.' }), retry: vi.fn(), quit: vi.fn(), subscribe: () => () => undefined }

describe('accessibility contract', () => {
  it('has no serious or critical automated violations', async () => {
    render(<App bridge={bridge} />)
    const results = await axe.run(document, { runOnly: { type: 'tag', values: ['wcag2a', 'wcag2aa', 'wcag22aa'] }, rules: { 'color-contrast': { enabled: false } } })
    expect(results.violations.filter((item) => item.impact === 'serious' || item.impact === 'critical')).toEqual([])
  })

  it('returns dialog focus and provides landmarks and a live region', async () => {
    const user = userEvent.setup(); render(<App bridge={bridge} />)
    const trigger = screen.getByRole('button', { name: 'Connection details' })
    await user.click(trigger); await user.keyboard('{Escape}')
    expect(trigger).toHaveFocus()
    expect(screen.getByRole('navigation', { name: 'Application' })).toBeVisible()
    expect(screen.getByRole('main')).toBeVisible()
    expect(document.querySelector('[aria-live="polite"]')).not.toBeNull()
  })
})
