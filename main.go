package main

import (
	"fmt"
	"os"

	"Blog-aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("Error occured while reading config: %v\n", err)
		return
	}

	State := state{
		cfg: &cfg,
	}
	cmds := commands{
		handler: map[string]func(*state, command) error{
			"login": handlerLogin,
		},
	}

	cmds.register("login", cmds.handler["login"])

	// Args login
	args := os.Args
	if len(args) < 2 {
		fmt.Printf("not enought arguments")
		os.Exit(1)
	}

	userCmd := command{
		name: args[1],
		args: args[2:],
	}

	if err := cmds.run(&State, userCmd); err != nil {
		fmt.Printf("Error while running cmd: %v\n", err)
		return
	}

	fmt.Printf("Config: %+v\n", cfg)
}
