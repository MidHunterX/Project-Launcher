# Project Launcher

Project Launcher automatically detects your project type and sets up a complete tmux development environment with pre-configured windows, services, and commands.

## 🔥 Features

- **Auto-detection**: Automatically detects project type (Python, FastAPI, Django, Node.js, Next.js, Angular, Rust, Elixir)
- **Quick Setup**: Installs dependencies and sets up virtual environments automatically
- **Tmux Integration**: Creates organized tmux sessions with dedicated windows for different tasks
- **Service Management**: Automatically starts required system services (PostgreSQL, Docker, MongoDB, Redis, Nginx)
- **Highly Customizable**: Override layouts, environments, and post-initialization hooks
- **Session Management**: Reattaches to existing sessions or creates new ones

## 🚀 Quick Start

1. Make the script executable:

```bash
chmod +x run
```

2. Run in your project directory:

```bash
./run
```

The script will:

- Detect your project type automatically
- Install dependencies if needed
- Start required services
- Create a tmux session with appropriate windows
- Launch your development server

## ℹ️ Supported Project Types

| Project Type | Detection Method                  | Default Layout       |
| ------------ | --------------------------------- | -------------------- |
| **FastAPI**  | `main.py` with FastAPI imports    | Server, Editor, Test |
| **Python**   | `main.py` file                    | Server, Editor       |
| **Django**   | `manage.py` file                  | Server, Editor, Test |
| **Next.js**  | `next.config.js` + `package.json` | Server, Editor       |
| **Node.js**  | `package.json` file               | Server, Editor, Test |
| **Angular**  | `angular.json` file               | Server, Editor       |
| **Rust**     | `Cargo.toml` file                 | Run, Editor, Test    |
| **Elixir**   | `mix.exs` file                    | Server, Editor       |
| **Generic**  | No specific files found           | Command, Editor      |

## ⚙️ Configuration

### Basic Settings

Edit the configuration section at the top of the script:

```bash
# Project name (defaults to directory name)
PROJECT_NAME=""

# Force a specific project type
PROJECT_TYPE=""

# Enable system services
ENABLED_SERVICES=(
  postgresql
  docker
  mongod
  redis-server
  nginx
)

# Behavior settings
AUTOSTART_SERVER=true    # Automatically start development servers
AUTORUN_COMMANDS=true    # Automatically run commands in tmux windows
```

### Available Services

Uncomment services in `ENABLED_SERVICES` to auto-start them:

- `postgresql` - PostgreSQL database
- `docker` - Docker daemon
- `mongod` - MongoDB
- `redis-server` - Redis server
- `nginx` - Nginx web server

### Custom Environment Setup

Enable custom environment setup for complex dependency management:

```bash
USE_CUSTOM_ENV=true
setup_custom_env() {
  # Example: Install with specific npm flags
  setup_env "node_modules/" "npm install --force --legacy-peer-deps"
}
```

### Custom Layouts

Override the default layout for complex multi-project setups:

```bash
USE_CUSTOM_LAYOUT=true
setup_custom_layout() {
  # Example: FastAPI + Next.js setup

  # 1. Start FastAPI server
  cd "../api-project/"
  setup_python_env
  source "./venv/bin/activate"
  create_tmux_session "API Server" "fastapi dev main.py"

  # 2. Start Next.js server
  cd "$script_dir"
  setup_node_env
  create_window "Web Server" "npm run dev"

  # 3. Create editor windows
  cd "../api-project/"
  create_window "Editor (API)" "nvim"
  cd "$script_dir"
  create_window "Editor (Web)" "nvim"
}
```

### Post-Initialization Hook

Add custom actions after setup completion:

```bash
USE_POST_INITIALIZATION_HOOK=true
setup_post_init_hook() {
  # Example: Open browser automatically
  log "🔗 Launching browser..."
  (
    sleep 2  # Wait for server to start
    firefox-developer-edition -new-tab "http://localhost:3000" &
  ) &
}
```

## 📚 Commands Reference

### Tmux Session Functions

- `create_tmux_session(name, command)` - Create new tmux session
- `create_window(name, command)` - Add window to existing session
- `create_temp_window(name, command)` - Add window that runs command once

### Environment Setup Functions

- `setup_python_env()` - Create Python venv and install requirements
- `setup_node_env()` - Run npm install
- `setup_elixir_env()` - Run mix setup
- `setup_env(dir, command)` - Generic dependency setup

### Utility Functions

- `log(message)` - Timestamped logging
- `detect_project_type()` - Auto-detect project type
- `init_services()` - Start enabled system services

## 📝 Examples

### Python FastAPI Project

```bash
# Detected automatically from main.py with FastAPI imports
# Creates: Server (fastapi dev), Editor (nvim), Test (pytest)
./run
```

### Next.js Project

```bash
# Detected from next.config.js + package.json
# Creates: Server (npm run dev), Editor (nvim)
./run
```

### Custom Multi-Service Setup

```bash
# In configuration section:
PROJECT_TYPE="custom"
USE_CUSTOM_LAYOUT=true
ENABLED_SERVICES=(postgresql redis-server)

# Custom layout handles multiple related projects
```

## 🐞 Troubleshooting

### Common Issues

**Script exits with "command not found":**

- Ensure required tools are installed (`npm`, `python`, `pip`, etc.)
- The script checks for command availability before creating windows

**Services fail to start:**

- Check if you have sudo permissions for systemctl
- Verify services are installed: `systemctl list-unit-files | grep service-name`

**Session already exists:**

- The script automatically attaches to existing sessions
- Kill existing session: `tmux kill-session -t project-name`

### Dependencies

Required tools (auto-checked by script):

- `tmux` - Session management
- `curl` - Health checks (for post-init hooks)
- Project-specific tools (`python`, `npm`, `cargo`, etc.)

Install tmux:

```bash
# Ubuntu/Debian
sudo apt install tmux

# macOS
brew install tmux

# Arch Linux
sudo pacman -S tmux
```

## 🎯 Advanced Usage

### Multiple Project Environments

For complex setups with multiple related projects:

1. Place the script in a parent directory
2. Use `USE_CUSTOM_LAYOUT=true`
3. Navigate between project directories in your custom layout
4. Each project can have its own window with appropriate environments

### Integration with IDEs

The script works well with terminal-based editors (nvim, emacs) but can be adapted for GUI editors:

```bash
setup_post_init_hook() {
  # Launch VS Code in specific workspace
  code ./my-workspace.code-workspace &

  # Or launch IntelliJ IDEA
  idea . &
}
```

### CI/CD Integration

Use the script in development containers or CI environments by setting:

```bash
AUTOSTART_SERVER=false  # Don't auto-start in CI
AUTORUN_COMMANDS=false  # Manual command execution
```

## Contributing

This script is designed to be easily customizable for your specific workflow. Common improvements:

1. Add support for new project types in `detect_project_type()`
2. Create new layout functions for specific frameworks
3. Add new service integrations in `init_services()`
4. Enhance environment setup for complex dependency chains

## License

This script is distributed under the [MIT License](LICENSE).
Adapt and modify as needed for your workflow and needs.
