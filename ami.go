package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/s3"
)

func awsSession() (*session.Session, *aws.Config, error) {
	region := "ap-south-1"
	c := &aws.Config{
		Region: aws.String(region),
	}
	sess, err := session.NewSession(c)
	if err != nil {
		return nil, nil, err
	}
	return sess, c, err
}

// S3Session returns an s3 session
func S3Session() (*s3.S3, error) {
	sess, c, err := awsSession()
	if err != nil {
		return nil, err
	}
	return s3.New(sess, c), nil
}

// EC2Session returns an ec2 session
func EC2Session() (*ec2.EC2, error) {
	sess, c, err := awsSession()
	if err != nil {
		return nil, err
	}
	return ec2.New(sess, c), nil
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

// UploadCentOSImage uploads the ova file to the S3 bucket
func UploadCentOSImage(sess *s3.S3, filename, key string) (*s3.PutObjectOutput, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	input := &s3.PutObjectInput{
		Body:   aws.ReadSeekCloser(file),
		Bucket: aws.String("centos7amis"),
		Key:    aws.String(key),
	}
	output, err := sess.PutObject(input)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func RegisterAMI() {
	ec2sess, err := EC2Session()
	if err != nil {
		log.Fatal(err)
	}
	architecture := "x86_64"
	description := "CentOS7_std"
	format := "OVA"
	S3Bucket := "centos7amis"
	S3Key := "centos7_std"

	imagediskcont := &ec2.ImageDiskContainer{
		Format: &format,
		UserBucket: &ec2.UserBucket{
			S3Bucket: &S3Bucket,
			S3Key:    &S3Key,
		},
	}
	imagediskcontainers := append([]*ec2.ImageDiskContainer{}, imagediskcont)
	input := &ec2.ImportImageInput{
		Architecture:   &architecture,
		Description:    &description,
		DiskContainers: imagediskcontainers,
	}
	output, err := ec2sess.ImportImage(input)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(output)
DESCRIBEIMAGETASK:
	descripInput := &ec2.DescribeImportImageTasksInput{}
	descripOutput, err := ec2sess.DescribeImportImageTasks(descripInput)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(descripOutput)
	for _, out := range descripOutput.ImportImageTasks {
		active := "active"
		if out.Status == &active {
			goto DESCRIBEIMAGETASK
		}
	}
}

func main() {
	// s3sess, err := S3Session()
	// if err != nil {
	//	 log.Fatal(err)
	// }
	RegisterAMI()
}
