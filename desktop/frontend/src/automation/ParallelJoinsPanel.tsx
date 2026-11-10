import { useEffect, useRef, useState, type FormEvent } from 'react'
import { Button, Dialog, Field, Notice, StatePanel } from '../components'
import { blankJoin, type AutomationBridge, type AutomationWorkspace, type JoinDefinition, type JoinDraft, type JoinInstance, type OperationResult, type RecoveryCondition } from './model'

type Props = { workspace: AutomationWorkspace; bridge: AutomationBridge; disabled: boolean; available?: boolean; recoverySupported?: boolean; query: string; attention: boolean; scope: string; onResult(result: OperationResult): void; onRefresh(): void }
const unavailable = (action: string): OperationResult => ({ action, outcome: 'unavailable', message: 'This scheduler does not expose parallel join controls. Update the selected scheduler.' })
const recoveryConditions: Array<{ value: RecoveryCondition; label: string }> = [
  { value: 'failure', label: 'failure' }, { value: 'canceled', label: 'cancellation' },
  { value: 'skipped', label: 'overlap skip' }, { value: 'not_started', label: 'definitive no-start' },
]
const describeConditions = (values: RecoveryCondition[] | null | undefined) => values?.length ? values.map(v => recoveryConditions.find(c => c.value === v)?.label ?? v).join(', ') : 'No recovery causes selected'

export function ParallelJoinsPanel({ workspace, bridge, disabled, available = true, recoverySupported = true, query, attention, scope, onResult, onRefresh }: Props) {
  const [draft, setDraft] = useState<JoinDraft>()
  const [error, setError] = useState<OperationResult>()
  const [busy, setBusy] = useState(false)
  const lock = useRef(false)
  const current = useRef(scope); current.current = scope
  const mounted = useRef(true)
  const previousTarget = useRef(scope.split('/')[0])
  const [invoker, setInvoker] = useState<HTMLElement | null>(null)
  const [selected, setSelected] = useState('')
  const [instances, setInstances] = useState<JoinInstance[]>()
  const [exactID, setExactID] = useState('')
  const [offset, setOffset] = useState(0)
  useEffect(() => { mounted.current = true; return () => { mounted.current = false } }, [])
  useEffect(() => {
    const target = scope.split('/')[0]
    if (previousTarget.current !== target) { setDraft(undefined); previousTarget.current = target }
    setError(undefined); setExactID(''); setInstances(undefined); setSelected(''); setOffset(0)
  }, [scope])
  const joins = (workspace.joins ?? []).filter(v => [v.name, v.entry_task_name, v.target_task_name, v.recovery_task_name, v.cleanup_task_name, ...v.branches.map(b => b.name)].join(' ').toLowerCase().includes(query.toLowerCase()))
  const fetchedInstances = instances ?? workspace.joinInstances ?? []
  const visible = fetchedInstances.filter(v => (!selected || v.definition_id === selected) && (!attention || ['waiting', 'blocked', 'resolved', 'target_pending', 'recovery_pending', 'cleanup_pending'].includes(v.state) || [v.target_settlement, v.recovery_settlement, v.cleanup_settlement].some(s => s && s !== 'success')))
  const open = (value: JoinDefinition | undefined, target: HTMLElement) => {
    setInvoker(target); setError(undefined)
    setDraft(value ? {
      id: value.id, name: value.name, entryTaskId: value.entry_task_id, targetTaskId: value.target_task_id,
      recoveryTaskId: value.recovery_task_id ?? '', recoveryOn: value.recovery_on ?? (value.recovery_task_id ? ['failure'] : []), cleanupTaskId: value.cleanup_task_id ?? '',
      branches: value.branches.map(b => ({ ...b })), originalUpdatedAt: value.updated_at, isNew: false, overwriteStale: false,
    } : blankJoin())
  }
  const openEditable = async (value: JoinDefinition, target: HTMLElement) => {
    if (lock.current || !available || disabled) return
    lock.current = true; setBusy(true); setError(undefined)
    const captured = current.current
    try {
      const result = await (bridge.joinDefinition?.(value.id) ?? Promise.resolve(unavailable('load_join_definition')))
      if (!mounted.current || captured !== current.current || (result.daemonId && result.daemonId !== scope.split('/')[0])) return
      if (result.outcome === 'accepted' && result.join) open(result.join, target)
      else setError(result.outcome === 'accepted' ? { action: 'load_join_definition', outcome: 'unavailable', message: 'The scheduler did not return a complete editable definition. Refresh and try again.' } : result)
    } catch {
      if (mounted.current && captured === current.current) setError({ action: 'load_join_definition', outcome: 'unavailable', message: 'The editable definition could not be loaded. Refresh and try again.' })
    } finally { lock.current = false; if (mounted.current) setBusy(false) }
  }
  const run = async (action: () => Promise<OperationResult>, read = false) => {
    if (lock.current || !available || (!read && disabled)) return
    lock.current = true; setBusy(true)
    const captured = current.current
    try {
      const result = await action()
      if (!mounted.current || captured !== current.current || (result.daemonId && result.daemonId !== scope.split('/')[0])) return
      setError(result.outcome === 'accepted' ? undefined : result)
      if (read) { if (result.instances) setInstances(result.instances); if (result.instance) setInstances([result.instance]) }
      else {
        onResult(result)
        if (result.outcome === 'accepted') { setDraft(undefined); onRefresh() }
        else if (result.outcome === 'stale' && draft) setDraft({ ...draft, overwriteStale: false })
      }
    } catch {
      if (captured === current.current && mounted.current) {
        const result: OperationResult = { action: read ? 'load_join_instances' : 'save_join', outcome: read ? 'unavailable' : 'uncertain', message: read ? 'Join observations could not refresh.' : 'The response was not authoritative. Refresh this scheduler before deciding whether to retry.' }
        setError(result); if (!read) onResult(result)
      }
    } finally { lock.current = false; if (mounted.current) setBusy(false) }
  }
  const inspect = (id: string, next = 0) => { setSelected(id); setOffset(next); void run(() => bridge.joinInstances?.(id, next) ?? Promise.resolve(unavailable('load_join_instances')), true) }
  const submit = (event: FormEvent) => { event.preventDefault(); if (draft) void run(() => bridge.saveJoin?.(draft) ?? Promise.resolve(unavailable('save_join'))) }
  const options = workspace.tasks.map(t => <option key={t.id} value={t.id}>{t.name} ({t.id})</option>)
  return <section className="panel source-section" aria-labelledby="parallel-joins-heading">
    <div className="section-heading"><h2 id="parallel-joins-heading">Parallel joins</h2><Button disabled={disabled || busy} onClick={event => open(undefined, event.currentTarget)}>Create parallel join</Button></div>
    <p>A successful entry run creates these named invocations. The success target waits for every exact branch to succeed. Optional recovery and cleanup wait for actual settlement.</p>
    {error && !draft && <Notice title="Join observation" tone="error">{error.message}</Notice>}
    {joins.length ? <div className="source-grid">{joins.map(value => <article className="source-card" key={value.id}>
      <h3>{value.name}</h3><dl>
        <dt>Entry</dt><dd>{value.entry_task_name ? value.entry_task_name + ' (' + value.entry_task_id + ')' : value.entry_task_id}</dd>
        <dt>Required branches</dt><dd>{value.branches.map(branch => branch.name + ' (' + (branch.task_name ? branch.task_name + ' / ' + branch.task_id : branch.task_id) + ')').join(', ')}</dd>
        <dt>Success target</dt><dd>{value.target_task_name ? value.target_task_name + ' (' + value.target_task_id + ')' : value.target_task_id}</dd>
        {value.recovery_task_id && <><dt>Recovery</dt><dd>{value.recovery_task_name ?? value.recovery_task_id}: {describeConditions(value.recovery_on)}</dd></>}
        {value.cleanup_task_id && <><dt>Cleanup</dt><dd>{value.cleanup_task_name ?? value.cleanup_task_id}</dd></>}
      </dl><div className="actions">
        <Button variant="secondary" disabled={disabled || busy} onClick={event => void openEditable(value, event.currentTarget)}>Edit parallel join</Button>
        <Button variant="secondary" disabled={!available || busy} onClick={() => inspect(value.id)}>Inspect instances</Button>
        <Button variant="danger" disabled={disabled || busy} onClick={() => { if (window.confirm('Delete ' + value.name + '? Outstanding decisions will be resolved; unrelated runs are unchanged.')) void run(() => bridge.deleteJoin?.(value.id) ?? Promise.resolve(unavailable('delete_join'))) }}>Delete parallel join</Button>
      </div>
    </article>)}</div> : <StatePanel title="No matching parallel joins" detail="Create a block using existing tasks, or change the source filters." />}
    <details><summary>Join instances and exact evidence</summary>
      <form onSubmit={event => { event.preventDefault(); if (exactID.trim()) void run(() => bridge.joinInstance?.(exactID.trim()) ?? Promise.resolve(unavailable('load_join_instance')), true) }}>
        <Field label="Exact instance ID"><input value={exactID} onChange={event => setExactID(event.target.value)} /></Field><Button type="submit" disabled={!available || busy || !exactID.trim()}>Open exact instance</Button>
      </form>
      <div className="actions"><Button variant="secondary" disabled={!available || busy || offset === 0} onClick={() => inspect(selected, Math.max(0, offset - 100))}>Previous instances</Button><Button variant="secondary" disabled={!available || busy || fetchedInstances.length < 100} onClick={() => inspect(selected, offset + 100)}>Next instances</Button><Button variant="secondary" disabled={!available || busy} onClick={() => inspect(selected, offset)}>Refresh instances</Button></div>
      <p>This is a bounded page of 100 instances. An exact ID remains inspectable outside this page.</p>
      {visible.map(value => <JoinInstanceEvidence key={value.id} instance={value} />)}{visible.length === 0 && <p>No matching instance observations in this page.</p>}
    </details>
    <Dialog open={Boolean(draft)} title={draft?.isNew ? 'Create parallel join' : 'Edit parallel join'} invoker={invoker} onClose={() => { if (!busy) setDraft(undefined) }}>
      {draft && <form onSubmit={submit} aria-label="Parallel join editor">
        {error && <Notice title={error.outcome === 'stale' ? 'Newer definition found' : 'Review this join'} tone="error">{error.message}</Notice>}
        <Field label="Join name"><input required value={draft.name} onChange={event => setDraft({ ...draft, name: event.target.value })} /></Field>
        <Field label="Entry task"><select required value={draft.entryTaskId} onChange={event => setDraft({ ...draft, entryTaskId: event.target.value })}><option value="">Select task</option>{options}</select></Field>
        <fieldset><legend>Required named branches</legend>{draft.branches.map((branch, index) => <fieldset key={branch.id}><legend>Branch {index + 1}</legend>
          <Field label={'Branch ' + (index + 1) + ' name'}><input required value={branch.name} onChange={event => setDraft({ ...draft, branches: draft.branches.map(value => value.id === branch.id ? { ...value, name: event.target.value } : value) })} /></Field>
          <Field label={'Branch ' + (index + 1) + ' task'}><select required value={branch.task_id} onChange={event => setDraft({ ...draft, branches: draft.branches.map(value => value.id === branch.id ? { ...value, task_id: event.target.value } : value) })}><option value="">Select task</option>{options}</select></Field>
          <Button type="button" variant="secondary" disabled={disabled || busy || draft.branches.length === 1} onClick={() => setDraft({ ...draft, branches: draft.branches.filter(value => value.id !== branch.id) })}>Remove branch {index + 1}</Button>
        </fieldset>)}<Button type="button" variant="secondary" disabled={disabled || busy} onClick={() => setDraft({ ...draft, branches: [...draft.branches, { id: crypto.randomUUID(), name: 'Branch ' + (draft.branches.length + 1), task_id: '' }] })}>Add branch</Button></fieldset>
        <Field label="Joined target task"><select required value={draft.targetTaskId} onChange={event => setDraft({ ...draft, targetTaskId: event.target.value })}><option value="">Select task</option>{options}</select></Field>
        {!recoverySupported && <p>This scheduler supports success-only joins. Update it to configure recovery or cleanup.</p>}
        <fieldset disabled={disabled || busy || !recoverySupported}><legend>Optional recovery and cleanup</legend>
          <Field label="Recovery task" error={error?.field === 'recoveryTaskId' ? error.message : undefined}><select value={draft.recoveryTaskId} onChange={event => setDraft({ ...draft, recoveryTaskId: event.target.value, recoveryOn: event.target.value ? (draft.recoveryTaskId ? draft.recoveryOn : ['failure']) : [] })}><option value="">None</option>{options}</select></Field>
          <fieldset disabled={!draft.recoveryTaskId}><legend>Unsuccessful branch causes</legend>{recoveryConditions.map(condition => <label className="check-field" key={condition.value}><input type="checkbox" checked={draft.recoveryOn.includes(condition.value)} onChange={event => setDraft({ ...draft, recoveryOn: event.target.checked ? recoveryConditions.map(value => value.value).filter(value => value === condition.value || draft.recoveryOn.includes(value)) : draft.recoveryOn.filter(value => value !== condition.value) })} />Recover on {condition.label}</label>)}</fieldset>
          {draft.recoveryTaskId && <p>{describeConditions(draft.recoveryOn)}. Recovery waits for all branches to settle.</p>}
          <Field label="Cleanup task" error={error?.field === 'cleanupTaskId' ? error.message : undefined}><select value={draft.cleanupTaskId} onChange={event => setDraft({ ...draft, cleanupTaskId: event.target.value })}><option value="">None</option>{options}</select></Field>
        </fieldset>
        <ol aria-label="Routing order"><li>Wait for every exact branch to finish or definitively not start.</li><li>All successful: run the joined target. Otherwise run recovery only for a selected matching cause; no matching cause selects neither follow-up.</li><li>Cleanup waits for all branches and the selected follow-up to settle, including failure, cancellation or definitive no-start. A failed cancellation request is not termination.</li></ol>
        {error?.outcome === 'stale' && <label><input type="checkbox" checked={draft.overwriteStale} onChange={event => setDraft({ ...draft, overwriteStale: event.target.checked })} />Explicitly overwrite the newer definition</label>}
        <div className="actions"><Button type="submit" disabled={disabled || busy || (error?.outcome === 'stale' && !draft.overwriteStale)}>{busy ? 'Saving' : 'Save parallel join'}</Button><Button type="button" variant="secondary" disabled={busy} onClick={() => setDraft(undefined)}>Cancel</Button></div>
      </form>}
    </Dialog>
  </section>
}

export function JoinInstanceEvidence({ instance: value }: { instance: JoinInstance }) {
  const stages = [
    { label: 'Target', task: value.target_task_id, delivery: value.target_delivery_id, deliveryState: value.target_delivery_state, run: value.target_run_id, outcome: value.target_outcome, started: value.target_started_at, ended: value.target_ended_at, settlement: value.target_settlement, settled: value.target_settled_at, reason: value.target_reason, cancellation: value.target_cancellation_status, cancellationError: value.target_cancellation_error },
    { label: 'Recovery', task: value.recovery_task_id, delivery: value.recovery_delivery_id, deliveryState: value.recovery_delivery_state, run: value.recovery_run_id, outcome: value.recovery_outcome, started: value.recovery_started_at, ended: value.recovery_ended_at, settlement: value.recovery_settlement, settled: value.recovery_settled_at, reason: value.recovery_reason, cancellation: value.recovery_cancellation_status, cancellationError: value.recovery_cancellation_error },
    { label: 'Cleanup', task: value.cleanup_task_id, delivery: value.cleanup_delivery_id, deliveryState: value.cleanup_delivery_state, run: value.cleanup_run_id, outcome: value.cleanup_outcome, started: value.cleanup_started_at, ended: value.cleanup_ended_at, settlement: value.cleanup_settlement, settled: value.cleanup_settled_at, reason: value.cleanup_reason, cancellation: value.cleanup_cancellation_status, cancellationError: value.cleanup_cancellation_error },
  ]
  const route = !value.selected_route ? 'Waiting for branches to settle' : value.selected_route === 'none' ? 'No configured follow-up selected' : value.selected_route
  return <article className="panel"><h3>{value.definition_name}: {value.state.replaceAll('_', ' ')}</h3>
    <dl><dt>Instance ID</dt><dd><code>{value.id}</code></dd><dt>Definition ID</dt><dd><code>{value.definition_id}</code></dd><dt>Entry run</dt><dd><code>{value.entry_run_id}</code></dd><dt>Reason</dt><dd>{value.reason || 'No blocking reason reported.'}</dd><dt>Selected route</dt><dd>{route}</dd></dl>
    {stages.filter(stage => stage.task).map(stage => <section key={stage.label} aria-label={stage.label + ' exact evidence'}><h4>{stage.label}</h4><dl>
      <dt>{stage.label} task</dt><dd><code>{stage.task}</code></dd><dt>{stage.label} admission</dt><dd>{stage.delivery || 'Not admitted'} {stage.deliveryState || ''}</dd>
      <dt>{stage.label} run</dt><dd><code>{stage.run || 'No native run reported'}</code> {stage.started && !stage.ended ? 'running' : stage.outcome || ''}</dd>
      <dt>{stage.label} settlement</dt><dd>{stage.settlement || 'Not settled'}{stage.settled ? ' at ' + stage.settled : ''}</dd>
      {stage.reason && <><dt>{stage.label} reason</dt><dd>{stage.reason}</dd></>}
      {stage.cancellation && <><dt>{stage.label} cancellation</dt><dd>{stage.cancellation}{stage.cancellationError ? ': ' + stage.cancellationError : ''}</dd></>}
    </dl></section>)}
    <table><caption>{value.branches.length} of {value.branch_count} exact branch invocations</caption><thead><tr><th>Branch</th><th>Invocation</th><th>Task</th><th>Delivery</th><th>Exact run</th><th>Outcome or reason</th><th>Settlement</th></tr></thead><tbody>{value.branches.map(branch => <tr key={branch.id}>
      <td>{branch.name} ({branch.branch_id})</td><td><code>{branch.id}</code></td><td>{branch.task_name ? branch.task_name + ' (' + branch.task_id + ')' : branch.task_id}</td><td>{branch.delivery_id}: {branch.delivery_state}</td><td><code>{branch.run_id || 'Not started'}</code></td>
      <td>{branch.started_at && !branch.ended_at ? 'running' : branch.outcome || 'waiting'}{branch.reason ? ': ' + branch.reason : ''}{branch.cancellation_status ? ' (cancellation ' + branch.cancellation_status + ')' : ''}</td><td>{branch.settlement || 'Not settled'}{branch.settled_at ? ' at ' + branch.settled_at : ''}</td>
    </tr>)}</tbody></table>
  </article>
}
