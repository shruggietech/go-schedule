import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { expect, it } from 'vitest'
import type { GroupSummary, TaskBridge } from './model'
import { GroupsPanel } from './GroupsPanel'

it('uses full paths and excludes the edited group and descendants from parent choices', async () => {
  const user = userEvent.setup(); const groups: GroupSummary[] = [{ id: 'root', name: 'Duplicate', parentId: '', path: 'Duplicate', depth: 0, declaredEnabled: true, effectiveEnabled: true, effectiveReason: '', childCount: 1, taskCount: 0, descendantCount: 1, updatedAt: 'now' }, { id: 'child', name: 'Duplicate', parentId: 'root', path: 'Duplicate / Duplicate', depth: 1, declaredEnabled: true, effectiveEnabled: true, effectiveReason: '', childCount: 0, taskCount: 0, descendantCount: 0, updatedAt: 'now' }]
  render(<GroupsPanel groups={groups} bridge={{} as TaskBridge} onResult={() => undefined} />); await user.click(screen.getByRole('button', { name: 'Duplicate' })); const parent = screen.getByLabelText('Parent group'); expect(parent).not.toHaveTextContent('Duplicate / Duplicate'); expect(parent).not.toHaveTextContent(/^Duplicate$/)
})
