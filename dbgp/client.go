package dbgp

import (
	"fmt"
	"strings"
)

// NewClient creates a new client instance
func NewClient(conn *Connection) *Client {
	c := &Client{
		Session:    NewSession(),
		Connection: conn,
	}

	return c
}

// Init reads the init protocol from the xdebug server
func (c *Client) Init() error {
	msg, err := c.Connection.ReadMessage()
	if err != nil {
		return err
	}

	// expecting init
	init, ok := msg.(*ProtocolInit)
	if !ok {
		return fmt.Errorf("Expecting init protocol")
	}

	c.Session.State = StateStarting
	c.Session.CurrentFile = init.FileURI
	c.Session.SetTargetFiles(strings.TrimPrefix(c.Session.CurrentFile, "file://"))
	return nil
}

// GetContext returns the context by id
func (c *Client) GetContext(contextID int) ([]ProtocolProperty, error) {
	cmd := fmt.Sprintf("context_get -i %d -c %d", c.Session.NextTransactionID(), contextID)
	c.Connection.SendMessage(cmd)

	proto, err := c.Connection.GetResponse()
	if err != nil {
		return nil, err
	}

	return proto.PropertyList, nil
}

// GetContextNames returns the context name mapping
func (c *Client) GetContextNames() ([]ProtocolContext, error) {
	cmd := fmt.Sprintf("context_names -i %d -d 0", c.Session.NextTransactionID())
	c.Connection.SendMessage(cmd)

	proto, err := c.Connection.GetResponse()
	if err != nil {
		return nil, err
	}

	return proto.ContextList, nil
}

// GetBreakpointList returns the current list of activated bp's
func (c *Client) GetBreakpointList() ([]ProtocolBreakpoint, error) {
	bpl := []ProtocolBreakpoint{}
	cmd := fmt.Sprintf("breakpoint_list -i %d", c.Session.NextTransactionID())
	c.Connection.SendMessage(cmd)

	proto, err := c.Connection.GetResponse()
	if err != nil {
		return bpl, err
	}

	return proto.BreakpointList, nil
}

// SetExceptionBreakpoint sets a generic break to exceptions
func (c *Client) SetExceptionBreakpoint() error {
	cmd := fmt.Sprintf("breakpoint_set -i %d -t exception -x *", c.Session.NextTransactionID())
	c.Connection.SendMessage(cmd)

	proto, err := c.Connection.GetResponse()
	if err != nil {
		return err
	}

	if proto.HasError() {
		return fmt.Errorf(proto.Error.Message)
	}

	return nil
}

// SetBreakpointToCall creates a call breakpoint
func (c *Client) SetBreakpointToCall(funcName string) error {
	cmd := fmt.Sprintf(
		"breakpoint_set -i %d -t call -m %s",
		c.Session.NextTransactionID(),
		funcName)

	c.Connection.SendMessage(cmd)
	proto, err := c.Connection.GetResponse()
	if err != nil {
		return err
	}

	if proto.HasError() {
		return fmt.Errorf(proto.Error.Message)
	}

	return nil
}

// SetBreakpoint creates a breakpoint
func (c *Client) SetBreakpoint(file string, line int, expr string) error {
	fmtExpr := ""
	if len(expr) > 0 {
		fmtExpr = " -- " + expr
	}

	cmd := fmt.Sprintf("breakpoint_set -i %d -t line -f file://%s -n %d%s", c.Session.NextTransactionID(), file, line, fmtExpr)
	c.Connection.SendMessage(cmd)

	proto, err := c.Connection.GetResponse()
	if err != nil {
		return err
	}

	if proto.HasError() {
		return fmt.Errorf(proto.Error.Message)
	}

	return nil
}

// GetProperty tries fetch a variable by name
func (c *Client) GetProperty(s string) (*ProtocolResponse, error) {
	cmd := fmt.Sprintf("property_get -i %d -n %s", c.Session.NextTransactionID(), s)
	c.Connection.SendMessage(cmd)

	proto, err := c.Connection.GetResponse()
	if err != nil {
		return nil, err
	}

	if proto.HasError() {
		return nil, fmt.Errorf(proto.Error.Message)
	}

	return proto, nil
}

// Step into the next instruction
func (c *Client) Step() (*ProtocolResponse, error) {
	cmd := fmt.Sprintf("step_into -i %d", c.Session.NextTransactionID())
	c.Connection.SendMessage(cmd)

	proto, err := c.Connection.GetResponse()
	if err != nil {
		return nil, err
	}

	if proto.HasError() {
		return nil, fmt.Errorf(proto.Error.Message)
	}

	return proto, nil
}

// Next step over the next instruction
func (c *Client) Next() (*ProtocolResponse, error) {
	cmd := fmt.Sprintf("step_over -i %d", c.Session.NextTransactionID())
	c.Connection.SendMessage(cmd)

	proto, err := c.Connection.GetResponse()
	if err != nil {
		return nil, err
	}

	if proto.HasError() {
		return nil, fmt.Errorf(proto.Error.Message)
	}

	return proto, nil
}

// Finish stops the debugger
func (c *Client) Finish() (*ProtocolResponse, error) {
	cmd := fmt.Sprintf("stop -i %d", c.Session.NextTransactionID())
	c.Connection.SendMessage(cmd)

	proto, err := c.Connection.GetResponse()

	if err != nil {
		return nil, err
	}

	if proto.HasError() {
		return nil, fmt.Errorf(proto.Error.Message)
	}

	return proto, nil
}

// Run the debugger
func (c *Client) Run() (*ProtocolResponse, error) {
	cmd := fmt.Sprintf("run -i %d", c.Session.NextTransactionID())
	c.Connection.SendMessage(cmd)

	proto, err := c.Connection.GetResponse()

	if err != nil {
		return nil, err
	}

	if proto.HasError() {
		return nil, fmt.Errorf(proto.Error.Message)
	}

	return proto, nil
}
