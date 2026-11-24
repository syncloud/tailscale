import { defineConfig, devices } from '@playwright/test'
import { env } from './e2e/helpers/env'

export default defineConfig({
  testDir: './e2e/specs',
  fullyParallel: false,
  workers: 1,
  retries: process.env.CI ? 1 : 0,
  reporter: 'list',
  use: {
    baseURL: `https://tailscale.${env('PLAYWRIGHT_DOMAIN')}`,
    ignoreHTTPSErrors: true,
  },
  projects: [
    { name: 'desktop', use: { ...devices['Desktop Chrome'], viewport: { width: 1280, height: 900 } } },
    { name: 'mobile', use: { ...devices['Pixel 5'] } },
  ],
})
