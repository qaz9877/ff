package info

import (
	"context"
	"ff/pojo"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/tcolgate/mp3"
)

func GetMP3Duration33(filePath string) (int, error) {
	// 打开 MP3 文件
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// 创建解码器
	d := mp3.NewDecoder(file)

	// 计算时长
	var totalDuration time.Duration
	skipped := 0
	for {
		// 解码一帧
		frame := mp3.Frame{}
		if err := d.Decode(&frame, &skipped); err != nil {
			// 到达文件结尾
			break
		}

		// 累加帧时长
		totalDuration += frame.Duration()
	}
	durationSec := totalDuration.Seconds()
	durationInt := int(durationSec)
	return durationInt, nil
}
func GetMP4Duration(filePath string) (int, error) {
	// 构建 FFmpeg 命令
	ctx, cancel := context.WithTimeout(context.Background(), 70*time.Second)
	defer cancel()
	//-threads 5  -preset veryfast
	//ffmpeg  -i "1.成都“不用咬全靠吸”的蹄花，裹满辣椒水太满足了(Av1052272125,P1).mp4" -threads 5 -preset veryfast -f null null
	//ffmpeg  -i "1.成都“不用咬全靠吸”的蹄花，裹满辣椒水太满足了(Av1052272125,P1).mp4" -threads 5 -preset veryfast  -f null null
	cmd := exec.CommandContext(ctx, pojo.ResultDir+"\\bin\\ffmpeg.exe", "-i", filePath, "-f", "null", "NUL")

	// go func() {
	// 	<-ctx.Done()
	// 	if ctx.Err() == context.DeadlineExceeded {
	// 		fmt.Println("Context done due to timeout. Killing process...")
	// 		cmd.Process.Signal(syscall.SIGINT)
	// 	}
	// }()
	// 执行命令并捕获输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}

	// 解析输出以获取时长
	durationStr := extractDuration(string(output))
	if durationStr == "" {
		return 0, fmt.Errorf("Failed to extract duration from FFmpeg output")
	}

	// 将时长字符串转换为秒数
	duration, err := parseDuration(durationStr)
	if err != nil {
		return 0, err
	}

	return int(duration), nil
}

// 从 FFmpeg 输出中提取时长信息
func extractDuration(ffmpegOutput string) string {
	lines := strings.Split(ffmpegOutput, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Duration:") {
			parts := strings.Split(line, ",")
			durationPart := strings.TrimSpace(parts[0])
			duration := strings.TrimPrefix(durationPart, "Duration: ")
			return duration
		}
	}
	return ""
}

// 将时长字符串解析为秒数
func parseDuration(durationStr string) (int64, error) {
	parts := strings.Split(durationStr, ":")
	hours := parseDurationPart(parts[0])
	minutes := parseDurationPart(parts[1])
	seconds := parseDurationPart(parts[2])
	totalSeconds := hours*3600 + minutes*60 + seconds
	return int64(totalSeconds), nil
}

// 解析时长部分并转换为整数
func parseDurationPart(part string) int {
	var value int
	fmt.Sscanf(part, "%d", &value)
	return value
}
