package storage

import (
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/navinds25/mission-ctrl/pkg/common"
)

// Session returns an s3 session
func Session() (*s3.S3, error) {
	sess, c, err := common.AWSSession()
	if err != nil {
		return nil, err
	}
	return s3.New(sess, c), nil
}

// ListBuckets returns all the buckets in aws S3
func ListBuckets(sess *s3.S3) (*s3.ListBucketsOutput, error) {
	input := s3.ListBucketsInput{}
	resp, err := sess.ListBuckets(&input)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ListObjectsinBucket returns a list of objects for a given bucket.
func ListObjectsinBucket(sess *s3.S3, bucketname string) (*s3.ListObjectsOutput, error) {
	input := s3.ListObjectsInput{
		Bucket: aws.String(bucketname),
	}
	objects, err := sess.ListObjects(&input)
	if err != nil {
		return nil, err
	}
	return objects, nil
}

// CreateBucket creates an S3 bucket
func CreateBucket(sess *s3.S3, bucketname, acl string) (string, error) {
	input := s3.CreateBucketInput{
		Bucket: aws.String(bucketname),
		ACL:    aws.String(acl),
	}
	output, err := sess.CreateBucket(&input)
	if err != nil {
		if awserr, ok := err.(awserr.Error); ok {
			switch awserr.Code() {
			case s3.ErrCodeBucketAlreadyOwnedByYou:
				return s3.ErrCodeBucketAlreadyOwnedByYou, nil
			case s3.ErrCodeBucketAlreadyExists:
				return s3.ErrCodeBucketAlreadyExists, nil
			}
		}
		return "", err
	}
	return *output.Location, nil
}

// KeyExists checks if a key exists in a bucket
func KeyExists(sess *s3.S3, bucketname, key string) (bool, error) {
	_, err := sess.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucketname),
		Key:    aws.String(key),
	})
	if err != nil {
		if awserr, ok := err.(awserr.Error); ok {
			switch awserr.Code() {
			case s3.ErrCodeNoSuchBucket:
				return false, nil
			case s3.ErrCodeNoSuchKey:
				return false, nil
			case "NotFound":
				return false, nil
			default:
				return false, err
			}
		}
	}
	return true, nil
}

// DownloadFile downloads a file from S3
func DownloadFile(s3sess *s3.S3, bucketName, key, fileName string) error {
	downloader := s3manager.NewDownloaderWithClient(s3sess)
	tmpFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	_, err = downloader.Download(tmpFile, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}
	return nil
}

// UploadFile uploads the file to the S3 bucket
// filename is the local file to be uploaded
// key is the name of the file in thr bucket
func UploadFile(sess *s3.S3, bucket, filename string) (*s3.PutObjectOutput, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	key := filepath.Base(filename)
	input := &s3.PutObjectInput{
		Body:   aws.ReadSeekCloser(file),
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	output, err := sess.PutObject(input)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// MultiPartUploadFile uploads a large file.
func MultiPartUploadFile(bucket, filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	key := filepath.Base(filename)
	awsSess, _, err := common.AWSSession()
	if err != nil {
		return "", err
	}
	uploader := s3manager.NewUploader(awsSess, func(u *s3manager.Uploader) {
		u.PartSize = 5 * 1024 * 1024
		u.Concurrency = 4
	})
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		return "", err
	}
	return result.Location, nil
}
