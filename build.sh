GOOS=linux GOARCH=amd64 go build -o build/hbparser-linux-amd64
GOOS=linux GOARCH=arm go build -o build/hbparser-linux-arm
GOOS=darwin GOARCH=amd64 go build -o build/hbparser-macos-amd64
go build -o build/hbparser-macos-arm
GOOS=windows GOARCH=amd64 go build -o build/hbparser-windows-amd64
sudo chmod build/hbparser* 755
