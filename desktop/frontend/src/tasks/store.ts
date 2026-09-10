import { useCallback, useEffect, useRef, useState } from 'react'
import type { OperationResult, TaskBridge, Workspace } from './model'

export function useTaskWorkspace(bridge: TaskBridge, available = true, refreshToken = 0) {
  const [workspace, setWorkspace] = useState<Workspace>()
  const [status, setStatus] = useState<OperationResult>()
  const [selected, setSelected] = useState('')
  const requestSequence = useRef(0)
  const apply = useCallback((result: OperationResult) => {
    setStatus(result)
    if (!result.workspace) return
    setWorkspace(result.workspace)
    setSelected((current) => result.workspace?.tasks.some((task) => task.id === current) ? current : result.workspace?.tasks[0]?.id ?? '')
  }, [])
  const accept = useCallback((result: OperationResult) => { requestSequence.current++; apply(result) }, [apply])
  const load = useCallback(async () => {
	if (!available) return
    const request = ++requestSequence.current
    const result = await bridge.workspace()
    if (request === requestSequence.current) apply(result)
  }, [apply, available, bridge])
  useEffect(() => {
    if (available) return
    requestSequence.current++
    setStatus({ action: 'load', outcome: 'unavailable', message: 'Tasks are unavailable until the selected scheduler reconnects.' })
  }, [available])
  useEffect(() => {
    void load()
    let timer: ReturnType<typeof setTimeout> | undefined
    const unsubscribe = bridge.subscribe?.((event) => {
      if (event.kind.startsWith('task.') || event.kind.startsWith('group.')) {
        clearTimeout(timer)
        timer = setTimeout(() => void load(), 75)
      }
    })
    return () => { clearTimeout(timer); unsubscribe?.() }
  }, [bridge, load, refreshToken])
  return { workspace, status, selected, setSelected, load, accept, setStatus }
}
