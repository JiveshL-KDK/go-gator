# Go Gator

Go Gator is a command-line RSS feed aggregator built in Go, backed by a PostgreSQL database. It's a learning project built while following the [Boot.dev](https://boot.dev) backend track.

> **Status:** v0.1 — user management (register/login/list/reset) is implemented. Feed following, fetching, and aggregation are coming next.

## Prerequisites

- [Go](https://go.dev/doc/install) 1.27+
- [PostgreSQL](https://www.postgresql.org/download/) running locally or remotely
- [goose](https://github.com/pressly/goose) (optional, for running the SQL migrations in `sql/schema`)

## Installation

Clone the repo and build the binary:

```bash
git clone https://github.com/JiveshL-KDK/go-gator.git
cd go-gator
go build -o go-gator
```

Or install it directly with `go install`:

```bash
go install github.com/JiveshL-KDK/go-gator@latest
```

## Configuration

Go Gator reads its configuration from a `.gatorconfig.json` file in your home directory. Create it before running any command:

```bash
echo '{"db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable", "current_user_name": ""}' > ~/.gatorconfig.json
```

| Field                | Description                                  |
| --------------------- | --------------------------------------------- |
| `db_url`              | PostgreSQL connection string                  |
| `current_user_name`   | Set automatically by the `login` command      |

Once the config file exists and points at a valid database, run the migrations in `sql/schema` (e.g. with `goose`) to create the `users` table before using the CLI.

## Usage

```bash
./go-gator <command> [arguments]
```

### Commands

| Command    | Arguments     | Description                          |
| ---------- | ------------- | ------------------------------------- |
| `register` | `<username>`  | Create a new user and log in as them |
| `login`    | `[username]`  | Log in as an existing user (with no argument, prints the currently logged-in user) |
| `list`     | —             | List all registered users             |
| `reset`    | —             | Delete all users from the database    |

**Examples:**

```bash
./go-gator register alice
./go-gator login alice
./go-gator list
./go-gator reset
```

## Project Structure

```
.
├── main.go                     # Entry point — parses CLI args and dispatches commands
├── internal/
│   ├── commands/                # Command definitions and handlers
│   ├── config/                  # Reads/writes ~/.gatorconfig.json
│   ├── database/                # sqlc-generated database access code
│   └── state/                   # Shared app state (config + DB connection)
├── sql/
│   ├── schema/                  # goose migrations
│   └── queries/                 # SQL queries used by sqlc
└── sqlc.yaml                    # sqlc configuration
```

## Tech Stack

- [Go](https://go.dev)
- [PostgreSQL](https://www.postgresql.org)
- [sqlc](https://sqlc.dev) — generates type-safe Go from SQL
- [goose](https://github.com/pressly/goose) — SQL schema migrations
- [lib/pq](https://github.com/lib/pq) — PostgreSQL driver
- [google/uuid](https://github.com/google/uuid) — UUID generation

## Roadmap

- [ ] Add and follow RSS feeds
- [ ] Aggregate posts from followed feeds on a schedule
- [ ] Browse aggregated posts from the CLI

## License

This project is for educational purposes as part of the Boot.dev curriculum.
