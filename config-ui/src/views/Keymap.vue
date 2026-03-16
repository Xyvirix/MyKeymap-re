<script lang="ts" setup>
import Action from "@/components/actions/Action.vue";
import Key from "@/components/Key.vue";
import { parseKeyboardLayout, useConfigStore } from "@/store/config";
import { computed } from "vue";
import ActionCommentTable from "@/components/ActionCommentTable.vue";
import trimStart from "lodash-es/trimStart";

const store = useConfigStore();

const keyboardRows = computed(() => {
  const rows = parseKeyboardLayout(store.options.keyboardLayout, store.keymap!.hotkey)
  // console.log(rows)
  return rows
})

function getKeyText(hotkey: string) {
  const string = trimStart(hotkey, '*')
  return string.charAt(0).toUpperCase() + string.slice(1)
}

const keyWidthMap: Record<string, number> = {
  esc: 1.2,
  tab: 1.6,
  capslock: 1.85,
  enter: 1.9,
  backspace: 2.35,
  space: 2.8,
  singlepress: 2.4,
  lshift: 2.15,
  rshift: 2.45,
  lctrl: 1.45,
  rctrl: 1.45,
  lalt: 1.35,
  ralt: 1.35,
  lwin: 1.35,
  rwin: 1.35,
  printscreen: 1.8,
  scrolllock: 1.8,
  pause: 1.55,
  insert: 1.4,
  home: 1.35,
  pgup: 1.35,
  delete: 1.4,
  end: 1.2,
  pgdn: 1.35,
  numpadenter: 1.65,
  numpad0: 1.75,
}

function getKeyUnits(hotkey: string) {
  const normalized = trimStart(hotkey, "*").toLowerCase()
  if (keyWidthMap[normalized]) {
    return keyWidthMap[normalized]
  }
  if (normalized.includes("button")) {
    return 1.55
  }
  if (normalized.startsWith("wheel")) {
    return 1.45
  }
  if (normalized.length >= 9) {
    return 1.75
  }
  if (normalized.length >= 6) {
    return 1.45
  }
  if (normalized.length >= 4) {
    return 1.2
  }
  return 1
}
</script>

<template>
  <div class="workspace-shell" v-if="store.keymap">
    <section class="keyboard-stage">
      <div class="keyboard-panel">
        <div v-for="(row, index) in keyboardRows" :key="index" class="keyboard-row">
          <Key
            v-for="(hotkey, keyIndex) in row"
            :key="hotkey + keyIndex"
            :hotkey="hotkey"
            :label="getKeyText(hotkey)"
            :units="getKeyUnits(hotkey)"
          />
        </div>
      </div>
    </section>
    <aside class="inspector-panel">
      <Action class="action-panel" />
      <action-comment-table class="summary-panel">
        <template #keyText="{ hotkey }">
          {{ getKeyText(hotkey) }}
        </template>
      </action-comment-table>
    </aside>
  </div>
  <div v-else>Error: keymap not found</div>
</template>

<style scoped>
.workspace-shell {
  display: grid;
  grid-template-columns: minmax(0, 1.68fr) minmax(340px, 0.92fr);
  gap: 14px;
  align-items: start;
  min-width: 0;
}

.keyboard-stage {
  min-width: 0;
  padding: 12px 14px 14px;
  border: 1px solid rgba(127, 146, 184, 0.16);
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.88);
  box-shadow: 0 10px 26px rgba(39, 68, 120, 0.06);
}

.keyboard-panel {
  display: flex;
  flex-direction: column;
  gap: 7px;
  min-width: 0;
}

.keyboard-row {
  display: flex;
  flex-wrap: nowrap;
  align-items: center;
  gap: 6px;
  min-width: 0;
}

.inspector-panel {
  position: sticky;
  top: 0;
  display: grid;
  gap: 12px;
  min-width: 0;
}

.action-panel {
  margin-top: 0;
}

.summary-panel {
  min-width: 0;
}

.summary-panel :deep(.comment-card) {
  padding: 5px;
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.76);
  box-shadow: 0 10px 26px rgba(39, 68, 120, 0.04);
}

.summary-panel :deep(.Table) {
  border-radius: 16px;
  box-shadow: none;
}

.summary-panel :deep(th) {
  height: 2.85rem;
  font-size: 0.95rem;
}

.summary-panel :deep(td) {
  padding-top: 0.8rem;
  padding-bottom: 0.8rem;
  font-size: 0.95rem;
}

@media (max-width: 1500px) {
  .workspace-shell {
    grid-template-columns: 1fr;
  }

  .inspector-panel {
    position: static;
  }
}

@media (max-width: 960px) {
  .workspace-shell {
    gap: 12px;
  }

  .keyboard-stage {
    padding: 10px 10px 12px;
  }
}
</style>
