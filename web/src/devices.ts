import { computed, ref } from 'vue'
import { api, type Device } from '@/api'

const devices = ref<Device[]>([])
const loading = ref(true)
const error = ref('')

async function load() {
  try {
    devices.value = await api.devices()
    error.value = ''
  } catch (e) {
    error.value = (e as Error).message
  } finally {
    loading.value = false
  }
}

const self = computed(() => devices.value.find((d) => d.self))
const clients = computed(() => devices.value.filter((d) => !d.self))
const direct = computed(() => clients.value.filter((d) => d.connection === 'direct').length)
const relayed = computed(() => clients.value.filter((d) => d.connection === 'relay').length)

export function useDevices() {
  return { devices, loading, error, load, self, clients, direct, relayed }
}
