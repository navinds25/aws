package common

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// Region for specifying aws region.
var Region string

// AWSSession returns a session (region is hardcoded.)
func AWSSession() (*session.Session, *aws.Config, error) {
	if Region == "" {
		Region = "ap-south-1"
	}
	c := &aws.Config{
		Region: aws.String(Region),
	}
	sess, err := session.NewSession(c)
	if err != nil {
		return nil, nil, err
	}
	return sess, c, err
}
