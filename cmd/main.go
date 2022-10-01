package main

import (
	"flag"
	"google-photo-sync/pkg/photoprism"
	"log"
	"os"
	"time"

	"google-photo-sync/pkg/google_photos"
)

const (
	SyncImage = "SyncImage"
	ListAlbum = "ListAlbum"
	DEBUG     = "debug"

	GoogleService     = "google"
	PhotoprismService = "photoprism"
)

func main() {
	checkRequiredEnv()

	command := flag.String("command", SyncImage, "command to execute")
	flag.Parse()

	service := os.Getenv("SYNC_SERVICE")
	if service == "" {
		service = GoogleService
	}

	if *command == DEBUG {
		log.Println("Debugging mode")
		time.Sleep(1 * time.Hour)
		os.Exit(0)
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

func checkRequiredEnv() {
	if os.Getenv("CONFIG_PATH") == "" {
		log.Fatalln("CONFIG_PATH hasn't been set")
	}

	if os.Getenv("IMAGE_PATH") == "" {
		log.Fatalln("IMAGE_PATH hasn't been set")
	}
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
