package main

import (
	"fmt"

	"Blog-aggregator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	handler map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	fun, ok := c.handler[cmd.name]
	if !ok {
		return fmt.Errorf("\"%v\" cmd not found", cmd.name)
	}

	if err := fun(s, cmd); err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handler[name] = f
}
