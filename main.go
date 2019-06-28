package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/navinds25/mission-ctrl/pkg/mcs3"
	log "github.com/sirupsen/logrus"
)

// UserPass for the username, password.
type UserPass struct {
	Key         string `json:"key"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Description string `json:"description"`
	Website     string `json:"website"`
}

// PassDB is the Main Struct for password file
type PassDB struct {
	Data []*UserPass `json:"data"`
}

// BucketName is the name of the bucket that stores the data
var BucketName = "kms-pass-evive"

// PassFile is the name of the file that holds the passwords.
var PassFile = "dunno2"

func getInput(fieldName string) (fieldValue string) {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("%s->", fieldName)
		fieldValue, err := reader.ReadString('\n')
		if err != nil {
			log.Error(err)
			continue
		}
		fieldValue = strings.Replace(fieldValue, "\n", "", -1)
		return fieldValue
	}
}

// GetTmpFileData returns the struct containing the data from the temp file
func GetTmpFileData(fileName string) (*PassDB, error) {
	decodedData := &PassDB{}
	tmpFile, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	if err := json.NewDecoder(tmpFile).Decode(decodedData); err != nil {
		return nil, err
	}
	log.Printf("TmpData: %+v\n", decodedData)
	return decodedData, nil
}

func createPassFile(s3sess *s3.S3) error {
	dummyJSON, err := json.Marshal(PassDB{})
	if err != nil {
		return err
	}
	dummyfile := filepath.Join("/tmp", PassFile+"_dummy")
	if err := ioutil.WriteFile(dummyfile, dummyJSON, 0644); err != nil {
		return err
	}
	dummyfd, err := os.Open(dummyfile)
	if err != nil {
		return err
	}
	output, err := s3sess.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(BucketName),
		Key:    aws.String(PassFile),
		Body:   dummyfd,
	})
	if err != nil {
		return err
	}
	log.Info(output)
	return nil
}

// InitPassFile initializes the PassFile
func InitPassFile(s3sess *s3.S3) error {
	_, err := s3sess.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(BucketName),
		Key:    aws.String(PassFile),
	})
	if err != nil {
		if awserr, ok := err.(awserr.Error); ok {
			switch awserr.Code() {
			case s3.ErrCodeNoSuchBucket:
				log.Error("bucket does not exist")
			case s3.ErrCodeNoSuchKey:
				log.Error("No such key")
			case "NotFound":
				log.Error("NotFound")
				if err := createPassFile(s3sess); err != nil {
					return err
				}
			default:
				log.Error("Dunno")
			}
		}
	} else {
		log.Info("PassFile exists: ", PassFile)
	}
	return nil
}

// AddSecret adds a secret to the password file
func AddSecret(s3sess *s3.S3) error {
	userpass := &UserPass{
		Key:         getInput("key to identify"),
		Username:    getInput("username"),
		Password:    getInput("password"),
		Description: getInput("description"),
		Website:     getInput("website/url"),
	}
	log.Infof("userpass: %+v\n", userpass)

	// TMPFileName
	tempS3FileName := filepath.Join("/tmp", PassFile+"_s3")
	tempNewFileName := filepath.Join("/tmp", PassFile+"_new")

	// Downloading file from S3
	if err := DownloadFile(s3sess, BucketName, tempS3FileName); err != nil {
		return err
	}
	log.Info("Downloaded file")
	// Read downloaded data
	tmpData, err := GetTmpFileData(tempS3FileName)
	if err != nil {
		return err
	}
	newData := PassDB{}
	newData.Data = append(tmpData.Data, userpass)
	log.Infof("newData: %+v", newData)

	// Create updated file
	newDataBytes, err := json.Marshal(newData)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(tempNewFileName, newDataBytes, 0644); err != nil {
		return err
	}

	// Upload the new file
	tempNewFD, err := os.Open(tempNewFileName)
	output, err := s3sess.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(BucketName),
		Key:    aws.String(PassFile),
		Body:   tempNewFD,
	})
	log.Info(output)

	// Clean up temp files
	if err := os.Remove(tempNewFileName); err != nil {
		return err
	}
	if err := os.Remove(tempS3FileName); err != nil {
		return err
	}
	return nil
}

// DownloadFile downloads a file from S3
func DownloadFile(s3sess *s3.S3, bucketName string, fileName string) error {
	downloader := s3manager.NewDownloaderWithClient(s3sess)
	tmpFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	_, err = downloader.Download(tmpFile, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(PassFile),
	})
	if err != nil {
		return err
	}
	return nil
}

// ListSecrets lists the secrets
func ListSecrets(s3sess *s3.S3) error {
	return nil
}

func main() {
	s3sess, err := mcs3.Session()
	if err != nil {
		log.Fatal(err)
	}
	if err := InitPassFile(s3sess); err != nil {
		log.Fatal(err)
	}
	if err := AddSecret(s3sess); err != nil {
		log.Fatal(err)
	}
	if err := ListSecrets(s3sess); err != nil {
		log.Fatal(err)
	}
}
