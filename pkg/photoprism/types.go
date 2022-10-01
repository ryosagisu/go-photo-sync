package photoprism

import "google-photo-sync/configs"

type (
	Sync struct {
		Config configs.PhotoPrismConfig
	}

	Photos struct {
		PhotoPath string
		PhotoName string
	}
)
