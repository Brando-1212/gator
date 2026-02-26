package main

import (
	"database/sql"
	"log"
	"os"
	
	"gator/internal/config"
	"gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db *database.Queries
	cfg *config.Config
} 

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	dbQueries := database.New(db)
	State := &state{
		db:  dbQueries,
		cfg: &cfg,
	}


	cmds := commands{
		regCommands: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handlerAddFeed)
	cmds.register("feeds", handlerFeeds)

	if len(os.Args) < 2 {
		log.Fatal("Must provide command and args. usage: cli <command> [args]")
	}

	cmdname := os.Args[1]
	cmdargs := os.Args[2:]

	err = cmds.run(State, command{Name: cmdname, args: cmdargs})
	if err != nil {
		log.Fatal(err)
	}

}