<script setup lang="ts">

import Table from "@/components/Table.vue";
import { Action as IAction } from "@/types/config";
import trimEnd from "lodash-es/trimEnd";
import { useConfigStore } from "@/store/config";
import { computed } from "vue";

const store = useConfigStore();
const { translate } = useConfigStore();
const getActionAllComment = (actions: IAction[]) => {
  let comment = actions.reduce((pre, current) => {
    if (current.comment) {
      let groupName = current.windowGroupID == 0 ? "" : store.options.windowGroups.find(w => w.id == current.windowGroupID)?.name + ": "
      return pre + groupName + translate(current.comment) + "\r\n"
    }
    return pre
  }, "")

  return trimEnd(comment, "\n")
}

const sortedHotkeys = computed(() => {
  if (!store.hotkeys) {
    return []
  }
  return Object.entries(store.hotkeys)
    // .filter(([_, actions]) => showActionComment(actions))
    .map(([hotkey, actions]) => ({
      hotkey,
      comment: getActionAllComment(actions)
    }))
    .sort((a, b) => a.comment.localeCompare(b.comment))
})

</script>

<template>
  <v-card elevation="0" class="comment-card">
    <Table class="text-left" :titles="[translate('label:404'), translate('label:305')]">
      <tr v-for="(item, index) in sortedHotkeys" :key="index">
        <td v-if="item.comment" class="hotkey-cell">
          <slot name="keyText" :hotkey="item.hotkey"></slot>
        </td>
        <td v-if="item.comment" class="comment-cell">
          {{ item.comment }}
        </td>
      </tr>
    </Table>
  </v-card>
</template>

<style scoped>
.comment-card {
  width: 100%;
  padding: 6px;
  border: 1px solid rgba(127, 146, 184, 0.16);
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.82);
  box-shadow: 0 12px 30px rgba(39, 68, 120, 0.05);
}

.hotkey-cell {
  white-space: nowrap;
  width: 28%;
}

.comment-cell {
  white-space: pre-wrap;
  word-break: break-word;
  line-height: 1.45;
  color: #3f4c62;
}

:deep(td) {
  height: auto;
  padding-top: 10px;
  padding-bottom: 10px;
  border-bottom-color: #e4e4e4aa;
}
</style>
