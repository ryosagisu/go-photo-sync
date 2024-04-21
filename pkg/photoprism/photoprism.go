package photoprism

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"google-photo-sync/configs"
	"google-photo-sync/pkg/common"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func Init() *Sync {
	log.Println("Starting photo sync...")

	cfg := configs.ReadConfig(os.Getenv("CONFIG_PATH"))
	ppg := cfg.PhotoPrismConfig
	ppg.ImagePath = os.Getenv("IMAGE_PATH")
	return &Sync{
		Config: ppg,
	}
}

func (pps *Sync) SyncImage() {
	// Capture connection properties.

	// Get a database handle.
	for _, dbConfig := range pps.Config.Databases {
		db, err := sqlx.Connect("mysql", dbConfig.GetDSN())
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("Connected to %s\n", dbConfig.Name)
		favoritePhotos(db, pps.Config.SourcePath, pps.Config.ImagePath, dbConfig.Name)
		err = db.Close()
		if err != nil {
			log.Printf("failed to close connection: %v", err)
		}
	}
}

// favoritePhotos queries for albums that have the specified artist name.
func favoritePhotos(db *sqlx.DB, sourcePath, destinationPath, name string) {
	// An albums slice to hold data from returned rows.
	var photos []Photos
	imagePath := fmt.Sprintf("%s/%s", destinationPath, name)
	localImages := common.ListLocalImages(imagePath)

	err := db.Select(&photos, "select p.photo_path, p.photo_name, f.file_name from photos p left join files f on p.photo_uid = f.photo_uid where p.photo_favorite = 1 and p.photo_type = 'image' and f.file_root = '/'")
	if err != nil {
		log.Printf("query error: %v", err)
		return
	}

	log.Printf("Syncing %d images\n", len(photos))
	for _, photo := range photos {
		filename := filepath.Base(photo.FileName)
		targetPath := fmt.Sprintf("%s/%s", imagePath, filename)
		filePath := fmt.Sprintf("%s/%s", sourcePath, photo.FileName)
		err := CopyFile(filePath, targetPath)
		if err != nil {
			log.Printf("failed to copy file from %s to %s: %v", filePath, targetPath, err)
			continue
		}

		log.Printf("Copying %s.jpg\n", photo.PhotoName)
		if localImages[photo.PhotoName] {
			localImages[photo.PhotoName] = false
		}
	}

	log.Println("Delete missing images...")
	common.DeleteLocalFile(imagePath, localImages)
}

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = copyFileContents(src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func (ph *Photos) getFilePath(sourcePath string) string {
	return fmt.Sprintf("%s/%s/%s.jpg", sourcePath, ph.PhotoPath, ph.PhotoName)
}
