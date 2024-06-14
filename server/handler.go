package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"youtube/tce-server/shared"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("client connected", conn.RemoteAddr())

	//
	dirPath := filepath.Join("received_data")
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating directory %s: %v", err)
		return
	}

	// read metadata
	fileName, err := readFileName(conn) // get file name
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading file name: %v", err)
		return
	}

	fileSize, err := readFileSize(conn) // get file size
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading file size: %v", err)
		return
	}

	filePath := filepath.Join(dirPath, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create file: %s\n", filePath)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(conn)
	buffer := make([]byte, shared.BufferSize)
	var receivedBytes int64

	for {
		n, err := reader.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Fprintf(os.Stderr, "Error reading data from connection: %s\n", err)
			}
			break
		}

		writtenBytes, writeErr := file.Write(buffer[:n])
		if writeErr != nil {
			fmt.Fprintf(os.Stderr, "Error writing to file: %s\n", writeErr)
			return
		}
		receivedBytes += int64(writtenBytes)
		printProgress("written to the file: ", receivedBytes, fileSize)
	}

	fmt.Println("\n File received successfully and saved to", filePath)
}

func printProgress(prefix string, current, total int64) {
	percentage := float64(current) / float64(total) * 100
	fmt.Printf("\r%s: %.2f%%\n", prefix, percentage)
}

func readFileSize(conn net.Conn) (int64, error) {
	sizeBuffer := make([]byte, 8)
	_, err := io.ReadFull(conn, sizeBuffer)
	if err != nil {
		return 0, err
	}
	fileSize := int64(binary.LittleEndian.Uint64(sizeBuffer))
	return fileSize, nil
}

func readFileName(conn net.Conn) (string, error) {
	return readString(conn)
}

func readString(conn net.Conn) (string, error) {
	lenBuffer := make([]byte, 8)
	_, err := io.ReadFull(conn, lenBuffer)
	if err != nil {
		return "", err
	}
	strLen := binary.LittleEndian.Uint64(lenBuffer)
	strBuffer := make([]byte, strLen)
	_, err = io.ReadFull(conn, strBuffer)
	if err != nil {
		return "", err
	}
	return string(strBuffer), nil
}
