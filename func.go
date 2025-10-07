package main

import (
	"fmt"

	"Blog-aggregator/internal/config"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("username is required")
	}

	s.cfg.UserName = cmd.args[0]
	if err := config.SetUser(*s.cfg); err != nil {
		return err
	}

	fmt.Printf("username is now set to %v\n", cmd.args[0])
	return nil
}
