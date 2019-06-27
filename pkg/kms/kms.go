package kms

import (
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/navinds25/mission-ctrl/pkg/common"
)

// Session returns a kms session
func Session() (*kms.KMS, error) {
	sess, _, err := common.AWSSession()
	if err != nil {
		return nil, err
	}
	kmsSvc := kms.New(sess)
	return kmsSvc, nil
}

// ListCustomerKeys is for listing Customer Keys
func ListCustomerKeys(kmsSvc *kms.KMS) ([]*kms.DescribeKeyOutput, error) {
	res, err := kmsSvc.ListKeys(&kms.ListKeysInput{})
	if err != nil {
		return nil, err
	}
	customerKeys := []*kms.DescribeKeyOutput{}
	for _, key := range res.Keys {
		output, err := kmsSvc.DescribeKey(&kms.DescribeKeyInput{
			KeyId: key.KeyId,
		})
		if err != nil {
			return nil, err
		}
		outputKeyManager := *output.KeyMetadata.KeyManager
		if outputKeyManager == kms.KeyManagerTypeCustomer {
			customerKeys = append(customerKeys, output)
		}
	}
	return customerKeys, nil
}
