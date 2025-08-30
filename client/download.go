package client

import (
	"fmt"
	"io"
	"net"
	"os"
)

func Download(hash, addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("Failed to connect:", err)
		return
	}
	defer conn.Close()

	payload := "DOWNLOAD\n" + hash + "\n"

	conn.Write([]byte(payload))

	data, err := io.ReadAll(conn)
	if err != nil {
		fmt.Println("Failed to receive file:", err)
		return
	}

	if err := os.MkdirAll("client_store", 0755); err != nil {
		fmt.Println("Failed to create directory:", err)	
		return
	}
	os.WriteFile("client_store/downloaded_"+hash+".txt", data, 0644)
	fmt.Println("File downloaded and saved.")
}
