import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'
import { SearchPage } from './SearchPage'
import type { SearchBridge, SearchSnapshot } from './model'

const snapshot: SearchSnapshot = {
  generation: 2,
  query: 'archive',
  startedAt: '2026-09-21T12:00:00Z',
  completedAt: '2026-09-21T12:00:01Z',
  complete: true,
  observations: [
    { registration: { key: 'local', kind: 'local', label: 'Workshop', daemonId: 'daemon-a', shortDaemonId: 'daemon-a' }, state: 'connected', observedAt: '2026-09-21T12:00:00Z', truncated: false, matches: [{ registrationKey: 'local', expectedDaemonId: 'daemon-a', sourceLabel: 'Workshop', sourceShortId: 'daemon-a', result: { kind: 'task', object_id: 'task-a', task_id: 'task-a', name: 'Archive', context: 'Active', action_hints: ['run_now'] }, availableActions: ['open', 'run_now'] }] },
    { registration: { key: 'remote', kind: 'remote', label: 'Workshop', daemonId: 'daemon-b', shortDaemonId: 'daemon-b' }, state: 'connected', observedAt: '2026-09-21T12:00:00Z', truncated: false, matches: [{ registrationKey: 'remote', expectedDaemonId: 'daemon-b', sourceLabel: 'Workshop', sourceShortId: 'daemon-b', result: { kind: 'task', object_id: 'task-b', task_id: 'task-b', name: 'Archive', context: 'Disabled', action_hints: ['enable'] }, availableActions: ['open', 'enable'] }] },
  ],
}

function bridge(): SearchBridge {
  return { search: vi.fn().mockResolvedValue(snapshot), execute: vi.fn().mockResolvedValue({ action: 'run_now', outcome: 'accepted', message: '1 of 1 requested objects accepted.', outcomes: [] }), subscribe: () => () => undefined }
}

describe('SearchPage', () => {
  it('keeps duplicate names source-qualified and opens the exact daemon identity', async () => {
    const user = userEvent.setup()
    const onOpen = vi.fn().mockResolvedValue(undefined)
    render(<SearchPage bridge={bridge()} onOpen={onOpen} />)
    await user.type(screen.getByLabelText('Search all systems'), 'archive')
    await user.click(screen.getByRole('button', { name: 'Search' }))
    expect(await screen.findAllByText('Archive')).toHaveLength(2)
    expect(screen.getByText('daemon-a')).toBeVisible()
    expect(screen.getByText('daemon-b')).toBeVisible()
    await user.click(screen.getAllByRole('button', { name: 'Open' })[1])
    expect(onOpen).toHaveBeenCalledWith(expect.objectContaining({ registrationKey: 'remote', expectedDaemonId: 'daemon-b', taskId: 'task-b' }))
  })

  it('requires compatible selections and groups confirmation by source', async () => {
    const user = userEvent.setup()
    render(<SearchPage bridge={bridge()} onOpen={vi.fn()} />)
    await user.type(screen.getByLabelText('Search all systems'), 'archive')
    await user.click(screen.getByRole('button', { name: 'Search' }))
    await screen.findByText('2 systems searched; 2 results.')
    await user.click(screen.getAllByLabelText('Select Archive from Workshop')[0])
    const review = screen.getByRole('button', { name: 'Review 1 selected' })
    expect(review).toBeEnabled()
    await user.click(review)
    expect(screen.getByRole('dialog', { name: 'Confirm run now' })).toHaveTextContent('Workshop (daemon-a)')
    expect(screen.getByRole('button', { name: 'Cancel' })).toHaveFocus()
  })

  it('announces partial target failures without hiding successful results', async () => {
    const partial = { ...snapshot, observations: [...snapshot.observations, { registration: { key: 'lost', kind: 'remote' as const, label: 'Offline', shortDaemonId: 'daemon-c' }, state: 'timed_out' as const, truncated: false, matches: [], failure: { state: 'timed_out' as const, message: 'Timed out.', action: 'Try again.' } }] }
    const value = bridge(); value.search = vi.fn().mockResolvedValue(partial)
    const user = userEvent.setup(); render(<SearchPage bridge={value} onOpen={vi.fn()} />)
    await user.type(screen.getByLabelText('Search all systems'), 'archive'); await user.click(screen.getByRole('button', { name: 'Search' }))
    expect(await screen.findByText('Timed out. Try again.')).toBeVisible()
    expect(screen.getAllByText('Archive')).toHaveLength(2)
  })
})
