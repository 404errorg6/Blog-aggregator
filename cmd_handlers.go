package main

import (
	"context"
	"fmt"
	"time"

	"Blog-aggregator/internal/config"
	"Blog-aggregator/internal/database"

	"github.com/google/uuid"
)

func handlerUsers(s *state, cmd command) error {
	currentUser := s.cfg.UserName
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, userName := range users {
		if currentUser != userName {
			fmt.Printf("* %v \n", userName)
		} else {
			fmt.Printf("* %v (current)\n", userName)
		}
	}
	return nil
}

func handlerReset(s *state, cmd command) error {
	if err := s.db.VoidUser(context.Background()); err != nil {
		return err
	}
	fmt.Printf("successfully blanked out the table\n")
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("username is required\n")
	}

	name := cmd.args[0]
	if _, err := s.db.GetUser(context.Background(), name); err != nil {
		return fmt.Errorf("\"%v\" user doesn't exist\n", name)
	}

	s.cfg.UserName = name
	if err := config.SetUser(*s.cfg); err != nil {
		return err
	}

	fmt.Printf("username is now set to %v\n", cmd.args[0])
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("username is required\n")
	}

	Context := context.Background()
	name := cmd.args[0]
	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}

	if _, err := s.db.GetUser(Context, name); err == nil {
		return fmt.Errorf("\"%v\" already exists\n", name)
	}

	_, err := s.db.CreateUser(Context, params)
	if err != nil {
		return err
	}
	fmt.Printf("user \"%v\" successfully created\n", name)
	s.cfg.UserName = name
	config.SetUser(*s.cfg)

	return nil
}
