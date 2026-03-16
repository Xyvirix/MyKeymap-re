<template>
  <v-app v-if="configStore.config" class="desktop-app">
    <div class="desktop-shell">
      <navigation-drawer />

      <div class="desktop-body">
        <header class="desktop-topbar">
          <div>
            <p class="eyebrow">{{ isKeymapRoute ? "Keymap Editor" : "Modern Desktop Workspace" }}</p>
            <h1 class="page-title">{{ currentTitle }}</h1>
          </div>

          <div class="toolbar-actions">
            <v-chip :color="appStore.engineTone" variant="tonal" size="large">
              {{ appStore.engineLabel }}
            </v-chip>
            <v-chip v-if="appStore.dirty" color="warning" variant="outlined" size="large">
              有未应用修改
            </v-chip>
            <v-btn
              class="toolbar-btn"
              variant="outlined"
              color="secondary"
              prepend-icon="mdi-refresh"
              @click="appStore.restartEngine"
            >
              重启引擎
            </v-btn>
            <v-btn
              class="toolbar-btn"
              variant="outlined"
              :color="appStore.engineStatus.paused || !appStore.engineStatus.running ? 'success' : 'warning'"
              :prepend-icon="appStore.engineStatus.paused || !appStore.engineStatus.running ? 'mdi-play' : 'mdi-pause'"
              @click="appStore.toggleEnginePause"
            >
              {{ appStore.engineStatus.paused || !appStore.engineStatus.running ? "恢复映射" : "暂停映射" }}
            </v-btn>
            <v-btn
              class="toolbar-btn"
              variant="outlined"
              color="primary"
              prepend-icon="mdi-content-save-outline"
              :disabled="!appStore.canSave"
              :loading="appStore.isSaving"
              @click="configStore.saveConfig"
            >
              保存
            </v-btn>
            <v-btn
              class="toolbar-btn primary-apply-btn"
              color="primary"
              prepend-icon="mdi-lightning-bolt"
              :disabled="!appStore.canApply"
              :loading="appStore.isApplying"
              @click="configStore.applyConfig"
            >
              应用
            </v-btn>
          </div>
        </header>

        <div v-if="!isKeymapRoute" class="desktop-summary">
          <div class="summary-card">
            <span class="summary-label">版本</span>
            <strong>{{ appStore.info?.version || configStore.options.mykeymapVersion }}</strong>
          </div>
          <div class="summary-card">
            <span class="summary-label">运行模式</span>
            <strong>{{ appStore.info?.desktopMode ? "桌面宿主" : "浏览器调试" }}</strong>
          </div>
          <div class="summary-card">
            <span class="summary-label">进程状态</span>
            <strong>{{ appStore.engineStatus.managed ? "受控" : "外部/未知" }}</strong>
          </div>
          <div class="summary-card">
            <span class="summary-label">最近同步</span>
            <strong>{{ statusTimestamp }}</strong>
          </div>
        </div>

        <main class="desktop-main" :class="{ 'desktop-main--keymap': isKeymapRoute }">
          <router-view v-slot="{ Component }">
            <keep-alive include="HomeSettings">
              <component :is="Component" />
            </keep-alive>
          </router-view>
        </main>

        <footer class="desktop-statusbar">
          <span>Ctrl+S 保存</span>
          <span>Ctrl+Shift+S 保存并应用</span>
          <span v-if="appStore.info?.rootDir">根目录: {{ appStore.info.rootDir }}</span>
        </footer>
      </div>
    </div>

    <v-snackbar v-model="messageVisible" color="primary" location="bottom right" timeout="2600">
      {{ appStore.lastMessage }}
    </v-snackbar>

    <v-snackbar v-model="errorVisible" color="error" location="bottom right" timeout="4200">
      {{ appStore.lastError }}
    </v-snackbar>
  </v-app>

      <v-app v-else class="desktop-loading">
    <div class="loading-panel">
      <v-progress-circular indeterminate color="primary" size="52" />
      <h2>正在加载 MyKeymap-re 配置</h2>
      <p>桌面宿主正在准备配置和引擎状态。</p>
    </div>
  </v-app>
</template>

<script lang="ts" setup>
import { computed, onMounted } from "vue"
import { useRoute } from "vue-router"
import NavigationDrawer from "@/components/NavigationDrawer.vue"
import { useConfigStore } from "@/store/config"
import { useAppStore } from "@/store/app"

const configStore = useConfigStore()
const appStore = useAppStore()
const route = useRoute()

const routeTitleMap: Record<string, string> = {
  "/": "概览与指南",
  "/settings": "设置工作台",
}

const currentTitle = computed(() => {
  if (route.path.startsWith("/keymap/")) {
    return configStore.keymap?.name || "映射编辑器"
  }
  return routeTitleMap[route.path] ?? "MyKeymap-re Desktop"
})

const isKeymapRoute = computed(() => route.path.startsWith("/keymap/"))

const statusTimestamp = computed(() => {
  const stamp = appStore.engineStatus.updatedAt
  return stamp ? new Date(stamp).toLocaleTimeString() : "--"
})

const messageVisible = computed({
  get: () => Boolean(appStore.lastMessage),
  set: (value: boolean) => {
    if (!value) {
      appStore.clearMessage()
    }
  },
})

const errorVisible = computed({
  get: () => Boolean(appStore.lastError),
  set: (value: boolean) => {
    if (!value) {
      appStore.clearError()
    }
  },
})

onMounted(() => {
  void appStore.init()

  window.addEventListener("keydown", async evt => {
    if (evt.ctrlKey && evt.shiftKey && evt.key.toLowerCase() === "s") {
      evt.preventDefault()
      await configStore.applyConfig()
      return
    }

    if (evt.ctrlKey && evt.key.toLowerCase() === "s") {
      evt.preventDefault()
      await configStore.saveConfig()
    }
  })
})
</script>

<style scoped>
.desktop-app {
  background:
    radial-gradient(circle at top right, rgba(61, 116, 255, 0.14), transparent 28%),
    radial-gradient(circle at left bottom, rgba(27, 196, 125, 0.12), transparent 24%),
    linear-gradient(180deg, #f7f9fc 0%, #eef3f8 100%);
  overflow: hidden;
}

.desktop-shell {
  --drawer-width: 306px;
  display: flex;
  height: 100vh;
  overflow: hidden;
}

.desktop-body {
  display: flex;
  flex: 1;
  flex-direction: column;
  width: calc(100% - var(--drawer-width));
  margin-left: var(--drawer-width);
  min-width: 0;
  min-height: 0;
  padding: 18px 20px 12px 14px;
  overflow: hidden;
}

.desktop-topbar {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  padding: 8px 4px 14px;
  flex-shrink: 0;
}

.eyebrow {
  margin: 0 0 6px;
  color: #5c6b84;
  font-size: 0.82rem;
  font-weight: 700;
  letter-spacing: 0.12em;
  text-transform: uppercase;
}

.page-title {
  margin: 0;
  color: #142033;
  font-size: 1.86rem;
  font-weight: 700;
}

.toolbar-actions {
  display: flex;
  flex-wrap: wrap;
  justify-content: flex-end;
  gap: 10px;
  align-items: center;
}

.toolbar-btn {
  min-width: 0;
  padding-inline: 14px;
  letter-spacing: 0;
  white-space: nowrap;
}

.primary-apply-btn {
  min-width: 108px;
}

.desktop-summary {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
  margin-bottom: 14px;
  flex-shrink: 0;
}

.summary-card {
  padding: 14px 16px;
  border: 1px solid rgba(127, 146, 184, 0.18);
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.86);
  box-shadow: 0 18px 36px rgba(39, 68, 120, 0.08);
  backdrop-filter: blur(10px);
}

.summary-label {
  display: block;
  margin-bottom: 6px;
  color: #7a889f;
  font-size: 0.82rem;
  text-transform: uppercase;
  letter-spacing: 0.08em;
}

.desktop-main {
  flex: 1;
  min-height: 0;
  min-width: 0;
  padding-right: 4px;
  overflow: auto;
}

.desktop-main--keymap {
  padding-top: 2px;
}

.desktop-statusbar {
  display: flex;
  flex-wrap: wrap;
  gap: 14px;
  padding: 12px 8px 2px;
  color: #627086;
  font-size: 0.88rem;
  flex-shrink: 0;
}

.desktop-loading {
  background: linear-gradient(180deg, #f6f8fb 0%, #e8edf5 100%);
}

.loading-panel {
  display: flex;
  min-height: 100vh;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  gap: 14px;
  color: #1d2940;
}

.loading-panel h2,
.loading-panel p {
  margin: 0;
}

.desktop-app :deep(.v-application__wrap) {
  min-height: 100vh;
  overflow: hidden;
}

@media (max-width: 1320px) {
  .desktop-shell {
    --drawer-width: 306px;
  }

  .desktop-summary {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .desktop-topbar {
    flex-direction: column;
    align-items: stretch;
  }

  .toolbar-actions {
    justify-content: flex-start;
  }
}

@media (max-width: 960px) {
  .desktop-shell {
    --drawer-width: 0px;
  }

  .desktop-body {
    width: 100%;
    margin-left: 0;
    padding-inline: 12px;
  }
}
</style>
