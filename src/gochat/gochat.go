package main
import (
	"fmt"
	"os"
	"server" // gochat server package
	"client" // gochat client package
	
)
func main() {
	args := os.Args[1:] // get cmd args, agrs[0]=path
	fmt.Println("Input args", args)
	// when no args provided then start server
	if (len(args)<=0){
		fmt.Println("Starting server ....")
	   server.Run()
   // when 2  args Ip and port provided then start client
   }else if (len(args) == 2){
		fmt.Printf("Connecting to server on %s:%s ....", args[0], args[1])
		client.Run(args[0], args[1])
	// Print usage in case of invalid a
	}else {
			fmt.Println("\n\n ************ \n Invalid Arags \n ************ \n Usage: \n gochat or \n gochat <IP> <Port>")
			fmt.Println(" example: gochat 127.0.0.1 3200")
		}
}