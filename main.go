package main

import (
	"fmt"
	"net"
	"log"
	"os"
	"strings"
	"errors"
	"flag"
)

func main() {
	SERVER_HOST := flag.String("host", "localhost", "IP for the server.")
	SERVER_PORT := flag.String("port", "8080", "Port number for the server.")
	flag.Parse()

	fmt.Println("Running server..")

	server, err := net.Listen("tcp", *SERVER_HOST + ":" + *SERVER_PORT)
	checkErr(err)

	defer server.Close()

	fmt.Println("Listening on: " + *SERVER_HOST + ":" + *SERVER_PORT)
	
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
	checkErr(err)

	fileName := parseRequest(buffer)
	path, err := os.Getwd()
	checkErr(err)

	page, err := loadHTML(path + fileName)
	if err != nil {
		respond(404, connection, err.Error())
	} else {
		respond(200, connection, page)
	}

	connection.Close()
}

func loadHTML(path string) (string, error) {
	dat, err := os.ReadFile(path)
	if err != nil {
		return "", errors.New("File not found.")
	}
	html := string(dat[:])
	return html, nil
}

func parseRequest(request []byte) string {
	header := string(request[:])
	headers := strings.Split(header, "\n")
	fileName := strings.Split(headers[0], " ")[1]
	if (fileName == "/") {
		return "/index.html"
	} else {
		return fileName
	}
}

func respond(code int, connection net.Conn, message string) {
	switch code {
	case 200:
		_, err := connection.Write([]byte("HTTP/1.0 200 OK\n\n" + message))
		checkErr(err)
	case 404:
		_, err := connection.Write([]byte("HTTP/1.0 404 NOT FOUND\n\n" + message))
		checkErr(err)
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
