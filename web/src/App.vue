<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { api, logoutPath, type Me } from '@/api'
import { useDevices } from '@/devices'
import SummaryStrip from '@/components/SummaryStrip.vue'
import ConnectCard from '@/components/ConnectCard.vue'
import DevicesCard from '@/components/DevicesCard.vue'

const me = ref<Me | null>(null)
const { load } = useDevices()
let timer: number | undefined

onMounted(async () => {
  me.value = await api.me()
  await load()
  timer = window.setInterval(load, 15000)
})
onUnmounted(() => window.clearInterval(timer))
</script>

<template>
  <div class="shell">
    <header class="bar">
      <div class="bar-inner">
        <div class="brand">
          <svg class="logo" viewBox="0 0 24 24" aria-hidden="true">
            <circle cx="5" cy="5" r="2.2" />
            <circle cx="12" cy="5" r="2.2" />
            <circle cx="19" cy="5" r="2.2" opacity="0.35" />
            <circle cx="5" cy="12" r="2.2" />
            <circle cx="12" cy="12" r="2.2" />
            <circle cx="19" cy="12" r="2.2" />
            <circle cx="5" cy="19" r="2.2" opacity="0.35" />
            <circle cx="12" cy="19" r="2.2" />
            <circle cx="19" cy="19" r="2.2" opacity="0.35" />
          </svg>
          <span class="brand-name">Tailscale</span>
        </div>
        <div v-if="me" class="account">
          <span class="user" data-testid="user-name">{{ me.user }}</span>
          <a :href="logoutPath" class="logout" data-testid="logout">Sign out</a>
        </div>
      </div>
    </header>
    <main v-if="me" class="content">
      <SummaryStrip :url="me.url" />
      <div class="grid">
        <ConnectCard :url="me.url" />
        <DevicesCard />
      </div>
    </main>
  </div>
</template>

<style scoped>
.bar {
  background: var(--ts-surface);
  border-bottom: 1px solid var(--ts-border);
  position: sticky;
  top: 0;
  z-index: 10;
}
.bar-inner {
  max-width: 1040px;
  margin: 0 auto;
  padding: 12px 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.brand {
  display: flex;
  align-items: center;
  gap: 10px;
}
.logo {
  width: 26px;
  height: 26px;
  fill: var(--ts-accent);
}
.brand-name {
  font-size: 17px;
  font-weight: 700;
  letter-spacing: -0.01em;
}
.account {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 14px;
}
.user {
  color: var(--ts-text-2);
}
.logout {
  color: var(--ts-accent);
  text-decoration: none;
  font-weight: 500;
  padding: 6px 12px;
  border: 1px solid var(--ts-border);
  border-radius: 8px;
}
.logout:hover {
  background: var(--ts-accent-soft);
}
.content {
  max-width: 1040px;
  margin: 0 auto;
  padding: 20px 16px 40px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.grid {
  display: grid;
  grid-template-columns: minmax(0, 1fr);
  gap: 16px;
}
</style>
