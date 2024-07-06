package info

import (
	"bufio"
	"ff/pojo"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"errors"
)

// func main() {
// 	url := "http://localhost:6666"
// 	openBrowser(url)
// }

func OpenBrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return err
}

// func main() {
// 	// 调用设置环境变量的函数
// 	err := setECREnv()
// 	if err != nil {
// 		fmt.Println("Error setting ECR environment variable:", err)
// 		return
// 	}
// 	fmt.Println("ECR environment variable set successfully!")
// }

func SetECREnv() error {
	// 设置环境变量

	ffmpegdir := pojo.ResultDir + "\\bin"
	ffmpegdir = ffmpegdir
	//fmt.Println("sssssssssssssssss:" + ffmpegdir)
	err := os.Setenv("ECR", ffmpegdir)
	if err != nil {
		fmt.Println("Error setting environment variable:", err)
		return err
	}
	return nil
}

func ExtractFileName(filePath string) (string, string) {
	// 使用 filepath 包来提取文件名
	fileName := filepath.Base(filePath)
	//fmt.Println(fileName)
	// 从文件名中提取不带扩展名的部分
	extension := filepath.Ext(fileName)
	//fmt.Println(extension)
	nameWithoutExt := strings.TrimSuffix(fileName, extension)

	return nameWithoutExt, extension
}

func Work(video_addres, music_addres, audio_addres string, volume1, volume2 float64) error {

	output_video_temp := filepath.ToSlash(pojo.ResultDir + `\temp\output_video.mp4`)
	output_audio_temp := filepath.ToSlash(pojo.ResultDir + `\temp\output_audio.mp3`)
	_, extension := ExtractFileName(video_addres)
	name, _ := ExtractFileName(audio_addres)
	result_name := filepath.ToSlash(pojo.ResultDir + `\result\` + name + extension)
	fmt.Println("result", result_name)
	//判断文件是否存在
	err := FileIfNotExist(video_addres)
	if err != nil {
		return err
	}
	err = FileIfNotExist(music_addres)
	if err != nil {
		return err
	}
	err = FileIfNotExist(audio_addres)
	if err != nil {
		return err
	}

	mp3_duration, err := GetMP4Duration(music_addres)
	fmt.Printf("mp3时长为: %d 秒\n", mp3_duration)
	if err != nil {
		return err
	}

	audio_duration, err := GetMP4Duration(audio_addres)
	fmt.Printf("音频时长为: %d 秒\n", audio_duration)
	if err != nil {
		//fmt.Println(err)
		return err
	}

	//video_duration, err := GetMP4Duration2(video_addres)
	// video_duration, err := GetMP4Duration(video_addres)
	// fmt.Printf("视频时长为: %d 秒\n", video_duration)
	// if err != nil {
	// 	//fmt.Println(err)
	// 	return err
	// }

	if mp3_duration == 0 || audio_duration == 0 {
		return errors.New("时间获取错误")
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
		//fmt.Println(remainder)
		for i := 0; i < remainder; i++ {
			inputFiles = append(inputFiles, music_addres)
		}
		//fmt.Println(inputFiles)
	}

	fmt.Println("第一步")
	// //第一步
	err = TrimVideoMp4(video_addres, audio_duration, output_video_temp)
	if err != nil {
		//fmt.Println(err)
		return err
	}

	fmt.Println("第二步")
	// //第二步
	err = MergeAndTrimMP3(inputFiles, output_audio_temp, audio_duration)
	if err != nil {
		//fmt.Println(err)
		return err
	}

	//第三步
	err = MergeAndAdjustVolume(output_video_temp, output_audio_temp, audio_addres, volume1, volume2, result_name)
	if err != nil {
		//fmt.Println(err)
		return err
	}
	fmt.Println("success:     ", output_video_temp)
	return nil
}

func Work2(video_addres, music_addres, audio_addres string, volume1, volume2 float64, inputImage, textFile, fontFile string, fontsize int) error {

	output_video_temp := filepath.ToSlash(pojo.ResultDir + `\temp\output_video.mp4`)
	//	output_video_temp2 := filepath.ToSlash(pojo.ResultDir + `\temp\output2_video.mp4`)
	output_audio_temp := filepath.ToSlash(pojo.ResultDir + `\temp\output_audio.mp3`)
	outputImage := filepath.ToSlash(pojo.ResultDir + `\temp\output_image.jpg`)
	outputImage2 := filepath.ToSlash(pojo.ResultDir + `\temp\output_image2.jpg`)
	_, extension := ExtractFileName(video_addres)
	name, _ := ExtractFileName(audio_addres)
	result_name := filepath.ToSlash(pojo.ResultDir + `\result\` + name + extension)
	fontFile = "./ttf/" + fontFile
	textFile = strings.Replace(textFile, "\\", "\\\\", -1)
	textFile = strings.Replace(textFile, ":", "\\:", -1)
	fontFile = strings.Replace(fontFile, "\\", "\\\\", -1)
	fontFile = strings.Replace(fontFile, ":", "\\:", -1)
	fmt.Println("result", result_name)
	//判断文件是否存在
	err := FileIfNotExist(video_addres)
	if err != nil {
		return err
	}
	//ttf G:\信任文件\代码\ff\ttf
	err = FileIfNotExist(music_addres)
	if err != nil {
		return err
	}
	err = FileIfNotExist(audio_addres)
	if err != nil {
		return err
	}

	mp3_duration, err := GetMP4Duration(music_addres)
	fmt.Printf("mp3时长为: %d 秒\n", mp3_duration)
	if err != nil {
		return err
	}

	audio_duration, err := GetMP4Duration(audio_addres)
	fmt.Printf("音频时长为: %d 秒\n", audio_duration)
	if err != nil {
		//fmt.Println(err)
		return err
	}

	//video_duration, err := GetMP4Duration2(video_addres)
	// video_duration, err := GetMP4Duration(video_addres)
	// fmt.Printf("视频时长为: %d 秒\n", video_duration)
	// if err != nil {
	// 	//fmt.Println(err)
	// 	return err
	// }

	if mp3_duration == 0 || audio_duration == 0 {
		return errors.New("时间获取错误")
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
		//fmt.Println(remainder)
		for i := 0; i < remainder; i++ {
			inputFiles = append(inputFiles, music_addres)
		}
		//fmt.Println(inputFiles)
	}

	// fmt.Println("第一步")
	// //第一步
	// err = TrimVideoMp4(video_addres, audio_duration, output_video_temp)
	// if err != nil {
	// 	//fmt.Println(err)
	// 	return err
	// }

	fmt.Println("第二步")
	// //第二步
	err = MergeAndTrimMP3(inputFiles, output_audio_temp, audio_duration)
	if err != nil {
		//fmt.Println(err)
		return err
	}
	// fmt.Println("第三步")
	// //第三步
	// err = MergeAndAdjustVolume(output_video_temp, output_audio_temp, audio_addres, volume1, volume2, output_video_temp2)
	// if err != nil {
	// 	//fmt.Println(err)
	// 	return err
	// }
	fmt.Println("第四步")
	//第四步
	err = AddTextToImage(inputImage, textFile, outputImage, fontFile, fontsize)
	if err != nil {
		return err
	}
	fmt.Println("第五步")
	//第五步
	// err = AddAttachedPic(output_video_temp2, outputImage, result_name)
	// if err != nil {
	// 	return err
	// }
	fmt.Println("第六步")
	width, height, err := GetVideoResolution(video_addres)
	fmt.Printf("宽高 %s %s", width, height)
	if err != nil {
		return err
	}
	fmt.Println("第七步")
	err = ProcessImage(outputImage, outputImage2, width, height)
	if err != nil {
		return err
	}
	// fmt.Println("第八步")
	// err = ProcessVideoWithImage(output_video_temp2, outputImage2, result_name)
	// if err != nil {
	// 	return err
	// }
	fmt.Println("第九步")
	err = ProcessVideoWithImage2(video_addres, output_audio_temp, audio_addres, outputImage2, result_name, audio_duration, volume1, volume2)
	if err != nil {
		return err
	}

	fmt.Println("success:     ", output_video_temp)
	return nil
}

// CreateDirectoryIfNotExist 函数用于检查指定路径的目录是否存在，如果不存在则创建。
func CreateDirectoryIfNotExist(dirPath string) error {
	// 判断目录是否存在
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		// 如果目录不存在，则创建目录
		err := os.Mkdir(dirPath, 0755)
		if err != nil {
			return fmt.Errorf("创建目录失败: %v", err)
		}
		fmt.Println("目录已创建:", dirPath)
	} else {
		fmt.Println("目录已存在:", dirPath)
	}
	return nil
}

func FileIfNotExist(filePath string) error {
	// 判断文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("文件不存在: %s", filePath)
	}
	return nil
}

// CreateIfNotExist 创建文件，如果文件不存在的话
func CreateIfNotExist(filename string) (*os.File, error) {
	// 检查文件是否存在
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// 文件不存在，创建文件
		fmt.Println("文件不存在，正在创建...")
		file, err := os.Create(filename)
		if err != nil {
			return nil, fmt.Errorf("创建文件时出错: %v", err)
		}
		fmt.Println("文件创建成功.")
		return file, nil
	} else if err != nil {
		// 其他错误
		return nil, fmt.Errorf("检查文件时出错: %v", err)
	}

	// 文件已经存在，打开文件并返回
	fmt.Println("文件已存在.")
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("打开文件时出错: %v", err)
	}
	return file, nil
}

// ReadFileByLines 逐行读取文件并返回一个包含所有行的切片
func ReadFileByLines(filename string) ([]string, error) {
	// 打开文件
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("打开文件时出错: %v", err)
	}
	defer file.Close()

	// 创建一个切片用于存储文件的所有行
	lines := make([]string, 0)

	// 创建 Scanner 对象来逐行读取文件内容
	scanner := bufio.NewScanner(file)

	// 逐行读取文件内容并添加到切片中
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// 检查是否有读取错误
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取文件时出错: %v", err)
	}

	return lines, nil
}
