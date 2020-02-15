package main

import (
	"errors"
	"fmt"
	"log"

	"google.golang.org/api/drive/v3"

	"github.com/zrma/uds-go/pkg/api"
)

func getBaseFolder(driveService *api.Service) error {
	r, err := driveService.Files.List().
		Q("properties has {key='udsRoot' and value='true'} and trashed=false").
		PageSize(1).
		Fields("nextPageToken, files(id, name, properties)").Do()
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to retrieve files: %v", err))
	}
	fmt.Println("Files:")
	if len(r.Files) == 0 {
		fmt.Println("No files found.")

		file, err := driveService.Files.Create(&drive.File{
			Name:       "UDS Root",
			MimeType:   "application/vnd.google-apps.folder",
			Properties: map[string]string{"udsRoot": "true"},
			Parents:    []string{},
		}).Fields("id").Do()
		if err != nil {
			return errors.New(fmt.Sprintf("Unable to create folder: %v", err))
		}
		fmt.Println(file)
		fmt.Println(file.Name)

	} else {
		for _, i := range r.Files {
			fmt.Printf("%s (%s)\n", i.Name, i.Id)
		}
	}
	return nil
}

func main() {
	driveService, err := api.NewService()
	if err != nil {
		log.Fatalf("Unable to retrieve NewService: %v", err)
	}

	if err := getBaseFolder(driveService); err != nil {
		log.Fatalln(err)
	}
}
