import { useState } from 'react'
import { Button, Notice } from '../components'
import type { BundleBridge, BundlePlan, BundleResult } from './model'

export function BundlesPage({ bridge, targetName, manageAvailable }: { bridge: BundleBridge; targetName: string; manageAvailable: boolean }) {
  const [document, setDocument] = useState('')
  const [plan, setPlan] = useState<BundlePlan>()
  const [reviewed, setReviewed] = useState(false)
  const [result, setResult] = useState<BundleResult>()
  const run = async (action: () => Promise<BundleResult>) => { const next = await action(); setResult(next); if (next.document) setDocument(next.document); if (next.plan) { setPlan(next.plan); setReviewed(false) } }
  return <>
    <header className="page-header"><div><p className="eyebrow">{targetName}</p><h1>Portable bundles</h1><p>Export portable intent, inspect a target-bound plan, then explicitly apply that plan.</p></div><Button onClick={() => void run(() => bridge.exportBundle())}>Export bundle</Button></header>
    {result && <Notice title={result.outcome === 'accepted' ? 'Bundle action completed' : 'Bundle action needs attention'} tone={result.outcome === 'accepted' ? 'success' : 'error'}>{result.message}{result.issues?.length ? <ul>{result.issues.map((issue, index) => <li key={`${issue.kind}-${index}`}>{issue.kind}: {issue.message}</li>)}</ul> : null}</Notice>}
    <section className="card"><label htmlFor="bundle-document">Bundle JSON</label><textarea id="bundle-document" value={document || result?.document || ''} onChange={(event) => { setDocument(event.target.value); setPlan(undefined); setReviewed(false) }} rows={14} spellCheck={false} /><div className="actions"><Button variant="secondary" onClick={() => void run(() => bridge.validate(document || result?.document || ''))}>Validate</Button><Button variant="secondary" onClick={() => void run(() => bridge.compare(document || result?.document || ''))}>Compare drift</Button><Button onClick={() => void run(() => bridge.preview(document || result?.document || ''))}>Preview changes</Button></div></section>
    {plan && <section className="card"><h2>{result?.compare_only ? 'Drift comparison' : 'Reviewed plan'}</h2><p>Target daemon: <code>{plan.target_daemon_id}</code></p><ul>{plan.items.map((item) => <li key={`${item.kind}-${item.portable_id}`}><strong>{item.action}</strong> {item.kind} {item.name || item.portable_id}{item.message ? `: ${item.message}` : ''}</li>)}</ul>{!result?.compare_only && <><label><input type="checkbox" checked={reviewed} onChange={(event) => setReviewed(event.target.checked)} /> I reviewed this target-bound plan.</label><div className="actions"><Button disabled={!reviewed || !manageAvailable} onClick={() => void run(() => bridge.apply(plan))}>Apply reviewed plan</Button></div></>}</section>}
  </>
}
