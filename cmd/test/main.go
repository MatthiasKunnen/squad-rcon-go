package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gorcon/rcon"
)

func main() {
	conn, err := rcon.Dial("", "")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	response, err := conn.Execute("ListCommands 0")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response)

	for {
		fmt.Println(time.Now().Format("2006-01-02-15:04:05"))
		response, err = conn.Execute("ListPlayers")
		if err != nil {
			log.Println(err)
		} else {
			fmt.Println(response)
		}

		time.Sleep(time.Second)
	}
}
