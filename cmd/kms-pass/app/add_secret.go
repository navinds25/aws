package app

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/navinds25/mission-ctrl/pkg/storage"
)

// InitPassFile initializes the PassFile
func InitPassFile(s3sess *s3.S3, bucketname, passfile string) error {
	exists, err := storage.KeyExists(s3sess, bucketname, passfile)
	if err != nil {
		return err
	}
	if !exists {
		if err := createPassFile(s3sess); err != nil {
			return err
		}
	}
	return nil
}

func createPassFile(s3sess *s3.S3, bucketname, passfile string) error {
	dummyJSON, err := json.Marshal(PassDB{})
	if err != nil {
		return err
	}
	dummyfile := filepath.Join("/tmp", passfile+"_dummy")
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
	if err := storage.DownloadFile(s3sess, BucketName, PassFile, tempS3FileName); err != nil {
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
