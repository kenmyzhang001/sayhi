env GOOS=linux GOARCH=amd64 CGO_ENABLED=0  go  build -o sayhi
sudo docker restart sayhi
