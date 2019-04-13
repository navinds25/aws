package ec2

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/navinds25/aws-mission-ctrl/pkg/common"
)

// Session returns an ec2 session
func Session() (*ec2.EC2, error) {
	sess, c, err := common.AWSSession()
	if err != nil {
		return nil, err
	}
	return ec2.New(sess, c), nil
}
