package oss

import (
	"net/http"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var (
	endpoint   string
	bucketName string
)

func init() {
	godotenv.Load("../.env")

	endpoint = os.Getenv("OSS_ENDPOINT")
	bucketName = os.Getenv("OSS_BUCKET_NAME")
}

func TestUploadFile(t *testing.T) {
	t.Run("upload file failed", func(t *testing.T) {
		err := UploadFile("test_bear.mp4", "../data/fake.mp4")

		assert.NotNil(t, err)
	})

	t.Run("upload file successfully", func(t *testing.T) {
		err := UploadFile("test_bear.mp4", "../data/bear.mp4")
		if err != nil {
			t.Fatalf("upload file failed, err:%v\n", err)
		}

		url := "https://" + bucketName + "." + endpoint + "/test_bear.mp4"
		res, err := http.Get(url)
		if err != nil {
			t.Fatalf("http request failed, err:%v\n", err)
		}
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, "video/mp4", res.Header.Get("Content-Type"))
	})
}
