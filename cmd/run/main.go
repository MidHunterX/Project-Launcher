package main

import (
	"fmt"
	"os"

	"run/internal/config"
	"run/internal/project"
	"run/internal/tmux"
	"run/internal/util"
)

func main() {
	cfg := config.Load()
	project.InitServices(cfg.EnabledServices)

	// ====================== Auto Detection Section ====================== //

	if cfg.ProjectName == "" {
		pwd, _ := os.Getwd()
		cfg.ProjectName = project.DetectName(pwd)
	}

	if cfg.ProjectType == "" {
		cfg.ProjectType = project.DetectTech()
		util.Log(fmt.Sprintf("🔍 Detected: %s", cfg.ProjectType))
	}

	if cfg.URL == "" {
		cfg.URL = project.GetServerURL(cfg.ProjectType)
	}

	// ========================= Execution Section ========================= //

	if cfg.UseCustomEnv {
		project.ExecuteCustomEnv(cfg)
	} else {
		project.SetupEnv(cfg.ProjectType)
	}

	if cfg.URL != "" && cfg.URL != "none" {
		go util.OpenBrowser(cfg.URL, cfg.Browser)
	}

	// DEFENSE: Tmux Dependency or run directly (blocking)
	if !util.CommandExists("tmux") {
		cmdStr := project.GetServerCommand(cfg.ProjectType)
		if cmdStr == "" {
			util.Log("🚫 Error: No server found for this type")
			os.Exit(1)
		}
		util.Log("ℹ️ Running server without tmux...")
		util.RunCommand(cmdStr)
		os.Exit(0)
	}

	// CHECK: Session already exists
	if tmux.SessionExists(cfg.ProjectName) {
		util.Log(fmt.Sprintf("Session '%s' already exists. Attaching...", cfg.ProjectName))
		tmux.Attach(cfg.ProjectName)
		os.Exit(0)
	}

	// ACTION: TMUX window layout
	if cfg.UseCustomLayout {
		project.ExecuteCustomLayout(cfg)
	} else {
		tmux.SetupLayout(cfg)
	}

	// ACTION: Post-init hook
	if cfg.UsePostInitHook {
		project.ExecutePostInitHook(cfg)
	}

	tmux.Attach(cfg.ProjectName)
	util.Log("👋 Goodbye!")
}
