package google_photos

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"google-photo-sync/configs"
	"google-photo-sync/pkg/common"

	"google.golang.org/api/photoslibrary/v1"
)

func Init() *Sync {
	log.Println("Starting google-photo-sync")

	cfg := configs.ReadConfig(os.Getenv("CONFIG_PATH"))
	gpc := cfg.GooglePhotosConfig
	svc := GetService(gpc.CredentialPath)

	client := &http.Client{
		Timeout: 5 * time.Minute,
	}

	return &Sync{
		Config:        gpc,
		PhotosLibrary: svc,
		Client:        client,
	}
}

func (gp *Sync) SyncImage() {
	localImages := common.ListLocalImages(gp.Config.ImagePath)
	var pageToken string
	log.Println("Downloading images...")
	for {
		req := &photoslibrary.SearchMediaItemsRequest{
			PageSize:  100,
			AlbumId:   gp.Config.AlbumId,
			PageToken: pageToken,
		}

		items, err := gp.PhotosLibrary.MediaItems.Search(req).Do()
		if err != nil {
			log.Fatalf("failed to search media: %v", err)
		}

		pageToken = items.NextPageToken
		for _, item := range items.MediaItems {
			fileName := getImageName(item.Id, gp.Config.ImagePath)
			if localImages[item.Id] {
				localImages[item.Id] = false
			}

			err = downloadImage(fileName, item.BaseUrl)
			if err != nil {
				log.Printf("Failed to download: %v\n", err)
			}
		}

		if pageToken == "" {
			break
		}
	}

	log.Println("Delete missing images...")
	common.DeleteLocalFile(gp.Config.ImagePath, localImages)
}

func getImageName(id, outputPath string) string {
	return fmt.Sprintf("%s/%s.jpg", outputPath, id)
}

// Skip download if file exist
func downloadImage(fileName, baseUrl string) error {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		url := fmt.Sprintf("%v=d", baseUrl)
		output, err := os.Create(fileName)
		if err != nil {
			return err
		}
		defer output.Close()

		response, err := client.Get(url)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		n, err := io.Copy(output, response.Body)
		if err != nil {
			return err
		}
		log.Printf("Downloaded '%v' (%v)", fileName, uint64(n))
	}
	return nil
}

func (gp *Sync) ListAlbum() {
	albums, err := gp.PhotosLibrary.Albums.List().Do()
	if err != nil {
		log.Fatalf("failed to list album: %v", err)
	}

	for _, album := range albums.Albums {
		log.Println(album.Id, album.Title)
	}
}
