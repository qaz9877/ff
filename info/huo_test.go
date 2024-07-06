package info

import (
	"fmt"
	"testing"
)

func TestGetMP3PlayDuration(t *testing.T) {
	filePath := "G:\\信任文件\\代码\\ff\\新建文件夹\\歌曲.mp3" // Replace this with the actual file path
	duration, err := GetMP3Duration33(filePath)
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Printf("Duration: %.2fseconds\n", duration.Seconds())
	fmt.Printf("音频时长为: %d 秒\n", duration)
}

func TestGetMP4Duration(t *testing.T) {
	filePath := `G:\信任文件\代码\ff\新建文件夹\8.mp4` // Replace this with the actual file path
	duration, err := GetMP4Duration(filePath)
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Printf("Duration: %.2fseconds\n", duration.Seconds())
	fmt.Printf("视频时长为: %d 秒\n", duration)
}

func TestGetMP4Duration2(t *testing.T) {
	filePath := `G:\信任文件\代码\ff\新建文件夹\8.mp4` // Replace this with the actual file path
	duration, err := GetMP4Duration2(filePath)
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Printf("Duration: %.2fseconds\n", duration.Seconds())
	fmt.Printf("视频时长为: %d 秒\n", duration)
}
