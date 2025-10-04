package project

import (
	"fmt"
	"os"
	"os/exec"

	"run/internal/config"
	"run/internal/util"
)

func ExecuteCustomEnv(cfg *config.Config) {
	executeHook(cfg, "setup_env_custom")
}

func ExecuteCustomLayout(cfg *config.Config) {
	cmd := executeHook(cfg, "setup_layout_custom")
	if cmd != nil {
		cmd.Env = append(os.Environ(),
			fmt.Sprintf("PROJECT_NAME=%s", cfg.ProjectName),
			fmt.Sprintf("PROJECT_TYPE=%s", cfg.ProjectType),
		)
		cmd.Run()
	}
}

func ExecutePostInitHook(cfg *config.Config) {
	executeHook(cfg, "setup_post_init_hook")
}

// TODO: handle inline hooks by exec in project path
func executeHook(cfg *config.Config, funcName string) *exec.Cmd {
	if !util.FileExists(cfg.OverrideFile) {
		util.Log(fmt.Sprintf("🔴 Error: %s%s()%s not defined because override file is missing", util.Red, funcName, util.Reset))
		return nil
	}

	// The script checks if the function exists before calling it.
	script := fmt.Sprintf("source %s && type %s >/dev/null 2>&1 && %s", cfg.OverrideFile, funcName, funcName)
	cmd := exec.Command("bash", "-c", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}
