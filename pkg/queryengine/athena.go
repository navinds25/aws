package queryengine

import (
	"github.com/aws/aws-sdk-go/service/athena"
	"github.com/navinds25/mission-ctrl/pkg/common"
)

// Session returns a session for AWS Athena
func Session() (*athena.Athena, error) {
	sess, _, err := common.AWSSession()
	if err != nil {
		return nil, err
	}
	athenaSess := athena.New(sess)
	return athenaSess, nil
}
