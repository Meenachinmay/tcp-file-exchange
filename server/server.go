package main

import (
	"fmt"
	"net"
	"os"
)

const (
	port               = ":8080"
	numWorkers         = 10 // number of workers in the pool (no of goroutines)
	maxNoOfConnections = 1000
)

type Job struct {
	conn net.Conn
}

var JobQueue = make(chan Job, maxNoOfConnections)

func StartServer() error {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("error creating listener: %w", err)
	}
	defer listener.Close()

	fmt.Printf("Listening on port %s\n", port)

	// worker pool start
	for i := 0; i < numWorkers; i++ {
		go worker(i)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error accepting connection: %v\n", err)
			continue
		}

		// Inform the client that it is in the queue
		_, err = conn.Write([]byte("Your file transfer is in the queue. Please wait...\n"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error sending queue message: %v", err)
			conn.Close()
			continue
		}
		// Submit the job to the job queue
		JobQueue <- Job{conn: conn}
	}

}

func worker(id int) {
	for job := range JobQueue {
		fmt.Printf("Worker %d received job %s\n", id, job.conn.RemoteAddr())
		handleConnection(job.conn)
	}
}
