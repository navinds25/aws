package mcec2

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/navinds25/aws-mission-ctrl/pkg/mccfn"
)

// EC2Session returns an ec2 session
func EC2Session() (*ec2.EC2, error) {
	sess, c, err := mccfn.AWSSession()
	if err != nil {
		return nil, err
	}
	return ec2.New(sess, c), nil
}
