package configs

type (
	Config struct {
		GooglePhotosConfig GooglePhotosConfig `toml:"google_photos"`
		PhotoPrismConfig   PhotoPrismConfig   `toml:"photoprism"`
	}

	GooglePhotosConfig struct {
		CredentialPath string `toml:"credential_path"`
		ImagePath      string `toml:"image_path"`
		AlbumId        string `toml:"album_id"`
	}

	PhotoPrismConfig struct {
		SourcePath      string     `toml:"source_path"`
		DestinationPath string     `toml:"destination_path"`
		Databases       []Database `toml:"databases"`
	}

	Database struct {
		Host     string `toml:"host"`
		Port     string `toml:"port"`
		User     string `toml:"user"`
		Password string `toml:"password"`
		DBName   string `toml:"dbname"`
	}
)
