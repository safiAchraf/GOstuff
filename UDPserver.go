package main

import (
	"fmt"
	"net"
)


func main() {
	addr, err := net.ResolveUDPAddr("udp", ":8080")
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer conn.Close()

	fmt.Println("UDP server listening on port 8080")

	// Buffer to hold incoming data
	buffer := make([]byte, 1024)
	newaddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:3434")

	for {
		// Read from the connection
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			continue
		}
		
		fmt.Printf("Received %s from %s\n", string(buffer[:n]), addr)
		_, err = conn.WriteToUDP([]byte("Hello from UDP server"), newaddr)

		if err != nil {
			fmt.Println("Error sending data:", err)
		}

		fmt.Println("Data sent to client")
		_, err = conn.WriteToUDP([]byte("Hello client , this is the UDP server"), addr)
	
	}
}