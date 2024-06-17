package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"youtube/tce-server/shared"
)

func StartClient(serverAddress, filePath string) error {
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return fmt.Errorf("error connecting to server: %v", err)
	}
	defer conn.Close()

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("error getting file info: %v", err)
	}

	fileSize := fileInfo.Size()
	fileName := fileInfo.Name()

	err = sendMetaData(conn, fileName, fileSize)
	if err != nil {
		return fmt.Errorf("error sending metadata: %v", err)
	}
	var sentBytes int64

	buffer := make([]byte, shared.BufferSize)
	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return fmt.Errorf("error reading file: %v", err)
		}
		if n == 0 {
			break
		}

		_, err = conn.Write(buffer[:n])
		if err != nil {
			return fmt.Errorf("error writing to server: %v", err)
		}
		sentBytes += int64(n)
		printProgress("Sent", sentBytes, fileSize)
	}

	// Wait for server acknowledgment before closing
	ackBuffer := make([]byte, 256)
	_, err = conn.Read(ackBuffer)
	if err != nil && err != io.EOF {
		fmt.Println("Error reading acknowledgment from server:", err)
	} else {
		fmt.Println("acknowledgment received from server:\n")
	}

	fmt.Println("File sent successfully")
	return nil
}

func printProgress(prefix string, current, total int64) {
	percentage := float64(current) / float64(total) * 100
	fmt.Printf("\r%s: %.2f%%\n", prefix, percentage)
}

func sendFileSize(conn net.Conn, fileSize int64) error {
	sizeBuffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(sizeBuffer, uint64(fileSize))
	_, err := conn.Write(sizeBuffer)
	return err
}

func sendMetaData(conn net.Conn, fileName string, fileSize int64) error {
	// send file name
	err := sendString(conn, fileName)
	if err != nil {
		return err
	}

	// send file size
	err = sendFileSize(conn, fileSize)
	if err != nil {
		return err
	}
	return nil
}

func sendString(conn net.Conn, str string) error {
	strLen := uint64(len(str))
	lenBuffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(lenBuffer, strLen)
	_, err := conn.Write(lenBuffer)
	if err != nil {
		return err
	}
	_, err = conn.Write([]byte(str))
	return err
}
