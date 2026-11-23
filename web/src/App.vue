<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { api, logoutPath, type Me } from '@/api'
import ConnectCard from '@/components/ConnectCard.vue'
import DevicesCard from '@/components/DevicesCard.vue'

const me = ref<Me | null>(null)

onMounted(async () => {
  me.value = await api.me()
})
</script>

<template>
  <div class="page">
    <header class="header">
      <h1 class="title">Tailscale</h1>
      <div v-if="me" class="account">
        <span data-testid="user-name">{{ me.user }}</span>
        <a :href="logoutPath" data-testid="logout">Sign out</a>
      </div>
    </header>
    <main v-if="me" class="main">
      <ConnectCard :url="me.url" />
      <DevicesCard />
    </main>
  </div>
</template>

<style>
body {
  margin: 0;
  font-family: system-ui, -apple-system, 'Segoe UI', Roboto, sans-serif;
  background: var(--el-bg-color-page);
  color: var(--el-text-color-primary);
}
</style>

<style scoped>
.page {
  max-width: 960px;
  margin: 0 auto;
  padding: 16px;
}
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}
.title {
  font-size: 22px;
  margin: 8px 0;
}
.account {
  display: flex;
  gap: 12px;
  font-size: 14px;
}
.main {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
</style>
