package main 

import (
	"context"
	"fmt"
	"time"

	"gator/internal/database"
	"github.com/google/uuid"	
)

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Enter name")
	}
	name := cmd.args[0]

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
	})
	if err != nil {
		return fmt.Errorf("Failed to create user %w", err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("Failed to set user: %w", err)
	}

	fmt.Println("User created:")
	printUser(user)
	return nil


}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Enter username")
	}
	name := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), name)

	if err != nil {
		return fmt.Errorf("User does not exist: %w", err)
	}


	err = s.cfg.SetUser(name)

	if err != nil {
		return fmt.Errorf("Error in setting user: %w", err)
	}
	fmt.Println("User set successful")
	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("No arguments needed for reset")
	}
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error while reseting database: %w", err)
	}
	
	fmt.Println("Database reset")
	return nil
	
}

func handlerUsers(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("No arguments needed for users command")
	}
	users, err := s.db.GetUsers(context.Background()) 
	if err != nil {
		return fmt.Errorf("Error while getting users: %w", err)
	}
	for _, user := range users {
		if s.cfg.CurrentUserName == user.Name {
			fmt.Printf("* %s (current)\n", user.Name)
		}else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("Error while running aggergation: %w", err)
	}
	fmt.Printf("Feed: %+v\n", feed)
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {


	if len(cmd.args) != 2 {
		return fmt.Errorf("Need name and url of feed")
	}

	name := cmd.args[0]
	url := cmd.args[1]

	
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: user.ID,
		Name: name,
		Url: url,
	})
	if err != nil {
		return fmt.Errorf("Error while creating feed: %w", err)
	}
	fmt.Println("Feed created")

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:user.ID,
		FeedID:feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Error when creating feed follow: %w", err)
	}

	//fmt.Printf("Feed: %v\n", feed)
	printFeed(feed)
	fmt.Println("=================================")
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("No arguments needed for feeds command")
	}
	Feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Error while getting feeds: %w", err)
	}
	if len(Feeds) == 0 {
		fmt.Println("No feeds in database.")
		return nil
	}

	fmt.Println("feeds found")

	for _, feed := range Feeds {
		fmt.Printf("Feed: (Name: %s, Url: %s, User: %s)\n", feed.Name, feed.Url, feed.Name_2)
		fmt.Println("======================================")
	}

	return nil



	//fmt.Println(Feeds)
	//return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Need to provide URL")
	}
	url := cmd.args[0]


	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Error when getting feed: %w", err)
	}

	follows, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:user.ID,
		FeedID:feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Error when creating feed follow: %w", err)
	}

	fmt.Printf("%s followed by %s\n",follows.FeedName, follows.UserName)
	return nil

}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("No args needed for following command")
	}

	followed, err := s.db.GETFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("Error when getting user follows: %w", err)
	}

	fmt.Println(followed)

	return nil

}

func printUser(user database.User) {
	fmt.Printf(" * ID:      %v\n", user.ID)
	fmt.Printf(" * Name:    %v\n", user.Name)
}


func printFeed(feed database.Feed) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* UserID:        %s\n", feed.UserID)
}