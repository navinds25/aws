package app

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/navinds25/mission-ctrl/pkg/queryengine"
)

// ListSecrets lists the secrets
func ListSecrets(s3sess *s3.S3) error {
	queryengine.Session()
	return nil
}
