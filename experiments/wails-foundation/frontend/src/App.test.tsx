import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it } from 'vitest'
import { App } from './App'

describe('control center', () => {
  it('keeps target identity and one page heading visible while navigating', async () => {
    const user = userEvent.setup()
    render(<App />)
    expect(screen.getAllByText('This computer').length).toBeGreaterThan(0)
    expect(screen.getByRole('heading', { level: 1, name: 'Tasks' })).toBeVisible()
    await user.click(screen.getByRole('button', { name: 'Activity' }))
    expect(screen.getByRole('heading', { level: 1, name: 'Activity' })).toBeVisible()
    expect(screen.getAllByText('This computer').length).toBeGreaterThan(0)
  })

  it('exposes all required states through the prototype controls', async () => {
    const user = userEvent.setup()
    render(<App />)
    await user.selectOptions(screen.getByLabelText('State'), 'disconnected')
    expect(screen.getByRole('heading', { name: 'This computer is offline' })).toBeVisible()
    expect(screen.getByRole('button', { name: 'Try again' })).toBeVisible()
  })

  it('opens the task workflow from its primary action', async () => {
    const user = userEvent.setup()
    render(<App />)
    await user.click(screen.getByRole('button', { name: 'Create task' }))
    expect(screen.getByRole('heading', { level: 1, name: 'Task editor' })).toBeVisible()
    expect(screen.getByRole('button', { name: 'Save task' })).toBeVisible()
  })
})
