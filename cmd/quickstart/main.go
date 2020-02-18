package main

import (
	"log"

	"github.com/zrma/uds-go/pkg/api"
)

func main() {
	driveService, err := api.NewService()
	if err != nil {
		log.Fatalf("Unable to retrieve NewService: %v", err)
	}

	if _, err := driveService.GetBaseFolder(); err != nil {
		log.Fatalln(err)
	}
}
