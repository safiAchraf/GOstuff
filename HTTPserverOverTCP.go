package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Server is listening on port 4221")

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Connection accepted")

	reader := bufio.NewReader(conn)

	headers, err := parseHeaders(reader)
	if err != nil {
		fmt.Println("Error parsing headers:", err.Error())
		return
	}

	fmt.Println("Method:", headers["Method"])
	fmt.Println("Path:", headers["Path"])
	fmt.Println("Version:", headers["Version"])

	userAgent := headers["User-Agent"]
	if userAgent != "" {
		fmt.Println("User-Agent:", userAgent)
	} else {
		fmt.Println("User-Agent header not found")
	}

	contentLength, _ := strconv.Atoi(headers["Content-Length"])
	if contentLength > 0 {
		body, err := readBody(reader, contentLength)
		if err != nil {
			fmt.Println("Error reading body:", err.Error())
			return
		}
		fmt.Println("Body:", string(body))
	}

	conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Length: 11\r\n\r\nHello World"))
}

func parseHeaders(reader *bufio.Reader) (map[string]string, error) {
	headers := make(map[string]string)
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	requestParts := strings.Fields(requestLine)
	if len(requestParts) < 3 {
		return nil, fmt.Errorf("invalid request line")
	}

	headers["Method"] = requestParts[0]
	headers["Path"] = requestParts[1]
	headers["Version"] = requestParts[2]

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		if line == "\r\n" {
			break
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			headers[key] = value
		}
	}
	return headers, nil
}

func readBody(reader *bufio.Reader, contentLength int) ([]byte, error) {
	body := make([]byte, contentLength)
	_, err := io.ReadFull(reader, body)
	if err != nil {
		return nil, err
	}
	return body, nil
}


func sendFile(conn net.Conn, filePath string) {

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err.Error())
		sendResponce(conn, "Not Found", 404)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err.Error())
		sendResponce(conn, "Internal Server Error", 500)
		return
	}
	fileSize := fileInfo.Size()

	headers := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Length: %d\r\n\r\n", fileSize)
	conn.Write([]byte(headers))

	_, err = io.Copy(conn, file)
	if err != nil {
		fmt.Println("Error sending file:", err.Error())
		return
	}
}

func sendResponce (conn net.Conn, responce string , statusCode int) {
	conn.Write([]byte("HTTP/1.1 " + strconv.Itoa(statusCode) + " " + responce + "\r\nContent-Length: 0\r\n\r\n"))
}