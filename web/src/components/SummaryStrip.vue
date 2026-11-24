<script setup lang="ts">
import { useDevices } from '@/devices'

defineProps<{ url: string }>()
const { self, clients, direct, relayed } = useDevices()
</script>

<template>
  <section class="hero" data-testid="summary">
    <div class="hero-main">
      <div class="hero-label">This device</div>
      <div class="hero-name">{{ self?.name ?? '…' }}</div>
      <div class="hero-state" data-testid="summary-state">
        <span class="dot" :class="{ on: self?.online }" />
        {{ self?.online ? 'Online and accepting connections' : 'Starting…' }}
      </div>
    </div>
    <div class="stats">
      <div class="stat">
        <div class="stat-value" data-testid="stat-devices">{{ clients.length }}</div>
        <div class="stat-label">connected devices</div>
      </div>
      <div class="stat">
        <div class="stat-value" data-testid="stat-direct">{{ direct }}</div>
        <div class="stat-label">direct</div>
      </div>
      <div class="stat">
        <div class="stat-value" data-testid="stat-relayed">{{ relayed }}</div>
        <div class="stat-label">relayed</div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.hero {
  background: var(--ts-hero);
  color: #fff;
  border-radius: 16px;
  padding: 22px 24px;
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  box-shadow: var(--ts-shadow);
}
.hero-label {
  font-size: 12px;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  opacity: 0.75;
}
.hero-name {
  font-size: 24px;
  font-weight: 700;
  margin: 2px 0 6px;
}
.hero-state {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  opacity: 0.95;
}
.dot {
  width: 9px;
  height: 9px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.5);
}
.dot.on {
  background: #5cf2a0;
  box-shadow: 0 0 0 4px rgba(92, 242, 160, 0.25);
}
.stats {
  display: flex;
  gap: 12px;
}
.stat {
  background: rgba(255, 255, 255, 0.14);
  border-radius: 12px;
  padding: 10px 16px;
  min-width: 84px;
  text-align: center;
}
.stat-value {
  font-size: 22px;
  font-weight: 700;
}
.stat-label {
  font-size: 12px;
  opacity: 0.85;
}
@media (max-width: 640px) {
  .stats {
    width: 100%;
  }
  .stat {
    flex: 1;
    min-width: 0;
    padding: 10px 6px;
  }
}
</style>
