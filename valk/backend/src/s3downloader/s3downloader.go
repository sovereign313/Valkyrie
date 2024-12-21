package s3downloader

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const (
	maxPartSize        = int64(512 * 1000)
	maxRetries         = 3
	awsAccessKeyID     = "AWS_ACCESS_KEY_ID"
	awsSecretAccessKey = "AWS_SECRET_ACCESS_KEY"
)

func CheckForBucket(bucketname string, region string) (bool, error) {
	creds := credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, "")
	_, err := creds.Get()
	if err != nil {
		return false, err
	}

	cfg := aws.NewConfig().WithRegion(region).WithCredentials(creds)
	sess, err := session.NewSession(cfg)
	svc := s3.New(sess)

	result, err := svc.ListBuckets(nil)
	if err != nil {
		return false, err
	}

	flag := false
	for _, b := range result.Buckets {
		if *b.Name == bucketname {
			flag = true
		}
	}

	return flag, nil
}


func DownloadFile(fname string, bucketname string, region string) (bool, error) {
	creds := credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, "")
	_, err := creds.Get()
	if err != nil {
		return false, err
	}

	cfg := aws.NewConfig().WithRegion(region).WithCredentials(creds)
	cfg.Region = aws.String(region)
	sess, err := session.NewSession(cfg)
	svc := s3manager.NewDownloader(sess)

	file, err := os.Create(fname)
	if err != nil {
		return false, err
	}
	defer file.Close()

	path := "/" + file.Name()

	_, err = svc.Download(file, &s3.GetObjectInput {
		Bucket: aws.String(bucketname),
		Key: aws.String(path),
	})

	if err != nil {
		return false, err
	}

	return true, nil	
}

