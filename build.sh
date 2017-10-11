GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -x -o proxyamd64 main.go
GOOS=linux GOARCH=386 go build -ldflags "-s -w" -x -o proxy386 main.go
GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -x -o proxywindows_amd64.exe main.go
GOOS=windows GOARCH=386 go build -ldflags "-s -w" -x -o proxywindows_386.exe main.go
