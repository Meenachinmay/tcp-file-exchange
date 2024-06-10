package main

import "log"

func main() {
	err := StartServer()

	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
