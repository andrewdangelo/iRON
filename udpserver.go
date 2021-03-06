package main

import (
	"crypto/aes"
	"crypto/cipher"
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

func forwardMessage(destination string, message []byte) {

	println("-------------------")
	println("forwaring to:" + destination)
	println(message)
	println("-------------------")

	conn, err := net.Dial("udp", destination)
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}

	conn.Write(message)
}

func unfoldPacket(packet []byte) (int, string, []byte) {
	var IP_LENGTH = 14
	numberOfDestinations := packet[0]
	numberOfDestinations = numberOfDestinations - 1
	var destinations []byte
	var totalDestinations int = 0

	if numberOfDestinations != 0 {
		destinations = packet[1 : numberOfDestinations*byte(IP_LENGTH)]
		totalDestinations = (len(destinations) + 1) / IP_LENGTH
	}

	if totalDestinations != 0 {
		nextdestination := packet[1 : IP_LENGTH+1]
		newPacket := packet[IP_LENGTH+1 : len(packet)]
		packetwithcount := append([]byte{numberOfDestinations}, newPacket...)
		return totalDestinations, string(nextdestination), packetwithcount
	}
	return 0, "", packet
}

func decrypt(ciphertext []byte, key []byte) []byte {
	length := Find(ciphertext, 0)
	trimmed := ciphertext[0:length]

	c, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		fmt.Println(err)
	}
	nonceSize := gcm.NonceSize()
	if len(trimmed) < nonceSize {
		fmt.Println(err)
	}
	nonce, trimmed := ciphertext[:nonceSize], trimmed[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, trimmed, nil)

	return plaintext
}

func Find(arr []byte, key byte) int {
	for i, n := range arr {
		if key == n {
			return i
		}
	}
	return -1
}

func decryptPacket(packet []byte, key []byte) []byte {
	fmt.Println("decrypting packet")
	numberOfHops := packet[0]
	payload := packet[1:len(packet)]
	decryptedPayload := decrypt(payload, key)

	newPacket := append([]byte{numberOfHops}, decryptedPayload...)
	return newPacket
}

func processRequest(conn *net.UDPConn, packet []byte) {
	hopsleft, nextdestination, newpacket := unfoldPacket(packet)

	if hopsleft != 0 {
		fmt.Println("hops left: ", hopsleft)
		fmt.Println("destination: ", nextdestination)
		forwardMessage(nextdestination, newpacket)
	} else {
		fmt.Println("Hit destination!")
	}
}

//tor network like node
func processAnonymousRequest(conn *net.UDPConn, encryptedPacket []byte) {

	fmt.Println("processing anonymous request")
	key := []byte("passphrasewhichneedstobe32bytes!") //needs to be 32 bytes
	decryptedpacket := decryptPacket(encryptedPacket, key)

	fmt.Println("decrypted packet")
	fmt.Println(decryptedpacket)

	hopsleft, nextdestination, newpacket := unfoldPacket(decryptedpacket)

	fmt.Println("hops left: ", hopsleft)
	fmt.Println("Next destination: ", nextdestination)
	fmt.Println("New packet", newpacket)

	fmt.Println("Decrypted packet")
	fmt.Println(decryptedpacket)

	if hopsleft != 0 {
		fmt.Println("hops left: ", hopsleft)
		fmt.Println("destination: ", nextdestination)
		//forwardMessage(nextdestination, newpacket)
	} else {
		fmt.Println("Hit destination!")
	}
}

func main() {
	p := make([]byte, 2048) //creates a "slice"

	var PORT int = 1234
	enable_logging := false
	anonymous := false

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

	fmt.Printf("Enable Anon Node (yes or no): ")
	option := ""
	fmt.Scanln(&option)
	if option == "yes" {
		anonymous = true
	}

	addr := net.UDPAddr{
		Port: PORT,
		IP:   net.ParseIP("127.0.0.1"),
	}

	fmt.Println("==============================================")
	fmt.Println("")
	fmt.Println(" Started udpserver at ", addr.IP, ":", addr.Port)
	fmt.Println(" Logging enabled:", enable_logging)
	fmt.Println(" Anonymous Node: ", anonymous)
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

		if anonymous {
			go processAnonymousRequest(server, p)
		} else {
			go processRequest(server, p)
		}

		//every request gets a response back
		go sendResponse(server, remoteaddr)
	}

}
