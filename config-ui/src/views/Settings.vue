<script setup lang="ts">
import Table from "@/components/Table.vue";
import Tip from "@/components/Tip.vue";

import { storeToRefs } from "pinia";
import { ref } from "vue";
import { useConfigStore } from "@/store/config";
import { Keymap } from "@/types/config";
// import PathDialog from "@/components/dialog/PathDialog.vue";
import WindowGroupDialog from "@/components/dialog/WindowGroupDialog.vue";
import findLastIndex from "lodash-es/findLastIndex";
import { server } from "@/store/server";
import { computed } from "vue";
import { languageList } from "@/store/language-map";
import { useAppStore } from "@/store/app";

const { customKeymaps, customParentKeymaps, customSonKeymaps, options, keymaps } = storeToRefs(useConfigStore())
const { translate } = useConfigStore()
const appStore = useAppStore()


const currId = ref(0)

const showMouseOption = ref(false)
const showLanguageOption = ref(false)
const showKeyboardLayout = ref(false)
const showKeymapDelay = ref(true)
const showSkin = ref(false)
const resetOtherToFalse = (newValue: boolean) => {
  [showMouseOption, showLanguageOption, showKeyboardLayout, showKeymapDelay, showSkin].forEach(x => x.value = false)
  return newValue
}

const skin = computed(() => {
  return [
    [
      { key: "windowWidth", label: translate('label:743'), },
      { key: "windowYPos", label: translate('label:744'), },
      { key: "borderRadius", label: translate('label:745'), },
      { key: "hideAnimationDuration", label: translate('label:746'), },
    ],
    [
      { key: "backgroundColor", label: translate('label:747'), },
      { key: "backgroundOpacity", label: translate('label:748'), },
      { key: "gridlineColor", label: translate('label:749'), },
      { key: "gridlineOpacity", label: translate('label:748'), },
    ],
    [
      { key: "borderWidth", label: translate('label:750'), },
      { key: "borderColor", label: translate('label:751'), },
      { key: "borderOpacity", label: translate('label:748'), },
    ],
    [
      { key: "keyColor", label: translate('label:752') },
      { key: "keyOpacity", label: translate('label:748'), },
      { key: "cornerColor", label: translate('label:753'), },
      { key: "cornerOpacity", label: translate('label:748'), },
    ],
    [
      { key: "windowShadowSize", label: translate('label:754'), },
      { key: "windowShadowColor", label: translate('label:755'), },
      { key: "windowShadowOpacity", label: translate('label:748'), },
    ],
]
})

const checkKeymapData = (keymap: Keymap) => {
  if (keymap.hotkey == "") {
    currId.value = keymap.id
  }
  // 判断当前热键是否已存在，已存在删除当前模式
  const f = keymaps.value.find(k => k.hotkey == keymap.hotkey && k.parentID == keymap.parentID)!
  if (f.id != keymap.id && keymap.hotkey) {
    removeKeymap(keymap.id)
  }
  currId.value = f.id
  // 把 bs, esc 这样的非标准键名替换成 Backspace, Escape
  keymap.hotkey = normalizeKeyName(keymap.hotkey)
}

const disabledKeymapOption = (keymap: Keymap) => {
  // 状态为启动时、被作为前置键不允许删除
  if (keymap.enable) {
    return true
  }
  return customSonKeymaps.value.findIndex(k => k.parentID == keymap.id) != -1
}

const hasSubKeymap = (keymap: Keymap) => {
  return customSonKeymaps.value.findIndex(k => k.parentID == keymap.id) != -1
}

const deleteBtnTip = (keymap: Keymap) => {
  if (keymap.enable) {
    return '开启时不允许删除'
  }
  if (customSonKeymaps.value.findIndex(k => k.parentID == keymap.id) != -1) {
    return '被依赖时不允许删除'
  }
  return ''
}

function toggleKeymapEnable(keymap: Keymap) {
  // 开启的keymap有前置键连同前置键一块开启
  if (!keymap.enable && keymap.parentID != 0) {
    customParentKeymaps.value.find(k => k.id == keymap.parentID)!.enable = true
  }

  // 关闭的时候连同子键一块关闭
  if (keymap.enable) {
    customSonKeymaps.value.filter(k => k.parentID == keymap.id).forEach(k => k.enable = false)
  }

  keymap.enable = !keymap.enable
  useConfigStore().changeAbbrEnable()
}

function nextKeymapId() {
  const length = customKeymaps.value.length;
  if (length == 0) {
    return 5
  }

  return customKeymaps.value[length - 1].id + 1
}

function addKeymap() {
  const newKeymap: Keymap = {
    id: nextKeymapId(),
    name: "",
    enable: false,
    hotkey: "",
    parentID: 0,
    delay: 0,
    isNew: true,
    hotkeys: {}
  }

  keymaps.value.splice(customKeymaps.value.length, 0, newKeymap)
}

function removeKeymap(id: number) {
  removeKeymapByIndex(findLastIndex(keymaps.value, k => k.id == id))
}

function removeKeymapByIndex(index: number) {
  keymaps.value.splice(index, 1)
}

async function onStartupChange() {
  options.value.startup = !options.value.startup
  if (options.value.startup) {
    appStore.ingestResult(await server.enableRunAtStartup())
  } else {
    appStore.ingestResult(await server.disableRunAtStartup())
  }
}

function normalizeKeyName(hotkey: string) : string {
  const m = {
    'esc': 'Escape',
    'bs': 'Backspace',
    'del': 'Delete',
    'ins': 'Insert',
    'lctrl': 'LControl',
    'rctrl': 'RControl',
  } as any

  for (const [k,v] of Object.entries(m)) {
    m[ '*' + k] = '*' + v
  }

  const v = m[hotkey.toLowerCase()]
  return v ? v : hotkey
}
</script>

<template>
  <v-container :fluid="true" class="settings-shell">
    <v-row class="settings-layout">
      <v-col cols="12" lg="7" xl="7" class="settings-primary">
        <v-card elevation="0" class="settings-card keymap-card">
          <Table class="text-left keymap-table" :titles="[translate('label:501'), translate('label:502'), translate('label:503'), translate('label:504')]">
            <tr :class="currId == keymap.id ? '' : ''"
                @click="currId = keymap.id"
                v-for="keymap in customKeymaps" :key="keymap.id">
              <td>
                <v-text-field v-model.lazy="keymap.name" @blur="checkKeymapData(keymap)"
                              variant="plain" class="keymap-name-field"></v-text-field>
              </td>
              <td>
                <v-text-field v-model.lazy="keymap.hotkey" @blur="checkKeymapData(keymap)"
                              variant="plain" class="keymap-hotkey-field"></v-text-field>
              </td>
              <td>
                <v-select v-model="keymap.parentID" :items="customParentKeymaps.filter(c => c.id != keymap.id)"
                          :item-title="item => item.name"
                          :item-value="item => item.id" :disabled="hasSubKeymap(keymap)"
                          item-color="blue"
                          variant="plain" class="keymap-parent-field">
                </v-select>
              </td>
              <td>
                <div class="table-actions">
                  <v-switch hide-details color="primary" :model-value="keymap.enable"
                            @click="toggleKeymapEnable(keymap)"></v-switch>
                  <tip :text="deleteBtnTip(keymap)">
                    <v-btn icon="mdi-delete-outline" variant="text" width="40" height="40"
                           :disabled="disabledKeymapOption(keymap)"
                           @click="removeKeymap(keymap.id)"></v-btn>
                  </tip>
                </div>
              </td>
            </tr>
          </Table>

          <div class="d-flex justify-end">
            <v-btn class="ma-3 text-none" color="green" @click="addKeymap()">{{ translate('label:405') }}</v-btn>
          </div>
        </v-card>
      </v-col>
      <v-col cols="12" lg="5" xl="5" class="settings-secondary">
        <div class="otherSetting">
          <v-row :dense="true">
            <v-col>
              <v-card :title="translate('label:505')" class="settings-card section-card">
                <v-card-text>
                  <v-switch :label="translate('label:506')" color="primary"
                            :model-value="options.startup"
                            @change="onStartupChange"></v-switch>
                  <div class="option-shortcuts">
                    <v-btn class="option-trigger text-none" color="blue" variant="outlined" @click="showLanguageOption = resetOtherToFalse(!showLanguageOption)">{{ translate('label:781') }}</v-btn>
                    <window-group-dialog class="option-trigger-host"/>
                    <v-btn class="option-trigger text-none" color="blue" variant="outlined" @click="showMouseOption = resetOtherToFalse(!showMouseOption)">{{ translate('label:701') }}</v-btn>
                    <v-btn class="option-trigger text-none" color="blue" variant="outlined" @click="showKeyboardLayout = resetOtherToFalse(!showKeyboardLayout)">{{ translate('label:721') }}</v-btn>
                    <v-btn class="option-trigger text-none" color="blue" variant="outlined" @click="showSkin = resetOtherToFalse(!showSkin)">{{ translate('label:741') }}</v-btn>
                    <v-btn class="option-trigger text-none" color="blue" variant="outlined" @click="showKeymapDelay = resetOtherToFalse(!showKeymapDelay)">{{ translate('label:761') }}</v-btn>
                  </div>
                </v-card-text>
              </v-card>
            </v-col>
          </v-row>
            <v-row :dense="true" v-show="showMouseOption">
              <v-col>
                <v-card :title="translate('label:702')" class="settings-card section-card">
                  <v-card-text>
                    <v-row class="mouseRow" no-gutters>
                      <v-col>
                        <v-text-field v-model="options.mouse.delay1" variant="underlined"
                                      type="number" step=".01" maxlength="5" color="primary"
                                      :label="translate('label:703')"></v-text-field>
                      </v-col>
                      <v-col>
                        <v-text-field v-model="options.mouse.delay2" variant="underlined"
                                      type="number" step=".01" maxlength="5" color="primary"
                                      :label="translate('label:704')"></v-text-field>
                      </v-col>
                    </v-row>
                    <v-row class="mouseRow" no-gutters>
                      <v-col>
                        <v-text-field v-model="options.mouse.fastRepeat" variant="underlined"
                                      type="number" step="1" maxlength="5" color="primary"
                                      :label="translate('label:705')"></v-text-field>
                      </v-col>
                      <v-col>
                        <v-text-field v-model="options.mouse.fastSingle" variant="underlined"
                                      type="number" step="1" maxlength="5" color="primary"
                                      :label="translate('label:706')"></v-text-field>
                      </v-col>
                    </v-row>
                    <v-row class="mouseRow" no-gutters>
                      <v-col>
                        <v-text-field v-model="options.mouse.slowRepeat" variant="underlined"
                                      type="number" step="1" maxlength="5" color="primary"
                                      :label="translate('label:707')"></v-text-field>
                      </v-col>
                      <v-col>
                        <v-text-field v-model="options.mouse.slowSingle" variant="underlined"
                                      type="number" step="1" maxlength="5" color="primary"
                                      :label="translate('label:708')"></v-text-field>
                      </v-col>
                    </v-row>
                    <v-row class="mouseRow" no-gutters>
                      <v-col>
                        <v-text-field v-model="options.mouse.tipSymbol" variant="underlined" color="primary" :label="translate('label:709')"></v-text-field>
                      </v-col>
                      <v-col>
                        <br>
                        <!-- <v-label>备选符号: 🖱️🔘</v-label> -->
                      </v-col>
                    </v-row>
                    <v-row>
                      <v-col>
                        <v-checkbox :label="translate('label:710')" color="secondary" hide-details density="compact" v-model="options.mouse.showTip" />
                        <v-checkbox :label="translate('label:711')" color="secondary" hide-details density="compact" v-model="options.mouse.keepMouseMode" />
                      </v-col>
                    </v-row>
                  </v-card-text>
                </v-card>
              </v-col>
              <v-col>
                <v-card :title="translate('label:712')" class="settings-card section-card compact-card">
                  <v-card-text>
                    <v-text-field v-model="options.scroll.delay1" variant="underlined"
                                  type="number" step=".01" maxlength="5" color="primary"
                                  :label="translate('label:713')"></v-text-field>
                    <v-text-field v-model="options.scroll.delay2" variant="underlined"
                                  type="number" step=".01" maxlength="5" color="primary"
                                  :label="translate('label:714')"></v-text-field>
                    <v-text-field v-model="options.scroll.onceLineCount" variant="underlined"
                                  type="number" step="1" maxlength="5" color="primary"
                                  :label="translate('label:715')"></v-text-field>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>
            <v-row v-show="showKeyboardLayout">
              <v-col>
                <v-card :title="translate('label:722')" elevation="0" class="settings-card section-card">
                  <v-card-text>
                    <v-textarea color="primary" variant="underlined" auto-grow rows="4" v-model="options.keyboardLayout"></v-textarea>
                  </v-card-text>
                  <v-card-actions class="d-flex justify-end">
                    <v-btn class="text-none" variant="outlined" color="green" @click="useConfigStore().resetKeyboardLayout(0)">{{ translate('label:723') }}</v-btn>
                    <v-btn class="text-none" variant="outlined" color="green" @click="useConfigStore().resetKeyboardLayout(74)">{{ translate('label:724') }}</v-btn>
                    <v-btn class="text-none" variant="outlined" color="green" @click="useConfigStore().resetKeyboardLayout(104)">{{ translate('label:725') }}</v-btn>
                    <v-btn class="text-none" variant="outlined" color="blue" @click="useConfigStore().resetKeyboardLayout(1)">{{ translate('label:726') }}</v-btn>
                  </v-card-actions>
                </v-card>
              </v-col>
            </v-row>
            <v-row v-show="showLanguageOption">
              <v-col>
                <v-card elevation="0" class="settings-card section-card compact-card">
                  <v-card-text>
                    <v-select :items="languageList" v-model="options.language" variant="outlined"></v-select>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>
            <v-row v-show="showKeymapDelay">
              <v-col>
                <v-card :title="translate('label:762')" elevation="0" class="settings-card section-card">
                  <v-card-text>
                    {{ translate('label:763') }}<br>
                    {{ translate('label:764') }}<br>
                    {{ translate('label:765') }}<br>
                    &nbsp;
                    <v-row>
                      <v-col cols="3" v-for="keymap in customKeymaps" :key="keymap.id">
                        <v-text-field v-model.number="keymap.delay" variant="underlined"
                                      type="number" step="1" maxlength="5" min="0" color="primary"
                                      :label="keymap.name" :class="{'positive-number': keymap.delay > 0}"></v-text-field>
                      </v-col>
                    </v-row>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>
            <v-row v-show="showSkin">
              <v-col>
                <v-card :title="translate('label:742')" elevation="0" class="settings-card section-card">
                  <v-card-text>
                    <v-row v-for="(row, index) in skin" :key="index">
                      <v-col cols="3" v-for="item in row" :key="item.key">
                        <v-text-field v-model="options.commandInputSkin[item.key]" variant="underlined" color="primary" :label="item.label"></v-text-field>
                      </v-col>
                    </v-row>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>
        </div>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>
.settings-shell {
  width: 100%;
  max-width: none;
  padding: 0 2px 8px 4px;
}

.settings-layout {
  margin: 0 !important;
  align-items: start;
}

.settings-primary,
.settings-secondary {
  min-width: 0;
  padding: 5px !important;
}

.settings-card {
  border: 1px solid rgba(127, 146, 184, 0.16);
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.92);
  box-shadow: 0 8px 22px rgba(39, 68, 120, 0.06);
}

.keymap-card {
  width: 100%;
  max-width: 100%;
  overflow: hidden;
}

.keymap-card :deep(.Table) {
  border: 0;
  border-radius: 0;
  background: transparent;
  box-shadow: none;
}

.keymap-card :deep(table) {
  table-layout: fixed;
}

.keymap-card :deep(th),
.keymap-card :deep(td) {
  padding-inline: 10px;
}

.keymap-card :deep(th:nth-child(1)),
.keymap-card :deep(td:nth-child(1)) {
  width: 34%;
}

.keymap-card :deep(th:nth-child(2)),
.keymap-card :deep(td:nth-child(2)) {
  width: 18%;
}

.keymap-card :deep(th:nth-child(3)),
.keymap-card :deep(td:nth-child(3)) {
  width: 24%;
}

.keymap-card :deep(th:nth-child(4)),
.keymap-card :deep(td:nth-child(4)) {
  width: 24%;
}

.keymap-name-field,
.keymap-hotkey-field,
.keymap-parent-field {
  width: 100%;
  min-width: 0;
}

.table-actions {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 6px;
}

.section-card :deep(.v-card-text) {
  padding: 6px 10px 8px;
}

.compact-card :deep(.v-card-text) {
  padding: 5px 10px 6px;
}

.option-shortcuts {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 6px;
  margin-top: 6px;
  align-items: start;
}

.option-trigger {
  justify-content: flex-start;
  min-height: 36px;
  padding-inline: 10px;
  font-size: 0.95rem;
  letter-spacing: 0;
  white-space: nowrap;
}

.option-trigger-host {
  display: contents;
}

table .v-text-field :deep(input) {
  min-height: auto;
  padding: 0 !important;
}

table .v-text-field :deep(.v-input__details) {
  min-height: auto;
  height: 0 !important;
}

table .v-text-field :deep(.v-field--disabled) {
  opacity: 1 !important;
}

.v-text-field :deep(label),.v-switch :deep(label),.v-checkbox :deep(label) {
  color: black;
  opacity: 1;
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

.otherSetting .v-card {
  margin-bottom: 5px;
}

.mouseRow .v-col:first-child {
  padding-right: 10px;
}

.positive-number {
  color: #d05;
}

:deep(.v-field) {
  border-radius: 16px;
}

:deep(.v-card-title) {
  font-weight: 700;
  padding: 12px 14px 2px;
}

.settings-shell :deep(.v-row) {
  margin: 0;
}

.settings-shell :deep(.v-col) {
  min-width: 0;
  padding: 5px;
}

.settings-shell :deep(.v-btn) {
  letter-spacing: 0.01em;
}

.settings-shell :deep(.v-field__input) {
  min-height: 38px;
}

.settings-shell :deep(.v-switch .v-selection-control) {
  min-height: 34px;
}

.settings-shell :deep(.v-selection-control) {
  align-items: center;
}

@media (max-width: 1380px) {
  .option-shortcuts {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 1440px) {
  .settings-primary,
  .settings-secondary {
    flex: 0 0 100% !important;
    max-width: 100% !important;
  }
}

@media (max-width: 960px) {
  .settings-shell {
    padding-inline: 2px;
  }

  .option-shortcuts {
    grid-template-columns: 1fr;
  }
}
</style>
