<script setup lang="ts">

import Table from "@/components/Table.vue";
import { storeToRefs } from "pinia";
import { useConfigStore } from "@/store/config";
import { ref } from "vue";
import ActionView from "@/components/actions/Action.vue";
import { Action } from "@/types/config";

const { hotkeys, windowGroupID } = storeToRefs(useConfigStore())
const { translate } = useConfigStore()

const currHotkey = ref<string>("")
const changeCustomHotkey = (hotkey: string, newHotkey: string) => {
  useConfigStore().changeHotkey(hotkey, newHotkey)
  checkRow(newHotkey)
}

const checkRow = (hotkey: string, windowGroupId?: number) => {
  currHotkey.value = hotkey
  useConfigStore().hotkey = hotkey
  if (windowGroupId != undefined) {
    useConfigStore().windowGroupID = windowGroupId
  }
}

const removeCustomHotkey = (hotkey: string) => {
  useConfigStore().removeHotkey(hotkey);
  if (currHotkey.value == hotkey) {
    checkRow("")
  }
}

const getActionComment = (action: Array<Action>) => {
  return action.find(a => !a.isEmpty && a.windowGroupID == windowGroupID.value)?.comment ?? ''
}

const getActionWindowGroupId = () => {
  return windowGroupID.value
}

</script>

<template>
  <div class="custom-hotkey-shell">
    <section class="hotkey-list-panel">
      <v-card elevation="0" class="workspace-card hotkey-list-card">
        <Table class="text-left" :titles="[translate('label:404'), translate('label:305'), '']">
          <tr :class="currHotkey == hotkey ? 'bg-blue-lighten-4' : ''"
              @click="checkRow(hotkey as string, getActionWindowGroupId())"
              v-for="(action, hotkey, index) in hotkeys" :key="index">
            <td style="width: 20%">
              <v-text-field :model-value="hotkey" placeholder="此处修改"
                            @change="changeCustomHotkey(hotkey as string, $event.target.value)"
                            variant="plain" class="hotkey-field"></v-text-field>
            </td>
            <td style="width: 60%; cursor: pointer;" class="hotkey-comment-cell"><div class="hotkey-comment-text">{{ translate(getActionComment(action)) }}</div></td>
            <td style="width: 20%">
              <v-btn icon="mdi-delete-outline" variant="text" width="40" height="40"
                     @click.stop="removeCustomHotkey(hotkey as string)"></v-btn>
            </td>
          </tr>
        </Table>

        <div class="d-flex justify-end">
          <v-btn class="ma-3 text-none" color="green" @click="useConfigStore().addHotKey()">{{ translate('label:405') }}</v-btn>
        </div>
      </v-card>
    </section>
    <section class="custom-editor-panel">
      <div class="example-grid">
        <v-card elevation="0" class="workspace-card example-card" title="Example 1">
          <v-card-text>
              <p>!c = Alt + C</p>
              <p>#c = Win + C</p>
              <p>^c = Ctrl + C</p>
              <p>^!c = Ctrl + Alt + C</p>
              <p>^+c = Ctrl + Shift + C</p>
              <p>+!c = Shift + Alt + C</p>
          </v-card-text>
        </v-card>
        <v-card elevation="0" class="workspace-card example-card" title="Example 2">
          <v-card-text>
              <p>F11 = F11</p>
              <p>!1 &nbsp;= Alt + 1</p>
              <p>+F2 = Shift + F2</p>
              <p>!space = Alt + Space</p>

              <br>
              <p>更多特殊按键参考: <a target="_blank"
                                  href="https://wyagd001.github.io/v2/docs/KeyList.htm#keyboard"
                                  style="color: green; text-decoration: none">reference</a></p>
          </v-card-text>
        </v-card>
      </div>
      <ActionView class="custom-action-view" />
    </section>
  </div>
</template>

<style scoped>
.custom-hotkey-shell {
  display: grid;
  grid-template-columns: minmax(320px, 420px) minmax(0, 1fr);
  gap: 16px;
  align-items: start;
  min-width: 0;
}

.workspace-card {
  border: 1px solid rgba(127, 146, 184, 0.16);
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.86);
  box-shadow: 0 12px 30px rgba(39, 68, 120, 0.07);
}

.hotkey-list-panel,
.custom-editor-panel {
  min-width: 0;
}

.hotkey-list-card {
  overflow: hidden;
}

.example-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
  margin-bottom: 14px;
}

.example-card :deep(.v-card-title) {
  padding-bottom: 0;
  font-weight: 700;
}

.custom-action-view {
  min-width: 0;
}

.hotkey-field {
  width: 100%;
  min-width: 0;
}

.hotkey-comment-cell {
  cursor: pointer;
}

.hotkey-comment-text {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

table .v-text-field :deep(input) {
  min-height: auto;
  padding: 0 !important;
}

table .v-text-field :deep(.v-input__details) {
  min-height: auto;
  height: 0 !important;
}

table .v-switch :deep(.v-selection-control) {
  min-height: auto;
}

table .v-select :deep(.v-field__input) {
  padding: 0;
}

table .v-autocomplete :deep(input) {
  top: 13px
}

@media (max-width: 1440px) {
  .custom-hotkey-shell {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 900px) {
  .example-grid {
    grid-template-columns: 1fr;
  }
}
</style>
