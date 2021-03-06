package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	connect("0.0.0.0", 8123)

}

func connect(host string, port int) {
	ip := &net.TCPAddr{IP: net.ParseIP(host)}
	portNum := &net.TCPAddr{Port: port}
	conn, err := net.DialTCP("tcp", ip, portNum)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	go receiveRemoteServerMsgs(conn)
	go localClientInput(conn)
	for {
	}
}

// Reads data from local client Stdin and sends across conn
func localClientInput(conn net.Conn) {
	fmt.Print("> ")
	scanner := bufio.NewScanner(os.Stdin)
	writer := bufio.NewWriter(conn)
	for scanner.Scan() {
		fmt.Print("> ")
		_, err := writer.WriteString(scanner.Text() + "\n")
		if err != nil {
			log.Println(err)
		}
		writer.Flush()
		if err != nil {
			log.Println(err)
			conn.Close()
		}
	}
}

// listens continuously for messages from server
func receiveRemoteServerMsgs(conn net.Conn) {
	for {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Print("\r")
			fmt.Print(scanner.Text() + "\n")
			fmt.Print("> ")
		}
		if scanner.Err() != nil {
			log.Printf("error: %s\n", scanner.Err())
			conn.Close()
			os.Exit(1)
		}
	}
}
