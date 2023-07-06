deploy:
	sudo systemctl stop words
	# git pull
	go build
	sudo cp ./words.service /etc/systemd/system/words.service
	systemctl enable words
	systemctl start words
	systemctl -l status words