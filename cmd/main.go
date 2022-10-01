package main

import (
	"flag"
	"google-photo-sync/pkg/photoprism"
	"log"
	"os"

	"google-photo-sync/pkg/google_photos"
)

const (
	SyncImage = "SyncImage"
	ListAlbum = "ListAlbum"

	GoogleService     = "Google"
	PhotoprismService = "Photoprism"
)

func main() {
	command := flag.String("command", SyncImage, "command to execute")
	flag.Parse()

	service := os.Getenv("SYNC_SERVICE")
	if service == "" {
		service = GoogleService
	}

	switch service {
	case GoogleService:
		googleSync(*command)
	case PhotoprismService:
		photoprismSync(*command)
	default:
		log.Printf("unspported service: %s\n", service)
	}

	log.Println("Good bye...")
}

func googleSync(cmd string) {
	log.Println("Initializing google sync")
	gp := google_photos.Init()

	switch cmd {
	case SyncImage:
		gp.SyncImage()
	case ListAlbum:
		gp.ListAlbum()
	default:
		log.Printf("unsupported command: %s\n", cmd)
	}
}

func photoprismSync(cmd string) {
	log.Println("Initializing photoprism sync")
	pp := photoprism.Init()

	switch cmd {
	case SyncImage:
		pp.SyncImage()
	default:
		log.Printf("unsupported command: %s\n", cmd)
	}
}
