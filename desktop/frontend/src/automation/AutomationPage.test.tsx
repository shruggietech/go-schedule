import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import axe from 'axe-core'
import { describe, expect, it, vi } from 'vitest'
import { AutomationPage } from './AutomationPage'
import type { AutomationBridge, AutomationWorkspace, OperationResult, SecretResult } from './model'

const workspace: AutomationWorkspace = {
  tasks: [{ id: 'task-1', name: 'Build', readiness: 'ready', reason: 'Ready.' }, { id: 'task-2', name: 'Deploy', readiness: 'ready', reason: 'Ready.' }],
  chains: [{ id: 'chain-1', sourceTaskId: 'task-1', sourceTaskName: 'Build', targetTaskId: 'task-2', targetTaskName: 'Deploy', onOutcome: 'success', readiness: 'ready', reason: 'Runs the target.', updatedAt: '2026-09-07T00:00:00Z' }],
  triggers: [{ id: 'trigger-1', name: 'Webhook', targetTaskId: 'task-2', targetTaskName: 'Deploy', enabled: true, readiness: 'target_missing', reason: 'Target task is missing.', updatedAt: '2026-09-07T00:00:00Z' }],
  triggerSets: [{ id: 'set-1', name: 'Fleet', targetTaskId: 'task-2', targetTaskName: 'Deploy', memberCount: 2, enabledCount: 2, readiness: 'ready', reason: 'All enabled members are ready.', updatedAt: '2026-09-07T00:00:00Z', members: [{ id: 'm1', name: 'Fleet 1', position: 1, enabled: true, readiness: 'ready', reason: '' }, { id: 'm2', name: 'Fleet 2', position: 2, enabled: true, readiness: 'ready', reason: '' }] }],
  watchers: [{ id: 'watcher-1', name: 'Inbox', kind: 'directory', path: 'C:\\a\\very\\long\\inbox', pattern: '*.csv', recursive: true, debounce: '500ms', stability: '1s', targetTaskId: 'task-1', targetTaskName: 'Build', enabled: true, health: 'degraded', healthReason: 'Permission denied.', readiness: 'ready', reason: '', updatedAt: '2026-09-07T00:00:00Z' }],
  loadedAt: '2026-09-07T00:00:00Z',
}
const accepted: OperationResult = { action: 'load', outcome: 'accepted', message: 'Loaded.', workspace }
function mockBridge(overrides: Partial<AutomationBridge> = {}): AutomationBridge {
  const operation = vi.fn().mockResolvedValue(accepted); const secret = vi.fn().mockResolvedValue({ action: 'reveal', outcome: 'accepted', message: 'Revealed.', title: 'Webhook', secrets: [{ label: 'Webhook', key: 'secret-value', command: 'gosched trigger fire secret-value' }] } satisfies SecretResult)
  return { workspace: vi.fn().mockResolvedValue(accepted), saveChain: operation, deleteChain: operation, saveTrigger: secret, setTriggerEnabled: operation, revealTrigger: secret, rotateTrigger: secret, fireTrigger: operation, deleteTrigger: operation, createTriggerSet: secret, retargetTriggerSet: operation, setTriggerSetEnabled: operation, revealTriggerSet: secret, rotateTriggerSet: secret, deleteTriggerSet: operation, saveWatcher: operation, setWatcherEnabled: operation, deleteWatcher: operation, ...overrides }
}

describe('AutomationPage', () => {
  it('presents all sources with textual state and passes accessibility checks', async () => {
    const { container } = render(<AutomationPage bridge={mockBridge()} />)
    expect(await screen.findByRole('heading', { name: 'Completion chains' })).toBeVisible()
    expect(screen.getByRole('heading', { name: 'External triggers' })).toBeVisible()
    expect(screen.getByRole('heading', { name: 'Trigger Sets' })).toBeVisible()
    expect(screen.getByRole('heading', { name: 'Filesystem watchers' })).toBeVisible()
    expect(screen.getByText('Permission denied.')).toBeVisible()
    expect(screen.getByText('Target task is missing.')).toBeVisible()
    const results = await axe.run(container, { rules: { 'color-contrast': { enabled: false } } })
    expect(results.violations.filter((item) => item.impact === 'serious' || item.impact === 'critical')).toEqual([])
  })

  it('filters a large collection without losing search focus', async () => {
    const user = userEvent.setup(); const large = { ...workspace, triggers: Array.from({ length: 100 }, (_, i) => ({ ...workspace.triggers[0], id: `trigger-${i}`, name: `Hook ${i}`, readiness: i === 99 ? 'target_missing' : 'ready' })) }
    render(<AutomationPage bridge={mockBridge({ workspace: vi.fn().mockResolvedValue({ ...accepted, workspace: large }) })} />)
    const search = await screen.findByRole('searchbox', { name: 'Search' }); await user.type(search, 'Hook 99')
    expect(screen.getByRole('heading', { name: 'Hook 99' })).toBeVisible(); expect(screen.queryByRole('heading', { name: 'Hook 1' })).not.toBeInTheDocument(); expect(search).toHaveFocus()
    await user.click(screen.getByRole('checkbox', { name: 'Needs attention only' })); expect(screen.getByRole('heading', { name: 'Hook 99' })).toBeVisible()
  })

  it('retains the last complete workspace as read-only context', async () => {
    render(<AutomationPage bridge={mockBridge()} available={false} />)
    expect(await screen.findByText('Read-only context')).toBeVisible()
    expect(screen.getByRole('heading', { name: 'Webhook' })).toBeVisible()
    expect(screen.getAllByRole('button', { name: 'Create' })[0]).toBeDisabled()
  })

  it('shows secrets only after an explicit reveal and clears the dialog', async () => {
    const user = userEvent.setup(); render(<AutomationPage bridge={mockBridge()} />)
    await screen.findByRole('heading', { name: 'Webhook' })
    expect(screen.queryByText('secret-value')).not.toBeInTheDocument()
    await user.click(screen.getByRole('button', { name: 'Reveal' })); expect(await screen.findByText('secret-value')).toBeVisible()
    await user.click(screen.getByRole('button', { name: 'Close' })); expect(screen.queryByText('secret-value')).not.toBeInTheDocument()
  })

  it('suppresses duplicate fire submissions while pending', async () => {
    const user = userEvent.setup(); let finish!: (value: OperationResult) => void; const fire = vi.fn().mockReturnValue(new Promise<OperationResult>((resolve) => { finish = resolve }))
    const readyWorkspace = { ...workspace, triggers: [{ ...workspace.triggers[0], readiness: 'ready', reason: 'Ready.' }] }
    render(<AutomationPage bridge={mockBridge({ workspace: vi.fn().mockResolvedValue({ ...accepted, workspace: readyWorkspace }), fireTrigger: fire })} />)
    const button = await screen.findByRole('button', { name: 'Fire now' }); await user.click(button); await user.click(button)
    expect(fire).toHaveBeenCalledTimes(1); finish({ action: 'fire_trigger', outcome: 'accepted', message: 'Fired.', workspace: readyWorkspace }); await waitFor(() => expect(button).toBeEnabled())
  })

  it('confirms key rotations with the affected source identity', async () => {
    const user = userEvent.setup(); const rotate = vi.fn().mockResolvedValue({ action: 'rotate_trigger', outcome: 'accepted', message: 'Rotated.', secrets: [{ label: 'Webhook', key: 'replacement', command: 'fire replacement' }] }); const confirm = vi.spyOn(window, 'confirm').mockReturnValue(false)
    render(<AutomationPage bridge={mockBridge({ rotateTrigger: rotate })} />); const button = await screen.findByRole('button', { name: 'Rotate' }); await user.click(button); expect(rotate).not.toHaveBeenCalled(); expect(confirm).toHaveBeenCalledWith(expect.stringContaining('Webhook'))
    confirm.mockReturnValue(true); await user.click(button); expect(rotate).toHaveBeenCalledTimes(1); confirm.mockRestore()
  })

  it('preserves trigger-save success when its refresh fails', async () => {
    const user = userEvent.setup(); const load = vi.fn().mockResolvedValueOnce(accepted).mockResolvedValue({ action: 'load', outcome: 'unavailable', message: 'Offline.' }); const save = vi.fn().mockResolvedValue({ action: 'save_trigger', outcome: 'accepted', message: 'Trigger saved.', entityId: 'trigger-1' })
    render(<AutomationPage bridge={mockBridge({ workspace: load, saveTrigger: save })} />); const trigger = await screen.findByRole('heading', { name: 'Webhook' }); await user.click(trigger.closest('article')!.querySelector<HTMLButtonElement>('button')!); await user.click(screen.getByRole('button', { name: 'Save' })); expect(await screen.findByText(/Trigger saved\. Refresh the workspace/)).toBeVisible(); expect(screen.getByRole('heading', { name: 'Webhook' })).toBeVisible()
  })

  it('traps focus inside the source editor and restores its invoker', async () => {
    const user = userEvent.setup(); render(<AutomationPage bridge={mockBridge()} />); const heading = await screen.findByRole('heading', { name: 'External triggers' }); const create = heading.closest('section')!.querySelector<HTMLButtonElement>('button')!; await user.click(create)
    expect(screen.getByRole('dialog', { name: 'Create external trigger' })).toBeVisible(); expect(screen.getByLabelText('Name')).toHaveFocus(); await user.keyboard('{Shift>}{Tab}{/Shift}'); expect(screen.getByRole('button', { name: 'Cancel' })).toHaveFocus(); await user.keyboard('{Escape}'); expect(create).toHaveFocus()
  })
})
