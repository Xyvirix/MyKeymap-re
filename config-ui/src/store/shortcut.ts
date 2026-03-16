import { defineStore } from 'pinia'
import { ref } from 'vue'
import { desktopApi } from './server'

export const useShortcutStore = defineStore('shortcut', () => {
  const shortcuts = fetchShortcuts()
  return { shortcuts }
})

const fetchShortcuts = () => {
  const shortcuts = ref<string[]>()
  void desktopApi.listShortcuts().then((items) => {
    shortcuts.value = items.map(x => x.path)
  }).catch((error) => {
    console.error(error)
    shortcuts.value = []
  })
  return shortcuts
}
