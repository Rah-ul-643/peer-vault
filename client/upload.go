package client

import (
	"fmt"
	"io"
	"net"
	"os"
	"peerVault/encryption"
	"strings"
)

func UploadFile(filePath string, addr string) {
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Failed to open file:", err)
		return
	}
	defer f.Close()
	
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("Failed to connect:", err)
		return
	}
	defer conn.Close()
	
	filenameSlices := strings.Split(filePath, ".")
	fileType := filenameSlices[len(filenameSlices)-1]

	conn.Write([]byte("UPLOAD\n"))
	conn.Write([]byte(fileType + "\n"))

	n, err := io.Copy(conn, f)
	if err != nil {
		fmt.Println("Failed to upload file:", err)
		return
	}

	filehash, err := encryption.HashFile(filePath)
	if err != nil {
		fmt.Println("Failed to hash file:", err)
		return
	}

	fmt.Printf("File uploaded, %d bytes sent. File hash: %s", n, filehash)
}
