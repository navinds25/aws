package lambda

import (
	"log"

	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/navinds25/aws-mission-ctrl/pkg/common"
)

// Session returns a session for AWS lambda
func Session() (*lambda.Lambda, error) {
	sess, _, err := common.AWSSession()
	if err != nil {
		return nil, err
	}
	lambdaSess := lambda.New(sess)
	return lambdaSess, nil
}

// AddLambda adds a lambda function
func AddLambda() error {
	lsess, err := Session()
	if err != nil {
		return err
	}
	log.Println(lsess.APIVersion)
	//lambda.CreateFunctionInput{}
	//lsess.CreateFunction()
	return nil
}
