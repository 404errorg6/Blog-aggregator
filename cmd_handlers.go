package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"Blog-aggregator/internal/config"
	"Blog-aggregator/internal/database"

	"github.com/google/uuid"
)

func handlerBrowser(s *state, cmd command, user database.User) error {
	var limit int
	if len(cmd.args) < 1 {
		fmt.Printf("limit argument not provided, defaulting to 2 posts\n")
		limit = 2
	} else {
		tmp, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return err
		}
		if tmp < 1 {
			return fmt.Errorf("limit argument cannot be less than 1\n")
		}
		limit = tmp
	}

	params := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	}

	posts, err := s.db.GetPostsForUser(context.Background(), params)
	if err != nil {
		return err
	}
	if len(posts) < 1 {
		return fmt.Errorf("No posts found\n")
	}

	for i, post := range posts {
		fmt.Printf("Post # %v:\n\n", i+1)
		fmt.Printf("Title:        %v\n", post.Title)
		fmt.Printf("Description:  %v\n", post.Description)
		fmt.Printf("URL:          %v\n", post.Url)
		fmt.Printf("PublishedAt:  %v\n", post.PublishedAt)
		fmt.Printf("-----------------------------\n\n")
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("url is required")
	}

	url := cmd.args[0]

	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return err
	}

	params := database.DelFeedFollowEntryParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	if err := s.db.DelFeedFollowEntry(context.Background(), params); err != nil {
		return err
	}

	fmt.Printf("successfully unfollowed \"%v\"\n", feed.Name)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
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

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("url is required")
	}

	url := cmd.args[0]
	ctx := context.Background()

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

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("\"%v\" requires two arguments\n", cmd.name)
	}

	_, err := s.db.GetFeedByURL(context.Background(), cmd.args[1])
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

	if _, err := s.db.CreateFeedFollow(context.Background(), paramsFF); err != nil {
		return err
	}

	fmt.Printf("successfully added to \"%v\": %v\n", user.Name, cmd.args[1])
	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("time_between_reqs argument is required")
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Collecting feeds every %v\n", timeBetweenReqs)
	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		if err := scrapeFeeds(s); err != nil {
			return err
		}
	}
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
	fmt.Printf("successfully blanked out the database\n")
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
