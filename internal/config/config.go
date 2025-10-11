package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	GlobalConfigFileName = "config.yaml"
	LocalConfigFileName  = "run.yaml"
)

func Load() *Config {
	cfg := &Config{}
	cfg.loadDefaults()
	cfg.loadGlobalConfig()
	cfg.loadLocalConfig()
	return cfg
}

func (c *Config) loadDefaults() {
	c.Project = ProjectConfig{Name: "", Type: ""}
	c.System = SystemConfig{EnabledServices: []string{}}
	c.Browser = BrowserConfig{URL: "", Command: "xdg-open"}
	c.Behavior = BehaviorConfig{
		AutoRunCommands: true,
		UseCustomLayout: false,
		UseCustomEnv:    false,
		UsePostInitHook: false,
	}
	c.EnvSetup = nil
	c.Layout = nil
	c.ScriptHook = ""
}

type Config struct {
	Project    ProjectConfig  `yaml:"project"`
	System     SystemConfig   `yaml:"system"`
	Browser    BrowserConfig  `yaml:"browser"`
	Behavior   BehaviorConfig `yaml:"behavior"`
	EnvSetup   []EnvSetupRule `yaml:"env_setup"`
	Layout     []LayoutWindow `yaml:"layout"`
	ScriptHook string         `yaml:"script_hook"`
}

type ProjectConfig struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

type SystemConfig struct {
	EnabledServices []string `yaml:"enabled_services"`
}

type BrowserConfig struct {
	URL     string `yaml:"url"`
	Command string `yaml:"command"`
}

type BehaviorConfig struct {
	AutoRunCommands bool `yaml:"autorun_commands"`
	UseCustomLayout bool `yaml:"use_custom_layout"`
	UseCustomEnv    bool `yaml:"use_custom_env"`
	UsePostInitHook bool `yaml:"use_post_init_hook"`
}

type EnvSetupRule struct {
	Check   string `yaml:"check"`
	Command string `yaml:"command"`
}

type LayoutWindow struct {
	Name    string `yaml:"name"`
	Path    string `yaml:"path"`
	Env     string `yaml:"env,omitempty"`
	Command string `yaml:"command"`
}

func (c *Config) loadGlobalConfig() {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		home, _ := os.UserHomeDir()
		configHome = filepath.Join(home, ".config")
	}
	configFile := filepath.Join(configHome, "run", GlobalConfigFileName)
	c.loadYAMLFile(configFile)
}

func (c *Config) loadLocalConfig() {
	c.loadYAMLFile(LocalConfigFileName)
}

func (c *Config) loadYAMLFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return // file not found or unreadable — silently ignore
	}
	_ = yaml.Unmarshal(data, c) // override existing fields
}
