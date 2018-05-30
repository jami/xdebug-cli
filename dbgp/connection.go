package dbgp

import (
	"bytes"
	"fmt"
	"io"
	"net"
)

const (
	dbgpBufferSize = 128
)

// GetResponse read a message response
func (c *Connection) GetResponse() (*ProtocolResponse, error) {
	proto, err := c.ReadMessage()
	if err != nil {
		return nil, err
	}

	resp, castable := proto.(*ProtocolResponse)
	if !castable {
		return nil, fmt.Errorf("expecting response protocol")
	}

	return resp, nil
}

// ReadMessage read a message
func (c *Connection) ReadMessage() (interface{}, error) {
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
func (c *Connection) SendMessage(msg string) {
	c.connection.Write([]byte(msg))
	c.connection.Write([]byte{0})
}

// Close the connection
func (c *Connection) Close() {
	c.connection.Close()
}

// NewConnection constructor
func NewConnection(conn net.Conn) *Connection {
	c := &Connection{
		connection:       conn,
		sendHistory:      []string{},
		transactionIndex: 1,
	}

	return c
}
