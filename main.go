package main

import (
	"fmt"
	"net"
	"log"
	"os"
	"strings"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "8080"
	SERVER_TYPE = "tcp"
)

func main() {
	fmt.Println("Running server..")

	server, err := net.Listen(SERVER_TYPE, SERVER_HOST+":"+SERVER_PORT)
	checkErr(err)

	defer server.Close()

	fmt.Println("Listening on: " + SERVER_HOST + ":" + SERVER_PORT)
	
	for {
		connection, err := server.Accept()
		checkErr(err)
		fmt.Println("Client connected")
		go processClient(connection)
	}
}

func processClient(connection net.Conn) {
	buffer := make([]byte, 1024)
	_, err := connection.Read(buffer)

	fileName := requestParse(buffer)
	if (fileName == "/") {
		fileName = "/index.html"
	}

	page := loadHTML("htdocs" + fileName)
	_, err = connection.Write([]byte("HTTP/1.1 200 OK\n\n" + page))
	checkErr(err)

	connection.Close()
}

func loadHTML(path string) string {
	dat, err := os.ReadFile(path)
	html := string(dat[:])
	checkErr(err)
	return html
}

func requestParse(request []byte) string {
	header := string(request[:])
	headers := strings.Split(header, "\n")
	fileName := strings.Split(headers[0], " ")[1]
	return fileName
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
