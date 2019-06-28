package mcs3

import (
	"os"

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
func CreateBucket(sess *s3.S3, bucketname, acl string) (*s3.CreateBucketOutput, error) {
	input := s3.CreateBucketInput{
		Bucket: &bucketname,
		ACL:    &acl,
	}
	output, err := sess.CreateBucket(&input)
	if err != nil {
		return nil, err
	}
	return output, nil
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
