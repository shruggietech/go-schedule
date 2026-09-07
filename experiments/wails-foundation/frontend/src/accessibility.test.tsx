import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import axe from 'axe-core'
import { describe, expect, it } from 'vitest'
import { App } from './App'

describe('accessibility contract', () => {
  it('has no serious or critical automated violations in the connected view', async () => {
    render(<App />)
    const results = await axe.run(document, { runOnly: { type: 'tag', values: ['wcag2a', 'wcag2aa', 'wcag22aa'] }, rules: { 'color-contrast': { enabled: false } } })
    expect(results.violations.filter((violation) => violation.impact === 'serious' || violation.impact === 'critical')).toEqual([])
  })

  it('closes the target dialog with Escape and returns focus', async () => {
    const user = userEvent.setup()
    render(<App />)
    const trigger = screen.getByRole('button', { name: 'Switch target' })
    await user.click(trigger)
    expect(screen.getByRole('dialog', { name: 'Switch target' })).toBeVisible()
    expect(screen.getByRole('button', { name: 'Close' })).toHaveFocus()
    await user.keyboard('{Escape}')
    expect(screen.queryByRole('dialog')).not.toBeInTheDocument()
    expect(trigger).toHaveFocus()
  })

  it('provides named navigation, main content, and a polite live region', () => {
    render(<App />)
    expect(screen.getByRole('navigation', { name: 'Application' })).toBeVisible()
    expect(screen.getByRole('main')).toBeVisible()
    expect(document.querySelector('[aria-live="polite"]')).not.toBeNull()
  })
})
