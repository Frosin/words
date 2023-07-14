deploy:
	sudo systemctl stop words
	git pull
	go test ./...
	go build -o words ./cmd/main.go
	sudo cp ./words.service /etc/systemd/system/words.service
	systemctl enable words
	systemctl start words
	systemctl -l status words