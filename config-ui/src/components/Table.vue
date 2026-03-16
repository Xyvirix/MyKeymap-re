<script setup lang="ts">
defineProps<{
  titles?: string[]
}>()
</script>

<template>
  <div class="Table" >
    <table class="my-table">
      <slot name="title">
        <tr>
          <th v-for="title in titles" :key="title">{{title}}</th>
        </tr>
      </slot>
      <slot></slot>
    </table>
  </div>
</template>

<style scoped>
.Table {
  overflow-x: auto;
  overflow-y: hidden;
  border: 1px solid rgba(127, 146, 184, 0.16);
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.92);
  box-shadow: 0 10px 28px rgba(39, 68, 120, 0.05);
}

table {
  border-collapse: collapse;
  width: 100%;
  min-width: 0;
  table-layout: fixed;
}

th, :slotted(th) {
  font-weight: 500;
  color: rgba(var(--v-theme-on-surface),var(--v-medium-emphasis-opacity));

  padding: 0 0.875rem;
  height: 3.15rem;
  border-bottom: 1px solid #e4e4e4;
}

:slotted(td) {
  padding: 0 0.875rem;
  height: 3.15rem;
  border-bottom: 1px solid #e4e4e4;
}

:slotted(tr:hover:not(:first-child)) {
  background: rgba(76, 111, 255, 0.1);
}

th {
  color: #18253a;
  font-size: 1rem;
  background-color: #f7f9fc;
}

:deep(.v-select .v-field__append-inner i) {
  visibility: hidden;
}

:deep(tr) {
  color: rgba(0, 0, 0, 0.96);
}

:deep(.v-field--disabled .v-select__selection) {
  opacity: 0.5;
}

</style>
