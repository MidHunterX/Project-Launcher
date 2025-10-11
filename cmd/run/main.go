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
	project.InitServices(cfg.System.EnabledServices)

	// ====================== Auto Detection Section ====================== //

	if cfg.Project.Name == "" {
		pwd, _ := os.Getwd()
		cfg.Project.Name = project.DetectName(pwd)
	}

	if cfg.Project.Type == "" {
		cfg.Project.Type = project.DetectTech()
		util.Log(fmt.Sprintf("🔍 Detected: %s", cfg.Project.Type))
	}

	if cfg.Browser.URL == "" {
		cfg.Browser.URL = project.GetServerURL(cfg.Browser.URL)
	}

	// ========================= Execution Section ========================= //

	if cfg.Behavior.UseCustomEnv {
		project.ExecuteCustomEnv(cfg)
	} else {
		project.SetupEnv(cfg.Project.Type)
	}

	if cfg.Browser.URL != "" && cfg.Browser.URL != "none" {
		go util.OpenBrowser(cfg.Browser.URL, cfg.Browser.Command)
	}

	// DEFENSE: Tmux Dependency or run directly (blocking)
	if !util.CommandExists("tmux") {
		cmdStr := project.GetServerCommand(cfg.Project.Type)
		if cmdStr == "" {
			util.Log("🚫 Error: No server found for this type")
			os.Exit(1)
		}
		util.Log("ℹ️ Running server without tmux...")
		util.RunCommand(cmdStr)
		os.Exit(0)
	}

	// CHECK: Session already exists
	if tmux.SessionExists(cfg.Project.Name) {
		util.Log(fmt.Sprintf("Session '%s' already exists. Attaching...", cfg.Project.Name))
		tmux.Attach(cfg.Project.Name)
		os.Exit(0)
	}

	// ACTION: TMUX window layout
	if cfg.Behavior.UseCustomLayout {
		project.ExecuteCustomLayout(cfg)
	} else {
		tmux.SetupLayout(cfg)
	}

	// ACTION: Post-init hook
	if cfg.Behavior.UsePostInitHook {
		project.ExecutePostInitHook(cfg)
	}

	tmux.Attach(cfg.Project.Name)
	util.Log("👋 Goodbye!")
}
