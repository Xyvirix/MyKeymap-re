package desktop

import (
	"os"
	"sync"

	"github.com/getlantern/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type trayController struct {
	once      sync.Once
	app       *DesktopApp
	pauseItem *systray.MenuItem
}

func (a *DesktopApp) startTray() {
	controller := &trayController{app: a}
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.ctx == nil {
		return
	}
	controller.once.Do(func() {
		go systray.Run(controller.onReady, func() {})
	})
}

func (a *DesktopApp) stopTray() {
	systray.Quit()
}

func (t *trayController) onReady() {
	if icon, err := os.ReadFile(t.app.paths.TrayIcon); err == nil {
		systray.SetIcon(icon)
	}
	systray.SetTooltip("MyKeymap Desktop")
	systray.SetTitle("MyKeymap")

	openItem := systray.AddMenuItem("打开 MyKeymap", "显示主窗口")
	reloadItem := systray.AddMenuItem("热重载配置", "重启映射引擎")
	t.pauseItem = systray.AddMenuItem("暂停映射", "暂停或恢复映射引擎")
	systray.AddSeparator()
	quitItem := systray.AddMenuItem("退出", "退出桌面应用")

	t.refreshPauseLabel()

	go func() {
		for {
			select {
			case <-openItem.ClickedCh:
				t.app.ShowWindow()
			case <-reloadItem.ClickedCh:
				_, _ = t.app.RestartEngine()
			case <-t.pauseItem.ClickedCh:
				status, _ := t.app.engine.Status()
				if status.Paused || !status.Running {
					_, _ = t.app.ResumeEngine()
				} else {
					_, _ = t.app.PauseEngine()
				}
				t.refreshPauseLabel()
			case <-quitItem.ClickedCh:
				if t.app.ctx != nil {
					runtime.Quit(t.app.ctx)
				}
				return
			}
		}
	}()
}

func (t *trayController) refreshPauseLabel() {
	if t.pauseItem == nil {
		return
	}
	status, _ := t.app.engine.Status()
	if status.Paused || !status.Running {
		t.pauseItem.SetTitle("恢复映射")
		return
	}
	t.pauseItem.SetTitle("暂停映射")
}
