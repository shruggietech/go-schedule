import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { expect, it, vi } from 'vitest'
import type { TaskBridge, TaskSummary } from './model'
import { TasksPage } from './TasksPage'

it('filters by search and effective state while preserving full group context', async () => {
  const tasks: TaskSummary[] = [{ id: 'one', name: 'Nightly backup', groupId: 'nightly', groupPath: 'Operations / Nightly', commandConfigured: true, declaredEnabled: true, effectiveState: 'runnable', effectiveReason: 'Ready to run.', lifecycle: 'active', timezone: 'UTC', scheduleSummary: 'Every day at 09:00', policySummary: '', nextRuns: [], updatedAt: 'now' }, { id: 'two', name: 'Manual report', groupId: '', groupPath: 'Not assigned', commandConfigured: true, declaredEnabled: false, effectiveState: 'manual_only', effectiveReason: 'No automatic activation source is configured.', lifecycle: 'active', timezone: 'UTC', scheduleSummary: 'Manual only', policySummary: '', nextRuns: [], updatedAt: 'now' }]
  const bridge = { workspace: vi.fn().mockResolvedValue({ action: 'load', outcome: 'accepted', message: 'Loaded.', workspace: { tasks, groups: [], loadedAt: 'now' } }) } as unknown as TaskBridge
  const user = userEvent.setup(); render(<TasksPage bridge={bridge} platform="linux" onActivity={() => undefined} />); await waitFor(() => expect(screen.getByRole('button', { name: 'Nightly backup' })).toBeVisible()); expect(screen.getAllByText('Operations / Nightly')).toHaveLength(2); await user.selectOptions(screen.getByLabelText('State'), 'manual_only'); expect(screen.queryByRole('button', { name: 'Nightly backup' })).not.toBeInTheDocument(); expect(screen.getByRole('button', { name: 'Manual report' })).toBeVisible(); await user.type(screen.getByRole('searchbox', { name: 'Search' }), 'missing'); expect(screen.getByRole('heading', { name: 'No matching tasks' })).toBeVisible()
})

it('forwards remote target and capability boundaries through the populated workspace', async () => {
  const tasks: TaskSummary[] = [{ id: 'one', name: 'Remote backup', groupId: '', groupPath: 'Not assigned', commandConfigured: true, declaredEnabled: true, effectiveState: 'runnable', effectiveReason: 'Ready to run.', lifecycle: 'active', timezone: 'UTC', scheduleSummary: 'Every day', policySummary: '', nextRuns: [], updatedAt: 'now' }]
  const bridge = { workspace: vi.fn().mockResolvedValue({ action: 'load', outcome: 'accepted', message: 'Loaded.', workspace: { tasks, groups: [], loadedAt: 'now' } }) } as unknown as TaskBridge
  const user = userEvent.setup(); render(<TasksPage bridge={bridge} platform="linux" available manageAvailable editAvailable={false} groupAvailable={false} previewAvailable={false} targetName="Production (https://example.test, daemon-i)" onActivity={() => undefined} />)
  await user.click(await screen.findByRole('button', { name: 'Remote backup' }))
  expect(screen.getByRole('button', { name: 'New group' })).toBeDisabled()
  expect(screen.getByRole('button', { name: 'Edit' })).toBeDisabled()
  await user.click(screen.getByRole('button', { name: 'Create task' }))
  expect(screen.getByRole('button', { name: 'Preview' })).toBeDisabled()
  expect(screen.getAllByText('Production (https://example.test, daemon-i)').length).toBeGreaterThan(0)
})

it('opens task creation in a modal and returns focus to its invoker', async () => {
  const bridge = { workspace: vi.fn().mockResolvedValue({ action: 'load', outcome: 'accepted', message: 'Loaded.', workspace: { tasks: [], groups: [], loadedAt: 'now' } }) } as unknown as TaskBridge
  const user = userEvent.setup(); render(<TasksPage bridge={bridge} platform="linux" onActivity={() => undefined} />)
  const create = await screen.findByRole('button', { name: 'Create task' })
  await user.click(create)
  expect(screen.getByRole('dialog', { name: 'Create task' })).toBeVisible()
  expect(screen.getByText('Advanced settings').closest('details')).not.toHaveAttribute('open')
  await user.click(screen.getByRole('button', { name: 'Cancel' }))
  expect(screen.queryByRole('dialog')).not.toBeInTheDocument()
  expect(create).toHaveFocus()
})

it('blocks another mutation until an uncertain result is reconciled by a fresh workspace', async () => {
  const tasks: TaskSummary[] = [{ id: 'one', name: 'Remote backup', groupId: '', groupPath: 'Not assigned', commandConfigured: true, declaredEnabled: true, effectiveState: 'runnable', effectiveReason: 'Ready to run.', lifecycle: 'active', timezone: 'UTC', scheduleSummary: 'Manual only', policySummary: '', nextRuns: [], updatedAt: 'now' }]
  const bridge = { workspace: vi.fn().mockResolvedValue({ action: 'load', outcome: 'accepted', message: 'Loaded.', workspace: { tasks, groups: [], loadedAt: 'now' } }), runTask: vi.fn().mockResolvedValue({ action: 'run_task', outcome: 'uncertain', message: 'The remote request may have completed. Refresh before trying again.' }) } as unknown as TaskBridge
  const user = userEvent.setup()
  const view = render(<TasksPage bridge={bridge} platform="linux" available targetName="Remote" refreshToken={1} onActivity={() => undefined} />)
  await user.click(await screen.findByRole('button', { name: 'Remote backup' }))
  await user.click(screen.getByRole('button', { name: 'Run now' }))
  await user.click(screen.getByRole('button', { name: 'Confirm run' }))
  await waitFor(() => expect(screen.getByText(/may have completed/)).toBeVisible())
  expect(screen.getByRole('button', { name: 'Run now' })).toBeDisabled()
  view.rerender(<TasksPage bridge={bridge} platform="linux" available targetName="Remote" refreshToken={2} onActivity={() => undefined} />)
  await waitFor(() => expect(screen.getByRole('button', { name: 'Run now' })).toBeEnabled())
})
