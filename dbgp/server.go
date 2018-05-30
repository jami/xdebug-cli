package dbgp

import (
	"fmt"
	"net"
)

// Listen start the listening
func (s *Server) Listen() error {
	var err error

	addr := fmt.Sprintf("%s:%d", s.Address, s.Port)
	if s.listener, err = net.Listen("tcp", addr); err != nil {
		return err
	}

	return nil
}

// Accept connections and start handler
func (s *Server) Accept(h ConnectionHandler) {
	for {
		conn, err := s.listener.Accept()
		fmt.Println("(xdbg-cli) start session")
		if err != nil {
			fmt.Println(err)
			continue
		}

		go h(NewConnection(conn))
	}
}

// NewServer creates a new server
func NewServer(a string, p uint16) *Server {
	s := &Server{
		Address: a,
		Port:    p,
	}
	return s
}
