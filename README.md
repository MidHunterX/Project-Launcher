# Project Launcher

Project Launcher automatically detects your project type and sets up a complete tmux development environment with pre-configured windows, services, and commands.

## 🔥 Features

- **Auto-detection**: Automatically detects project type (Python, FastAPI, Django, Rust, Node.js, Next.js etc.)
- **Quick Setup**: Installs dependencies and sets up virtual environments automatically
- **Tmux Integration**: Creates organized tmux sessions with dedicated windows for different tasks
- **Service Management**: Automatically starts required system services (PostgreSQL, Docker, MongoDB, Redis, Nginx etc.)
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
- Run post-initialization commands if enabled

## ⚙️ Configuration

### Basic Settings

Edit the configuration section at the top of the script:

```bash
# Project name (defaults to using directory name)
PROJECT_NAME=""

# Force a specific project type (defaults to using auto-detection)
PROJECT_TYPE=""

# Enable system services
ENABLED_SERVICES=(
  postgresql
)

# Behavior settings
AUTOSTART_SERVER=true    # Automatically start development servers
AUTORUN_COMMANDS=true    # Automatically execute commands in tmux
```

### Supported Project Types

- `none` - Generic project (default if not detected)
- `python` - Basic Python project with venv
- `django` - Django Full-Stack Framework
- `fastapi` - FastAPI Framework
- `nodejs` - Node.js Backend with NPM
- `nextjs` - Next.js Framework
- `angular` - Angular Framework
- `elixir` - Elixir Framework with Mix
- `rust` - Rust Project with Cargo

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

Override the default layout for complex multi-project setups (or just use Docker at this point):

```bash
USE_CUSTOM_LAYOUT=true
setup_custom_layout() {
  # Example: FastAPI + Next.js setup

  # 1. Start FastAPI server
  cd "../fastapi-project/"
  setup_python_env
  create_tmux_session "API Server" "fastapi dev main.py"

  # 2. Start Next.js server
  cd "$script_dir"
  setup_node_env
  create_window "Web Server" "npm run dev"

  # 3. Create editor windows
  cd "../fastapi-project/"
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
USE_POST_INITIALIZATION_HOOK=true
setup_post_init_hook() {
  # Launch VS Code in specific workspace
  code . &
}
```

## 🤝 Contributing

This script is designed to be easily customizable for your specific workflow. Common improvements:

1. Add support for new project types in `detect_project_type()`
2. Create new layout functions for specific frameworks
3. Add new service integrations in `init_services()`
4. Enhance environment setup for complex dependency chains

## 📜 License

This script is distributed under the [MIT License](LICENSE).

Adapt and modify as needed for your workflow and needs.
