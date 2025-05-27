$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -o steam-list main.go
