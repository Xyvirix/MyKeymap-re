import type { Config } from "./config"

export interface EngineStatus {
  running: boolean
  paused: boolean
  managed: boolean
  pid?: number
  lastError?: string
  updatedAt: string
}

export interface ActionResult {
  message: string
  status: EngineStatus
}

export interface Shortcut {
  path: string
}

export interface AppInfo {
  name: string
  version: string
  rootDir: string
  desktopMode: boolean
  runAtStartup: boolean
  frontendSource: string
}

export interface DesktopBinding {
  GetConfig(): Promise<Config>
  SaveConfig(config: Config): Promise<ActionResult>
  ApplyConfig(config: Config): Promise<ActionResult>
  GetEngineStatus(): Promise<EngineStatus>
  StartEngine(): Promise<ActionResult>
  RestartEngine(): Promise<ActionResult>
  PauseEngine(): Promise<ActionResult>
  ResumeEngine(): Promise<ActionResult>
  ListShortcuts(): Promise<Shortcut[]>
  RunWindowSpy(): Promise<ActionResult>
  SetRunAtStartup(enabled: boolean): Promise<ActionResult>
  GetAppInfo(): Promise<AppInfo>
}
