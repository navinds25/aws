package ami

import (
	"os"

	mcec2 "github.com/navinds25/aws-mission-ctrl/pkg/ec2"
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/s3"
)

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

// RegisterAMI creates an ami
func RegisterAMI() {
	ec2sess, err := mcec2.Session()
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
