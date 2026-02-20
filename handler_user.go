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



func printUser(user database.User) {
	fmt.Printf(" * ID:      %v\n", user.ID)
	fmt.Printf(" * Name:    %v\n", user.Name)
}