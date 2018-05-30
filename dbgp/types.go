package dbgp

import (
	"net"
)

// ProtocolInit data struct
type ProtocolInit struct {
	FileURI  string `xml:"fileuri,attr"`
	Language string `xml:"language,attr"`
	AppID    string `xml:"appid,attr"`
	IDEKey   string `xml:"idekey,attr"`
}

// ProtocolBreakpoint data struct
type ProtocolBreakpoint struct {
	Type     string `xml:"type,attr"`
	FileName string `xml:"filename,attr"`
	Line     int    `xml:"lineno,attr"`
	State    string `xml:"state,attr"`
	HitCount int    `xml:"hit_count,attr"`
}

// ProtocolContext data struct
type ProtocolContext struct {
	ID   string `xml:"id,attr"`
	Name string `xml:"name,attr"`
}

// ProtocolStack data struct
type ProtocolStack struct {
	Where    string `xml:"where,attr"`
	Level    int    `xml:"level,attr"`
	Type     string `xml:"type,attr"`
	Filename string `xml:"filename,attr"`
	Line     int    `xml:"lineno,attr"`
}

// ProtocolProperty data struct
type ProtocolProperty struct {
	Name        string            `xml:"name,attr"`
	Fullname    string            `xml:"fullname,attr"`
	Type        string            `xml:"type,attr"`
	Children    int               `xml:"children,attr"`
	NumChildren int               `xml:"numchildren,attr"`
	Page        int               `xml:"page,attr"`
	PageSize    int               `xml:"pagesize,attr"`
	Content     string            `xml:",innerxml"`
	Property    *ProtocolProperty `xml:"property"`
}

// ProtocolMessage data struct
type ProtocolMessage struct {
	Filename  string `xml:"filename,attr"`
	Line      int    `xml:"lineno,attr"`
	Exception string `xml:"exception,attr"`
	Value     string `xml:",chardata"`
}

// ProtocolError data struct
type ProtocolError struct {
	Code    int    `xml:"code,attr"`
	Message string `xml:"message"`
}

// ProtocolResponse data struct
type ProtocolResponse struct {
	Command        string               `xml:"command,attr"`
	Context        string               `xml:"context,attr"`
	TransactionID  string               `xml:"transaction_id,attr"`
	Reason         string               `xml:"reason,attr"`
	Status         string               `xml:"status,attr"`
	Error          ProtocolError        `xml:"error"`
	BreakpointList []ProtocolBreakpoint `xml:"breakpoint"`
	ContextList    []ProtocolContext    `xml:"context"`
	PropertyList   []ProtocolProperty   `xml:"property"`
	StackList      []ProtocolStack      `xml:"stack"`
	Message        ProtocolMessage      `xml:"message"`
}

// Connection model
type Connection struct {
	connection       net.Conn
	sendHistory      []string
	transactionIndex int
}

// ConnectionHandler handles srv accepts
type ConnectionHandler func(*Connection)

// SessionStateType enumeration type
type SessionStateType int

const (
	// StateStarting initial state
	StateStarting SessionStateType = iota
	// StateStopping remote debugger is trying to stop the process
	StateStopping
	// StateStopped remote debugger has stopped
	StateStopped
	// StateRunning interpreter is running until next breakpoint or EOF
	StateRunning
	// StateBreak reached breakpoint
	StateBreak
	// StateNone undefined state
	StateNone
)

// Session debug session
type Session struct {
	TransactionID int
	State         SessionStateType
	CurrentFile   string
	CurrentLine   int
	TargetFiles   []string
}

// Server model
type Server struct {
	Address  string
	Port     uint16
	listener net.Listener
}

// Client for xdbg srv
type Client struct {
	Connection *Connection
	Session    *Session
}
