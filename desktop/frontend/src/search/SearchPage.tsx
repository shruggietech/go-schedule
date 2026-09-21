import { useEffect, useMemo, useRef, useState } from 'react'
import { Button, DataTable, Dialog, Notice, StatePanel, StatusBadge } from '../components'
import type { Drilldown } from '../systems/SystemsPage'
import type { SearchAction, SearchActionResult, SearchBridge, SearchKind, SearchMatch, SearchSelection, SearchSnapshot } from './model'

const kinds: SearchKind[] = ['task', 'group', 'failure', 'schedule', 'alert']
const actions: Array<{ id: Exclude<SearchAction, 'open'>; label: string }> = [{ id: 'acknowledge', label: 'Acknowledge' }, { id: 'enable', label: 'Enable' }, { id: 'disable', label: 'Disable' }, { id: 'run_now', label: 'Run now' }]
const keyOf = (match: SearchMatch) => `${match.registrationKey}:${match.result.kind}:${match.result.object_id}`
const selectionOf = (match: SearchMatch): SearchSelection => ({ registrationKey: match.registrationKey, expectedDaemonId: match.expectedDaemonId, kind: match.result.kind, objectId: match.result.object_id, displayName: match.result.name })

export function SearchPage({ bridge, onOpen }: { bridge: SearchBridge; onOpen(intent: Drilldown): Promise<void> }) {
  const [query, setQuery] = useState('')
  const [selectedKinds, setSelectedKinds] = useState<SearchKind[]>(kinds)
  const [snapshot, setSnapshot] = useState<SearchSnapshot>()
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [selected, setSelected] = useState<Set<string>>(new Set())
  const [action, setAction] = useState<Exclude<SearchAction, 'open'>>('run_now')
  const [confirming, setConfirming] = useState(false)
  const [invoker, setInvoker] = useState<HTMLElement | null>(null)
  const [result, setResult] = useState<SearchActionResult>()
  const latest = useRef(0)
  const preserveDuringRefresh = useRef<SearchSnapshot | undefined>(undefined)
  const accept = (next: SearchSnapshot) => {
    if (next.generation < latest.current) return
    latest.current = next.generation
    const preserved = preserveDuringRefresh.current
    const merged = preserved && !next.complete ? { ...next, observations: next.observations.map((observation) => observation.state === 'connecting' ? preserved.observations.find((previous) => previous.registration.key === observation.registration.key) ?? observation : observation) } : next
    if (next.complete) preserveDuringRefresh.current = undefined
    setSnapshot(merged); setLoading(!next.complete)
  }
  useEffect(() => bridge.subscribe(accept), [bridge])
  const submit = async () => {
    if (!query.trim()) { setError('Enter a search term.'); return }
    setError(''); setSelected(new Set()); setLoading(true)
    try { accept(await bridge.search({ query: query.trim(), kinds: selectedKinds, limit: 50 })) }
    catch { setError('Registered schedulers could not be searched. Check the connections and try again.'); setLoading(false) }
  }
  const matches = useMemo(() => (snapshot?.observations ?? []).flatMap((observation) => observation.matches).sort((a, b) => a.result.kind.localeCompare(b.result.kind) || a.result.name.localeCompare(b.result.name) || a.sourceLabel.localeCompare(b.sourceLabel) || a.result.object_id.localeCompare(b.result.object_id)), [snapshot])
  const chosen = matches.filter((match) => selected.has(keyOf(match)))
  const compatible = chosen.length > 0 && chosen.every((match) => match.availableActions.includes(action))
  const open = (match: SearchMatch) => {
    const destination = match.result.kind === 'failure' || match.result.kind === 'alert' ? 'activity' : match.result.kind === 'schedule' ? 'schedule' : 'tasks'
    void onOpen({ registrationKey: match.registrationKey, expectedDaemonId: match.expectedDaemonId, label: `${match.sourceLabel} (${match.sourceShortId})`, destination, taskId: match.result.task_id, recordId: match.result.object_id, recordKind: match.result.kind, occurredAt: match.result.occurred_at, context: `${match.result.kind}: ${match.result.name}` })
  }
  const execute = async () => {
    setConfirming(false); setLoading(true)
    try { const outcome = await bridge.execute({ action, selections: chosen.map(selectionOf) }); setResult(outcome); preserveDuringRefresh.current = snapshot; setSelected(new Set()); await submit() }
    catch { setError('The action outcomes could not be confirmed. Refresh before retrying.'); setLoading(false) }
  }
  const grouped = useMemo(() => Object.entries(chosen.reduce<Record<string, SearchMatch[]>>((acc, match) => { (acc[`${match.sourceLabel} (${match.sourceShortId}, ${match.registrationKey})`] ??= []).push(match); return acc }, {})), [chosen])
  return <>
    <header className="page-header"><div><p className="eyebrow">Registered schedulers</p><h1>Search</h1><p>Find tasks, groups, failures, schedules, and alerts without changing the selected scheduler.</p></div></header>
    <form className="search-controls panel" onSubmit={(event) => { event.preventDefault(); void submit() }}>
      <label className="search-query">Search all systems<input value={query} maxLength={200} onChange={(event) => setQuery(event.target.value)} /></label>
      <fieldset><legend>Include</legend><div className="search-kind-list">{kinds.map((kind) => <label key={kind}><input type="checkbox" checked={selectedKinds.includes(kind)} onChange={(event) => setSelectedKinds((current) => event.target.checked ? [...current, kind] : current.filter((value) => value !== kind))} />{kind}</label>)}</div></fieldset>
      <Button type="submit" pending={loading} disabled={selectedKinds.length === 0}>Search</Button>
    </form>
    {error && <Notice title="Search needs attention" tone="error" identity={error}>{error}</Notice>}
    <p className="systems-announcement" role="status" aria-live="polite" aria-atomic="true">{loading ? `${snapshot?.observations.filter((item) => item.state !== 'connecting').length ?? 0} systems responded; search still running.` : snapshot ? `${snapshot.observations.length} systems searched; ${matches.length} results.${snapshot.observations.some((item) => item.truncated) ? ' Some schedulers have additional matches.' : ''}` : 'Enter a query to search every registered scheduler.'}</p>
    {(snapshot?.observations ?? []).some((item) => item.truncated) && <Notice title="Additional matches exist" tone="info">{snapshot?.observations.filter((item) => item.truncated).map((item) => item.registration.label).join(', ')} returned the per-scheduler result limit. Refine the query to narrow those results.</Notice>}
    {(snapshot?.observations ?? []).some((item) => item.failure) && <section className="search-targets panel" aria-label="Target search status">{snapshot?.observations.filter((item) => item.failure).map((item) => <div key={item.registration.key}><span>{item.registration.label} ({item.registration.shortDaemonId || 'identity unavailable'})</span><StatusBadge state={item.state} /><span>{item.failure?.message} {item.failure?.action}</span></div>)}</section>}
    {result && <Notice title="Action outcomes" tone={result.outcome === 'accepted' ? 'success' : 'warning'} identity={`${result.action}:${result.message}`}>{result.message}{result.outcomes.map((outcome) => <div key={`${outcome.registrationKey}:${outcome.objectId}`}>{outcome.displayName}: {outcome.message}</div>)}</Notice>}
    {snapshot && matches.length === 0 && snapshot.complete ? <StatePanel title="No matching automation" detail="Try a broader query or include more result types." /> : matches.length > 0 && <section className="panel search-results"><div className="section-heading"><h2>Results</h2><div className="actions"><label>Action<select value={action} onChange={(event) => setAction(event.target.value as Exclude<SearchAction, 'open'>)}>{actions.map((item) => <option key={item.id} value={item.id}>{item.label}</option>)}</select></label><Button variant="secondary" disabled={!compatible} onClick={(event) => { setInvoker(event.currentTarget); setConfirming(true) }}>Review {chosen.length || ''} selected</Button></div></div><DataTable caption="Cross-daemon search results" headings={['Select', 'Result', 'Type', 'Source', 'Observed context', 'Open']} rows={matches.map((match) => [<input aria-label={`Select ${match.result.name} from ${match.sourceLabel}`} type="checkbox" checked={selected.has(keyOf(match))} disabled={!match.availableActions.some((value) => value !== 'open')} title={match.disabledReason} onChange={(event) => setSelected((current) => { const next = new Set(current); if (event.target.checked) next.add(keyOf(match)); else next.delete(keyOf(match)); return next })} />, <strong>{match.result.name}</strong>, match.result.kind, <span>{match.sourceLabel}<small className="search-source-id">{match.sourceShortId}</small></span>, match.result.context || match.result.occurred_at || 'Current', <Button variant="secondary" onClick={() => open(match)}>Open</Button>])} /></section>}
    <Dialog open={confirming} title={`Confirm ${action.replace('_', ' ')}`} invoker={invoker} onClose={() => setConfirming(false)} actions={<><Button variant="secondary" onClick={() => setConfirming(false)}>Cancel</Button><Button variant="affirmative" onClick={() => void execute()}>Confirm {action.replace('_', ' ')}</Button></>}>{grouped.map(([source, values]) => <section key={source}><h3>{source}</h3><ul>{values.map((match) => <li key={keyOf(match)}>{match.result.name} ({match.result.kind})</li>)}</ul></section>)}<p>Each target is revalidated independently. Successful actions are not rolled back when another target fails.</p></Dialog>
  </>
}
