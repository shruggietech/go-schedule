import { useCallback, useEffect, useId, useRef, useState } from 'react'
import { Appearance, Condition, Page, fixtureFor, loadSnapshot, showAbout, subscribeToProofEvents, type ProofSnapshot } from './proof'

const pages: Array<{ id: Page; label: string }> = [
  { id: 'tasks', label: 'Tasks' },
  { id: 'editor', label: 'Task editor' },
  { id: 'schedule', label: 'Schedule' },
  { id: 'activity', label: 'Activity' },
  { id: 'targets', label: 'Targets' },
]

const conditions: Condition[] = ['connected', 'new', 'empty', 'loading', 'disconnected', 'degraded', 'destructive', 'validation-error', 'success']

function Status({ condition }: { condition: Condition }) {
  const label = condition === 'validation-error' ? 'Needs correction' : condition[0].toUpperCase() + condition.slice(1)
  return <span className={`status status-${condition}`}><span aria-hidden="true" className="status-shape" />{label}</span>
}

function Tasks({ snapshot, condition, onEdit, onCancel, onRetry }: { snapshot: ProofSnapshot; condition: Condition; onEdit(): void; onCancel(): void; onRetry(): void }) {
  if (condition === 'loading') return <StatePanel title="Loading tasks" detail="Connecting to this computer and reading scheduled work." action="Cancel" onAction={onCancel} />
  if (condition === 'disconnected') return <StatePanel title="This computer is offline" detail={snapshot.health.message} action="Try again" onAction={onRetry} />
  if (condition === 'empty' || condition === 'new') return <StatePanel title="No tasks yet" detail="Create a task to run recurring work on this computer." action="Create task" onAction={onEdit} />
  return <div className="split"><section aria-labelledby="task-list-title" className="panel list-panel"><div className="panel-heading"><div><p className="eyebrow">3 tasks</p><h2 id="task-list-title">Scheduled work</h2></div><button className="quiet-button">Filter</button></div><div className="task-list">{snapshot.tasks.map((task, index) => <button className={`task-row ${index === 0 ? 'selected' : ''}`} key={task.id} onClick={onEdit}><span><strong>{task.name}</strong><small>{task.schedule}</small></span><span><Status condition={task.lastResult === 'failed' ? 'degraded' : task.enabled ? 'connected' : 'empty'} /><small>{task.nextRun || 'Paused'}</small></span></button>)}</div></section><aside aria-labelledby="task-detail-title" className="panel inspector"><p className="eyebrow">Selected task</p><h2 id="task-detail-title">Nightly backup</h2><dl><div><dt>Schedule</dt><dd>Every day at 23:00</dd></div><div><dt>Next run</dt><dd>Today, 23:00</dd></div><div><dt>Last result</dt><dd><Status condition="success" /></dd></div></dl><button className="secondary-button" onClick={onEdit}>Edit task</button></aside></div>
}

function Editor({ condition }: { condition: Condition }) {
  const errorId = useId()
  return <section className="panel form-panel" aria-labelledby="editor-title"><div className="panel-heading"><div><p className="eyebrow">Runs on This computer</p><h2 id="editor-title">Task details</h2></div>{condition === 'success' && <Status condition="success" />}</div><form onSubmit={(event) => event.preventDefault()}><label>Task name<input defaultValue={condition === 'new' ? '' : 'Nightly backup'} aria-invalid={condition === 'validation-error'} aria-describedby={condition === 'validation-error' ? errorId : undefined} /></label>{condition === 'validation-error' && <p className="field-error" id={errorId}>Enter a task name before saving.</p>}<label>Command<input className="mono" defaultValue="gosched backup --all" /></label><div className="form-grid"><label>Schedule<input defaultValue="Every day at 23:00" /></label><label>Time zone<select defaultValue="local"><option value="local">Local time</option><option value="utc">UTC</option></select></label></div><label className="check"><input type="checkbox" defaultChecked />Enable this task after saving</label><div className="actions"><button type="button" className="quiet-button">Cancel</button><button type="submit" className="primary-button">Save task</button></div></form></section>
}

function ScheduleView() {
  return <section className="panel" aria-labelledby="schedule-title"><div className="panel-heading"><div><p className="eyebrow">September 7 to 13</p><h2 id="schedule-title">Upcoming work</h2></div><button className="quiet-button">Today</button></div><ol className="timeline"><li><time>Today, 23:00</time><strong>Nightly backup</strong><span>Expected duration 12 minutes</span></li><li><time>Monday, 08:30</time><strong>Publish weekly reports</strong><span>Runs after source refresh</span></li><li><time>Friday, 18:00</time><strong>Clear temporary files</strong><span>Paused</span></li></ol></section>
}

function Activity({ condition }: { condition: Condition }) {
  return <section className="panel" aria-labelledby="activity-title"><div className="panel-heading"><div><p className="eyebrow">Latest first</p><h2 id="activity-title">Run history</h2></div><button className="quiet-button">Export</button></div>{condition === 'degraded' && <Notice title="Live updates delayed" detail="History is still available. Reconnect to receive new results." />}<div className="activity-list"><article><Status condition="success" /><div><h3>Nightly backup</h3><p>Completed today at 02:12 in 11m 42s</p></div><button className="quiet-button">View output</button></article><article><Status condition="degraded" /><div><h3>Publish weekly reports</h3><p>Exited with status 1 yesterday at 08:31</p></div><button className="quiet-button">View output</button></article></div></section>
}

function Targets({ snapshot }: { snapshot: ProofSnapshot }) {
  return <section className="panel" aria-labelledby="targets-title"><div className="panel-heading"><div><p className="eyebrow">Execution context</p><h2 id="targets-title">Choose a target</h2></div></div><button className="target-card selected"><span className="target-mark" aria-hidden="true">L</span><span><strong>{snapshot.target.displayName}</strong><small>{snapshot.target.detail}</small></span><Status condition={snapshot.target.connection as Condition} /><span className="selection-label">Current target</span></button><Notice title="Mutations stay explicit" detail="Create, edit, run, and delete actions always name their destination before confirmation." /></section>
}

function Notice({ title, detail }: { title: string; detail: string }) {
  return <div className="notice" role="status"><strong>{title}</strong><span>{detail}</span></div>
}

function StatePanel({ title, detail, action, onAction }: { title: string; detail: string; action: string; onAction?(): void }) {
  return <section className="panel state-panel"><span className="state-symbol" aria-hidden="true">○</span><h2>{title}</h2><p>{detail}</p><button className="primary-button" onClick={onAction}>{action}</button></section>
}

export function App({ snapshotLoader = loadSnapshot }: { snapshotLoader?: typeof loadSnapshot }) {
  const [page, setPage] = useState<Page>('tasks')
  const [condition, setCondition] = useState<Condition>('connected')
  const [appearance, setAppearance] = useState<Appearance>('light')
  const [snapshot, setSnapshot] = useState(() => fixtureFor('connected'))
  const [announcement, setAnnouncement] = useState('')
  const [targetDialogOpen, setTargetDialogOpen] = useState(false)
  const mainRef = useRef<HTMLElement>(null)
  const targetButtonRef = useRef<HTMLButtonElement>(null)
  const dialogCloseRef = useRef<HTMLButtonElement>(null)

  const applySnapshot = useCallback((loaded: ProofSnapshot) => {
    setSnapshot(loaded)
    setCondition(loaded.target.connection === 'connected' && loaded.tasks.length === 0 ? 'empty' : loaded.target.connection)
  }, [])
  const refreshSnapshot = useCallback(() => {
    void snapshotLoader().then(applySnapshot).catch(() => applySnapshot(fixtureFor('disconnected')))
  }, [applySnapshot, snapshotLoader])
  const applyPrototypeCondition = (next: Condition) => {
    setCondition(next)
    setSnapshot(fixtureFor(next))
  }

  useEffect(refreshSnapshot, [refreshSnapshot])
  useEffect(() => subscribeToProofEvents((event) => setAnnouncement(event.message)), [])
  useEffect(() => {
    if (!targetDialogOpen) return
    dialogCloseRef.current?.focus()
    const closeOnEscape = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        setTargetDialogOpen(false)
        targetButtonRef.current?.focus()
      }
    }
    window.addEventListener('keydown', closeOnEscape)
    return () => window.removeEventListener('keydown', closeOnEscape)
  }, [targetDialogOpen])

  const selectPage = (next: Page) => { setPage(next); requestAnimationFrame(() => mainRef.current?.focus()) }
  const title = pages.find((candidate) => candidate.id === page)!.label
  const primary = page === 'tasks' ? 'Create task' : page === 'targets' ? 'Switch target' : 'View tasks'
  const runPrimaryAction = () => {
    if (page === 'tasks') selectPage('editor')
    else if (page === 'targets') setTargetDialogOpen(true)
    else selectPage('tasks')
  }

  return <div className="app" data-appearance={appearance}>
    <a className="skip-link" href="#main-content">Skip to main content</a>
    <nav aria-label="Application" className="rail"><div className="brand"><img src="/go-schedule-mark.svg" alt="" /><span>go-schedule</span></div><div className="nav-links">{pages.map((item) => <button aria-current={page === item.id ? 'page' : undefined} key={item.id} onClick={() => selectPage(item.id)}><span className="nav-glyph" aria-hidden="true">{item.label[0]}</span><span>{item.label}</span></button>)}</div><button className="about" onClick={() => void showAbout()}>About</button></nav>
    <div className="workspace"><header className="target-bar" aria-label="Target context"><div><span className="target-dot" aria-hidden="true" /><strong>{snapshot.target.displayName}</strong><span>{snapshot.target.platform}</span></div><Status condition={snapshot.target.connection as Condition} /><button ref={targetButtonRef} onClick={() => setTargetDialogOpen(true)}>Switch target</button></header>
      <main id="main-content" ref={mainRef} tabIndex={-1}><header className="page-header"><div><p className="eyebrow">This computer</p><h1>{title}</h1><p>{page === 'tasks' ? 'Review, create, and maintain scheduled work.' : page === 'editor' ? 'Define what runs and when it runs.' : page === 'schedule' ? 'See when enabled work will run next.' : page === 'activity' ? 'Understand recent results and investigate failures.' : 'Keep execution context visible and deliberate.'}</p></div>{page !== 'editor' && <button className="primary-button" onClick={runPrimaryAction}>{primary}</button>}</header>
        {condition === 'destructive' && <Notice title="Delete Nightly backup?" detail="This removes its schedule from This computer. Run history remains available." />}
        {page === 'tasks' && <Tasks snapshot={snapshot} condition={condition} onEdit={() => selectPage('editor')} onCancel={() => applyPrototypeCondition('connected')} onRetry={refreshSnapshot} />}
        {page === 'editor' && <Editor condition={condition} />}
        {page === 'schedule' && <ScheduleView />}
        {page === 'activity' && <Activity condition={condition} />}
        {page === 'targets' && <Targets snapshot={snapshot} />}
      </main>
    </div>
    {targetDialogOpen && <div className="dialog-backdrop"><section aria-labelledby="switch-target-title" aria-modal="true" className="dialog" role="dialog"><div className="panel-heading"><div><p className="eyebrow">Execution context</p><h2 id="switch-target-title">Switch target</h2></div><button ref={dialogCloseRef} className="quiet-button" onClick={() => { setTargetDialogOpen(false); targetButtonRef.current?.focus() }}>Close</button></div><p>Choose where subsequent task changes will run.</p><button className="target-card selected"><span className="target-mark" aria-hidden="true">L</span><span><strong>This computer</strong><small>Local scheduler service</small></span><Status condition="connected" /><span className="selection-label">Current target</span></button></section></div>}
    <aside className="proof-controls" aria-label="Prototype controls"><label>State<select value={condition} onChange={(event) => applyPrototypeCondition(event.target.value as Condition)}>{conditions.map((item) => <option key={item}>{item}</option>)}</select></label><label>Appearance<select value={appearance} onChange={(event) => setAppearance(event.target.value as Appearance)}><option value="system">System</option><option value="light">Light</option><option value="dark">Dark</option></select></label></aside>
    <div className="sr-only" aria-live="polite">{announcement}</div>
  </div>
}
