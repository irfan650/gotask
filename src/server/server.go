
package server
import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
)

func Run() {
	clientCount := 0
	Clients := make(map[net.Conn]int) //Store number of clients Map<Conn Object, Clinet ID>
	// Channel into which the TCP server will push new connections.
	newConnections := make(chan net.Conn)  // channel to push all new connections
	deadConnections := make(chan net.Conn) // Channel to rmove clients connection when disconnected
	messages := make(chan string)          // channel to push messages to all connected clients like a broadcast
	// Start the TCP server port
	port := 3200
	IP:= "127.0.0.1"
	server, err := net.Listen("tcp", ":"+ strconv.Itoa(port))
	defer server.Close()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Printf("Gochat Server Started: IP %s, Port %d", IP, port)
	// go routine to accept new connections
	// It is lisning for new clients conections forever
	go func() {
		for {
			conn, err := server.Accept()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			newConnections <- conn
		}
	}()

	// Handle in case of new connections
	// Handle in case of dead connections
	// Handle to send/recieve messeages between connections
	for {
		select {
		// Accept new clients
		case conn := <-newConnections:
			clientCount++
			log.Printf("Accepted new client, #%d", clientCount)
			Clients[conn] = clientCount
			//go routine to recive incomming messages from connection
			go func(conn net.Conn, clientId int) {
				reader := bufio.NewReader(conn)
				for {
					message_rec, err := reader.ReadString('\n')
					if err != nil {
						log.Printf("Cannot read message at this time")
						break
					}
					messages <- fmt.Sprintf("Client %d > %s", clientId, message_rec)
				}
				deadConnections <- conn
			}(conn, Clients[conn])
			// send message to all connections in case of any nw message
		case message := <-messages:
			// Loop over all connected clients
			log.Printf("<< Recived message: <%s>", message)
			for conn, _ := range Clients {
				// go-routine to bradcast message so the network operation cannot block
				go func(conn net.Conn, message string) {
					_, err := conn.Write([]byte(message))
					// in case of error then the connection is dead.
					if err != nil {
						deadConnections <- conn
					} else {
						log.Printf(">> send message <%s> to: <Client %d>", message, Clients[conn])
					}
				}(conn, message)
			}
		// Remove dead clients
		case conn := <-deadConnections:
			log.Printf("Client %d disconnected", Clients[conn])
			delete(Clients, conn)
			clientCount--
		}
	}
}
