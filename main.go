package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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

const file_name = "app.json"

func main() {
	// apps := request_json()
	apps := read_json()
	a := 0
	// 打印解析后的数据
	for _, app := range apps {
		if a > 10 {
			break
		}
		a++
		fmt.Printf("AppID: %d, Name: %s\n", app.Appid, app.Name)
	}

	// write_json(apps)
}

func request_json() []App {
	resp, err := http.Get("https://api.steampowered.com/ISteamApps/GetAppList/v2/")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil
	}
	apps := response.Applist.Apps
	return apps
}

func write_json(apps []App) {

	// 创建文件（会覆盖原文件）
	file, err := os.Create(file_name)
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
}

func read_json() []App {
	file, err := os.Open(file_name)
	if err != nil {
		fmt.Println("打开文件失败:", err)
		return nil
	}
	defer file.Close()

	// 解码 JSON
	var apps []App
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&apps)
	if err != nil {
		fmt.Println("解析 JSON 失败:", err)
		return nil
	}
	return apps
}
