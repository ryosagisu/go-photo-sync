package photoprism

import "google-photo-sync/configs"

type (
	Sync struct {
		Config configs.PhotoPrismConfig
	}

	Photos struct {
		PhotoPath string `db:"photo_path"`
		PhotoName string `db:"photo_name"`
	}
)
