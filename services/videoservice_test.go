package services

import (
	"testing"

	"github.com/Ghjattu/tiny-tiktok/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewVideoWithEmptyTitle(t *testing.T) {
	models.InitDatabase(true)

	vs := &VideoService{}

	status_code, statue_msg := vs.CreateNewVideo("test", "", 1, "test")

	assert.Equal(t, int32(1), status_code)
	assert.Equal(t, "video title is empty", statue_msg)
}

func TestCreateNewVideoWithCorrectVideo(t *testing.T) {
	models.InitDatabase(true)

	vs := &VideoService{}

	status_code, statue_msg := vs.CreateNewVideo("test", "test", 1, "test")

	assert.Equal(t, int32(0), status_code)
	assert.Equal(t, "create new video successfully", statue_msg)
}
