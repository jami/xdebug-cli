package dbgp

import (
	"bytes"
	"fmt"
	"io"
	"net"
)

type dbgpStateType int

const (
	dbgpBufferSize = 128
)

const (
	dbgpStateStarting dbgpStateType = iota
	dbgpStateStopping
	dbgpStateStopped
	dbgpStateRunning
	dbgpStateBreak
	dbgpStateNone
)

// DBGPMessageType enum
type DBGPMessageType int

const (
	// DBGPMessageStatus status message
	DBGPMessageStatus DBGPMessageType = iota
	// DBGPMessageInit init message
	DBGPMessageInit
)

// DBGPMessage sent by the debugger
type DBGPMessage struct {
	MessageType DBGPMessageType
}

type dbgpProtocolInit struct {
	FileURI  string `xml:"fileuri,attr"`
	Language string `xml:"language,attr"`
	AppID    string `xml:"appid,attr"`
	IDEKey   string `xml:"idekey,attr"`
}

// DBGPConnection model
type DBGPConnection struct {
	connection       net.Conn
	sendHistory      []string
	state            dbgpStateType
	transactionIndex int
}

// ReadMessage read a message
func (c *DBGPConnection) ReadMessage() (interface{}, error) {
	buffer := make([]byte, dbgpBufferSize)
	bufferMessage := []byte{}
	dbgpMessageSize := ""
	dbgpMessageContent := ""

	for {
		_, err := c.connection.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, err
		}

		idx := bytes.IndexByte(buffer, byte(0))
		if idx == -1 {
			bufferMessage = append(bufferMessage, buffer...)
			continue
		}

		if len(dbgpMessageSize) == 0 {
			dbgpMessageSize = string(bufferMessage) + string(buffer[:idx])
			bufferMessage = make([]byte, len(buffer[idx+1:]))
			copy(bufferMessage, buffer[idx+1:])
			continue
		}

		dbgpMessageContent = string(bufferMessage) + string(buffer[:idx])
		break
	}

	return CreateProtocolFromXML(dbgpMessageContent)
}

// SendMessage writes a message
func (c *DBGPConnection) SendMessage(msg string) {
	c.connection.Write([]byte(msg))
	c.connection.Write([]byte{0})
}

// NewDBGPConnection constructor
func NewDBGPConnection(conn net.Conn) *DBGPConnection {
	c := &DBGPConnection{
		connection:       conn,
		sendHistory:      []string{},
		transactionIndex: 1,
		state:            dbgpStateNone,
	}

	return c
}

// DBGPServer model
type DBGPServer struct {
	Address  string
	Port     uint16
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
		fmt.Println("(xdbg-cli) start session")
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleDGBPConnection(conn)
	}
}

// NewDBGPServer creates a new server
func NewDBGPServer(a string, p uint16) *DBGPServer {
	s := &DBGPServer{
		Address: a,
		Port:    p,
	}
	return s
}

// handleDGBPConnection
func handleDGBPConnection(c net.Conn) {
	defer c.Close()

	dbgpConnection := NewDBGPConnection(c)
	msg, err := dbgpConnection.ReadMessage()
	if err != nil {
		fmt.Println("Connection error: ", err)
		return
	}

	// expecting init
	if _, ok := msg.(ProtocolInit); !ok {
		fmt.Println("Expecting init protocol")
		return
	}

	cp := NewCommandProcessor()

	for {
		dbgpConnection.transactionIndex++
		cmd := cp.GetCommand(dbgpConnection.transactionIndex)

		if cmd == "q" || cmd == "quit" {
			fmt.Println("Quitting debugger")
			break
		}

		dbgpConnection.SendMessage(cmd)
		msg, err = dbgpConnection.ReadMessage()

		if err != nil {
			fmt.Println("Error: ", err)
			break
		}
	}
	/*
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
			}
		}
	*/
	fmt.Println("Close connection")
}