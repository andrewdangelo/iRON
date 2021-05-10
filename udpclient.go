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
)

var clear map[string]func() //create a map for storing clear funcs

var IP_LENGTH = 14

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
	c, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		fmt.Println(err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println(err)
	}
	encryptedMessage := gcm.Seal(nonce, nonce, text, nil)
	return encryptedMessage
}

func decrypt(ciphertext []byte, key []byte) []byte {
	c, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		fmt.Println(err)
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		fmt.Println(err)
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	return plaintext
}

func encryptionTest() {
	text := []byte("This message is encrypted!!! hiiii!!! Yayyyy!!")
	key := []byte("passphrasewhichneedstobe32bytes!") //needs to be 32 bytes
	fmt.Println("Before Encryption:")
	fmt.Println(string(text))
	fmt.Println(text)
	encryptedMessage := encrypt(text, key)
	fmt.Println("After Encryption:")
	fmt.Println(string(encryptedMessage))
	fmt.Println(encryptedMessage)
	fmt.Println("After decryption:")
	decryptedMessage := decrypt(encryptedMessage, key)
	fmt.Println(string(decryptedMessage))
	fmt.Println(decryptedMessage)
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

func encryptPacket(packet []byte, key []byte) []byte {
	numberOfHops := packet[0]
	payload := packet[1:len(packet)]
	encryptPayload := encrypt(payload, key)
	newPacket := append([]byte{numberOfHops}, encryptPayload...)
	return newPacket
}

func decryptPacket(packet []byte, key []byte) []byte {
	fmt.Println("decrypting packet")
	numberOfHops := packet[0]
	payload := packet[1:len(packet)]
	decryptedPayload := decrypt(payload, key)
	newPacket := append([]byte{numberOfHops}, decryptedPayload...)
	return newPacket
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
		encryptionTest()
	case 4:
		ddos()
	case 5:
		sendThrough3Servers()
	case 6:
		sendthroughmany()
	case 7:
		sendThrough3AnonServers()
	default:
		fmt.Println("Invalid option please choose again")
	}
	fmt.Println("-----------------------")
}

func sendThrough3Servers() {
	packet := buildPacketFromString("hello world!")
	packetwithdestination := addDestination(packet, "127.0.0.1:1234")
	packet2withdestination := addDestination(packetwithdestination, "127.0.0.1:1235")
	packet3withdestination := addDestination(packet2withdestination, "127.0.0.1:1236")
	sendPacket("127.0.0.1:1234", packet3withdestination)
}

func sendThrough3AnonServers() {
	key := []byte("passphrasewhichneedstobe32bytes!") //needs to be 32 bytes
	packet := buildPacketFromString("hello world!")
	packetwithdestination := addDestination(packet, "127.0.0.1:1234")
	packetwithdestinationencrypted := encryptPacket(packetwithdestination, key)

	packet2withdestination := addDestination(packetwithdestinationencrypted, "127.0.0.1:1235")
	packet2withdestinationencrypted := encryptPacket(packet2withdestination, key)

	packet3withdestination := addDestination(packet2withdestinationencrypted, "127.0.0.1:1236")
	packet3ithdestinationencrypted := encryptPacket(packet3withdestination, key)

	sendPacket("127.0.0.1:1234", packet3ithdestinationencrypted)
}

func sendthroughmany() {

	for i := 0; i < 100; i++ {
		message := generatePacketPayload(1000)
		packet := buildPacketFromBytes(message)
		packetwithdestination := addDestination(packet, "127.0.0.1:1234")
		packet2withdestination := addDestination(packetwithdestination, "127.0.0.1:1235")
		packet3withdestination := addDestination(packet2withdestination, "127.0.0.1:1236")
		getDestinations(packet3withdestination)
		sendPacket("127.0.0.1:1234", packet3withdestination)
	}
}

func sendPacket(ip string, packet []byte) {
	conn, err := net.Dial("udp", ip)
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
	conn.Write(packet)
	conn.Close()
}

func buildPacketFromString(payload string) []byte {
	//convert the message in bytes
	byteMessage := []byte(payload)
	//sets the first byte to 0 which means it has no destination after being sent once
	packet := append([]byte{0}, byteMessage...)
	return packet
}

func buildPacketFromBytes(payload []byte) []byte {
	//sets the first byte to 0 which means it has no destination after being sent once
	packet := append([]byte{0}, payload...)
	return packet
}

//the ip is 4 bytes
func addDestination(packet []byte, ip string) []byte {
	numberOfDestinations := packet[0]
	numberOfDestinations++ //update the value
	ipBytes := []byte(ip)

	packetpayload := packet[1:len(packet)]

	fmt.Println("adding destination:")
	fmt.Println(ipBytes)
	packet[0] = numberOfDestinations
	packetwithdestination := append(ipBytes, packetpayload...)
	packetwithdestinationandcount := append([]byte{numberOfDestinations}, packetwithdestination...)
	return packetwithdestinationandcount
}

func getDestinations(packet []byte) {
	numberOfDestinations := packet[0]
	destinations := packet[1 : numberOfDestinations*byte(IP_LENGTH)]
	totalDestinations := (len(destinations) + 1) / IP_LENGTH

	fmt.Println("Total destinations", totalDestinations)

	for i := 0; i < totalDestinations; i++ {

		lowerBound := i * IP_LENGTH
		upperBound := i*IP_LENGTH + IP_LENGTH

		bytes := destinations[lowerBound:upperBound]
		str := string(bytes[:])
		fmt.Print(str)

		if i != totalDestinations-1 {
			fmt.Print(" ---> ")
		}
	}
	fmt.Print("\n")
}

func main() {

	CallClear()

	for {
		var option int = 0
		fmt.Println("-----------------------")
		fmt.Println("|      iRON Menu      |")
		fmt.Println("-----------------------")
		fmt.Println("0) To quit the program")
		fmt.Println("1) a message to the server")
		fmt.Println("2) Send in ordered messages to the server")
		fmt.Println("3) Test AES encryption")
		fmt.Println("4) DDos a server at a specific port")
		fmt.Println("5) Send Message through 3 servers")
		fmt.Println("6) Send many messages through 3 servers")
		fmt.Println("7) Send message through 3 anonymous servers")
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
