package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"
)

// в этой функции идет регистрация нового пользователя
func registeringClient(conn net.Conn, server Server, deadConn chan net.Conn) {
	// прилетает сообщение приветствия
	conn.Write(greetingMessage())
	// выводит строку, чтобы ввести в поле имя клиента
	conn.Write([]byte(enteringName))
	buf := make([]byte, 1024)
	nbyte, err := conn.Read(buf)
	// подгружаются прошлые сообщения из файла archive.txt
	priorMsg := loadPriorMsg()
	conn.Write(priorMsg)

	if err != nil {
		deadConn <- conn
	} else {

		user := make([]byte, nbyte-1)
		copy(user, buf[:nbyte])

		server.Clients[conn] = Client{
			id:   1,
			name: string(user),
		}
		for conn := range server.Clients {
			text := fmt.Sprintf("\n%s has joined our chat...\n", user)
			conn.Write([]byte(text))
			time := getTime()
			text = fmt.Sprintf("[%s][%s]:", time, server.Clients[conn].name)
			conn.Write([]byte(text))
		}

	}
}

func acceptingMsg(conn net.Conn, messages chan Message, server Server, deadConn chan net.Conn) {
	client := server.Clients[conn]
	if client.id == 0 {
		registeringClient(conn, server, deadConn)
	}

	buf := make([]byte, 1024)
	for {
		nbyte, err := conn.Read(buf)
		if err != nil {
			deadConn <- conn
			fmt.Println("connection terminated")

			break
		} else {
			message := make([]byte, nbyte)
			copy(message, buf[:nbyte])
			if flag := checkEmptyMsg(message); flag {
				clientName := server.Clients[conn].name

				readyMessage := Message{
					Text:       string(message),
					SenderName: clientName,
					Sender:     conn,
				}
				messages <- readyMessage
			}
		}
	}
}

func checkEmptyMsg(b []byte) bool {
	flag := false
	for i := 0; i <= len(b)-2; i++ {
		if b[i] != ' ' {
			flag = true
		}
	}
	return flag
}

func getTime() string {
	currentTime := time.Now()
	time := currentTime.Format(timeFormat)
	return time
}

func greetingMessage() []byte {
	data, err := os.ReadFile(greetingPath)
	if err != nil {
		fmt.Println(err)
	}
	return data
}

// функция для создания и записи архива сообщений
func logConfig() {
	file, err := os.Create(archivePath)
	if err != nil {
		fmt.Println(err)
	}
	log.SetOutput(file)
	log.SetFlags(0)
}

func loadPriorMsg() []byte {
	file, err := os.Open(archivePath)
	if err != nil {
		fmt.Println(err)
	}
	text, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	return text
}
