import { useEffect, useMemo, useRef, useState } from 'react'
import { Button, Notice, StatePanel, StatusBadge } from '../components'
import type { DesktopBridge, Route, SystemObservation, SystemsSnapshot } from '../connection/model'

type SortKey = 'label' | 'state' | 'freshness' | 'upcoming' | 'failures' | 'alerts' | 'notifications'
export interface Drilldown { registrationKey: string; label: string; destination: Extract<Route, 'tasks' | 'schedule' | 'activity' | 'notifications'>; taskId?: string; recordId?: string; context: string }

const severity: Record<string, number> = { connected: 0, degraded: 1, recovering: 2, timed_out: 3, unavailable: 4, access_denied: 5, unauthorized: 6, forbidden: 7, revoked: 8, incompatible: 9, trust_changed: 10, identity_changed: 11, connecting: 12 }
const dateValue = (value?: string) => value ? Date.parse(value) || 0 : 0
const needsAttention = (item: SystemObservation) => item.state !== 'connected' || item.stale || (item.summary?.recent_failure_count ?? 0) > 0 || (item.summary?.unacknowledged_alert_count ?? 0) > 0 || (item.summary?.notification_problem_count ?? 0) > 0
const when = (value?: string) => value ? new Date(value).toLocaleString() : 'None'

export function SystemsPage({ bridge, onOpen }: { bridge: DesktopBridge; onOpen(intent: Drilldown): Promise<void> }) {
  const [snapshot, setSnapshot] = useState<SystemsSnapshot>()
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [query, setQuery] = useState('')
  const [state, setState] = useState('all')
  const [attention, setAttention] = useState(false)
  const [sort, setSort] = useState<SortKey>('state')
  const latestGeneration = useRef(0)
  const requestSequence = useRef(0)
  const accept = (next: SystemsSnapshot) => {
    if (next.generation >= latestGeneration.current) { latestGeneration.current = next.generation; setSnapshot(next) }
  }
  const refresh = async () => {
    const request = ++requestSequence.current
    setLoading(true); setError('')
    try {
      const next = await bridge.allSystems?.() ?? { generation: 0, startedAt: '', completedAt: '', complete: true, observations: [] }
      accept(next)
    }
    catch { if (request === requestSequence.current) setError('All Systems could not be refreshed. Check the desktop connection and try again.') }
    finally { if (request === requestSequence.current) setLoading(false) }
  }
  useEffect(() => {
    const unsubscribe = bridge.subscribeSystems?.(accept) ?? (() => undefined)
    void refresh()
    return unsubscribe
  }, [bridge])
  const observations = useMemo(() => {
    const needle = query.trim().toLocaleLowerCase()
    return [...(snapshot?.observations ?? [])].filter((item) => {
      const registration = item.registration
      const text = [registration.label, registration.endpoint, registration.daemonId, registration.platform, registration.version].join(' ').toLocaleLowerCase()
      return (!needle || text.includes(needle)) && (state === 'all' || item.state === state) && (!attention || needsAttention(item))
    }).sort((a, b) => {
      const byLabel = a.registration.label.localeCompare(b.registration.label) || a.registration.key.localeCompare(b.registration.key)
      if (sort === 'label') return byLabel
      if (sort === 'state') return (severity[a.state] ?? 99) - (severity[b.state] ?? 99) || byLabel
      if (sort === 'freshness') return dateValue(b.observedAt) - dateValue(a.observedAt) || byLabel
      if (sort === 'upcoming') return (dateValue(a.summary?.next_occurrence?.scheduled_for) || Number.MAX_SAFE_INTEGER) - (dateValue(b.summary?.next_occurrence?.scheduled_for) || Number.MAX_SAFE_INTEGER) || byLabel
      const field = sort === 'failures' ? 'recent_failure_count' : sort === 'alerts' ? 'unacknowledged_alert_count' : 'notification_problem_count'
      return (b.summary?.[field] ?? 0) - (a.summary?.[field] ?? 0) || byLabel
    })
  }, [attention, query, snapshot, sort, state])
  const states = [...new Set((snapshot?.observations ?? []).map((item) => item.state))].sort()
  const open = (item: SystemObservation, destination: Drilldown['destination'], taskId: string | undefined, recordId: string | undefined, context: string) => void onOpen({ registrationKey: item.registration.key, label: item.registration.label, destination, taskId, recordId, context })
  return <>
    <header className="page-header"><div><p className="eyebrow">Registered schedulers</p><h1>All Systems</h1><p>Current, bounded observations from independent daemons. Refreshing here never changes the selected scheduler.</p></div><Button pending={loading} onClick={() => void refresh()}>Refresh all</Button></header>
    {error && <Notice title="Refresh failed" tone="error" identity={error}>{error}</Notice>}
    <div className="systems-filters panel" aria-label="Filter and sort systems">
      <label>Search<input value={query} onChange={(event) => setQuery(event.target.value)} /></label>
      <label>Connection state<select value={state} onChange={(event) => setState(event.target.value)}><option value="all">All states</option>{states.map((value) => <option key={value} value={value}>{value.replaceAll('_', ' ')}</option>)}</select></label>
      <label>Sort by<select value={sort} onChange={(event) => setSort(event.target.value as SortKey)}><option value="state">State severity</option><option value="label">Display label</option><option value="freshness">Freshness</option><option value="upcoming">Upcoming work</option><option value="failures">Recent failures</option><option value="alerts">Alerts</option><option value="notifications">Notification problems</option></select></label>
      <label className="check-field"><input type="checkbox" checked={attention} onChange={(event) => setAttention(event.target.checked)} />Needs attention only</label>
    </div>
    <p className="systems-announcement" role="status" aria-live="polite" aria-atomic="true">{loading ? 'Refreshing registered systems.' : `${snapshot?.observations.length ?? 0} registered systems refreshed; ${observations.length} shown.`}</p>
    {loading && !snapshot ? <StatePanel title="Refreshing registered systems" detail="Each daemon has an independent five-second deadline." busy /> : observations.length === 0 ? <StatePanel title="No matching systems" detail="Change the filters or register a remote scheduler from Connections." /> : <div className="systems-grid">{observations.map((item) => {
      const summary = item.summary
      return <article className={`panel system-card ${needsAttention(item) ? 'system-attention' : ''}`} key={item.registration.key}>
        <header><div><h2>{item.registration.label}</h2><p>{item.registration.kind === 'local' ? 'This computer' : item.registration.endpoint} · {item.registration.shortDaemonId || 'Identity unavailable'}</p></div><StatusBadge state={item.state} /></header>
        <p className="system-meta">{[item.registration.platform, item.registration.architecture].filter(Boolean).join('/') || 'Platform unavailable'}{item.registration.version ? ` · ${item.registration.version}` : ''} · {item.stale ? `Stale since ${when(item.observedAt)}` : `Observed ${when(item.observedAt)}`}</p>
        {item.failure && <Notice title={item.stale ? 'Showing last session result' : 'System unavailable'} tone="warning" identity={`${item.registration.key}:${item.state}`}>{item.failure.message} {item.failure.action}</Notice>}
        {summary ? <dl className="system-metrics"><div><dt>Active tasks</dt><dd>{summary.active_task_count}</dd></div><div><dt>Next 24 hours</dt><dd>{summary.next_occurrence ? `${summary.next_occurrence.task_name}, ${when(summary.next_occurrence.scheduled_for)}` : 'None'}</dd></div><div><dt>Recent failures</dt><dd>{summary.recent_failure_count}</dd></div><div><dt>Unacknowledged alerts</dt><dd>{summary.unacknowledged_alert_count}</dd></div><div><dt>Notification problems</dt><dd>{summary.notification_problem_count}</dd></div></dl> : <p>No operational values are available yet.</p>}
        <div className="actions">
          <Button variant="secondary" onClick={() => open(item, 'tasks', summary?.next_occurrence?.task_id, undefined, summary?.next_occurrence ? `Upcoming task: ${summary.next_occurrence.task_name}` : 'Task list')}>Tasks</Button>
          <Button variant="secondary" onClick={() => open(item, 'schedule', summary?.next_occurrence?.task_id, undefined, summary?.next_occurrence ? `Upcoming at ${when(summary.next_occurrence.scheduled_for)}` : 'Schedule')}>Schedule</Button>
          {summary?.recent_failure && <Button variant="secondary" onClick={() => open(item, 'activity', summary.recent_failure?.task_id, summary.recent_failure?.run_id, `Failed run ${summary.recent_failure?.run_id}`)}>Failed run</Button>}
          {summary?.unacknowledged_alert && <Button variant="secondary" onClick={() => open(item, 'activity', summary.unacknowledged_alert?.task_id, summary.unacknowledged_alert?.alert_id, `Alert ${summary.unacknowledged_alert?.alert_id}`)}>Alert</Button>}
          {!summary?.recent_failure && !summary?.unacknowledged_alert && <Button variant="secondary" onClick={() => open(item, 'activity', undefined, undefined, 'Activity')}>Activity</Button>}
          <Button variant="secondary" onClick={() => open(item, item.registration.kind === 'remote' ? 'activity' : 'notifications', summary?.notification_problem?.task_id, summary?.notification_problem?.delivery_id, summary?.notification_problem ? `Notification delivery ${summary.notification_problem.delivery_id}` : 'Notifications')}>{item.registration.kind === 'remote' ? 'Notification issue' : 'Notifications'}</Button>
        </div>
      </article>
    })}</div>}
  </>
}
