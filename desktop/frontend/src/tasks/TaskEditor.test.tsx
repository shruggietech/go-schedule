import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'
import { blankTask, type TaskBridge } from './model'
import { TaskEditor } from './TaskEditor'

describe('task editor', () => {
  it('exposes complete basic and advanced intent and previews the exact command boundary', async () => {
    const user = userEvent.setup(); const previewTask = vi.fn().mockResolvedValue({ action: 'preview', outcome: 'accepted', message: 'Every day', command: { program: 'echo', args: ['two words'] } }); const bridge = { previewTask } as unknown as TaskBridge
    render(<TaskEditor initial={{ ...blankTask(), commandLine: `echo "two words"` }} groups={[]} platform="linux" bridge={bridge} onSaved={() => undefined} onCancel={() => undefined} />)
    await user.click(screen.getByText('Advanced settings'))
    for (const label of ['Working directory', 'Run identity', 'Standard input', 'Environment', 'Overlap policy', 'Catch-up policy', 'Missing-date policy', 'Time basis', 'DST gap', 'DST overlap']) expect(screen.getByLabelText(label)).toBeVisible()
    await user.click(screen.getByRole('button', { name: 'Preview' }))
    expect(previewTask).toHaveBeenCalledWith(expect.objectContaining({ commandLine: `echo "two words"` }))
    expect(screen.getByRole('heading', { name: 'Exact launch preview' })).toBeVisible()
    expect(screen.getByText('two words', { exact: false })).toBeVisible()
  })

  it('makes stale overwrite an explicit choice', async () => {
    const user = userEvent.setup(); const saveTask = vi.fn().mockResolvedValue({ action: 'save_task', outcome: 'stale', message: 'Newer version.' }); const bridge = { saveTask } as unknown as TaskBridge
    render(<TaskEditor initial={{ ...blankTask(), commandLine: 'echo' }} groups={[]} platform="linux" bridge={bridge} onSaved={() => undefined} onCancel={() => undefined} />)
    await user.click(screen.getByRole('button', { name: 'Save inactive task' })); expect(screen.getByRole('button', { name: 'Overwrite newer version' })).toBeVisible()
  })

  it('keeps one-off wall time in the selected task timezone', async () => {
    const user = userEvent.setup(); const previewTask = vi.fn().mockResolvedValue({ action: 'preview', outcome: 'accepted', message: 'Valid.' }); const bridge = { previewTask } as unknown as TaskBridge
    render(<TaskEditor initial={{ ...blankTask(), commandLine: 'echo', mode: 'one_off', timezone: 'America/New_York', at: '2030-01-15T14:00:00Z' }} groups={[]} platform="linux" bridge={bridge} onSaved={() => undefined} onCancel={() => undefined} />)
    const runAt = screen.getByLabelText('Run at')
    expect(runAt).toHaveValue('2030-01-15T09:00')
    await user.clear(runAt); await user.type(runAt, '2030-01-16T10:30'); await user.click(screen.getByRole('button', { name: 'Preview' }))
    expect(previewTask).toHaveBeenCalledWith(expect.objectContaining({ at: '2030-01-16T10:30', timezone: 'America/New_York' }))
    expect(screen.getByText('Interpreted in America/New_York.')).toBeVisible()
  })

  it('retains the draft after an uncertain remote save', async () => {
    const user = userEvent.setup(); const saveTask = vi.fn().mockResolvedValue({ action: 'save_task', outcome: 'uncertain', message: 'The remote request may have completed. Refresh before retrying.' }); const onSaved = vi.fn(); const bridge = { saveTask } as unknown as TaskBridge
    render(<TaskEditor initial={{ ...blankTask(), commandLine: 'echo' }} groups={[]} platform="linux" bridge={bridge} onSaved={onSaved} onCancel={() => undefined} />)
    await user.type(screen.getByLabelText('Name'), 'Remote draft')
    await user.click(screen.getByRole('button', { name: 'Save inactive task' }))
    expect(await screen.findByText('The remote request may have completed. Refresh before retrying.')).toBeVisible()
    expect(screen.getByLabelText('Name')).toHaveValue('Remote draft')
    expect(onSaved).not.toHaveBeenCalled()
  })

  it('opens advanced settings before focusing an invalid advanced field', async () => {
    const user = userEvent.setup(); const previewTask = vi.fn().mockResolvedValue({ action: 'preview', outcome: 'rejected', message: 'Working directory is invalid.', field: 'working_dir' }); const bridge = { previewTask } as unknown as TaskBridge
    render(<TaskEditor initial={{ ...blankTask(), commandLine: 'echo' }} groups={[]} platform="linux" bridge={bridge} onSaved={() => undefined} onCancel={() => undefined} />)
    await user.click(screen.getByRole('button', { name: 'Preview' }))
    const field = await screen.findByLabelText('Working directory')
    await waitFor(() => expect(field.closest('details')).toHaveAttribute('open'))
    await waitFor(() => expect(field).toHaveFocus())
  })

  it('cannot be dismissed while a save is pending', async () => {
    const user = userEvent.setup(); let finishSave: ((value: { action: string; outcome: 'accepted'; message: string }) => void) | undefined; const saveTask = vi.fn().mockReturnValue(new Promise((resolve) => { finishSave = resolve })); const onSaved = vi.fn(); const onCancel = vi.fn(); const bridge = { saveTask } as unknown as TaskBridge
    render(<TaskEditor initial={{ ...blankTask(), commandLine: 'echo' }} groups={[]} platform="linux" bridge={bridge} onSaved={onSaved} onCancel={onCancel} />)
    await user.click(screen.getByRole('button', { name: 'Save inactive task' }))
    expect(screen.getByRole('button', { name: 'Cancel' })).toBeDisabled()
    await user.keyboard('{Escape}')
    expect(onCancel).not.toHaveBeenCalled()
    finishSave?.({ action: 'save_task', outcome: 'accepted', message: 'Saved.' })
    await waitFor(() => expect(onSaved).toHaveBeenCalledOnce())
  })
})
