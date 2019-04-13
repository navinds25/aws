package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/service/apigatewayv2"
	"github.com/navinds25/aws-mission-ctrl/pkg/common"
)

func main() {
	sess, _, err := common.AWSSession()
	if err != nil {
		log.Fatal(err)
	}
	agsess := apigatewayv2.New(sess)
	//createApiInput := apigatewayv2.CreateApiInput{
	//
	//}
	//agsess.CreateApi()
	output, err := agsess.GetApis(&apigatewayv2.GetApisInput{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(output)
}
