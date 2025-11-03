# mo

A lightweight CLI tool to streamline common development tasks. Built to scratch my own itch after typing the same commands hundreds of times.

## What it does

`mo` handles the repetitive stuff - creating databases, switching environment configs, syncing files between local and remote servers. Nothing fancy, just saves time.

## Installation

```bash
git clone https://github.com/lpheller/mogo.git
cd mogo
go build -o mo
sudo mv mo /usr/local/bin/
```

Or with Go installed:

```bash
go install github.com/lpheller/mogo@latest
```

## Usage

### Database commands

```bash
mo db create mydb          # Create a new MySQL database
mo db list                 # List all databases
mo db open                 # Open database client
```

Aliases: `db:create`, `createdb`, etc. - whatever feels natural.

### Laravel commands

```bash
mo l:clear                 # Clear all caches (or: mo lc)
mo l:fresh                 # Migrate fresh with seed (or: mo lf)
mo l:fresh --no-seed       # Without seeding

# Or with subcommands:
mo l clear
mo l fresh
```

### Environment management

Quick switches for common Laravel .env configurations:

```bash
mo env sqlite              # Switch to SQLite
mo env mailtrap            # Configure Mailtrap for emails
mo env maildev             # Use local Maildev
mo env sync                # Sync .env with .env.example
```

### Config shortcuts

Open config files without remembering where they are:

```bash
mo config nvim             # Opens your nvim config
mo config git              # Opens .gitconfig
mo config -e vim myapp     # Use vim instead of default editor
```

First time? Run `mo config:edit` to set up your config paths.

### Project setup

```bash
mo setup                   # Auto-detect and setup project (composer, npm, migrations, etc.)
```

Detects Laravel, Node.js projects and runs the appropriate setup steps.

### Remote sync

Sync databases and storage folders between local and remote servers (Laravel Projects): 

```bash
mo pull --database         # Pull database from remote server
mo pull --storage          # Pull storage folder
mo push --database         # Push database to remote
mo push --storage          # Push storage folder
```

Requires SSH config in your projects `.env` file. When no `PULL_SSH_HOST`, `PULL_SSH_HOST`, and `PULL_PROJECT_DIR` are set, `mo` will prompt you to enter them. These values are then stored in the local `.env` for future use. Same goes for `mo push`. 

## Configuration

Config lives in `~/.config/mortimer/config.json`:

```json
{
  "db_user": "root",
  "db_password": "",
  "db_host": "127.0.0.1",
  "db_port": "3306",
  "editor": "vscode",
  "config_paths": {
    "nvim": "/Users/you/.config/nvim/init.vim",
    "git": "/Users/you/.gitconfig"
  }
}
```

Edit with `mo config:edit` or add your own shortcuts.

## Why "mo"?

Short for Mortimer/ Morty. Needed a CLI sidekick that's short to type and doesn't clash with existing commands. Plus, typing `mo` hundreds of times a day just feels right.

## Development

```bash
go test ./...              # Run tests
go build                   # Build binary
```

## License

MIT - do whatever you want with it.
