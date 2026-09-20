export interface AgentHTTPStatus { enabled: boolean; endpoint?: string; allowedOrigins: string[]; credentialFingerprint?: string; enabledAt?: string; clientName?: string; permission?: string; lastAccessedAt?: string; requestCount: number; requireConfirmation?: boolean }
export interface AgentAuthority { name: string; status: 'available' | 'future'; description: string }
export interface AgentAccessWorkspace { stdioDescription: string; http: AgentHTTPStatus; authorities: AgentAuthority[] }
export interface AgentAccessDraft { clientName: string; port: number; allowedOrigins: string[]; permission: 'observe' | 'operate' | 'manage'; requireConfirmation?: boolean }
export interface AgentAccessResult { action: string; outcome: 'accepted' | 'rejected' | 'unavailable'; message: string; workspace?: AgentAccessWorkspace }
export interface AgentAccessBridge { workspace(): Promise<AgentAccessResult>; enable(draft: AgentAccessDraft): Promise<AgentAccessResult>; rotate(): Promise<AgentAccessResult>; revoke(): Promise<AgentAccessResult>; openGuide(): Promise<AgentAccessResult> }
