import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'
import { App } from './App'
import { fixtureFor } from './proof'

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

  it('derives disconnected and empty conditions from real snapshots', async () => {
    const disconnected = fixtureFor('disconnected')
    const { rerender } = render(<App snapshotLoader={() => Promise.resolve(disconnected)} />)
    expect(await screen.findByRole('heading', { name: 'This computer is offline' })).toBeVisible()
    rerender(<App snapshotLoader={() => Promise.resolve(fixtureFor('empty'))} />)
    expect(await screen.findByRole('heading', { name: 'No tasks yet' })).toBeVisible()
  })

  it('makes every page-level primary action operational', async () => {
    const user = userEvent.setup()
    render(<App />)
    await user.click(screen.getByRole('button', { name: 'Schedule' }))
    await user.click(screen.getByRole('button', { name: 'View tasks' }))
    expect(screen.getByRole('heading', { level: 1, name: 'Tasks' })).toBeVisible()
    await user.click(screen.getByRole('button', { name: 'Targets' }))
    await user.click(screen.getAllByRole('button', { name: 'Switch target' }).at(-1)!)
    expect(screen.getByRole('dialog', { name: 'Switch target' })).toBeVisible()
  })

  it('cancels loading and retries a disconnected snapshot', async () => {
    const user = userEvent.setup()
    const loader = vi.fn().mockRejectedValueOnce(new Error('offline')).mockResolvedValue(fixtureFor('connected'))
    render(<App snapshotLoader={loader} />)
    expect(await screen.findByRole('heading', { name: 'This computer is offline' })).toBeVisible()
    await user.click(screen.getByRole('button', { name: 'Try again' }))
    await waitFor(() => expect(screen.getByRole('heading', { name: 'Scheduled work' })).toBeVisible())
    await user.selectOptions(screen.getByLabelText('State'), 'loading')
    await user.click(screen.getByRole('button', { name: 'Cancel' }))
    expect(screen.getByRole('heading', { name: 'Scheduled work' })).toBeVisible()
  })
})
