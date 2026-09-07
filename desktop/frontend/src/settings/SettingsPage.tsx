import { Button, Notice, StatePanel } from '../components'
import type { Appearance } from '../connection/model'
import type { SettingsWorkspace } from './model'

const transitionCopy: Record<string, { title: string; detail: string }> = {
  migrated: { title: 'Appearance migrated', detail: 'Your previous light, dark, or system choice was carried into the new desktop.' },
  not_found: { title: 'System appearance selected', detail: 'No earlier desktop appearance was found, so go-schedule follows your operating system.' },
  invalid: { title: 'Earlier preference could not be used', detail: 'The previous appearance value was invalid. go-schedule uses system appearance and remains usable.' },
  unreadable: { title: 'Earlier preference could not be read', detail: 'The previous preference file was unavailable. go-schedule uses system appearance and remains usable.' },
  not_required: { title: 'Desktop defaults restored', detail: 'The current desktop preference file is authoritative and follows your operating system.' },
}

export function SettingsPage({ workspace, message, pending, onAppearance, onRestore, onCopy, onOpen, onConnections }: { workspace?: SettingsWorkspace; message: string; pending: boolean; onAppearance(value: Appearance): void; onRestore(): void; onCopy(id: string): void; onOpen(key: string): void; onConnections(): void }) {
  if (!workspace) return <><PageHeader /><StatePanel title="Desktop settings are unavailable" detail={message || 'Check access to the user configuration directory, then try again.'} busy={pending} action="Restore desktop defaults" onAction={onRestore} /></>
  const transition = transitionCopy[workspace.preferences.transition.status] ?? transitionCopy.not_found
  return <><PageHeader />
    <div className="settings-grid">
      <section className="panel settings-section" aria-labelledby="appearance-heading"><div><p className="eyebrow">Desktop preference</p><h2 id="appearance-heading">Appearance</h2><p>Choose a durable color mode for this desktop.</p></div><label className="settings-control">Appearance<select value={workspace.preferences.appearance} disabled={pending} onChange={(event) => onAppearance(event.target.value as Appearance)}><option value="system">System</option><option value="light">Light</option><option value="dark">Dark</option></select></label><Button variant="secondary" disabled={pending} onClick={onRestore}>Restore desktop defaults</Button><p className="path-text">Preference file: {workspace.preferencePath}</p></section>
      <section className="panel settings-section" aria-labelledby="transition-heading"><div><p className="eyebrow">Upgrade</p><h2 id="transition-heading">{transition.title}</h2><p>{transition.detail}</p></div><p>Legacy font selection and scroll sensitivity were retired because browser typography and native scrolling now provide those behaviors.</p></section>
    </div>
    {!workspace.daemonAvailable && <Notice title="Daemon storage details are unavailable" tone="warning">Local preferences and product information are still available. Daemon-owned paths are not guessed. <Button variant="secondary" onClick={onConnections}>Open Connections</Button></Notice>}
    <section className="panel settings-section" aria-labelledby="storage-heading"><div><p className="eyebrow">Local files</p><h2 id="storage-heading">Application storage</h2><p>Paths are read-only. Ownership and removal behavior remain explicit, including locations configured outside application-owned data.</p></div><div className="storage-list">{workspace.storage.map((record) => <article className={`storage-record storage-${record.existence}`} key={record.id}><div className="storage-heading"><h3>{record.label}</h3><span className="status-label">{record.existence === 'present' ? 'Present' : record.existence === 'absent' ? 'Not present' : 'Unavailable'}</span></div><p className="path-text">{record.path || 'Path unavailable'}</p><dl><dt>Owner</dt><dd>{record.owner}</dd><dt>Scope</dt><dd>{record.scope}</dd><dt>Normal removal</dt><dd>{record.normalRemoval}</dd><dt>Explicit wipe</dt><dd>{record.explicitWipe}</dd></dl>{record.copyable && <Button variant="secondary" disabled={pending} onClick={() => onCopy(record.id)}>Copy path</Button>}</article>)}</div></section>
    <section className="panel settings-section about-section" aria-labelledby="about-heading"><img src="/go-schedule-mark.svg" alt="" /><div><p className="eyebrow">Product information</p><h2 id="about-heading">{workspace.product.name}</h2><p>Version {workspace.product.version}</p><p>Built and maintained by {workspace.product.publisher}.</p><div className="product-links">{workspace.product.links.map((link) => <Button key={link.key} variant="secondary" disabled={pending} onClick={() => onOpen(link.key)}>{link.label}</Button>)}</div></div></section>
  </>
}

function PageHeader() { return <header className="page-header"><div><p className="eyebrow">This computer</p><h1>Settings</h1><p>Desktop preferences, storage ownership, and product information.</p></div></header> }
