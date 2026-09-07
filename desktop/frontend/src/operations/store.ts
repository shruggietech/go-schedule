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

export function useActivity(bridge: OperationsBridge, available: boolean, refreshToken: number) {
  const [workspace, setWorkspace] = useState<ActivityWorkspace>()
  const [status, setStatus] = useState<OperationResult>()
  const [pending, setPending] = useState(false)
  const sequence = useRef(0)
  const apply = useCallback((result: OperationResult) => { setStatus(result); if (result.activity) setWorkspace(result.activity) }, [])
  const load = useCallback(async () => {
    if (!available) return
    const request = ++sequence.current
    const result = await bridge.activityWorkspace()
    if (request === sequence.current) apply(result)
  }, [apply, available, bridge])
  const acknowledge = useCallback(async (ids: string[]) => {
    if (!available || pending) return
    setPending(true)
    const request = ++sequence.current
    try {
      const result = ids.length === 1 ? await bridge.acknowledgeAlert(ids[0]) : await bridge.acknowledgeAlerts(ids)
      if (request === sequence.current) apply(result)
    } finally { setPending(false) }
  }, [apply, available, bridge, pending])
  useEffect(() => { void load() }, [load, refreshToken])
  useEffect(() => {
    let timer: ReturnType<typeof setTimeout> | undefined
    const unsubscribe = bridge.subscribe?.((event) => { if (relevant(event.kind)) { clearTimeout(timer); timer = setTimeout(() => void load(), 75) } })
    return () => { clearTimeout(timer); unsubscribe?.() }
  }, [bridge, load])
  return { workspace, status, pending, load, acknowledge }
}
