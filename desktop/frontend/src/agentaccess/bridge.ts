import type { AgentAccessBridge, AgentAccessDraft, AgentAccessResult } from './model'

type NativeWindow = Window & { go?: { main?: { App?: Record<string, (...args: unknown[]) => Promise<AgentAccessResult>> } } }
const unavailable = (action: string): AgentAccessResult => ({ action, outcome: 'unavailable', message: 'Agent Access is available in the installed desktop application.' })

export function createAgentAccessBridge(nativeWindow: NativeWindow = window): AgentAccessBridge {
  const app = nativeWindow.go?.main?.App
  const call = (name: string, action: string, ...args: unknown[]) => app?.[name]?.(...args) ?? Promise.resolve(unavailable(action))
  return {
    workspace: () => call('AgentAccessWorkspace', 'load_agent_access'),
    enable: (draft: AgentAccessDraft) => call('EnableAgentAccess', 'enable_agent_access', draft),
    rotate: () => call('RotateAgentAccess', 'rotate_agent_access'),
    revoke: () => call('RevokeAgentAccess', 'revoke_agent_access'),
    openGuide: () => call('OpenAgentAccessGuide', 'open_agent_access_guide'),
  }
}

export const agentAccessBridge = createAgentAccessBridge()
