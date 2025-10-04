package project

import (
	"os"
	"path/filepath"
	"strings"

	"run/internal/util"
)

func DetectName(pwd string) string {
	return filepath.Base(pwd)
}

func DetectTech() string {
	// CHECK: Python/FastAPI/Flask
	if util.FileExists("main.py") {
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
	if util.FileExists("manage.py") {
		return "django"
	}

	// CHECK: Node.js projects
	if util.FileExists("package.json") {
		if util.FileExists("angular.json") {
			return "angular"
		}
		if util.FileExists("next.config.ts") || util.FileExists("next.config.js") {
			return "nextjs"
		}
		return "nodejs"
	}

	// CHECK: Elixir
	if util.FileExists("mix.exs") {
		return "elixir"
	}

	// CHECK: Rust
	if util.FileExists("Cargo.toml") {
		return "rust"
	}

	return "none"
}

func GetServerURL(projectType string) string {
	urls := map[string]string{
		"flask":   "http://localhost:5000",
		"django":  "http://localhost:8000",
		"nextjs":  "http://localhost:3000",
		"angular": "http://localhost:4200",
		"elixir":  "http://localhost:4000",
	}
	return urls[projectType] // Returns "" if not found
}

func GetServerCommand(projectType string) string {
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
	return commands[projectType] // Returns "" if not found
}
