package main

import (
	"errors"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
)

const (
	US_EAST_VA       = "us-east-1"
	US_WEST_OR       = "us-west-2"
	US_WEST_CA       = "us-west-1"
	EU_IR            = "eu-west-1"
	EU_FK            = "eu-central-1"
	ASIA_SG          = "ap-southeast-1"
	ASIA_SD          = "ap-southeast-2"
	ASIA_TK          = "ap-northeast-1"
	SOUTH_AMERICA_SP = "sa-east-1"
)

const (
	SCHED_BUCKET = "api-schedule"
	//CRED_PROFILE = "aetworker"
	//API_REGION       = ASIA_SG
	//S3_BUCKET_REGION = US_WEST_OR
	// FOR Framwork
	CRED_FILE_PATH = "/Users/hawk/.aws/credentials"
	//FOR AWS_TEST
	//CRED_FILE_PATH = "../../../../../.aws/credentials"
	//FOR AWS_TEST MAC
	//CRED_FILE_PATH = "/Users/rachael_pai/.aws/credentials"
)

var Profile string = "goapi"

var Astra_AcessKey string
var Astra_SecretKey string

func GetCredentialChain() (*credentials.Credentials, *aws.Config) {
	config := aws.NewConfig()
	ec2m := ec2metadata.New(session.New(), config)
	var ProviderList []credentials.Provider = []credentials.Provider{
		&credentials.EnvProvider{},
		&credentials.SharedCredentialsProvider{Filename: CRED_FILE_PATH, Profile: Profile},
		&ec2rolecreds.EC2RoleProvider{
			Client: ec2m,
		},
	}
	creds := credentials.NewChainCredentials(ProviderList)
	return creds, config
	//return credentials.NewStaticCredentials(accessKey, secretKey, ``)
}

func GetCredential(accessKey, secretKey string) *credentials.Credentials {
	return credentials.NewStaticCredentials(accessKey, secretKey, ``)
}

func GetCredentialShared() *credentials.Credentials {
	return credentials.NewSharedCredentials(CRED_FILE_PATH, Profile)
}

func InitConfig(region string) (*aws.Config, error) {
	//creds := GetCredentialShared()
	creds, conf := GetCredentialChain()
	val, err := creds.Get()
	if err != nil {
		log.Println("InitConfig:", err)
	}
	log.Println("Cred AccessKeyID:", val.AccessKeyID)
	conf.WithRegion(region).WithCredentials(creds)
	return conf, nil
}

func InitConfigByKey(region string, access string, secret string) (*aws.Config, error) {
	if access != "" && secret != "" {
		creds := GetCredential(access, secret)
		config := aws.NewConfig().WithRegion(region).WithCredentials(creds)
		return config, nil
	} else {
		return nil, errors.New("Invalid access or secret key")
	}
}
