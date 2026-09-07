import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { expect, it, vi } from 'vitest'
import type { GroupSummary, OperationResult, TaskBridge } from './model'
import { GroupsPanel } from './GroupsPanel'

it('uses full paths and excludes the edited group and descendants from parent choices', async () => {
  const user = userEvent.setup(); const groups: GroupSummary[] = [{ id: 'root', name: 'Duplicate', parentId: '', path: 'Duplicate', depth: 0, declaredEnabled: true, effectiveEnabled: true, effectiveReason: '', childCount: 1, taskCount: 0, descendantCount: 1, updatedAt: 'now' }, { id: 'child', name: 'Duplicate', parentId: 'root', path: 'Duplicate / Duplicate', depth: 1, declaredEnabled: true, effectiveEnabled: true, effectiveReason: '', childCount: 0, taskCount: 0, descendantCount: 0, updatedAt: 'now' }]
  render(<GroupsPanel groups={groups} bridge={{} as TaskBridge} onResult={() => undefined} />); await user.click(screen.getByRole('button', { name: 'Duplicate' })); const parent = screen.getByLabelText('Parent group'); expect(parent).not.toHaveTextContent('Duplicate / Duplicate'); expect(parent).not.toHaveTextContent(/^Duplicate$/)
})

it('submits a group save only once while the request is pending', async () => {
  const user = userEvent.setup(); let resolveSave: ((result: OperationResult) => void) | undefined
  const saveGroup = vi.fn(() => new Promise<OperationResult>((resolve) => { resolveSave = resolve }))
  render(<GroupsPanel groups={[]} bridge={{ saveGroup } as unknown as TaskBridge} onResult={() => undefined} />)
  await user.click(screen.getByRole('button', { name: 'New group' })); await user.type(screen.getByLabelText('Group name'), 'Operations')
  const save = screen.getByRole('button', { name: 'Save group' }); await user.dblClick(save)
  expect(saveGroup).toHaveBeenCalledTimes(1); expect(save).toBeDisabled()
  resolveSave?.({ action: 'save_group', outcome: 'accepted', message: 'Saved.' })
})
