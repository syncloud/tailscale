<script setup lang="ts">
import { computed, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { Cellphone, Iphone, Monitor, House, Connection } from '@element-plus/icons-vue'
import { api, type Device } from '@/api'
import { useDevices } from '@/devices'

const { devices, loading, error, load } = useDevices()
const removing = ref<Device | null>(null)
const busy = ref(false)

const sorted = computed(() => devices.value)

function icon(device: Device) {
  if (device.self) return House
  const os = device.os.toLowerCase()
  if (os === 'android') return Cellphone
  if (os === 'ios') return Iphone
  if (os) return Monitor
  return Connection
}

function lastSeen(device: Device): string {
  if (device.online) return 'Active now'
  if (!device.last_seen) return 'Never connected'
  const minutes = Math.round((Date.now() / 1000 - device.last_seen) / 60)
  if (minutes < 60) return `Last seen ${minutes} min ago`
  if (minutes < 60 * 48) return `Last seen ${Math.round(minutes / 60)} h ago`
  return `Last seen ${new Date(device.last_seen * 1000).toLocaleDateString()}`
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
</script>

<template>
  <section class="ts-card" data-testid="devices-card">
    <h2 class="ts-card-title">Devices</h2>
    <div v-loading="loading" class="list" data-testid="devices-list">
      <p v-if="error" class="empty">Unable to load devices: {{ error }}</p>
      <p v-else-if="!loading && sorted.length <= 1" class="empty">
        No phones or laptops yet. Connect one with the steps above.
      </p>
      <div v-for="row in sorted" :key="row.id" class="row" :class="{ self: row.self }">
        <div class="avatar" :class="{ online: row.online }">
          <el-icon :size="20"><component :is="icon(row)" /></el-icon>
        </div>
        <div class="info">
          <div class="name-line">
            <span class="name" :data-testid="`device-name-${row.name}`">{{ row.name }}</span>
            <span v-if="row.self" class="ts-pill accent" :data-testid="`device-self-${row.name}`">this device</span>
          </div>
          <div class="meta">
            {{ row.user }} · {{ row.addresses[0] }} · {{ lastSeen(row) }}
          </div>
        </div>
        <div class="status" :data-testid="`device-status-${row.name}`">
          <span v-if="row.online" class="ts-pill good"><span class="ts-dot" />online</span>
          <span v-else class="ts-pill off"><span class="ts-dot" />offline</span>
          <span v-if="row.connection === 'direct'" class="ts-pill good" :title="row.endpoint">direct</span>
          <span v-if="row.connection === 'relay'" class="ts-pill warn">relayed</span>
        </div>
        <div class="actions">
          <el-button
            v-if="!row.self"
            type="danger"
            plain
            size="small"
            :data-testid="`device-remove-${row.name}`"
            @click="removing = row"
          >
            Remove
          </el-button>
        </div>
      </div>
    </div>
    <el-dialog
      :model-value="removing !== null"
      title="Remove device"
      width="min(420px, 92vw)"
      data-testid="remove-dialog"
      @close="removing = null"
    >
      <p>
        Remove <b>{{ removing?.name }}</b>? It will need to sign in again to reconnect.
      </p>
      <template #footer>
        <el-button data-testid="remove-cancel" @click="removing = null">Cancel</el-button>
        <el-button type="danger" :loading="busy" data-testid="remove-confirm" @click="confirmRemove">Remove</el-button>
      </template>
    </el-dialog>
  </section>
</template>

<style scoped>
.list {
  display: flex;
  flex-direction: column;
  gap: 10px;
  min-height: 60px;
}
.row {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 12px 14px;
  border: 1px solid var(--ts-border);
  border-radius: 12px;
  background: var(--ts-surface-2);
  flex-wrap: wrap;
}
.row.self {
  border-color: var(--ts-accent);
  background: var(--ts-accent-soft);
}
.avatar {
  flex: none;
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--ts-surface);
  color: var(--ts-off);
  border: 1px solid var(--ts-border);
}
.avatar.online {
  color: var(--ts-accent);
}
.info {
  flex: 1;
  min-width: 160px;
}
.name-line {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}
.name {
  font-weight: 600;
  font-size: 15px;
}
.meta {
  font-size: 13px;
  color: var(--ts-text-2);
  margin-top: 2px;
}
.status {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}
.actions {
  min-width: 76px;
  display: flex;
  justify-content: flex-end;
}
.empty {
  color: var(--ts-text-2);
  font-size: 14px;
  margin: 4px 0 8px;
}
@media (max-width: 640px) {
  .status {
    order: 3;
    margin-left: 54px;
  }
  .actions {
    order: 4;
    margin-left: auto;
    min-width: 0;
  }
}
</style>
