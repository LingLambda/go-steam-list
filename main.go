package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Response struct {
	Applist Applist `json:"applist"`
}

type Applist struct {
	Apps []App `json:"apps"`
}

type App struct {
	Appid int    `json:"appid"`
	Name  string `json:"name"`
}

const FILE_NAME = "app.json"
const INTERVAL = 120 * time.Second

var is_init = false

func main() {
	log.Printf("正在进行首次运行...")
	update_json() //启动时运行一次

	log.Printf("启动定时更新 间隔%v\n", INTERVAL)
	ticker := time.NewTicker(INTERVAL)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				update_json()
			}
		}
	}()
	select {} // 阻塞，程序不会退出，直到外部终止
}

func update_json() {
	local_apps, err := read_json()
	if err != nil {
		log.Println(err)
		if is_init {
			panic("已初始化但仍无法读取 app.json ，异常退出")
		}
		log.Println("执行初始化函数...")
		if !init_json() {
			panic("初始化失败!")
		}
	}
	local_map := make(map[int]string)
	for _, app := range local_apps {
		local_map[app.Appid] = app.Name
	}
	log.Println("请求最新 applist...")
	req_apps, err := request_json()
	if err != nil {
		log.Println(err)
	}
	log.Println("对比请求与本地文件...")
	var new_app []App
	has_new := false
	for _, app := range req_apps {
		_, ok := local_map[app.Appid]
		if !ok {
			has_new = true
			new_app = append(new_app, app)
		}
	}

	if has_new {
		new_app_count := len(new_app)
		local_apps = append(local_apps, new_app...)
		log.Printf("新增 %d 个应用\n", new_app_count)
		write_json(local_apps)
	} else {
		log.Printf("太好了，没有新应用")
	}
}

func init_json() bool {
	apps, err := request_json()
	if err != nil {
		log.Println("请求 JSON 失败：", err)
		return false
	}
	write_json(apps)
	is_init = true
	return true
}

func request_json() ([]App, error) {
	resp, err := http.Get("https://api.steampowered.com/ISteamApps/GetAppList/v2/")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Println("Response status:", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
	}
	apps := response.Applist.Apps
	return apps, nil
}

func write_json(apps []App) {

	// 创建文件（会覆盖原文件）
	file, err := os.Create(FILE_NAME)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 使用 json.NewEncoder 将结构体写入文件
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // 可选：格式化输出
	err = encoder.Encode(apps)
	if err != nil {
		panic(err)
	}
	log.Printf("已将数据写入到文件 %v\n", FILE_NAME)
}

func read_json() ([]App, error) {
	file, err := os.Open(FILE_NAME)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	// 解码 JSON
	var apps []App
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&apps)
	if err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %v", err)
	}
	return apps, nil
}
