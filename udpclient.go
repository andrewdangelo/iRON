package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

var clear map[string]func() //create a map for storing clear funcs

func init() {
	fmt.Println("init")
	clear = make(map[string]func()) //Initialize it

	clear["darwin"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

}

func CallClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

func encrypt(text []byte, key []byte) []byte {

	// generate a new aes cipher using our 32 byte long key
	c, err := aes.NewCipher(key)
	// if there are any errors, handle them
	if err != nil {
		fmt.Println(err)
	}

	// gcm or Galois/Counter Mode, is a mode of operation
	// for symmetric key cryptographic block ciphers
	gcm, err := cipher.NewGCM(c)
	// if any error generating new GCM
	// handle them
	if err != nil {
		fmt.Println(err)
	}

	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, gcm.NonceSize())
	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println(err)
	}

	// here we encrypt our text using the Seal function
	// Seal encrypts and authenticates plaintext, authenticates the
	// additional data and appends the result to dst, returning the updated
	// slice. The nonce must be NonceSize() bytes long and unique for all
	// time, for a given key.

	encryptedMessage := gcm.Seal(nonce, nonce, text, nil)

	return encryptedMessage

}

func sendEncryptedMessage() {
	text := []byte("This message is encrypted")
	key := []byte("passphrasewhichneedstobe32bytes!") //needs to be 32 bytes
	fmt.Println(text)
	encryptedMessage := encrypt(text, key)
	fmt.Println(encryptedMessage)

	//send the encrypted message
	p := make([]byte, 2048)

	conn, err := net.Dial("udp", "127.0.0.1:1234")

	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}

	conn.Write(encryptedMessage)

	_, err = bufio.NewReader(conn).Read(p)
	if err == nil {
		fmt.Printf("%s\n", p)
	} else {
		fmt.Printf("Some error %v\n", err)
	}
	conn.Close()

}
func sendMessage() {
	p := make([]byte, 2048)
	conn, err := net.Dial("udp", "127.0.0.1:1234")
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}

	fmt.Fprintf(conn, "Hi UDP Server, How are you doing?")
	_, err = bufio.NewReader(conn).Read(p)
	if err == nil {
		fmt.Printf("%s\n", p)
	} else {
		fmt.Printf("Some error %v\n", err)
	}
	conn.Close()
}

func sendInOrderMessages() {
	p := make([]byte, 2048)
	conn, err := net.Dial("udp", "127.0.0.1:1234")
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}

	for i := 0; i < 100; i++ {
		var msg string = strconv.Itoa(i)
		fmt.Fprintf(conn, msg)
	}

	_, err = bufio.NewReader(conn).Read(p)
	if err == nil {
		fmt.Printf("%s\n", p)
	} else {
		fmt.Printf("Some error %v\n", err)
	}
	conn.Close()
}

func sendThroughTwoServers() {
	//the message to send
	messageToSend := "127.0.0.1:1235~Hello this message is going through 2 servers"

	//convert the message in bytes
	byteMessage := []byte(messageToSend)

	fmt.Println("Btye Message: ", byteMessage)

	//put the message in a packet
	packet := gopacket.NewPacket(byteMessage, layers.LayerTypeUDP, gopacket.Default)

	fmt.Println(packet.String())

	//call up the udp server
	conn, err := net.Dial("udp", "127.0.0.1:1234")

	conn.Write(packet.Data())
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}

	conn.Close()
}

func ddos() {

	port := 1234
	packetsToSend := 1000
	millisecondDelay := 0

	fmt.Printf("Port to ddos: ")
	fmt.Scanln(&port)
	fmt.Printf("How many packets?: ")
	fmt.Scanln(&packetsToSend)
	fmt.Printf("How much delay(ms) between packets?: ")
	fmt.Scanln(&millisecondDelay)

	server := "127.0.0.1:" + strconv.Itoa(port)

	conn, err := net.Dial("udp", server)
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}

	//maximum size of a UDP packet is 65507 Bytes so we want to generate random messages that are this size to make
	//the network suffer
	for i := 0; i < packetsToSend; i++ {
		byteMessage := generatePacketPayload(1000)
		conn.Write(byteMessage)
		println("Packets sent: ", i+1)
		time.Sleep(time.Duration(millisecondDelay) * time.Millisecond)
	}

	conn.Close()
}

func generatePacketPayload(byteSize int) []byte {
	token := make([]byte, byteSize)
	if _, err := rand.Read(token); err != nil {
		// Handle err
	}
	return token
}

func menu(option int) {
	fmt.Println("-----------------------")
	fmt.Println("|        Output        |")
	fmt.Println()
	switch option {
	case 1:
		sendMessage()
	case 2:
		sendInOrderMessages()
	case 3:
		sendEncryptedMessage()
	case 4:
		sendThroughTwoServers()
	case 5:
		ddos()
	default:
		fmt.Println("Invalid option please choose again")
	}
	fmt.Println("-----------------------")
}

func main() {

	CallClear()

	for {
		var option int = 0
		fmt.Println("-----------------------")
		fmt.Println("|      iRON Menu      |")
		fmt.Println("-----------------------")
		fmt.Println("0) To quit the program")
		fmt.Println("1) Send a message to the server")
		fmt.Println("2) Send in ordered messages to the server")
		fmt.Println("3) Send an encrypted message to the server")
		fmt.Println("4) Send a message through 2 servers")
		fmt.Println("5) DDos a server at a specific port")
		fmt.Println("")
		fmt.Println("Type the number of the option you would like to select")
		fmt.Scanln(&option)

		CallClear()

		if option == 0 {
			break
		} else {
			menu(option)
		}

	}

	fmt.Println("Closing the iRON client...")

}
