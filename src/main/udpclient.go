package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
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

func menu(option int) {
	fmt.Println("-----------------------")
	fmt.Println("|        Output        |")
	fmt.Println()
	switch option {
	case 1:
		sendMessage()
	case 2:
		sendInOrderMessages()
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
