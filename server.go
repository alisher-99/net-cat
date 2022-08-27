package main

import (
	"fmt"
	"log"
	"net"
)

func serving(server *Server, port string) {
	newConn := make(chan net.Conn)
	msgs := make(chan Message)
	endConn := make(chan net.Conn)
	logConfig()

	listener, err := net.Listen(networkType, port)
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				endConn <- conn
			}
			newConn <- conn

		}
	}()

	for {
		var logger string
		select {
		case conn := <-newConn:
			fmt.Println("new connection")
			server.ClientCount++
			if server.ClientCount > 10 {
				fmt.Println("Client max count limit exceed")
			} else {
				go acceptingMsg(conn, msgs, *server, endConn)
			}
		case message := <-msgs:
			for conn := range server.Clients {
				time := getTime()
				logger = fmt.Sprintf("[%s][%s]:%s", time, message.SenderName, message.Text)

				if conn == message.Sender {
					text := fmt.Sprintf("[%s][%s]:", time, message.SenderName)
					message.Sender.Write([]byte(text))
				} else {
					if len(message.Text) == 1 && message.Text[0] == 10 {
						continue
					}
					text := fmt.Sprintf("\n[%s][%s]:%s", time, message.SenderName, message.Text)
					conn.Write([]byte(text))

					text = fmt.Sprintf("[%s][%s]:", time, server.Clients[conn].name)
					conn.Write([]byte(text))
				}

			}
			log.Print(logger)
		case deadConnection := <-endConn:
			clientName := server.Clients[deadConnection].name
			delete(server.Clients, deadConnection)
			for conn, client := range server.Clients {
				server.ClientCount--

				text := fmt.Sprintf("\n%s has left our chat...\n", clientName)
				conn.Write([]byte(text))

				time := getTime()
				text = fmt.Sprintf("[%s][%s]:", time, client.name)
				conn.Write([]byte(text))

			}
		}
	}
}
