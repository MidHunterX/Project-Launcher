# Project Launcher

A universal, intelligent project launcher that automatically detects your project type and sets up an optimized development environment with tmux sessions, proper virtual environments, and optional system services.

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

- **Automatic Project Detection**: Intelligently detects project type based on files present
- **Multi-Technology Support**: Rust, Python, Django, FastAPI, Node.js, Next.js, Elixir, Angular
- **Environment Management**: Automatically sets up virtual environments and dependencies
- **Tmux Integration**: Creates organized tmux sessions with dedicated windows for development
- **System Service Management**: Start required services (PostgreSQL, Docker, MongoDB, etc.)
- **Custom Overrides**: Project-specific customizations
- **Flexible Configuration**: Extensive configuration options for different workflows

### Session Management

- **Existing Sessions**: Automatically attaches to existing session
- **Clean Shutdown**: Proper cleanup and deactivation
- **Multiple Projects**: Each project gets its own isolated session

### Error Handling

- **Missing Dependencies**: Graceful fallback without tmux
- **Command Not Found**: Skips unavailable commands
- **Service Failures**: Continues execution with warnings

## 📦 Installation

Make the runner globally accessible:

```bash
# Copy to a directory in your PATH
sudo cp run /usr/local/bin/
sudo chmod +x /usr/local/bin/run
```

### Dependencies

- **Required**: `bash`, `systemctl` (for services)
- **Optional**: `tmux` (for session management)
- **Per-project**: Technology-specific tools (npm, pip, cargo, etc.)

## 🚀 Quick Start

Navigate to any project directory and run:

```bash
run
```

The runner will:

1. Detect your project type automatically
2. Install dependencies if needed
3. Start required services
4. Create a tmux session (if available) with appropriate windows
5. Launch your development server
6. Run post-initialization commands if enabled (need tmux)

## ⚙️ Full Configuration

Create a `.run_env` file in your project root and add following overrides for
project-specific customization. All overrides are optional.

```bash
#!/usr/bin/env bash

# -----------------------------[ CONFIGURATION ]----------------------------- #

# Default: project dir name
PROJECT_NAME=""
# Default: automatic detection
# Available values:
# none | rust | python | django | fastapi | nodejs | nextjs | elixir | angular
PROJECT_TYPE=""

# Services: add any valid systemd service names here to start
ENABLED_SERVICES=(
  # postgresql
  # docker
  # nginx
)

# Settings: [ true | false ]
AUTOSTART_SERVER=true
AUTORUN_COMMANDS=true

# ---------------------------[ CUSTOM OVERRIDES ]--------------------------- #

current_dir="$(pwd)"

USE_CUSTOM_ENV=false
setup_env_custom() {
  # Project specific custom init setup example
  setup_env "node_modules/" "npm install --force --legacy-peer-deps"
}

USE_CUSTOM_LAYOUT=false
setup_layout_custom() {
  # Project specific custom layout example

  # 1. Server (FastAPI)
  cd "../backend/"
  setup_env "fastapi"
  create_tmux_session "API Server" "fastapi dev main.py"

  # 2. Server (NextJS)
  cd "$current_dir"
  setup_env "nextjs"
  create_window "Web Server" "npm run dev"

  # 3. Editor (FastAPI)
  cd "../backend/"
  create_window "Editor (API)" "nvim"
  deactivate

  # 4. Editor (NextJS)
  cd "$current_dir"
  create_window "Editor (Web)" "nvim"
}

USE_POST_INITIALIZATION_HOOK=false
setup_post_init_hook() {
  # Project specific post execution hook example
  log "🔗 Launching browser..."
  (
    new_tab_url="http://localhost:3000"

    # Wait for the server to start
    local timeout=20
    while ! curl -sf "$new_tab_url" >/dev/null; do
      sleep 1
      ((timeout--)) || { exit 1; }
    done

    # Launch the browser
    firefox-developer-edition -P Personal -no-remote -new-tab $new_tab_url &
  ) &
}
```

## 🔧 Troubleshooting

### Common Issues

1. Runner won't start services:

   ```bash
   # Check if user has sudo privileges for systemctl
   sudo systemctl start postgresql
   ```

2. Tmux not found:

- Runner falls back to running server directly
- Install tmux: `sudo apt install tmux` or `brew install tmux`

3. Virtual environment issues:

   ```bash
   # Clean up and retry
   rm -rf venv/
   run
   ```

## 📜 License

This runner is distributed under the [MIT License](LICENSE).
