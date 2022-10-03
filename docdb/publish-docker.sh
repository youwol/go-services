#bash
env GOOS=linux GARCH=amd64 CGO_ENABLED=0 go build ./cmd/docdb-server
sudo docker build -t registry.gitlab.com/youwol/platform/docdb:0.3.42 .
sudo docker push registry.gitlab.com/youwol/platform/docdb:0.3.42