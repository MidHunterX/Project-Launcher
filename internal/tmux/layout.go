package tmux

import (
	"run/internal/config"
	"run/internal/project"
)

func SetupLayout(cfg *config.Config) {
	switch cfg.ProjectType {
	case "fastapi":
		setupFastAPILayout(cfg)
	case "django":
		setupDjangoLayout(cfg)
	case "rust":
		setupRustLayout(cfg)
	default:
		setupBaseLayout(cfg)
	}
}

func setupBaseLayout(cfg *config.Config) {
	serverCmd := project.GetServerCommand(cfg.ProjectType)
	CreateWindow(cfg, "Cmd", serverCmd)
	CreateWindow(cfg, "Editor", "nvim")
}

func setupFastAPILayout(cfg *config.Config) {
	setupBaseLayout(cfg)
	CreateWindow(cfg, "Test", "clear && pytest")
}

func setupDjangoLayout(cfg *config.Config) {
	setupBaseLayout(cfg)
	CreateWindow(cfg, "Test", "python manage.py test")
}

func setupRustLayout(cfg *config.Config) {
	setupBaseLayout(cfg)
	CreateWindow(cfg, "Test", "cargo test")
}
