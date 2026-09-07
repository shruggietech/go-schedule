import { render, screen } from '@testing-library/react'
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
})
