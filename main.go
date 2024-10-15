package main

import "log"

func main() {
	srv := NewServer(":3000")

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
