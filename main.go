//go:generate goversioninfo -icon=icon.ico
package main

import (
	"ff/info"
	"ff/pojo"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/juju/ratelimit"
)

type TTTask struct {
	ID           int     `json:"id"`
	AudioAddress string  `json:"audio_address"`
	VideoAddress string  `json:"video_address"`
	AudioVolume  float64 `json:"audio_volume"`
	MusicAddress string  `json:"music_address"`
	MusicVolume  float64 `json:"music_volume"`
}

var ch = make(chan pojo.Task, 1024)

var clients = make(map[*websocket.Conn]bool) // 连接的客户端
var broadcast = make(chan Message)           // 从连接的客户端发送消息的广播通道
var broadcast2 = make(chan map[int]pojo.Task)

// 消息对象
type Message struct {
	TaskID   int `json:"taskId"`
	Progress int `json:"progress"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// go generate   go build -o 音频处理软件.exe
func main() {
	//设置环境变量
	//E:\ffmpeg-6.0-essentials_build\ffmpeg-6.0-essentials_build\bin
	//ffmpeg.exe
	info.SetECREnv()
	r := gin.Default()

	go func() {
		processTasks2()
	}()
	// 设置静态文件目录
	r.Use(RateLimitMiddleware(time.Second, 100, 100))
	fmt.Println(info.GlobleObject.ConfFilePath)
	fmt.Println(info.GlobleObject.Audio_address)
	fmt.Println(info.GlobleObject.Video_address)
	fmt.Println(info.GlobleObject.Image_Address)
	fmt.Println(info.GlobleObject.Txt_Address)
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")
	//超时处理
	r.Use(timeoutMiddleware())
	// 设置路由
	r.POST("/send-tasks", sendTasksHandler)
	r.GET("/check-task-status", checkTaskStatusHandler)
	r.GET("/ws-endpoint", handleConnections)
	go handleMessages()
	r.GET("/clear-task", clearTaskStatusHandler)

	// GET 请求，返回全局 Address 对象
	r.GET("/getAddress", getAddress)

	// GET 请求，返回全局 Address 对象
	r.GET("/getAddressinfo", getAddressinfo)
	r.GET("/getAddressinfo2", getAddressinfo2)

	// POST 请求，修改全局 Address 对象
	r.POST("/updateAddress", updateAddress)
	// 设置路由
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, welcome to my website!",
		})
	})

	// 实现 /index 接口，返回 index.html 页面
	r.GET("/index", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"title": "Gin Web Framework",
		})
	})
	// 启动服务
	fmt.Println("访问地址：http://localhost:8081/index")
	info.OpenBrowser("http://localhost:8081/index")
	// if err := r.Run(":8081"); err != nil {
	// 	log.Fatalf("Failed to start server: %v", err)
	// }
	server := &http.Server{
		Addr:         ":8081",          // 监听端口
		Handler:      r,                // 使用 Gin 路由
		ReadTimeout:  10 * time.Second, // 设置读取超时时间为 10 秒
		WriteTimeout: 10 * time.Second, // 设置写入超时时间为 10 秒
	}

	// 启动 HTTP 服务器
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		// 打印错误信息
		panic(err)
	}
}

func handleMessages() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop() // 在程序退出时停止定时器
	for {
		select {
		case <-ticker.C:
			// 在这里执行你想要定时执行的任务
			// 从广播通道接收消息
			targetMap := getAllTasks()
			// 向所有连接的客户端发送消息
			for client := range clients {
				err := client.WriteJSON(targetMap)
				if err != nil {
					log.Printf("error: %v", err)
					client.Close()
					delete(clients, client)
				}
			}
		}

	}
}

func handleConnections(c *gin.Context) {
	// 升级 HTTP 连接为 WebSocket 连接
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal(err)
	}
	// 确保在函数结束时关闭连接
	defer ws.Close()

	// 将新连接添加到 clients 映射
	clients[ws] = true
	log.Println("WebSocket 服务器启动，监听端口 8081...")
	// 循环监听客户端发送的消息
	for {
		var msg Message
		// 读取消息
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		// 将消息发送到广播通道
		broadcast <- msg
	}
}

func RateLimitMiddleware(fillInterval time.Duration, cap, quantum int64) gin.HandlerFunc {
	bucket := ratelimit.NewBucketWithQuantum(fillInterval, cap, quantum)
	return func(c *gin.Context) {
		if bucket.TakeAvailable(1) < 1 {
			c.String(http.StatusForbidden, "rate limit...")
			c.Abort()
			return
		}
		c.Next()
	}
}
func getAddressinfo(c *gin.Context) {
	arr, err := ListFilesInDirectory(info.GlobleObject.Audio_address)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, arr)
	}

}
func getAddressinfo2(c *gin.Context) {
	arr, err := ListFilesInDirectory(info.GlobleObject.Txt_Address)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, arr)
	}

}

func ListFilesInDirectory(dirPath string) ([]string, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}

	return fileNames, nil
}

func SetGlobalObjectValues(obj *info.GlobleObj) {
	rv := reflect.ValueOf(obj).Elem()
	rt := rv.Type()

	for i := 0; i < rv.NumField(); i++ {
		field := rt.Field(i)
		value := rv.Field(i).String()
		fmt.Println(value)
		if value != "" && IsDirectory(value) {
			switch field.Name {
			case "ConfFilePath":
				info.GlobleObject.ConfFilePath = value
			case "Audio_address":
				info.GlobleObject.Audio_address = value
			case "Video_address":
				info.GlobleObject.Video_address = value
			case "Music_address":
				info.GlobleObject.Music_address = value
			case "Image_Address":
				info.GlobleObject.Image_Address = value
			case "Txt_Address":
				info.GlobleObject.Txt_Address = value
			}
		}
	}
}

func IsDirectory(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		//fmt.Println("'sssssssssssss")
		return false
	}
	return fi.Mode().IsDir()
}

// getAddress 返回全局 Address 对象
func getAddress(c *gin.Context) {
	c.JSON(http.StatusOK, info.GlobleObject)
}

// updateAddress 修改全局 Address 对象
func updateAddress(c *gin.Context) {
	// 从请求体中解析 JSON 数据
	var newAddress info.GlobleObj
	if err := c.ShouldBindJSON(&newAddress); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		// 更新全局 Address 对象
		err := info.Gugai(newAddress)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			// 返回成功消息
			SetGlobalObjectValues(&newAddress)
			c.JSON(http.StatusOK, info.GlobleObject)
		}

	}

}
func sendTasksHandler(c *gin.Context) {
	// Simulate sending tasks to backend
	var tasks []pojo.Task
	//var taskd []TTTask
	// 解析前端发送的JSON数据
	if err := c.BindJSON(&tasks); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		// 打印接收到的任务数据
		//fmt.Println("接收到的任务数据：", tasks)

		// 将接收到的任务传递给处理函数
		processTasks(tasks)
		c.JSON(http.StatusOK, gin.H{"message": "Tasks sent successfully."})
	}

}

func checkTaskStatusHandler(c *gin.Context) {
	targetMap := getAllTasks()
	if targetMap != nil {
		c.JSON(http.StatusOK, targetMap)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "没有"})
	}

	// 执行耗时操作

}

func testResponse(c *gin.Context) {
	c.JSON(http.StatusGatewayTimeout, gin.H{
		"code":    http.StatusGatewayTimeout,
		"message": "timeout",
	})
}

func timeoutMiddleware() gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(3000*time.Millisecond),
		timeout.WithHandler(func(c *gin.Context) {
			c.Next()
		}),
		timeout.WithResponse(testResponse),
	)
}

// clearTaskStatusHandler 是清空任务状态的处理函数
func clearTaskStatusHandler(c *gin.Context) {
	// 清空任务状态
	clearMap()
	// 返回成功消息
	fmt.Println("清空任务状态")
	c.JSON(http.StatusOK, gin.H{"message": "任务状态已清空"})
}

// 处理任务的函数，您可以根据需要编写具体的逻辑
func processTasks(tasks []pojo.Task) {
	pojo.TasksLock.Lock()
	defer pojo.TasksLock.Unlock()
	for i := 0; i < len(tasks); i++ {
		tasks[i].Status = false
		tasks[i].Progress = 0
		tasks[i].AudioAddress = info.GlobleObject.Audio_address + "\\" + tasks[i].AudioAddress
		tasks[i].VideoAddress = info.GlobleObject.Video_address + "\\" + tasks[i].VideoAddress
		tasks[i].MusicAddress = info.GlobleObject.Music_address + "\\" + tasks[i].MusicAddress
		if tasks[i].InputImage != "" {
			tasks[i].InputImage = info.GlobleObject.Image_Address + "\\" + tasks[i].InputImage
			tasks[i].TextFile = info.GlobleObject.Txt_Address + "\\" + tasks[i].TextFile
		}
		pojo.Tasks[tasks[i].ID] = tasks[i]

		ch <- tasks[i]
		fmt.Println("状态更新")

	}

	// for i := 0; i < len(tasks); i++ {
	// 	// _, task := range tasks {
	// 	fmt.Println("任务id:", tasks[i].ID)
	// 	// 在这里编写处理任务的逻辑
	// 	// 您可以访问任务的各个字段，如 task.AudioAddress、task.VideoAddress 等
	// 	// 示例：fmt.Printf("处理任务 ID：%d\n", task.ID)
	// 	//fmt.Printf("处理任务 ID：%d\n", task.ID)

	// 	err := info.Work(tasks[i].VideoAddress, tasks[i].MusicAddress, tasks[i].AudioAddress, tasks[i].MusicVolume, tasks[i].AudioVolume)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		tasks[i].Status = false
	// 		tasks[i].Progress = 2
	// 		pojo.Tasks[tasks[i].ID] = &tasks[i]
	// 		continue
	// 	}
	// 	tasks[i].Status = true
	// 	tasks[i].Progress = 1

	// 	pojo.Tasks[tasks[i].ID] = &tasks[i]
	// 	fmt.Println(pojo.Tasks)
	// }
}

func processTasks2() {
	for {
		select {
		case v := <-ch:
			//  ... foo
			// _, task := range tasks {
			//fmt.Println("任务id:", v.ID)
			// 在这里编写处理任务的逻辑
			// 您可以访问任务的各个字段，如 task.AudioAddress、task.VideoAddress 等
			// 示例：fmt.Printf("处理任务 ID：%d\n", task.ID)
			fmt.Printf("处理任务")
			fmt.Printf("处理任务 ID：%d\n", v.ID)
			var err error
			if v.FontFile == "" || v.InputImage == "" || v.TextFile == "" {
				err = info.Work(v.VideoAddress, v.MusicAddress, v.AudioAddress, v.MusicVolume, v.AudioVolume)
			} else {
				//fmt.Println("进入2")
				//fmt.Println(v)
				err = info.Work2(v.VideoAddress, v.MusicAddress, v.AudioAddress, v.MusicVolume, v.AudioVolume, v.InputImage,
					v.TextFile, v.FontFile, v.Fontsize)
			}
			ddd(v, err)
		default:
			time.Sleep(1 * time.Second)
		}
	}
}

func ddd(v pojo.Task, err error) {
	if pojo.TasksLock.TryLock() {
		defer pojo.TasksLock.Unlock()
		if err != nil {
			fmt.Println(err)
			v.Status = false
			v.Progress = 2
			pojo.Tasks[v.ID] = v

		} else {
			v.Status = true
			v.Progress = 1
			//pojo.TasksLock.Lock()
			pojo.Tasks[v.ID] = v
		}
		fmt.Println(pojo.Tasks)
	} else {
		fmt.Println("有锁")
	}

}
func getAllTasks2() map[int]pojo.Task {

	targetMap := make(map[int]pojo.Task)
	fmt.Println("任务数据长度：", len(targetMap))
	// 遍历 tasks map，将任务信息添加到切片中
	pojo.TasksLock.RLock()
	// pojo.TasksLock.RUnlock()
	// pojo.TasksLock.RLock()
	for key, value := range pojo.Tasks {
		targetMap[key] = value
	}
	pojo.TasksLock.RUnlock()
	//fmt.Println("1任务数据：", targetMap)

	return targetMap
}

func getAllTasks() map[int]pojo.Task {
	//fmt.Println("任务数据长度：", "555555555555555")
	var targetMap = make(map[int]pojo.Task)
	// 遍历 tasks map，将任务信息添加到切片中
	if pojo.TasksLock.TryRLock() {
		// pojo.TasksLock.RUnlock()
		// pojo.TasksLock.RLock()
		for key, value := range pojo.Tasks {
			targetMap[key] = value
		}
		defer pojo.TasksLock.RUnlock()
		//fmt.Println("1任务数据：", targetMap)
		fmt.Println("任务数据长度：", len(targetMap))
	}
	return targetMap
}

// clearMap 函数接收一个map作为参数，并清空该map
func clearMap() {
	pojo.TasksLock.Lock()
	for key := range pojo.Tasks {
		delete(pojo.Tasks, key)
	}
	defer pojo.TasksLock.Unlock()

}

func m4() error {
	video_addres := filepath.ToSlash(`C:\Users\ASUS\Desktop\视频\1.盘点游戏史上最具影响力的 10大独立游戏！(Av1302304460,P1).mp4`)
	music_addres := filepath.ToSlash(`C:\Users\ASUS\Desktop\音乐\1320525681.wma`)
	audio_addres := filepath.ToSlash(`C:\Users\ASUS\Desktop\音频\1.我只唱一句，你们自己找差距，王者素人被团灭的7大神级现场(Av1752382067,P1).mp3`)
	volume1 := 0.3
	volume2 := 2.2
	err := info.Work(video_addres, music_addres, audio_addres, volume1, volume2)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// 执行 FFmpeg 命令

// func m1() {
// 	// 打开 MP3 文件
// 	file, err := os.Open("your_file.mp3")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer file.Close()

// 	// 解析文件
// 	metadata, err := tag.ReadFrom(file)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// 获取时长
// 	duration := metadata.Duration()

// 	fmt.Println("Duration:", duration)
// }

// func m3() {
// 	videos := []string{"input1.mp4", "input2.mp4", "input3.mp4"}
// 	output := "output4.mp4"

// 	err := info.ConcatVideos(output, videos)
// 	if err != nil {
// 		fmt.Printf("Error: %s", err.Error())
// 		return
// 	}

// 	fmt.Println("Video concatenation completed successfully!")
// }

// 超时控制中间件

func timedHandler(duration time.Duration) func(c *gin.Context) {
	return func(c *gin.Context) {
		// 获取替换之后的context 它具备了超时控制
		ctx := c.Request.Context()
		// 定义响应struct
		type responseData struct {
			status int
			body   map[string]interface{}
		}
		// 创建一个done chan表明request要完成了
		doneChan := make(chan responseData)
		// 模拟API耗时的处理
		go func() {
			time.Sleep(duration)
			doneChan <- responseData{
				status: 200,
				body:   gin.H{"hello": "world"},
			}
		}()
		// 监听两个chan谁先到达
		select {
		// 超时
		case <-ctx.Done():
			return
			// 请求完成
		case res := <-doneChan:
			c.JSON(res.status, res.body)
		}
	}
}
