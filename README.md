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
- **Custom Overrides**: Project-specific customization via `run_lib.sh`
- **Flexible Configuration**: Extensive configuration options for different workflows

## 📦 Installation

Make the runner globally accessible:

```bash
# Copy to a directory in your PATH
sudo cp run /usr/local/bin/
sudo chmod +x /usr/local/bin/run
```

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

## ⚙️ Basic Configuration Overrides

Create a `run_lib.sh` file in your project root and add following overrides for
project-specific customization:

### TMUX Session Name

By default, the runner uses the project directory name as the tmux session name.

```bash
PROJECT_NAME="your-project-name"
```

### Project Type

By default, the runner tries to detect the project type based on files present
in the directory. If you need to override this detection, you can specify the
project type to use:

```bash
# Available values:
# none | rust | python | django | fastapi | nodejs | nextjs | elixir | angular
PROJECT_TYPE="type-to-use"
```

### Service Management (Systemd)

You can specify valid systemd service names to start using the `ENABLED_SERVICES` array.

```bash
ENABLED_SERVICES=(
  postgresql
  docker
  nginx
)
```

### Behavior settings

```bash
AUTOSTART_SERVER=true
AUTORUN_COMMANDS=true
```

## ⚙️ Advanced Configuration Overrides

### Custom Environment Setup

```bash
USE_CUSTOM_ENV=true
setup_env_custom() {
  # Custom dependency installation
  setup_env "node_modules/" "npm install --legacy-peer-deps"

  # Additional setup steps
  npm run build:deps
}
```

### Post-Initialization Hook

Example 1: Launch Browser

```bash
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

Example 2: GUI IDE Integration

```bash
USE_POST_INITIALIZATION_HOOK=true
setup_post_init_hook() {
  log "🔗 Launching browser & VSCode..."
  (
    # Launch browser automatically
    sleep 3
    firefox --new-tab http://localhost:3000 &

    # Open IDE
    code . &
  ) &
}
```

## Advanced Features

### Environment Detection Logic

The runner uses intelligent detection of project type to optimize setup.

### Session Management

- **Existing Sessions**: Automatically attaches to existing session
- **Clean Shutdown**: Proper cleanup and deactivation
- **Multiple Projects**: Each project gets its own isolated session

### Error Handling

- **Missing Dependencies**: Graceful fallback without tmux
- **Command Not Found**: Skips unavailable commands
- **Service Failures**: Continues execution with warnings

## Troubleshooting

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

4. Session already exists:

   ```bash
   # Kill existing session
   tmux kill-session -t project-name
   run
   ```

## Contributing

The runner is designed to be extensible. To add support for new technologies:

1. Add detection logic to `detect_project_type()`
2. Add environment setup in `setup_env()`
3. Add server command in `run_server_command()`
4. Optionally add custom layout in `setup_layout()`

## Dependencies

- **Required**: `bash`, `systemctl` (for services)
- **Optional**: `tmux` (for session management)
- **Per-project**: Technology-specific tools (npm, pip, cargo, etc.)

## 📜 License

This runner is distributed under the [MIT License](LICENSE).
