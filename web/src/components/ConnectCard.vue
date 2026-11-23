<script setup lang="ts">
import { onMounted, ref } from 'vue'
import QRCode from 'qrcode'

const props = defineProps<{ url: string }>()
const qr = ref('')
const platform = ref('android')

onMounted(async () => {
  qr.value = await QRCode.toDataURL(props.url, { margin: 1, width: 180 })
})
</script>

<template>
  <el-card data-testid="connect-card">
    <template #header>Connect a phone or laptop</template>
    <div class="connect">
      <img v-if="qr" :src="qr" alt="Server address QR code" class="qr" data-testid="server-qr" />
      <div class="steps">
        <p>
          Server:
          <code data-testid="server-url">{{ url }}</code>
        </p>
        <el-tabs v-model="platform">
          <el-tab-pane label="Android" name="android" data-testid="steps-android">
            <ol>
              <li>Install Tailscale from Google Play.</li>
              <li>Before signing in, open the menu, Accounts, Use an alternate server.</li>
              <li>Enter the server above and sign in with your device login.</li>
            </ol>
          </el-tab-pane>
          <el-tab-pane label="iOS" name="ios" data-testid="steps-ios">
            <ol>
              <li>Install Tailscale from the App Store.</li>
              <li>Tap Log in, open the menu, Use custom coordination server.</li>
              <li>Enter the server above and sign in with your device login.</li>
            </ol>
          </el-tab-pane>
          <el-tab-pane label="Windows, macOS, Linux" name="desktop" data-testid="steps-desktop">
            <p>Install Tailscale, then run:</p>
            <code>tailscale login --login-server {{ url }}</code>
          </el-tab-pane>
        </el-tabs>
        <p class="hint">Once connected, your usual Syncloud addresses go straight to this device.</p>
      </div>
    </div>
  </el-card>
</template>

<style scoped>
.connect {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  align-items: flex-start;
}
.qr {
  width: 180px;
  height: 180px;
}
.steps {
  flex: 1;
  min-width: 240px;
}
code {
  word-break: break-all;
}
.hint {
  color: var(--el-text-color-secondary);
  font-size: 13px;
}
</style>
