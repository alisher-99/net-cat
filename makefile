pusk:
	go run .
newclient:
	nc localhost 8989
docker:
	docker run -d -p 8989:8989 net-cat