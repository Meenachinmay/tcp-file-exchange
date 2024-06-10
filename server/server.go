package main

import (
	"fmt"
	"net"
	"os"
)

const (
	port = ":8080"
)

func StartServer() error {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("error creating listener: %w", err)
	}
	defer listener.Close()

	fmt.Printf("Listening on port %s\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error accepting connection: %v\n", err)
			continue
		}

		go handleConnection(conn)
	}

}
