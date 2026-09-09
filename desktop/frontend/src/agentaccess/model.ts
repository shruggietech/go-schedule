export interface AgentHTTPStatus { enabled: boolean; endpoint?: string; allowedOrigins: string[]; credentialFingerprint?: string; enabledAt?: string; clientName?: string; lastAccessedAt?: string; requestCount: number }
export interface AgentAuthority { name: string; status: 'available' | 'future'; description: string }
export interface AgentAccessWorkspace { stdioDescription: string; http: AgentHTTPStatus; authorities: AgentAuthority[] }
export interface AgentAccessDraft { clientName: string; port: number; allowedOrigins: string[] }
export interface AgentAccessResult { action: string; outcome: 'accepted' | 'rejected' | 'unavailable'; message: string; workspace?: AgentAccessWorkspace }
export interface AgentAccessBridge { workspace(): Promise<AgentAccessResult>; enable(draft: AgentAccessDraft): Promise<AgentAccessResult>; rotate(): Promise<AgentAccessResult>; revoke(): Promise<AgentAccessResult>; openGuide(): Promise<AgentAccessResult> }
