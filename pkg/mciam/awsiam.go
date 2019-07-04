package mciam

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/navinds25/mission-ctrl/pkg/common"

	"github.com/aws/aws-sdk-go/aws"
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

type vmimportTrustPolicy interface{}

// CreateRole creates a role
func CreateRole(sess *iam.IAM) error {
	//vmimportTrustPolicyDoc := `{"Version": "2012-10-17", "Statement": [{"Effect": "Allow", "Principal": { "Service": "vmie.amazonaws.com" }, "Action": "sts:AssumeRole", "Condition": {"StringEquals":{"sts:Externalid": "vmimport"}}}]}`
	file, err := os.Open("config/iam/roles/vmimport_trustpolicy.json")
	if err != nil {
		return err
	}
	vmimportTrustPolicyDoc, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	var vmimportTrustPolicyInst vmimportTrustPolicy
	if err := json.Unmarshal([]byte(vmimportTrustPolicyDoc), &vmimportTrustPolicyInst); err != nil {
		return err
	}
	assumeRolePolicyDoc := fmt.Sprintf("%v", vmimportTrustPolicyInst)
	output, err := sess.CreateRole(&iam.CreateRoleInput{
		RoleName:                 aws.String("vmimport"),
		Description:              aws.String("for creating AMIs"),
		AssumeRolePolicyDocument: aws.String(assumeRolePolicyDoc),
	})
	if err != nil {
		return err
	}
	log.Println(output)
	return nil
}
