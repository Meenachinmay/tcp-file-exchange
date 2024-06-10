package main

import "log"

func main() {
	err := StartClient("localhost:8080", "./testFile.mp4")
	if err != nil {
		log.Fatal("Error starting client:", err)
	}
}
