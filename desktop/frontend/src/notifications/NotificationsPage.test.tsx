import { render, screen, waitFor, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'
import { NotificationsPage } from './NotificationsPage'
import type { NotificationBridge, NotificationResult, NotificationWorkspace, Policy } from './model'

const workspace: NotificationWorkspace = {
  channels: [{ id: 'c1', name: 'Ops hook', kind: 'webhook', endpointSummary: 'https://example.test/...', hasAuthorization: true, enabled: true, updatedAt: '2026-09-07T12:00:00Z' }, { id: 'c2', name: 'Disabled hook', kind: 'webhook', endpointSummary: 'https://disabled.test/...', hasAuthorization: false, enabled: false, updatedAt: '2026-09-07T12:00:00Z' }],
  tasks: [{ type: 'task', id: 't1', name: 'Backup', context: 'Operations / Nightly' }],
  groups: [{ type: 'group', id: 'g1', name: 'Operations', context: 'Operations' }],
  deliveries: [
    { id: 'd1', channelId: 'c1', channelName: 'Ops hook', destinationSummary: 'https://example.test/...', kind: 'test', state: 'queued', attempts: 0, createdAt: '2026-09-07T12:00:00Z' },
    { id: 'd2', channelId: 'c1', channelName: 'Ops hook', destinationSummary: 'https://example.test/...', kind: 'task_outcome', taskId: 't1', taskName: 'Backup', runId: 'r1', state: 'retrying', attempts: 1, nextAttemptAt: '2026-09-07T12:05:00Z', createdAt: '2026-09-07T12:01:00Z' },
    { id: 'd3', channelId: 'c1', channelName: 'Ops hook', destinationSummary: 'https://example.test/...', kind: 'task_outcome', state: 'sending', attempts: 1, createdAt: '2026-09-07T12:02:00Z' },
    { id: 'd4', channelId: 'c1', channelName: 'Ops hook', destinationSummary: 'https://example.test/...', kind: 'task_outcome', state: 'successful', attempts: 1, createdAt: '2026-09-07T12:03:00Z' },
    { id: 'd5', channelId: 'c1', channelName: 'Ops hook', destinationSummary: 'https://example.test/...', kind: 'task_outcome', state: 'failed', attempts: 3, lastStatus: 503, lastError: 'receiver unavailable', createdAt: '2026-09-07T12:04:00Z' },
  ], loadedAt: '2026-09-07T12:05:00Z',
}
const policy: Policy = { scope: workspace.tasks[0], directAssignments: [], effectiveSourceType: 'group', effectiveSourceId: 'g1', effectiveSourceName: 'Operations', effectiveAssignments: [{ channelId: 'c1', onFailure: true, onSuccess: false }] }

function bridge(overrides: Partial<NotificationBridge> = {}): NotificationBridge {
  const accepted = (action: string, extra: Partial<NotificationResult> = {}): Promise<NotificationResult> => Promise.resolve({ action, outcome: 'accepted', message: 'Done.', ...extra })
  return { workspace: () => accepted('load_notifications', { workspace }), saveChannel: () => accepted('save_notification_channel', { workspace }), setChannelEnabled: () => accepted('toggle_notification_channel', { workspace }), testChannel: () => accepted('test_notification_channel', { workspace }), deleteChannel: () => accepted('delete_notification_channel', { workspace }), policy: () => accepted('load_notification_policy', { policy }), savePolicy: () => accepted('save_notification_policy', { policy }), ...overrides }
}

describe('NotificationsPage', () => {
  it('never redisplays secrets and requires explicit replacement intent', async () => {
    const user = userEvent.setup(); const saveChannel = vi.fn(bridge().saveChannel)
    render(<NotificationsPage bridge={bridge({ saveChannel })} available refreshToken={1} />)
    await screen.findAllByText('Ops hook'); await user.click(screen.getAllByRole('button', { name: 'Edit' })[0])
    expect(screen.queryByDisplayValue('Bearer secret')).not.toBeInTheDocument(); expect(screen.queryByDisplayValue('https://example.test/hook')).not.toBeInTheDocument()
    await user.click(screen.getByLabelText('Replace stored endpoint')); expect(screen.getByLabelText('New HTTPS endpoint')).toHaveValue('')
    await user.type(screen.getByLabelText('New HTTPS endpoint'), 'https://new.test/hook')
    await user.click(screen.getByLabelText('Replace or remove stored authorization')); await user.click(screen.getByRole('button', { name: 'Save channel' }))
    await waitFor(() => expect(saveChannel).toHaveBeenCalledWith(expect.objectContaining({ replaceEndpoint: true, endpoint: 'https://new.test/hook', replaceAuthorization: true, authorization: '' })))
    expect(screen.getByLabelText('HTTPS endpoint')).toHaveValue(''); expect(screen.getByLabelText('Authorization (optional)')).toHaveValue('')
  })

  it('supports tests, disabled explanation, and confirmed removal', async () => {
    const user = userEvent.setup(); const testChannel = vi.fn(bridge().testChannel); const deleteChannel = vi.fn(bridge().deleteChannel)
    render(<NotificationsPage bridge={bridge({ testChannel, deleteChannel })} available refreshToken={1} />); await screen.findAllByText('Ops hook')
    const cards = screen.getAllByRole('article'); expect(within(cards[1]).getByRole('button', { name: 'Send test' })).toBeDisabled()
    await user.click(within(cards[0]).getByRole('button', { name: 'Send test' })); await waitFor(() => expect(testChannel).toHaveBeenCalledWith('c1'))
    expect(await screen.findByRole('status', { name: 'Notification action complete' })).toHaveTextContent('Done.')
    await user.click(within(cards[0]).getByRole('button', { name: 'Remove' })); expect(screen.getByRole('dialog', { name: 'Remove webhook channel?' })).toBeInTheDocument()
    await user.click(screen.getByRole('button', { name: 'Remove channel' })); await waitFor(() => expect(deleteChannel).toHaveBeenCalledWith('c1'))
  })

  it('explains inheritance and warns without blocking success assignments', async () => {
    const user = userEvent.setup(); const savePolicy = vi.fn(bridge().savePolicy)
    render(<NotificationsPage bridge={bridge({ savePolicy })} available refreshToken={1} />); await screen.findAllByText('Ops hook')
    await user.selectOptions(screen.getByLabelText('Task or group'), 'task:t1'); expect(await screen.findByText(/inherits from Operations/)).toBeInTheDocument()
    const policyEditor = screen.getByText('Direct assignments for Backup').closest('fieldset')!; await user.click(within(policyEditor).getAllByLabelText('Success')[0])
    expect(screen.getByText('Success notifications can be noisy')).toBeInTheDocument(); await user.click(screen.getByRole('button', { name: 'Save direct policy' }))
    await waitFor(() => expect(savePolicy).toHaveBeenCalledWith(expect.objectContaining({ scopeType: 'task', scopeId: 't1', assignments: [expect.objectContaining({ channelId: 'c1', onSuccess: true })] })))
  })

  it('hides a prior policy immediately when the selected scope changes or clears', async () => {
    const user = userEvent.setup(); const next = new Promise<NotificationResult>(() => undefined); const loadPolicy = vi.fn().mockResolvedValueOnce({ action: 'load_notification_policy', outcome: 'accepted', message: 'Done.', policy }).mockReturnValueOnce(next)
    render(<NotificationsPage bridge={bridge({ policy: loadPolicy })} available refreshToken={1} />); await screen.findAllByText('Ops hook')
    await user.selectOptions(screen.getByLabelText('Task or group'), 'task:t1'); expect(await screen.findByText('Direct assignments for Backup')).toBeInTheDocument()
    await user.selectOptions(screen.getByLabelText('Task or group'), 'group:g1'); expect(screen.queryByText('Direct assignments for Backup')).not.toBeInTheDocument(); expect(screen.getByText('Loading policy')).toBeInTheDocument()
    await user.selectOptions(screen.getByLabelText('Task or group'), ''); expect(screen.queryByText('Direct assignments for Backup')).not.toBeInTheDocument()
  })

  it('distinguishes every delivery state and test detail from task outcomes', async () => {
    const user = userEvent.setup(); render(<NotificationsPage bridge={bridge()} available refreshToken={1} />)
    expect(await screen.findByText('5 matching deliveries')).toBeInTheDocument()
    for (const state of ['Queued', 'Retrying', 'Sending', 'Successful', 'Failed']) expect(screen.getAllByText(state).length).toBeGreaterThan(0)
    expect(screen.getByText('Test')).toBeInTheDocument(); expect(screen.getAllByText('Task outcome')).toHaveLength(4)
    await user.click(screen.getAllByRole('button', { name: /9\/7\/2026/ })[0]); expect(screen.getByText('Test delivery')).toBeInTheDocument(); expect(screen.getByText(/not a task completion/)).toBeInTheDocument()
    await user.selectOptions(screen.getByLabelText('State'), 'failed'); expect(screen.getByText('1 matching deliveries')).toBeInTheDocument(); await user.click(screen.getByRole('button', { name: /9\/7\/2026/ })); expect(screen.getByText('receiver unavailable')).toBeInTheDocument()
  })
})
