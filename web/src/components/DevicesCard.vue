<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { api, type Device } from '@/api'

const devices = ref<Device[]>([])
const loading = ref(true)
const removing = ref<Device | null>(null)
const busy = ref(false)
let timer: number | undefined

async function load() {
  try {
    devices.value = await api.devices()
  } catch (e) {
    ElMessage.error(`Unable to load devices: ${(e as Error).message}`)
  } finally {
    loading.value = false
  }
}

async function confirmRemove() {
  if (!removing.value) return
  busy.value = true
  try {
    await api.removeDevice(removing.value.id)
    ElMessage.success(`${removing.value.name} removed`)
    removing.value = null
    await load()
  } catch (e) {
    ElMessage.error(`Unable to remove: ${(e as Error).message}`)
  } finally {
    busy.value = false
  }
}

function lastSeen(device: Device): string {
  if (device.online) return 'now'
  if (!device.last_seen) return 'never'
  return new Date(device.last_seen * 1000).toLocaleString()
}

onMounted(() => {
  load()
  timer = window.setInterval(load, 15000)
})
onUnmounted(() => window.clearInterval(timer))
</script>

<template>
  <el-card data-testid="devices-card">
    <template #header>Devices</template>
    <el-table v-loading="loading" :data="devices" row-key="id" data-testid="devices-table" empty-text="No devices yet">
      <el-table-column label="Name" min-width="140">
        <template #default="{ row }">
          <span :data-testid="`device-name-${row.name}`">{{ row.name }}</span>
          <el-tag v-if="row.self" size="small" class="self" :data-testid="`device-self-${row.name}`">this device</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="user" label="User" min-width="100" />
      <el-table-column label="Address" min-width="120">
        <template #default="{ row }">{{ row.addresses[0] }}</template>
      </el-table-column>
      <el-table-column label="Status" min-width="160">
        <template #default="{ row }">
          <span :data-testid="`device-status-${row.name}`">
            <el-tag v-if="row.online" type="success" size="small">online</el-tag>
            <el-tag v-else type="info" size="small">offline</el-tag>
            <el-tag v-if="row.connection === 'direct'" type="success" size="small" class="path">direct</el-tag>
            <el-tag v-if="row.connection === 'relay'" type="warning" size="small" class="path">relayed</el-tag>
          </span>
        </template>
      </el-table-column>
      <el-table-column label="Last seen" min-width="140">
        <template #default="{ row }">{{ lastSeen(row) }}</template>
      </el-table-column>
      <el-table-column width="100" align="right">
        <template #default="{ row }">
          <el-button v-if="!row.self" type="danger" link :data-testid="`device-remove-${row.name}`" @click="removing = row">
            Remove
          </el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-dialog
      :model-value="removing !== null"
      title="Remove device"
      width="min(420px, 92vw)"
      data-testid="remove-dialog"
      @close="removing = null"
    >
      <p>
        Remove <b>{{ removing?.name }}</b>? It will have to sign in again to reconnect.
      </p>
      <template #footer>
        <el-button data-testid="remove-cancel" @click="removing = null">Cancel</el-button>
        <el-button type="danger" :loading="busy" data-testid="remove-confirm" @click="confirmRemove">Remove</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<style scoped>
.self,
.path {
  margin-left: 6px;
}
</style>
