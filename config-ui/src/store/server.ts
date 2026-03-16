import type { Config } from "@/types/config"
import type { ActionResult, AppInfo, DesktopBinding, EngineStatus, Shortcut } from "@/types/desktop"

declare global {
  interface Window {
    go?: {
      desktop?: {
        DesktopApp?: DesktopBinding
      }
      main?: {
        DesktopApp?: DesktopBinding
      }
    }
  }
}

const fallbackStatus = (): EngineStatus => ({
  running: true,
  paused: false,
  managed: false,
  updatedAt: new Date().toISOString(),
})

const desktopBinding = (): DesktopBinding | undefined =>
  window.go?.desktop?.DesktopApp ?? window.go?.main?.DesktopApp

export const isDesktopRuntime = () => Boolean(desktopBinding())

const isWailsRuntime = () => typeof window !== "undefined" && window.location.origin.startsWith("http://wails.localhost")

const wait = (ms: number) => new Promise(resolve => window.setTimeout(resolve, ms))

const resolveDesktopBinding = async (): Promise<DesktopBinding | undefined> => {
  const binding = desktopBinding()
  if (binding) {
    return binding
  }

  if (!isWailsRuntime()) {
    return undefined
  }

  for (let i = 0; i < 20; i++) {
    await wait(100)
    const next = desktopBinding()
    if (next) {
      return next
    }
  }

  throw new Error("Wails 桌面绑定未初始化")
}

const browserJson = async <T>(path: string, init?: RequestInit): Promise<T> => {
  const response = await fetch(path, init)
  if (!response.ok) {
    throw new Error(`${response.status} ${response.statusText}`)
  }
  return response.json() as Promise<T>
}

const browserCommand = async (path: string, init?: RequestInit) => {
  const response = await fetch(path, init)
  if (!response.ok) {
    throw new Error(`${response.status} ${response.statusText}`)
  }
}

const withStatus = async (message: string): Promise<ActionResult> => ({
  message,
  status: await desktopApi.getEngineStatus(),
})

export const desktopApi = {
  async getConfig(): Promise<Config> {
    const binding = await resolveDesktopBinding()
    if (binding) {
      return binding.GetConfig()
    }
    return browserJson<Config>("/config")
  },

  async saveConfig(config: Config): Promise<ActionResult> {
    const binding = await resolveDesktopBinding()
    if (binding) {
      return binding.SaveConfig(config)
    }
    await browserCommand("/config", {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(config),
    })
    return withStatus("配置已保存")
  },

  async applyConfig(config: Config): Promise<ActionResult> {
    const binding = await resolveDesktopBinding()
    if (binding) {
      return binding.ApplyConfig(config)
    }
    await browserCommand("/config", {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(config),
    })
    return withStatus("配置已应用")
  },

  async getEngineStatus(): Promise<EngineStatus> {
    const binding = await resolveDesktopBinding()
    if (binding) {
      return binding.GetEngineStatus()
    }
    return fallbackStatus()
  },

  async startEngine(): Promise<ActionResult> {
    const binding = await resolveDesktopBinding()
    if (binding) {
      return binding.StartEngine()
    }
    return withStatus("开发模式不接管引擎")
  },

  async restartEngine(): Promise<ActionResult> {
    const binding = await resolveDesktopBinding()
    if (binding) {
      return binding.RestartEngine()
    }
    return withStatus("开发模式不接管引擎")
  },

  async pauseEngine(): Promise<ActionResult> {
    const binding = await resolveDesktopBinding()
    if (binding) {
      return binding.PauseEngine()
    }
    return withStatus("开发模式不接管引擎")
  },

  async resumeEngine(): Promise<ActionResult> {
    const binding = await resolveDesktopBinding()
    if (binding) {
      return binding.ResumeEngine()
    }
    return withStatus("开发模式不接管引擎")
  },

  async listShortcuts(): Promise<Shortcut[]> {
    const binding = await resolveDesktopBinding()
    if (binding) {
      return binding.ListShortcuts()
    }
    return browserJson<Shortcut[]>("/shortcuts")
  },

  async runWindowSpy(): Promise<ActionResult> {
    const binding = await resolveDesktopBinding()
    if (binding) {
      return binding.RunWindowSpy()
    }
    await browserCommand("/server/command/2", { method: "POST" })
    return withStatus("WindowSpy 已启动")
  },

  async setRunAtStartup(enabled: boolean): Promise<ActionResult> {
    const binding = await resolveDesktopBinding()
    if (binding) {
      return binding.SetRunAtStartup(enabled)
    }
    await browserCommand(enabled ? "/server/command/3" : "/server/command/4", { method: "POST" })
    return withStatus(enabled ? "已开启开机启动" : "已关闭开机启动")
  },

  async getAppInfo(): Promise<AppInfo> {
    const binding = await resolveDesktopBinding()
    if (binding) {
      return binding.GetAppInfo()
    }
    return {
      name: "MyKeymap Browser Settings",
      version: "",
      rootDir: "",
      desktopMode: false,
      runAtStartup: false,
      frontendSource: location.origin,
    }
  },
}

export const server = {
  runWindowSpy: () => desktopApi.runWindowSpy(),
  enableRunAtStartup: () => desktopApi.setRunAtStartup(true),
  disableRunAtStartup: () => desktopApi.setRunAtStartup(false),
}
