import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  build: { target: 'es2022' },
  test: { environment: 'jsdom', setupFiles: './src/test-setup.ts', css: true, exclude: ['e2e/**', 'node_modules/**'] },
})
