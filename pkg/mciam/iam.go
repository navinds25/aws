package mciam

import (
	"github.com/navinds25/mission-ctrl/pkg/common"

	"github.com/aws/aws-sdk-go/service/iam"
)

// Session returns an iam object
func Session() (*iam.IAM, error) {
	sess, c, err := common.AWSSession()
	if err != nil {
		return nil, err
	}
	svc := iam.New(sess, c)
	return svc, nil
}

// ListRoles lists the roles in AWS IAM.
func ListRoles(sess *iam.IAM) (*iam.ListRolesOutput, error) {
	resp, err := sess.ListRoles(&iam.ListRolesInput{})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ListUsers lists the IAM users in AWS.
func ListUsers(sess *iam.IAM) (*iam.ListUsersOutput, error) {
	resp, err := sess.ListUsers(&iam.ListUsersInput{})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// CreateUser creates an IAM user.
func CreateUser(sess *iam.IAM, userConf *iam.CreateUserInput) (*iam.CreateUserOutput, error) {
	resp, err := sess.CreateUser(userConf)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
