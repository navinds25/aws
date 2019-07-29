package main

import (
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/navinds25/mission-ctrl/pkg/storage"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("not enough arguments. arg1 = filename, arg2 = bucketname")
	}
	filename := os.Args[2]
	_, err := os.Stat(filename)
	if err != nil {
		log.Fatal(err)
	}
	bucketname := os.Args[1]
	loc, err := storage.MultiPartUploadFile(bucketname, filename)
	if err != nil {
		log.Fatal(err)
	}
	log.Info(loc)
}
