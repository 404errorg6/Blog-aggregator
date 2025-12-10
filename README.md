# Blog-aggregator

## Overview
The **Blog-aggregator** is a Go-based application designed to manage and aggregate blog feeds. It provides functionality for user management, feed subscriptions, and post aggregation. The project is structured to ensure scalability and maintainability, leveraging SQL for database interactions and Go for backend logic.

---

## Features
- **User Management**: Create and manage user accounts.
- **Feed Subscriptions**: Follow and unfollow blog feeds.
- **Post Aggregation**: Aggregate posts from subscribed feeds.
- **Database Integration**: Uses PostgreSQL for data storage.
- **Modular Design**: Organized into distinct modules for configuration, database, and SQL queries.

---

## Project Structure
The project is organized as follows:

```
cmd_handlers.go       # Command handlers for CLI or API
consts.go             # Constants used across the application
func.go               # Utility functions
go.mod                # Go module dependencies
main.go               # Application entry point
README.md             # Project documentation
sqlc.yaml             # SQLC configuration file

internal/
  config/             # Configuration-related files
    consts.go
    funcs.go
  database/           # Database-related files
    db.go
    feed_follows.sql.go
    feeds.sql.go
    models.go
    posts.sql.go
    users.sql.go

sql/
  queries/            # SQL query files
    feed_follows.sql
    feeds.sql
    posts.sql
    users.sql
  schema/             # Database schema files
    001_users.sql
    002_feeds.sql
    003_feed_follows.sql
    004_feeds.sql
    005_posts.sql
```

---

## Prerequisites
- **Go**: Version 1.25.1 or higher.
- **PostgreSQL**: Ensure a running PostgreSQL instance.
- **Dependencies**:
  - `github.com/google/uuid`
  - `github.com/lib/pq`

---

## Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/404errorg6/Blog-aggregator.git
   cd Blog-aggregator
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Set up the database:
   - Create a PostgreSQL database.
   - Apply the schema files in `sql/schema/` in order.

4. Configure the application:
   - Update the configuration files in `internal/config/` as needed.

---

## Usage
1. Run the application:
   ```bash
   go run main.go
   ```

2. Interact with the application via the CLI.

---

## Development
### Generating SQL Code
This project uses [SQLC](https://sqlc.dev/) for generating Go code from SQL queries. To regenerate the database code, run it from root of this repo:
```bash
sqlc generate
```

---
