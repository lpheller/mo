# Personal Helper CLI Tool

A personal helper CLI tool with a few commands that support me in my daily development workflow.

## Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/lpheller/mogo.git
   ```
2. Navigate to the project directory:
   ```sh
   cd mogo
   ```
3. Build the project:
   ```sh
   go build -o mo main.go
   ```

## Usage

### Available Commands

#### Database Management

- `db create` (Alias: `db c`): Create a new database
- `db list` (Alias: `db l`): List all databases
- `db open` (Alias: `db o`): Open the database in the default editor
- `db:open` (Alias: `opendb`): Open the database in the default editor
- `db:list` (Alias: `listdb`): List all databases
- `db:create` (Alias: `createdb`): Create a new database

#### Environment Management

- `env sqlite` (Alias: `env:sqlite`): Set the DB_CONNECTION to sqlite
- `env mailtrap` (Alias: `env:mailtrap`): Set the mail driver to mailtrap
- `env maildev` (Alias: `env:maildev`): Set the mail driver to mail-dev
- `env sync` (Alias: `sync:env`): Sync the .env file with .env.example

#### Configuration

- `config:edit` (Alias: `edit:config`): Edit the Mortimer config file
- `config` (Aliases: `cfg`, `qc`): Quickly open configuration files
  - Options:
    - `--editor` (Alias: `-e`): Specify a custom editor

#### Project Setup

- `setup` (Alias: `s`): Setup a project by running appropriate commands
