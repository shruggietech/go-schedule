import type { AgentAccessBridge, AgentAccessDraft, AgentAccessResult, AgentActionsResult, AgentGrantDraft, AgentGrantEditDraft } from './model'

type NativeWindow = Window & { go?: { main?: { App?: Record<string, (...args: unknown[]) => Promise<AgentAccessResult | AgentActionsResult>> } } }
const unavailable = (action: string): AgentAccessResult => ({ action, outcome: 'unavailable', message: 'Agent Access is available in the installed desktop application.' })
const actionsUnavailable = (): AgentActionsResult => ({ action: 'load_agent_actions', outcome: 'unavailable', message: 'Agent activity is available in the installed desktop application.', actions: [] })

export function createAgentAccessBridge(nativeWindow: NativeWindow = window): AgentAccessBridge {
  const app = nativeWindow.go?.main?.App
  const call = (name: string, action: string, ...args: unknown[]) => (app?.[name]?.(...args) as Promise<AgentAccessResult> | undefined) ?? Promise.resolve(unavailable(action))
  return {
    workspace: () => call('AgentAccessWorkspace', 'load_agent_access'),
    enable: (draft: AgentAccessDraft) => call('EnableAgentAccess', 'enable_agent_access', draft),
    rotate: () => call('RotateAgentAccess', 'rotate_agent_access'),
    revoke: () => call('RevokeAgentAccess', 'revoke_agent_access'),
    createGrant: (draft: AgentGrantDraft) => call('CreateAgentGrant', 'create_agent_grant', draft),
    editGrant: (draft: AgentGrantEditDraft) => call('EditAgentGrant', 'edit_agent_grant', draft),
    revokeGrant: (actorId: string) => call('RevokeAgentGrant', 'revoke_agent_grant', actorId),
    actions: (actorId: string) => (app?.AgentGrantActions?.(actorId) as Promise<AgentActionsResult> | undefined) ?? Promise.resolve(actionsUnavailable()),
    openGuide: () => call('OpenAgentAccessGuide', 'open_agent_access_guide'),
  }
}

export const agentAccessBridge = createAgentAccessBridge()
