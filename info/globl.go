package info

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type GlobleObj struct {
	ConfFilePath  string
	Audio_address string `json:"audio_address"`
	Video_address string `json:"video_address"`
	Music_address string `json:"music_address"`
	Image_Address string `json:"image_address"`
	Txt_Address   string `json:"txt_address"`
}

var GlobleObject *GlobleObj

// 判断一个文件是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 读取用户的配置文件
func (g *GlobleObj) Reload() {

	if confFileExists, _ := PathExists(g.ConfFilePath); confFileExists != true {
		fmt.Println("Config File ", g.ConfFilePath, " is not exist!!")
		return
	}

	data, err := ioutil.ReadFile(g.ConfFilePath)
	if err != nil {
		panic(err)
	}
	//将json数据解析到struct中
	err = json.Unmarshal(data, g)
	if err != nil {
		panic(err)
	}
}

func Gugai(address GlobleObj) error {
	if address.Audio_address != "" && address.Music_address != "" && address.Video_address != "" && address.Image_Address != "" && address.Txt_Address != "" {
		jsonData, err := json.Marshal(address)
		if err != nil {
			fmt.Println("序列化对象时出错:", err)
			return err
		}

		// 打开文件进行写入
		file, err := os.Create("./info/56.json")
		if err != nil {
			fmt.Println("创建文件时出错:", err)
			return err
		}
		defer file.Close()

		_, err = file.Write(jsonData)
		if err != nil {
			fmt.Println("写入文件时出错:", err)
			return err
		}
		return nil
	}
	return fmt.Errorf("修改错误")
}

func init() {
	GlobleObject = &GlobleObj{
		ConfFilePath:  "./info/56.json",
		Audio_address: `F:\视频专用\配音`,
		Video_address: `F:\视频专用\视频`,
		Music_address: `F:\视频专用\音乐`,
		Image_Address: `F:\视频专用\图片`,
		Txt_Address:   `F:\视频专用\主题`,
	}
	GlobleObject.Reload()
}
