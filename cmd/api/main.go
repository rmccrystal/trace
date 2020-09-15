package main

import (
	"github.com/sirupsen/logrus"
	"trace/pkg/api"
)

func main() {
	/*
	config := database.Config{
		MongoURI:     "",
		DatabaseName: "",
	}

	db, err := database.Connect(config)

	if err != nil {
		logrus.Fatalf("Could not connect to database: %s", err)
	}
	 */

	addr := "0.0.0.0:8080"

	if err := api.Listen(addr, nil); err != nil {
		logrus.Fatalf("Failed to listen on %s", addr)
	}
}
