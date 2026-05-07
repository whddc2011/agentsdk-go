//go:build windows

package toolbuiltin

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

func bashOutputBaseDir() string {
	return filepath.Join(os.TempDir(), "agentsdk", "bash-output")
}

func newBashExecCmd(ctx context.Context, command string) *exec.Cmd {
	shell, args := resolveWindowsShell(command)
	cmd := exec.CommandContext(ctx, shell, args...)
	// 这是 Windows 下**彻底静默无窗口**的终极组合
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000 | 0x20000000 | 0x00000008,
	}
	return cmd
}

func resolveWindowsShell(command string) (exe string, argv []string) {
	if path, err := exec.LookPath("bash"); err == nil {
		return path, []string{"-c", command}
	}
	comspec := os.Getenv("ComSpec")
	if comspec == "" {
		if path, err := exec.LookPath("cmd.exe"); err == nil {
			comspec = path
		} else {
			comspec = "cmd.exe"
		}
	}
	return comspec, []string{"/d", "/s", "/c", command}
}
