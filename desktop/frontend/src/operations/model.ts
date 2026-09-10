export type ScheduleOccurrence = { id: string; taskId: string; taskName: string; runId?: string; time: string; kind: 'prediction' | 'recorded'; state: string; outcome?: string }
export type ScheduleSnapshot = { from: string; to: string; loadedAt: string; occurrences: ScheduleOccurrence[] }
export type RunRecord = { id: string; taskId: string; scheduledFor: string; startedAt?: string; endedAt?: string; state: string; outcome?: string; exitCode?: number; output: string; outputTruncated: boolean; trigger: string; sourceTaskId?: string; sourceRunId?: string; sourceTriggerId?: string; sourceWatcherId?: string }
export type LogRecord = { id: string; time: string; severity: string; source: string; message: string; taskId?: string; runId?: string; detail?: string }
export type AlertRecord = { id: string; taskId?: string; runId?: string; time: string; severity: string; kind: string; message: string; acknowledged: boolean }
export type ActivityWorkspace = { runs: RunRecord[]; logs: LogRecord[]; alerts: AlertRecord[]; logPath: string; loadedAt: string }
export type OperationResult = { action: string; outcome: 'accepted' | 'rejected' | 'unavailable' | 'uncertain'; message: string; field?: string; schedule?: ScheduleSnapshot; activity?: ActivityWorkspace }
export interface OperationsBridge {
  scheduleWindow(days: number): Promise<OperationResult>
  activityWorkspace(): Promise<OperationResult>
  acknowledgeAlert(id: string): Promise<OperationResult>
  acknowledgeAlerts(ids: string[]): Promise<OperationResult>
  subscribe?(listener: (event: { kind: string }) => void): () => void
}
