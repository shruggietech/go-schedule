export type BundleItem = { kind: string; portable_id: string; name: string; action: string; message?: string }
export type BundlePlan = { id: string; bundle_digest: string; target_daemon_id: string; target_fingerprint: string; items: BundleItem[] }
export type BundleIssue = { kind: string; identity?: string; message: string }
export type BundleResult = { action: string; outcome: 'accepted' | 'rejected' | 'unavailable' | 'uncertain'; message: string; document?: string; plan?: BundlePlan; compare_only?: boolean; issues?: BundleIssue[]; items?: BundleItem[] }
export type BundleBridge = { exportBundle(): Promise<BundleResult>; validate(document: string): Promise<BundleResult>; compare(document: string): Promise<BundleResult>; preview(document: string, watcherPaths?: Record<string, string>): Promise<BundleResult>; apply(plan: BundlePlan): Promise<BundleResult> }
