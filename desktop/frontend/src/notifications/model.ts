export type Channel = { id: string; name: string; kind: string; endpointSummary: string; hasAuthorization: boolean; enabled: boolean; updatedAt: string }
export type ChannelDraft = { id: string; name: string; endpoint: string; authorization: string; replaceEndpoint: boolean; replaceAuthorization: boolean; isNew: boolean }
export type Scope = { type: 'task' | 'group'; id: string; name: string; context: string }
export type Assignment = { channelId: string; onSuccess: boolean; onFailure: boolean }
export type PolicyDraft = { scopeType: 'task' | 'group'; scopeId: string; assignments: Assignment[] }
export type Policy = { scope: Scope; directAssignments: Assignment[]; effectiveSourceType: 'task' | 'group' | 'none'; effectiveSourceId?: string; effectiveSourceName: string; effectiveAssignments: Assignment[] }
export type ConfiguredScope = { type: 'task' | 'group'; id: string; name: string; context: string; sourceType: 'task' | 'group' | 'none'; sourceName: string; onSuccess: boolean; onFailure: boolean; destinationCount: number; enabledDestinationCount: number; enabledSuccessDestinationCount: number; enabledFailureDestinationCount: number }
export type Delivery = { id: string; channelId?: string; channelName: string; destinationSummary: string; kind: 'test' | 'task_outcome'; taskId?: string; taskName?: string; groupId?: string; groupName?: string; runId?: string; state: 'queued' | 'retrying' | 'sending' | 'successful' | 'failed'; attempts: number; nextAttemptAt?: string; createdAt: string; claimedAt?: string; completedAt?: string; lastStatus?: number; lastError?: string }
export type NotificationWorkspace = { channels: Channel[]; tasks: Scope[]; groups: Scope[]; coverage?: ConfiguredScope[]; coverageComplete?: boolean; deliveries: Delivery[]; loadedAt: string }
export type NotificationResult = { action: string; outcome: 'accepted' | 'rejected' | 'conflict' | 'stale' | 'unavailable'; message: string; field?: string; entityId?: string; workspace?: NotificationWorkspace; policy?: Policy }
export interface NotificationBridge {
  workspace(): Promise<NotificationResult>
  saveChannel(draft: ChannelDraft): Promise<NotificationResult>
  setChannelEnabled(id: string, enabled: boolean): Promise<NotificationResult>
  testChannel(id: string): Promise<NotificationResult>
  deleteChannel(id: string): Promise<NotificationResult>
  policy(scopeType: 'task' | 'group', scopeId: string): Promise<NotificationResult>
  savePolicy(draft: PolicyDraft): Promise<NotificationResult>
  subscribe?(listener: (event: { kind: string }) => void): () => void
}
