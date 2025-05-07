# go-steam-list

## 介绍

steam 平台有这样一个 api 可供获取所有 staem app 的 id name 对：`https://api.steampowered.com/ISteamApps/GetAppList/v2/`，但其有缺陷，每次获取时总会少将近 10%的应用信息，原因不明。本项目
用于定时访问接口，在多次访问的冗余中补全缺失的部分应用信息。

## 使用方式

使用 go 运行，并等待你想等待的时间（建议一直挂着），获取到的所有应用会存到 app.json 文件中

```go
go run .
```

## 许可证

本项目使用 MIT 许可证
