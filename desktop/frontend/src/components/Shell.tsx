import { useRef, useState, type ReactNode } from 'react'
import { Button, Dialog, StatusBadge, ToastRegion } from '.'
import type { Appearance, ConnectionSnapshot, Route } from '../connection/model'

const routes: Array<{ id: Route; label: string }> = [
  { id: 'tasks', label: 'Tasks' }, { id: 'automation', label: 'Automation Sources' }, { id: 'schedule', label: 'Schedule' }, { id: 'activity', label: 'Activity' }, { id: 'connections', label: 'Connections' }, { id: 'settings', label: 'Settings' },
]

export function Shell({ route, onRoute, appearance, onAppearance, connection, announcement, onRetry, onQuit, children }: { route: Route; onRoute(route: Route): void; appearance: Appearance; onAppearance(value: Appearance): void; connection: ConnectionSnapshot; announcement: string; onRetry(): void; onQuit(): void; children: ReactNode }) {
  const [dialogOpen, setDialogOpen] = useState(false)
  const [invoker, setInvoker] = useState<HTMLElement | null>(null)
  const mainRef = useRef<HTMLElement>(null)
  const selectRoute = (next: Route) => { onRoute(next); requestAnimationFrame(() => mainRef.current?.focus()) }
  return <div className="app" data-appearance={appearance}>
    <a className="skip-link" href="#main-content">Skip to main content</a>
    <nav className="rail" aria-label="Application"><div className="brand"><img src="/go-schedule-mark.svg" alt="" /><span>go-schedule</span></div><div className="nav-links">{routes.map((item) => <button key={item.id} aria-current={route === item.id ? 'page' : undefined} onClick={() => selectRoute(item.id)}>{item.label}</button>)}</div><Button variant="quiet" onClick={onQuit}>Exit</Button></nav>
    <div className="workspace"><header className="target-bar" aria-label="Target context"><div><strong>{connection.target.displayName}</strong><span>{connection.target.platform}{connection.target.version ? ` · ${connection.target.version}` : ''}</span></div><StatusBadge state={connection.state} /><Button variant="secondary" onClick={(event) => { setInvoker(event.currentTarget); setDialogOpen(true) }}>Connection details</Button></header>
      <main id="main-content" ref={mainRef} tabIndex={-1}>{children}</main>
      <footer className="preferences"><label>Appearance<select value={appearance} onChange={(event) => onAppearance(event.target.value as Appearance)}><option value="system">System</option><option value="light">Light</option><option value="dark">Dark</option></select></label>{connection.action && <Button variant="secondary" onClick={onRetry}>Try again</Button>}</footer>
    </div>
    <Dialog open={dialogOpen} title="This computer connection" invoker={invoker} onClose={() => setDialogOpen(false)}><p>{connection.message}</p><dl><dt>Platform</dt><dd>{connection.target.platform}</dd><dt>Capabilities</dt><dd>{connection.target.capabilities.length ? connection.target.capabilities.join(', ') : 'Unavailable'}</dd></dl></Dialog>
    <ToastRegion message={announcement} />
  </div>
}
