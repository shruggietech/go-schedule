import { useCallback, useEffect, useRef, useState } from 'react'
import type { AutomationBridge, AutomationWorkspace, OperationResult } from './model'

export function useAutomationWorkspace(bridge: AutomationBridge) {
  const [workspace, setWorkspace] = useState<AutomationWorkspace>()
  const [status, setStatus] = useState<OperationResult>()
  const sequence = useRef(0)
  const apply = useCallback((result: OperationResult) => { setStatus(result); if (result.workspace) setWorkspace(result.workspace) }, [])
  const accept = useCallback((result: OperationResult) => { sequence.current++; apply(result) }, [apply])
  const load = useCallback(async () => { const request = ++sequence.current; const result = await bridge.workspace(); if (request === sequence.current) apply(result) }, [apply, bridge])
  useEffect(() => {
    void load(); let timer: ReturnType<typeof setTimeout> | undefined
    const unsubscribe = bridge.subscribe?.((event) => { if (['chain.', 'trigger.', 'trigger_set.', 'filesystem_watcher.', 'task.'].some((prefix) => event.kind.startsWith(prefix))) { clearTimeout(timer); timer = setTimeout(() => void load(), 75) } })
    return () => { clearTimeout(timer); unsubscribe?.() }
  }, [bridge, load])
  return { workspace, status, setStatus, load, accept }
}
