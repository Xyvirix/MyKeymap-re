<script lang="ts" setup>
import { computed } from "vue";
import { useConfigStore } from "@/store/config";

const store = useConfigStore()
const props = defineProps<{
  hotkey: string;
  label: string;
  units?: number;
}>();

const isSelected = computed(() => store.hotkey === props.hotkey)
const isMapped = computed(() => !store.keymap?.hotkey.includes("Abbr") && !store.getAction(props.hotkey).isEmpty)

const disabled = computed(() => {
  return store.disabledKeys[store.keymap!.id][props.hotkey.toLowerCase()]
})

const displayUnits = computed(() => props.units ?? 1)

const labelClass = computed(() => {
  if (props.label.length >= 10) {
    return "key-label-xs"
  }
  if (props.label.length >= 7) {
    return "key-label-sm"
  }
  if (props.label.length >= 4) {
    return "key-label-md"
  }
  return "key-label-base"
})

function click(hotkey: string) {
  store.hotkey = hotkey
}


</script>

<template>
  <v-hover v-slot:default="{ isHovering, props }">
    <v-card v-bind="props"
            :elevation="0"
            :disabled="disabled"
            @click="click(hotkey)"
            class="keycap d-flex justify-center align-center"
            :style="{ '--key-units': displayUnits }"
            :class="{
              'is-hovered': isHovering,
              'is-selected': isSelected,
              'is-mapped': isMapped,
              'is-disabled': disabled,
            }">
      <div class="key-label" :class="labelClass">{{ label }}</div>
    </v-card>
  </v-hover>
</template>

<style scoped>
/* 鼠标在 card 之上 hover 时 vuetify 会加一个变暗遮罩, 去掉这个东西 */
:deep(.v-card__overlay) {
  background-color: unset;
}

.keycap {
  --key-unit: 43px;
  --key-height: 42px;
  --key-gap: 6px;
  width: calc((var(--key-unit) * var(--key-units)) + (var(--key-gap) * (var(--key-units) - 1)));
  min-width: calc((var(--key-unit) * var(--key-units)) + (var(--key-gap) * (var(--key-units) - 1)));
  height: var(--key-height);
  padding: 0 9px;
  border: 1px solid rgba(127, 146, 184, 0.18);
  border-radius: 13px;
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.96), rgba(244, 248, 251, 0.94));
  box-shadow: inset 0 -1px 0 rgba(163, 179, 209, 0.18), 0 5px 12px rgba(39, 68, 120, 0.04);
  transition: transform 120ms ease, border-color 120ms ease, background-color 120ms ease;
}

.key-label {
  overflow: hidden;
  text-align: center;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: #26354b;
  line-height: 1;
}

.key-label-base {
  font-size: 1.08rem;
  font-weight: 600;
}

.key-label-md {
  font-size: 0.95rem;
  font-weight: 600;
}

.key-label-sm {
  font-size: 0.84rem;
  font-weight: 600;
}

.key-label-xs {
  font-size: 0.74rem;
  font-weight: 600;
}

.is-hovered:not(.is-disabled) {
  transform: translateY(-1px);
  border-color: rgba(82, 118, 190, 0.32);
}

.is-selected {
  border-color: rgba(76, 111, 255, 0.42);
  background: linear-gradient(180deg, rgba(243, 247, 255, 0.98), rgba(228, 238, 255, 0.96));
  box-shadow: inset 0 -1px 0 rgba(88, 123, 214, 0.22), 0 8px 18px rgba(76, 111, 255, 0.08);
}

.is-mapped:not(.is-selected):not(.is-disabled) {
  border-color: rgba(79, 178, 145, 0.28);
  background: linear-gradient(180deg, rgba(248, 252, 250, 0.98), rgba(237, 247, 243, 0.96));
}

.is-disabled {
  border-color: rgba(161, 172, 191, 0.22);
  background: linear-gradient(180deg, rgba(241, 244, 248, 0.98), rgba(228, 234, 242, 0.96));
  box-shadow: none;
}

.is-disabled .key-label {
  color: rgba(72, 84, 106, 0.7);
}

@media (max-width: 1260px) {
  .keycap {
    --key-unit: 40px;
    --key-height: 40px;
  }
}
</style>
