import { cloneElement, useEffect, useId, useRef, useState, type ButtonHTMLAttributes, type InputHTMLAttributes, type ReactElement, type ReactNode, type SelectHTMLAttributes, type TextareaHTMLAttributes } from 'react'
import type { ConnectionState } from '../connection/model'

export function Button({ variant = 'primary', pending = false, disabled, ...props }: ButtonHTMLAttributes<HTMLButtonElement> & { variant?: 'primary' | 'secondary' | 'subtle' | 'quiet' | 'affirmative' | 'danger'; pending?: boolean }) {
  const resolvedVariant = variant === 'quiet' ? 'subtle' : variant
  return <button {...props} aria-busy={pending || undefined} disabled={disabled || pending} className={`button button-${resolvedVariant} ${props.className ?? ''}`} />
}

export function Link({ href, children }: { href: string; children: ReactNode }) { return <a className="link" href={href}>{children}</a> }

export function StatusBadge({ state }: { state: ConnectionState }) {
  const labels: Record<ConnectionState, string> = { connecting: 'Connecting', connected: 'Connected', degraded: 'Degraded', recovering: 'Recovering', unavailable: 'Unavailable', access_denied: 'Access denied', unauthorized: 'Unauthorized', revoked: 'Credential revoked', forbidden: 'Forbidden', incompatible: 'Incompatible', trust_changed: 'Trust changed', identity_changed: 'Identity changed', timed_out: 'Timed out' }
  return <span className={`status status-${state}`}><span className="status-shape" aria-hidden="true" />{labels[state]}</span>
}

export function Notice({ title, children, tone = 'info', dismissible = tone === 'warning' || tone === 'error' }: { title: string; children: ReactNode; tone?: 'info' | 'warning' | 'error' | 'success'; dismissible?: boolean }) {
  const [dismissed, setDismissed] = useState(false)
  const contentKey = `${title}:${tone}:${typeof children === 'string' ? children : ''}`
  useEffect(() => setDismissed(false), [contentKey])
  if (tone === 'success') return <ToastRegion message={`${title}: ${typeof children === 'string' ? children : ''}`} />
  if (dismissed) return null
  return <section className={`notice notice-${tone}`} role={tone === 'error' || tone === 'warning' ? 'alert' : 'status'}><div className="notice-heading"><strong>{title}</strong>{dismissible && <Button variant="subtle" aria-label={`Dismiss ${title}`} onClick={() => setDismissed(true)}>Dismiss</Button>}</div><div>{children}</div></section>
}

export function StatePanel({ title, detail, busy = false, action, onAction }: { title: string; detail: string; busy?: boolean; action?: string; onAction?(): void }) {
  return <section className="panel state-panel" aria-busy={busy}><span className="state-symbol" aria-hidden="true">○</span><h2>{title}</h2><p>{detail}</p>{action && <Button onClick={onAction}>{action}</Button>}</section>
}

type FieldControlProps = InputHTMLAttributes<HTMLInputElement> | SelectHTMLAttributes<HTMLSelectElement> | TextareaHTMLAttributes<HTMLTextAreaElement>

export function Field({ label, help, error, children }: { label: string; help?: string; error?: string; children: ReactElement<FieldControlProps> }) {
  const generatedControlID = useId(); const helpID = useId(); const errorID = useId()
  const controlID = children.props.id ?? generatedControlID
  const describedBy = [help && helpID, error && errorID].filter(Boolean).join(' ') || undefined
  const control = cloneElement(children, { id: controlID, 'aria-describedby': describedBy, 'aria-invalid': error ? true : undefined })
  return <div className="field"><label htmlFor={controlID}>{label}</label>{control}{help && <small id={helpID}>{help}</small>}{error && <small className="field-error" id={errorID}>{error}</small>}</div>
}

export function DataTable({ caption, headings, rows }: { caption: string; headings: string[]; rows: ReactNode[][] }) {
  if (rows.length === 0) return <StatePanel title="No results" detail="There is no information to show yet." />
  return <div className="table-scroll" tabIndex={0}><table><caption>{caption}</caption><thead><tr>{headings.map((heading) => <th scope="col" key={heading}>{heading}</th>)}</tr></thead><tbody>{rows.map((row, index) => <tr key={index}>{row.map((cell, cellIndex) => <td key={cellIndex}>{cell}</td>)}</tr>)}</tbody></table></div>
}

export function Disclosure({ summary, children }: { summary: string; children: ReactNode }) { return <details><summary>{summary}</summary><div>{children}</div></details> }

export function Dialog({ open, title, children, actions, invoker, onClose }: { open: boolean; title: string; children: ReactNode; actions?: ReactNode; invoker: HTMLElement | null; onClose(): void }) {
  const ref = useRef<HTMLDivElement>(null)
  useEffect(() => {
    if (!open) return
    const dialog = ref.current
    const focusable = () => Array.from(dialog?.querySelectorAll<HTMLElement>('button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])') ?? []).filter((node) => !node.hasAttribute('disabled'))
    focusable()[0]?.focus()
    const handle = (event: KeyboardEvent) => {
      if (event.key === 'Escape') { event.preventDefault(); onClose(); invoker?.focus(); return }
      if (event.key !== 'Tab') return
      const nodes = focusable(); if (nodes.length === 0) return
      const first = nodes[0]; const last = nodes[nodes.length - 1]
      if (event.shiftKey && document.activeElement === first) { event.preventDefault(); last.focus() }
      else if (!event.shiftKey && document.activeElement === last) { event.preventDefault(); first.focus() }
    }
    document.addEventListener('keydown', handle)
    return () => document.removeEventListener('keydown', handle)
  }, [invoker, onClose, open])
  if (!open) return null
  return <div className="dialog-backdrop"><div aria-modal="true" className="dialog" ref={ref} role="dialog" aria-labelledby="dialog-title"><h2 id="dialog-title">{title}</h2><div className="dialog-content">{children}</div><div className="dialog-actions">{actions}<Button variant="secondary" onClick={() => { onClose(); invoker?.focus() }}>Close</Button></div></div></div>
}

export function ToastRegion({ message, duration = 5_000 }: { message: string; duration?: number }) {
  const [visible, setVisible] = useState(message)
  const [paused, setPaused] = useState(false)
  useEffect(() => { setVisible(message) }, [message])
  useEffect(() => {
    if (!visible || paused) return
    const timer = window.setTimeout(() => setVisible(''), duration)
    return () => window.clearTimeout(timer)
  }, [duration, paused, visible])
  if (!visible) return <div className="toast-live" role="status" aria-live="polite" aria-atomic="true" />
  return <div className="toast-region" role="status" aria-live="polite" aria-atomic="true" onMouseEnter={() => setPaused(true)} onMouseLeave={() => setPaused(false)} onFocusCapture={() => setPaused(true)} onBlurCapture={(event) => { if (!event.currentTarget.contains(event.relatedTarget)) setPaused(false) }}><span>{visible}</span><Button variant="subtle" aria-label="Dismiss notification" onClick={() => setVisible('')}>Dismiss</Button></div>
}
