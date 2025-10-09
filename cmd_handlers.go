package main

import (
	"context"
	"fmt"
	"html"
	"time"

	"Blog-aggregator/internal/config"
	"Blog-aggregator/internal/database"

	"github.com/google/uuid"
)

func handlerFollowing(s *state, cmd command) error {
	user, err := s.db.GetUser(context.Background(), s.cfg.UserName)
	if err != nil {
		return err
	}
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	if len(feeds) < 1 {
		return fmt.Errorf("\"%v\" not following any feed\n", user.Name)
	}

	fmt.Printf("\"%v\" is following:\n", user.Name)
	for _, feed := range feeds {
		fmt.Printf("  - %v\n", feed.FeedName)
	}
	return nil
}

func handlerFollow(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("url is required")
	}

	url := cmd.args[0]
	ctx := context.Background()

	user, err := s.db.GetUser(ctx, s.cfg.UserName)
	if err != nil {
		return err
	}
	feed, err := s.db.GetFeedByURL(ctx, url)
	if err != nil {
		return err
	}

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	_, err = s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Printf("follower: %v\n", user.Name)
	fmt.Printf("feed: %v\n", feed.Name)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for i, feed := range feeds {
		user, err := s.db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return err
		}

		fmt.Printf("Feed # %v:\n", i+1)
		fmt.Printf("	Name:    %v\n", feed.Name)
		fmt.Printf("	URL:     %v\n", feed.Url)
		fmt.Printf("	Creater: %v\n", user.Name)
	}
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("\"%v\" requires two arguments\n", cmd.name)
	}

	currentUserName := s.cfg.UserName
	user, err := s.db.GetUser(context.Background(), currentUserName)
	if err != nil {
		return err
	}

	_, err = s.db.GetFeedByURL(context.Background(), cmd.args[1])
	if err == nil {
		return fmt.Errorf("url already exist\nUse \"follow\" cmd instead")
	}

	paramsF := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	}

	_, err = s.db.CreateFeed(context.Background(), paramsF)
	if err != nil {
		return err
	}

	paramsFF := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    paramsF.ID,
	}

	s.db.CreateFeedFollow(context.Background(), paramsFF)

	fmt.Printf("successfully added to \"%v\": %v\n", currentUserName, cmd.args[1])
	return nil
}

func handlerAgg(s *state, cmd command) error {
	//if len(cmd.args) < 1 {
	//	return fmt.Errorf("URL is required")
	//}
	url := "https://www.wagslane.dev/index.xml"

	rss, err := fetchFeed(context.Background(), url)
	if err != nil {
		return err
	}

	rss.Channel.Description = html.UnescapeString(rss.Channel.Description)
	rss.Channel.Title = html.UnescapeString(rss.Channel.Title)
	for _, item := range rss.Channel.Item { // Pass by val??
		title, desc := item.Title, item.Description
		item.Description = html.UnescapeString(desc)
		item.Title = html.UnescapeString(title)
	}
	fmt.Printf("%+v\n", rss)
	return nil
}

func handlerUsers(s *state, cmd command) error {
	currentUser := s.cfg.UserName
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	if len(users) < 1 {
		return fmt.Errorf("No users available\n")
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
