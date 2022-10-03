#bash
env GOOS=linux GARCH=amd64 CGO_ENABLED=0 go build ./cmd/storage-server
sudo docker build -t registry.gitlab.com/youwol/platform/storage:0.2.13 .
sudo docker push registry.gitlab.com/youwol/platform/storage:0.2.13