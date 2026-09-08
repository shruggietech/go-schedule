import { useCallback, useEffect, useRef, useState } from 'react'
import type { ChannelDraft, NotificationBridge, NotificationResult, NotificationWorkspace, Policy, PolicyDraft } from './model'

const relevant = (kind: string) => ['notification.', 'task.', 'group.', 'run.'].some((prefix) => kind.startsWith(prefix))

export function useNotifications(bridge: NotificationBridge, available: boolean, refreshToken: number) {
  const [workspace, setWorkspace] = useState<NotificationWorkspace>()
  const [policy, setPolicy] = useState<Policy>()
  const [status, setStatus] = useState<NotificationResult>()
  const [pending, setPending] = useState(false)
  const [policyPending, setPolicyPending] = useState(false)
  const workspaceSequence = useRef(0)
  const policySequence = useRef(0)
  const selectedScope = useRef<{ type: 'task' | 'group'; id: string } | undefined>(undefined)
  const mutationPending = useRef(false)
  const policyMutationPending = useRef(false)
  const apply = useCallback((result: NotificationResult) => {
    setStatus(result)
    if (result.workspace) setWorkspace(result.workspace)
    if (result.policy) setPolicy(result.policy)
  }, [])
  const load = useCallback(async () => {
    if (!available) return
    const request = ++workspaceSequence.current
    const result = await bridge.workspace()
    if (request === workspaceSequence.current) apply(result)
  }, [apply, available, bridge])
  const mutate = useCallback(async (work: () => Promise<NotificationResult>) => {
    if (!available || mutationPending.current) return undefined
    mutationPending.current = true
    setPending(true)
    const request = ++workspaceSequence.current
    try {
      const result = await work()
      if (request === workspaceSequence.current) apply(result)
      return result
    } finally { mutationPending.current = false; setPending(false) }
  }, [apply, available])
  const selectPolicy = useCallback(async (type: 'task' | 'group', id: string) => {
    selectedScope.current = { type, id }
    setPolicyPending(true)
    const request = ++policySequence.current
    try {
      const result = await bridge.policy(type, id)
      if (request === policySequence.current) apply(result)
    } finally { if (request === policySequence.current) setPolicyPending(false) }
  }, [apply, bridge])
  const savePolicy = useCallback(async (draft: PolicyDraft) => {
    if (!available || policyMutationPending.current) return
    policyMutationPending.current = true
    setPolicyPending(true)
    const request = ++policySequence.current
    try {
      const result = await bridge.savePolicy(draft)
      if (request === policySequence.current) apply(result)
    } finally { policyMutationPending.current = false; if (request === policySequence.current) setPolicyPending(false) }
  }, [apply, available, bridge])
  useEffect(() => { void load() }, [load, refreshToken])
  useEffect(() => {
    let timer: ReturnType<typeof setTimeout> | undefined
    const unsubscribe = bridge.subscribe?.((event) => {
      if (!relevant(event.kind)) return
      clearTimeout(timer)
      timer = setTimeout(() => { void load(); const scope = selectedScope.current; if (scope) void selectPolicy(scope.type, scope.id) }, 75)
    })
    return () => { clearTimeout(timer); unsubscribe?.() }
  }, [bridge, load, selectPolicy])
  return {
    workspace, policy, status, pending, policyPending, load, selectPolicy, savePolicy,
    saveChannel: (draft: ChannelDraft) => mutate(() => bridge.saveChannel(draft)),
    setChannelEnabled: (id: string, enabled: boolean) => mutate(() => bridge.setChannelEnabled(id, enabled)),
    testChannel: (id: string) => mutate(() => bridge.testChannel(id)),
    deleteChannel: (id: string) => mutate(() => bridge.deleteChannel(id)),
  }
}
