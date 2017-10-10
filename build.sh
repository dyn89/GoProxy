GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -x -o proxy_amd64 main.go
GOOS=linux GOARCH=386 go build -ldflags "-s -w" -x -o proxy_386 main.go
GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -x -o proxy_windows_amd64.exe main.go
GOOS=windows GOARCH=386 go build -ldflags "-s -w" -x -o proxy_windows_386.exe main.go
