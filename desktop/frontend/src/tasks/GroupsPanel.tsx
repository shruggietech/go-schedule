import { useRef, useState } from 'react'
import { Button, Dialog, Field, Notice } from '../components'
import type { GroupDraft, GroupSummary, OperationResult, TaskBridge } from './model'

type ConfirmedAction = { group: GroupSummary; kind: 'toggle' | 'delete'; invoker: HTMLElement }

export function GroupsPanel({ groups, bridge, available = true, targetName = 'This computer', onResult }: { groups: GroupSummary[]; bridge: TaskBridge; available?: boolean; targetName?: string; onResult(result: OperationResult): void }) {
  const empty: GroupDraft = { id: '', name: '', parentId: '', enabled: false, isNew: true, originalUpdatedAt: '', overwriteStale: false }
  const [draft, setDraft] = useState<GroupDraft | null>(null)
  const [confirmedAction, setConfirmedAction] = useState<ConfirmedAction | null>(null)
  const [saveResult, setSaveResult] = useState<OperationResult>()
  const [pending, setPending] = useState(false)
  const pendingRef = useRef(false)
  const save = async () => {
    if (!draft || pendingRef.current) return
    pendingRef.current = true
    setPending(true)
    try {
      const result = await bridge.saveGroup(draft)
      setSaveResult(result)
      onResult(result)
      if (result.outcome === 'accepted') setDraft(null)
    } finally {
      pendingRef.current = false
      setPending(false)
    }
  }
  const confirm = async () => {
    if (!confirmedAction) return
    const { group, kind } = confirmedAction
    const result = kind === 'delete' ? await bridge.deleteGroup(group.id) : await bridge.setGroupEnabled(group.id, !group.declaredEnabled)
    setConfirmedAction(null)
    onResult(result)
  }
  const editedPath = groups.find((group) => group.id === draft?.id)?.path
  const parents = groups.filter((group) => group.id !== draft?.id && !(editedPath && group.path.startsWith(editedPath + ' / ')))
  return <aside className="panel groups" aria-labelledby="groups-title">
    <div className="section-heading"><h2 id="groups-title">Groups</h2><Button variant="secondary" disabled={!available} onClick={() => { setSaveResult(undefined); setDraft(empty) }}>New group</Button></div>
    {groups.length === 0 ? <p>No groups yet.</p> : <ul className="group-tree">{groups.map((group) => <li key={group.id} style={{ paddingInlineStart: `${group.depth}rem` }}>
      <button className="group-name" onClick={() => { setSaveResult(undefined); setDraft({ id: group.id, name: group.name, parentId: group.parentId, enabled: group.declaredEnabled, isNew: false, originalUpdatedAt: group.updatedAt, overwriteStale: false }) }}>{group.path}</button>
      <span>{group.taskCount} tasks, {group.descendantCount} descendants</span><span>{group.effectiveEnabled ? 'Enabled' : group.effectiveReason}</span>
      <Button variant="quiet" disabled={!available} onClick={(event) => setConfirmedAction({ group, kind: 'toggle', invoker: event.currentTarget })}>{group.declaredEnabled ? 'Disable' : 'Enable'}</Button>
      <Button variant="quiet" disabled={!available} onClick={(event) => setConfirmedAction({ group, kind: 'delete', invoker: event.currentTarget })}>Delete</Button>
    </li>)}</ul>}
    {draft && <div className="group-editor"><h3>{draft.isNew ? 'Create group' : 'Edit group'}</h3>
      {saveResult && saveResult.outcome !== 'accepted' && <Notice title={saveResult.outcome === 'stale' ? 'Newer changes found' : 'Could not save group'} tone="error">{saveResult.message}{saveResult.outcome === 'stale' && <Button variant="secondary" onClick={() => setDraft({ ...draft, overwriteStale: true })}>Overwrite newer version</Button>}</Notice>}
      <Field label="Group name"><input value={draft.name} onChange={(event) => setDraft({ ...draft, name: event.target.value })} /></Field>
      <Field label="Parent group"><select value={draft.parentId} onChange={(event) => setDraft({ ...draft, parentId: event.target.value })}><option value="">Root</option>{parents.map((group) => <option key={group.id} value={group.id}>{group.path}</option>)}</select></Field>
      <label><input type="checkbox" checked={draft.enabled} onChange={(event) => setDraft({ ...draft, enabled: event.target.checked })} /> Enabled</label>
      <div className="actions"><Button disabled={!available || pending} onClick={() => void save()}>Save group</Button><Button variant="quiet" onClick={() => setDraft(null)}>Cancel</Button></div>
    </div>}
    <Dialog open={confirmedAction !== null} title={confirmedAction?.kind === 'delete' ? `Delete ${confirmedAction.group.name}?` : `${confirmedAction?.group.declaredEnabled ? 'Disable' : 'Enable'} ${confirmedAction?.group.name}?`} invoker={confirmedAction?.invoker ?? null} onClose={() => setConfirmedAction(null)}>
      <p>{confirmedAction?.kind === 'delete' ? `On ${targetName}, this removes ${confirmedAction.group.descendantCount} descendant groups and leaves assigned tasks ungrouped.` : `On ${targetName}, this changes effective scheduling for ${confirmedAction?.group.descendantCount ?? 0} descendant groups and their assigned tasks while preserving each declared setting.`}</p>
      <Button variant={confirmedAction?.kind === 'delete' ? 'danger' : 'primary'} onClick={() => void confirm()}>Confirm {confirmedAction?.kind === 'delete' ? 'delete' : confirmedAction?.group.declaredEnabled ? 'disable' : 'enable'}</Button>
    </Dialog>
  </aside>
}
