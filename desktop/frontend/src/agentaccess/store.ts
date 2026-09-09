import { useCallback, useEffect, useRef, useState } from 'react'
import type { AgentAccessBridge, AgentAccessDraft, AgentAccessResult, AgentAccessWorkspace } from './model'

const unavailable = (action: string): AgentAccessResult => ({ action, outcome: 'unavailable', message: 'Agent Access is unavailable. Check the local scheduler connection, then try again.' })

export function useAgentAccess(bridge: AgentAccessBridge, available: boolean, refreshToken: number) {
  const [workspace, setWorkspace] = useState<AgentAccessWorkspace>()
  const [status, setStatus] = useState<AgentAccessResult>()
  const [pending, setPending] = useState(false)
  const sequence = useRef(0)
  const mutationPending = useRef(false)
  const apply = useCallback((result: AgentAccessResult) => { setStatus(result); if (result.workspace) setWorkspace(result.workspace) }, [])
  const load = useCallback(async () => {
    if (!available || mutationPending.current) return
    const request = ++sequence.current
    try {
      const result = await bridge.workspace()
      if (request === sequence.current) apply(result)
    } catch {
      if (request === sequence.current) apply(unavailable('load_agent_access'))
    }
  }, [apply, available, bridge])
  const mutate = useCallback(async (work: () => Promise<AgentAccessResult>) => {
    if (!available || mutationPending.current) return
    mutationPending.current = true
    setPending(true)
    const request = ++sequence.current
    try { const result = await work(); if (request === sequence.current) apply(result) } catch { if (request === sequence.current) apply(unavailable('agent_access_action')) } finally { mutationPending.current = false; setPending(false) }
  }, [apply, available])
  useEffect(() => { void load() }, [load, refreshToken])
  useEffect(() => {
    if (!available || !workspace?.http.enabled) return
    const timer = setInterval(() => void load(), 5000)
    return () => clearInterval(timer)
  }, [available, load, workspace?.http.enabled])
  return { workspace, message: status?.message ?? '', pending, load, enable: (draft: AgentAccessDraft) => mutate(() => bridge.enable(draft)), rotate: () => mutate(bridge.rotate), revoke: () => mutate(bridge.revoke), openGuide: () => mutate(bridge.openGuide) }
}
