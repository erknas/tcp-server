package main

import (
	"fmt"
	"log"
)

func main() {
	srv := NewServer(":3000")

	go func() {
		for msg := range srv.msgch {
			fmt.Printf("Received from: %s\nMessage: %s\n", msg.from, msg.payload)
		}
	}()

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
