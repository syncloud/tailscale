export type Me = {
  user: string
  url: string
}

export type Device = {
  id: number
  name: string
  user: string
  addresses: string[]
  online: boolean
  last_seen: number
  self: boolean
  connection: '' | 'direct' | 'relay'
  endpoint: string
  os: string
}

export const loginPath = '/app/auth/login'
export const logoutPath = '/app/auth/logout'

async function json<T>(response: Response): Promise<T> {
  if (response.status === 401) {
    window.location.href = loginPath
    throw new Error('sign in required')
  }
  if (!response.ok) {
    throw new Error(`${response.status}: ${await response.text()}`)
  }
  return response.json() as Promise<T>
}

export const api = {
  me: () => fetch('/app/api/me').then(json<Me>),
  devices: () => fetch('/app/api/devices').then(json<Device[]>),
  removeDevice: async (id: number) => {
    const response = await fetch(`/app/api/devices/${id}`, { method: 'DELETE' })
    if (!response.ok) {
      throw new Error(`${response.status}: ${await response.text()}`)
    }
  },
}
