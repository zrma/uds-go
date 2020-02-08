package main

import (
	"fmt"
	"log"

	"google.golang.org/api/drive/v3"

	"github.com/zrma/uds-go/pkg/api"
)

func main() {
	driveService, err := api.NewService()
	if err != nil {
		log.Fatalf("Unable to retrieve NewService: %v", err)
	}

	r, err := driveService.Files.List().
		Q("properties has {key='udsRoot' and value='true'} and trashed=false").
		PageSize(1).
		Fields("nextPageToken, files(id, name, properties)").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
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
			log.Fatalf("Unable to create folder: %v", err)
		}
		fmt.Println(file)
		fmt.Println(file.Name)

	} else {
		for _, i := range r.Files {
			fmt.Printf("%s (%s)\n", i.Name, i.Id)
		}
	}
}
