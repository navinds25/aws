package main

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/navinds25/mission-ctrl/pkg/servers"
	"github.com/navinds25/mission-ctrl/pkg/storage"
	log "github.com/sirupsen/logrus"
)

func main() {
	// cli/config should replace
	//filename := os.Args[1]
	//_, err := os.Stat(filename)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//iamSess, err := mciam.Session()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//roles, err := mciam.ListRoles(iamSess)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Println(roles)
	//if err := mciam.CreateRole(iamSess); err != nil {
	//	log.Fatal(err)
	//}
	//os.Exit(0)
	bucketName := "mordor.amis"
	filename := "epp4.4.0.6.raw"
	// Takes filename, format
	s3sess, err := storage.Session()
	if err != nil {
		log.Fatal(err)
	}
	createBucketOutput, err := storage.CreateBucket(s3sess, bucketName, s3.BucketCannedACLPrivate)
	if err != nil {
		log.Fatal(err)
	}
	log.Info(createBucketOutput)
	//uploadFileOutput, err := storage.MultiPartUploadFile(bucketName, filename)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Info(uploadFileOutput)
	ec2sess, err := servers.Session()
	if err != nil {
		log.Fatal(err)
	}
	if err := servers.RegisterAMI(ec2sess, &servers.RegisterAMIOptions{
		Architecture: "x86_64",
		Description:  "EPP Server",
		Format:       "raw",
		Bucket:       bucketName,
		Key:          filename,
	}); err != nil {
		log.Fatal(err)
	}
}
