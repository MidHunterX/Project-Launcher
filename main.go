package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Configuration holds all runtime settings
type Config struct {
	ProjectName     string
	ProjectType     string
	EnabledServices []string
	URL             string
	Browser         string
	AutoRunCommands bool
	UseCustomLayout bool
	UseCustomEnv    bool
	UsePostInitHook bool
	OverrideFile    string
}

// Colors for terminal output
const (
	Red   = "\033[1;31m"
	Green = "\033[1;32m"
	Cyan  = "\033[1;36m"
	Reset = "\033[0;0m"
)

var config Config

func main() {
	// ====================== Initialization Section ====================== //

	config = Config{
		ProjectName:     "",
		ProjectType:     "",
		EnabledServices: []string{},
		URL:             "",
		Browser:         "",
		AutoRunCommands: true,
		UseCustomLayout: false,
		UseCustomEnv:    false,
		UsePostInitHook: false,
		OverrideFile:    ".run_env",
	}

	loadGlobalConfig()

	loadProjectOverrides()

	initServices()

	// ====================== Auto Detection Section ====================== //

	if config.ProjectName == "" {
		pwd, _ := os.Getwd()
		config.ProjectName = filepath.Base(pwd)
	}

	if config.ProjectType == "" {
		config.ProjectType = detectProjectType()
		log(fmt.Sprintf("🔍 Detected: %s", config.ProjectType))
	}

	if config.URL == "" {
		config.URL = getServerURL(config.ProjectType)
	}

	// ========================= Execution Section ========================= //

	if config.UseCustomEnv {
		executeCustomEnv()
	} else {
		setupEnv(config.ProjectType)
	}

	if config.URL != "" && config.URL != "none" {
		// go routine to open browser in background
		go openBrowser(config.URL, config.Browser)
	}

	// ----------------------------------------------- DEFENSE: Tmux Dependency
	if !commandExists("tmux") {
		// FALLBACK: run directly on current terminal (blocking)
		cmd := getServerCommand(config.ProjectType)
		if cmd == "" {
			log("🚫 Error: No server found for this type")
			os.Exit(1)
		}
		log("ℹ️ Running server without tmux...")
		runCommand(cmd)
		os.Exit(0)
	}

	// CHECK: Session already exists
	if tmuxSessionExists(config.ProjectName) {
		log(fmt.Sprintf("Session '%s' already exists. Attaching...", config.ProjectName))
		tmuxAttach(config.ProjectName)
		os.Exit(0)
	}

	// ACTION: TMUX window layout
	if config.UseCustomLayout {
		executeCustomLayout()
	} else {
		setupLayout(config.ProjectType)
	}

	// ACTION: Post-init hook
	if config.UsePostInitHook {
		executePostInitHook()
	}

	tmuxAttach(config.ProjectName)

	log("👋 Goodbye!")
}

// ========================= [ HELPER FUNCTIONS ] ========================= //

func loadGlobalConfig() {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		home, _ := os.UserHomeDir()
		configHome = filepath.Join(home, ".config")
	}
	configFile := filepath.Join(configHome, "run", "config.conf")

	if _, err := os.Stat(configFile); err == nil {
		loadConfigFile(configFile)
	}
}

func loadProjectOverrides() {
	if _, err := os.Stat(config.OverrideFile); err == nil {
		loadConfigFile(config.OverrideFile)
	}
}

func loadConfigFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse simple key=value pairs
		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.Trim(strings.TrimSpace(parts[1]), "\"")

				switch key {
				case "PROJECT_NAME":
					config.ProjectName = value
				case "PROJECT_TYPE":
					config.ProjectType = value
				case "URL":
					config.URL = value
				case "BROWSER":
					config.Browser = value
				case "AUTORUN_COMMANDS":
					config.AutoRunCommands = value == "true"
				case "USE_CUSTOM_LAYOUT":
					config.UseCustomLayout = value == "true"
				case "USE_CUSTOM_ENV":
					config.UseCustomEnv = value == "true"
				case "USE_POST_INITIALIZATION_HOOK":
					config.UsePostInitHook = value == "true"
				}
			}
		}

		// TODO: Parse EnabledServices array
	}
}

// ======================== [ SERVICE MANAGEMENT ] ======================== //

func initServices() {
	for _, service := range config.EnabledServices {
		if !isServiceActive(service) {
			log(fmt.Sprintf("🔧 Starting %s...", service))
			cmd := exec.Command("sudo", "systemctl", "start", service)
			cmd.Run()
		}
	}
}

func isServiceActive(service string) bool {
	cmd := exec.Command("systemctl", "is-active", "--quiet", service)
	return cmd.Run() == nil
}

// ========================= [ PROJECT HANDLING ] ========================= //

func detectProjectType() string {
	// CHECK: Python/FastAPI
	if fileExists("main.py") {
		content, _ := os.ReadFile("main.py")
		if strings.Contains(string(content), "from fastapi import FastAPI") {
			return "fastapi"
		}
		if strings.Contains(string(content), "from flask import Flask") {
			return "flask"
		}
		return "python"
	}

	// CHECK: Python Django
	if fileExists("manage.py") {
		return "django"
	}

	// CHECK: Node.js projects
	if fileExists("package.json") {
		if fileExists("angular.json") {
			return "angular"
		}
		if fileExists("next.config.ts") || fileExists("next.config.js") {
			return "nextjs"
		}
		return "nodejs"
	}

	// CHECK: Elixir
	if fileExists("mix.exs") {
		return "elixir"
	}

	// CHECK: Rust
	if fileExists("Cargo.toml") {
		return "rust"
	}

	return "none"
}

// --------------------------------------------------------- PROJECT SERVER URL

func getServerURL(projectType string) string {
	urls := map[string]string{
		"flask":   "http://localhost:5000",
		"django":  "http://localhost:8000",
		"nextjs":  "http://localhost:3000",
		"angular": "http://localhost:4200",
		"elixir":  "http://localhost:4000",
	}

	if url, ok := urls[projectType]; ok {
		return url
	}
	return ""
}

// --------------------------------------------------------- PROJECT SERVER CMD

func getServerCommand(projectType string) string {
	commands := map[string]string{
		"python":  "python main.py",
		"fastapi": "fastapi dev main.py",
		"flask":   "python main.py",
		"django":  "python manage.py runserver",
		"nodejs":  "npm start",
		"nextjs":  "npm run dev",
		"angular": "ng serve",
		"elixir":  "iex -S mix phx.server",
		"rust":    "cargo run",
	}

	if cmd, ok := commands[projectType]; ok {
		return cmd
	}
	return ""
}

// ---------------------------------------------------------- PROJECT SETUP ENV

func setupEnv(projectType string) {
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
			log(fmt.Sprintf("🚧 Unknown technology: %s", projectType))
		}
	}
}

func setupPythonEnv() {
	dependencyDir := "venv/"
	if !dirExists(dependencyDir) {
		log("🚀 Setting up dependencies...")
		if !commandExists("python") {
			log("❌ Error: python not found")
			os.Exit(1)
		}
		runCommand("python -m venv venv")

		log("📦 Installing requirements...")
		if fileExists("requirements.txt") {
			runCommand("venv/bin/pip install -r requirements.txt")
		} else {
			log("🚧 Warning: No requirements.txt found")
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
	if !dirExists(dependencyDir) {
		log("🚀 Setting up dependencies...")
		cmdParts := strings.Fields(dependencyCmd)
		if !commandExists(cmdParts[0]) {
			log(fmt.Sprintf("❌ Error: %s not found", cmdParts[0]))
			os.Exit(1)
		}
		runCommand(dependencyCmd)
	}
}

// --------------------------------------------------------- PROJECT INJECT ENV

func injectEnv(projectType string, cmd string) string {
	switch projectType {
	case "python", "flask", "django", "fastapi":
		return inject_venv(cmd)
	default:
		return cmd
	}
}

func inject_venv(cmd string) string {
	if dirExists("venv/") {
		return "source venv/bin/activate && " + cmd
	}
	return cmd
}

// ====================== [ TMUX SESSION MANAGEMENT ] ====================== //

func setupLayout(projectType string) {
	switch projectType {
	case "fastapi":
		setupFastAPILayout()
	case "django":
		setupDjangoLayout()
	case "rust":
		setupRustLayout()
	default:
		setupBaseLayout(projectType)
	}
}

func setupBaseLayout(projectType string) {
	createWindow("Cmd", getServerCommand(projectType))
	createWindow("Editor", "nvim")
}

func setupFastAPILayout() {
	setupBaseLayout(config.ProjectType)
	createWindow("Test", "clear && pytest")
}

func setupDjangoLayout() {
	setupBaseLayout(config.ProjectType)
	createWindow("Test", "python manage.py test")
}

func setupRustLayout() {
	setupBaseLayout(config.ProjectType)
	createWindow("Test", "cargo test")
}

// ------------------------------------------------------------- TMUX FUNCTIONS

func createWindow(name, command string) {
	// DEFENSE: Invalid Name
	if name == "" {
		log("🚫 Error: Missing window name")
		return
	}

	// DEFENSE: Window without command
	if command == "" {
		ensureWindow(name)
		return
	}

	// DEFENSE: Invalid Command
	cmdParts := strings.Fields(command)
	if !commandExists(cmdParts[0]) {
		return
	}

	ensureWindow(name)

	command = injectEnv(config.ProjectName, command)

	// "clear" fixes send-keys doubling text due to shell race condition
	tmuxSendKeys(config.ProjectName, name, "clear")
	tmuxSendKeysEnter(config.ProjectName, name)

	tmuxSendKeys(config.ProjectName, name, command)
	if config.AutoRunCommands {
		tmuxSendKeysEnter(config.ProjectName, name)
	}
}

func ensureWindow(name string) {
	if tmuxSessionExists(config.ProjectName) {
		cmd := exec.Command("tmux", "new-window", "-t", config.ProjectName, "-n", name)
		cmd.Run()
	} else {
		cmd := exec.Command("tmux", "new-session", "-d", "-s", config.ProjectName, "-n", name)
		cmd.Run()
	}
}

func tmuxSessionExists(session string) bool {
	cmd := exec.Command("tmux", "has-session", "-t", session)
	return cmd.Run() == nil
}

func tmuxSendKeys(session, window, keys string) {
	cmd := exec.Command("tmux", "send-keys", "-t", fmt.Sprintf("%s:%s", session, window), keys)
	cmd.Run()
}

func tmuxSendKeysEnter(session, window string) {
	cmd := exec.Command("tmux", "send-keys", "-t", fmt.Sprintf("%s:%s", session, window), "Enter")
	cmd.Run()
}

func tmuxAttach(session string) {
	cmd := exec.Command("tmux", "attach", "-t", session)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// ========================== [ EXTRA FEATURES ] ========================== //

func openBrowser(url, browser string) {
	if url == "" {
		return
	}

	if browser == "" {
		browser = "xdg-open"
	} else {
		cmdParts := strings.Fields(browser)
		if !commandExists(cmdParts[0]) {
			log(fmt.Sprintf("❌ Error: %s not found", cmdParts[0]))
			os.Exit(1)
		}
	}

	log("🔗 Launching browser...")

	// Wait for server to start
	timeout := 20
	for timeout > 0 {
		resp, err := http.Get(url)
		if err == nil {
			resp.Body.Close()
			break
		}
		time.Sleep(1 * time.Second)
		timeout--
	}

	if timeout <= 0 {
		return
	}

	// Launch browser
	var cmd *exec.Cmd
	if browser == "xdg-open" {
		cmd = exec.Command(browser, url)
	} else {
		parts := strings.Fields(browser)
		parts = append(parts, url)
		cmd = exec.Command(parts[0], parts[1:]...)
	}
	cmd.Start()
}

func executeCustomEnv() {
	if !fileExists(config.OverrideFile) {
		log(fmt.Sprintf("🔴 Error: %s%s%s not defined after %sUSE_CUSTOM_ENV=true%s",
			Red, "setup_env_custom()", Reset, Cyan, Reset))
		return
	}

	// Execute the .run_env file which should define setup_env_custom
	cmd := exec.Command("bash", "-c", fmt.Sprintf("source %s && type setup_env_custom >/dev/null 2>&1 && setup_env_custom", config.OverrideFile))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func executeCustomLayout() {
	if !fileExists(config.OverrideFile) {
		log(fmt.Sprintf("🔴 Error: %s%s%s not defined after %sUSE_CUSTOM_LAYOUT=true%s",
			Red, "setup_layout_custom()", Reset, Cyan, Reset))
		return
	}

	// Execute the .run_env file which should define setup_layout_custom
	cmd := exec.Command("bash", "-c", fmt.Sprintf("source %s && type setup_layout_custom >/dev/null 2>&1 && setup_layout_custom", config.OverrideFile))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("PROJECT_NAME=%s", config.ProjectName),
		fmt.Sprintf("PROJECT_TYPE=%s", config.ProjectType),
	)
	cmd.Run()
}

func executePostInitHook() {
	if !fileExists(config.OverrideFile) {
		log(fmt.Sprintf("🔴 Error: %s%s%s not defined after %sUSE_POST_INITIALIZATION_HOOK=true%s",
			Red, "setup_post_init_hook()", Reset, Cyan, Reset))
		return
	}

	// Execute the .run_env file which should define setup_post_init_hook
	cmd := exec.Command("bash", "-c", fmt.Sprintf("source %s && type setup_post_init_hook >/dev/null 2>&1 && setup_post_init_hook", config.OverrideFile))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// ========================= [ UTILITY FUNCTIONS ] ========================= //

func log(message string) {
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf("[%s%s%s] %s\n", Green, timestamp, Reset, message)
}

func commandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func dirExists(dirname string) bool {
	info, err := os.Stat(dirname)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func runCommand(command string) {
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Run()
}
