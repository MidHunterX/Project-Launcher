# Project Launcher

A universal, intelligent project runner that automatically detects your project type and sets up an optimized development environment with tmux sessions, proper virtual environments, and optional system services.

### Supported Technologies

<p>
  <img src="https://img.shields.io/badge/Rust-black?style=for-the-badge&logo=rust&logoColor=#E57324" />
  <img src="https://img.shields.io/badge/Elixir-4B275F?style=for-the-badge&logo=elixir&logoColor=white" />

  <img src="https://img.shields.io/badge/Node%20js-339933?style=for-the-badge&&logo=nodedotjs&logoColor=white" />
  <img src="https://img.shields.io/badge/next%20js-000000?style=for-the-badge&logo=nextdotjs&logoColor=white" />
  <img src="https://img.shields.io/badge/Angular-DD0031?style=for-the-badge&logo=angular&logoColor=white" />

  <img src="https://img.shields.io/badge/Python-FFD43B?style=for-the-badge&logo=python&logoColor=blue" />
  <img src="https://img.shields.io/badge/flask-%23000.svg?style=for-the-badge&logo=flask&logoColor=white" />
  <img src="https://img.shields.io/badge/fastapi-109989?style=for-the-badge&logo=FASTAPI&logoColor=white" />
  <img src="https://img.shields.io/badge/Django-092E20?style=for-the-badge&logo=django&logoColor=green" />
</p>

## 🔥 Features

- **Automatic Project Detection**: Intelligently detects project type based on files present
- **Multi-Technology Support**: Rust, Python, Django, FastAPI, Node.js, Next.js, Elixir, Angular
- **Environment Management**: Automatically sets up virtual environments and dependencies
- **Tmux Integration**: Creates organized tmux sessions with dedicated windows for development
- **System Service Management**: Start required services (PostgreSQL, Docker, MongoDB, etc.)
- **Custom Overrides**: Project-specific customizations with `.run_env`
- **Flexible Configuration**: Extensive configuration options for different workflows

## 📦 Installation

Clone the repository and run the install script:

```bash
git clone "https://github.com/MidHunterX/Project-Launcher" --depth 1
cd Project-Launcher
bash ./install.sh
```

### Requirements

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
4. Create a tmux session with appropriate windows
5. Launch your development server
6. Run post-initialization commands if enabled
7. Open server URL (if any) in browser

> [!NOTE]
> If TMUX is not installed, the runner will just detect project, run the
> appropriate server command and open server URL in browser.

## ⚙️ Configuration Settings

Create a `.run_env` file in your project root and add following overrides for
project-specific customization. All settings are optional.

### Per-Project Configuration

```bash
# ------------------------------------------------------------- PROJECT DETAILS

# Default: project dir name
PROJECT_NAME=""

# Default: automatic detection
# Available values:
# none | rust | python | django | fastapi | nodejs | nextjs | elixir | angular
PROJECT_TYPE=""

# ------------------------------------------------------------- SYSTEM SERVICES

# Services: add any valid systemd service names here to start
ENABLED_SERVICES=(
  # postgresql
  # docker
  # nginx
)

# ------------------------------------------------------------ BROWSER SETTINGS

# Default: automatic detection
# Available values: none | {custom url}
URL=""

# Default: xdg-open
# Available values: {custom command, arguments and newtab flag at the end}
BROWSER=""

# --------------------------------------------------------- BEHAVIORAL SETTINGS

# Settings: [ true | false ]
AUTORUN_COMMANDS=true
```

### Global Configuration

Create `run/config.conf` in your config directory.

Example 1: Just Browser

```bash
# custom browser, arguments and newtab flag at the end
BROWSER="firefox --new-tab"
```

Example 2: Window Manager + Custom Timezone + Custom Browser Profile

```bash
# custom browser, arguments and newtab flag at the end
BROWSER="hyprctl dispatch -- exec [workspace 2] TZ=Asia/Dubai firefox-developer-edition -P Personal -no-remote -new-tab"
```

The global configuration file sets default values for all projects.

### Per-Project Logic Overriding

Use the following overrides in `.run_env` for more granular customization.

These can also be used to extend support for new technologies as well.

Add `#!/usr/bin/env bash` at the top of the file to make it play nice with LSP.

```bash
# ------------------------------------------------------ OVERRIDE PROJECT SETUP

USE_CUSTOM_ENV=false
setup_env_custom() {
  # Project specific custom init setup example
  # Syntax: setup_env "env_directory" "command if env_directory not found"
  setup_env "node_modules/" "npm install --force --legacy-peer-deps"
}

# -------------------------------------------------------- OVERRIDE TMUX LAYOUT

USE_CUSTOM_LAYOUT=false
setup_layout_custom() {
  # Project specific custom layout example
  local current_dir=$(pwd)
  local api_dir="../example-fastapi-project"

  # 1. Server (FastAPI)
  cd "$api_dir"
  setup_python_env
  create_window "API Server" "fastapi dev main.py"

  # 2. Server (NextJS)
  cd "$current_dir"
  setup_node_env
  create_window "Web Server" "npm run dev"

  # 3. Editor (FastAPI)
  cd "$api_dir"
  create_window "Editor (API)" "nvim"
  deactivate

  # 4. Editor (NextJS)
  cd "$current_dir"
  create_window "Editor (Web)" "nvim"
}

# ------------------------------------------------------- CUSTOM POST INIT HOOK

USE_POST_INITIALIZATION_HOOK=false
setup_post_init_hook() {
  # Project specific post execution hook example
  # Custom browser launch logic
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

### Environment Variables

You can also add environment variables to `.run_env`

```bash
export RUST_LOG=debug
```

## 📜 License

This runner is distributed under the [MIT License](LICENSE).
