package main

import (
	"fmt"

	"Blog-aggregator/internal/config"
	"Blog-aggregator/internal/database"
)

var cmds = commands{
	handler: make(map[string]func(*state, command) error),
}

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

type state struct {
	db  *database.Queries
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
