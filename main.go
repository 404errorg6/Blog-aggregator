package main

import (
	"database/sql"
	"fmt"
	"os"

	"Blog-aggregator/internal/config"
	"Blog-aggregator/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("Error occured while reading config: %v\n", err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.DB_URL)
	dbQueries := database.New(db)

	State := state{
		db:  dbQueries,
		cfg: &cfg,
	}
	// Register commands
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)

	// Args login
	args := os.Args
	if len(args) < 2 {
		fmt.Printf("not enought arguments\n")
		os.Exit(1)
	}

	userCmd := command{
		name: args[1],
		args: args[2:],
	}

	if err := cmds.run(&State, userCmd); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
