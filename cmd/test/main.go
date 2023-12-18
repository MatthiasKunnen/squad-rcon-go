package main

import (
	"fmt"
	"log"
	"squad-rcon-go/pkg/squadrcon"
	"strings"
)

func main() {
	conn, err := squadrcon.Connect("", "", squadrcon.Settings{
		DialTimeout:   0,
		PacketIdStart: 10000,
		WriteTimeout:  0,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	response, err := conn.Execute("ListCommands 1")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(
		"%[1]s Execute response %[1]s\n%s\n%[1]s END OF RESPONSE %[1]s\n",
		strings.Repeat("=", 8),
		response,
	)
	//
	//for {
	//	fmt.Println(time.Now().Format("2006-01-02-15:04:05"))
	//	response, err = conn.Execute("ListPlayers")
	//	if err != nil {
	//		log.Println(err)
	//	} else {
	//		fmt.Println(response)
	//	}
	//
	//	time.Sleep(time.Second)
	//}
}
