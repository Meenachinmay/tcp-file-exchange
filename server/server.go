package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
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
var wg sync.WaitGroup

func StartServer() error {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("error creating listener: %w", err)
	}
	defer listener.Close()

	fmt.Printf("Listening on port %s\n", port)

	// worker pool start
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
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

		// can write a wait group here

		// Submit the job to the job queue
		JobQueue <- Job{conn: conn}
	}
	wg.Wait()
	close(JobQueue)

	// here we can wait for all tasks to be completed before finishing the method.
	return nil
}

func worker(id int) {
	defer wg.Done()
	for job := range JobQueue {
		log.Printf("Worker %d received job %s\n", id, job.conn.RemoteAddr())
		handleConnection(job.conn)
	}
}
