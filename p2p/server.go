package server

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"peerVault/storage"
	"strings"
)

type Server struct {
	port     int
	protocol string
	Listener net.Listener
	Storage  *storage.Storage
}

func NewServer(port int, protocol string, store *storage.Storage) *Server {
	return &Server{
		port:     port,
		protocol: protocol,
		Storage:  store,
	}
}

func (s *Server) ListenAndServe() error {
	listener, err := net.Listen(s.protocol, fmt.Sprintf(":%d", s.port))

	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	defer listener.Close()

	s.Listener = listener
	fmt.Printf("Server is listening on %s:%d...\n", s.protocol, s.port)

	s.StartAcceptConnectionsLoop()

	return nil
}

func (s *Server) StartAcceptConnectionsLoop() {
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) error {
	remoteAddr := conn.RemoteAddr().String()
	fmt.Println("New connection from", remoteAddr)

	defer conn.Close()
	
	reader := bufio.NewReader(conn)

	// Read command
	commandLine, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Failed to read command:", err)
		return err
	}
	command := strings.TrimSpace(commandLine)

	switch  command {
		case "UPLOAD":
			fmt.Println("Handling UPLOAD operation")
			// Read file type
			fileTypeLine, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Failed to read file type:", err)
				return err
			}
			
			fileType := strings.TrimSpace(fileTypeLine)
			fmt.Println("File type:", fileType)

			if err := s.Storage.SaveFile(conn, reader, fileType); err!=nil {
				fmt.Println("Error saving file:", err)
				return err
			}
		case "DOWNLOAD":
			fmt.Println("Handling DOWNLOAD operation")
			s.SendFileToClient(conn, reader)
		default:
			return fmt.Errorf("unknown operation type: %s", command)
	}
	
	return nil
}

	
func (s *Server) SendFileToClient(conn net.Conn, reader *bufio.Reader) error {
	
	hashByte, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println("Error reading hash:", err)	
		return err
	}
	hash := strings.TrimSpace(string(hashByte))

	data, err := s.Storage.ReadFile(hash)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return err
	}

	if _, err = conn.Write(data); err != nil {
		fmt.Printf("Error sending message: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("File sent to client", conn.RemoteAddr().String())

	return nil
}


func (s *Server) Close() error {
	if s.Listener != nil {
		if err := s.Listener.Close(); err != nil {
			return fmt.Errorf("failed to close listener: %w", err)
		}
		fmt.Println("Server closed successfully")
	}
	return nil
}
