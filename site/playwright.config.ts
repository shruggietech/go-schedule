import { defineConfig } from '@playwright/test';
export default defineConfig({ testDir: './e2e', workers: 1, retries: 0, reporter: 'list', use: { baseURL: 'http://127.0.0.1:4174', browserName: 'chromium', headless: true, trace: 'retain-on-failure' } });
