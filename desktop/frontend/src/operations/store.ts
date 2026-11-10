import { useCallback, useEffect, useRef, useState } from 'react'
import type { ActivityWorkspace, OperationResult, OperationsBridge, ScheduleSnapshot } from './model'

const relevant = (kind: string) => ['run.', 'log.', 'alert.', 'task.'].some((prefix) => kind.startsWith(prefix))

export function useSchedule(bridge: OperationsBridge, days: number, available: boolean, refreshToken: number) {
  const [snapshot, setSnapshot] = useState<ScheduleSnapshot>()
  const [status, setStatus] = useState<OperationResult>()
  const sequence = useRef(0)
  const load = useCallback(async () => {
    if (!available) return
    const request = ++sequence.current
    const result = await bridge.scheduleWindow(days)
    if (request !== sequence.current) return
    setStatus(result)
    if (result.schedule) setSnapshot(result.schedule)
  }, [available, bridge, days])
  useEffect(() => { void load() }, [load, refreshToken])
  useEffect(() => {
    let timer: ReturnType<typeof setTimeout> | undefined
    const unsubscribe = bridge.subscribe?.((event) => { if (relevant(event.kind)) { clearTimeout(timer); timer = setTimeout(() => void load(), 75) } })
    return () => { clearTimeout(timer); unsubscribe?.() }
  }, [bridge, load])
  return { snapshot, status, load }
}

export function useActivity(bridge: OperationsBridge, available: boolean, refreshToken: number, initialRecord?: { kind: string; id: string }, daemonID = '') {
  const [workspace, setWorkspace] = useState<ActivityWorkspace>()
  const [workspaceScope, setWorkspaceScope] = useState('')
  const [status, setStatus] = useState<OperationResult>()
  const [pending, setPending] = useState(false)
  const [cancelingRunID, setCancelingRunID] = useState<string>()
  const [cancellationNotice, setCancellationNotice] = useState('')
  const sequence = useRef(0)
  const mounted = useRef(true)
  const mutation = useRef<object | undefined>(undefined)
  const refreshAfterMutation = useRef(false)
  const uncertain = useRef(false)
  const cancellationRun = useRef<string | undefined>(undefined)
  const polling = useRef<{ token: object; timer?: ReturnType<typeof setTimeout>; expiry?: ReturnType<typeof setTimeout> } | undefined>(undefined)
  const scope = daemonID + ':' + refreshToken + ':' + available
  const currentScope = useRef(scope)
  currentScope.current = scope
  const valid = useCallback((expected: string) => mounted.current && currentScope.current === expected, [])
  const stopPolling = useCallback(() => {
    const previous = polling.current
    if (previous) { clearTimeout(previous.timer); clearTimeout(previous.expiry) }
    polling.current = undefined
  }, [])
  useEffect(() => {
    mounted.current = true
    sequence.current++
    mutation.current = undefined
    refreshAfterMutation.current = false
    uncertain.current = false
    cancellationRun.current = undefined
    setPending(false); setCancelingRunID(undefined); setCancellationNotice('')
    return () => { mounted.current = false; sequence.current++; stopPolling() }
  }, [scope, stopPolling])
  const apply = useCallback((result: OperationResult, expected: string) => {
    if (!valid(expected)) return
    const authoritative = result.outcome === 'accepted' && Boolean(result.activity)
    if (result.outcome === 'uncertain') uncertain.current = true
    else if (authoritative) uncertain.current = false
    setStatus(uncertain.current && result.outcome !== 'uncertain' ? { ...result, outcome: 'uncertain', message: result.message + ' The previous mutation remains uncertain. Refresh Activity for an authoritative observation before trying again.' } : result)
    if (result.activity) {
      setWorkspace(result.activity); setWorkspaceScope(expected)
      const run = result.activity.runs.find((value) => value.id === cancellationRun.current)
      if (run && (run.endedAt || (run.outcome && run.state !== 'running'))) {
        stopPolling(); setCancelingRunID(undefined)
        setCancellationNotice('Run ' + run.id + ' finished with outcome ' + run.state.replaceAll('_', ' ') + '.')
      }
    }
  }, [stopPolling, valid])
  const read = useCallback(() => initialRecord && bridge.activityRecord ? bridge.activityRecord(initialRecord.kind, initialRecord.id) : bridge.activityWorkspace(), [bridge, initialRecord?.id, initialRecord?.kind])
  const load = useCallback(async () => {
    if (!available) return
    if (mutation.current) { refreshAfterMutation.current = true; return }
    const request = ++sequence.current, expected = currentScope.current
    let result = await read()
    const selected = cancellationRun.current
    if (request === sequence.current && valid(expected) && result.activity && selected && !result.activity.runs.some((run) => run.id === selected) && bridge.activityRecord) {
      const exact = await bridge.activityRecord('run', selected)
      const run = exact.activity?.runs.find((value) => value.id === selected)
      if (run && exact.outcome === 'accepted') result = { ...result, activity: { ...result.activity, runs: [...result.activity.runs, run] } }
    }
    if (request === sequence.current) apply(result, expected)
  }, [apply, available, bridge, read, valid])
  const acknowledge = useCallback(async (ids: string[]) => {
    if (!available || mutation.current || uncertain.current) return
    const token = {}, expected = currentScope.current, request = ++sequence.current
    mutation.current = token; setPending(true)
    try {
      const result = ids.length === 1 ? await bridge.acknowledgeAlert(ids[0]) : await bridge.acknowledgeAlerts(ids)
      if (request === sequence.current) apply(result, expected)
    } finally {
      if (mutation.current === token) { mutation.current = undefined; if (valid(expected)) setPending(false) }
      if (valid(expected) && refreshAfterMutation.current) { refreshAfterMutation.current = false; void load() }
    }
  }, [apply, available, bridge, load, valid])
  const reconcile = useCallback((id: string, expected: string) => {
    stopPolling()
    const token = {}, session: NonNullable<typeof polling.current> = { token }
    polling.current = session
    setCancelingRunID(id)
    session.expiry = setTimeout(() => {
      if (polling.current?.token !== token || !valid(expected)) return
      stopPolling(); setCancelingRunID(undefined)
      setCancellationNotice('Termination is not yet confirmed for run ' + id + '. Refresh Activity to inspect its current result.')
    }, 8000)
    const poll = async () => {
      if (polling.current?.token !== token || !valid(expected)) return
      try {
        const result = bridge.activityRecord ? await bridge.activityRecord('run', id) : await bridge.activityWorkspace()
        if (polling.current?.token !== token || !valid(expected)) return
        const run = result.activity?.runs.find((value) => value.id === id)
        if (result.outcome === 'accepted' && run) {
          setWorkspace((previous) => previous ? { ...previous, runs: previous.runs.some((value) => value.id === id) ? previous.runs.map((value) => value.id === id ? run : value) : [...previous.runs, run] } : result.activity)
          setWorkspaceScope(expected)
          if (run.endedAt || (run.outcome && run.state !== 'running')) {
            stopPolling(); setCancelingRunID(undefined)
            setCancellationNotice('Run ' + id + ' finished with outcome ' + run.state.replaceAll('_', ' ') + '.')
            return
          }
          if (run.cancellationStatus === 'failed') {
            stopPolling(); setCancelingRunID(undefined)
            setCancellationNotice('Cancellation cleanup is unconfirmed for run ' + id + '. ' + (run.cancellationError || 'Refresh Activity for current evidence.'))
            return
          }
        }
      } catch { /* Retain the last complete observation; expiry reports unconfirmed termination. */ }
      if (polling.current?.token === token && valid(expected)) session.timer = setTimeout(() => void poll(), 250)
    }
    void poll()
  }, [bridge, stopPolling, valid])
  const cancelRun = useCallback(async (id: string) => {
    if (!available || !bridge.cancelRun || mutation.current || polling.current || uncertain.current || workspaceScope !== currentScope.current) return
    const token = {}, expected = currentScope.current, request = ++sequence.current
    mutation.current = token; setPending(true)
    try {
      const result = await bridge.cancelRun(id, daemonID)
      if (request !== sequence.current || !valid(expected)) return
      apply(result, expected)
      if (result.outcome === 'accepted' && result.cancellation) {
        cancellationRun.current = id
        setCancellationNotice(result.message)
        if (result.cancellation.status === 'requested') reconcile(id, expected)
        else { await load() }
      }
    } catch {
      if (valid(expected)) apply({ action: 'cancel_run', outcome: 'uncertain', message: 'The cancellation response was not authoritative. Refresh the selected scheduler before trying again.' }, expected)
    } finally {
      if (mutation.current === token) { mutation.current = undefined; if (valid(expected)) setPending(false) }
      if (valid(expected) && refreshAfterMutation.current) { refreshAfterMutation.current = false; void load() }
    }
  }, [apply, available, bridge, daemonID, load, reconcile, valid, workspaceScope])
  useEffect(() => { void load() }, [load, refreshToken])
  useEffect(() => {
    if (initialRecord && bridge.activityRecord) return
    let timer: ReturnType<typeof setTimeout> | undefined
    const unsubscribe = bridge.subscribe?.((event) => { if (relevant(event.kind)) { clearTimeout(timer); timer = setTimeout(() => void load(), 75) } })
    return () => { clearTimeout(timer); unsubscribe?.() }
  }, [bridge, load, initialRecord?.id, initialRecord?.kind])
  return { workspace, status, pending, load, acknowledge, cancelRun, cancelingRunID, cancellationNotice, current: workspaceScope === scope }
}
