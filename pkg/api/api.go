package api

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/spf13/afero"
	"golang.org/x/net/context"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"

	"github.com/zrma/uds-go/pkg/uds"
)

var AppFs afero.Fs = nil

func init() {
	AppFs = afero.NewOsFs()
}

// NewService function returns initialized Service object's pointer
func NewService() (*Service, error) {
	api := &Service{}
	if err := api.Init(); err != nil {
		return nil, err
	}
	return api, nil
}

// Service struct is google api service wrapper.
type Service struct {
	*drive.Service
	ctx context.Context
}

// Init works internally but public(export) for using in apt_test package
func (api *Service) Init() error {
	_, caller, _, _ := runtime.Caller(2)
	basePath := filepath.Dir(caller)

	afs := &afero.Afero{Fs: AppFs}
	b, err := afs.ReadFile(filepath.Join(basePath, credentialFile))
	if err != nil {
		return err
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := Helper.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		return err
	}

	token, err := Helper.GetToken(config, filepath.Join(basePath, tokenFile))
	if err != nil {
		return err
	}

	api.ctx = context.Background()
	driveService, err := drive.NewService(
		api.ctx,
		option.WithTokenSource(config.TokenSource(api.ctx, token)),
	)
	if err != nil {
		return err
	}

	api.Service = driveService
	return nil
}

// GetBaseFolder locate the base UDS folder
func (api *Service) GetBaseFolder() (*drive.File, error) {
	r, err := api.Files.List().
		Q("properties has {key='udsRoot' and value='true'} and trashed=false").
		PageSize(1).
		Fields("nextPageToken, files(id, name, properties)").Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve files: %v", err)
	}

	fileLength := len(r.Files)
	if fileLength == 0 {
		fmt.Println("No files found.")
		return api.createRootFolder()
	} else if fileLength == 1 {
		for _, i := range r.Files {
			fmt.Printf("%s (%s)\n", i.Name, i.Id)
		}
		return r.Files[0], nil
	}
	return nil, fmt.Errorf("multiple UDS Roots found")
}

func (api *Service) createRootFolder() (*drive.File, error) {
	return api.Files.Create(&drive.File{
		Name:       "UDS Root",
		MimeType:   "application/vnd.google-apps.folder",
		Properties: map[string]string{"udsRoot": "true"},
		Parents:    []string{},
	}).Fields("id").Do()
}

func (api *Service) CreateMediaFolder(media *uds.File) (*drive.File, error) {
	return api.Files.Create(&drive.File{
		Name:     media.Name,
		MimeType: "application/vnd.google-apps.folder",
		Properties: map[string]string{
			"udsRoot":      "true",
			"size":         media.Size,
			"size_numeric": media.SizeNumeric,
			"encoded_size": media.EncodedSize,
			"md5":          media.MD5,
		},
		Parents: media.Parents,
	}).Fields("id").Do()
}

func (api *Service) ListFiles(query string) ([]*uds.File, error) {
	q := "properties has {key='uds' and value='true'} and trashed=false"
	if query != "" {
		q += fmt.Sprintf(" and name contains '%s'", query)
	}

	r, err := api.Files.List().
		Q(q).
		PageSize(1000).
		Fields("nextPageToken, files(id, name, properties, mimeType)").Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve files: %v", err)
	}

	var files []*uds.File
	for _, f := range r.Files {
		props := f.Properties
		file := &uds.File{
			Name:        props["name"],
			Mime:        props["mimeType"],
			Size:        props["size"],
			EncodedSize: props["encoded_size"],
			SizeNumeric: props["size_numeric"],
			ID:          props["id"],
			MD5:         props["md5Checksum"],
			Shared:      props["shared"] == "true",
		}
		file.Init()
		files = append(files, file)
	}
	return files, nil
}
