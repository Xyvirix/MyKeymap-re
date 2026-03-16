package desktop

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type Paths struct {
	RootDir      string
	BinDir       string
	DataDir      string
	SiteDir      string
	TemplatesDir string
	ConfigFile   string
	MyKeymapExe  string
	TrayIcon     string
}

func DiscoverPaths() (Paths, error) {
	candidates := []string{}

	if wd, err := os.Getwd(); err == nil {
		candidates = append(candidates, wd)
	}
	if exe, err := os.Executable(); err == nil {
		dir := filepath.Dir(exe)
		candidates = append(candidates, dir, filepath.Dir(dir), filepath.Dir(filepath.Dir(dir)))
	}

	for _, start := range candidates {
		if root, ok := findRoot(start); ok {
			return buildPaths(root), nil
		}
	}

	return Paths{}, fmt.Errorf("cannot locate MyKeymap root directory")
}

func (p Paths) FrontendFS() (fs.FS, error) {
	for _, dir := range []string{
		p.SiteDir,
		filepath.Join(p.RootDir, "config-ui", "dist"),
	} {
		if ok, err := hasFile(dir, "index.html"); err == nil && ok {
			return os.DirFS(dir), nil
		}
	}

	return nil, fmt.Errorf("frontend assets not found in %s or config-ui/dist", p.SiteDir)
}

func buildPaths(root string) Paths {
	templatesDir := filepath.Join(root, "config-server", "templates")
	if ok, _ := hasFile(filepath.Join(root, "bin", "templates"), "MyKeymap.tmpl"); ok {
		templatesDir = filepath.Join(root, "bin", "templates")
	}

	return Paths{
		RootDir:      root,
		BinDir:       filepath.Join(root, "bin"),
		DataDir:      filepath.Join(root, "data"),
		SiteDir:      filepath.Join(root, "bin", "site"),
		TemplatesDir: templatesDir,
		ConfigFile:   filepath.Join(root, "data", "config.json"),
		MyKeymapExe:  filepath.Join(root, "MyKeymap.exe"),
		TrayIcon:     filepath.Join(root, "bin", "icons", "logo.ico"),
	}
}

func findRoot(start string) (string, bool) {
	current := filepath.Clean(start)
	for {
		if isRootDir(current) {
			return current, true
		}

		parent := filepath.Dir(current)
		if parent == current {
			return "", false
		}
		current = parent
	}
}

func isRootDir(dir string) bool {
	for _, path := range []string{
		filepath.Join(dir, "data", "config.json"),
		filepath.Join(dir, "MyKeymap.exe"),
		filepath.Join(dir, "bin"),
	} {
		if _, err := os.Stat(path); err != nil {
			return false
		}
	}
	return true
}

func hasFile(dir, file string) (bool, error) {
	info, err := os.Stat(filepath.Join(dir, file))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return !info.IsDir(), nil
}
