package main

import "net"

// структура клиента
type Client struct {
	name string
	id   int
}

// структура сервера
type Server struct {
	Clients     map[net.Conn]Client
	ClientCount int
}

// структура сообщений
type Message struct {
	Text       string
	SenderName string
	Sender     net.Conn
}

// пути и неизменяемые сообщения при событиях
const (
	enteringName = "[ENTER YOUR NAME]: "
	timeFormat   = "2006-01-02 15:04:05"
	greetingPath = "./files/greeting.txt"
	archivePath  = "archive.txt"
	networkType  = "tcp"
	errArgs      = "[USAGE]: ./TCPChat $port"
)
