package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

// IAMSession returns an iam object
func IAMSession(region string) *iam.IAM {
	c := &aws.Config{
		Region: aws.String(region),
	}
	sess, err := session.NewSession(c)
	if err != nil {
		log.Fatal(err)
	}
	svc := iam.New(sess, c)
	return svc
}

func main() {
	sess := IAMSession("ap-south-1")

	//	input := iam.GetAccountAuthorizationDetailsInput{}
	//	resp, err := svc.GetAccountAuthorizationDetails(&input)
	//	if err != nil {
	//		log.Fatal(err.Error())
	//	}
	//	//log.Println(resp.Policies)
	//	for _, policyDetail := range resp.Policies {
	//		log.Println(policyDetail)
	//	}

	// list users
	//input := iam.ListUsersInput{}
	//resp, _ := sess.ListUsers(&input)
	//input := iam.ListRolesInput{}
	//resp, _ := sess.ListRoles(&input)
	//input := iam.CreateUserInput{
	//	UserName: "vmimport"
	//	PermissionsBoundary:
	//	Tags:
	//	Path: "/system/"
	//}
	//resp, err := svc.CreateUser(&input)

	log.Println(resp)
}
