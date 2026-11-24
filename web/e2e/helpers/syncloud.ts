import { expect, type Page } from '@playwright/test'
import { env } from './env'

export async function signIn(page: Page): Promise<void> {
  const username = page.locator('#username-textfield')
  const app = page.getByTestId('devices-card')

  await page.goto('/')
  await expect(username.or(app)).toBeVisible({ timeout: 60_000 })
  if (await app.isVisible()) return

  await username.fill(env('DEVICE_USER'))
  await page.locator('#password-textfield').fill(env('DEVICE_PASSWORD'))
  await page.locator('#sign-in-button').click()

  const notNow = page.getByRole('button', { name: 'Not now' })
  const consent = page.locator('#accept-button')
  await expect(notNow.or(consent).or(app)).toBeVisible({ timeout: 60_000 })
  if (await notNow.isVisible()) {
    await notNow.click()
  }
  if (await consent.isVisible()) {
    await consent.click()
  }
  await expect(app).toBeVisible({ timeout: 60_000 })
}
