<script setup lang="ts">

import { useConfigStore } from "@/store/config";
import { storeToRefs } from "pinia";
import Action from "@/components/actions/Action.vue";
import Key from "@/components/Key.vue";
import { ref } from "vue";
import trimEnd from "lodash-es/trimEnd";
import ActionCommentTable from "@/components/ActionCommentTable.vue";

const { hotkeys } = storeToRefs(useConfigStore())
const { removeHotkey, changeHotkey, translate } = useConfigStore()

const cmd = ref("")

const runCmd = () => {
  let cmdStr = cmd.value.toLowerCase()
  if (!cmdStr || !cmdStr.trim()) {
    return
  }

  if (cmdStr.startsWith("del ")) {
    // 删除
    console.log(cmdStr.substring(4))
    removeHotkey(cmdStr.substring(4))
  } else if (cmdStr.startsWith("rn ")) {
    // 重命名
    cmdStr = cmdStr.substring(3)
    changeHotkey(useConfigStore().hotkey, cmdStr)
  }

  useConfigStore().hotkey = cmdStr.startsWith("del ") ? "" : cmdStr
  cmd.value = ""
}

const formatSpace = (hotkey: string) => {
  const trimmed = trimEnd(hotkey, ' ')
  return trimmed + '◻️'.repeat(hotkey.length - trimmed.length)
}

</script>

<template>
  <div class="abbr-shell">
    <section class="abbr-primary">
      <v-card elevation="0" class="workspace-card keycaps-card">
        <v-card-text class="keycaps-body">
          <div class="abbr-key-grid">
            <key
              v-for="(action, hotkey) in hotkeys"
              :key="hotkey"
              :hotkey="hotkey as string"
              :label="formatSpace(hotkey as string)"
            />
          </div>

          <v-text-field v-model="cmd" @keydown.enter="runCmd()"
                        class="abbr-command-field" variant="underlined" color="primary"
                        :label="translate('label:406')">
          </v-text-field>
        </v-card-text>
      </v-card>
      <action class="abbr-action-panel"></action>
    </section>
    <aside class="abbr-secondary">
      <action-comment-table class="abbr-summary-panel">
        <template #keyText="{hotkey}">
          {{ formatSpace(hotkey as string) }}
        </template>
      </action-comment-table>
    </aside>
  </div>

</template>

<style scoped>
.abbr-shell {
  display: grid;
  grid-template-columns: minmax(0, 1.55fr) minmax(300px, 0.9fr);
  gap: 14px;
  align-items: start;
  min-width: 0;
}

.abbr-primary,
.abbr-secondary {
  min-width: 0;
}

.abbr-secondary {
  position: sticky;
  top: 0;
}

.workspace-card {
  border: 1px solid rgba(127, 146, 184, 0.16);
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.88);
  box-shadow: 0 10px 26px rgba(39, 68, 120, 0.06);
}

.keycaps-card {
  margin-bottom: 12px;
}

.keycaps-body {
  padding: 16px 16px 12px;
}

.abbr-key-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 10px 12px;
}

.abbr-command-field {
  margin-top: 16px;
}

.abbr-action-panel {
  min-width: 0;
}

.abbr-summary-panel {
  min-width: 0;
}

.abbr-summary-panel :deep(.comment-card) {
  padding: 5px;
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.8);
  box-shadow: 0 10px 26px rgba(39, 68, 120, 0.04);
}

.abbr-summary-panel :deep(.Table) {
  border-radius: 16px;
  box-shadow: none;
}

.abbr-summary-panel :deep(th) {
  height: 2.85rem;
  font-size: 0.95rem;
}

.abbr-summary-panel :deep(td) {
  padding-top: 0.8rem;
  padding-bottom: 0.8rem;
  font-size: 0.95rem;
}

@media (max-width: 1600px) {
  .abbr-shell {
    grid-template-columns: 1fr;
  }

  .abbr-secondary {
    position: static;
  }
}

@media (max-width: 960px) {
  .keycaps-body {
    padding: 12px 12px 10px;
  }

  .abbr-key-grid {
    gap: 8px 10px;
  }
}
</style>
