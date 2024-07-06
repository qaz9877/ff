package info

import (
	"bytes"
	"context"
	"errors"
	"ff/pojo"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// 分离视频流,并且裁剪一定时间  第一步
func TrimVideoMp4(inputFile string, duration int, outputFile string) error {
	//cmdStr := fmt.Sprintf(`ffmpeg  -i "%s"  -ss 0 -t %d -an -c:v libx264 -crf 18 -c:a aac -strict -2 -y %s`, inputFile, duration, outputFile)
	// 使用exec.Command执行命令字符串
	//fmt.Println(cmdStr)

	//ffmpeg -ss 0 -t 274  -accurate_seek -i "C:\Users\ASUS\Desktop\视频\1.盘点游戏史上最具影响力的 10大独立游戏！(Av1302304460,P1).mp4"    -an
	//-codec copy -crf 18   -y C:\Users\ASUS\Desktop\temp\output_video.mp4
	cmdArgs := []string{
		pojo.ResultDir + "\\bin\\ffmpeg.exe",
		"-ss", "0",
		"-t", strconv.Itoa(duration),
		"-i", inputFile,
		"-an",
		"-codec", "copy",
		"-crf", "18",
		// "-c:a", "aac",
		// "-strict", "-2",
		"-y", outputFile,
	}
	sr := ""
	for _, v := range cmdArgs {
		sr += v + " "
	}
	fmt.Println(sr)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, cmdArgs[0], cmdArgs[1:]...)
	// go func() {
	// 	<-ctx.Done()
	// 	if ctx.Err() == context.DeadlineExceeded {
	// 		fmt.Println("Context done due to timeout. Killing process...")
	// 		cmd.Process.Kill()
	// 	}
	// }()
	//fmt.Println("分离视频流,并且裁剪一定时间  第一步")
	// 执行命令并捕获输出
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	// output, err := cmd.CombinedOutput()
	// if err != nil {
	// 	// 打印标准错误输出
	// 	fmt.Println("FFmpeg command failed:", string(output))
	// 	return err
	// }
	return nil
}

// 合并mp3 并且裁剪   第二步
func MergeAndTrimMP3(inputFiles []string, outputFile string, duration int) error {
	// 构建 FFmpeg 命令
	// cmdArgs := fmt.Sprintf("-i %s ", inputFiles[0])
	cmdArgs := []string{}
	//cmdArgs := ""
	// "", "auto",
	// 	"-threads", "10",
	cmdArgs = append(cmdArgs, "-hwaccel")
	cmdArgs = append(cmdArgs, "auto")

	for _, inputFile := range inputFiles {
		//cmdArgs += fmt.Sprintf("-i \"%s\" ", inputFile)
		cmdArgs = append(cmdArgs, "-i")
		cmdArgs = append(cmdArgs, inputFile)
	}
	cmdArgs = append(cmdArgs, "-threads")
	cmdArgs = append(cmdArgs, "auto")
	cmdArgs = append(cmdArgs, "-filter_complex")
	cmdArgs = append(cmdArgs, fmt.Sprintf("concat=n=%d:v=0:a=1", len(inputFiles)))
	cmdArgs = append(cmdArgs, "-ss")
	cmdArgs = append(cmdArgs, "0")
	cmdArgs = append(cmdArgs, "-t")
	cmdArgs = append(cmdArgs, strconv.Itoa(duration))
	//	"-preset", "",  "-gpu", "1",
	cmdArgs = append(cmdArgs, "-gpu")
	cmdArgs = append(cmdArgs, "1")
	cmdArgs = append(cmdArgs, "-preset")
	cmdArgs = append(cmdArgs, "ultrafast")
	cmdArgs = append(cmdArgs, "-y")
	cmdArgs = append(cmdArgs, outputFile)
	//-i anullsrc=cl=stereo:r=44100:d=0.01 -filter_complex
	//cmdArgs += fmt.Sprintf(" -filter_complex concat=n=%d:v=0:a=1 -ss 0 -t %d -y \"%s\"", len(inputFiles)+1, duration, outputFile)
	//fmt.Println(cmdArgs)

	// 使用 exec.Command 执行命令字符串
	sr := ""
	for _, v := range cmdArgs {
		sr += v + " "
	}
	fmt.Println(sr)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, pojo.ResultDir+"\\bin\\ffmpeg.exe", cmdArgs...)

	// go func() {
	// 	<-ctx.Done()
	// 	if ctx.Err() == context.DeadlineExceeded {
	// 		fmt.Println("Context done due to timeout. Killing process...")
	// 		cmd.Process.Kill()
	// 	}
	// }()

	// 执行命令
	//fmt.Println("合并mp3 并且裁剪   第二步")
	_, err := cmd.CombinedOutput()
	//time.Sleep(45 * time.Second)
	if err != nil {
		return err
	}

	return nil
}

// 融合 视频  音频  音乐  音量比例  音乐比例  第三步
func MergeAndAdjustVolume(inputVideo string, inputAudio1 string, inputAudio2 string, volume1 float64, volume2 float64, outputFile string) error {
	// 构建 FFmpeg 命令
	// cmdArgs := fmt.Sprintf("-i \"%s\" -i \"%s\" -i \"%s\" -filter_complex "+
	// 	"'[0:v]scale=1280:720:force_original_aspect_ratio=decrease,pad=1280:720:(ow-iw)/2:(oh-ih)/2[v]; "+
	// 	"[1:a]volume=%.2f[a1]; [2:a]volume=%.2f[a2]; "+
	// 	"[a1][a2]amix=inputs=2:duration=longest[a]' "+
	// 	"-map [v] -map [a] -c:v libx264 -c:a aac -strict -2 -y \"%s\"", inputVideo, inputAudio1, inputAudio2, volume1, volume2, outputFile)

	//-force_key_frames "expr:gte(t,n_forced*2)" -filter_complex "[1:a]volume=0.5[a1];[2:a]volume=2[a2];[a1][a2]amix=inputs=2:duration=longest[a]"
	//-map 0:v -map "[a]" -c:v copy -c:a aac -strict -2 -strict -2  output.mp4
	fmt.Println("[1:a]volume=" + fmt.Sprintf("%.3f", volume1) + "[a1]; [2:a]volume=" + fmt.Sprintf("%.3f", volume2) + "[a2]")

	cmdArgs := []string{
		//-hwaccel_device 1
		"-hwaccel", "auto",

		"-i", inputVideo,
		"-i", inputAudio1, //音乐
		"-i", inputAudio2, // 配音
		"-threads", "auto",
		"-filter_complex",
		"[1:a]volume=" + fmt.Sprintf("%.3f", volume1) + "[a1];[2:a]volume=" + fmt.Sprintf("%.3f", volume2) + "[a2];" +
			"[a1][a2]amix=inputs=2:duration=longest[a]",
		"-map", "0:v",
		"-map", "[a]",
		"-c:v", "copy",
		"-preset", "ultrafast",
		"-c:a", "aac",
		"-gpu", "1",
		"-y", outputFile,
	}
	sr := ""
	for _, v := range cmdArgs {
		sr += v + " "
	}
	fmt.Println(sr)
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	// 使用 exec.Command 执行命令字符串
	cmd := exec.CommandContext(ctx, pojo.ResultDir+"\\bin\\ffmpeg.exe", cmdArgs...)
	// go func() {
	// 	<-ctx.Done()
	// 	if ctx.Err() == context.DeadlineExceeded {
	// 		fmt.Println("Context done due to timeout. Killing process...")
	// 		cmd.Process.Kill()
	// 	}
	// }()
	// str := ""
	// for _, same := range cmd.Args {
	// 	str += same + " "
	// }
	// fmt.Println(str)
	// 执行命令
	//fmt.Println("融合 视频  音频  音乐  音量比例  音乐比例  第三步")
	// 执行命令并捕获输出
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}

// 第四步
func AddTextToImage(inputImage, textFile, outputImage, fontFile string, fontsize int) error {
	// 构造ffmpeg命令参数

	// -frames:v 1 -update 1
	cmdArgs := []string{
		"-hwaccel", "auto",
		"-i", inputImage,
		"-threads", "auto",
		"-vf", fmt.Sprintf("drawtext=textfile='%s':fontcolor=yellow:fontsize=%d:fontfile='%s':x=(W-text_w)/2:y=(H-text_h)/2:shadowx=2:shadowy=2:shadowcolor=black", textFile, fontsize, fontFile),
		// "-frames:v", "1",
		// "-update", "1",
		// "-strict", "-2",
		"-gpu", "1",
		"-preset", "ultrafast",
		"-y", outputImage,
	}
	sr := ""
	for _, v := range cmdArgs {
		sr += v + " "
	}
	fmt.Println(sr)
	// ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	// defer cancel()
	cmd := exec.Command(pojo.ResultDir+"\\bin\\ffmpeg.exe", cmdArgs...)

	err := cmd.Start()
	if err != nil {
		fmt.Println("启动命令时出错:", err)
		return err
	}

	done := make(chan error, 1)

	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(2 * time.Second):
		fmt.Println("已等待1秒，仍然在执行中，程序将退出。")
		// 如果希望杀死正在执行的命令，可以使用以下代码
		// 注意：这里只是示例，实际中可能需要根据具体情况进行处理
		if err := cmd.Process.Kill(); err != nil {
			fmt.Println("杀死命令时出错:", err)
			return err
		}
	case err := <-done:
		if err != nil {
			fmt.Println("命令执行出错:", err)
			return err
		} else {
			fmt.Println("命令执行完成.")
			return nil
		}
	}
	return nil
}

// 第五
func AddAttachedPic(inputVideo, inputImage, outputVideo string) error {
	cmdArgs := []string{
		"-i", inputVideo,
		"-i", inputImage,
		"-map", "1",
		"-map", "0",
		"-c", "copy",
		"-disposition:0", "attached_pic",
		// "-strict", "-2",
		"-y", outputVideo,
	}
	sr := ""
	for _, v := range cmdArgs {
		sr += v + " "
	}
	fmt.Println(sr)

	cmd := exec.Command(pojo.ResultDir+"\\bin\\ffmpeg.exe", cmdArgs...)
	err := cmd.Start()
	if err != nil {
		fmt.Println("启动命令时出错:", err)
		return err
	}

	done := make(chan error, 1)

	go func() {
		//time.Sleep(30 * time.Second)
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(2 * time.Second):
		fmt.Println("已等待1秒，仍然在执行中，程序将退出。")
		// 如果希望杀死正在执行的命令，可以使用以下代码
		// 注意：这里只是示例，实际中可能需要根据具体情况进行处理
		if err := cmd.Process.Kill(); err != nil {
			fmt.Println("杀死命令时出错:", err)

		}
	case err := <-done:
		if err != nil {
			fmt.Println("命令执行出错:", err)
		} else {
			fmt.Println("命令执行完成.")
		}
	}
	return nil
}
func ffprobeCommand(inputFile string) *exec.Cmd {
	return exec.Command(pojo.ResultDir+"\\bin\\ffprobe.exe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height", "-of", "csv=s=x:p=0", inputFile)
}

// 第六
func GetVideoResolution(inputFile string) (width, height string, err error) {
	cmd := ffprobeCommand(inputFile)
	sr := ""
	for _, v := range cmd.Args {
		sr += v + " "
	}
	fmt.Println(sr)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return "", "", err
	}
	res := strings.Split(strings.TrimSpace(out.String()), "x")
	if len(res) != 2 {
		return "", "", errors.New("unable to parse video resolution")
	}
	return res[0], res[1], nil
}

// 第七步
func ffmpegCommand1(inputFile, outputFile, width, height string) *exec.Cmd {
	return exec.Command(pojo.ResultDir+"\\bin\\ffmpeg.exe", "-i", inputFile, "-vf", fmt.Sprintf("scale=%s:%s,setdar=4:3", width, height), "-strict", "-2", "-y", outputFile)
}

func ProcessImage(inputFile, outputFile, width, height string) error {
	cmd := ffmpegCommand1(inputFile, outputFile, width, height)

	sr := ""
	for _, v := range cmd.Args {
		sr += v + " "
	}
	fmt.Println(sr)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("执行命令时出错: %v", err)
	}
	return nil
}

// 第八步
func ffmpegCommand(inputVideo, inputImage, outputVideo string) *exec.Cmd {
	return exec.Command(pojo.ResultDir+"\\bin\\ffmpeg.exe", "-i", inputVideo, "-i", inputImage, "-filter_complex", "[0:v][1:v]overlay=0:0:enable='eq(n,0)'", "-map", "0:a", "-c:v", "libx264", "-preset", "medium", "-crf", "23", "-c:a", "copy", "-y", outputVideo)
}

func ProcessVideoWithImage(inputVideo, inputImage, outputVideo string) error {
	cmd := exec.Command(pojo.ResultDir+"\\bin\\ffmpeg.exe",
		"-hwaccel", "auto",

		"-i", inputVideo,
		"-i", inputImage,
		"-threads", "auto",
		"-filter_complex", "[0:v][1:v]overlay=0:0:enable='eq(n,0)'",
		//-c:v h264_nvenc
		// "-c:v", "h264_nvenc",
		// "-c:a", "copy",
		"-gpu", "1",
		"-preset", "ultrafast",
		"-y",
		outputVideo,
	)
	sr := ""
	for _, v := range cmd.Args {
		sr += v + " "
	}
	fmt.Println(sr)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// 第九步
func ProcessVideoWithImage2(inputVideo, inputmusic, inputaudio, inputImage, outputVideo string, duration int, volume1 float64, volume2 float64) error {
	/* ffmpeg -hwaccel auto -i a.mp4 -i output_audio.mp3 -i c.wav -i output_image2.jpg
	-threads auto -filter_complex
	"[0:v]trim=start=0:end=283,setpts=PTS-STARTPTS[v2];
	[v2][3:v]overlay=0:0:enable='eq(n,0)'[v3];
	[1:a]volume=0.5[a1];
	[2:a]volume=1.0[a2];
	[a1][a2]amix=inputs=2:duration=longest[a]"
	 -map "[v3]" -map "[a]"
	 -c:a aac
	 -q:v 18
	  -preset veryfast

	  -strict -2  -gpu 1 -preset ultrafast -y result_2.mp4
	*/
	cmd := exec.Command(pojo.ResultDir+"\\bin\\ffmpeg.exe",
		"-hwaccel", "auto",

		"-i", inputVideo,
		"-i", inputmusic,
		"-i", inputaudio,
		"-i", inputImage,
		"-threads", "auto",
		"-filter_complex",
		"[0:v]trim=start=0:end="+strconv.Itoa(duration)+",setpts=PTS-STARTPTS[v2];"+"[v2][3:v]overlay=0:0:enable='eq(n,0)'[v3];"+"[1:a]volume="+fmt.Sprintf("%.3f", volume1)+"[a1];[2:a]volume="+fmt.Sprintf("%.3f", volume2)+"[a2];"+"[a1][a2]amix=inputs=2:duration=longest[a]",
		"-map", "[v3]", "-map", "[a]",
		"-c:a", "aac",
		"-q:v", "18",
		"-gpu", "1",
		"-strict", "-2",
		"-preset", "veryfast",
		"-y", outputVideo,
	)
	sr := ""
	for _, v := range cmd.Args {
		sr += v + " "
	}
	fmt.Println(sr)
	err := cmd.Run()
	// fmt.Println("7777")
	// fmt.Println(string(out))
	if err != nil {
		return err
	}
	return nil
}

// mp3裁剪
func TrimMP3(inputFile string, duration int, outputFile string) error {
	// 拼接命令字符串
	cmdStr := fmt.Sprintf("ffmpeg -i '%s' -ss 0 -t %d -y %s", inputFile, duration, outputFile)

	// 使用exec.Command执行命令字符串
	cmd := exec.Command("cmd", "/C", cmdStr)

	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// 视频裁剪  可以有开始时间
func RunFFmpegCommand(inputFile, outputFile, startTime, duration string) error {
	cmd := exec.Command("ffmpeg", "-i", inputFile, "-ss", startTime, "-t", duration, "-c:v", "libx264", "-crf", "18", "-c:a", "aac", "-strict", "-2", outputFile)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// 生成 FFmpeg 命令字符串
func executeFFmpegCommand(cmdStr string) []string {
	// 使用 strings.Fields 将命令字符串分割为命令和参数
	args := strings.Fields(cmdStr)

	return args
}
