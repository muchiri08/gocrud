package main

import (
	"flag"
	"fmt"
	"github.com/muchiri08/crud/api"
	"github.com/muchiri08/crud/storage"
)

func main() {
	migrate := flag.Bool("migrate", false, "initialise database tables")
	flag.Parse()
	store, err := storage.NewPostgresStore()
	if err != nil {
		fmt.Println(err)
		return
	}
	address := api.Address{
		Host: "localhost",
		Port: ":8080",
	}

	server := api.NewApiServer(address, store)

	if *migrate {
		server.Store.RunMigrationScript()
		server.Store.InitAdmin()
	}

	server.Run()
}
