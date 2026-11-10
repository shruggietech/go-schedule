import { useCallback, useEffect, useRef, useState } from 'react'
import type { AutomationBridge, AutomationWorkspace, OperationResult } from './model'

export function useAutomationWorkspace(bridge: AutomationBridge, daemonID = '', generation = 0) {
  const [workspace, setWorkspace] = useState<AutomationWorkspace>()
  const [status, updateStatus] = useState<OperationResult>()
  const sequence = useRef(0)
  const scope = useRef(daemonID)
  const workspaceScope = useRef(daemonID)
  const uncertain = useRef<OperationResult | undefined>(undefined)
  if (scope.current !== daemonID) { scope.current = daemonID; sequence.current++; uncertain.current = undefined }
  const setStatus = useCallback((result: OperationResult) => { if (result.outcome === 'uncertain') uncertain.current = result; updateStatus(result) }, [])
  const apply = useCallback((result: OperationResult, authoritative = false) => {
    if (result.daemonId && result.daemonId !== scope.current) return
    if (result.workspace && result.outcome === 'accepted') { workspaceScope.current = scope.current; setWorkspace(result.workspace); if (authoritative) uncertain.current = undefined }
    if (uncertain.current && !authoritative) { updateStatus(uncertain.current); return }
    if (uncertain.current && (!result.workspace || result.outcome !== 'accepted')) { updateStatus(uncertain.current); return }
    setStatus(result)
  }, [setStatus])
  const accept = useCallback((result: OperationResult) => { sequence.current++; apply(result) }, [apply])
  const load = useCallback(async () => { const request = ++sequence.current; const target = daemonID; const result = await bridge.workspace(target); if (request === sequence.current && target === scope.current) apply(result, result.outcome === 'accepted' && Boolean(result.workspace)) }, [apply, bridge, daemonID])
  const refreshAfter = useCallback(async (saved: OperationResult) => { const request = ++sequence.current; const target = daemonID; apply(saved); const refresh = await bridge.workspace(target); if (request !== sequence.current || target !== scope.current) return; if (saved.outcome === 'uncertain') { apply(refresh, refresh.outcome === 'accepted' && Boolean(refresh.workspace)); return }; apply(refresh.workspace ? { ...saved, workspace: refresh.workspace } : { ...saved, message: saved.message + ' Refresh the workspace to see the latest state.' }) }, [apply, bridge, daemonID])
  useEffect(() => { void load(); let timer: ReturnType<typeof setTimeout> | undefined
    const unsubscribe = bridge.subscribe?.((event) => { if (['chain.', 'join.', 'join_instance.', 'run.', 'trigger.', 'trigger_set.', 'filesystem_watcher.', 'task.'].some((prefix) => event.kind.startsWith(prefix))) { clearTimeout(timer); timer = setTimeout(() => void load(), 75) } })
    return () => { sequence.current++; clearTimeout(timer); unsubscribe?.() }
  }, [bridge, load, generation])
  return { workspace: workspaceScope.current === daemonID ? workspace : undefined, status, setStatus, load, accept, refreshAfter, reconciling: status?.outcome === 'uncertain' }
}
