# **_Note: this is mock documentation_**

# rss 😎
rss is a multi-user feed aggregator CLI application written in Go. It is only intented for local use and mainly allows the user to:

- Add RSS feeds from across the internet to be collected
- Store the collected posts in a PostgreSQL database
- Follow and unfollow RSS feeds that other users have added
- View summaries of the aggregated posts in the terminal, with a link to the full post

The program uses and generates a single JSON file to keep track of two things:
1. Who is currently logged in
2. The connection credentials for the PostgreSQL database

```json
{
  "db_url": "connection_string_goes_here",
  "current_user_name": "username_goes_here"
}
```
The file is generated at `$HOME/.rssconfig.json`
There's no user-based authentication for this program. If someone has the database credentials, they can act as any user.

## Installation

### 1. Install Go
This super program requires Go to be installed. Here is an easy command for installing it:

```bash
curl -sS https://webi.sh/golang | sh
```

### 2. Install PostgreSQL
You may install it in Linux/WSL (Debian) with:

```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
```
Ensure the installation worked with:
```bash
postgres --version
```

### 3. Install rss
This command will download, build and install the `rss` command into your Go toolchain's `bin` directory:
he

```bash
go install github.com/joseflores1/rss@latest
```
\
**_This would work if I somehow could, through the binary, create and populate the needed database (is that even legal?)_**


Because this is not the case, the following must be done:

1. Create a PostgreSQL database named `rss`

2. Install goose
    ```bash
    go install github.com/pressly/goose/v3/cmd/goose@latest
    ```

3. Migrate database
    ```bash
    # From the root of the repository
    cd sql/schema/
    goose postgres <connection_string> up
    ```
    Where the connection string is of the type `protocol://username:password@host:port/database` (e.g. `postgres://postgres:postgres@localhost:5432/rss`)

\
These steps should create and populate the `rss` database needed for the program to work

## Usage
You just write commands in the following way:

```bash
rss <command_name> [<args>...]
```

## Available commands

- `addfeed <feed_name> <feed_url>` – Add a new RSS feed using the given URL and names it according to `<feed_name>` (requires login)
- `agg <time_between_requests>` – Scrap the added feeds for posts to save into the database, check [time.ParseDuration](https://pkg.go.dev/time#ParseDuration) for the accepted `<time_between_requests>` formats (requires login and an added feed)
- `browse [<limit>]` – Browses the posts of the feeds that the user follows, a `<limit>` argument can be passed to limit the number of posts shown in console (default is 2)
- `feeds` – List all added feeds
- `follow <feed_url>` – Follow a feed (requires login)
- `following` – Show the feeds you're following (requires login)
- `login <user_name>` – Log into your account
- `register <user_name>` – Create a new user account
- `reset` – Delete all the contents from the database, it is useful for development
- `unfollow <feed_url>` – Unfollow a feed (requires login)
- `users` – List registered users
