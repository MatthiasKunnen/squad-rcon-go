package main

import (
	"fmt"
	"log"
	"squad-rcon-go/pkg/squadrcon"
	"time"
)

func main() {
	conn, err := squadrcon.Connect("", "", squadrcon.RconSettings{
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

	fmt.Println(response)
	time.Sleep(5 * time.Second)
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
