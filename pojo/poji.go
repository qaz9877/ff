package pojo

import (
	"os"
	"sync"
)

var ResultDir, _ = os.Getwd()
var (
	Tasks     = make(map[int]Task)
	TasksLock = &sync.RWMutex{}

	// TasksLock2 sync.Mutex
)

type Sumupinfo struct {
	ID        int
	Mp3       MP3info
	Mp4       MP4info
	Audioinfo audioinfo
}

type MP3info struct {
	time    int
	name    string
	address string
}

type MP4info struct {
	time         int
	name         string
	address      string
	output_video string
}

type audioinfo struct {
	time    int
	name    string
	address string
}

// Task represents a single task
type Task struct {
	ID           int `json:"id"`
	Status       bool
	AudioAddress string  `json:"audio_address"`
	VideoAddress string  `json:"video_address"`
	AudioVolume  float64 `json:"audio_volume"`
	MusicAddress string  `json:"music_address"`
	MusicVolume  float64 `json:"music_volume"`
	Progress     int     `json:"progress"`
	InputImage   string  `json:"input_image"`
	TextFile     string  `json:"text_file"`
	FontFile     string  `json:"font_file"`
	Fontsize     int     `json:"font_size"`
}

type Address struct {
	Audio_address string `json:"audio_address"`
	Video_address string `json:"video_address"`
	Music_address string `json:"music_address"`
	ImageAddress  string `json:"image_address"`
	TtfAddress    string `json:"ttf_address"`
}

// TaskStatus 结构体用于表示任务状态
type TaskStatus struct {
	ID     int  `json:"id"`
	Status bool `json:"status"`
}
