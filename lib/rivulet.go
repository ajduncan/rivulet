package rivulet

import (
	"fmt"
)

func NewRivulet(pwd string) {
	// we need a database to store things
	db, err := NewDatabase(pwd)
	if err != nil {
		fmt.Println("Error initializing database.")
		return
	}

	// we need a raw server to pass things around
	server := NewRivuletServer("default", *db)

	// we need a web api thing.
	api := NewAPI(server, db)

	// setup the handlers
	api.Init()

	api.Run()
	server.Run()
}
