package xdebugcli

import (
	"fmt"
	"io"
	"net"
)

// DBGPServer model
type DBGPServer struct {
	Address  string
	Port     int
	listener net.Listener
}

// Listen start the listening
func (s *DBGPServer) Listen() error {
	var err error

	addr := fmt.Sprintf("%s:%d", s.Address, s.Port)
	if s.listener, err = net.Listen("tcp", addr); err != nil {
		return err
	}

	return nil
}

// Accept connections and start handler
func (s *DBGPServer) Accept() {
	for {
		conn, err := s.listener.Accept()
		fmt.Println("Start session")
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleDGBPConnection(conn)
	}
}

// NewDBGPServer creates a new server
func NewDBGPServer(a string, p int) *DBGPServer {
	s := &DBGPServer{
		Address: a,
		Port:    p,
	}
	return s
}

// handleDGBPConnection
func handleXDebugConnection(c net.Conn) {
	defer c.Close()

	dbgpMessage := ""
	buffer := make([]byte, 512)
	cp := NewCommandProcessor()

	for {
		bytesRead, err := c.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			fmt.Println(err)
			break
		}
		dbgpMessage = dbgpMessage + string(buffer[:bytesRead])

		if bytesRead < len(buffer) || (bytesRead == len(buffer) && buffer[bytesRead-1] == 0) {
			cp.Process(dbgpMessage)
			dbgpMessage = ""
			/*
				reader := bufio.NewReader(os.Stdin)
				text, _ := reader.ReadString('\n')
				text = strings.Trim(text, "\n")
				c.Write([]byte(text))
				c.Write([]byte{0})
			*/
		}
	}

	fmt.Println("Close session")
}
