# Project Launcher

Project Launcher automatically detects your project type and sets up a complete tmux development environment with pre-configured windows, services, and commands.

### Supported Technologies

<span>
  <img src="https://img.shields.io/badge/Rust-black?style=for-the-badge&logo=rust&logoColor=#E57324" />

  <img src="https://img.shields.io/badge/Elixir-4B275F?style=for-the-badge&logo=elixir&logoColor=white" />

  <img src="https://img.shields.io/badge/Python-FFD43B?style=for-the-badge&logo=python&logoColor=blue" />
  <img src="https://img.shields.io/badge/fastapi-109989?style=for-the-badge&logo=FASTAPI&logoColor=white" />
  <img src="https://img.shields.io/badge/Django-092E20?style=for-the-badge&logo=django&logoColor=green" />

  <img src="https://img.shields.io/badge/Node%20js-339933?style=for-the-badge&logo=nodedotjs&logoColor=white" />
  <img src="https://img.shields.io/badge/next%20js-000000?style=for-the-badge&logo=nextdotjs&logoColor=white" />
  <img src="https://img.shields.io/badge/Angular-DD0031?style=for-the-badge&logo=angular&logoColor=white" />
</span>

## 🔥 Features

- **Auto-detection**: Automatically detects project type (Python, FastAPI, Django, Rust, Node.js, Next.js etc.)
- **Quick Setup**: Installs dependencies and sets up virtual environments automatically
- **Tmux (optional) Integration**: Creates organized tmux sessions with dedicated windows for different tasks
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
- Create a tmux session (if available) with appropriate windows
- Launch your development server
- Run post-initialization commands if enabled (need tmux)

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

### Post-Initialization Hook Examples

#### Opening Browser

Enable a post-initialization hook to open browser with specific URL.

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

#### Integration with IDEs

The script works well with terminal-based editors (nvim, emacs) but can be adapted for GUI editors:

```bash
USE_POST_INITIALIZATION_HOOK=true
setup_post_init_hook() {
  # Launch VS Code in specific workspace
  code . &
}
```

## 📜 License

This script is distributed under the [MIT License](LICENSE).

Adapt and modify as needed for your workflow and needs.
