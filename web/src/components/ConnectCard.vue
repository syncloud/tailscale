<script setup lang="ts">
import { onMounted, ref } from 'vue'
import QRCode from 'qrcode'
import { ElMessage } from 'element-plus'
import { CopyDocument, Iphone, Monitor, Cellphone } from '@element-plus/icons-vue'

const props = defineProps<{ url: string }>()
const qr = ref('')
const platform = ref('android')

const steps: Record<string, string[]> = {
  android: [
    'Install Tailscale from Google Play.',
    'Before signing in, open the ⋮ menu → Accounts → Use an alternate server.',
    'Enter the server address and sign in with your device login.',
  ],
  ios: [
    'Install Tailscale from the App Store.',
    'Tap Log in, then the ⋯ menu → Use custom coordination server.',
    'Enter the server address and sign in with your device login.',
  ],
}

async function copy() {
  try {
    await navigator.clipboard.writeText(props.url)
    ElMessage.success('Server address copied')
  } catch {
    ElMessage.warning('Copy is not available here, select the address instead')
  }
}

onMounted(async () => {
  qr.value = await QRCode.toDataURL(props.url, { margin: 1, width: 200 })
})
</script>

<template>
  <section class="ts-card" data-testid="connect-card">
    <h2 class="ts-card-title">Connect a phone or laptop</h2>
    <div class="connect">
      <div class="qr-frame">
        <img v-if="qr" :src="qr" alt="Server address QR code" class="qr" data-testid="server-qr" />
        <div class="qr-caption">Scan to get the server address</div>
      </div>
      <div class="body">
        <div class="server">
          <div class="server-label">Server address</div>
          <div class="server-row">
            <code class="server-url" data-testid="server-url">{{ url }}</code>
            <el-button :icon="CopyDocument" data-testid="copy-server-url" @click="copy">Copy</el-button>
          </div>
        </div>
        <el-tabs v-model="platform" class="tabs">
          <el-tab-pane name="android" data-testid="steps-android">
            <template #label><span class="tab"><el-icon><Cellphone /></el-icon>Android</span></template>
            <ol class="steps">
              <li v-for="(step, i) in steps.android" :key="i"><span class="num">{{ i + 1 }}</span>{{ step }}</li>
            </ol>
          </el-tab-pane>
          <el-tab-pane name="ios" data-testid="steps-ios">
            <template #label><span class="tab"><el-icon><Iphone /></el-icon>iOS</span></template>
            <ol class="steps">
              <li v-for="(step, i) in steps.ios" :key="i"><span class="num">{{ i + 1 }}</span>{{ step }}</li>
            </ol>
          </el-tab-pane>
          <el-tab-pane name="desktop" data-testid="steps-desktop">
            <template #label><span class="tab"><el-icon><Monitor /></el-icon>Desktop</span></template>
            <ol class="steps">
              <li><span class="num">1</span>Install Tailscale for Windows, macOS or Linux.</li>
              <li><span class="num">2</span>Run this and sign in with your device login:</li>
            </ol>
            <code class="command">tailscale login --login-server {{ url }}</code>
          </el-tab-pane>
        </el-tabs>
        <div class="hint">Once connected, your usual Syncloud addresses go straight to this device.</div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.connect {
  display: flex;
  flex-wrap: wrap;
  gap: 24px;
  align-items: flex-start;
}
.qr-frame {
  background: #fff;
  border: 1px solid var(--ts-border);
  border-radius: 14px;
  padding: 12px;
  text-align: center;
}
.qr {
  width: 172px;
  height: 172px;
  display: block;
}
.qr-caption {
  font-size: 12px;
  color: #5b6478;
  margin-top: 6px;
}
.body {
  flex: 1;
  min-width: 260px;
}
.server {
  background: var(--ts-surface-2);
  border: 1px solid var(--ts-border);
  border-radius: 12px;
  padding: 12px 14px;
}
.server-label {
  font-size: 12px;
  color: var(--ts-text-2);
  margin-bottom: 6px;
}
.server-row {
  display: flex;
  align-items: center;
  gap: 10px;
}
.server-url {
  flex: 1;
  font-size: 14px;
  word-break: break-all;
  color: var(--ts-text);
}
.tabs {
  margin-top: 8px;
}
.tab {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}
.steps {
  list-style: none;
  padding: 0;
  margin: 4px 0 0;
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.steps li {
  display: flex;
  gap: 10px;
  align-items: flex-start;
  line-height: 1.5;
  font-size: 14px;
}
.num {
  flex: none;
  width: 22px;
  height: 22px;
  border-radius: 50%;
  background: var(--ts-accent-soft);
  color: var(--ts-accent);
  font-size: 12px;
  font-weight: 700;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  margin-top: 1px;
}
.command {
  display: block;
  margin-top: 10px;
  padding: 10px 12px;
  border-radius: 10px;
  background: #0f131b;
  color: #d6dcff;
  font-size: 13px;
  word-break: break-all;
}
.hint {
  margin-top: 14px;
  font-size: 13px;
  color: var(--ts-text-2);
}
@media (max-width: 640px) {
  .qr-frame {
    margin: 0 auto;
  }
}
</style>
