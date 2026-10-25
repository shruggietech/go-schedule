import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  build: { target: 'es2022', emptyOutDir: false },
  test: { testTimeout: 60_000, hookTimeout: 60_000, teardownTimeout: 60_000, environment: 'jsdom', setupFiles: './src/test-setup.ts', css: true, exclude: ['e2e/**', 'node_modules/**'] },
})
