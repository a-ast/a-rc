package launchd

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

const label = "com.a-rc.daemon"

//go:embed plist.tmpl
var plistTmpl string

// Adapter implements core.ProcessManager using macOS launchd.
type Adapter struct {
	plistPath string
}

func New() *Adapter {
	home, _ := os.UserHomeDir()
	return &Adapter{
		plistPath: filepath.Join(home, "Library", "LaunchAgents", label+".plist"),
	}
}

type plistData struct {
	Label      string
	BinaryPath string
	ConfigPath string
	LogDir     string
}

func (a *Adapter) Install(binaryPath, configPath, logDir string) error {
	// Ensure log directory exists.
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return fmt.Errorf("creating log dir: %w", err)
	}

	// If already running, send SIGHUP to reload config.
	if running, pid, err := a.Status(); err == nil && running && pid > 0 {
		proc, err := os.FindProcess(pid)
		if err == nil {
			_ = proc.Signal(os.Interrupt) // SIGHUP via kill -HUP
			// Use kill directly for SIGHUP since os.Interrupt is SIGINT on Unix.
			_ = exec.Command("kill", "-HUP", strconv.Itoa(pid)).Run()
			fmt.Println("daemon reloaded (SIGHUP sent)")
			return nil
		}
	}

	// Render plist.
	tmpl, err := template.New("plist").Parse(plistTmpl)
	if err != nil {
		return fmt.Errorf("parsing plist template: %w", err)
	}
	f, err := os.Create(a.plistPath)
	if err != nil {
		return fmt.Errorf("writing plist: %w", err)
	}
	if err := tmpl.Execute(f, plistData{
		Label:      label,
		BinaryPath: binaryPath,
		ConfigPath: configPath,
		LogDir:     logDir,
	}); err != nil {
		f.Close()
		return fmt.Errorf("rendering plist: %w", err)
	}
	f.Close()

	// Bootstrap the agent for the current user session.
	uid := os.Getuid()
	out, err := exec.Command(
		"launchctl", "bootstrap",
		fmt.Sprintf("gui/%d", uid),
		a.plistPath,
	).CombinedOutput()
	if err != nil {
		return fmt.Errorf("launchctl bootstrap: %w\n%s", err, out)
	}
	return nil
}

func (a *Adapter) Uninstall() error {
	uid := os.Getuid()
	out, err := exec.Command(
		"launchctl", "bootout",
		fmt.Sprintf("gui/%d", uid),
		a.plistPath,
	).CombinedOutput()
	if err != nil {
		return fmt.Errorf("launchctl bootout: %w\n%s", err, out)
	}
	if err := os.Remove(a.plistPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("removing plist: %w", err)
	}
	return nil
}

// Status queries launchctl for the daemon's running state and PID.
func (a *Adapter) Status() (running bool, pid int, err error) {
	out, err := exec.Command("launchctl", "list", label).Output()
	if err != nil {
		// Not found — not running.
		return false, 0, nil
	}
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, `"PID"`) {
			// Format: "PID" = 12345;
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				pidStr := strings.TrimSuffix(parts[2], ";")
				if p, e := strconv.Atoi(pidStr); e == nil && p > 0 {
					return true, p, nil
				}
			}
		}
	}
	return true, 0, nil // listed but not running (crashed/stopped)
}
