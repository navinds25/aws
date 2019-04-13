package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/navinds25/aws-mission-ctrl/pkg/common"
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
