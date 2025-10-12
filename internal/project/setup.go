package project

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"run/internal/util"
)

func InitServices(services []string) {
	for _, service := range services {
		if !isServiceActive(service) {
			util.Log(fmt.Sprintf("🔧 Starting %s...", service))
			cmd := exec.Command("sudo", "systemctl", "start", service)
			cmd.Run() // Errors are ignored for simplicity
		}
	}
}

func isServiceActive(service string) bool {
	cmd := exec.Command("systemctl", "is-active", "--quiet", service)
	return cmd.Run() == nil
}

func SetupEnv(projectType string) {
	switch projectType {
	case "python", "flask", "django", "fastapi":
		setupPythonEnv()
	case "nodejs", "nextjs", "angular":
		setupNodeEnv()
	case "elixir":
		setupElixirEnv()
	case "rust":
		// Rust sets up its own environment with cargo
	default:
		if projectType != "none" {
			util.Log(fmt.Sprintf("🚧 Unknown technology for env setup: %s", projectType))
		}
	}
}

func setupPythonEnv() {
	if !util.DirExists("venv/") {
		util.Log("🚀 Setting up python venv...")
		if !util.CommandExists("python") {
			util.Log("❌ Error: python not found")
			os.Exit(1)
		}
		util.RunCommand("python -m venv venv")

		util.Log("📦 Installing requirements...")
		if util.FileExists("requirements.txt") {
			util.RunCommand("venv/bin/pip install -r requirements.txt")
		} else {
			util.Log("🚧 Warning: No requirements.txt found")
		}
	}
}

func setupNodeEnv() {
	setupBaseEnv("node_modules/", "npm install")
}

func setupElixirEnv() {
	setupBaseEnv("deps/", "mix setup")
}

func setupBaseEnv(dependencyDir, dependencyCmd string) {
	if !util.DirExists(dependencyDir) {
		util.Log("🚀 Setting up dependencies...")
		cmdParts := strings.Fields(dependencyCmd)
		if !util.CommandExists(cmdParts[0]) {
			util.Log(fmt.Sprintf("❌ Error: %s not found", cmdParts[0]))
			os.Exit(1)
		}
		util.RunCommand(dependencyCmd)
	}
}

func InjectEnv(projectType string, cmd string) string {
	switch projectType {
	case "python", "flask", "django", "fastapi":
		return injectVenv(cmd)
	default:
		return cmd
	}
}

func injectVenv(cmd string) string {
	if util.DirExists("venv/") {
		return fmt.Sprintf(`bash -c "source venv/bin/activate && %s"`, cmd)
	}
	return cmd
}
