package configs

type (
	Config struct {
		GooglePhotosConfig GooglePhotosConfig `toml:"google_photos"`
		PhotoPrismConfig   PhotoPrismConfig   `toml:"photoprism"`
	}

	GooglePhotosConfig struct {
		CredentialPath string `toml:"credential_path"`
		AlbumId        string `toml:"album_id"`
		ImagePath      string
	}

	PhotoPrismConfig struct {
		SourcePath string     `toml:"source_path"`
		Databases  []Database `toml:"databases"`
		ImagePath  string
	}

	Database struct {
		Host     string `toml:"host"`
		Port     string `toml:"port"`
		User     string `toml:"user"`
		Password string `toml:"password"`
		Name     string `toml:"dbname"`
	}
)
