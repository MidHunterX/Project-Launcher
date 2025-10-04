package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

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

func Load() *Config {
	cfg := &Config{
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

	cfg.loadGlobalConfig()
	cfg.loadProjectOverrides()

	return cfg
}

func (c *Config) loadGlobalConfig() {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		home, _ := os.UserHomeDir()
		configHome = filepath.Join(home, ".config")
	}
	configFile := filepath.Join(configHome, "run", "config.conf")

	if _, err := os.Stat(configFile); err == nil {
		c.loadConfigFile(configFile)
	}
}

func (c *Config) loadProjectOverrides() {
	if _, err := os.Stat(c.OverrideFile); err == nil {
		c.loadConfigFile(c.OverrideFile)
	}
}

// TODO: Make config use YAML instead
func (c *Config) loadConfigFile(filename string) {
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

		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.Trim(strings.TrimSpace(parts[1]), "\"")

				switch key {
				case "PROJECT_NAME":
					c.ProjectName = value
				case "PROJECT_TYPE":
					c.ProjectType = value
				case "URL":
					c.URL = value
				case "BROWSER":
					c.Browser = value
				case "AUTORUN_COMMANDS":
					c.AutoRunCommands = value == "true"
				case "USE_CUSTOM_LAYOUT":
					c.UseCustomLayout = value == "true"
				case "USE_CUSTOM_ENV":
					c.UseCustomEnv = value == "true"
				case "USE_POST_INITIALIZATION_HOOK":
					c.UsePostInitHook = value == "true"
				}
			}
		}
		// TODO: Parse EnabledServices array or YAML
	}
}
