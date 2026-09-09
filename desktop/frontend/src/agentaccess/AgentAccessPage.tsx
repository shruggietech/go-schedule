import { useState } from 'react'
import { Button, Notice, StatePanel } from '../components'
import type { AgentAccessBridge } from './model'
import { useAgentAccess } from './store'

export function AgentAccessPage({ bridge, available, refreshToken }: { bridge: AgentAccessBridge; available: boolean; refreshToken: number }) {
  const access = useAgentAccess(bridge, available, refreshToken)
  const [name, setName] = useState('Codex')
  const [port, setPort] = useState('43123')
  const [origins, setOrigins] = useState('')
  if (!available) return <StatePanel title="Agent Access unavailable" detail="Reconnect to the local scheduler service to inspect or control agent access." />
  if (!access.workspace && access.message) return <StatePanel title="Agent Access unavailable" detail={access.message} action="Try again" onAction={() => void access.load()} />
  if (!access.workspace) return <StatePanel title="Loading Agent Access" detail="Reading the local MCP access boundary." busy />
  const http = access.workspace.http
  return <section aria-labelledby="agent-access-title">
    <header className="page-header"><div><p className="eyebrow">This computer</p><h1 id="agent-access-title">Agent Access</h1><p>Understand exactly what local agents can read and close optional network access.</p></div><Button variant="secondary" onClick={() => void access.openGuide()}>Setup guide</Button></header>
    {access.message && <Notice title="Agent Access status" tone={access.message.includes('could not') || access.message.includes('unavailable') ? 'warning' : 'info'}>{access.message}</Notice>}
    <div className="settings-grid">
      <article className="panel"><div className="section-heading"><div><p className="eyebrow">Process access</p><h2>stdio</h2></div><span>Available on demand</span></div><p>{access.workspace.stdioDescription}</p></article>
      <article className="panel"><div className="section-heading"><div><p className="eyebrow">Optional network access</p><h2>Localhost HTTP</h2></div><span>{http.enabled ? 'Active' : 'Off'}</span></div>
        {http.enabled ? <><dl><dt>Client</dt><dd>{http.clientName}</dd><dt>Endpoint</dt><dd><code>{http.endpoint}</code></dd><dt>Allowed browser origins</dt><dd>{http.allowedOrigins.length ? http.allowedOrigins.join(', ') : 'None'}</dd><dt>Credential fingerprint</dt><dd><code>{http.credentialFingerprint}</code></dd><dt>Enabled</dt><dd>{http.enabledAt}</dd><dt>Last successful access</dt><dd>{http.lastAccessedAt ?? 'No successful requests yet'}</dd><dt>Successful requests</dt><dd>{http.requestCount}</dd></dl><div className="button-row"><Button disabled={access.pending} onClick={() => void access.rotate()}>Rotate credential</Button><Button variant="secondary" disabled={access.pending} onClick={() => void access.revoke()}>Revoke access</Button></div></> : <form onSubmit={(event) => { event.preventDefault(); void access.enable({ clientName: name, port: Number(port), allowedOrigins: origins.split(/\r?\n/).map((value) => value.trim()).filter(Boolean) }) }}><label>Client name<input value={name} maxLength={64} required onChange={(event) => setName(event.target.value)} /></label><label>Loopback port<input value={port} type="number" min="1" max="65535" required onChange={(event) => setPort(event.target.value)} /></label><label>Allowed browser origins (one per line)<textarea value={origins} placeholder="http://127.0.0.1:3000" onChange={(event) => setOrigins(event.target.value)} /></label><Button disabled={access.pending} type="submit">Enable and copy credential</Button></form>}
      </article>
      <article className="panel"><p className="eyebrow">Authority</p><h2>What agents can do</h2><div className="authority-list">{access.workspace.authorities.map((authority) => <section key={authority.name}><h3>{authority.name} <span className="muted">({authority.status === 'available' ? 'Available' : 'Future, unavailable'})</span></h3><p>{authority.description}</p></section>)}</div></article>
    </div>
  </section>
}
