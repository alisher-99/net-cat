package main

import (
	"fmt"
	"net"
	"os"
)

var port = ":8989"

func main() {
	args := os.Args[1:]
	if len(args) == 1 {
		port = ":" + args[0]
	} else if len(args) > 1 {
		fmt.Println(errArgs)
		return
	}
	fmt.Println("Server is listening on port", port, "using tcp")

	server := new(Server)
	server.Clients = make(map[net.Conn]Client)

	serving(server, port)
}
