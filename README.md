# Go Gator

Go Gator (`gator`) is a command-line RSS feed aggregator built in Go, backed by a PostgreSQL database. It's a learning project built while following the [Boot.dev](https://boot.dev) backend track.

> **Status:** v0.2 — user management, feed following, and RSS aggregation are all implemented.

## Prerequisites

Before installing, make sure you have the following installed on your machine:

- [Go](https://go.dev/doc/install) 1.27+
- [PostgreSQL](https://www.postgresql.org/download/) running locally or remotely
- [goose](https://github.com/pressly/goose) (optional, for running the SQL migrations in `sql/schema`)

## Installation

Install the `gator` CLI directly with `go install`:

```bash
go install github.com/JiveshL-KDK/go-gator@latest
```

This compiles a standalone `gator` binary and places it in your `$GOPATH/bin` (or `$HOME/go/bin` by default) — make sure that directory is on your `$PATH`. Since Go programs compile to a single static binary, once installed you can run `gator` directly without needing the Go toolchain around.

Alternatively, clone the repo and build it yourself:

```bash
git clone https://github.com/JiveshL-KDK/go-gator.git
cd go-gator
go build -o gator
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
gator <command> [arguments]
```

(If you built locally instead of using `go install`, run `./gator <command> [arguments]` from the repo directory instead.)

### Commands

| Command      | Arguments                     | Description                                                                        |
| ------------ | ------------------------------ | ----------------------------------------------------------------------------------- |
| `register`   | `<username>`                   | Create a new user and log in as them                                                |
| `login`      | `[username]`                   | Log in as an existing user (with no argument, prints the currently logged-in user)  |
| `list`       | —                               | List all registered users                                                           |
| `reset`      | —                               | Delete all users from the database                                                  |
| `add`        | `<name> <url>`                 | Add a new RSS feed and automatically follow it                                      |
| `list-feeds` | —                               | List all feeds added by any user                                                    |
| `follow`     | `<url>`                        | Follow an existing feed by its URL                                                  |
| `following`  | —                               | Show the feeds the current user is following                                        |
| `unfollow`   | `<url>`                        | Unfollow a feed by its URL                                                          |
| `agg`        | —                               | Continuously fetch the next feed due for a refresh and store new posts             |
| `browse`     | `[limit]`                      | Show recent posts from feeds the current user follows (defaults to a set limit)     |

A few commands (`add`, `list-feeds`, `follow`, `following`, `unfollow`, `browse`) require you to be logged in first with `login`.

**Examples:**

```bash
gator register alice
gator login alice
gator add "Boot.dev Blog" https://blog.boot.dev/index.xml
gator following
gator agg
gator browse 5
```

## Project Structure

```
.
├── main.go                     # Entry point — parses CLI args and dispatches commands
├── internal/
│   ├── commands/                # Command definitions and handlers
│   ├── config/                  # Reads/writes ~/.gatorconfig.json
│   ├── constants/                # Shared constant values
│   ├── database/                # sqlc-generated database access code
│   ├── rss/                     # RSS feed fetching and parsing
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

- [x] Add and follow RSS feeds
- [x] Aggregate posts from followed feeds
- [x] Browse aggregated posts from the CLI
- [ ] Run aggregation on a recurring schedule

## License

This project is for educational purposes as part of the Boot.dev curriculum.
