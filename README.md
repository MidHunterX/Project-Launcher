# Project Launcher

A development environment orchestrator that automatically detects your project type and sets up appropriate tmux sessions, environment variables, optional system services and other script based automations.

### Supported Technologies

<p>
  <img src="https://img.shields.io/badge/Rust-black?style=for-the-badge&logo=rust&logoColor=#E57324" />
  <img src="https://img.shields.io/badge/Elixir-4B275F?style=for-the-badge&logo=elixir&logoColor=white" />
  <img src="https://img.shields.io/badge/Flutter-02569B?style=for-the-badge&logo=flutter&logoColor=white" />

  <img src="https://img.shields.io/badge/Node%20js-339933?style=for-the-badge&&logo=nodedotjs&logoColor=white" />
  <img src="https://img.shields.io/badge/next%20js-000000?style=for-the-badge&logo=nextdotjs&logoColor=white" />
  <img src="https://img.shields.io/badge/Angular-DD0031?style=for-the-badge&logo=angular&logoColor=white" />

  <img src="https://img.shields.io/badge/Python-FFD43B?style=for-the-badge&logo=python&logoColor=blue" />
  <img src="https://img.shields.io/badge/HTML5-E34F26?style=for-the-badge&logo=html5&logoColor=white" />
  <img src="https://img.shields.io/badge/flask-%23000.svg?style=for-the-badge&logo=flask&logoColor=white" />
  <img src="https://img.shields.io/badge/fastapi-109989?style=for-the-badge&logo=FASTAPI&logoColor=white" />
  <img src="https://img.shields.io/badge/Django-092E20?style=for-the-badge&logo=django&logoColor=green" />

  <img src="https://img.shields.io/badge/Astal-5190cf?style=for-the-badge&logo=astral&logoColor=white" />
</p>

### What to expect

<table>
  <tr>
    <td>Without run</td>
    <td>With run</td>
  </tr>
  <tr>
    <td>

```bash
> cd django-project
# Now.. what was that venv command?
> python -m venv venv
> source venv/bin/activate
# Phew almost done
> python manage.py runserver
# Error.. dependency not installed
> pip install -r requirements.txt
# Finally done
> python manage.py runserver
# Migration error dammit
> python manage.py migrate
# Open a new terminal
> nvim
# Now need to see output
> python manage.py runserver
# Open up browser and go to localhost:8000
# Done! Start working on the project.

# Finished development. Now for the next
> cd fastapi-project
> python -m venv venv
> source venv/bin/activate
> pip install -r requirements.txt
# What was the run command again?
> fastapi dev main.py
# Open up a new terminal
> cd nextjs-project
> npm install
> npm run dev
# Open up browser and go to localhost:3000
# Done! Wait.. features are not working
# Go to fastapi-project
# Debugging.. ahh found  new model changes
> ^C
# Look into README for migration command
> alembic upgrade head
# Ohh.. this project needs database server
> systemctl start postgresql
> fastapi dev main.py
# Finally time to develop new features.
> nvim
# Now what was I supposed to build?
```

</td>
  <td>

```bash
> cd django-project
> run
# Start coding right away!

# Add a onetime mono-repo config
> cd nextjs-project
> run
# Start coding right away again!
```

  </td>
</tr>
</table>

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
3. Migrate database if any unapplied migrations are detected
4. Start required services
5. Create a tmux session with appropriate windows
6. Launch your development server
7. Run post-initialization commands if enabled
8. Open server URL (if any) in browser
9. Does a post validation and warns if models are left changed without creating migrations on exit

> [!NOTE]
> If TMUX is not installed, the runner will just detect project, run the
> appropriate server command and open server URL in browser.

## ⚙️ Configuration Settings

Create a `.runrc` file in your project root and add following overrides for
project-specific customization. All settings are optional.

### Per-Project Configuration

#### Minimal Configuration

The most minimal configuration is actually NO configuration at all. But, if you
are like me who don't want PostgreSQL running in the background all the time
and only wants it running when developing a specific project only then, your
minimal config for that project would just be:

```bash
ENABLED_SERVICES=(
  postgresql
)
```

#### Full Configuration

```bash
# ------------------------------------------------------------- PROJECT DETAILS

# Default: project dir name
PROJECT_NAME=""

# Default: automatic detection
# Available values:
# none | rust | python | django | fastapi | flask | html | nodejs | nextjs | elixir | angular | flutter | astal-gtk
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

Example 2: Custom Timezone + Browser

```bash
BROWSER="TZ=Asia/Dubai firefox --new-tab"
```

Example 3: Window Manager + Custom Timezone + Custom Browser Profile

```bash
BROWSER="hyprctl dispatch -- exec [workspace 2] TZ=Asia/Dubai firefox-developer-edition -P Personal -no-remote -new-tab"
```

The global configuration file sets default values for all projects.

### Per-Project Logic Overriding

Use the following overrides in `.runrc` for more granular customization.

These can also be used to extend support for new technologies as well.

Add `#!/usr/bin/env bash` at the top of the file to make it play nice with LSP.

```bash
# ------------------------------------------------------ OVERRIDE PROJECT SETUP

USE_CUSTOM_ENV=false
setup_env_custom() {
  # Project specific custom init setup example
  # Syntax: setup_base_env "env_directory" "command to init env_directory"
  setup_base_env "node_modules/" "npm install --force --legacy-peer-deps"
}

# -------------------------------------------------------- OVERRIDE TMUX LAYOUT

USE_CUSTOM_LAYOUT=false
setup_layout_custom() {
  # Project specific custom layout example
  local current_dir=$(pwd)
  local api_dir="../example-fastapi-project"

  # 1. Server (FastAPI)
  cd "$api_dir"
  setup_env "python"
  create_window "API Server" "fastapi dev main.py"

  # 2. Server (NextJS)
  cd "$current_dir"
  setup_env "nextjs"
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

You can also add environment variables to `.runrc`

```bash
export RUST_LOG=debug
```

### Example Scenarios

A python project's requirements only works with an older python version.

```sh
USE_CUSTOM_ENV=true
setup_env_custom() {
  setup_python_env "venv" "python3.11 -m venv venv"
}
```

## 📜 License

This runner is distributed under the [MIT License](LICENSE).
