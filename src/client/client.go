package client
import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func Run(ip string, port string ) {
	// create connection with server
	conn, err := net.Dial("tcp", ip + ":" + port)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Printf("Gochat client connected")
	// go routine to listen incomming messages from serve
	go ListenIncommingMessage(conn)
	for true {
		//User's message from cmd to send to server
		ConsoleInput(conn)
	}
}
// listen for any meesage that come from the chat server
func ListenIncommingMessage(conn net.Conn) {
	for true {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Println(err)
			break
		}
		fmt.Print("Message from server: " + message + "\n")
		fmt.Print("Text to send: \n")
	}
}

// keep watching for console input
func ConsoleInput(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	// read in input from stdistdin
	fmt.Print("Text to send: \n")
	text, err := reader.ReadString('\n')
	if err != nil {
		log.Println(err)
		return
	}
	// send to socket
	conn.Write([]byte(text))
}
