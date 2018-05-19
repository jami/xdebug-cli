package dbgp

import (
	"encoding/xml"
	"fmt"
	"strings"

	"golang.org/x/net/html/charset"
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

// HasError getter
func (p ProtocolResponse) HasError() bool {
	return p.Error.Code > 0 || len(p.Error.Message) > 0
}

// CreateProtocolFromXML creator
func CreateProtocolFromXML(xmlString string) (interface{}, error) {
	fmt.Println("Protocol from xml ", xmlString)
	decoder := xml.NewDecoder(strings.NewReader(xmlString))
	decoder.CharsetReader = charset.NewReaderLabel

	for {
		t, err := decoder.Token()
		if err != nil {
			return nil, err
		}
		se, ok := t.(xml.StartElement)
		if ok {
			switch name := se.Name.Local; name {
			case "init":
				intiProto := &ProtocolInit{}
				err := decoder.DecodeElement(intiProto, &se)
				if err == nil {
					return intiProto, nil
				}
				return nil, err
			case "response":
				respProto := &ProtocolResponse{}
				err := decoder.DecodeElement(respProto, &se)
				if err == nil {
					return respProto, nil
				}
				return nil, err
			default:
				return nil, fmt.Errorf("Invalid protocol xml")
			}
		}
	}
}
