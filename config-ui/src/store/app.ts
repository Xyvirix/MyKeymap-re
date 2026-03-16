import { defineStore } from "pinia"
import { computed, ref } from "vue"
import { desktopApi } from "./server"
import type { Config } from "@/types/config"
import type { ActionResult, AppInfo, EngineStatus } from "@/types/desktop"

export const useAppStore = defineStore("app", () => {
  const info = ref<AppInfo>()
  const engineStatus = ref<EngineStatus>({
    running: false,
    paused: false,
    managed: false,
    updatedAt: new Date().toISOString(),
  })
  const lastMessage = ref("")
  const lastError = ref("")
  const isLoading = ref(false)
  const isSaving = ref(false)
  const isApplying = ref(false)
  const dirty = ref(false)
  const lastPersistedSnapshot = ref("")
  const pollHandle = ref<number>()

  const engineLabel = computed(() => {
    if (engineStatus.value.paused) {
      return "已暂停"
    }
    return engineStatus.value.running ? "运行中" : "未运行"
  })

  const engineTone = computed(() => {
    if (engineStatus.value.paused) {
      return "warning"
    }
    return engineStatus.value.running ? "success" : "error"
  })

  const canApply = computed(() => dirty.value && !isApplying.value)
  const canSave = computed(() => dirty.value && !isSaving.value)

  async function init() {
    isLoading.value = true
    try {
      info.value = await desktopApi.getAppInfo()
      engineStatus.value = await desktopApi.getEngineStatus()
      startPolling()
    } catch (error) {
      handleError(error)
    } finally {
      isLoading.value = false
    }
  }

  async function refreshStatus() {
    try {
      engineStatus.value = await desktopApi.getEngineStatus()
    } catch (error) {
      handleError(error)
    }
  }

  function startPolling() {
    if (pollHandle.value) {
      window.clearInterval(pollHandle.value)
    }
    pollHandle.value = window.setInterval(() => {
      void refreshStatus()
    }, 3500)
  }

  function setBaseline(config: Config) {
    lastPersistedSnapshot.value = JSON.stringify(config)
    dirty.value = false
  }

  function updateDirty(config: Config) {
    dirty.value = JSON.stringify(config) !== lastPersistedSnapshot.value
  }

  async function saveConfig(config: Config) {
    isSaving.value = true
    try {
      const result = await desktopApi.saveConfig(config)
      consumeResult(result)
      setBaseline(config)
    } catch (error) {
      handleError(error)
      throw error
    } finally {
      isSaving.value = false
    }
  }

  async function applyConfig(config: Config) {
    isApplying.value = true
    try {
      const result = await desktopApi.applyConfig(config)
      consumeResult(result)
      setBaseline(config)
    } catch (error) {
      handleError(error)
      throw error
    } finally {
      isApplying.value = false
    }
  }

  async function toggleEnginePause() {
    try {
      const result = engineStatus.value.paused || !engineStatus.value.running
        ? await desktopApi.resumeEngine()
        : await desktopApi.pauseEngine()
      consumeResult(result)
    } catch (error) {
      handleError(error)
    }
  }

  async function restartEngine() {
    try {
      const result = await desktopApi.restartEngine()
      consumeResult(result)
    } catch (error) {
      handleError(error)
    }
  }

  async function startEngine() {
    try {
      const result = await desktopApi.startEngine()
      consumeResult(result)
    } catch (error) {
      handleError(error)
    }
  }

  function consumeResult(result: ActionResult) {
    engineStatus.value = result.status
    lastMessage.value = result.message
    lastError.value = result.status.lastError ?? ""
  }

  function ingestResult(result: ActionResult) {
    consumeResult(result)
  }

  function handleError(error: unknown) {
    lastError.value = error instanceof Error ? error.message : String(error)
    if (!lastMessage.value) {
      lastMessage.value = "操作失败"
    }
  }

  function clearMessage() {
    lastMessage.value = ""
  }

  function clearError() {
    lastError.value = ""
  }

  return {
    info,
    engineStatus,
    lastMessage,
    lastError,
    dirty,
    isLoading,
    isSaving,
    isApplying,
    engineLabel,
    engineTone,
    canApply,
    canSave,
    init,
    refreshStatus,
    setBaseline,
    updateDirty,
    saveConfig,
    applyConfig,
    toggleEnginePause,
    restartEngine,
    startEngine,
    ingestResult,
    clearMessage,
    clearError,
  }
})
