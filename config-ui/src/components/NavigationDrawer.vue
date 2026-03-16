<script lang="ts" setup>
import { useConfigStore } from '@/store/config';
import { storeToRefs } from 'pinia';
import { Keymap } from "@/types/config";
import { useAppStore } from "@/store/app";
import { server } from "@/store/server";

const { enabledKeymaps, customParentKeymaps, options } = storeToRefs(useConfigStore())
const { translate } = useConfigStore()
const appStore = useAppStore()

async function openWindowSpy() {
  appStore.ingestResult(await server.runWindowSpy())
}

const colors = ["#4c6fff", "#11b67a", "#f06a6a", "#8f5dff", "#1c9ed8", "#ff9f1a", "#ea4c89"]

function getHotkey(keymap: Keymap) {
  return keymap.parentID ? customParentKeymaps.value.find(k => k.id == keymap.parentID)!.hotkey : keymap.hotkey;
}

const getColor = (keymap: Keymap) => {
  let hotkey = getHotkey(keymap)

  let hash = 0;
  for (let i = 0; i < hotkey.length; i++) {
    hash = hotkey.charCodeAt(i) + ((hash << 5) - hash);
  }

  return colors[Math.abs(hash) % colors.length];
}

const getIcon = (keymap: Keymap) => {
  let icon = "mdi-"
  let hotkey = getHotkey((keymap))

  if (hotkey == "settings") {
    return icon + "cog-outline"
  } else if (hotkey == "customHotkeys") {
    return icon + "keyboard-outline"
  } else if (hotkey == "capslockAbbr") {
    return icon + 'rocket-launch-outline'
  } else if (hotkey == "semicolonAbbr") {
    return icon + 'format-text-variant-outline'
  } else if (hotkey.toLowerCase().includes('button')) {
    return icon + 'cursor-default-outline'
  }

  hotkey = hotkey.replace(/^[^!#^+\w]/, '')
  let key = hotkey.substring(0, 1)
  if (/[LlRr]/.test(key)) {
    key = hotkey.substring(1, 2)
  }
  key = key.toLowerCase()

  if (/[^#!^+a-z0-9]/.test(key)) {
    return icon + "rhombus"
  }

  if (/\d/.test(key)) {
    return icon + "numeric-" + key + "-box-outline"
  }

  if (key == "!") {
    key = "a"
  } else if (key == "#") {
    key = "w"
  } else if (key == "^") {
    key = "c"
  } else if (key == "+") {
    key = "s"
  }

  if (/[a-zA-Z]/.test(key)) {
    return icon + "alpha-" + key + "-box"
  }
  return icon + "rhombus"
}
</script>

<template>
  <v-navigation-drawer permanent width="306" class="desktop-drawer" :scrim="false">
    <div class="brand-panel">
      <div class="brand-top">
        <v-avatar size="62" rounded="xl" class="brand-avatar">
          <v-img src="@/assets/logo.png"></v-img>
        </v-avatar>
        <div>
          <div class="brand-label">MyKeymap-re</div>
          <div class="brand-meta">Desktop Remap Studio</div>
        </div>
      </div>

      <div class="brand-stats">
        <div>
          <span>版本</span>
          <strong>{{ options.mykeymapVersion }}</strong>
        </div>
        <div>
          <span>状态</span>
          <strong>{{ appStore.engineLabel }}</strong>
        </div>
      </div>
    </div>

    <div class="drawer-section-label">Workspace</div>
    <v-list class="drawer-list" nav density="comfortable">
      <v-list-item to="/" rounded="xl" prepend-icon="mdi-view-dashboard-outline" title="概览与指南" />
      <v-list-item to="/settings" rounded="xl" prepend-icon="mdi-tune-variant" title="设置工作台" />
    </v-list>

    <div class="drawer-section-label">Keymaps</div>
    <v-list class="drawer-list flex-grow-1" nav density="comfortable">
      <v-virtual-scroll :items="enabledKeymaps" height="calc(100vh - 424px)">
        <template #default="{ item: keymap, index }">
          <v-list-item :key="index" :value="keymap" rounded="xl"
                       :to="keymap.id != 4 ? '/keymap/' + keymap.id : '/' + keymap.hotkey">
            <template #prepend>
              <div class="nav-icon-shell" :style="{ background: getColor(keymap) + '15', color: getColor(keymap) }">
                <v-icon :icon="getIcon(keymap)" size="22"></v-icon>
              </div>
            </template>
            <v-list-item-title>{{ keymap.name }}</v-list-item-title>
          </v-list-item>
        </template>
      </v-virtual-scroll>
    </v-list>

    <div class="drawer-actions">
      <v-btn block variant="outlined" color="secondary" prepend-icon="mdi-radar" @click="openWindowSpy">
        {{ translate('label:309') }}
      </v-btn>
      <v-btn block color="primary" prepend-icon="mdi-lightning-bolt" :disabled="!appStore.canApply" @click="useConfigStore().applyConfig">
        应用修改
      </v-btn>
    </div>
  </v-navigation-drawer>
</template>

<style scoped>
.desktop-drawer {
  border-right: 1px solid rgba(255, 255, 255, 0.08);
  background:
    linear-gradient(180deg, rgba(13, 22, 38, 0.96) 0%, rgba(17, 31, 55, 0.96) 100%),
    radial-gradient(circle at top, rgba(78, 123, 255, 0.16), transparent 25%);
  color: #eef3ff;
  flex-shrink: 0;
}

.brand-panel {
  margin: 14px;
  padding: 16px;
  border: 1px solid rgba(164, 183, 255, 0.16);
  border-radius: 22px;
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.07), rgba(255, 255, 255, 0.03));
}

.brand-top {
  display: flex;
  align-items: center;
  gap: 12px;
}

.brand-avatar {
  background: rgba(255, 255, 255, 0.08);
}

.brand-label {
  color: #f7fbff;
  font-size: 1.28rem;
  font-weight: 700;
}

.brand-meta {
  color: rgba(232, 240, 255, 0.72);
  font-size: 0.82rem;
}

.brand-stats {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
  margin-top: 14px;
}

.brand-stats div {
  padding: 9px 11px;
  border-radius: 14px;
  background: rgba(255, 255, 255, 0.05);
}

.brand-stats span {
  display: block;
  margin-bottom: 4px;
  color: rgba(232, 240, 255, 0.72);
  font-size: 0.78rem;
}

.brand-stats strong {
  color: #ffffff;
}

.drawer-section-label {
  margin: 8px 18px 6px;
  color: rgba(232, 240, 255, 0.58);
  font-size: 0.76rem;
  font-weight: 700;
  letter-spacing: 0.12em;
  text-transform: uppercase;
}

.drawer-list {
  padding: 0 12px;
}

.nav-icon-shell {
  display: flex;
  width: 36px;
  height: 36px;
  align-items: center;
  justify-content: center;
  border-radius: 12px;
}

.drawer-actions {
  display: grid;
  gap: 8px;
  padding: 14px;
}

.desktop-drawer :deep(.v-list-item) {
  min-height: 48px;
  margin-bottom: 4px;
  color: rgba(243, 247, 255, 0.9);
}

.desktop-drawer :deep(.v-list-item--active) {
  background: linear-gradient(90deg, rgba(74, 109, 255, 0.24), rgba(74, 109, 255, 0.1));
  color: #ffffff;
}

.desktop-drawer :deep(.v-list-item-title) {
  font-size: 0.98rem;
  font-weight: 500;
}

.desktop-drawer :deep(.v-navigation-drawer__content) {
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}
</style>
