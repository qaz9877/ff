package info

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

func TestTrimMP3(t *testing.T) {
	//go test -v -run TestTrimMP3

	//go test -v trimMP3_test.go

	inputFile := "歌曲.mp3"
	duration := 60
	outputFile := "output.mp3"

	err := TrimMP3(inputFile, duration, outputFile)
	if err != nil {
		t.Errorf("trimMP3(%s, %d, %s) returned error: %v", inputFile, duration, outputFile, err)
	}
}

// func TestMergeAndTrimMP3(t *testing.T) {
// 	inputFile := "歌曲1.mp3"
// 	outputFile := "output2.mp3"
// 	duration := 240

// 	err := MergeAndTrimMP3(inputFile, outputFile, duration)
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return
// 	}

// 	fmt.Println("MP3 files merged and trimmed successfully!")

// }

func TestTrimVideoMp4(t *testing.T) {
	inputFile := `C:\Users\ASUS\Desktop\视频\1.盘点游戏史上最具影响力的 10大独立游戏！(Av1302304460,P1).mp4`
	duration := 15
	outputFile := `C:\Users\ASUS\Desktop\result\output.mp4`

	err := TrimVideoMp4(inputFile, duration, outputFile)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Video trimmed successfully!")
}

func TestTask(t *testing.T) {
	video_addres := filepath.ToSlash(`C:\Users\ASUS\Desktop\视频\1.盘点游戏史上最具影响力的 10大独立游戏！(Av1302304460,P1).mp4`)
	music_addres := filepath.ToSlash(`C:\Users\ASUS\Desktop\音乐\1320525681.wma`)
	audio_addres := filepath.ToSlash(`C:\Users\ASUS\Desktop\音频\1.我只唱一句，你们自己找差距，王者素人被团灭的7大神级现场(Av1752382067,P1).mp3`)
	output_video_temp := filepath.ToSlash(`C:\Users\ASUS\Desktop\temp\output_video.mp4`)
	output_audio_temp := filepath.ToSlash(`C:\Users\ASUS\Desktop\temp\output_audio.mp3`)
	name, extension := ExtractFileName(video_addres)
	result_name := filepath.ToSlash(`C:\Users\ASUS\Desktop\result\` + name + extension)
	volume1 := 0.6
	volume2 := 2.2
	mp3_duration, err := GetMP4Duration(music_addres)
	fmt.Printf("mp3时长为: %d 秒\n", mp3_duration)
	if err != nil {
		fmt.Println(err)
	}

	audio_duration, err := GetMP3Duration33(audio_addres)
	fmt.Printf("音频时长为: %d 秒\n", audio_duration)
	if err != nil {
		fmt.Println(err)
	}

	//video_duration, err := GetMP4Duration2(video_addres)
	video_duration, err := GetMP4Duration2(video_addres)
	fmt.Printf("视频时长为: %d 秒\n", video_duration)
	if err != nil {
		fmt.Println(err)
	}

	//求inputFiles中的个数
	inputFiles := []string{}
	if audio_duration < mp3_duration {
		for i := 0; i < 2; i++ {
			inputFiles = append(inputFiles, music_addres)
		}
	} else {
		remainder := audio_duration / mp3_duration
		remainder++
		for i := 0; i < remainder; i++ {
			inputFiles = append(inputFiles, music_addres)
		}
	}

	//第一步
	// err = TrimVideoMp4(video_addres, audio_duration, output_video_temp)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	//第二步
	err = MergeAndTrimMP3(inputFiles, output_audio_temp, audio_duration)
	if err != nil {
		fmt.Println(err)
	}

	//第三步
	err = MergeAndAdjustVolume(output_video_temp, output_audio_temp, audio_addres, volume1, volume2, result_name)
	if err != nil {
		fmt.Println(err)
	}
	// "33:47" "36"
	fmt.Println("结束")
}

func TestAddTextToImage(t *testing.T) {
	inputImage := `C:\Users\ASUS\Desktop\temp\芙莉莲.jpg`
	textFile := `C:\Users\ASUS\Desktop\temp\text.txt`
	textFile = strings.Replace(textFile, "\\", "\\\\", -1)
	textFile = strings.Replace(textFile, ":", "\\:", -1)

	// fmt.Println(textFile)
	// textFile = `C\:\\Users\\ASUS\\Desktop\\temp\\text.txt`
	// fmt.Println(textFile)
	outputImage := `C:\Users\ASUS\Desktop\temp\out3.jpg`
	fontFile := `C:\Users\ASUS\Desktop\temp\YaHeiMonacoHybrid.ttf`
	fontFile = strings.Replace(fontFile, "\\", "\\\\", -1)
	fontFile = strings.Replace(fontFile, ":", "\\:", -1)
	err := AddTextToImage(inputImage, textFile, outputImage, fontFile, 50)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("image trimmed successfully!")
}

func TestAddAttachedPic(t *testing.T) {
	inputVideo := `C:\Users\ASUS\Desktop\temp\66601.mp4`
	inputImage := `C:\Users\ASUS\Desktop\temp\out3.jpg`
	outputVideo := `C:\Users\ASUS\Desktop\temp\test_output.mp4`

	// 测试执行命令
	err := AddAttachedPic(inputVideo, inputImage, outputVideo)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// // 检查输出文件是否存在
	// if _, err := os.Stat(outputVideo); os.IsNotExist(err) {
	// 	fmt.Println.Errorf("output video file does not exist: %v", err)
	// }

	// // 清理测试生成的输出文件
	// err = os.Remove(outputVideo)
	// if err != nil {
	// 	t.Errorf("error cleaning up test output file: %v", err)
	// }
}
