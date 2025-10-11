package tmux

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"run/internal/config"
	"run/internal/project"
	"run/internal/util"
)

func CreateWindow(cfg *config.Config, name, command string) {
	if name == "" {
		util.Log("🚫 Error: Missing window name for tmux")
		return
	}
	if command == "" {
		ensureWindow(cfg.Project.Name, name)
		return
	}

	cmdParts := strings.Fields(command)
	if !util.CommandExists(cmdParts[0]) {
		util.Log(fmt.Sprintf("🚫 Command '%s' not found, skipping window '%s'", cmdParts[0], name))
		return
	}

	ensureWindow(cfg.Project.Name, name)
	injectedCmd := project.InjectEnv(cfg.Project.Type, command)

	SendKeys(cfg.Project.Name, name, "clear")
	SendKeysEnter(cfg.Project.Name, name)
	SendKeys(cfg.Project.Name, name, injectedCmd)

	if cfg.Behavior.AutoRunCommands {
		SendKeysEnter(cfg.Project.Name, name)
	}
}

func ensureWindow(session, name string) {
	if SessionExists(session) {
		exec.Command("tmux", "new-window", "-t", session, "-n", name).Run()
	} else {
		exec.Command("tmux", "new-session", "-d", "-s", session, "-n", name).Run()
	}
}

func SessionExists(session string) bool {
	cmd := exec.Command("tmux", "has-session", "-t", session)
	return cmd.Run() == nil
}

func SendKeys(session, window, keys string) {
	target := fmt.Sprintf("%s:%s", session, window)
	exec.Command("tmux", "send-keys", "-t", target, keys).Run()
}

func SendKeysEnter(session, window string) {
	target := fmt.Sprintf("%s:%s", session, window)
	exec.Command("tmux", "send-keys", "-t", target, "Enter").Run()
}

func Attach(session string) {
	cmd := exec.Command("tmux", "attach", "-t", session)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
