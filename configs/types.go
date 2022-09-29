package configs

type Config struct {
	CredentialFile string `yaml:"credential_file"`
	CredentialPath string `yaml:"credential_path"`
	OutputPath     string `yaml:"output_path"`
	AlbumId        string `yaml:"album_id"`
	AlbumName      string `yaml:"album_name"`
}
