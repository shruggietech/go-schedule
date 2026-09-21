import type { BundleBridge, BundleResult } from './model'

type NativeWindow = Window & { go?: { main?: { App?: Record<string, (...args: unknown[]) => Promise<BundleResult>> } } }
const unavailable = (action: string): BundleResult => ({ action, outcome: 'unavailable', message: 'Portable bundles are unavailable.' })

export function createBundleBridge(nativeWindow: NativeWindow = window): BundleBridge {
  const app = nativeWindow.go?.main?.App
  const call = (name: string, action: string, ...args: unknown[]) => app?.[name]?.(...args) ?? Promise.resolve(unavailable(action))
  return { exportBundle: () => call('PortableBundleExport', 'export_bundle'), validate: (document) => call('PortableBundleValidate', 'validate_bundle', document), compare: (document) => call('PortableBundleCompare', 'compare_bundle', document), preview: (document) => call('PortableBundlePreview', 'preview_bundle', document), apply: (plan) => call('PortableBundleApply', 'apply_bundle', plan) }
}

export const bundleBridge = createBundleBridge()
