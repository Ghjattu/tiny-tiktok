package ffmpeg

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/disintegration/imaging"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

// GetSnapshot get snapshot from video.
//
//	@param videoPath string
//	@param coverFileName string
//	@param frameNum int
//	@return error
func GetSnapshot(videoPath, coverFileName string, frameNum int) error {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(videoPath).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		log.Println("generate snapshot failed: ", err.Error())
		return err
	}

	img, err := imaging.Decode(buf)
	if err != nil {
		log.Println("decode snapshot failed: ", err.Error())
		return err
	}

	err = imaging.Save(img, coverFileName)
	if err != nil {
		log.Println("save snapshot failed: ", err.Error())
		return err
	}

	return nil
}
