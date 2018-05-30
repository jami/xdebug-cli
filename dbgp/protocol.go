package dbgp

import (
	"encoding/xml"
	"fmt"
	"strings"

	"golang.org/x/net/html/charset"
)

// HasError getter
func (p ProtocolResponse) HasError() bool {
	return p.Error.Code > 0 || len(p.Error.Message) > 0
}

// CreateProtocolFromXML creator
func CreateProtocolFromXML(xmlString string) (interface{}, error) {
	//fmt.Println("Protocol from xml ", xmlString)
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
