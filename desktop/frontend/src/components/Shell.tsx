import { useEffect, useRef, useState, type ReactNode } from 'react'
import { Button, Dialog, StatusBadge, ToastRegion } from '.'
import type { Appearance, ConnectionSnapshot, Route } from '../connection/model'

const routes: Array<{ id: Route; label: string }> = [
  { id: 'tasks', label: 'Tasks' }, { id: 'automation', label: 'Automation Sources' }, { id: 'schedule', label: 'Schedule' }, { id: 'activity', label: 'Activity' }, { id: 'notifications', label: 'Notifications' }, { id: 'agentAccess', label: 'Agent Access' }, { id: 'connections', label: 'Connections' }, { id: 'settings', label: 'Settings' },
]

export function Shell({ route, onRoute, appearance, onAppearance, appearancePending = false, connection, announcement, onRetry, onQuit, children }: { route: Route; onRoute(route: Route): void; appearance: Appearance; onAppearance(value: Appearance): void; appearancePending?: boolean; connection: ConnectionSnapshot; announcement: string; onRetry(): void; onQuit(): void; children: ReactNode }) {
  const [dialogOpen, setDialogOpen] = useState(false)
  const [invoker, setInvoker] = useState<HTMLElement | null>(null)
  const [systemDark, setSystemDark] = useState(() => window.matchMedia?.('(prefers-color-scheme: dark)').matches ?? false)
  const mainRef = useRef<HTMLElement>(null)
  useEffect(() => {
    const query = window.matchMedia?.('(prefers-color-scheme: dark)')
    if (!query) return
    const update = (event: MediaQueryListEvent) => setSystemDark(event.matches)
    setSystemDark(query.matches); query.addEventListener('change', update)
    return () => query.removeEventListener('change', update)
  }, [])
  const resolvedAppearance = appearance === 'system' ? systemDark ? 'dark' : 'light' : appearance
  const selectRoute = (next: Route) => { onRoute(next); requestAnimationFrame(() => mainRef.current?.focus()) }
  return <div className="app" data-appearance={appearance} data-resolved-appearance={resolvedAppearance}>
    <a className="skip-link" href="#main-content">Skip to main content</a>
    <nav className="rail" aria-label="Application"><div className="brand"><img src="/go-schedule-mark.svg" alt="" /><span>go-schedule</span></div><div className="nav-links">{routes.map((item) => <button key={item.id} aria-current={route === item.id ? 'page' : undefined} onClick={() => selectRoute(item.id)}>{item.label}</button>)}</div><Button variant="quiet" onClick={onQuit}>Exit</Button></nav>
    <div className="workspace"><header className="target-bar" aria-label="Target context"><div><strong>{connection.target.displayName}</strong><span>{connection.target.platform}{connection.target.architecture ? `/${connection.target.architecture}` : ''}{connection.target.version ? ` · ${connection.target.version}` : ''}</span>{connection.stale && <span role="status">Data may be stale</span>}</div><StatusBadge state={connection.state} /><Button variant="secondary" onClick={(event) => { setInvoker(event.currentTarget); setDialogOpen(true) }}>Connection details</Button></header>
      <main id="main-content" ref={mainRef} tabIndex={-1}>{children}</main>
      <footer className="preferences"><label>Appearance<select value={appearance} disabled={appearancePending} onChange={(event) => onAppearance(event.target.value as Appearance)}><option value="system">System</option><option value="light">Light</option><option value="dark">Dark</option></select></label>{connection.action && <Button variant="secondary" onClick={onRetry}>Try again</Button>}</footer>
    </div>
    <Dialog open={dialogOpen} title={`${connection.target.displayName} connection`} invoker={invoker} onClose={() => setDialogOpen(false)}><p>{connection.message}</p><dl>{connection.target.endpoint && <><dt>Endpoint</dt><dd>{connection.target.endpoint}</dd></>}<dt>Platform</dt><dd>{connection.target.platform}{connection.target.architecture ? `/${connection.target.architecture}` : ''}</dd><dt>Capabilities</dt><dd>{connection.target.capabilities.length ? connection.target.capabilities.join(', ') : 'Unavailable'}</dd><dt>Data freshness</dt><dd>{connection.stale ? 'Last-known data, refresh pending' : 'Current'}</dd><dt>Last successful connection</dt><dd>{connection.lastSuccessfulAt || 'Not available'}</dd>{connection.recovery === 'automatic' && <><dt>Automatic retry</dt><dd>Attempt {connection.retryAttempt || 1}{connection.nextRetryAt ? ` scheduled for ${connection.nextRetryAt}` : ' scheduled'}</dd></>}<dt>Recovery</dt><dd>{connection.recovery === 'automatic' ? 'Automatic' : connection.recovery === 'manual' ? 'Manual action required' : 'None'}</dd></dl></Dialog>
    <ToastRegion message={announcement} />
  </div>
}
