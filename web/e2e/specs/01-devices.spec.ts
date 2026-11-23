import { test, expect } from '../fixtures'
import { shoot } from '../helpers/screenshot'
import { signIn } from '../helpers/syncloud'
import { env } from '../helpers/env'

const domain = env('PLAYWRIGHT_DOMAIN')
const device = domain.split('.')[0]

test('signs in and shows how to connect', async ({ page }, info) => {
  await signIn(page)

  await expect(page.getByTestId('user-name')).toHaveText(env('DEVICE_USER'))
  await expect(page.getByTestId('server-url')).toHaveText(`https://tailscale.${domain}`)
  await expect(page.getByTestId('server-qr')).toBeVisible()
  await page.getByRole('tab', { name: 'iOS' }).click()
  await expect(page.getByTestId('steps-ios')).toBeVisible()
  await shoot(page, `connect-${info.project.name}`)
})

test('lists this device as online and not removable', async ({ page }, info) => {
  await signIn(page)

  await expect(page.getByTestId(`device-name-${device}`)).toBeVisible()
  await expect(page.getByTestId(`device-self-${device}`)).toBeVisible()
  await expect(page.getByTestId(`device-status-${device}`)).toContainText('online')
  await expect(page.getByTestId(`device-remove-${device}`)).toHaveCount(0)
  await expect(page.getByTestId('summary-state')).toContainText('Online')
  await shoot(page, `devices-${info.project.name}`)
})

test('mobile shows devices as a list with every column on screen', async ({ page }, info) => {
  test.skip(info.project.name !== 'mobile', 'mobile layout')
  await signIn(page)

  await expect(page.getByTestId('devices-list')).toBeVisible()
  const status = page.getByTestId(`device-status-${device}`)
  await status.scrollIntoViewIfNeeded()
  await expect(status).toBeInViewport({ ratio: 1 })
  await shoot(page, 'devices-list-mobile')
})

test('removes a client device', async ({ page }, info) => {
  test.skip(info.project.name !== 'desktop', 'one removal per run')
  await signIn(page)

  await expect(page.getByTestId('device-name-ci')).toBeVisible()
  await page.getByTestId('device-remove-ci').click()
  await expect(page.getByTestId('remove-confirm')).toBeVisible()
  await shoot(page, 'remove-confirm')
  await page.getByTestId('remove-confirm').click()

  await expect(page.getByTestId('device-name-ci')).toHaveCount(0)
  await expect(page.getByTestId(`device-name-${device}`)).toBeVisible()
  await shoot(page, 'device-removed')
})
