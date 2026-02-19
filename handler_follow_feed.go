package main

import (
	"context"
	"fmt"
	"time"

	"github.com/duuuuuuke/gator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, currentUser database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <feed_url>", cmd.Name)
	}

	feed_url := cmd.Args[0]
	feed, err := s.db.GetFeedByURL(context.Background(), feed_url)
	if err != nil {
		return fmt.Errorf("error getting feed by url: %w", err)
	}

	res, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    currentUser.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating feed follow: %w", err)
	}

	fmt.Println("Feed follow created:")
	printFeedFollow(res.UserName, res.FeedName)
	return nil
}

func handlerListFollowing(s *state, cmd command, currentUser database.User) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	follows, err := s.db.GetFeedFollowsForUser(context.Background(), currentUser.ID)
	if err != nil {
		return fmt.Errorf("error getting feed follows for user: %w", err)
	}

	if len(follows) == 0 {
		fmt.Println("You are not following any feeds.")
		return nil
	}

	fmt.Printf("Feed follows for user %s:\n", currentUser.Name)
	for _, follow := range follows {
		fmt.Printf("* %s\n", follow.FeedName)
	}

	return nil
}

func handlerUnfollowFeed(s *state, cmd command, currentUser database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <feed_url>", cmd.Name)
	}

	feed_url := cmd.Args[0]
	feed, err := s.db.GetFeedByURL(context.Background(), feed_url)
	if err != nil {
		return fmt.Errorf("error getting feed by url: %w", err)
	}

	err = s.db.DeleteFeedFollowByIDs(context.Background(), database.DeleteFeedFollowByIDsParams{
		UserID: currentUser.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error deleting feed follow: %w", err)
	}

	fmt.Printf("Unfollowed feed %s successfully!\n", feed.Name)
	return nil
}

func printFeedFollow(userName, feedName string) {
	fmt.Printf("* User:          %s\n", userName)
	fmt.Printf("* Feed:          %s\n", feedName)
}
