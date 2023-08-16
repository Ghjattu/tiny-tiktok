package oss

import (
	"os"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/joho/godotenv"
)

var (
	client *oss.Client
	bucket *oss.Bucket
)

func init() {
	godotenv.Load("../.env")

	bucketName := os.Getenv("OSS_BUCKET_NAME")
	endpoint := os.Getenv("OSS_ENDPOINT")
	accessKeyID := os.Getenv("OSS_ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("OSS_ACCESS_KEY_SECRET")
	var err error

	client, err = oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		panic(err)
	}
	bucket, err = client.Bucket(bucketName)
	if err != nil {
		panic(err)
	}
}

// UploadFile upload file to oss and set public read permission.
//
//	@param objectKey string
//	@param localFilePath string
//	@return error
func UploadFile(objectKey, localFilePath string) error {
	err := bucket.PutObjectFromFile(objectKey, localFilePath)
	if err != nil {
		return err
	}

	err = bucket.SetObjectACL(objectKey, oss.ACLPublicRead)
	if err != nil {
		bucket.DeleteObject(objectKey)
		return err
	}

	return nil
}
