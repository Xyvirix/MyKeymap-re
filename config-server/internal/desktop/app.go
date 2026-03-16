package desktop

import (
	"context"
	"os"
	"path/filepath"
	"sync"

	"settings/internal/script"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type ActionResult struct {
	Message string       `json:"message"`
	Status  EngineStatus `json:"status"`
}

type Shortcut struct {
	Path string `json:"path"`
}

type AppInfo struct {
	Name           string `json:"name"`
	Version        string `json:"version"`
	RootDir        string `json:"rootDir"`
	DesktopMode    bool   `json:"desktopMode"`
	RunAtStartup   bool   `json:"runAtStartup"`
	FrontendSource string `json:"frontendSource"`
}

type DesktopApp struct {
	ctx   context.Context
	paths Paths

	engine *EngineController

	mu sync.Mutex
}

func NewDesktopApp() (*DesktopApp, error) {
	paths, err := DiscoverPaths()
	if err != nil {
		return nil, err
	}

	return &DesktopApp{
		paths:  paths,
		engine: NewEngineController(paths),
	}, nil
}

func (a *DesktopApp) Paths() Paths {
	return a.paths
}

func (a *DesktopApp) Startup(ctx context.Context) {
	a.ctx = ctx
	_, _ = a.GetAppInfo()
	if status, err := a.engine.Status(); err == nil && !status.Running && !status.Paused {
		_, _ = a.StartEngine()
	}
}

func (a *DesktopApp) DomReady(ctx context.Context) {
	a.emitStatus()
}

func (a *DesktopApp) Shutdown(ctx context.Context) {
	_ = a.engine.Stop()
}

func (a *DesktopApp) BeforeClose(ctx context.Context) bool {
	if a.ctx == nil {
		return false
	}

	runtime.WindowHide(a.ctx)
	runtime.EventsEmit(a.ctx, "app-hidden")
	return true
}

func (a *DesktopApp) ShowWindow() {
	if a.ctx == nil {
		return
	}
	runtime.WindowUnminimise(a.ctx)
	runtime.WindowShow(a.ctx)
	runtime.WindowCenter(a.ctx)
}

func (a *DesktopApp) GetConfig() (*script.Config, error) {
	return script.ParseConfig(a.paths.ConfigFile)
}

func (a *DesktopApp) SaveConfig(config *script.Config) (*ActionResult, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if err := script.SaveConfigFileTo(a.paths.ConfigFile, config); err != nil {
		return a.actionResult("保存失败", err)
	}

	return a.actionResult("配置已保存", nil)
}

func (a *DesktopApp) ApplyConfig(config *script.Config) (*ActionResult, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if err := script.SaveConfigFileTo(a.paths.ConfigFile, config); err != nil {
		return a.actionResult("保存失败", err)
	}
	if err := script.GenerateDesktopScriptsAt(a.paths.TemplatesDir, a.paths.BinDir, config); err != nil {
		return a.actionResult("生成脚本失败", err)
	}
	if err := a.engine.Restart(); err != nil {
		return a.actionResult("热重载失败", err)
	}

	return a.actionResult("配置已应用并热重载", nil)
}

func (a *DesktopApp) GetEngineStatus() (*EngineStatus, error) {
	status, err := a.engine.Status()
	return &status, err
}

func (a *DesktopApp) StartEngine() (*ActionResult, error) {
	if err := a.engine.Start(); err != nil {
		return a.actionResult("启动失败", err)
	}
	return a.actionResult("MyKeymap 正在运行", nil)
}

func (a *DesktopApp) RestartEngine() (*ActionResult, error) {
	if err := a.engine.Restart(); err != nil {
		return a.actionResult("重启失败", err)
	}
	return a.actionResult("MyKeymap 已重启", nil)
}

func (a *DesktopApp) PauseEngine() (*ActionResult, error) {
	if err := a.engine.Pause(); err != nil {
		return a.actionResult("暂停失败", err)
	}
	return a.actionResult("MyKeymap 已暂停", nil)
}

func (a *DesktopApp) ResumeEngine() (*ActionResult, error) {
	if err := a.engine.Resume(); err != nil {
		return a.actionResult("恢复失败", err)
	}
	return a.actionResult("MyKeymap 已恢复", nil)
}

func (a *DesktopApp) ListShortcuts() ([]Shortcut, error) {
	pattern := filepath.Join(a.paths.RootDir, "shortcuts", "*.lnk")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	res := make([]Shortcut, 0, len(files))
	for _, f := range files {
		res = append(res, Shortcut{
			Path: f[len(a.paths.RootDir)+1:],
		})
	}
	return res, nil
}

func (a *DesktopApp) RunWindowSpy() (*ActionResult, error) {
	cmd := hiddenCommand(a.paths.RootDir, a.paths.MyKeymapExe, "/script", "./bin/WindowSpy.ahk")
	if err := cmd.Start(); err != nil {
		return a.actionResult("WindowSpy 启动失败", err)
	}
	return a.actionResult("WindowSpy 已启动", nil)
}

func (a *DesktopApp) SetRunAtStartup(enabled bool) (*ActionResult, error) {
	state := "Off"
	if enabled {
		state = "On"
	}

	cmd := hiddenCommand(a.paths.RootDir, a.paths.MyKeymapExe, "/script", "./bin/MiscTools.ahk", "RunAtStartup", state)
	if err := cmd.Run(); err != nil {
		return a.actionResult("开机启动设置失败", err)
	}

	msg := "已关闭开机启动"
	if enabled {
		msg = "已开启开机启动"
	}
	return a.actionResult(msg, nil)
}

func (a *DesktopApp) GetAppInfo() (*AppInfo, error) {
	config, err := script.ParseConfig(a.paths.ConfigFile)
	if err != nil {
		return nil, err
	}

	frontendSource := a.paths.SiteDir
	if ok, _ := hasFile(a.paths.SiteDir, "index.html"); !ok {
		frontendSource = filepath.Join(a.paths.RootDir, "config-ui", "dist")
	}

	return &AppInfo{
		Name:           "MyKeymap-re Desktop",
		Version:        config.Options.MykeymapVersion,
		RootDir:        a.paths.RootDir,
		DesktopMode:    true,
		RunAtStartup:   isRunAtStartupEnabled(),
		FrontendSource: frontendSource,
	}, nil
}

func (a *DesktopApp) actionResult(message string, err error) (*ActionResult, error) {
	status, statusErr := a.engine.Status()
	if statusErr != nil && err == nil {
		err = statusErr
	}
	if err != nil {
		status.LastError = err.Error()
	}
	a.emitStatus()
	return &ActionResult{
		Message: message,
		Status:  status,
	}, err
}

func (a *DesktopApp) emitStatus() {
	if a.ctx == nil {
		return
	}
	status, err := a.engine.Status()
	if err == nil {
		runtime.EventsEmit(a.ctx, "engine-status", status)
	}
}

func isRunAtStartupEnabled() bool {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		return false
	}
	linkFile := filepath.Join(appData, "Microsoft", "Windows", "Start Menu", "Programs", "Startup", "MyKeymap.lnk")
	_, err := os.Stat(linkFile)
	return err == nil
}
