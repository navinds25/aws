package servers

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/navinds25/mission-ctrl/pkg/common"
)

// Session returns an ec2 session
func Session() (*ec2.EC2, error) {
	sess, c, err := common.AWSSession()
	if err != nil {
		return nil, err
	}
	return ec2.New(sess, c), nil
}

// RegisterAMIOptions is the struct for the options to RegisterAMI
type RegisterAMIOptions struct {
	Architecture string // eg: x86_64
	Description  string
	Format       string
	Bucket       string
	Key          string
}

// RegisterAMI creates an ami
func RegisterAMI(sess *ec2.EC2, opts *RegisterAMIOptions) error {
	imagediskcont := &ec2.ImageDiskContainer{
		Format: &opts.Format,
		UserBucket: &ec2.UserBucket{
			S3Bucket: aws.String(opts.Bucket),
			S3Key:    aws.String(opts.Key),
		},
	}
	imagediskcontainers := append([]*ec2.ImageDiskContainer{}, imagediskcont)
	input := &ec2.ImportImageInput{
		Architecture:   aws.String(opts.Architecture),
		Description:    aws.String(opts.Description),
		DiskContainers: imagediskcontainers,
	}
	output, err := sess.ImportImage(input)
	if err != nil {
		return err
	}
	log.Println(output)
DESCRIBEIMAGETASK:
	descripInput := &ec2.DescribeImportImageTasksInput{}
	descripOutput, err := sess.DescribeImportImageTasks(descripInput)
	if err != nil {
		return err
	}
	log.Println(descripOutput)
	for _, out := range descripOutput.ImportImageTasks {
		active := "active"
		if out.Status == &active {
			goto DESCRIBEIMAGETASK
		}
	}
	return nil
}
