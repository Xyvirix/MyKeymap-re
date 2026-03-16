package desktop

import (
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"settings/internal/script"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

type EngineStatus struct {
	Running   bool   `json:"running"`
	Paused    bool   `json:"paused"`
	Managed   bool   `json:"managed"`
	PID       int    `json:"pid,omitempty"`
	LastError string `json:"lastError,omitempty"`
	UpdatedAt string `json:"updatedAt"`
}

type EngineController struct {
	mu        sync.Mutex
	paths     Paths
	process   *exec.Cmd
	paused    bool
	lastError string
}

func NewEngineController(paths Paths) *EngineController {
	return &EngineController{paths: paths}
}

func (e *EngineController) Status() (EngineStatus, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	running, pid, err := lookupProcess("MyKeymap.exe")
	if err != nil {
		e.lastError = err.Error()
	}

	managed := false
	if e.process != nil && e.process.Process != nil {
		managed = processAlive(e.process.Process.Pid)
	}

	return EngineStatus{
		Running:   running,
		Paused:    e.paused,
		Managed:   managed,
		PID:       pid,
		LastError: e.lastError,
		UpdatedAt: time.Now().Format(time.RFC3339),
	}, err
}

func (e *EngineController) Start() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	running, _, err := lookupProcess("MyKeymap.exe")
	if err == nil && running {
		e.paused = false
		e.lastError = ""
		return nil
	}

	if err := e.prepareRuntimeAssets(); err != nil {
		e.lastError = err.Error()
		return err
	}

	cmd := hiddenCommand(e.paths.RootDir, e.paths.MyKeymapExe, "/script", "./bin/MyKeymap.ahk")
	if err := cmd.Start(); err != nil {
		e.lastError = err.Error()
		return err
	}

	e.process = cmd
	e.paused = false
	e.lastError = ""

	go func(cmd *exec.Cmd) {
		_ = cmd.Wait()
	}(cmd)

	return nil
}

func (e *EngineController) Stop() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := killImage("MyKeymap-CommandInput.exe"); err != nil {
		e.lastError = err.Error()
	}
	if err := killImage("MyKeymap.exe"); err != nil {
		e.lastError = err.Error()
		return err
	}

	e.process = nil
	e.lastError = ""
	return nil
}

func (e *EngineController) Restart() error {
	if err := e.Stop(); err != nil {
		return err
	}
	time.Sleep(350 * time.Millisecond)
	return e.Start()
}

func (e *EngineController) Pause() error {
	if err := e.Stop(); err != nil {
		return err
	}

	e.mu.Lock()
	e.paused = true
	e.mu.Unlock()
	return nil
}

func (e *EngineController) Resume() error {
	e.mu.Lock()
	e.paused = false
	e.mu.Unlock()
	return e.Start()
}

func lookupProcess(image string) (bool, int, error) {
	cmd := hiddenCommand("", "tasklist", "/FO", "CSV", "/NH", "/FI", "IMAGENAME eq "+image)
	out, err := cmd.Output()
	if err != nil {
		return false, 0, err
	}

	text := strings.TrimSpace(string(out))
	if text == "" || !strings.Contains(strings.ToLower(text), strings.ToLower(image)) {
		return false, 0, nil
	}

	reader := csv.NewReader(strings.NewReader(text))
	record, err := reader.Read()
	if err != nil {
		return true, 0, nil
	}
	if len(record) < 2 {
		return true, 0, nil
	}

	pid, _ := strconv.Atoi(record[1])
	return true, pid, nil
}

func killImage(image string) error {
	cmd := hiddenCommand("", "taskkill", "/IM", image, "/F")
	out, err := cmd.CombinedOutput()
	if err == nil {
		return nil
	}

	text := strings.ToLower(string(out))
	if strings.Contains(text, "not found") || strings.Contains(text, "没有运行的任务") {
		return nil
	}

	return fmt.Errorf("taskkill %s failed: %s", image, strings.TrimSpace(string(out)))
}

func processAlive(pid int) bool {
	return pid > 0
}

func hiddenCommand(dir, name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd
}

func (e *EngineController) prepareRuntimeAssets() error {
	config, err := script.ParseConfig(e.paths.ConfigFile)
	if err != nil {
		return err
	}

	if err := script.GenerateDesktopScriptsAt(e.paths.TemplatesDir, e.paths.BinDir, config); err != nil {
		return err
	}

	shortcuts, err := filepath.Glob(filepath.Join(e.paths.RootDir, "shortcuts", "*.lnk"))
	if err != nil {
		return err
	}
	if len(shortcuts) == 0 {
		if _, err := os.Stat(filepath.Join(e.paths.RootDir, "shortcuts")); os.IsNotExist(err) {
			if mkErr := os.MkdirAll(filepath.Join(e.paths.RootDir, "shortcuts"), 0755); mkErr != nil {
				return mkErr
			}
		}

		cmd := hiddenCommand(e.paths.RootDir, e.paths.MyKeymapExe, "/script", "./bin/MiscTools.ahk", "GenerateShortcuts")
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}
