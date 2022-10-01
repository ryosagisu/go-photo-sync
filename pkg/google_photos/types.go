package google_photos

import (
	"net/http"

	"google-photo-sync/configs"

	"google.golang.org/api/photoslibrary/v1"
)

var (
	client *http.Client
)

type (
	Sync struct {
		Config        configs.GooglePhotosConfig
		PhotosLibrary *photoslibrary.Service
		Client        *http.Client
	}
)
