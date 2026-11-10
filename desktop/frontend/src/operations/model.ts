export type ScheduleOccurrence = { id: string; taskId: string; taskName: string; runId?: string; time: string; kind: 'prediction' | 'recorded'; state: string; outcome?: string }
export type ScheduleSnapshot = { from: string; to: string; loadedAt: string; occurrences: ScheduleOccurrence[] }
export type RunRecord = { deliveryId?: string; joinInstanceId?: string; joinBranchId?: string; joinRole?: string; cancellationStatus?: string; cancellationError?: string; id: string; taskId: string; scheduledFor: string; startedAt?: string; endedAt?: string; state: string; outcome?: string; exitCode?: number; output: string; outputTruncated: boolean; trigger: string; sourceTaskId?: string; sourceRunId?: string; sourceTriggerId?: string; sourceWatcherId?: string }
export type LogRecord = { id: string; time: string; severity: string; source: string; message: string; taskId?: string; runId?: string; detail?: string }
export type AlertRecord = { id: string; taskId?: string; runId?: string; time: string; severity: string; kind: string; message: string; acknowledged: boolean }
export type ActivityWorkspace = { runs: RunRecord[]; logs: LogRecord[]; alerts: AlertRecord[]; logPath: string; loadedAt: string }
export type RunCancellation = { run_id: string; task_id: string; status: 'requested' | 'already_finished'; outcome?: string; cancellation_status?: string; cancellation_error?: string }
export type OperationResult = { cancellation?: RunCancellation; action: string; outcome: 'accepted' | 'rejected' | 'unavailable' | 'uncertain'; message: string; field?: string; schedule?: ScheduleSnapshot; activity?: ActivityWorkspace }
export interface OperationsBridge {
  scheduleWindow(days: number): Promise<OperationResult>
  activityWorkspace(): Promise<OperationResult>
	activityRecord?(kind: string, id: string): Promise<OperationResult>
  cancelRun?(id: string, daemonID: string): Promise<OperationResult>
  acknowledgeAlert(id: string): Promise<OperationResult>
  acknowledgeAlerts(ids: string[]): Promise<OperationResult>
  subscribe?(listener: (event: { kind: string }) => void): () => void
}
