package app

import "github.com/aws/aws-sdk-go/service/s3"

// UserPass for the username, password.
type UserPass struct {
	Key         string `json:"key"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Description string `json:"description"`
}

// PassDB is the Main Struct for password file
type PassDB struct {
	Data []*UserPass `json:"data"`
}

// SecretInstance is the main struct for secrets
type SecretInstance struct {
	S3Instance s3.S3
	Bucket     string
	Passfile   string
}
