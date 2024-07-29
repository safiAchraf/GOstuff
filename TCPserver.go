package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	// Start the TCP server
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		log.Fatal("Failed to bind to port 4221")
	}
	defer l.Close()
	fmt.Println("Server is listening on port 4221")

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Connection accepted")

	// Send initial message to client
	conn.Write([]byte("Hello, client! \n"))

	// Buffer to store incoming data
	buf := make([]byte, 1024)

	// Read the incoming message
	n, err := conn.Read(buf)
	if err != nil {
		log.Println("Error reading:", err.Error())
		return
	}

	// Print the received message
	fmt.Println("Received message:", string(buf[:n]))

	// Echo the message back to the client
	conn.Write(buf[:n])
}