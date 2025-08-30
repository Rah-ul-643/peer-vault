package main

import (
	"flag"
	"fmt"
	"os"
	
	"peerVault/client"
	"peerVault/p2p"
	"peerVault/storage"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Expected 'serve', 'upload', or 'get' subcommands.")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "serve":
		serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)
		port := serveCmd.Int("port", 8080, "Port to listen on")
		serveCmd.Parse(os.Args[2:])

		store := storage.NewStorage("peerVault_storage")
		s := server.NewServer(*port, "tcp", store)
		if err := s.ListenAndServe(); err != nil {
			fmt.Println("Server error:", err)
		}

	case "upload":
		uploadCmd := flag.NewFlagSet("upload", flag.ExitOnError)
		file := uploadCmd.String("file", "", "File to upload")
		to := uploadCmd.String("to", "", "Target peer address")
		uploadCmd.Parse(os.Args[2:])

		if *file == "" || *to == "" {
			fmt.Println("Usage: upload --file <file> --to <address>")
			os.Exit(1)
		}
		client.UploadFile(*file, *to)

	case "download":
		getCmd := flag.NewFlagSet("download", flag.ExitOnError)
		hash := getCmd.String("filehash", "", "Hash of file to retrieve")
		from := getCmd.String("from", "", "Source peer address")
		getCmd.Parse(os.Args[2:])

		if *hash == "" || *from == "" {
			fmt.Println("Usage: get --hash <hash> --from <address>")
			os.Exit(1)
		}
		client.Download(*hash, *from)

	default:
		fmt.Println("Unknown command:", os.Args[1])
	}
}