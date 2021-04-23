package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

const ENABLE_LOGGING_SETTING = true

func sendResponse(conn *net.UDPConn, addr *net.UDPAddr) {
	_, err := conn.WriteToUDP([]byte("From server: Hello I got your message "), addr)
	if err != nil {
		fmt.Printf("Couldn't send response %v", err)
	}
}

func forwardMessage(destination string, message string) {

	println("----------")
	println("forwaring to:" + destination + "333")
	println("----------")

	conn, err := net.Dial("udp", destination)
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}

	fmt.Println("attempting to forward a message to " + destination)

	//convert the message in bytes
	byteMessage := []byte(message)
	conn.Write(byteMessage)

}

func main() {
	p := make([]byte, 2048) //creates a "slice"

	var PORT int = 1234
	enable_logging := false

	fmt.Printf("Port Number: ")
	fmt.Scanln(&PORT)

	if ENABLE_LOGGING_SETTING {
		fmt.Printf("Enable Logging (yes or no): ")
		option := ""
		fmt.Scanln(&option)
		if option == "yes" {
			enable_logging = true
		}
	}

	addr := net.UDPAddr{
		Port: PORT,
		IP:   net.ParseIP("127.0.0.1"),
	}

	fmt.Println("==============================================")
	fmt.Println("")
	fmt.Println(" Started udpserver at ", addr.IP, ":", addr.Port)
	fmt.Println(" Logging enabled:", enable_logging)
	fmt.Println("")
	fmt.Println("==============================================")

	server, err := net.ListenUDP("udp", &addr)

	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}

	//initialize logging
	f, err := os.Create("logs.txt")

	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}

	for {
		_, remoteaddr, err := server.ReadFromUDP(p)
		now := time.Now()

		message := now.String() + " | Read a message from | " + remoteaddr.String() + " | " + string(p)

		fmt.Println(message)

		if enable_logging {
			fmt.Fprintln(f, message)
			if err != nil {
				fmt.Println(err)

			}
		}
		if err != nil {
			fmt.Printf("Some error  %v", err)
			continue
		}

		/*
			stringPayload := string(p)

			split := s.Split(stringPayload, "~")


				if len(split) >= 2 {
					go forwardMessage(split[0], split[1])
				}
		*/

		go sendResponse(server, remoteaddr)
	}

}
